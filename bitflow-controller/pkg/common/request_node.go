package common

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func RequestNode(cli client.Client, nodeName string) (*corev1.Node, error) {
	node := &corev1.Node{}
	err := cli.Get(context.TODO(), types.NamespacedName{Name: nodeName}, node)
	return node, err
}

func RequestReadyNode(cli client.Client, nodeName string) (*corev1.Node, error) {
	node, err := RequestNode(cli, nodeName)
	if err != nil {
		return node, err
	}
	if !IsNodeReady(node) {
		node = nil
		err = fmt.Errorf("Node %v is not ready", node.Name)
	}
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

func RequestReadyNodes(cli client.Client) (*corev1.NodeList, error) {
	nodes, err := RequestNodes(cli)
	if err != nil {
		return nodes, err
	}

	// Remove the nodes from the slice, that are not marked as ready
	for i := 0; i < len(nodes.Items) && len(nodes.Items) > 0; {
		if !IsNodeReady(&nodes.Items[i]) {
			copy(nodes.Items[i:], nodes.Items[i+1:])
			nodes.Items = nodes.Items[:len(nodes.Items)-1]
		} else {
			i++
		}
	}
	return nodes, nil
}

func IsNodeReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}
