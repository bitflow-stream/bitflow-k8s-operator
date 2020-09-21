package v1

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/ryanuber/go-glob"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StepsKind             = "BitflowStep"
	StepsResourceSingular = "bitflow-step"
	StepsResource         = StepsResourceSingular + "s"

	MatchCheckWildcard = "wildcard"
	MatchCheckExact    = "exact"
	MatchCheckRegex    = "regex"
	MatchCheckPresent  = "present"
	MatchCheckAbsent   = "absent"

	StepTypeOneToOne  = "one-to-one" // The default, when the field is empty
	StepTypeAllToOne  = "all-to-one"
	StepTypeSingleton = "singleton"
)

func init() {
	SchemeBuilder.Register(&BitflowStep{}, &BitflowStepList{})
}

// BitflowStep is the Schema for the bitflowsteps API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ValidationError",type="string",JSONPath=".status.validationError"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:path=bitflow-steps,shortName=bst;bfst
// +kubebuilder:singular=bitflow-step
type BitflowStep struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BitflowStepSpec         `json:"spec,omitempty"`
	Status BitflowValidationStatus `json:"status,omitempty"`
}

// BitflowStepSpec defines the desired state of BitflowStep
// +k8s:openapi-gen=true
type BitflowStepSpec struct {
	// +kubebuilder:validation:Enum=one-to-one;all-to-one;"";singleton
	Type string `json:"type,omitempty"`

	NodeLabels   map[string][]string `json:"nodeLabels,omitempty"`
	Ingest       []*IngestMatch      `json:"ingest,omitempty"`
	Outputs      []*StepOutput       `json:"outputs,omitempty"`
	MinResources corev1.ResourceList `json:"minResources,omitempty"`
	Template     *v1.Pod             `json:"template"`
}

// +k8s:openapi-gen=true
type StepOutput struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// +kubebuilder:validation:MinLength=1
	URL string `json:"url"`

	Labels map[string]string `json:"labels"`

	parsedUrl *url.URL `json:"-"`
}

func (output *StepOutput) GetOutputSourceURL(pod *corev1.Pod) string {
	parsedUrl := output.GetUrl()
	if parsedUrl == nil {
		return ""
	}
	urlCopy := *parsedUrl
	podIP := pod.Status.PodIP
	urlCopy.Host = net.JoinHostPort(podIP, urlCopy.Port())
	return urlCopy.String()
}

func (output *StepOutput) GetUrl() *url.URL {
	parsed, err := output.ParseUrl()
	if err != nil {
		log.Errorf("Failed to parse URL of output '%v' (%v). Should have been validated before. %v", output.Name, output.URL, err)
	}
	return parsed
}

func (output *StepOutput) ParseUrl() (*url.URL, error) {
	if output.parsedUrl != nil {
		return output.parsedUrl, nil
	}
	parsed, err := url.Parse(output.URL)
	if err == nil {
		output.parsedUrl = parsed
	}
	return parsed, err
}

// +k8s:openapi-gen=true
type IngestMatch struct {
	// +kubebuilder:validation:MinLength=1
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
	// +kubebuilder:validation:Enum=wildcard;exact;regex;present;"";absent
	Check string `json:"check,omitempty"`

	// These fields are initialized in Init(), if Regex == true
	compiledKey   *regexp.Regexp `json:"-"`
	compiledValue *regexp.Regexp `json:"-"`
}

func (s *IngestMatch) Matches(labels map[string]string) bool {
	if s.Check == MatchCheckPresent || s.Check == MatchCheckAbsent {
		_, present := labels[s.Key]
		return present == (s.Check == MatchCheckPresent)
	}
	for key, val := range labels {
		if s.matches(key, val) {
			return true
		}
	}
	return false
}

func (s *IngestMatch) matches(key, val string) bool {
	switch s.Check {
	case "", MatchCheckWildcard:
		return glob.Glob(s.Key, key) && glob.Glob(s.Value, val)
	case MatchCheckExact:
		return key == s.Key && val == s.Value
	case MatchCheckRegex:
		if s.compiledKey == nil || s.compiledValue == nil {
			// This indicates a bug, since Init() should have compiled the regexes. Try to compile now.
			if err := s.compileRegexes(); err != nil {
				log.Errorln(err, "Error while compiling Regex")
				return false
			}
		}
		return s.compiledKey.MatchString(key) && s.compiledValue.MatchString(val)
	default:
		return false // The default case should not occur due to previous checks
	}
}

