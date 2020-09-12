package bitflow

import (
	"context"
	"strconv"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	// TODO what about restarting pods?

	type requiredSource struct {
		source       *bitflowv1.BitflowSource
		pod          *PodStatus
		waitingForIp bool
		logger       *log.Entry
	}

	// Instead of working with actually running pods, check the pods that we WANT/PLAN to be running.
	requiredSources := make(map[string]*requiredSource)
	r.pods.Read(func() {
		for _, pod := range r.pods.pods {
			for _, out := range pod.step.Spec.Outputs {
				sourceObject := createOutputSource(pod.step, pod.pod, out, pod.inputSources, r.idLabels)
				if sourceObject != nil {
					requiredSources[sourceObject.Name] = &requiredSource{
						source:       sourceObject,
						pod:          pod,
						waitingForIp: pod.pod.Status.PodIP == "",
						logger:       sourceObject.Log().WithFields(log.Fields{"step": pod.step.Name, "pod": pod.pod.Name, "output": out.Name}),
					}
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

	// Check existing sources and delete those that are not necessary. Try to update instead of re-creating, if possible.
	for _, existing := range allOutputSources {
		deleteSource := false
		if required, ok := requiredSources[existing.Name]; ok {
			if r.compareSources(required.source, existing) == "" || required.waitingForIp {
				delete(requiredSources, existing.Name)
			} else if r.canSourceBeUpdated(required.source, existing) {
				// Update source without deleting it
				delete(requiredSources, existing.Name)
				required.logger.Info("Updating output source")
				if err := r.client.Update(context.TODO(), existing); err != nil {
					required.logger.Errorln("Failed to update output source:", err)
				}
			} else {
				// Updating the source is not possible. Delete it and recreate it afterwards.
				deleteSource = true
				required.logger.Info("Re-creating output source")
			}
		} else {
			log.WithFields(log.Fields{"source": existing.Name}).Info("Deleting dangling output source")
			deleteSource = true
		}
		if deleteSource {
			r.deleteObject(existing, "Failed to delete output source")
		}
	}

	// Create missing sources
	for _, source := range requiredSources {
		if source.waitingForIp {
			continue
		}

		source.logger.Info("Creating output source")
		if err := r.client.Create(context.TODO(), source.source); err != nil {
			source.logger.Errorln("Error creating output source:", err)
		}
	}
}

func (r *BitflowReconciler) compareSources(source1, source2 *bitflowv1.BitflowSource) string {
	return r.compareObjects(source1.TypeMeta, source2.TypeMeta, source1.ObjectMeta, source2.ObjectMeta, source1.Spec, source2.Spec)
}

func (r *BitflowReconciler) canSourceBeUpdated(required, existing *bitflowv1.BitflowSource) bool {
	// Simple update is only possible, if the .Spec field is the only changed field. The TypeMeta, ObjectMeta and Status require deletion of the object.
	// TODO check if Labels can be updated without deletion.

	return r.compareMetaData(required.TypeMeta, existing.TypeMeta, required.ObjectMeta, existing.ObjectMeta) == ""
}

func createOutputSource(step *bitflowv1.BitflowStep, pod *corev1.Pod, out *bitflowv1.StepOutput, matchedInputSources []*bitflowv1.BitflowSource, extraLabels map[string]string) *bitflowv1.BitflowSource {
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

	source := &bitflowv1.BitflowSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: pod.Namespace,
			Labels:    labels,
		},
		Spec: bitflowv1.BitflowSourceSpec{
			URL: url,
		},
	}

	// TODO set the controller pod AND the step as owner references?? Set the analysis pod as owner??
	// if err := controllerutil.SetControllerReference(step, source, r.scheme); err != nil {
	//		logger.Errorln("Error setting controller ref on source:", err)
	// return
	// }

	return source
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
