package bitflow

import (
	"fmt"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PatchPodTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestPatchOneToOne() {
	s.SubTestSuite(new(PatchPodTestSuite))
}

func (s *PatchPodTestSuite) TestPatchOneToOne() {
	pod := s.Pod("pod1")
	step := s.Step("step11", "")
	source := s.Source("source11", map[string]string{"expectLabel": "Hello"})

	port := 6666
	ip := "127.0.0.2"
	PatchOneToOnePod(pod, step, source, map[string]string{"label1": "param1"}, map[string]string{"env1": "param1"}, ip, port)

	dataStr := fmt.Sprintf("dynamic://%v:%v/dataSources/%v", ip, port, pod.Name)
	expectedEnv := map[string]string{
		"env1":                     "param1",
		PodEnvStepName:             "step11",
		PodEnvStepType:             "one-to-one",
		PodEnvOneToOneSourceName:   "source11",
		PodEnvOneToOneSourceLabels: "expectLabel=Hello",
		PodEnvDataSource:           dataStr,
		PodEnvOneToOneSourceUrl:    "http://example.com"}
	expectedLabels := map[string]string{"label1": "param1", bitflowv1.LabelStepName: "step11",
		bitflowv1.LabelStepType: "one-to-one", bitflowv1.PodLabelOneToOneSourceName: "source11"}

	s.assertPodEnv(pod, expectedEnv)
	s.assertLabels(expectedLabels, pod.Labels)
}

func (s *PatchPodTestSuite) TestPatchSingleton() {
	fmt.Println("Running TestPatchSingleton")
	pod := s.Pod("pod1")
	step := s.Step("step11", bitflowv1.StepTypeAllToOne)

	port := 6666
	ip := "127.0.0.2"
	PatchSingletonPod(pod, step, map[string]string{"label1": "param1"}, map[string]string{"env1": "param1"}, ip, port)
	dataStr := buildDataSource(pod.Name, ip, port)
	expectedEnv := map[string]string{"env1": "param1", PodEnvStepName: "step11",
		PodEnvStepType: bitflowv1.StepTypeAllToOne, PodEnvDataSource: dataStr}
	expectedLabels := map[string]string{
		"label1":                "param1",
		bitflowv1.LabelStepName: "step11",
		bitflowv1.LabelStepType: bitflowv1.StepTypeAllToOne}

	s.assertPodEnv(pod, expectedEnv)
	s.assertLabels(expectedLabels, pod.Labels)
}

func (s *PatchPodTestSuite) TestPatchSingleton2() {
	pod := s.Pod("pod1")
	step := s.Step("step11", bitflowv1.StepTypeSingleton)

	port := 6666
	ip := "127.0.0.2"
	PatchSingletonPod(pod, step, map[string]string{"label1": "param1"}, map[string]string{"env1": "param1"}, ip, port)
	expectedEnv := map[string]string{"env1": "param1", PodEnvStepName: "step11",
		PodEnvStepType: bitflowv1.StepTypeSingleton}
	expectedLabels := map[string]string{"label1": "param1", bitflowv1.LabelStepName: "step11",
		bitflowv1.LabelStepType: bitflowv1.StepTypeSingleton}

	s.assertPodEnv(pod, expectedEnv)
	s.assertLabels(expectedLabels, pod.Labels)
}

func (s *PatchPodTestSuite) TestPatchSource() {
	out := &bitflowv1.StepOutput{Name: "stepout", URL: "tcp://:9000", Labels: map[string]string{"Hello": "world", "overwritten-in": "new", bitflowv1.LabelStepType: "overwritten"}}
	step := &bitflowv1.BitflowStep{ObjectMeta: metav1.ObjectMeta{Name: "teststep-1"}, Spec: bitflowv1.BitflowStepSpec{Type: "one-to-one", Outputs: []*bitflowv1.StepOutput{out}}}
	step.Validate()
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "the-input-pod"}}

	input1 := s.Source("insource1", map[string]string{"in1": "Hello", "overwritten-in": "gone", "same-in": "x", "different-in": "a", bitflowv1.PipelinePathLabelPrefix + "1": "other-step", bitflowv1.PipelineDepthLabel: "1"})
	input2 := s.Source("insource2", map[string]string{"in2": "Hello", "overwritten-in": "gone", "same-in": "x", "different-in": "b", bitflowv1.PipelinePathLabelPrefix + "1": "other-step", bitflowv1.PipelineDepthLabel: "1"})

	source := createOutputSource(step, pod, out, []*bitflowv1.BitflowSource{input1, input2}, map[string]string{"id": "extra"})
	s.NotNil(source)

	expectedLabels := map[string]string{
		"id":                                    "extra",
		"Hello":                                 "world",
		"overwritten-in":                        "new",
		"same-in":                               "x",
		bitflowv1.PipelineDepthLabel:            "2",
		bitflowv1.PipelinePathLabelPrefix + "1": "other-step",
		bitflowv1.PipelinePathLabelPrefix + "2": "teststep-1",
		bitflowv1.LabelStepType:                 "one-to-one",
		bitflowv1.LabelStepName:                 "teststep-1",
		bitflowv1.SourceLabelPodName:            "the-input-pod",
	}

	s.assertLabels(expectedLabels, source.Labels)
}

func (s *PatchPodTestSuite) assertLabels(expected, actual map[string]string) {
	for key, val := range expected {
		actual, ok := actual[key]
		s.True(ok, "Missing label %v", key)
		s.Equal(val, actual, "Wrong value '%v' for expected label %v=%v", actual, key, val)
	}
}

func (s *PatchPodTestSuite) assertPodEnv(pod *corev1.Pod, entries map[string]string) {
	for key, val := range entries {
		found := false
		for _, env := range pod.Spec.Containers[0].Env {
			if env.Name == key && val == env.Value {
				found = true
			}
		}
		s.True(found, "Did not find expected env var %v=%v in pod %v", key, val, pod.Name)
	}
}
