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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TODO if possible, avoid this hack with the fake step name. It is necessary to execute a regularly recurring reconcile.
const STATE_VALIDATION_FAKE_STEP_NAME = "bitflow-validate-state-identifier-fake-step"

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
	start := time.Now()
	stepName := req.Name
	isStepReconcile := stepName == STATE_VALIDATION_FAKE_STEP_NAME

	// TODO probably reconcile all steps here to avoid unnecessary pod deletions due to delayed updates
	if isStepReconcile {
		logger := log.WithField("step", stepName)
		logger.Debugf("Reconciling step")
		r.pods.Debug()

		if err := r.updatePodStatus(stepName, logger); err != nil {
			// TODO correctly/fully integrate the statistics
			r.statistic.ErrorOccurred()
			logger.Errorf("Error reconciling step: %v", err)
		}
	} else {
		log.Debugln("Recurring node resource reconciliation triggered")
	}

	// After updating the pod status, do the actual pod modifications (start/delete)
	r.reconcileNodeResources()

	// Manage automatically created output data sources
	r.reconcileOutputSources()

	// Make sure the regular automatic reconcile is triggered again
	var result reconcile.Result
	if !isStepReconcile {
		if heartbeat := r.config.GetReconcileHeartbeat(); heartbeat > 0 {
			result.RequeueAfter = heartbeat
		}
	}
	r.statistic.PodsReconciled(time.Now().Sub(start))
	return result, nil
}

func (r *BitflowReconciler) updatePodStatus(stepName string, logger *log.Entry) error {
	// Load and validate the step
	step, err := common.GetStep(r.client, stepName, r.namespace)
	if err != nil {
		err := r.handleStepError(stepName, err, logger)
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
	if stepType == bitflowv1.StepTypeSingleton {
		return r.ReconcileSingletonPod(step, nil)
	} else if stepType == bitflowv1.StepTypeOneToOne {
		return r.ReconcileOneToOnePods(step, matchedSources)
	} else if step.Type() == bitflowv1.StepTypeAllToOne {
		return r.ReconcileAllToOnePod(step, matchedSources)
	} else {
		// Unknown step type, should have been detected in validateStep()
		return fmt.Errorf("Cannot handle unknown step type: %v", stepType)
	}
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

	r.pods.DeletePodsWithLabel(bitflowv1.LabelStepName, stepName)
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
