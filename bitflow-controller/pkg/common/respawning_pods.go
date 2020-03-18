package common

import (
	"sync"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type RespawningStatus struct {
	previousIP string
	Pod        *corev1.Pod
}

func (spec RespawningStatus) CountContainer() int {
	return len(spec.Pod.Spec.Containers)
}

type RespawningPods struct {
	elements map[string]RespawningStatus
	sync.RWMutex
}

func NewRespawningPods() *RespawningPods {
	return &RespawningPods{
		elements: make(map[string]RespawningStatus),
	}
}

func (respawning *RespawningPods) Size() int {
	respawning.RLock()
	defer respawning.RUnlock()
	return len(respawning.elements)
}

func (respawning *RespawningPods) Add(pod *corev1.Pod) {
	respawning.Lock()
	defer respawning.Unlock()
	respawning.elements[pod.Name] = RespawningStatus{pod.Status.PodIP, pod}
}

func (respawning *RespawningPods) Delete(name string) {
	respawning.Lock()
	defer respawning.Unlock()
	respawning.delete(name)
}

func (respawning *RespawningPods) delete(name string) {
	delete(respawning.elements, name)
}

func (respawning *RespawningPods) DeletePodsWithLabel(labelKey string, labelValue string) {
	respawning.Lock()
	defer respawning.Unlock()
	for key, status := range respawning.elements {
		if status.Pod.Labels[labelKey] == labelValue {
			respawning.delete(key)
		}
	}
}

func (respawning *RespawningPods) DeletePodsWithLabelExcept(labelKey string, labelValue string, valid []string) {
	respawning.Lock()
	defer respawning.Unlock()
	for key, status := range respawning.elements {
		if status.Pod.Labels[labelKey] == labelValue {
			found := false
			for _, validPod := range valid {
				if validPod == key {
					found = true
				}
			}
			if !found {
				respawning.delete(key)
			}
		}
	}
}

func (respawning *RespawningPods) IsPodRestarting(name string) (RespawningStatus, bool) {
	respawning.RLock()
	defer respawning.RUnlock()
	res, ok := respawning.elements[name]
	return res, ok
}

func (respawning *RespawningPods) IsPodRestartingOnNode(podName, nodeName string) (RespawningStatus, bool) {
	respawning.RLock()
	defer respawning.RUnlock()
	for key, value := range respawning.elements {
		if podName == key && value.Pod.Spec.NodeName == nodeName {
			return value, true
		}
	}
	return RespawningStatus{}, false
}

func (respawning *RespawningPods) ListPods() []string {
	respawning.RLock()
	defer respawning.RUnlock()
	podList := make([]string, 0, len(respawning.elements))
	for key := range respawning.elements {
		podList = append(podList, key)
	}
	return podList
}

func (respawning *RespawningPods) CountRestarting(pods []*corev1.Pod, currentPod, currentNode string) int {
	var count int
	var found bool
	respawning.RLock()
	defer respawning.RUnlock()
	for key, value := range respawning.elements {
		found = false
		for _, pod := range pods {
			if pod.Name == key {
				found = true
			}
		}
		if !found && key != currentPod && GetNodeName(value.Pod) == currentNode {
			log.Debugln("count containers ", key)
			count += value.CountContainer()
		}
	}
	return count
}

func (respawning *RespawningPods) Debug() {
	if !log.IsLevelEnabled(log.DebugLevel) {
		return
	}
	respawning.RLock()
	defer respawning.RUnlock()
	for key := range respawning.elements {
		log.Debugln("RespawningEntry", key)
	}
}
