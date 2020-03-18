package common

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UnpackPodList(list *corev1.PodList, err error) ([]*corev1.Pod, error) {
	if err != nil {
		return nil, err
	}
	pods := make([]*corev1.Pod, len(list.Items))
	for i, pod := range list.Items {
		podCopy := pod
		pods[i] = &podCopy
	}
	return pods, nil
}

func RequestPod(cli client.Client, podName, namespace string) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := cli.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: namespace}, pod)
	return pod, err
}

func RequestPods(cli client.Client, namespace string) ([]*corev1.Pod, error) {
	return RequestSelectedPods(cli, namespace, nil)
}

func RequestSelectedPods(cli client.Client, namespace string, selector labels.Selector) ([]*corev1.Pod, error) {
	podList := &corev1.PodList{}
	err := cli.List(context.TODO(), &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	}, podList)
	return UnpackPodList(podList, err)
}

func RequestAllPodsOnNode(cli client.Client, nodeName, namespace string, podLabels map[string]string) ([]*corev1.Pod, error) {
	podList := &corev1.PodList{}
	err := cli.List(context.TODO(), &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(podLabels),
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": nodeName}),
		Namespace:     namespace,
	}, podList)
	return UnpackPodList(podList, err)
}

func RequestPodResources(cli client.Client, podName, namespace string) (ResourceWrapper, error) {
	var resources ResourceWrapper
	pod, err := RequestPod(cli, podName, namespace)
	if err != nil {
		return resources, err
	}
	resources.AddPod(pod)
	return resources, nil
}

type ResourceWrapper struct {
	TotalCpuLimit      resource.Quantity
	TotalMemoryLimit   resource.Quantity
	TotalCpuRequest    resource.Quantity
	TotalMemoryRequest resource.Quantity
}

func (wrapper *ResourceWrapper) AddPod(pod *corev1.Pod) {
	for _, container := range pod.Spec.Containers {
		wrapper.AddLimits(container.Resources.Limits)
		wrapper.AddRequests(container.Resources.Requests)
	}
}

func (wrapper *ResourceWrapper) AddLimits(resourceLimits corev1.ResourceList) {
	wrapper.TotalCpuLimit.Add(*(resourceLimits.Cpu()))
	wrapper.TotalMemoryLimit.Add(*(resourceLimits.Memory()))
}

func (wrapper *ResourceWrapper) AddRequests(resourceRequests corev1.ResourceList) {
	wrapper.TotalCpuRequest.Add(*(resourceRequests.Cpu()))
	wrapper.TotalMemoryRequest.Add(*(resourceRequests.Memory()))
}

func (wrapper *ResourceWrapper) AddResources(resources ResourceWrapper) {
	wrapper.TotalCpuLimit.Add(resources.TotalCpuLimit)
	wrapper.TotalMemoryLimit.Add(resources.TotalMemoryLimit)
	wrapper.TotalCpuRequest.Add(resources.TotalCpuRequest)
	wrapper.TotalMemoryRequest.Add(resources.TotalMemoryRequest)
}
