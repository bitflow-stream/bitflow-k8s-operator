// NOTE: Boilerplate only.  Ignore this file.

// Package v1 contains API Schema definitions for the bitflow v1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=bitflow.com
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
)

const (
	Group = "bitflow.com"
	Version = "v1"
	GroupVersion = Group + "/" + Version
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)
