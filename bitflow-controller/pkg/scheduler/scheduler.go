package scheduler

import (
	"github.com/antongulenko/golib"
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/config"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Scheduler struct {
	Client    client.Client
	Config    *config.Config
	Namespace string
	IdLabels  map[string]string
}

func (s *Scheduler) SchedulePod(pod *corev1.Pod, step *bitflowv1.BitflowStep, sources []*bitflowv1.BitflowSource) (*corev1.Node, string) {
	logger := step.Log().WithField("pod", pod.Name).WithField("num-sources", len(sources))
	task := schedulingTask{Scheduler: s, logger: logger, pod: pod, sources: sources}

	schedulers := s.getSchedulerList(step)
	logger.Debugln("Running schedulers:", schedulers)

	var nodeList *corev1.NodeList
	var err error
	if step.Spec.NodeLabels != nil {
		nodeList, err = common.RequestReadyNodesByLabels(s.Client, step.Spec.NodeLabels)
	} else {
		nodeList, err = common.RequestReadyNodes(s.Client)
	}

	if err != nil {
		logger.Errorf("Failed to request available nodes: %v", err)
		return nil, ""
	}

	if len(nodeList.Items) == 0 {
		logger.Errorln("List of available nodes is empty")
		return nil, ""
	}

	var node *corev1.Node
	successfulScheduler := ""
	for _, schedulerName := range schedulers {
		node = task.schedule(schedulerName, nodeList)
		successfulScheduler = schedulerName
		if node != nil {
			break
		}
	}

	if node != nil {
		logger.Debugln("Scheduling on node", node.Name, "selected by scheduler:", successfulScheduler)
	} else {
		successfulScheduler = ""
		logger.Warnf("Failed to select node, pod will not have scheduling affinity")
	}
	return node, successfulScheduler
}

func (s *Scheduler) getSchedulerList(step *bitflowv1.BitflowStep) []string {
	schedulers := golib.ParseSlice(step.Spec.Scheduler)
	return append(schedulers, s.Config.GetDefaultScheduler()...)
}

type schedulingTask struct {
	*Scheduler
	logger  *log.Entry
	pod     *corev1.Pod
	sources []*bitflowv1.BitflowSource
}

func (s schedulingTask) schedule(schedulerName string, nodeList *corev1.NodeList) *corev1.Node {
	switch schedulerName {
	case "first":
		return s.getFirstNode(nodeList)
	case "random":
		return s.getRandomNode(nodeList)
	case "leastContainers":
		return s.getNodeWithLeastContainers(nodeList)
	case "mostCPU":
		return s.getNodeWithMostFreeCPU(nodeList)
	case "mostMem":
		return s.getNodeWithMostFreeMemory(nodeList)
	case "sourceAffinity":
		return s.getNodeNearSource(nodeList)
	case "lowestPenalty":
		return s.getNodeWithLowestPenalty(nodeList)
	default:
		s.logger.Debugln("Unknown scheduler:", schedulerName)
		return nil
	}
}
