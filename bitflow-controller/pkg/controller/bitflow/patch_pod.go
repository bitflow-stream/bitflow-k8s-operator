package bitflow

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/antongulenko/golib"
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/go-bitflow/steps"
	corev1 "k8s.io/api/core/v1"
)

const (
	replacementTemplate = "{%s}"

	PodEnvStepType             = "BITFLOW_STEP_TYPE"
	PodEnvStepName             = "BITFLOW_STEP_NAME"
	PodEnvDataSource           = "BITFLOW_SOURCE"
	PodEnvOneToOneSourceUrl    = "BITFLOW_SINGLE_SOURCE_URL"
	PodEnvOneToOneSourceName   = "BITFLOW_SINGLE_SOURCE_NAME"
	PodEnvOneToOneSourceLabels = "BITFLOW_SINGLE_SOURCE_LABELS"
)

func copyMap(input map[string]string) map[string]string {
	result := make(map[string]string, len(input))
	for k, v := range input {
		result[k] = v
	}
	return result
}

func PatchOneToOnePod(pod *corev1.Pod, step *bitflowv1.BitflowStep, source *bitflowv1.BitflowSource, extraLabels, extraEnv map[string]string, apiIP string, apiPort int) {
	extraLabels = copyMap(extraLabels)
	extraEnv = copyMap(extraEnv)

	sourceString := buildDataSource(pod.Name, apiIP, apiPort)
	extraEnv[PodEnvOneToOneSourceName] = source.Name
	extraEnv[PodEnvOneToOneSourceLabels] = golib.FormatSortedMap(source.Labels)
	extraEnv[PodEnvOneToOneSourceUrl] = source.Spec.URL
	extraLabels[bitflowv1.PodLabelOneToOneSourceName] = source.Name
	patchPodMeta(pod, step, sourceString, extraLabels, extraEnv)
}

func PatchSingletonPod(pod *corev1.Pod, step *bitflowv1.BitflowStep, extraLabels, extraEnvVars map[string]string, apiIP string, apiPort int) {
	dataSource := ""
	if step.Spec.Type == bitflowv1.StepTypeAllToOne {
		dataSource = buildDataSource(pod.Name, apiIP, apiPort)
	}
	patchPodMeta(pod, step, dataSource, copyMap(extraLabels), copyMap(extraEnvVars))
}

func buildDataSource(podName, ownPodIP string, apiPort int) string {
	patchedUrl := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(ownPodIP, strconv.Itoa(apiPort)),
	}
	patchedUrl.Scheme = steps.DynamicSourceEndpointType
	patchedUrl.Path = RestApiDataSourcesPath + "/" + podName
	return patchedUrl.String()
}

func patchPodMeta(pod *corev1.Pod, step *bitflowv1.BitflowStep, sourceString string, extraLabels map[string]string, extraEnv map[string]string) {
	patchPodLabels(pod, step, extraLabels)
	env := patchPodEnv(pod, step, sourceString, extraEnv)
	patchTemplates(pod, env)
}

func patchPodLabels(pod *corev1.Pod, step *bitflowv1.BitflowStep, extraLabels map[string]string) {
	if pod.Labels == nil {
		pod.Labels = make(map[string]string)
	}
	for key, value := range extraLabels {
		pod.Labels[key] = value
	}
	pod.Labels[bitflowv1.LabelStepType] = step.Type()
	pod.Labels[bitflowv1.LabelStepName] = step.Name
}

func patchPodEnv(pod *corev1.Pod, step *bitflowv1.BitflowStep, sourceString string, extraEnvVars map[string]string) map[string]string {
	envVars := map[string]string{
		PodEnvStepType: step.Type(),
		PodEnvStepName: step.Name,
	}

	if sourceString != "" {
		envVars[PodEnvDataSource] = sourceString
	}
	for key, val := range extraEnvVars {
		envVars[key] = val
	}

	for key, val := range envVars {
		addEnv(pod, key, val)
	}
	return envVars
}

func patchTemplates(pod *corev1.Pod, envVars map[string]string) {
	replacements := make([]string, 0, len(envVars))
	for key, val := range envVars {
		replacements = append(replacements, fmt.Sprintf(replacementTemplate, key), val)
	}
	templateReplacer := strings.NewReplacer(replacements...)
	for i, container := range pod.Spec.Containers {
		for j, command := range container.Command {
			command = templateReplacer.Replace(command)
			pod.Spec.Containers[i].Command[j] = command
		}
		for j, arg := range container.Args {
			arg = templateReplacer.Replace(arg)
			pod.Spec.Containers[i].Args[j] = arg
		}
		for j, env := range container.Env {
			envVal := templateReplacer.Replace(env.Value)
			pod.Spec.Containers[i].Env[j].Value = envVal
		}
	}
}

func addEnv(pod *corev1.Pod, key, value string) {
	for i := range pod.Spec.Containers {
		varExists := false
		for j, env := range pod.Spec.Containers[i].Env {
			if env.Name == key {
				pod.Spec.Containers[i].Env[j].Value = value
				varExists = true
			}
		}
		if !varExists {
			// The variable does not exist yet - append it to the list
			pod.Spec.Containers[i].Env = append(pod.Spec.Containers[i].Env, corev1.EnvVar{Name: key, Value: value})
		}
	}
}
