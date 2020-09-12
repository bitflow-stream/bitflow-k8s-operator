package bitflow

import (
	"sync"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// Deleted/Dangling/Obsolete pods are entirely removed from the ManagedPods map
type PodStatus struct {
	pod          *corev1.Pod
	step         *bitflowv1.BitflowStep
	inputSources []*bitflowv1.BitflowSource

	targetNode *corev1.Node
	resources  *corev1.ResourceList
	respawning bool
}

func (spec *PodStatus) CountContainer() int {
	return len(spec.pod.Spec.Containers)
}

func (spec *PodStatus) TargetNode() string {
	if spec.targetNode != nil {
		return spec.targetNode.Name
	}

	// Usually, the spec.targetNode field should be set, but in case the controller is restarted, pods might already
	// be running before the first scheduling.
	return common.GetTargetNode(spec.pod)
}

func (spec *PodStatus) Log() *log.Entry {
	logger := log.WithFields(log.Fields{
		"pod":  spec.pod.Name,
		"step": spec.step.Name,
		"type": spec.step.Type(),
	})
	if stepType := spec.step.Type(); stepType == bitflowv1.StepTypeOneToOne {
		logger = logger.WithField("source", spec.inputSources[0].Name)
	} else if stepType == bitflowv1.StepTypeAllToOne {
		logger = logger.WithField("num-initial-sources", len(spec.inputSources))
	}
	if targetNode := spec.TargetNode(); targetNode != "" {
		logger = logger.WithField("node", targetNode)
	}
	return logger
}

func (spec *PodStatus) Clone() *PodStatus {
	if spec == nil {
		return &PodStatus{} // Default values
	}
	return &PodStatus{
		pod:          spec.pod,
		step:         spec.step,
		inputSources: spec.inputSources,
		targetNode:   spec.targetNode,
		resources:    spec.resources,
		respawning:   spec.respawning,
	}
}

type ManagedPods struct {
	pods map[string]*PodStatus
	lock sync.RWMutex
}

func NewManagedPods() *ManagedPods {
	return &ManagedPods{
		pods: make(map[string]*PodStatus),
	}
}

func (p *ManagedPods) Modify(f func()) {
	p.lock.Lock()
	defer p.lock.Unlock()
	f()
}

func (p *ManagedPods) Read(f func()) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	f()
}

func (p *ManagedPods) Len() (res int) {
	p.Read(func() {
		res = len(p.pods)
	})
	return
}

func (p *ManagedPods) CleanupStep(stepName string, existingPodNames map[string]bool) {
	p.Modify(func() {
		for name, pod := range p.pods {
			if pod.step.Name == stepName && !existingPodNames[name] {
				delete(p.pods, name)
			}
		}
	})
}

func (p *ManagedPods) Put(pod *corev1.Pod, step *bitflowv1.BitflowStep, inputSources []*bitflowv1.BitflowSource) {
	p.Modify(func() {
		// Re-allocate to avoid data race
		entry := p.pods[pod.Name].Clone()
		entry.pod = pod
		entry.step = step
		entry.inputSources = inputSources
		p.pods[pod.Name] = entry
	})
}

func (p *ManagedPods) UpdateExistingPod(pod *corev1.Pod) {
	p.Modify(func() {
		if entry, exists := p.pods[pod.Name]; exists {
			// Re-allocate to avoid data race
			entry = entry.Clone()
			entry.pod = pod.DeepCopy()
			p.pods[pod.Name] = entry
		}
	})
}

func (p *ManagedPods) MarkRespawning(pod *corev1.Pod, isRespawning bool) {
	p.Modify(func() {
		if status, ok := p.pods[pod.Name]; ok {
			status.respawning = isRespawning
		}
	})
}

func (p *ManagedPods) ListRespawningPods() (result []string) {
	p.Read(func() {
		result = make([]string, 0, len(p.pods))
		for name, pod := range p.pods {
			if pod.respawning {
				result = append(result, name)
			}
		}
	})
	return
}
