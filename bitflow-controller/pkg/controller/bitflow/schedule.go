package bitflow

import (
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/scheduler"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const (
	SchedulerNameRandom        = "random"
	SchedulerNameLeastOccupied = "least-occupied"
	SchedulerNameDefault       = SchedulerNameLeastOccupied
)

func (r *BitflowReconciler) schedulePods(nodes map[string]*corev1.Node) {
	schedulerName := r.config.GetSchedulerName()
	logger := log.WithField("schedule-strategy", schedulerName)
	sched := r.createScheduler(schedulerName, nodes, logger)

	// Run the scheduler and assign the results
	logger.Debug("Executing scheduler")
	changed, assignedNodes, err := sched.Schedule()
	if err != nil {
		logger.Errorln("Schedule routine returned error:", err)
		return
	}
	if changed {
		r.assignSchedulingResults(assignedNodes, nodes, logger)
	}
}

func (r *BitflowReconciler) createScheduler(schedulerName string, nodes map[string]*corev1.Node, logger *log.Entry) scheduler.Scheduler {
	var sched scheduler.Scheduler
	r.pods.Read(func() {
		switch schedulerName {
		case SchedulerNameRandom:
			sched = r.createRandomScheduler(r.pods.pods, nodes)
		default:
			logger.Warnf("Unknown scheduling strategy, falling back to default (%v)", SchedulerNameDefault)
			fallthrough
		case SchedulerNameLeastOccupied:
			sched = r.createLeastOccupiedScheduler(r.pods.pods, nodes)
		}
	})
	return sched
}

func (r *BitflowReconciler) assignSchedulingResults(assignedNodes map[string]string, allNodes map[string]*corev1.Node, logger *log.Entry) {
	r.pods.Modify(func() {
		for podName, nodeName := range assignedNodes {
			scheduledPod, ok1 := r.pods.pods[podName]
			targetNode, ok2 := allNodes[nodeName]
			if ok1 && ok2 {
				scheduledPod.targetNode = targetNode
			} else {
				logger.WithField("node", targetNode).WithField("pod", scheduledPod).
					Errorf("Scheduler returned unknown node and/or pod")
			}
		}
	})
}

func (r *BitflowReconciler) createLeastOccupiedScheduler(pods map[string]*PodStatus, nodes map[string]*corev1.Node) *scheduler.LeastOccupiedStaticScheduler {
	return &scheduler.LeastOccupiedStaticScheduler{
		StaticScheduler: r.createStaticScheduler(pods, nodes),
	}
}

func (r *BitflowReconciler) createRandomScheduler(pods map[string]*PodStatus, nodes map[string]*corev1.Node) *scheduler.RandomStaticScheduler {
	return &scheduler.RandomStaticScheduler{
		StaticScheduler: r.createStaticScheduler(pods, nodes),
	}
}

func (r *BitflowReconciler) createStaticScheduler(pods map[string]*PodStatus, nodes map[string]*corev1.Node) scheduler.StaticScheduler {
	currentStateNodes := r.getCurrentNodesForPods(pods)
	availableNodes := r.getAvailableNodesForPods(pods, nodes)
	dataSourceNodes := r.getDataSourceNodesForPods(pods, nodes)
	return scheduler.StaticScheduler{
		CurrentState:    currentStateNodes,
		DataSourceNodes: dataSourceNodes,
		AvailableNodes:  availableNodes,
	}
}

func (r *BitflowReconciler) getCurrentNodesForPods(pods map[string]*PodStatus) map[string]string {
	result := make(map[string]string)
	for name, pod := range pods {
		nodeName := ""
		if pod.targetNode != nil {
			nodeName = pod.targetNode.Name
		}
		result[name] = nodeName
	}
	return result
}

func (r *BitflowReconciler) getAvailableNodesForPods(pods map[string]*PodStatus, nodes map[string]*corev1.Node) map[string][]string {
	result := make(map[string][]string)
	for name, pod := range pods {
		for _, node := range nodes {
			if r.nodeMatchesLabels(node, pod.step.Spec.NodeLabels) {
				result[name] = append(result[name], node.Name)
			}
		}
	}
	return result
}

func (r *BitflowReconciler) nodeMatchesLabels(node *corev1.Node, labels map[string][]string) bool {
	nodeLabels := node.Labels
	for key, values := range labels {
		nodeValue := nodeLabels[key]
		if !r.stringSliceContains(nodeValue, values) {
			return false
		}
	}
	return true
}

func (r *BitflowReconciler) stringSliceContains(value string, stringSlice []string) bool {
	for _, sliceValue := range stringSlice {
		if sliceValue == value {
			return true
		}
	}
	return false
}

func (r *BitflowReconciler) getDataSourceNodesForPods(pods map[string]*PodStatus, nodes map[string]*corev1.Node) map[string][]string {
	result := make(map[string][]string, len(pods))
	for name, pod := range pods {
		for _, source := range pod.inputSources {
			node := r.findNodeForDataSource(source, nodes, pods)
			if node != nil {
				result[name] = append(result[name], node.Name)
			}
		}
	}
	return result
}

func (r *BitflowReconciler) findNodeForDataSource(source *bitflowv1.BitflowSource, nodes map[string]*corev1.Node, pods map[string]*PodStatus) *corev1.Node {
	// The source might contain the information what node it is associated with
	nodeLabel := r.config.GetStandaloneSourceLabel()
	if nodeName, ok := source.Labels[nodeLabel]; ok {
		if node, ok := nodes[nodeName]; ok {
			return node
		}
	}

	// Try to find the Pod that produces this data source
	podName := source.Labels[bitflowv1.SourceLabelPodName]
	if podName != "" {
		pod, ok := pods[podName]
		if ok && pod.targetNode != nil {
			podNodeName := pod.targetNode.Name
			if node, ok := nodes[podNodeName]; ok {
				return node
			}
		}
	}

	// No node could be found for this data source. It seems to be an external data source with no further info attached.
	return nil
}
