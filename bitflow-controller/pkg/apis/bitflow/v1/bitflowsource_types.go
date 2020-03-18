package v1

import (
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DataSourcesKind             = "BitflowSource"
	DataSourcesResourceSingular = "bitflow-source"
	DataSourcesResource         = DataSourcesResourceSingular + "s"
)

func init() {
	SchemeBuilder.Register(&BitflowSource{}, &BitflowSourceList{})
}

// BitflowSource is the Schema for the BitflowSource API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".spec.url"
// +kubebuilder:printcolumn:name="ValidationError",type="string",JSONPath=".status.validationError"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:path=bitflow-sources,shortName=bso;bfso,singular=bitflow-source
type BitflowSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BitflowSourceSpec       `json:"spec,omitempty"`
	Status BitflowValidationStatus `json:"status,omitempty"`
}

// BitflowSourceSpec defines the desired state of BitflowSource
// +k8s:openapi-gen=true
type BitflowSourceSpec struct {
	// +kubebuilder:validation:MinLength=1
	URL string `json:"url"`
}

// BitflowSourceList contains a list of BitflowSource
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BitflowSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BitflowSource `json:"items"`
}

func (l *BitflowSourceList) GetItems() []*BitflowSource {
	res := make([]*BitflowSource, len(l.Items))
	for i, item := range l.Items {
		itemCopy := item
		res[i] = &itemCopy
	}
	return res
}

func (s *BitflowSource) String() string {
	return fmt.Sprintf("%v '%v'", DataSourcesKind, s.Name)
}

func (s *BitflowSource) Log() *log.Entry {
	return s.LogFields(log.WithFields(nil))
}

func (s *BitflowSource) LogFields(entry *log.Entry) *log.Entry {
	return entry.WithFields(log.Fields{
		"source": s.Name,
		"url":    s.Spec.URL,
	})
}

func (s *BitflowSource) Validate() {
	s.Status.ValidationError = ""
	if s.Spec.URL == "" || len(s.Labels) == 0 {
		s.Status.ValidationError += fmt.Sprintf("# %v '%v': .spec.url and .metadata.labels[] must not be empty", DataSourcesKind, s.Name)
	}
	_, err := url.Parse(s.Spec.URL)
	if err != nil {
		s.Status.ValidationError += fmt.Sprintf("# Error parsing URL '%v': %v", s.Spec.URL, err)
	}
}

func (s *BitflowSource) EqualSpec(other *BitflowSource) bool {
	if s.Spec.URL != other.Spec.URL {
		return false
	}
	if len(s.Labels) != len(other.Labels) {
		return false
	}
	for key, val := range s.Labels {
		if other.Labels[key] != val {
			return false
		}
	}
	return true
}
