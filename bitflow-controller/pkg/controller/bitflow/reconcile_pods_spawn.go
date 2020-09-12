package bitflow

import (
	"context"
	"fmt"
	"strconv"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *BitflowReconciler) ensurePods() {
	nodeList, err := common.RequestReadyNodes(r.client)
	if err != nil {
		log.Errorln("Failed to query ready nodes, cannot spawn/delete pods:", err)
		return
	}
	nodes := make(map[string]*corev1.Node, len(nodeList.Items))
	for i, node := range nodeList.Items {
		nodes[node.Name] = &nodeList.Items[i]
	}

	// First assign nodes and resources to the pods. These two steps do not create pods, but only edit the ManagedPods field.
	r.schedulePods(nodes)
	r.assignPodResources(nodes)

	// Map steps to their associated pods.
	// Also, prepare the pods by patching their target node and resources (computed above)
	stepsAndPods := make(map[string]map[string]*PodStatus)
	r.pods.Modify(func() {
		for _, pod := range r.pods.pods {
			// Copy the target node and resource list into the Pod spec
			// DeepCopy-ing the pod should not be necessary
			common.SetTargetNode(pod.pod, pod.targetNode)
			patchPodResourceLimits(pod.pod, pod.resources)

			pods, ok := stepsAndPods[pod.step.Name]
			if !ok {
				pods = make(map[string]*PodStatus)
				stepsAndPods[pod.step.Name] = pods
			}
			pods[pod.pod.Name] = pod
		}
	})

	// Now all pods are prepared - make sure the existing pods match the planned pods.
	expectedSteps := make([]string, 0, len(stepsAndPods))
	for stepName, pods := range stepsAndPods {
		log.WithFields(log.Fields{"step": stepName, "num-pods": strconv.Itoa(len(pods))}).Debugf("Ensuring pods")
		r.ensurePodsForStep(stepName, pods)
		expectedSteps = append(expectedSteps, stepName)
	}

	// Delete Pods that do not belong to any managed step - e.g. when a step was entirely deleted
	r.cleanupDanglingPods(expectedSteps)

	// Make sure the created pods are visible on the next Reconcile loop
	r.waitForCacheSync()
}

func (r *BitflowReconciler) ensurePodsForStep(stepName string, pods map[string]*PodStatus) {
	logger := log.WithField("step", stepName)

	// First check running pods for correctness
	existingPods, err := r.listPodsForStep(stepName)
	if err != nil {
		logger.Errorln("Failed to list pods:", err)
		return
	}
	for i := range existingPods.Items {
		existingPod := &existingPods.Items[i]
		if plannedPod, ok := pods[existingPod.Name]; ok {
			if diff := r.comparePods(plannedPod.pod, existingPod); diff == "" {
				// The pod is correct, but the status might have been updated by the system.
				// Especially, the Pod might have been assigned an IP.
				r.pods.UpdateExistingPod(existingPod)
			} else if !plannedPod.respawning && !common.IsBeingDeleted(existingPod) {
				// Pod does not match expectations - delete it and mark it as respawning. Later, it will be re-created correctly.
				if r.deletePod(existingPod, logger, "respawn, "+diff) {
					r.pods.MarkRespawning(existingPod, true)
				}
				log.Debugf("Existing pod: %v", existingPod)
				log.Debugf("Planned  pod: %v", plannedPod.pod)
			}

			// Make sure this pod is not re-created again below
			delete(pods, existingPod.Name)
		} else {
			// Pod is dangling - delete it
			r.deletePod(existingPod, logger, "dangling")
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
		logger.Info("Respawning pod")
	} else {
		logger.Info("Spawning pod")
	}
	err := r.client.Create(context.TODO(), pod.pod)
	if err != nil {
		logger.Errorf("Error spawning pod: %v", err)
	} else {
		r.pods.MarkRespawning(pod.pod, false)
	}
}

func (r *BitflowReconciler) cleanupDanglingPods(expectedSteps []string) {
	// Delete all pods that do not belong to any of the listed steps

	labelSelector := make(labels.Set)
	for k, v := range r.idLabels {
		labelSelector[k] = v
	}
	selector := labels.SelectorFromSet(labelSelector)

	if len(expectedSteps) > 0 {
		req, err := labels.NewRequirement(bitflowv1.LabelStepName, selection.NotIn, expectedSteps)
		if err != nil {
			log.Errorf("Cleaning up dangling pods, failed to construct label-selector-requirement for %v expected step(s): %v", len(expectedSteps), err)
			return
		}
		selector = selector.Add(*req)
	}

	// Query dangling pods...
	var danglingPods corev1.PodList
	err := r.client.List(context.TODO(), &client.ListOptions{Namespace: r.namespace, LabelSelector: selector}, &danglingPods)
	if err != nil {
		err = fmt.Errorf("Failed to query matching pods: %v", err)
	}

	// ... and delete them one by one
	for _, pod := range danglingPods.Items {
		r.deletePod(&pod, log.WithField("pod", pod.Name), "dangling")
	}
}
