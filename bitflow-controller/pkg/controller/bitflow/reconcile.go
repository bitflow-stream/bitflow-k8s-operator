package bitflow

import (
	"context"
	"fmt"
	"time"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TODO if possible, avoid this hack with the fake step name. It is necessary to execute a regularly recurring reconcile.
const STATE_VALIDATION_FAKE_STEP_NAME = "bitflow-validate-state-identifier-fake-step"

type RequeueError struct {
	err error
}

func NewRequeueError(err error) error {
	return &RequeueError{err}
}

func (err *RequeueError) Error() string {
	return err.err.Error()
}

func isRequeueError(err error) bool {
	_, ok := err.(*RequeueError)
	return ok
}

// Reconcile checks all cluster resources associated with a BitflowStep object and makes sure that the cluster corresponds
// to the desired state.
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *BitflowReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	if req.Namespace != r.namespace {
		log.WithFields(log.Fields{"namespace": req.Namespace, "step": req.Name}).Warnf("Ignoring reconcile request for foreign namespace")
		return reconcile.Result{}, nil
	}

	if req.Name == STATE_VALIDATION_FAKE_STEP_NAME {
		return r.recurringReconcileNodeResources(), nil
	}

	start := time.Now()
	err := r.doReconcile(req.Name)
	r.statistic.PodsReconciled(time.Now().Sub(start))

	if err != nil {
		r.statistic.ErrorOccurred()
		if !isRequeueError(err) {
			log.WithField("step", req.Name).Errorf("Error reconciling step: %v", err)
			err = nil
		}
	}
	return reconcile.Result{}, err
}

func (r *BitflowReconciler) doReconcile(stepName string) error {
	logger := log.WithField("step", stepName)
	logger.Debugf("Reconciling step")
	r.respawning.Debug()

	// Load and validate the step
	step, err := common.GetStep(r.client, stepName, r.namespace)
	if err != nil {
		err := r.handleStepError(stepName, err, logger)
		r.reconcileNodeResources()
		return err
	}
	if isStepValid := r.validateStep(step); !isStepValid {
		step.Log().Debugln("Invalid step. Validation error:", step.Status.ValidationError)
		return nil
	}

	// Reconcile all pods for the step
	matchedSources, err := r.listMatchingSources(step)
	if err != nil {
		step.Log().Errorf("Failed to query matching sources for reconciling: %v", err)
	}
	stepType := step.Type()
	var pods []*corev1.Pod
	if stepType == bitflowv1.StepTypeSingleton {
		pods, err = r.GetSingletonPod(step, nil)
	} else if stepType == bitflowv1.StepTypeOneToOne {
		pods, err = r.GetOneToOnePods(step, matchedSources)
	} else if step.Type() == bitflowv1.StepTypeAllToOne {
		pods, err = r.GetAllToOnePod(step, matchedSources)
	} else {
		// Unknown step type, should have been detected in validateStep()
		return fmt.Errorf("Cannot handle unknown step type: %v", stepType)
	}
	if err != nil {
		return err
	}

	// After adding/removing pods, validate resources
	r.reconcileNodeResources()

	// Reconcile the output sources for the step
	return r.reconcileOutputSources(step, pods, matchedSources)
}

func (r *BitflowReconciler) handleStepError(stepName string, err error, logger *log.Entry) error {
	if errors.IsNotFound(err) {
		r.cleanupStep(stepName, "Step deleted")
		return nil
	}
	logger.Errorln("Error fetching step:", err)
	return err
}

func (r *BitflowReconciler) cleanupStep(stepName, reason string) {
	logger := log.WithField("step", stepName)
	logger.Debugln("Cleaning up step:", reason)

	r.cleanupOutputSourcesForStep(stepName, logger)
	podsDeleted := r.cleanupPodsForStep(stepName, logger, "cleaning up step")

	r.respawning.DeletePodsWithLabel(bitflowv1.LabelStepName, stepName)
	if podsDeleted {
		r.reconcileNodeResources()
	}
}

