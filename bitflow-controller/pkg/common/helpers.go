package common

import (
	"context"
	"crypto/sha1"
	"fmt"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	HashSuffixLength = 8

	nonDns1123CharReplacement  = "-"
	dns1123StartEndReplacement = "a" // DNS names must start and end with lower-case letters, so there is no real valid replacement for invalid characters...
)

var (
	nonDns1123Chars = regexp.MustCompile("[^-a-z0-9.]") // a dot is added because it does not seem forbidden in Kubernetes, although it is not allows in DNS 1123
	dns1123StartEnd = regexp.MustCompile("[a-z0-9]")
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

func GetNumberOfPodsForNode(cli client.Client, nodeName string) (int, error) {
	var podList corev1.PodList
	err := cli.List(context.TODO(), &client.ListOptions{}, &podList)

	if err != nil {
		return 0, err
	}

	count := 0
	for _, pod := range podList.Items {
		if pod.Spec.Affinity != nil {
			if pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0] == nodeName {
				count++
			}
		}
	}

	return count, nil
}
