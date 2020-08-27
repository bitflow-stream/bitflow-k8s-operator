package bitflow

import (
	"context"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

func (r *BitflowReconciler) spawnPods() {
	// TODO schedule and limit resources
	assignedNode, _ := r.scheduler.SchedulePod(pod, model.step, model.sources)
	if assignedNode != nil {
		common.SetTargetNode(assignedNode, pod)
		r.resourceLimiter.AssignResources(pod, assignedNode)

		// Bei restarting:
		r.resourceLimiter.AssignResourcesNodeImplicit(pod)
	}

	for stepName, pods := range r.pods.GetSteps() {
		r.spawnPodsForStep(stepName, pods)
	}

	// Make sure the created pods are visible on the next Reconcile loop
	r.waitForCacheSync()
}

func (r *BitflowReconciler) spawnPodsForStep(stepName string, pods map[string]*PodStatus) {
	logger := log.WithField("step", stepName)

	// First check running pods for correctness
	existingPods, err := r.listPodsForStep(stepName)
	if err != nil {
		logger.Errorln("Failed to list pods:", err)
		return
	}
	for _, existingPod := range existingPods.Items {
		if plannedPod, ok := pods[existingPod.Name]; ok {
			if !r.isPodCorrect(plannedPod, &existingPod) {
				// Pod does not match expectations - delete it and mark it as respawning. Later, it will be re-created correctly.
				r.deletePod(&existingPod, logger, "respawn")
				r.pods.MarkRespawning(&existingPod, true)
			}

			// Make sure this pod is not re-created again below
			delete(pods, existingPod.Name)
		} else {
			// Pod is dangling - delete it
			r.deletePod(&existingPod, logger, "dangling")
		}
	}

	// Then create missing pods
	for _, pod := range pods {
		r.spawnPod(pod)
	}
}

func (r *BitflowReconciler) spawnPod(pod *PodStatus) {
	logger := pod.Log()
	if pod.respawning {
		logger.Info("Spawning pod")
	} else {
		logger.Info("Respawning pod")
	}
	err := r.client.Create(context.TODO(), pod.pod)
	if err != nil {
		logger.Errorf("Error spawning pod: %v", err)
	} else {
		r.pods.MarkRespawning(pod.pod, false)
	}
}

func (r *BitflowReconciler) isPodCorrect(pod *PodStatus, running *corev1.Pod) bool {
	// TODO Check: node, resources, spec
}

func CompareSingletonSpec(pod *corev1.Pod, step *bitflowv1.BitflowStep) bool {
	return step.Type() == pod.Labels[bitflowv1.LabelStepType] &&
		step.Name == pod.Labels[bitflowv1.LabelStepName]
}

func CompareOneToOneSpec(source *bitflowv1.BitflowSource, pod *corev1.Pod, step *bitflowv1.BitflowStep) bool {
	return source.Name == pod.Labels[bitflowv1.PodLabelOneToOneSourceName] &&
		step.Type() == pod.Labels[bitflowv1.LabelStepType] &&
		step.Name == pod.Labels[bitflowv1.LabelStepName]
}
