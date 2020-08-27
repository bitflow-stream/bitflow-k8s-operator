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
const ReconcileLoopFakeStepName = "bitflow-trigger-reconcile-fake-step"

// Reconcile checks all cluster resources associated with a BitflowStep object and makes sure that the cluster corresponds
// to the desired state.
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *BitflowReconciler) Reconcile(req reconcile.Request) (result reconcile.Result, err error) {
	logger := log.WithFields(log.Fields{"namespace": req.Namespace, "step": req.Name})
	if req.Namespace != r.namespace {
		logger.Warnf("Ignoring reconcile request for foreign namespace")
		return
	}
	start := time.Now()
	stepName := req.Name
	isStepReconcile := stepName != ReconcileLoopFakeStepName

	if isStepReconcile {
		logger.Debugf("Reconciling step")
		r.updatePodStatus(stepName, logger)
	} else {
		log.Debugln("Auto-reconcile loop triggered")
	}

	// Manage automatically created output data sources
	r.reconcileOutputSources()
	r.statistic.PodsUpdated(time.Now().Sub(start))

	// TODO control, how often and when the schedule/spawn routine is triggered
	// TODO update pods for ALL steps before making modifications to pods?

	// After updating the pod status, do the actual pod modifications (start/delete)
	r.spawnPodsPeriodically()

	// Make sure the regular automatic reconcile is triggered again
	if !isStepReconcile {
		if heartbeat := r.config.GetReconcileHeartbeat(); heartbeat > 0 {
			result.RequeueAfter = heartbeat
		}
	}
	return
}

func (r *BitflowReconciler) updatePodStatus(stepName string, logger *log.Entry) {
	pods, step, err := r.constructPodsForStep(stepName)
	if err != nil && !errors.IsNotFound(err) {
		logger.Errorln("Failed to construct pods for step:", err)
		// Do NOT return here, clean up (delete) pods from failed step
	}

	// Create/update pods for this step
	for pod, inputSources := range pods {
		r.pods.Put(pod, step, inputSources)
	}

	// Delete other (dangling) pods for this step
	existingPodNames := make(map[string]bool, len(pods))
	for existingPod := range pods {
		existingPodNames[existingPod.Name] = true
	}
	r.pods.CleanupStep(stepName, existingPodNames)
}

func (r *BitflowReconciler) constructPodsForStep(stepName string) (map[*corev1.Pod][]*bitflowv1.BitflowSource, *bitflowv1.BitflowStep, error) {
	// Load and validate the step
	step, err := common.GetStep(r.client, stepName, r.namespace)
	if err != nil {
		return nil, nil, err
	}
	logger := step.Log()
	if isStepValid := r.validateStep(step, logger); !isStepValid {
		logger.Debugln("Step validation failed:", step.Status.ValidationError)
		return nil, nil, fmt.Errorf("Validation error in step: %v", step.Status.ValidationError)
	}

	// Construct all pods for the step
	matchedSources, err := r.listMatchingSources(step)
	if err != nil {
		logger.Errorf("Failed to query matching sources for: %v", err)
		// Continue, so that at least singleton pods can still be created, and other pods are cleaned up
	}
	switch stepType := step.Type(); stepType {
	case bitflowv1.StepTypeSingleton:
		return r.constructSingletonPod(step, nil), step, nil
	case bitflowv1.StepTypeOneToOne:
		return r.constructOneToOnePods(step, matchedSources), step, nil
	case bitflowv1.StepTypeAllToOne:
		return r.constructAllToOnePod(step, matchedSources), step, nil
	default:
		// Unknown step type, should have been detected in validateStep()
		return nil, nil, fmt.Errorf("Cannot handle unknown step type: %v", stepType)
	}
}

func (r *BitflowReconciler) validateStep(step *bitflowv1.BitflowStep, logger *log.Entry) bool {
	validationMsg := step.Status.ValidationError
	step.Validate()
	if validationMsg != step.Status.ValidationError {
		err := r.client.Status().Update(context.TODO(), step)
		if err != nil {
			logger.Errorf("Failed to update validation error status: %v", err)
		}
	}
	return step.Status.ValidationError == ""
}

func (r *BitflowReconciler) spawnPodsPeriodically() {
	now := time.Now()
	period := r.config.GetSpawnPeriod()

	// Keep clean, so read and write of lastSpawnRoutine are as close as possible to each other
	last := r.lastSpawnRoutine
	if last.IsZero() || now.Sub(last) >= period {
		r.lastSpawnRoutine = now
		log.Debugln("Scheduling/deleting pods...")
		startTimestamp := time.Now()
		r.spawnPods()
		r.statistic.PodsSpawned(time.Now().Sub(startTimestamp))
	}
}
