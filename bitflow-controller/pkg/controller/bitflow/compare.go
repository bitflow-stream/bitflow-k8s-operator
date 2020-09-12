package bitflow

import (
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *BitflowReconciler) comparePods(pod1, pod2 *corev1.Pod) string {
	if diff := r.compareMetaData(pod1.TypeMeta, pod2.TypeMeta, pod1.ObjectMeta, pod2.ObjectMeta); diff != "" {
		return diff
	}
	if len(pod1.Spec.Containers) != len(pod2.Spec.Containers) {
		return "num-containers"
	}
	for i, container1 := range pod1.Spec.Containers {
		if diff := r.compareContainers(&container1, &pod2.Spec.Containers[i]); diff != "" {
			return fmt.Sprintf("container[%v] %v", i, diff)
		}
	}
	if !reflect.DeepEqual(pod1.Spec.Affinity, pod2.Spec.Affinity) {
		return "node-affinity"
	}

	// Compare the rest of the pod
	pod1Spec := pod1.Spec.DeepCopy()
	pod2Spec := pod2.Spec.DeepCopy()
	pod1Spec.Containers = nil
	pod2Spec.Containers = nil
	pod1Spec.Affinity = nil
	pod2Spec.Affinity = nil
	return r.compare(pod1Spec, pod2Spec, "pod-spec")
}

func (r *BitflowReconciler) compareContainers(container1, container2 *corev1.Container) string {
	// Compare these fields separately to make the ordering irrelevant. E.g., env-vars seem to come in random order,
	// although they are represented as a slice. TODO there are likely more efficient ways to do this.
	if diff := r.compareContainerPorts(container1.Ports, container2.Ports); diff != "" {
		return diff
	}
	if diff := r.compareContainerEnv(container1.Env, container2.Env); diff != "" {
		return diff
	}
	if diff := r.compareContainerEnvFrom(container1.EnvFrom, container2.EnvFrom); diff != "" {
		return diff
	}
	if diff := r.compare(container1.Resources, container2.Resources, "resources"); diff != "" {
		return diff
	}

	// Compare the remaining fields of the container
	copy1 := container1.DeepCopy()
	copy2 := container2.DeepCopy()
	copy1.Ports = nil
	copy2.Ports = nil
	copy1.Env = nil
	copy2.Env = nil
	copy1.EnvFrom = nil
	copy2.EnvFrom = nil
	copy1.Resources = corev1.ResourceRequirements{}
	copy2.Resources = corev1.ResourceRequirements{}
	return r.compare(copy1, copy2, "spec")
}

func (r *BitflowReconciler) compareContainerPorts(ports1, ports2 []corev1.ContainerPort) string {
	map1 := make(map[string]corev1.ContainerPort)
	for _, obj := range ports1 {
		map1[obj.String()] = obj
	}
	map2 := make(map[string]corev1.ContainerPort)
	for _, obj := range ports2 {
		map2[obj.String()] = obj
	}

	if !reflect.DeepEqual(map1, map2) {
		return "container-ports"
	}
	return ""
}

func (r *BitflowReconciler) compareContainerEnv(env1, env2 []corev1.EnvVar) string {
	map1 := make(map[string]corev1.EnvVar)
	for _, obj := range env1 {
		map1[obj.String()] = obj
	}
	map2 := make(map[string]corev1.EnvVar)
	for _, obj := range env2 {
		map2[obj.String()] = obj
	}

	if !reflect.DeepEqual(map1, map2) {
		return "container-env"
	}
	return ""
}

func (r *BitflowReconciler) compareContainerEnvFrom(envFrom1, envFrom2 []corev1.EnvFromSource) string {
	map1 := make(map[string]corev1.EnvFromSource)
	for _, obj := range envFrom1 {
		map1[obj.String()] = obj
	}
	map2 := make(map[string]corev1.EnvFromSource)
	for _, obj := range envFrom2 {
		map2[obj.String()] = obj
	}

	if !reflect.DeepEqual(map1, map2) {
		return "container-env-from"
	}
	return ""
}

func (r *BitflowReconciler) compareObjects(type1, type2 metav1.TypeMeta, meta1, meta2 metav1.ObjectMeta, spec1, spec2 interface{}) string {
	if diff := r.compareMetaData(type1, type2, meta1, meta2); diff != "" {
		return diff
	}
	return r.compare(spec1, spec2, "spec")
}

func (r *BitflowReconciler) compareMetaData(type1, type2 metav1.TypeMeta, meta1, meta2 metav1.ObjectMeta) string {
	// The meta data contains many fields that are populated by the system. Ignore those fields, and explicitly compare the ones that
	// were potentially populated by Bitflow or the user.
	return r.compareMultipleAspects(
		[3]interface{}{type1, type2, "TypeMeta"},
		[3]interface{}{meta1.Name, meta2.Name, "Name"},
		[3]interface{}{meta1.GenerateName, meta2.GenerateName, "GenerateName"},
		[3]interface{}{meta1.Namespace, meta2.Namespace, "Namespace"},
		[3]interface{}{meta1.ClusterName, meta2.ClusterName, "ClusterName"},
		[3]interface{}{meta1.DeletionGracePeriodSeconds, meta2.DeletionGracePeriodSeconds, "DeletionGracePeriodSeconds"},
		[3]interface{}{meta1.Labels, meta2.Labels, "Labels"},
		[3]interface{}{meta1.Annotations, meta2.Annotations, "Annotations"},
		[3]interface{}{meta1.OwnerReferences, meta2.OwnerReferences, "OwnerReferences"},
		[3]interface{}{meta1.Initializers, meta2.Initializers, "Initializers"},
		[3]interface{}{meta1.Finalizers, meta2.Finalizers, "Finalizers"})
}

func (r *BitflowReconciler) compareMultipleAspects(aspects ...[3]interface{}) string {
	for _, aspect := range aspects {
		if diff := r.compare(aspect[0], aspect[1], aspect[2].(string)); diff != "" {
			return diff
		}
	}
	return ""
}

func (r *BitflowReconciler) compare(obj1, obj2 interface{}, description string) string {
	if !reflect.DeepEqual(obj1, obj2) {
		return description
	}
	return ""
}
