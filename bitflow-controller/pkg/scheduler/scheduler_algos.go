package scheduler

import (
	"math"
	"math/rand"
	"time"

	corev1 "k8s.io/api/core/v1"
)

var nodePickerRand = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

func (s schedulingTask) getFirstNode(nodes *corev1.NodeList) *corev1.Node {
	if nodes == nil {
		return nil
	}

	return &nodes.Items[0]
}

func (s schedulingTask) getRandomNode(nodes *corev1.NodeList) *corev1.Node {
	if nodes == nil {
		return nil
	}

	return &nodes.Items[nodePickerRand.Intn(len(nodes.Items))]
}

func (s schedulingTask) getNodeWithLeastContainers(nodes *corev1.NodeList) *corev1.Node {
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}

	pods, err := s.listAllBitflowPods()
	if err != nil {
		s.logger.Errorln("Failed to get Bitflow pods", err)
		return nil
	}

	nodeCountMap := make(map[string]int)

	for _, pod := range pods {
		nodeCountMap[pod.Spec.NodeName] += 1
	}

	var min = math.MaxInt32
	var minNode *corev1.Node
	for _, node := range nodes.Items {
		if minNode == nil || nodeCountMap[node.Name] < min {
			min = nodeCountMap[node.Name]
			minNode = node.DeepCopy()
		}
	}

	return minNode
}

func (s schedulingTask) getNodeWithMostFreeCPU(nodes *corev1.NodeList) *corev1.Node {
	if nodes == nil {
		return nil
	}

	availableCpu := nodes.Items[0].Status.Allocatable.Cpu()
	maxIndex := 0
	for i, node := range nodes.Items {
		if availableCpu.Cmp(*node.Status.Allocatable.Cpu()) < 0 {
			availableCpu = node.Status.Allocatable.Cpu()
			maxIndex = i
		}
	}
	return &nodes.Items[maxIndex]
}

func (s schedulingTask) getNodeWithMostFreeMemory(nodes *corev1.NodeList) *corev1.Node {
	if nodes == nil {
		return nil
	}

	availableMem := nodes.Items[0].Status.Allocatable.Memory()
	maxIndex := 0
	for i, node := range nodes.Items {
		if availableMem.Cmp(*node.Status.Allocatable.Memory()) < 0 {
			availableMem = node.Status.Allocatable.Memory()
			maxIndex = i
		}
	}
	return &nodes.Items[maxIndex]
}

func (s schedulingTask) getNodeNearSource(nodes *corev1.NodeList) *corev1.Node {
	var node *corev1.Node
	var err error
	switch len(s.sources) {
	case 0:
		return nil
	case 1:
		node, err = s.findNodeForDataSource(s.sources[0])
	default:
		node, err = s.findNodeForDataSources(s.sources)
	}
	if err != nil {
		s.logger.Errorln("Failed to query node for data source(s)", err)
	}

	if node != nil {
		for _, n := range nodes.Items {
			if n.Name == node.Name {
				return node
			}
		}
	}
	return node
}
