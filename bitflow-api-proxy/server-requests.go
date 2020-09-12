package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/controller/bitflow"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
)

func (s *ProxyServer) namespace(c *gin.Context) string {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = s.KubeNamespace
	}
	return namespace
}

func (s *ProxyServer) ping(c *gin.Context) {
	bitflow.ReplyJSON(c, map[string]string{"message": "pong"}, nil)
}

func (s *ProxyServer) health(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (s *ProxyServer) getConfig(c *gin.Context) {
	conf, err := s.controllerConfig.GetParameterMap()
	bitflow.ReplyJSON(c, conf, err)
}

func (s *ProxyServer) getBitflowStep(c *gin.Context) {
	step, err := common.GetStep(s.client, c.Param("stepName"), s.namespace(c))
	bitflow.ReplyJSON(c, step, err)
}

func (s *ProxyServer) getBitflowSteps(c *gin.Context) {
	steps, err := common.GetSteps(s.client, s.namespace(c))
	bitflow.ReplyJSON(c, steps, err)
}

func (s *ProxyServer) setBitflowStep(c *gin.Context) {
	var step = bitflowv1.BitflowStep{}
	if err := c.ShouldBindJSON(&step); err != nil {
		bitflow.ReplyError(c, http.StatusBadRequest, err)
		return
	}

	currentStep, err := common.GetStep(s.client, step.Name, s.namespace(c))
	if err != nil {
		bitflow.ReplyJSON(c, nil, err)
		return
	}

	step.SetResourceVersion(currentStep.GetResourceVersion())
	step.SetUID(currentStep.GetUID())
	updatedObject := step.DeepCopyObject()

	err = s.client.Update(c, updatedObject)
	bitflow.ReplyJSON(c, nil, err)
}

func (s *ProxyServer) getBitflowSource(c *gin.Context) {
	source, err := common.GetSource(s.client, c.Param("sourceName"), s.namespace(c))
	bitflow.ReplyJSON(c, source, err)
}

func (s *ProxyServer) getBitflowSources(c *gin.Context) {
	sources, err := common.GetSources(s.client, s.namespace(c))
	bitflow.ReplyJSON(c, sources, err)
}

func (s *ProxyServer) setBitflowSource(c *gin.Context) {
	var source = bitflowv1.BitflowSource{}
	if err := c.ShouldBindJSON(&source); err != nil {
		bitflow.ReplyError(c, http.StatusBadRequest, err)
		return
	}

	currentSource, err := common.GetSource(s.client, source.Name, s.namespace(c))
	if err != nil {
		bitflow.ReplyJSON(c, nil, err)
		return
	}

	source.SetResourceVersion(currentSource.GetResourceVersion())
	source.SetUID(currentSource.GetUID())
	updatedObject := source.DeepCopyObject()

	err = s.client.Update(c, updatedObject)
	bitflow.ReplyJSON(c, nil, err)
}

func (s *ProxyServer) getMatchingSources(c *gin.Context) {
	path := bitflow.RestApiMatchingSourcesPath
	if stepName := c.Param("stepName"); stepName != "" {
		path += "/" + stepName
	}
	s.forwardOperatorRequest(c, s.namespace(c), path)
}

func (s *ProxyServer) forwardOperatorRequest(c *gin.Context, namespace string, path string) {
	reverseProxy, operatorEndpoint, err := s.getOperatorReverseProxy(c, namespace)
	if err != nil {
		bitflow.ReplyError(c, http.StatusInternalServerError, err)
		return
	}
	requestUrl := *c.Request.URL // Copy the URL struct to adjust the path
	requestUrl.Host = operatorEndpoint
	requestUrl.Path = path
	proxyRequest, err := http.NewRequest(c.Request.Method, requestUrl.String(), c.Request.Body)
	reverseProxy.ServeHTTP(c.Writer, proxyRequest)
}

func (s *ProxyServer) getOperatorReverseProxy(c *gin.Context, namespace string) (*httputil.ReverseProxy, string, error) {
	operatorEndpoint, err := s.getOperatorEndpoint(namespace)
	if err != nil {
		return nil, "", err
	}
	scheme := c.Request.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}
	parsedUrl, err := url.Parse(fmt.Sprintf("%v://%v", scheme, operatorEndpoint))
	if err != nil {
		return nil, "", err
	}
	return httputil.NewSingleHostReverseProxy(parsedUrl), operatorEndpoint, nil
}

