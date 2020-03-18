package bitflow

import (
	"crypto/sha1"
	"fmt"
	"regexp"
	"strings"
)

const (
	// TODO make these prefix/suffixes configurable
	STEP_OUTPUT_PREFIX = "output"
	STEP_POD_SUFFIX    = "pod"
	HASH_SUFFIX_LENGTH = 8
	POD_SEPARATOR      = "-"
	SOURCE_SEPARATOR   = "."

	nonDns1123CharReplacement  = "-"
	dns1123StartEndReplacement = "a" // DNS names must start and end with lower-case letters, so there is no real valid replacement for invalid characters...
)

var (
	nonDns1123Chars = regexp.MustCompile("[^-a-z0-9.]") // a dot is added because it does not seem forbidden in Kubernetes, although it is not allows in DNS 1123
	dns1123StartEnd = regexp.MustCompile("[a-z0-9]")
)

func ConstructReproduciblePodName(stepName, sourceName string) string {
	return HashName(ConstructSingletonPodName(stepName)+POD_SEPARATOR, stepName, sourceName)
}

func ConstructSingletonPodName(stepName string) string {
	return CleanDnsName(stepName + POD_SEPARATOR + STEP_POD_SUFFIX)
}

func ConstructSourceName(podName, sourceName string) string {
	return CleanDnsName(STEP_OUTPUT_PREFIX + SOURCE_SEPARATOR + podName + SOURCE_SEPARATOR + sourceName)
}

func HashName(prefix string, hashInputStrings ...string) string {
	hashInputString := strings.Join(hashInputStrings, "--")
	hashString := fmt.Sprintf("%x", sha1.Sum([]byte(hashInputString)))
	return CleanDnsName(prefix + hashString[:HASH_SUFFIX_LENGTH])
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
