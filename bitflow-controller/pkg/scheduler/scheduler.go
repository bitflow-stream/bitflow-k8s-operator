package scheduler

import (
	"github.com/antongulenko/golib"
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
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

	var node *corev1.Node
	successfulScheduler := ""
	for _, schedulerName := range schedulers {
		node = task.schedule(schedulerName)
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

func (s schedulingTask) schedule(schedulerName string) *corev1.Node {
	switch schedulerName {
	case "first":
		return s.getFirstNode()
	case "random":
		return s.getRandomNode()
	case "leastContainers":
		return s.getNodeWithLeastContainers()
	case "mostCPU":
		return s.getNodeWithMostFreeCPU()
	case "mostMem":
		return s.getNodeWithMostFreeMemory()
	case "sourceAffinity":
		return s.getNodeNearSource()
	default:
		s.logger.Debugln("Unknown scheduler:", schedulerName)
		return nil
	}
}
