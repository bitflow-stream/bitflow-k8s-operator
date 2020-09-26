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

	// Compare the rest of the pod spec
	clean := func(pod *corev1.Pod) *corev1.PodSpec {
		spec := pod.Spec.DeepCopy()
		spec.Containers = nil
		spec.Affinity = nil

		// The following fields are partially automatically populated by Kubernetes
		// TODO they should still be compared somehow
		spec.Volumes = nil
		spec.RestartPolicy = ""
		spec.TerminationGracePeriodSeconds = nil
		spec.DNSPolicy = ""
		spec.ServiceAccountName = ""
		spec.DeprecatedServiceAccount = ""
		spec.NodeName = ""
		spec.SecurityContext = nil
		spec.SchedulerName = ""
		spec.Tolerations = nil
		spec.Priority = nil
		spec.EnableServiceLinks = nil
		return spec
	}
	return r.compare(clean(pod1), clean(pod2), "pod-spec")
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
	clean := func(container *corev1.Container) *corev1.Container {
		container = container.DeepCopy()
		container.Ports = nil
		container.Env = nil
		container.EnvFrom = nil
		container.Resources = corev1.ResourceRequirements{}
		// Ignore VolumeMounts, because Kubernetes adds some automatically. TODO should still be checked somehow.
		container.VolumeMounts = nil
		return container
	}
	return r.compare(clean(container1), clean(container2), "spec")
}

func (r *BitflowReconciler) compareContainerPorts(ports1, ports2 []corev1.ContainerPort) string {
	map1 := make(map[string]corev1.ContainerPort)
	for _, obj := range ports1 {
		if obj.Protocol == "" {
			obj.Protocol = corev1.ProtocolTCP
		}
		map1[obj.String()] = obj
	}
	map2 := make(map[string]corev1.ContainerPort)
	for _, obj := range ports2 {
		if obj.Protocol == "" {
			obj.Protocol = corev1.ProtocolTCP
		}
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
		[3]interface{}{meta1.OwnerReferences, meta2.OwnerReferences, "OwnerReferences"},
		[3]interface{}{meta1.Initializers, meta2.Initializers, "Initializers"},
		[3]interface{}{meta1.Finalizers, meta2.Finalizers, "Finalizers"},

		// Annotations are excluded, because Kubernetes adds some Annotations automatically
		// TODO find way to still check, if the necessary annotations exist
		// [3]interface{}{meta1.Annotations, meta2.Annotations, "Annotations"},
	)
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
