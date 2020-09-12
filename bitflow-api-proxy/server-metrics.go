package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/antongulenko/golib"
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/controller/bitflow"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

const (
	TIME_LAYOUT = "15:04:05"
	PRECISION   = 3
)

type MetricsRecordingServer struct {
	proxy             *ProxyServer
	RecordMetricsFile string
	metrics           *metrics.Clientset

	fileNumber int
	stopper    golib.StopChan
}

func (s *MetricsRecordingServer) registerFlags() {
	flag.StringVar(&s.RecordMetricsFile, "record-metrics", "", "Path to the metrics csv output file")
}

func (s *MetricsRecordingServer) init(proxy *ProxyServer) {
	s.proxy = proxy
	if err := s.doInit(); err != nil {
		log.Warnf("Failed to initialized metrics server. Metrics endpoints will not be available: %v", err)
		s.metrics = nil
	}
}

func (s *MetricsRecordingServer) doInit() error {
	kubeConf, err := common.GetKubernetesConfig(s.proxy.KubeConfig)
	if err != nil {
		return err
	}
	s.metrics, err = metrics.NewForConfig(kubeConf)
	return err
}

func (s *MetricsRecordingServer) registerEndpoints(gin *gin.Engine) {
	if s.metrics != nil {
		gin.GET("/metrics/record/start", s.startRecord)
		gin.GET("/metrics/record/stop", s.stopRecord)
		gin.GET("/metrics/node/:nodeName", s.currentNodeMetric)
		gin.GET("/metrics/pod/:namespace/:podName", s.currentPodMetric)
	}
}

func (s *MetricsRecordingServer) currentNodeMetric(c *gin.Context) {
	nodeName := c.Param("nodeName")
	sample, err := s.metrics.MetricsV1beta1().NodeMetricses().Get(nodeName, metav1.GetOptions{})

	ret := gin.H{}
	if err != nil {
		ret["Error"] = err.Error()
	} else {
		ret["Usage"] = sample.Usage
	}
	c.JSON(200, ret)
}

func (s *MetricsRecordingServer) currentPodMetric(c *gin.Context) {
	podName := c.Param("podName")
	namespace := s.proxy.namespace(c)
	sample, err := s.metrics.MetricsV1beta1().PodMetricses(namespace).Get(podName, metav1.GetOptions{})

	ret := gin.H{}
	if err != nil {
		ret["Error"] = err.Error()
	} else {
		ret["Usage"] = sample.Containers
	}
	c.JSON(200, ret)
}

func (s *MetricsRecordingServer) startRecord(c *gin.Context) {
	if !s.stopper.IsNil() && !s.stopper.Stopped() {
		fmt.Println("Record still running. Stop the current record before starting a new one")
	}
	s.stopper = golib.NewStopChan()
	go s.recordMetrics()
	c.JSON(200, gin.H{})
}

func (s *MetricsRecordingServer) stopRecord(c *gin.Context) {
	if !s.stopper.IsNil() {
		s.stopper.Stop()
		s.fileNumber++
	} else {
		fmt.Println("Start the record before you stop it")
	}
	c.JSON(200, gin.H{})
}

func (s *MetricsRecordingServer) handleMetricDataError(msg string, err error) bool {
	if err == nil {
		return false
	}
	log.Errorln(msg, err)
	s.stopper.Stop()
	return true
}

func (s *MetricsRecordingServer) recordMetrics() {
	file, err := os.Create(s.RecordMetricsFile + strconv.Itoa(s.fileNumber) + ".csv")
	if s.handleMetricDataError("Error creating file", err) {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	nodes, err := common.RequestNodes(s.proxy.client)
	if s.handleMetricDataError("Error requesting all nodes", err) {
		return
	}

	var sample *metricsv1.NodeMetrics
	stringValues := make([]string, 6)
	var absoluteCPU int64 = 2000
	var absoluteMem int64 = 1024 * 1024 * 1024 * 4
	var resourceLimit float64
	var podCount int
	var pods []*corev1.Pod
	s.writeHeader(csvWriter)
	for {
		if s.stopper.Stopped() {
			return
		}
		for _, node := range nodes.Items {
			resourceLimit = bitflow.GetNodeResourceLimit(&node, s.proxy.controllerConfig)

			// TODO check if this label selector is valid for selecting all Bitflow pods. Otherwise, somehow obtain the ID labels of the controller here.
			pods, err = common.RequestAllPodsOnNode(s.proxy.client, node.Name, s.proxy.KubeNamespace, map[string]string{bitflowv1.LabelStepName: "*"})
			if err != nil {
				log.Errorln("Error requesting all Bitflow pods on ", node.Name, err)
				podCount = -1
			} else {
				podCount = len(pods)
			}
			absoluteCPU = node.Status.Allocatable.Cpu().MilliValue()
			absoluteMem = node.Status.Allocatable.Memory().Value()
			sample, err = s.metrics.MetricsV1beta1().NodeMetricses().Get(node.Name, metav1.GetOptions{})
			if err != nil {
				log.Errorln("Error requesting node metrics ", node.Name, err)
			} else {
				stringValues[0] = sample.Name
				stringValues[1] = sample.Timestamp.Format(TIME_LAYOUT)
				stringValues[2] = s.calculatePercentage(absoluteCPU, sample.Usage.Cpu().MilliValue())
				stringValues[3] = s.calculatePercentage(absoluteMem, sample.Usage.Memory().Value())
				stringValues[4] = strconv.Itoa(podCount)
				stringValues[5] = strconv.FormatFloat(resourceLimit, 'f', -1, 64)
				csvWriter.Write(stringValues)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func (s *MetricsRecordingServer) writeHeader(csvWriter *csv.Writer) {
	header := []string{"Node", "Time", "Cpu", "Memory", "Pods", "Limit"}
	csvWriter.Write(header)
}

func (s *MetricsRecordingServer) calculatePercentage(absolute, current int64) string {
	perc := float64(current) / float64(absolute)
	presision := math.Pow10(PRECISION)
	perc = math.Round(perc*presision) / presision
	return strconv.FormatFloat(perc, 'f', -1, 64)
}
