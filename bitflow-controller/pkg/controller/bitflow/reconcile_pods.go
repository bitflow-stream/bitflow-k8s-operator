package bitflow

import (
	"context"
	"fmt"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/scheduler"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type PodCreation struct {
	name     string
	step     *bitflowv1.BitflowStep
	sources  []*bitflowv1.BitflowSource
	oneToOne bool
}

func (c *PodCreation) Log(node *corev1.Node) *log.Entry {
	entry := c.step.Log().WithFields(log.Fields{
		"pod":  c.name,
		"node": node.Name,
	})
	if c.sources != nil {
		if c.oneToOne {
			entry = c.sources[0].LogFields(entry)
		} else {
			entry = entry.WithField("sources", len(c.sources))
		}
	}
	return entry
}

func (r *BitflowReconciler) validateStep(step *bitflowv1.BitflowStep) bool {
	validationMsg := step.Status.ValidationError

	step.Validate()

	if validationMsg != step.Status.ValidationError {
		err := r.client.Status().Update(context.TODO(), step)
		if err != nil {
			step.Log().Errorf("Failed to update validation error status: %v", err)
		}
	}

	if step.Status.ValidationError != "" {
		r.cleanupStep(step.Name, fmt.Sprintf("Step validation error: %v", step.Status.ValidationError))
		return false
	}
	return true
}

func (r *BitflowReconciler) GetAllToOnePod(step *bitflowv1.BitflowStep, matchedSources []*bitflowv1.BitflowSource) ([]*corev1.Pod, error) {
	validSources := 0
	for _, source := range matchedSources {
		if source.Status.ValidationError == "" {
			validSources += 1
		}
	}
	if validSources == 0 {
		r.cleanupStep(step.Name, "Step has no valid sources")
		return nil, nil
	}
	return r.GetSingletonPod(step, matchedSources)
}

func (r *BitflowReconciler) GetSingletonPod(step *bitflowv1.BitflowStep, matchedSources []*bitflowv1.BitflowSource) ([]*corev1.Pod, error) {
	var podList []*corev1.Pod

	name := ConstructSingletonPodName(step.Name)
	found := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: r.namespace}, found)
	if err != nil {
		found = r.handleMissingSingletonPod(err, step, name, matchedSources)
	}

	// TODO instead of simply comparing the name and some labels, compare the entire Pod struct with the desired pod (fully resolved template)
	// Same for one-to-one pods and output sources.

	if found != nil && !CompareSingletonSpec(found, step) {
		r.deletePod(found, step.Log(), "does not match step spec")
	} else if found != nil && !common.IsBeingDeleted(found) {
		podList = []*corev1.Pod{found}
	}
	r.cleanupPodsForStep(step.Name, step.Log(), "dangling", podList...)
	return podList, nil
}

func CompareSingletonSpec(pod *corev1.Pod, step *bitflowv1.BitflowStep) bool {
	return step.Type() == pod.Labels[bitflowv1.LabelStepType] &&
		step.Name == pod.Labels[bitflowv1.LabelStepName]
}

func (r *BitflowReconciler) handleMissingSingletonPod(err error, step *bitflowv1.BitflowStep, podName string, matchedSources []*bitflowv1.BitflowSource) *corev1.Pod {
	if errors.IsNotFound(err) {
		return r.createPod(&PodCreation{
			name:     podName,
			step:     step,
			oneToOne: false,
			sources:  matchedSources,
		})
	}
	step.Log().WithField("pod", podName).Errorf("Failed to query pod: %v", err)
	return nil
}

func (r *BitflowReconciler) GetOneToOnePods(step *bitflowv1.BitflowStep, matchedSources []*bitflowv1.BitflowSource) ([]*corev1.Pod, error) {
	podList := make([]*corev1.Pod, 0, len(matchedSources))
	validRestartingPods := make([]string, 0, len(matchedSources))

	for _, source := range matchedSources {
		podName := ConstructReproduciblePodName(step.Name, source.Name)
		found := &corev1.Pod{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: r.namespace}, found)
		if err != nil {
			found = r.handleMissingPod(err, step, source, podName)
		}
		if found != nil {
			if !CompareOneToOneSpec(source, found, step) {
				r.deletePod(found, step.Log(), "does not match step spec")
			} else if source.Status.ValidationError != "" {
				r.deletePod(found, step.Log(), fmt.Sprintf("source '%v' has validation error: %v", source.Name, source.Status.ValidationError))
			} else if !common.IsBeingDeleted(found) {
				podList = append(podList, found)
			} else if common.IsBeingDeleted(found) {
				if _, ok := r.respawning.IsPodRestarting(found.Name); ok {
					validRestartingPods = append(validRestartingPods, found.Name)
				}
			}
		}
	}
	r.respawning.DeletePodsWithLabelExcept(bitflowv1.LabelStepName, step.Name, validRestartingPods)
	r.cleanupPodsForStep(step.Name, step.Log(), "dangling", podList...)
	return podList, nil
}

func CompareOneToOneSpec(source *bitflowv1.BitflowSource, pod *corev1.Pod, step *bitflowv1.BitflowStep) bool {
	return source.Name == pod.Labels[bitflowv1.PodLabelOneToOneSourceName] &&
		step.Type() == pod.Labels[bitflowv1.LabelStepType] &&
		step.Name == pod.Labels[bitflowv1.LabelStepName]
}

