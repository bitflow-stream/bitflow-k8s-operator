package bitflow

import (
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func (r *BitflowReconciler) constructAllToOnePod(step *bitflowv1.BitflowStep, matchedSources []*bitflowv1.BitflowSource) map[*corev1.Pod][]*bitflowv1.BitflowSource {
	validSources := 0
	for _, source := range matchedSources {
		if source.Status.ValidationError == "" {
			validSources += 1
		}
	}
	if validSources > 0 {
		return r.constructSingletonPod(step, matchedSources)
	}
	return nil
}

func (r *BitflowReconciler) constructSingletonPod(step *bitflowv1.BitflowStep, matchedSources []*bitflowv1.BitflowSource) map[*corev1.Pod][]*bitflowv1.BitflowSource {
	name := ConstructSingletonPodName(step.Name)
	pod := r.constructPod(name, false, step, matchedSources)
	return map[*corev1.Pod][]*bitflowv1.BitflowSource{
		pod: matchedSources,
	}
}

func (r *BitflowReconciler) constructOneToOnePods(step *bitflowv1.BitflowStep, matchedSources []*bitflowv1.BitflowSource) map[*corev1.Pod][]*bitflowv1.BitflowSource {
	pods := make(map[*corev1.Pod][]*bitflowv1.BitflowSource, len(matchedSources))
	for _, source := range matchedSources {
		name := ConstructReproduciblePodName(step.Name, source.Name)
		pod := r.constructPod(name, true, step, []*bitflowv1.BitflowSource{source})
		pods[pod] = []*bitflowv1.BitflowSource{source}
	}
	return pods
}

func (r *BitflowReconciler) constructPod(name string, oneToOne bool, step *bitflowv1.BitflowStep, sources []*bitflowv1.BitflowSource) *corev1.Pod {
	pod := step.Spec.Template.DeepCopy()
	pod.Name = name
	pod.Namespace = r.namespace

	extraLabels := r.idLabels
	extraEnv := r.config.GetPodEnvVars()
	if oneToOne {
		// TODO rethink owner references
		// setOwnerReferenceForPod(pod, &model.sources.Items[0])
		PatchOneToOnePod(pod, step, sources[0], extraLabels, extraEnv, r.ownPodIP, r.apiPort)
	} else {
		PatchSingletonPod(pod, step, extraLabels, extraEnv, r.ownPodIP, r.apiPort)
	}
	// TODO set owner references for all managed pods and sources
	// pod.OwnerReferences = r.makeOwnerReferences()

	return pod
}

// TODO rethink what owner references to set
// This method is duplicated in tags-change-notify-http.go
func (r *BitflowReconciler) makeOwnerReferences(pod *corev1.Pod, isController, blockOwnerDeletion bool) []metav1.OwnerReference {
	if pod == nil || pod.Name == "" || pod.UID == "" {
		return nil
	}
	// For some reason these meta fields are not set when querying the pod
	apiVersion := pod.APIVersion
	if apiVersion == "" {
		apiVersion = "v1"
	}
	kind := pod.Kind
	if kind == "" {
		kind = "pod"
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
		// For some reason these meta fields are not set when querying the pod
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
				log.Infoln("There already is an Ownerreference flags as controller, we will ignore it", "pod", pod, "Owner", oRef)
			} else {
				ref = append(ref, oRef)
			}
		}
		pod.OwnerReferences = ref
	}
}
