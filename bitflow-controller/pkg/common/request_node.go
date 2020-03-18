package common

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func RequestNode(cli client.Client, nodeName string) (*corev1.Node, error) {
	node := &corev1.Node{}
	err := cli.Get(context.TODO(), types.NamespacedName{Name: nodeName}, node)
	return node, err
}

func RequestNodes(cli client.Client) (*corev1.NodeList, error) {
	nodeList := &corev1.NodeList{}
	err := cli.List(context.TODO(), &client.ListOptions{}, nodeList)
	if err != nil {
		return nil, err
	}
	return nodeList, nil
}
