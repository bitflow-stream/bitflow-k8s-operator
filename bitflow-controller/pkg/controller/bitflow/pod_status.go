package bitflow

import (
	"sync"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// Deleted/Dangling/Obsolete pods are entirely removed from the ManagedPods map
type PodStatus struct {
	pod          *corev1.Pod
	step         *bitflowv1.BitflowStep
	inputSources []*bitflowv1.BitflowSource

	respawning bool
	previousIP string
}

func (spec *PodStatus) CountContainer() int {
	return len(spec.pod.Spec.Containers)
}

func (spec *PodStatus) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"pod":  spec.pod,
		"step": spec.step.Name,
		"type": spec.step.Type(),
	})
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
		entry := p.pods[pod.Name]

		// Re-allocate to avoid data race
		p.pods[pod.Name] = &PodStatus{
			pod:          pod,
			step:         step,
			inputSources: inputSources,

			// Copy these values, if an entry existed. Otherwise filled with default values.
			respawning: entry.respawning,
			previousIP: entry.previousIP,
		}
	})
}

func (p *ManagedPods) MarkRespawning(pod *corev1.Pod, isRespawning bool) {
	p.Modify(func() {
		if status, ok := p.pods[pod.Name]; ok {
			status.respawning = isRespawning
			if pod.Status.PodIP != "" {
				status.previousIP = pod.Status.PodIP
			}
		}
	})
}

func (p *ManagedPods) GetSteps() (result map[string]map[string]*PodStatus) {
	p.Read(func() {
		result = make(map[string]map[string]*PodStatus)
		for _, pod := range p.pods {
			pods, ok := result[pod.step.Name]
			if !ok {
				pods = make(map[string]*PodStatus)
			}
			pods[pod.pod.Name] = pod
		}
	})
	return
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

func (p *ManagedPods) IsPodRespawning(name string) (pod *PodStatus, isRespawning bool) {
	p.Read(func() {
		pod, isRespawning = p.pods[name]
		isRespawning = isRespawning && pod.respawning
	})
	return
}

func (p *ManagedPods) IsPodRestartingOnNode(podName, nodeName string) (PodStatus, bool) {
	p.RLock()
	defer p.RUnlock()
	for key, value := range p.pods {
		// Direktzugriff auf Spec.NodeName anstelle von GetNodeName(), da die Scheduling Affinity nicht berücksichtigt werden soll
		if podName == key && value.pod.Spec.NodeName == nodeName {
			return value, true
		}
	}
	return PodStatus{}, false
}

func (p *ManagedPods) CountRestarting(pods []*corev1.Pod, currentPod, currentNode string) int {
	var count int
	var found bool
	p.RLock()
	defer p.RUnlock()
	for key, value := range p.pods {
		found = false
		for _, pod := range pods {
			if pod.Name == key {
				found = true
			}
		}
		// Direktzugriff auf Spec.NodeName anstelle von GetNodeName(), da die Scheduling Affinity nicht berücksichtigt werden soll
		if !found && key != currentPod && value.pod.Spec.NodeName == currentNode {
			log.Debugln("count containers ", key)
			count += value.CountContainer()
		}
	}
	return count
}
