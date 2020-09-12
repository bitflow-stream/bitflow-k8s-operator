package bitflow

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/antongulenko/golib"
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/config"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	toolscache "k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var idLabels = map[string]string{"reconciler": "bitflow-test"}

type BitflowControllerTestSuite struct {
	common.AbstractTestSuite
}

func TestBitflowController(t *testing.T) {
	suite.Run(t, new(BitflowControllerTestSuite))
}

type BitflowControllerTestHelpers struct {
	common.AbstractTestSuite
}

func (s *BitflowControllerTestHelpers) SetupTest() {
	// Set debug log level for tests
	golib.LogVerbose = true
	golib.ConfigureLogging()
}

func (s *BitflowControllerTestHelpers) performReconcile(r *BitflowReconciler, stepName string) (reconcile.Result, error) {
	return r.Reconcile(reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      stepName,
			Namespace: common.TestNamespace,
		},
	})
}

func (s *BitflowControllerTestHelpers) testReconcile(r *BitflowReconciler, stepName string) {
	res, err := s.performReconcile(r, stepName)
	s.NoError(err)
	s.False(res.Requeue)
	s.Zero(res.RequeueAfter)
}

func (s *BitflowControllerTestHelpers) initReconciler(objects ...runtime.Object) *BitflowReconciler {
	configMap := s.ConfigMap("bitflow-config")
	objects = append(objects, configMap)
	cl := s.MakeFakeClient(objects...)
	conf := config.NewConfig(cl, common.TestNamespace, "bitflow-config")
	conf.SilentlyUseAllDefaults()
	return &BitflowReconciler{
		client: cl,
		scheme: scheme.Scheme,
		cache:  new(MockCache),

		namespace: common.TestNamespace,
		idLabels:  idLabels,
		ownPodIP:  "",
		apiPort:   0,

		pods:      NewManagedPods(),
		config:    conf,
		statistic: nil,
	}
}

func (s *BitflowControllerTestHelpers) assignIPToPods(cl client.Client) {
	var list corev1.PodList
	s.NoError(cl.List(context.TODO(), &client.ListOptions{}, &list))
	for _, pod := range list.Items {
		pod.Status.PodIP = "127.0.0.1"
		s.NoError(cl.Update(context.TODO(), &pod))
	}
}

func (s *BitflowControllerTestHelpers) addSources(namePrefix string, numSources int, labels map[string]string, objects ...runtime.Object) []runtime.Object {
	for i := 0; i < numSources; i++ {
		objects = append(objects, s.Source(namePrefix+strconv.Itoa(i), labels))
	}
	return objects
}

func (s *BitflowControllerTestHelpers) deletePodForSource(cl client.Client, sourceName string) {
	var list corev1.PodList
	s.NoError(cl.List(context.TODO(), &client.ListOptions{}, &list))
	numPods := 0
	for _, pod := range list.Items {
		if pod.Labels[bitflowv1.PodLabelOneToOneSourceName] == sourceName {
			s.NoError(cl.Delete(context.TODO(), &pod))
			numPods++
		}
	}
	s.Equal(1, numPods, "Should have deleted exactly one pod for source %v", sourceName)
}

type MockCache struct {
}

var _ cache.Cache = &MockCache{}

func (c *MockCache) Get(_ context.Context, _ client.ObjectKey, _ runtime.Object) error {
	return nil
}

func (c *MockCache) List(_ context.Context, _ *client.ListOptions, _ runtime.Object) error {
	return nil
}

func (c *MockCache) GetInformer(_ runtime.Object) (toolscache.SharedIndexInformer, error) {
	return nil, nil
}

func (c *MockCache) GetInformerForKind(_ schema.GroupVersionKind) (toolscache.SharedIndexInformer, error) {
	return nil, nil
}

func (c *MockCache) Start(_ <-chan struct{}) error {
	return nil
}

func (c *MockCache) WaitForCacheSync(_ <-chan struct{}) bool {
	return true
}

func (c *MockCache) IndexField(_ runtime.Object, _ string, _ client.IndexerFunc) error {
	return nil
}

func (s *BitflowControllerTestHelpers) assertNoPodsExist(cl client.Client) {
	var podList corev1.PodList
	s.NoError(cl.List(context.TODO(), &client.ListOptions{}, &podList))
	s.Empty(podList.Items)
}

func (s *BitflowControllerTestHelpers) assertNoSourceExists(cl client.Client) {
	var sourceList bitflowv1.BitflowSourceList
	s.NoError(cl.List(context.TODO(), &client.ListOptions{}, &sourceList))
	s.Empty(sourceList.Items)
}

func (s *BitflowControllerTestHelpers) assertMissingPod(cl client.Client, podName string) {
	var found corev1.Pod
	err := cl.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: common.TestNamespace}, &found)
	s.Error(err, "pod exists, but should not")
	s.True(errors.IsNotFound(err))
}

func (s *BitflowControllerTestHelpers) assertPodsForStep(cl client.Client, stepName string, podCount int) {
	var list corev1.PodList
	err := cl.List(context.TODO(), &client.ListOptions{}, &list)
	s.NoError(err)

	count := 0
	for _, pod := range list.Items {
		if strings.Index(pod.Name, stepName) == 0 && !common.IsBeingDeleted(&pod) {
			count++
		}
	}
	s.Equal(podCount, count, "Expected a different number of pods for step %v", stepName)
}

func (s *BitflowControllerTestHelpers) assertOutputSources(cl client.Client, count int) {
	var sourceList bitflowv1.BitflowSourceList
	err := cl.List(context.TODO(), &client.ListOptions{}, &sourceList)
	s.NoError(err)
	s.NotEmpty(sourceList.Items)

	counter := 0
	for _, source := range sourceList.Items {
		if strings.Index(source.Name, STEP_OUTPUT_PREFIX) == 0 {
			counter++
			_, err := url.Parse(source.Spec.URL)
			s.NoError(err)
		}
	}
	s.Equal(count, counter, "Expected to find output sources")
}

func (s *BitflowControllerTestHelpers) assertRespawningPods(r *BitflowReconciler, count int) {
	s.Len(r.pods.ListRespawningPods(), count, "Wrong number of respawning pods")
}

func (s *BitflowControllerTestHelpers) assertNumberOfPodsForNode(cl client.Client, nodeName string, expectedNumberOfPods int) {
	var list corev1.PodList
	err := cl.List(context.TODO(), &client.ListOptions{}, &list)
	s.NoError(err)

	actualNumberOfPods, err := s.getNumberOfPodsForNode(cl, nodeName)
	s.NoError(err)
	s.Equal(expectedNumberOfPods, actualNumberOfPods)
}

func (s *BitflowControllerTestHelpers) getNumberOfPodsForNode(cli client.Client, nodeName string) (int, error) {
	var podList corev1.PodList
	err := cli.List(context.TODO(), &client.ListOptions{}, &podList)

	if err != nil {
		return 0, err
	}

	count := 0
	for _, pod := range podList.Items {
		if common.GetTargetNode(&pod) == nodeName {
			count++
		}
	}

	return count, nil
}