func (r *BitflowReconciler) handleMissingPod(err error, step *bitflowv1.BitflowStep, source *bitflowv1.BitflowSource, podName string) *corev1.Pod {
	if errors.IsNotFound(err) && source.Status.ValidationError == "" {
		return r.createPod(&PodCreation{
			name:     podName,
			step:     step,
			oneToOne: true,
			sources:  []*bitflowv1.BitflowSource{source},
		})
	} else if source.Status.ValidationError != "" {
		source.Log().Infoln("Source has validation error and will not be processed:", source.Status.ValidationError)
	}
	step.Log().WithField("pod", podName).Errorf("Failed to query pod: %v", err)
	return nil
}

func (r *BitflowReconciler) createPod(model *PodCreation) *corev1.Pod {
	if spec, ok := r.respawning.IsPodRestarting(model.name); ok {
		return r.createRestartedPod(model.name, spec)
	}
	return r.createNewPod(model)
}

func (r *BitflowReconciler) createRestartedPod(podName string, status common.RespawningStatus) *corev1.Pod {
	pod := status.Pod
	r.resourceLimiter.AssignResourcesNodeImplicit(pod)
	log.WithField("pod", pod.Name).Info("Restarting pod")
	err := r.client.Create(context.TODO(), pod)
	if err != nil {
		log.WithField("pod", pod.Name).Errorf("Error creating restarted pod: %v", err)
		return nil
	}
	r.waitForCacheSync()
	r.respawning.Delete(podName) // TODO define max retries, if errors occur
	r.statistic.PodRespawned()
	return pod
}

func (r *BitflowReconciler) createNewPod(model *PodCreation) *corev1.Pod {
	pod := model.step.Spec.Template.DeepCopy()
	pod.Name = model.name
	pod.Namespace = r.namespace

	extraLabels := r.idLabels
	extraEnv := r.config.GetPodEnvVars()
	if model.oneToOne {
		// TODO rethink owner references
		// setOwnerReferenceForPod(pod, &model.sources.Items[0])
		PatchOneToOnePod(pod, model.step, model.sources[0], extraLabels, extraEnv, r.ownPodIP, r.apiPort)
	} else {
		PatchSingletonPod(pod, model.step, extraLabels, extraEnv, r.ownPodIP, r.apiPort)
	}
	// TODO set owner references for all managed pods and sources
	// pod.OwnerReferences = r.makeOwnerReferences()

	assignedNode, _ := r.scheduler.SchedulePod(pod, model.step, model.sources)
	if assignedNode != nil {
		scheduler.SetPodNodeAffinityRequired(assignedNode, pod)
		r.resourceLimiter.AssignResources(pod, assignedNode)
	}
	logger := model.Log(assignedNode)
	logger.Infof("Spawning pod")
	err := r.client.Create(context.TODO(), pod)
	if err != nil {
		logger.Errorf("Error creating pod: %v", err)
		return nil
	}
	r.waitForCacheSync()
	return pod
}

// TODO rethink what owner references to set
// This method is duplicated in tags-change-notify-http.go
func (r *BitflowReconciler) makeOwnerReferences(pod *corev1.Pod, isController, blockOwnerDeletion bool) []metav1.OwnerReference {
	if pod == nil || pod.Name == "" || pod.UID == "" {
		return nil
	}
	// For some reason these meta fields are not set when querying the Pod
	apiVersion := pod.APIVersion
	if apiVersion == "" {
		apiVersion = "v1"
	}
	kind := pod.Kind
	if kind == "" {
		kind = "Pod"
	}
	return []metav1.OwnerReference{{
		APIVersion:         apiVersion,
		Kind:               kind,
		Name:               pod.Name,
		UID:                pod.UID,
		Controller:         &isController,
		BlockOwnerDeletion: &blockOwnerDeletion,
	}}
}

// TODO rethink what owner references to set
func setOwnerReferenceForPod(pod *corev1.Pod, obj runtime.Object) {
	var apiVersion string
	var kind string
	var name string
	var uid types.UID
	var controller bool
	var blockDel bool
	matched := false
	if source, ok := obj.(*bitflowv1.BitflowSource); ok && source.UID != "" {
		// For some reason these meta fields are not set when querying the Pod
		apiVersion = source.APIVersion
		if apiVersion == "" {
			apiVersion = bitflowv1.GroupVersion
		}
		kind = source.Kind
		if kind == "" {
			kind = bitflowv1.DataSourcesKind
		}
		name = source.Name
		uid = source.UID
		controller = true
		blockDel = true
		matched = true
	} else if step, ok := obj.(*bitflowv1.BitflowStep); ok && step.UID != "" {
		apiVersion = step.APIVersion
		if apiVersion == "" {
			apiVersion = bitflowv1.GroupVersion
		}
		kind = step.Kind
		if kind == "" {
			kind = bitflowv1.StepsKind
		}
		name = step.Name
		uid = step.UID
		controller = false
		blockDel = false
		matched = true
	}

	if matched {
		ref := []metav1.OwnerReference{{
			APIVersion:         apiVersion,
			Kind:               kind,
			Name:               name,
			UID:                uid,
			Controller:         &controller,
			BlockOwnerDeletion: &blockDel,
		}}
		for _, oRef := range pod.OwnerReferences {
			if *oRef.Controller && controller {
				log.Infoln("There already is an Ownerreference flags as controller, we will ignore it", "Pod", pod, "Owner", oRef)
			} else {
				ref = append(ref, oRef)
			}
		}
		pod.OwnerReferences = ref
	}
}

func (r *BitflowReconciler) waitForCacheSync() {
	log.Debugln("Started cache sync")
	stopper := make(chan struct{})
	r.cache.WaitForCacheSync(stopper)
	close(stopper)
	log.Debugln("Finished cache sync")
}
