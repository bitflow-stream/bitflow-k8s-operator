package bitflow

import (
	"context"
	"fmt"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *BitflowReconciler) waitForCacheSync() {
	stopper := make(chan struct{})
	if !r.cache.WaitForCacheSync(stopper) {
		log.Warnln("Could not sync caches")
	}
	close(stopper)
}

func (r *BitflowReconciler) deletePod(pod *corev1.Pod, logger *log.Entry, reason string) bool {
	if !common.IsBeingDeleted(pod) {
		gracePeriod := r.config.GetDeleteGracePeriod()
		var delOpt client.DeleteOptionFunc
		if gracePeriod >= 0 {
			delOpt = client.GracePeriodSeconds(int64(gracePeriod.Seconds()))
		}

		logger = logger.WithField("pod", pod.Name)
		logger.Infof("Deleting pod (%v)", reason)
		if err := r.client.Delete(context.TODO(), pod, delOpt); err != nil {
			logger.Errorf("Failed to delete pod (%v): %v", reason, err)
		} else {
			return true
		}
	}
	return false
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

func (r *BitflowReconciler) genericSelector() labels.Selector {
	return labels.SelectorFromSet(r.idLabels)
}

func (r *BitflowReconciler) selectorForStep(stepName string) labels.Selector {
	selector := make(labels.Set)
	if stepName != "" {
		selector[bitflowv1.LabelStepName] = stepName
	}
	for k, v := range r.idLabels {
		selector[k] = v
	}
	return labels.SelectorFromSet(selector)
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

func (r *BitflowReconciler) listAllSteps() ([]*bitflowv1.BitflowStep, error) {
	return r.listMatchingSteps(nil)
}

func (r *BitflowReconciler) listMatchingSteps(source *bitflowv1.BitflowSource) ([]*bitflowv1.BitflowStep, error) {
	allSteps, err := common.GetSteps(r.client, r.namespace)
	if err != nil {
		return nil, err
	}

	// TODO instead of loading ALL sources and filtering them manually, construct a selector from step.Spec.Ingest
	var matchedSteps []*bitflowv1.BitflowStep
	for _, step := range allSteps {
		if source == nil || step.Matches(source.Labels) {
			matchedSteps = append(matchedSteps, step)
		}
	}
	return matchedSteps, nil
}

func (r *BitflowReconciler) deleteObject(obj runtime.Object, errMsg string, fmtParams ...interface{}) (ok bool) {
	ok = true
	if err := r.client.Delete(context.TODO(), obj); err != nil {
		level := log.ErrorLevel
		if errors.IsNotFound(err) {
			level = log.DebugLevel
		} else {
			ok = false
		}
		log.StandardLogger().Logf(level, errMsg+": %v", append(fmtParams, err)...)
	}
	return
}
