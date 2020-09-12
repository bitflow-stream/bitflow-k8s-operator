package common

import (
	"crypto/sha1"
	"fmt"
	"regexp"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	HashSuffixLength = 8

	nonDns1123CharReplacement  = "-"
	dns1123StartEndReplacement = "a" // DNS names must start and end with lower-case letters, so there is no real valid replacement for invalid characters...

	HostnameLabel = "kubernetes.io/hostname"
)

var (
	nonDns1123Chars = regexp.MustCompile("[^-a-z0-9.]") // a dot is added because it does not seem forbidden in Kubernetes, although it is not allows in DNS 1123
	dns1123StartEnd = regexp.MustCompile("[a-z0-9]")
)

func IsBeingDeleted(pod *corev1.Pod) bool {
	return pod.DeletionTimestamp != nil
}

func HashName(prefix string, hashInputStrings ...string) string {
	hashInputString := strings.Join(hashInputStrings, "--")
	hashString := fmt.Sprintf("%x", sha1.Sum([]byte(hashInputString)))
	return CleanDnsName(prefix + hashString[:HashSuffixLength])
}

func CleanDnsName(name string) string {
	// Prepare the given name to conform with DNS-1123, which is enforced by Kubernetes
	if len(name) > 0 {
		name = strings.ToLower(name)
		name = nonDns1123Chars.ReplaceAllString(name, nonDns1123CharReplacement)

		// Check if first and last chars conform to the spec
		if !dns1123StartEnd.MatchString(name[0:1]) {
			name = dns1123StartEndReplacement + name[1:]
		}
		lastChar := len(name) - 1
		if !dns1123StartEnd.MatchString(name[lastChar : lastChar+1]) {
			name = name[:lastChar] + dns1123StartEndReplacement
		}
	}

	// TODO the resulting name might be too long for DNS-1123, which is enforced by Kubernetes and only allows names of 63 characters
	// TODO if the name is too long, strip characters from the original prefix (middle or end), but preserve the UUID
	return name
}

func GetTargetNode(pod *corev1.Pod) string {
	// Avoid panic, check the entire object path for nil values and empty slices
	// This expression must match SetTargetNode()
	a := pod.Spec.Affinity
	if a != nil && a.NodeAffinity != nil && a.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		t := a.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
		if len(t) > 0 && len(t[0].MatchExpressions) > 0 && len(t[0].MatchExpressions[0].Values) > 0 {
			return t[0].MatchExpressions[0].Values[0]
		}
	}

	return ""
}

func SetTargetNode(pod *corev1.Pod, node *corev1.Node) {
	SetTargetNodeName(pod, node.Labels[HostnameLabel])
}

func SetTargetNodeName(pod *corev1.Pod, nodeName string) {
	affinity := getRequiredNodeAffinity()
	setLabelSelector(affinity, HostnameLabel, nodeName)
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
	if in.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		in.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.
			NodeSelectorTerms[0].MatchExpressions[0].Key = key

		in.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.
			NodeSelectorTerms[0].MatchExpressions[0].Values[0] = value
	}
}
