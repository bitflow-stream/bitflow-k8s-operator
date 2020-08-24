package scheduler

import (
	"context"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (s schedulingTask) getAvailableNodes(cli client.Client) *corev1.NodeList {
	nodes, err := common.RequestReadyNodes(cli)
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
