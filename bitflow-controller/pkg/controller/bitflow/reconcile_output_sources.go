package bitflow

import (
	"context"
	"strconv"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type PodOutputPair struct {
	pod    *corev1.Pod
	output *bitflowv1.StepOutput
}

func (r *BitflowReconciler) validateSource(source *bitflowv1.BitflowSource) {
	validationMsg := source.Status.ValidationError
	source.Validate()
	if validationMsg != source.Status.ValidationError {
		err := r.client.Status().Update(context.TODO(), source)
		if err != nil {
			source.Log().Errorf("Failed to update validation error status: %v", err)
		}
	}
}

func (r *BitflowReconciler) reconcileOutputSources() {
	// Instead of working with actually running pods, check the pods that we WANT/PLAN to be running.
	requiredSources := make(map[string]*bitflowv1.BitflowSource)
	r.pods.Read(func() {
		for _, pod := range r.pods.pods {
			for _, out := range pod.step.Spec.Outputs {
				requiredSource := CreateOutputSource(pod.step, pod.pod, out, pod.inputSources, r.idLabels)
				if requiredSource != nil {
					requiredSources[requiredSource.Name] = requiredSource
				}
			}
		}
	})

	// Query all managed data sources
	allOutputSources, err := r.listOutputSources("")
	if err != nil {
		log.Error("Failed to query existing managed output sources:", err)
		return
	}

	// Check existing source and delete those that are not necessary. Try to update instead of re-creating.
	for _, existing := range allOutputSources {
		if required, ok := requiredSources[existing.Name]; ok {
			if CompareSources(required, existing) {
				delete(requiredSources, existing.Name)
			} else if CanSourceBeUpdated(required, existing) {
				delete(requiredSources, existing.Name)
				// TODO update
			} else {
				// TODO delete (and remain in requiredSources map)
			}
		} else {
			// TODO delete
		}
	}

	// Create missing sources
	for _, source := range requiredSources {
		// TODO create
	}
}

func (r *BitflowReconciler) xxx() {
	// TODO properly clean up left over sources...

	if len(step.Spec.Outputs) == 0 {
		return nil
	}

	// TODO: handle source validation errors!!!!!!!!!!!!!!!!!!!
	podsWithoutSource := make([]PodOutputPair, 0, len(podList))
	podsWaitingForIP := make(map[string]bool)

	matchedSources := make(map[string]*bitflowv1.BitflowSource)
	var found bool
	for _, pod := range podList {
		if pod.Status.PodIP == "" {
			podsWaitingForIP[pod.Name] = true
			step.Log().WithField("pod", pod.Name).Debugf("pod is missing an IP, not creating output sources...")
			continue
		}
		if pod.DeletionTimestamp != nil {
			step.Log().WithField("pod", pod.Name).Debugf("pod is scheduled for deletion, not creating output sources...")
			continue
		}
		for _, out := range step.Spec.Outputs {
			found = false
			// Create a temporary source object the same way as an actual source object would be created in createSource().
			// Compare all properties of this temporary object with existing sources.
			requiredOut := r.makeSourceObject(step, pod, out, inputSources)
			if requiredOut == nil {
				continue
			}
			for _, existOut := range allOutputSources {

				// TODO IMPORTANT: when a pod is being restarted, it will get a new IP. In this case the source should be updated, instead of being deleted and recreated.
				// Avoid pods all subsequent pipeline pods!

				if CompareSources(requiredOut, existOut) {
					found = true
					matchedSources[existOut.Name] = existOut
				}
			}
			if !found {
				podsWithoutSource = append(podsWithoutSource, PodOutputPair{pod, out})
			}
		}
	}

	for _, source := range allOutputSources {
		if _, ok := matchedSources[source.Name]; !ok {
			podName := source.Labels[bitflowv1.SourceLabelPodName]
			logger := step.LogFields(source.Log())
			if podName == "" {
				logger.Warnf("Source does not have valid '%v' label, ignoring...", bitflowv1.SourceLabelPodName)
				continue
			}
			logger = logger.WithField("pod", podName)
			if _, present := r.pods.IsPodRespawning(podName); podName != "" && present {
				// pod is currently restarting -> dont kill its output source yet
				logger.Debug("Missing pod is being restarted, not deleting output source")
				continue
			}
			if podsWaitingForIP[podName] {
				continue
			}

			logger.Info("Deleting dangling output source")
			err := r.client.Delete(context.TODO(), source)
			if err != nil && !errors.IsNotFound(err) {
				logger.Errorln("Error deleting dangling output source:", err)
			}
		}
	}

	for _, podSource := range podsWithoutSource {
		r.createSource(step, &podSource, inputSources)
	}
	return nil
}

func CompareSources(required *bitflowv1.BitflowSource, existing *bitflowv1.BitflowSource) bool {
	if required.Name != existing.Name || required.Spec.URL != existing.Spec.URL || len(required.Labels) != len(existing.Labels) {
		return false
	}
	for i, label := range required.Labels {
		if existing.Labels[i] != label {
			return false
		}
	}
	return true
}

func (r *BitflowReconciler) createSource(step *bitflowv1.BitflowStep, podSource *PodOutputPair, matchedInputSources []*bitflowv1.BitflowSource) {

	// TODO instead of re-creating the patched source here, re-use the source object created earlier in reconcileOutputSources
	source := r.makeSourceObject(step, podSource.pod, podSource.output, matchedInputSources)
	if source == nil {
		return
	}
	logger := step.LogFields(source.Log()).WithField("pod", podSource.pod.Name)

	// TODO set the controller pod AND the step as owner references?? Set the analysis pod as owner??

	if err := controllerutil.SetControllerReference(step, source, r.scheme); err != nil {
		logger.Errorln("Error setting controller ref on source:", err)
		return
	}

	logger.Info("Creating new output source")
	err := r.client.Create(context.TODO(), source)
	if err != nil {
		logger.Errorln("Error creating output source:", err)
	}
}

func CreateOutputSource(step *bitflowv1.BitflowStep, pod *corev1.Pod, out *bitflowv1.StepOutput, matchedInputSources []*bitflowv1.BitflowSource, extraLabels map[string]string) *bitflowv1.BitflowSource {
	name := ConstructSourceName(pod.Name, out.Name)
	url := out.GetOutputSourceURL(pod)
	if url == "" {
		logger := step.Log().WithField("pod", pod.Name).WithField("output", out.Name)
		logger.Errorln("Cannot create patched output source: StepOutput has not been validated (URL not parsed)")
		return nil
	}

	labels := MergeLabels(matchedInputSources, step.Name)
	for key, val := range out.Labels {
		labels[key] = val
	}
	for key, val := range extraLabels {
		labels[key] = val
	}
	labels[bitflowv1.LabelStepName] = step.Name
	labels[bitflowv1.LabelStepType] = step.Type()
	labels[bitflowv1.SourceLabelPodName] = pod.Name
	labels[bitflowv1.SourceLabelPodOutputName] = out.Name

	return &bitflowv1.BitflowSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: pod.Namespace,
			Labels:    labels,
		},
		Spec: bitflowv1.BitflowSourceSpec{
			URL: url,
		},
	}
}

