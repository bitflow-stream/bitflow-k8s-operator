package bitflow

import (
	"sync"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const (
	PodPlanned    = "planned"    // The controller will soon schedule and create the pod
	PodCreated    = "created"    // The pod has been scheduled and started
	PodRespawning = "respawning" // The pod was deleted and waits to be restarted

	// Deleted/Dangling/Obsolete pods are entirely removed from the ManagedPods map
)

type PodStatus struct {
	status string
	pod    *corev1.Pod

	step         *bitflowv1.BitflowStep
	inputSources []*bitflowv1.BitflowSource
	previousIP   string
}

func (spec *PodStatus) CountContainer() int {
	return len(spec.pod.Spec.Containers)
}

func (spec *PodStatus) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"pod":  spec.pod,
		"step": spec.step.Name,
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

func (p *ManagedPods) Count(status string) (res int) {
	p.Read(func() {
		for _, pod := range p.pods {
			if pod.status == status {
				res++
			}
		}
	})
	return
}

func (p *ManagedPods) Put(pod *corev1.Pod, status string) {
	p.Modify(func() {
		p.pods[pod.Name] = PodStatus{
			status:     status,
			pod:        pod,
			previousIP: pod.Status.PodIP,
		}
	})
}

func (p *ManagedPods) Delete(name string) {
	p.Modify(func() {
		delete(p.pods, name)
	})
}

func (p *ManagedPods) DeletePodsWithLabel(labelKey string, labelValue string) {
	p.Lock()
	defer p.Unlock()
	for key, status := range p.pods {
		if status.pod.Labels[labelKey] == labelValue {
			delete(p.pods, key)
		}
	}
}

func (p *ManagedPods) DeletePodsWithLabelExcept(labelKey string, labelValue string, valid []string) {
	p.Lock()
	defer p.Unlock()
	for key, status := range p.pods {
		if status.pod.Labels[labelKey] == labelValue {
			found := false
			for _, validPod := range valid {
				if validPod == key {
					found = true
				}
			}
			if !found {
				delete(p.pods, key)
			}
		}
	}
}

func (p *ManagedPods) IsPodRestarting(name string) (PodStatus, bool) {
	p.RLock()
	defer p.RUnlock()
	res, ok := p.pods[name]
	return res, ok
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

func (p *ManagedPods) ListPods() []string {
	p.RLock()
	defer p.RUnlock()
	podList := make([]string, 0, len(p.pods))
	for key := range p.pods {
		podList = append(podList, key)
	}
	return podList
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

func (p *ManagedPods) Debug() {
	if !log.IsLevelEnabled(log.DebugLevel) {
		return
	}
	p.RLock()
	defer p.RUnlock()
	for key := range p.pods {
		log.Debugln("RespawningEntry", key)
	}
}
