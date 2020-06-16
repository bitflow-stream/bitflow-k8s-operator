package common

import (
	"fmt"

	"github.com/antongulenko/golib"
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	TestNamespace = "default"

	TestNodeCpu = int64(2000)
	TestNodeMem = int64(4 * 1024 * 1024)

	TestNodeBufferInitSize = 4
	TestNodeBufferFactor   = 2.0
	TestNodeResourceLimit  = 0.1
)

type AbstractTestSuite struct {
	golib.AbstractTestSuite
}

func (suite *AbstractTestSuite) MakeFakeClient(objects ...runtime.Object) client.Client {
	suite.NoError(bitflowv1.SchemeBuilder.AddToScheme(scheme.Scheme))
	return fake.NewFakeClient(objects...)
}

func (suite *AbstractTestSuite) Pod(name string) *corev1.Pod {
	return suite.PodLabels(name, nil)
}

func (suite *AbstractTestSuite) PodInitializingLabelsSettingDefaultLabel(name string) *corev1.Pod {
	return suite.PodLabels(name, map[string]string{"bitflow": "true"})
}

func (suite *AbstractTestSuite) PodLabels(name string, labels map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: TestNamespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func (suite *AbstractTestSuite) Node(name string) *corev1.Node {
	annotations := map[string]string{"bitflow-resource-limit": "0.1"}
	labels := map[string]string{"test-node": "yes", "kubernetes.io/hostname": name}
	return suite.Node2(name, labels, annotations)
}

func (suite *AbstractTestSuite) Node2(name string, labels, annotations map[string]string) *corev1.Node {
	return suite.NodeWithResources(name, labels, annotations, TestNodeCpu, TestNodeMem)
}

func (suite *AbstractTestSuite) SchedulerNode() *corev1.Node {
	return suite.Node2("node1",
		map[string]string{"test-node": "yes", HostnameLabel: "node1"},
		map[string]string{"bitflow-resource-limit": "0.1"})
}

func (s *AbstractTestSuite) NodeWithCpuAndMemory(name string, cpu int64, memory int64, resourceLimit float64) *corev1.Node {
	return s.NodeWithResources(name,
		map[string]string{"test-node": "yes", HostnameLabel: name},
		map[string]string{"bitflow-resource-limit": fmt.Sprintf("%f", resourceLimit)},
		cpu,
		memory*1024*1024)
}

func (suite *AbstractTestSuite) NodeWithResources(name string, labels, annotations map[string]string, cpu int64, memory int64) *corev1.Node {
	var res = make(corev1.ResourceList)
	res[corev1.ResourceCPU] = *resource.NewMilliQuantity(cpu, resource.DecimalSI)
	res[corev1.ResourceMemory] = *resource.NewQuantity(memory, resource.BinarySI)
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.NodeSpec{},
		Status: corev1.NodeStatus{
			Allocatable: res,
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeReady,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}
}

func (suite *AbstractTestSuite) ConfigMap(name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: TestNamespace,
			Labels:    nil,
		},
		Data: map[string]string{
			"resource.buffer.init":   fmt.Sprintf("%v", TestNodeBufferInitSize),
			"resource.buffer.factor": fmt.Sprintf("%v", TestNodeBufferFactor),
			"resource.limit":         fmt.Sprintf("%v", TestNodeResourceLimit),
		},
	}
}

func (suite *AbstractTestSuite) Step(name string, stepType string, ingestMatches ...string) *bitflowv1.BitflowStep {
	return suite.StepCustomSchedulers(name, stepType, "sourceAffinity,first", ingestMatches...)
}

func (suite *AbstractTestSuite) StepCustomSchedulers(name string, stepType string, schedulerList string, ingestMatches ...string) *bitflowv1.BitflowStep {
	step := &bitflowv1.BitflowStep{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: TestNamespace,
		},
		Spec: bitflowv1.BitflowStepSpec{
			Type:      stepType,
			Scheduler: schedulerList,
			Template:  suite.PodLabels("", map[string]string{"step": name}),
		},
	}
	suite.AddStepIngest(step, ingestMatches...)
	return step
}

func (suite *AbstractTestSuite) StepWithOutput(name string, stepType string, outputName string, outputLabels map[string]string, ingestMatches ...string) *bitflowv1.BitflowStep {
	step := suite.Step(name, stepType)
	suite.AddStepOutput(step, outputName, outputLabels)
	suite.AddStepIngest(step, ingestMatches...)
	return step
}

func (suite *AbstractTestSuite) AddStepIngest(step *bitflowv1.BitflowStep, matches ...string) {
	if len(matches) > 0 {
		suite.Equal(0, len(matches)%2)
		ingests := make([]*bitflowv1.IngestMatch, 0, len(matches)/2)
		for i := 0; i < len(matches); i += 2 {
			ingests = append(ingests, &bitflowv1.IngestMatch{
				Key:   matches[i],
				Value: matches[i+1],
				Check: bitflowv1.MatchCheckWildcard,
			})
		}
		step.Spec.Ingest = ingests
	}
}

func (suite *AbstractTestSuite) AddStepOutput(step *bitflowv1.BitflowStep, name string, labels map[string]string) {
	step.Spec.Outputs = append(step.Spec.Outputs,
		&bitflowv1.StepOutput{
			Name:   name,
			URL:    "tcp://:9000",
			Labels: labels,
		})
}

func (suite *AbstractTestSuite) Source(name string, labels map[string]string) *bitflowv1.BitflowSource {
	return &bitflowv1.BitflowSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: TestNamespace,
			Labels:    labels,
		},
		Spec: bitflowv1.BitflowSourceSpec{
			URL: "http://example.com",
		},
	}
}
