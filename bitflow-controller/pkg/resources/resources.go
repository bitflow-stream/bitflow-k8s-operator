package resources

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/config"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MBytes = 1024 * 1024
	GBytes = 1024 * MBytes
)

type ResourceAssigner struct {
	Client     client.Client
	Config     *config.Config
	Respawning *common.RespawningPods
	Namespace  string
	PodLabels  map[string]string
}

// AssignResourcesNodeImplicit: Assumes the pod was once already deployed and still includes information
// about the former node on which this pod was placed
func (res *ResourceAssigner) AssignResourcesNodeImplicit(pod *corev1.Pod) *corev1.ResourceList {
	nodename := common.GetNodeName(pod)
	node, _ := common.RequestNode(res.Client, nodename)
	return res.AssignResources(pod, node)
}

func (res *ResourceAssigner) AssignResources(pod *corev1.Pod, node *corev1.Node) *corev1.ResourceList {
	initSize := res.Config.GetInitialResourceBufferSize()
	factor := res.Config.GetResourceBufferIncrementFactor()
	nodeInfo := res.RequestNodeInfo(node, pod.Name)
	numberOfContainers := len(pod.Spec.Containers)
	log.Debugf("Bitflow infos. buffer init size: %v, buffer increment factor: %v, node info: %v", initSize, factor, nodeInfo)

	resources := nodeInfo.GetCurrentResourceList(initSize, numberOfContainers, factor)
	if resources == nil {
		log.Warnf("Resource limit was set to %v, no resource limit will be enforced on pod", nodeInfo.TotalResourceLimit)
	}
	PatchPodResourceLimitList(pod, resources)
	return resources
}

func PatchPodResourceLimitList(pod *corev1.Pod, resourceList *corev1.ResourceList) {
	if resourceList == nil {
		return
	}
	resources := getContainerResourceLimits()
	resources.Limits = resourceList.DeepCopy()

	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].Resources = *resources
	}
}

func getContainerResourceLimits() *corev1.ResourceRequirements {
	resLimits := make(corev1.ResourceList)
	resLimits[corev1.ResourceCPU] = *resource.NewMilliQuantity(100, resource.DecimalSI)
	resLimits[corev1.ResourceMemory] = *resource.NewQuantity(512*MBytes, resource.BinarySI)

	resRequests := make(corev1.ResourceList)
	resRequests[corev1.ResourceCPU] = *resource.NewMilliQuantity(0, resource.DecimalSI)
	resRequests[corev1.ResourceMemory] = *resource.NewQuantity(0, resource.BinarySI)

	return &corev1.ResourceRequirements{
		Limits:   resLimits,
		Requests: resRequests,
	}
}

func (res *ResourceAssigner) GetCurrentResources(node *corev1.Node) *corev1.ResourceList {
	initSize := res.Config.GetInitialResourceBufferSize()
	factor := res.Config.GetResourceBufferIncrementFactor()
	nodeInfo := res.RequestNodeInfo(node, "")
	log.Debugf("Bitflow node infos during validation: init size: %v, factor: %v, Node infos: %v", initSize, factor, nodeInfo.String())
	return nodeInfo.GetCurrentResourceList(initSize, 0, factor)
}

func (res *ResourceAssigner) RequestNodeInfo(node *corev1.Node, podName string) NodeInfo {
	var nodeInfo NodeInfo
	totalLimit := RequestBitflowResourceLimitByNode(node, res.Config)
	nodeInfo.TotalResourceLimit = totalLimit
	nodeInfo.AllocatableResources = node.Status.Allocatable

	pods, err := common.RequestAllPodsOnNode(res.Client, node.Name, res.Namespace, res.PodLabels)
	if err != nil {
		log.Errorf("Failed to query all Bitflow pods on node %v: %v", node.Name, err)
		return nodeInfo
	}
	var count int
	for _, pod := range pods {
		log.Debugln("Found Bitflow pod", pod.Name)
		_, ok := res.Respawning.IsPodRestartingOnNode(pod.Name, node.Name)
		if pod.DeletionTimestamp != nil && !ok {
			log.Debugln("Not counting pod because", pod.Name, pod.DeletionTimestamp != nil, ok)
			continue
		}
		if podName != "" && pod.Name == podName {
			continue
		}
		log.Debugln("count pod")
		count += len(pod.Spec.Containers)
	}
	count += res.Respawning.CountRestarting(pods, podName, node.Name)

	nodeInfo.NumberOfBitflowContainers = count
	return nodeInfo
}