func MergeLabels(sources []*bitflowv1.BitflowSource, newStepName string) map[string]string {
	maxDepth := 0
	mergedLabels := make(map[string]string)
	for _, source := range sources {
		if source.Status.ValidationError != "" {
			continue
		}

		if depthStr, ok := source.Labels[bitflowv1.PipelineDepthLabel]; ok {
			newDepth, err := strconv.Atoi(depthStr)
			if err != nil {
				source.Log().Errorf("Failed to parse label value %v=%v as integer: %v. Assuming pipeline depth zero.", bitflowv1.PipelineDepthLabel, depthStr, err)
			} else if newDepth > maxDepth {
				maxDepth = newDepth
			}
		}

		if len(mergedLabels) == 0 {
			// For the first source, just copy all labels
			for key, val := range source.Labels {
				mergedLabels[key] = val
			}
		} else {
			// For all remaining sources, only keep the labels that are shared by all
			for key, val := range source.Labels {
				if oldVal, ok := mergedLabels[key]; ok && oldVal != val {
					delete(mergedLabels, key)
				}
			}
		}
	}

	// Put/Overwrite special labels
	maxDepth++
	depthStr := strconv.Itoa(maxDepth)
	mergedLabels[bitflowv1.PipelineDepthLabel] = depthStr
	mergedLabels[bitflowv1.PipelinePathLabelPrefix+depthStr] = newStepName

	return mergedLabels
}

func (r *BitflowReconciler) cleanupOutputSourcesForStep(stepName string, logger *log.Entry) {
	outputs, err := r.listOutputSources(stepName)
	if err != nil {
		logger.Errorln("Failed to query output sources for step:", err)
	} else {
		for _, source := range outputs {
			err = r.client.Delete(context.TODO(), source)
			if err != nil {
				source.Log().Errorln("Failed to delete source:", err)
			}
		}
	}
}
