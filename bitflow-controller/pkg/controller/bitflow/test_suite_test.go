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
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/resources"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/scheduler"
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

	idLabels := map[string]string{}
	conf := config.NewConfig(cl, common.TestNamespace, "bitflow-config")
	respawns := common.NewRespawningPods()
	sched := &scheduler.Scheduler{Client: cl, Config: conf, Namespace: common.TestNamespace, IdLabels: idLabels}
	res := &resources.ResourceAssigner{Client: cl, Config: conf, Respawning: respawns, Namespace: "default"}
	return &BitflowReconciler{
		client:          cl,
		scheme:          scheme.Scheme,
		cache:           &MockCache{},
		respawning:      respawns,
		config:          conf,
		scheduler:       sched,
		resourceLimiter: res,
		statistic:       nil,
		namespace:       common.TestNamespace,
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
	s.Error(err, "Pod exists, but should not")
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
	s.Equal(count, r.respawning.Size(), "Wrong number of respawning pods")
}

func (s *BitflowControllerTestHelpers) assertNumberOfPodsForNode(cl client.Client, nodeName string, expectedNumberOfPods int) {
	var list corev1.PodList
	err := cl.List(context.TODO(), &client.ListOptions{}, &list)
	s.NoError(err)

	actualNumberOfPods := 0
	for _, pod := range list.Items {
		if pod.Spec.Affinity != nil {
			if pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0] == nodeName {
				actualNumberOfPods++
			}
		}
	}
	s.Equal(expectedNumberOfPods, actualNumberOfPods)
}
