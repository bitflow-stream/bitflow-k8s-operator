package common

import (
	"context"
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetStep(cli client.Client, step, namespace string) (*bitflowv1.BitflowStep, error) {
	instance := &bitflowv1.BitflowStep{}
	objKey := client.ObjectKey{
		Namespace: namespace,
		Name:      step,
	}
	err := cli.Get(context.TODO(), objKey, instance)
	return instance, err
}

func GetSteps(cli client.Client, namespace string) ([]*bitflowv1.BitflowStep, error) {
	return GetSelectedSteps(cli, namespace, nil)
}

func GetSelectedSteps(cli client.Client, namespace string, selector labels.Selector) ([]*bitflowv1.BitflowStep, error) {
	list := &bitflowv1.BitflowStepList{}
	err := cli.List(context.TODO(), &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	}, list)
	return list.GetItems(), err
}

func GetSource(cli client.Client, source, namespace string) (*bitflowv1.BitflowSource, error) {
	instance := &bitflowv1.BitflowSource{}
	objKey := client.ObjectKey{
		Namespace: namespace,
		Name:      source,
	}
	err := cli.Get(context.TODO(), objKey, instance)
	return instance, err
}

func GetSources(cli client.Client, namespace string) ([]*bitflowv1.BitflowSource, error) {
	return GetSelectedSources(cli, namespace, nil)
}

func GetSelectedSources(cli client.Client, namespace string, selector labels.Selector) ([]*bitflowv1.BitflowSource, error) {
	list := &bitflowv1.BitflowSourceList{}
	err := cli.List(context.TODO(), &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	}, list)
	return list.GetItems(), err
}
