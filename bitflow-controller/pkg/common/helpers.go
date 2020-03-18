package common

import (
	corev1 "k8s.io/api/core/v1"
)

func IsBeingDeleted(pod *corev1.Pod) bool {
	return pod.DeletionTimestamp != nil
}

func GetNodeName(pod *corev1.Pod) string {
	if pod.Spec.NodeName != "" {
		return pod.Spec.NodeName
	}

	// TODO commented out, because this is actually a dependency to the 'scheduler' package.
	// TODO Check, in what cases pod.Spec.NodeName is not set
	/*
		// Avoid panic, check the entire object path for nil values and empty slices
		a := pod.Spec.Affinity
		if a != nil && a.NodeAffinity != nil && a.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
			t := a.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
			if len(t) > 0 && len(t[0].MatchExpressions) > 0 && len(t[0].MatchExpressions[0].Values) > 0 {
				return t[0].MatchExpressions[0].Values[0]
			}
		}
	*/

	return ""
}
