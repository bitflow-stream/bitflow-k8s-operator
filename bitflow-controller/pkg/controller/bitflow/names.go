package bitflow

import "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"

const (
	// TODO make these prefix/suffixes configurable
	STEP_OUTPUT_PREFIX = "output"
	STEP_POD_SUFFIX    = "pod"
	POD_SEPARATOR      = "-"
	SOURCE_SEPARATOR   = "."
)

func ConstructReproduciblePodName(stepName, sourceName string) string {
	return common.HashName(ConstructSingletonPodName(stepName)+POD_SEPARATOR, stepName, sourceName)
}

func ConstructSingletonPodName(stepName string) string {
	return common.CleanDnsName(stepName + POD_SEPARATOR + STEP_POD_SUFFIX)
}

func ConstructSourceName(podName, sourceName string) string {
	return common.CleanDnsName(STEP_OUTPUT_PREFIX + SOURCE_SEPARATOR + podName + SOURCE_SEPARATOR + sourceName)
}