func (r *BitflowReconciler) cleanupOutputSourcesForStep(stepName string, logger *log.Entry) {
	outputs, err := r.listOutputSources(stepName)
	if err != nil {
		logger.Errorln("Failed to query output sources for step:", err)
	} else {
		for _, source := range outputs {
			err = r.client.Delete(context.TODO(), source)
			if err != nil {
				source.Log().Errorln("Failed to delete source:", err)
			}
		}
	}
}

func (r *BitflowReconciler) cleanupPodsForStep(stepName string, logger *log.Entry, reason string, skipPods ...*corev1.Pod) (podsDeleted bool) {
	pods, err := r.listPodsForStep(stepName)
	if err != nil {
		logger.Errorf("Failed to list pods (%v): %v", reason, err)
		return
	}

	skipPodNames := make(map[string]bool, len(skipPods))
	for _, skipPod := range skipPods {
		skipPodNames[skipPod.Name] = true
	}

	for _, pod := range pods.Items {
		if !skipPodNames[pod.Name] {
			podsDeleted = podsDeleted || r.deletePod(&pod, logger, reason)
		}
	}
	return
}

func (r *BitflowReconciler) deletePod(pod *corev1.Pod, logger *log.Entry, reason string) bool {
	if !common.IsBeingDeleted(pod) {
		logger = logger.WithField("pod", pod.Name)
		logger.Infof("Deleting pod (%v)", reason)
		if err := r.client.Delete(context.TODO(), pod); err != nil {
			logger.Errorf("Failed to delete pod (%v): %v", reason, err)
		} else {
			return true
		}
	}
	return false
}

func (r *BitflowReconciler) selectorForStep(stepName string) labels.Selector {
	selector := labels.Set{bitflowv1.LabelStepName: stepName}
	for k, v := range r.idLabels {
		selector[k] = v
	}
	return labels.SelectorFromSet(selector)
}

func (r *BitflowReconciler) listPodsForStep(stepName string) (*corev1.PodList, error) {
	selector := r.selectorForStep(stepName)
	var allPods corev1.PodList
	err := r.client.List(context.TODO(), &client.ListOptions{Namespace: r.namespace, LabelSelector: selector}, &allPods)
	if err != nil {
		err = fmt.Errorf("Failed to query matching pods: %v", err)
	}
	return &allPods, err
}

func (r *BitflowReconciler) listOutputSources(stepName string) ([]*bitflowv1.BitflowSource, error) {
	return common.GetSelectedSources(r.client, r.namespace, r.selectorForStep(stepName))
}

func (r *BitflowReconciler) listMatchingSources(step *bitflowv1.BitflowStep) ([]*bitflowv1.BitflowSource, error) {
	sourceList, err := common.GetSources(r.client, r.namespace)
	if err != nil {
		return nil, err
	}

	// TODO instead of loading ALL sources and filtering them manually, construct a selector from step.Spec.Ingest
	var matchedSources []*bitflowv1.BitflowSource
	for _, source := range sourceList {
		if step.Matches(source.Labels) {
			matchedSources = append(matchedSources, source)
		}
	}
	return matchedSources, nil
}

func (r *BitflowReconciler) listMatchingSteps(source *bitflowv1.BitflowSource) ([]*bitflowv1.BitflowStep, error) {
	allSteps, err := common.GetSteps(r.client, r.namespace)
	if err != nil {
		return nil, err
	}

	// TODO instead of loading ALL sources and filtering them manually, construct a selector from step.Spec.Ingest
	var matchedSteps []*bitflowv1.BitflowStep
	for _, step := range allSteps {
		if step.Matches(source.Labels) {
			matchedSteps = append(matchedSteps, step)
		}
	}
	return matchedSteps, nil
}

func (r *BitflowReconciler) deleteObject(obj runtime.Object, errMsg string, fmtParams ...interface{}) {
	if err := r.client.Delete(context.TODO(), obj); err != nil {
		level := log.ErrorLevel
		if errors.IsNotFound(err) {
			level = log.DebugLevel
		}
		log.StandardLogger().Logf(level, errMsg+": %v", append(fmtParams, err)...)
	}
}
