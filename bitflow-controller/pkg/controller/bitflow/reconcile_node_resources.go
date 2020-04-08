package bitflow

import (
	"context"
	"time"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *BitflowReconciler) recurringReconcileNodeResources() reconcile.Result {
	log.Debugln("Recurring node resource reconciliation triggered")
	r.reconcileNodeResources()
	if heartbeat := r.config.GetValidationHeartbeat(); heartbeat <= 0 {
		return reconcile.Result{}
	} else {
		return reconcile.Result{RequeueAfter: heartbeat}
	}
}

func (r *BitflowReconciler) reconcileNodeResources() {
	now := time.Now()
	period := r.config.GetValidationPeriod()

	// Keep clean, so read and write of lastResourceReconciliation are as close as possible to each other
	last := r.lastResourceReconciliation
	shouldReconcile := last.IsZero() || now.Sub(last) >= period
	if shouldReconcile {
		r.lastResourceReconciliation = now
		log.Debugln("Reconciling node resources...")
		startTimestamp := time.Now()
		r.doReconcileResourcesOnAllNodes()
		r.statistic.NodeResourcesReconciled(time.Now().Sub(startTimestamp))
	}
}

func (r *BitflowReconciler) doReconcileResourcesOnAllNodes() {
	nodes, err := common.RequestReadyNodes(r.client)
	if err != nil {
		log.Errorln("Failed to retrieve all ready nodes:", err)
		return
	}
	for _, node := range nodes.Items {
		logger := log.WithField("node", node.Name)
		pods, err := common.RequestAllPodsOnNode(r.client, node.Name, r.namespace, r.idLabels)
		if err != nil {
			logger.Errorf("Failed to retrieve all Bitflow pods on node: %v", err)
			continue
		}
		resources := r.resourceLimiter.GetCurrentResources(&node)
		for _, pod := range pods {
			deletePod := false
			for _, container := range pod.Spec.Containers {
				logger := logger.WithField("pod", pod.Name).WithField("container", container.Name)
				containerCpu := container.Resources.Limits.Cpu().MilliValue()
				if resources == nil && !container.Resources.Limits.Cpu().IsZero() {
					logger.Infof("Container resource limit was removed (was %v mCPU), restarting pod", containerCpu)
					deletePod = true
					break
				}

				assignedCpu := resources.Cpu().MilliValue()
				logger = logger.WithField("currentCpu", containerCpu).WithField("assignedCpu", assignedCpu)
				logger.Debugf("Reconciling container resources")
				if containerCpu != assignedCpu {
					logger.Infof("Container resources do not match assigned resources, restarting pod")
					deletePod = true
					break
				}
			}
			if deletePod {
				r.deletePodForRestart(pod)
			}
		}
	}
}

func (r *BitflowReconciler) deletePodForRestart(pod *corev1.Pod) {
	if _, exists := r.respawning.IsPodRestarting(pod.Name); exists || common.IsBeingDeleted(pod) {
		return
	}

	gracePeriod := r.config.GetDeleteGracePeriod()
	var delOpt client.DeleteOptionFunc
	if gracePeriod >= 0 {
		delOpt = client.GracePeriodSeconds(int64(gracePeriod.Seconds()))
	}
	logger := log.WithField("pod", pod.Name)
	logger.Infoln("Deleting pod for restart")
	err := r.client.Delete(context.TODO(), pod, delOpt)

	if err != nil {
		logger.Errorf("Pod could not be deleted: %v", err)
	} else {
		newPod := pod.DeepCopy()
		newPod.ResourceVersion = ""
		r.respawning.Add(newPod)
	}
}
