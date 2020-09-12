package bitflow

import (
	"math"
	"strconv"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/config"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var zeroQuantity = resource.MustParse("0")

func (r *BitflowReconciler) assignPodResources(nodes map[string]*corev1.Node) {
	initSize := r.config.GetInitialResourceBufferSize()
	factor := r.config.GetResourceBufferIncrementFactor()

	r.pods.Modify(func() {
		// First, associate all pods with their respective nodes
		podsOnNodes := make(map[string]map[string]*PodStatus)
		for _, pod := range r.pods.pods {
			if targetNode := pod.TargetNode(); targetNode != "" {
				pods, ok := podsOnNodes[targetNode]
				if !ok {
					pods = make(map[string]*PodStatus)
					podsOnNodes[targetNode] = pods
				}
				pods[pod.pod.Name] = pod
			}
		}

		for nodeName, pods := range podsOnNodes {
			node := nodes[nodeName]
			numContainers := 0
			for _, pod := range pods {
				numContainers += len(pod.pod.Spec.Containers)
			}
			for _, pod := range pods {
				totalLimit := GetNodeResourceLimit(node, r.config)
				allocatableResources := node.Status.Allocatable
				resources := buildPodResourceList(initSize, factor, totalLimit, numContainers, allocatableResources)

				// Store the computed resources. We have exclusive write-access to the managed pods.
				pod.resources = resources
			}
		}
	})
}

func GetNodeResourceLimit(node *corev1.Node, conf *config.Config) float64 {
	annotation := conf.GetResourceLimitAnnotation()
	limit := node.Annotations[annotation]
	var limitFloat float64
	var err error
	if limit == "" {
		limitFloat = conf.GetDefaultNodeResourceLimit()
	} else {
		limitFloat, err = strconv.ParseFloat(limit, 64)
		if err != nil {
			log.Errorf("Error parsing resource limit (%v): %v, using default", limit, err)
			limitFloat = conf.GetDefaultNodeResourceLimit()
		}
	}
	if limitFloat <= 0 || limitFloat > 1 {
		return -1.0 // means no limit
	}
	return limitFloat
}

func buildPodResourceList(initSize int, factor float64, resourceLimit float64, numContainers int, allocatableResources corev1.ResourceList) *corev1.ResourceList {
	if resourceLimit <= 0 || resourceLimit > 1 {
		return nil
	}
	totalCount := float64(numContainers)
	bufferSize := float64(initSize)
	if factor <= 1.0 {
		log.Infoln("Factor is set to a value <= 1, using 2.0 instead", "Factor", factor)
		factor = 2.0
	}
	for bufferSize < totalCount {
		bufferSize = math.Round(bufferSize * factor)
	}

	cpuTotal := allocatableResources.Cpu()
	memoryTotal := allocatableResources.Memory()

	cpuLimit := float64(cpuTotal.ScaledValue(resource.Milli)) * (resourceLimit / bufferSize)
	memoryLimit := float64(memoryTotal.Value()) * (resourceLimit / bufferSize)
	cpuLimit = math.Round(cpuLimit)
	memoryLimit = math.Round(memoryLimit)

	// Use MustParse function to make the quantity comparable to quantities
	cpuResources := resource.MustParse(strconv.Itoa(int(cpuLimit)) + "m") // Rounded Milli-CPUs
	memResources := resource.MustParse(strconv.Itoa(int(memoryLimit)))

	return &corev1.ResourceList{
		corev1.ResourceCPU:    cpuResources,
		corev1.ResourceMemory: memResources,
	}
}

func patchPodResourceLimits(pod *corev1.Pod, resourceList *corev1.ResourceList) {
	resources := &corev1.ResourceRequirements{
		Limits: resourceList.DeepCopy(),
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    zeroQuantity,
			corev1.ResourceMemory: zeroQuantity,
		},
	}
	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].Resources = *resources
	}
}