func GetAllocatableResourcesAndTotalRecourceLimit(node *corev1.Node, config *config.Config) (corev1.ResourceList, float64) {
	allocatable := node.Status.Allocatable
	totalLimit := RequestBitflowResourceLimitByNode(node, config)
	return allocatable, totalLimit
}

func RequestBitflowResourceLimitByNode(node *corev1.Node, conf *config.Config) float64 {
	annotation := conf.GetResourceLimitAnnotation()
	limit := node.Annotations[annotation]
	var limitFloat float64
	var err error
	if limit == "" {
		limitFloat = conf.GetDefaultNodeResourceLimit()
	} else {
		limitFloat, err = strconv.ParseFloat(limit, 64)
		if err != nil {
			log.Errorln("Error Parsing float", err)
			limitFloat = conf.GetDefaultNodeResourceLimit()
		}
	}
	if limitFloat <= 0 || limitFloat > 1 {
		return -1.0 // means no limit
	}
	return limitFloat
}

type NodeInfo struct {
	NumberOfBitflowContainers int
	TotalResourceLimit        float64
	AllocatableResources      corev1.ResourceList
	// TODO validate total limit
}

func (nodeInfo *NodeInfo) String() string {
	var resources bytes.Buffer
	for name, limit := range nodeInfo.AllocatableResources {
		if resources.Len() > 0 {
			resources.WriteString(", ")
		}
		_, _ = fmt.Fprintf(&resources, "%v = %v", name, limit.String())
	}
	return fmt.Sprintf("Bitflow containers: %v, total resource limit: %v, resources: %v", nodeInfo.NumberOfBitflowContainers, nodeInfo.TotalResourceLimit, resources.String())
}

func (nodeInfo *NodeInfo) GetCurrentResourceList(initSize, spawning int, factor float64) *corev1.ResourceList {
	if nodeInfo.TotalResourceLimit <= 0 || nodeInfo.TotalResourceLimit > 1 {
		return nil
	}
	totalCount := float64(nodeInfo.NumberOfBitflowContainers + spawning)
	bufferSize := float64(initSize)
	if factor <= 1.0 {
		log.Infoln("Factor is set to a value <= 1, using 2.0 instead", "Factor", factor)
		factor = 2.0
	}
	for bufferSize < totalCount {
		bufferSize = math.Round(bufferSize * factor)
	}

	cpuTotal := nodeInfo.AllocatableResources.Cpu()
	memoryTotal := nodeInfo.AllocatableResources.Memory()

	cpuLimit := float64(cpuTotal.ScaledValue(resource.Milli)) * (nodeInfo.TotalResourceLimit / bufferSize)
	memoryLimit := float64(memoryTotal.Value()) * (nodeInfo.TotalResourceLimit / bufferSize)
	cpuLimit = math.Round(cpuLimit)
	memoryLimit = math.Round(memoryLimit)
	return nodeInfo.buildCurrentResourceLimits(int64(cpuLimit), int64(memoryLimit))
}

func (nodeInfo *NodeInfo) buildCurrentResourceLimits(cpu, memory int64) *corev1.ResourceList {
	return &corev1.ResourceList{
		corev1.ResourceCPU:    *resource.NewMilliQuantity(cpu, resource.DecimalSI),
		corev1.ResourceMemory: *resource.NewQuantity(memory, resource.BinarySI),
	}
}