func (s *ProxyServer) getOperatorEndpoint(namespace string) (string, error) {
	operatorSelector := labels.SelectorFromSet(labels.Set{
		s.OperatorSelectorKey: s.OperatorSelectorValue,
	})
	pods, err := common.RequestSelectedPods(s.client, namespace, operatorSelector)
	if err == nil && len(pods) != 1 {
		err = fmt.Errorf("Found %v pods that match the operator pod selector %v (expected exactly 1)", len(pods), operatorSelector)
	}
	if err != nil {
		return "", err
	}
	operatorPod := pods[0]
	ip := operatorPod.Status.PodIP
	if ip == "" {
		err = fmt.Errorf("Operator pod %v has no associated IP", operatorPod.Name)
	}
	return net.JoinHostPort(ip, s.OperatorApiPort), err
}

func (s *ProxyServer) getPod(c *gin.Context) {
	pod, err := common.RequestPod(s.client, c.Param("podName"), s.namespace(c))
	bitflow.ReplyJSON(c, pod, err)
}

func (s *ProxyServer) getPods(c *gin.Context) {
	pods, err := common.RequestPods(s.client, s.namespace(c))
	bitflow.ReplyJSON(c, pods, err)
}

func (s *ProxyServer) getPodResources(c *gin.Context) {
	podResources, err := common.RequestPodResources(s.client, c.Param("podName"), s.namespace(c))
	bitflow.ReplyJSON(c, podResources, err)
}

func (s *ProxyServer) getNode(c *gin.Context) {
	node, err := common.RequestNode(s.client, c.Param("nodeName"))
	bitflow.ReplyJSON(c, node, err)
}

func (s *ProxyServer) getNodes(c *gin.Context) {
	nodes, err := common.RequestNodes(s.client)
	bitflow.ReplyJSON(c, nodes, err)
}

func (s *ProxyServer) annotateNode(c *gin.Context) {
	nodeName := c.Param("nodeName")
	var annotations map[string]string
	if err := c.ShouldBindJSON(&annotations); err != nil {
		bitflow.ReplyError(c, http.StatusBadRequest, err)
		return
	}

	err := s.putNodeAnnotations(nodeName, annotations)
	bitflow.ReplyJSON(c, nil, err)
}

func (s *ProxyServer) putNodeAnnotations(nodeName string, annotations map[string]string) error {
	node, err := common.RequestNode(s.client, nodeName)
	if err != nil {
		return err
	}
	for key, value := range annotations {
		node.Annotations[key] = value
	}
	return s.client.Update(context.TODO(), node)
}

// TODO fix. We need access to ID labels of the controller. Maybe integrate this entire Proxy API with the controller?
/*
func (s *ProxyServer) getNodeResources(c *gin.Context) {
	nodeName := c.Param("nodeName")
	namespace := s.namespace(c)

	resourceList, err := common.RequestNodeResources(s.client, nodeName)
	if err != nil {
		bitflow.ReplyError(c, http.StatusInternalServerError, err)
		return
	}

	resources, errR := common.RequestPodResourcesOnNode(s.client, nodeName, namespace)
	bitflowLimit := common.RequestBitflowResourceLimitByName(s.client, s.controllerConfig, nodeName)
	bitflowPods, errBit := common.RequestAllBitflowPodsOnNode(s.client, nodeName, namespace)

	ret := make(map[string]interface{})
	ret["Allocatable CPU"] = resourceList.Cpu()
	ret["Allocatable Memory"] = resourceList.Memory() // .Value()

	//	if errM != nil {
	//		ret["Node Metrics - Error"] = errM.Error()
	//	} else {
	//		ret["Node Metrics"] = metrics
	//	}

	if errR != nil {
		ret["Nodes pod Resources - Error"] = errR.Error()
	} else {
		ret["Nodes pod Resources"] = resources
	}

	ret["Bitflow Resource limit"] = bitflowLimit

	if errBit != nil {
		ret["Bitflow Pods - Error"] = errBit.Error()
	} else {
		ret["Bitflow Pods"] = len(bitflowPods)
	}
	c.JSON(200, ret)
}

func RequestNodeResources(cli client.Client, nodeName string) (*corev1.ResourceList, error) {
	node, err := RequestNode(cli, nodeName)
	if err != nil {
		return nil, err
	}
	return &node.status.Allocatable, nil
}

func RequestBitflowResourceLimitByName(cli client.Client, config *config.Config, nodeName string) float64 {
	node, err := RequestNode(cli, nodeName)
	if err != nil {
		log.Errorln("Error getting Node", err)
		return -1.0
	}
	return RequestBitflowResourceLimitByNode(node, config)
}


func RequestPodResourcesOnNode(cli client.Client, nodeName, namespace string) (ResourceWrapper, error) {
	var resources ResourceWrapper

	pods, err := RequestAllPodsOnNode(cli, nodeName, namespace, nil)
	if err != nil {
		return resources, err
	}

	for _, pod := range pods {
		resources.AddResources(RequestResourcesOnPod(&pod))
	}
	return resources, nil
}

*/
