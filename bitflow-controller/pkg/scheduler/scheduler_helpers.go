package scheduler

import (
	"context"
	"fmt"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (s schedulingTask) getAvailableNodes(cli client.Client) *corev1.NodeList {
	nodes, err := common.RequestNodes(cli)
	if err != nil {
		s.logger.Errorln("Failed to request available nodes:", err)
		return nil
	} else if len(nodes.Items) == 0 {
		s.logger.Errorln("Zero nodes available for scheduling")
		return nil
	}
	return nodes
}

func (s schedulingTask) listAllBitflowPods() ([]*corev1.Pod, error) {
	podList := &corev1.PodList{}
	err := s.Client.List(context.TODO(), &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(s.IdLabels),
		Namespace:     s.Namespace,
	}, podList)
	return common.UnpackPodList(podList, err)
}

func (s schedulingTask) findNodeForDataSource(source *bitflowv1.BitflowSource) (*corev1.Node, error) {
	nodeLabel := s.Config.GetStandaloneSourceLabel()
	if nodeName, ok := source.Labels[nodeLabel]; ok {
		node, err := common.RequestNode(s.Client, nodeName)
		if err == nil {
			s.logger.Debugf("%v has label %v=%v, scheduling on node %v", source, nodeLabel, nodeName, node.Name)
			return node, nil
		} else {
			return nil, fmt.Errorf("Error requesting node '%v' based on label %v=%v of %v: %v", nodeName, nodeLabel, nodeName, source, err)
		}
	}
	return s.findNodeForDataSourcePod(source)
}

func (s schedulingTask) findNodeForDataSourcePod(source *bitflowv1.BitflowSource) (*corev1.Node, error) {
	podName := source.Labels[bitflowv1.SourceLabelPodName]
	if podName == "" {
		s.logger.Debugf("%v has no label %v, cannot use it for scheduling", source, bitflowv1.SourceLabelPodName)
		return nil, nil
	}
	pod, err := common.RequestPod(s.Client, podName, s.Namespace)
	if err != nil {
		return nil, fmt.Errorf("Error requesting pod '%v' based on label %v=%v of %v: %v",
			podName, bitflowv1.SourceLabelPodName, podName, source, err)
	}
	return common.RequestNode(s.Client, common.GetNodeName(pod))
}

func (s schedulingTask) findNodeForDataSources(sources []*bitflowv1.BitflowSource) (*corev1.Node, error) {
	sourcesOnNodes := make(map[string]int)
	nodes := make(map[string]*corev1.Node)

	var i int
	for _, source := range sources {
		node, err := s.findNodeForDataSource(source)
		if err != nil || node == nil {
			if err != nil {
				s.logger.Warnln(err)
			}
			continue
		}
		sourcesOnNodes[node.Name] = sourcesOnNodes[node.Name] + 1
		nodes[node.Name] = node
	}
	i = 0
	var maxKey string
	for key, value := range sourcesOnNodes {
		if value > i {
			i = value
			maxKey = key
		}
	}
	if i == 0 {
		return nil, nil
	}
	return nodes[maxKey], nil
}
