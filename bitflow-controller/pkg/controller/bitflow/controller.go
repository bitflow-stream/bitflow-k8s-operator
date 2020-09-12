package bitflow

import (
	"fmt"
	"sync"
	"time"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/config"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	kubeSource "sigs.k8s.io/controller-runtime/pkg/source"
)

// TODO make handling of namespaces consistent... All util-objects have their own namespaces configurations.

// Add creates a new Bitflow Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, watchNamespace string) error {
	return startReconciler(mgr, watchNamespace)
}

// Blank assignment to verify that BitflowReconciler implements reconcile.Reconciler
var _ reconcile.Reconciler = &BitflowReconciler{}

type BitflowReconciler struct {
	// This client is a split client that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	cache  cache.Cache

	// Static configuration through env vars
	namespace string
	ownPodIP  string
	apiPort   int
	idLabels  map[string]string

	pods      *ManagedPods
	config    *config.Config
	statistic *ReconcileStatistics

	recurringReconcileStarted sync.Once
	lastSpawnRoutine          time.Time
}

func startReconciler(mgr manager.Manager, watchNamespace string) error {
	params, err := readControllerEnvVars()
	if err != nil {
		return err
	}

	reconciler := &BitflowReconciler{
		client:    mgr.GetClient(),
		scheme:    mgr.GetScheme(), // include the cache for sync-function
		cache:     mgr.GetCache(),
		pods:      NewManagedPods(),
		namespace: watchNamespace,
		ownPodIP:  params.ownPodIP,
		apiPort:   params.apiPort,
		idLabels:  params.controllerIdLabels,
	}
	if params.recordStatistics {
		reconciler.statistic = new(ReconcileStatistics)
	}

	// Initialize various parts of the operator
	reconciler.config = config.NewConfig(mgr.GetClient(), reconciler.namespace, params.configMapName)
	reconciler.startRestApi(fmt.Sprintf(":%v", reconciler.apiPort))

	// Test if Kubernetes connection works
	kubeClient, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}
	configMap, err := kubeClient.CoreV1().ConfigMaps(watchNamespace).Get(params.configMapName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for key, value := range configMap.Data {
		log.Debugf("Entry in ConfigMap: %v = %v", key, value)
	}

	// Start watching relevant Kubernetes objects
	return reconciler.startWatchers(mgr, params)
}

func (r *BitflowReconciler) startWatchers(mgr manager.Manager, params ControllerParameters) error {
	err := mgr.GetFieldIndexer().IndexField(&corev1.Pod{}, "spec.nodeName", func(o runtime.Object) []string {
		return []string{common.GetTargetNode(o.(*corev1.Pod))}
	})
	if err != nil {
		return err
	}

	// Create a new controller
	c, err := controller.New(params.operatorName, mgr, controller.Options{Reconciler: r, MaxConcurrentReconciles: params.concurrentReconcile})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource: BitflowStep
	err = c.Watch(&kubeSource.Kind{Type: &bitflowv1.BitflowStep{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	err = r.addPodWatcher(c)
	if err != nil {
		return err
	}
	err = r.addSourceWatcher(c)
	if err != nil {
		return err
	}
	return r.startRecurringReconcile(c)
}

func (r *BitflowReconciler) addPodWatcher(ctl controller.Controller) error {
	var pred predicate.Funcs
	pred.UpdateFunc = func(e event.UpdateEvent) bool {
		// Only reconcile update-events if pod has an IP assigned. Skip reconciles when pod is spawning and frequently changes state.
		if pod, ok := e.ObjectNew.(*corev1.Pod); ok {
			if pod.Status.PodIP != "" {
				return true
			}
			log.Debugln("(watch-predicate) Ignoring pod without IP:", pod.Name)
		}
		return false
	}
	return ctl.Watch(&kubeSource.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(sourceObj handler.MapObject) (res []reconcile.Request) {
			if pod, ok := sourceObj.Object.(*corev1.Pod); ok {
				if !r.IsControlledObject(pod) {
					log.Debugf("(watch) Ignoring foreign pod %v (missing labels %v)", pod.Name, r.idLabels)
					return
				}

				step := pod.Labels[bitflowv1.LabelStepName]
				if step != "" {
					res = []reconcile.Request{{
						NamespacedName: types.NamespacedName{
							Name:      step,
							Namespace: sourceObj.Meta.GetNamespace(),
						}}}
				} else {
					log.Warnf("Controlled pod '%v' does not contain label '%v'. Deleting pod...", pod.Name, bitflowv1.LabelStepName)
					r.deleteObject(pod, "Failed to delete controlled pod '%v' with missing '%v' label", pod.Name, bitflowv1.LabelStepName)
				}
			}
			return
		},
		),
	}, pred)
}

func (r *BitflowReconciler) addSourceWatcher(ctl controller.Controller) error {
	return ctl.Watch(&kubeSource.Kind{Type: &bitflowv1.BitflowSource{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(sourceObj handler.MapObject) (res []reconcile.Request) {
			source, ok := sourceObj.Object.(*bitflowv1.BitflowSource)
			if !ok || !r.IsWatchedObject(source) {
				log.Debugf("(watch) Ignoring foreign source %v from namespace %v", sourceObj.Meta.GetName(), sourceObj.Meta.GetNamespace())
				return
			}

			r.validateSource(source)

			if r.IsControlledObject(source) {
				// It is an output source created and controlled by us, check the required label
				stepName := source.Labels[bitflowv1.LabelStepName]
				if stepName != "" {
					res = append(res, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Name:      stepName,
							Namespace: source.Namespace,
						}})
				} else {
					log.Warnf("Bogus controlled source '%v' does not contain label '%v'. Deleting source...", source.Name, bitflowv1.LabelStepName)
					r.deleteObject(source, "Failed to delete controlled source '%v' with missing '%v' label", source.Name, bitflowv1.LabelStepName)
				}
			}

			// Reconcile the steps that match this source
			if len(source.Labels) != 0 {
				matchingSteps, err := r.listMatchingSteps(source)
				if err != nil {
					log.Errorf("Failed to load steps that match %v: %v", source, err)
				} else {
					for _, step := range matchingSteps {
						res = append(res, reconcile.Request{
							NamespacedName: types.NamespacedName{
								Name:      step.Name,
								Namespace: step.Namespace,
							}})
					}
				}
			}
			return
		},
		),
	})
}

func (r *BitflowReconciler) startRecurringReconcile(ctl controller.Controller) error {
	// TODO the only purpose of this is to trigger a single reconcile.Request, that in turn will be re-queued in regular intervals
	// Find a more suitable mechanism to achieve a regular invocation of the reconcile method.
	return ctl.Watch(&kubeSource.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(sourceObj handler.MapObject) (requests []reconcile.Request) {
			r.recurringReconcileStarted.Do(func() {
				requests = []reconcile.Request{{
					NamespacedName: types.NamespacedName{
						Name:      ReconcileLoopFakeStepName,
						Namespace: r.namespace,
					}}}
			})
			return
		},
		),
	})
}

func (r *BitflowReconciler) IsWatchedObject(obj metav1.Object) bool {
	return obj.GetNamespace() == r.namespace
}

func (r *BitflowReconciler) IsControlledObject(obj metav1.Object) bool {
	if !r.IsWatchedObject(obj) {
		return false
	}

	objLabels := obj.GetLabels()
	for key, val := range r.idLabels {
		if objLabels[key] != val {
			return false
		}
	}
	return true
}