func (s *IngestMatch) compileRegexes() error {
	var err error
	s.compiledKey, err = regexp.Compile(s.Key)
	if err != nil {
		return fmt.Errorf("key failed to compile regex: %v", err)
	}
	s.compiledValue, err = regexp.Compile(s.Value)
	if err != nil {
		return fmt.Errorf("value failed to compile regex: %v", err)
	}
	return nil
}

// BitflowStepList contains a list of BitflowSteps
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BitflowStepList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BitflowStep `json:"items"`
}

func (l *BitflowStepList) GetItems() []*BitflowStep {
	res := make([]*BitflowStep, len(l.Items))
	for i, item := range l.Items {
		itemCopy := item
		res[i] = &itemCopy
	}
	return res
}

func (s *BitflowStep) Type() string {
	if s.Spec.Type == "" {
		return StepTypeOneToOne
	}
	return s.Spec.Type
}

func (s *BitflowStep) String() string {
	return fmt.Sprintf("%v step '%v'", s.Type(), s.Name)
}

func (s *BitflowStep) Matches(sourceLabels map[string]string) bool {
	return !s.IsRecursive(sourceLabels) && s.MatchesLabels(sourceLabels)
}

func (s *BitflowStep) MatchesLabels(sourceLabels map[string]string) bool {
	for _, match := range s.Spec.Ingest {
		if !match.Matches(sourceLabels) {
			return false
		}
	}
	return true
}

func (s *BitflowStep) IsRecursive(sourceLabels map[string]string) bool {
	for key, val := range sourceLabels {
		if strings.HasPrefix(key, PipelinePathLabelPrefix) && val == s.Name {
			return true
		}
	}
	return false
}

func (s *BitflowStep) Log() *log.Entry {
	return s.LogFields(log.WithFields(nil))
}

func (s *BitflowStep) LogFields(entry *log.Entry) *log.Entry {
	return entry.WithFields(log.Fields{
		"step": s.Name,
		"type": s.Type(),
	})
}

func (s *BitflowStep) Validate() {
	s.Status.ValidationError = ""
	if s.Spec.Type != "" && s.Spec.Type != StepTypeOneToOne && s.Spec.Type != StepTypeAllToOne && s.Spec.Type != StepTypeSingleton {
		s.Status.ValidationError += fmt.Sprintf("# %v '%v': Invalid value for .spec.type: %v", StepsKind, s.Name, s.Spec.Type)
	}
	if s.Spec.Type != StepTypeSingleton && len(s.Spec.Ingest) == 0 {
		s.Status.ValidationError += fmt.Sprintf("# %v '%v': .spec.ingest must not be empty", StepsKind, s.Name)
	}
	if s.Spec.Type == StepTypeSingleton && len(s.Spec.Ingest) > 0 {
		s.Status.ValidationError += fmt.Sprintf("# %v '%v': Analysis type %v must have empty .spec.ingest[]", StepsKind, s.Name, StepTypeSingleton)
	}
	if s.Spec.Type != StepTypeSingleton {
		for i, ingest := range s.Spec.Ingest {
			if ingest.Key == "" {
				s.Status.ValidationError += fmt.Sprintf("# %v '%v': .spec.ingest[%v].key is required", StepsKind, s.Name, i)
			}
			if ingest.Check == MatchCheckPresent || ingest.Check == MatchCheckAbsent {
				if ingest.Value != "" {
					s.Status.ValidationError += fmt.Sprintf("# %v '%v': .spec.ingest[%v].value cannot be defined for check=%v", StepsKind, s.Name, i, ingest.Check)
				}
			} else {
				if ingest.Value == "" {
					s.Status.ValidationError += fmt.Sprintf("# %v '%v': .spec.ingest[%v].value is required", StepsKind, s.Name, i)
				}
			}
			if ingest.Check == MatchCheckRegex {
				if err := ingest.compileRegexes(); err != nil {
					s.Status.ValidationError += fmt.Sprintf("# %v '%v': .spec.ingest[%v].%v", StepsKind, s.Name, i, err)
				}
			}
		}
	}
	for _, output := range s.Spec.Outputs {
		parsed, err := output.ParseUrl()
		if err != nil {
			s.Status.ValidationError += fmt.Sprintf("# Error parsing URL '%v': %v", output.URL, err)
		}
		output.parsedUrl = parsed
		if len(output.Labels) == 0 {
			s.Status.ValidationError += fmt.Sprintf("# Output must have at least one new label: %v", output.Labels)
		}
		if output.Name == "" {
			s.Status.ValidationError += fmt.Sprintf("# Output may not have empty name")
		}
	}
}
