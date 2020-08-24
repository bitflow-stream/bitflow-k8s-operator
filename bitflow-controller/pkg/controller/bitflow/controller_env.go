package bitflow

import (
	"fmt"
	"os"
	"strconv"

	"github.com/antongulenko/golib"
	log "github.com/sirupsen/logrus"
)

const (
	// In addition to these env vars, the variable WATCH_NAMESPACE and POD_NAME are important for the Operator framework
	EnvOperatorName        = "OPERATOR_NAME" // This variable is also used internally by the Kubernetes SDK
	EnvConfigMapName       = "CONFIG_MAP"
	EnvOwnPodIp            = "POD_IP"
	EnvRestApiPort         = "API_LISTEN_PORT"
	EnvConcurrentReconcile = "CONCURRENT_RECONCILE"
	EnvPodIdLabels         = "POD_ID_LABELS"
	EnvRecordStatistics    = "RECORD_STATISTICS"
)

type ControllerParameters struct {
	operatorName        string
	ownPodIP            string
	apiPort             int
	configMapName       string
	concurrentReconcile int
	controllerIdLabels  map[string]string // Attached to all managed objects to identify them
	recordStatistics    bool
}

func readControllerEnvVars() (ControllerParameters, error) {
	var r ControllerParameters

	// Required variables
	r.operatorName = os.Getenv(EnvOperatorName)
	r.configMapName = os.Getenv(EnvConfigMapName)
	r.ownPodIP = os.Getenv(EnvOwnPodIp)
	apiPortStr := os.Getenv(EnvRestApiPort)
	controllerLabelsStr := os.Getenv(EnvPodIdLabels)

	// Optional variables
	concurrentReconcileStr := os.Getenv(EnvConcurrentReconcile)
	recordStatisticsStr := os.Getenv(EnvRecordStatistics)

	// Make sure the required variables are present
	var missing []string
	if r.operatorName == "" {
		missing = append(missing, EnvOperatorName)
	}
	if r.configMapName == "" {
		missing = append(missing, EnvConfigMapName)
	}
	if r.ownPodIP == "" {
		missing = append(missing, EnvOwnPodIp)
	}
	if apiPortStr == "" {
		missing = append(missing, EnvRestApiPort)
	}
	if controllerLabelsStr == "" {
		missing = append(missing, EnvPodIdLabels)
	}
	if len(missing) > 0 {
		return r, fmt.Errorf("Missing one ore more required environment variable(s): %v", missing)
	}

	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		return r, fmt.Errorf("Failed to parse %v=%v as integer: %v", EnvRestApiPort, apiPortStr, err)
	}
	r.apiPort = apiPort

	controllerLabels := golib.ParseMap(controllerLabelsStr)
	if len(controllerLabels) == 0 {
		return r, fmt.Errorf("Need non-empty list of controller labels. Have: %v", controllerLabelsStr)
	}
	r.controllerIdLabels = controllerLabels

	if concurrentReconcileStr != "" {
		r.concurrentReconcile, err = strconv.Atoi(concurrentReconcileStr)
		if err != nil {
			return r, fmt.Errorf("Failed to parse %v=%v as integer: %v", EnvConcurrentReconcile, concurrentReconcileStr, err)
		}
	}
	if r.concurrentReconcile <= 0 {
		r.concurrentReconcile = 1
	}

	if recordStatisticsStr != "" {
		r.recordStatistics, err = strconv.ParseBool(recordStatisticsStr)
		if err != nil {
			return r, fmt.Errorf("Failed to parse %v=%v as bool: %v", EnvRecordStatistics, recordStatisticsStr, err)
		}
	}

	log.Infof("Loaded configuration from environment variables:")
	log.Infof("%v: %v = %v", "OperatorName", EnvOperatorName, r.operatorName)
	log.Infof("%v: %v = %v", "ConfigMap", EnvConfigMapName, r.configMapName)
	log.Infof("%v: %v = %v", "Controller ID labels", EnvPodIdLabels, r.controllerIdLabels)
	log.Infof("%v: %v = %v", "pod IP", EnvOwnPodIp, r.ownPodIP)
	log.Infof("%v: %v = %v", "REST API port", EnvRestApiPort, r.apiPort)
	log.Infof("%v: %v = %v", "Concurrent reconcile routines", EnvConcurrentReconcile, r.concurrentReconcile)
	log.Infof("%v: %v = %v", "Record statistics", EnvRecordStatistics, r.recordStatistics)
	return r, nil
}
