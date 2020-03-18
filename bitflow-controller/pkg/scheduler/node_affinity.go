package scheduler

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	HostnameLabel = "kubernetes.io/hostname"
)

func SetPodNodeAffinityPreferred(node *corev1.Node, pod *corev1.Pod) {
	affinity := getPreferredNodeAffinity()
	setLabelSelector(affinity, HostnameLabel, node.Labels[HostnameLabel])
	pod.Spec.Affinity = affinity
}

func SetPodNodeAffinityRequired(node *corev1.Node, pod *corev1.Pod) {
	affinity := getRequiredNodeAffinity()
	setLabelSelector(affinity, HostnameLabel, node.Labels[HostnameLabel])
	pod.Spec.Affinity = affinity
}

func getRequiredNodeAffinity() *corev1.Affinity {
	var affinity corev1.Affinity
	var nodeAffinity corev1.NodeAffinity
	var selector corev1.NodeSelector
	selectorTerms := make([]corev1.NodeSelectorTerm, 1, 1)
	term := getNodeSelectorTerm()

	selectorTerms[0] = term
	selector.NodeSelectorTerms = selectorTerms
	nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &selector
	affinity.NodeAffinity = &nodeAffinity
	return &affinity
}

func getPreferredNodeAffinity() *corev1.Affinity {
	var affinity corev1.Affinity
	var nodeAffinity corev1.NodeAffinity
	preference := make([]corev1.PreferredSchedulingTerm, 1, 1)
	var prefTerm corev1.PreferredSchedulingTerm
	selectorTerm := getNodeSelectorTerm()

	prefTerm.Weight = 1
	prefTerm.Preference = selectorTerm
	preference[0] = prefTerm
	nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = preference
	affinity.NodeAffinity = &nodeAffinity
	return &affinity
}

func getNodeSelectorTerm() corev1.NodeSelectorTerm {
	var term corev1.NodeSelectorTerm
	exp := make([]corev1.NodeSelectorRequirement, 1, 1)
	var exp1 corev1.NodeSelectorRequirement
	values := make([]string, 1, 1)

	op := corev1.NodeSelectorOpIn
	exp1.Operator = op
	exp1.Values = values
	exp[0] = exp1
	term.MatchExpressions = exp

	return term
}

func setLabelSelector(in *corev1.Affinity, key, value string) {
	if in.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
		in.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution[0].
			Preference.MatchExpressions[0].Key = key

		in.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution[0].
			Preference.MatchExpressions[0].Values[0] = value
	} else if in.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		in.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.
			NodeSelectorTerms[0].MatchExpressions[0].Key = key

		in.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.
			NodeSelectorTerms[0].MatchExpressions[0].Values[0] = value
	}
}
