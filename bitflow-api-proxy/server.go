package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ProxyServer struct {
	metrics MetricsRecordingServer

	KubeConfig        string
	KubeNamespace     string
	ConfigMapName     string
	ConfigCachePeriod time.Duration

	OperatorApiPort       string
	OperatorSelectorKey   string
	OperatorSelectorValue string

	controllerConfig *config.Config
	client           client.Client
}

func (s *ProxyServer) registerFlags() {
	flag.StringVar(&s.KubeConfig, "kube-config", "", "Path to the kubernetes config file - may be omitted if this endpoint is deployed inside a kubernetes cluster")
	flag.StringVar(&s.KubeNamespace, "kube-namespace", "default", "Kubernetes namespace to proxy objects for and retrieve configuration from")
	flag.StringVar(&s.ConfigMapName, "config-map", "bitflow-controller-config", "Name of the config map used by the Bitflow operator")
	flag.DurationVar(&s.ConfigCachePeriod, "config-cache", 0, "Amount of time a cached config value will be read from cache. If this time expires the next request to the parameter will read from kubernetes API and again cache the value. Use negative value to always load from kubernetes API.")

	flag.StringVar(&s.OperatorApiPort, "operator-api-port", "8888", "Port of the REST API of the operator")
	flag.StringVar(&s.OperatorSelectorKey, "operator-selector", "app.kubernetes.io/instance", "Kubernetes pod label selector key for finding the operator")
	flag.StringVar(&s.OperatorSelectorValue, "operator-selector-value", "bitflow-controller", "Kubernetes pod label selector value for finding the operator")

	s.metrics.registerFlags()
}

func (s *ProxyServer) init() error {
	if err := s.initKubernetes(); err != nil {
		return fmt.Errorf("Error initializing Kubernetes: %v", err)
	}
	if err := s.initBitflowConfig(); err != nil {
		return fmt.Errorf("Error initializing Bitflow Controller Config: %v", err)
	}
	return nil
}

func (s *ProxyServer) initKubernetes() error {
	var err error
	s.client, err = common.GetKubernetesClient(s.KubeConfig)
	if err != nil {
		return err
	}
	s.metrics.init(s)
	return nil
}

func (s *ProxyServer) initBitflowConfig() error {
	if s.ConfigCachePeriod > 0 {
		s.controllerConfig = config.NewConfigWithCache(s.client, s.KubeNamespace, s.ConfigMapName, s.ConfigCachePeriod)
	} else {
		s.controllerConfig = config.NewConfig(s.client, s.KubeNamespace, s.ConfigMapName)
	}
	return s.controllerConfig.Validate()
}

func (s *ProxyServer) registerEndpoints(gin *gin.Engine) {
	// Meta
	gin.GET("/ping", s.ping)
	gin.GET("/health", s.health)
	gin.GET("/config", s.getConfig)

	// Objects
	gin.GET("/steps", s.getBitflowSteps)
	gin.GET("/steps/:namespace", s.getBitflowSteps)
	gin.GET("/step/:namespace/:stepName", s.getBitflowStep)
	gin.POST("/step/:namespace", s.setBitflowStep)
	gin.GET("/datasources", s.getBitflowSources)
	gin.GET("/datasources/:namespace", s.getBitflowSources)
	gin.GET("/datasource/:namespace/:sourceName", s.getBitflowSource)
	gin.POST("/datasource/:namespace", s.setBitflowSource)
	gin.GET("/matchingSources", s.getMatchingSources)
	gin.GET("/matchingSources/:namespace", s.getMatchingSources)
	gin.GET("/matchingSources/:namespace/:stepName", s.getMatchingSources)
	gin.GET("/pods", s.getPods)
	gin.GET("/pods/:namespace", s.getPods)
	gin.GET("/pod/:namespace/:podName", s.getPod)
	gin.GET("/pod/:namespace/:podName/resources", s.getPodResources)

	// Nodes
	gin.GET("/nodes", s.getNodes)
	gin.GET("/node/:nodeName", s.getNode)
	gin.PUT("/node/:nodeName/annotate", s.annotateNode)

	// TODO fix
	// gin.GET("/node/:nodeName/resources/:namespace", s.getNodeResources)

	s.metrics.registerEndpoints(gin)
}

func (s *ProxyServer) getRequestLogLevel(ctx *gin.Context) (log.Level, string, bool) {
	if len(ctx.Errors) == 0 && ctx.Request.URL.Path == "/health" {
		// Hide health requests from the default logs, as they occur frequently
		return log.DebugLevel, "", true
	}
	return log.InfoLevel, "", false
}
