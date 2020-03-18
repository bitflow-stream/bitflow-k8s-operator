// Package v1 contains API Schema definitions for the bitflow v1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=bitflow.com
package v1

// To generate zz_generated.deepcopy.go:
// Execute: operator-sdk generate k8s
// Perform following changes manually (until we know how to adjust the code generation process accordingly):
//  - IngestMatch.DeepCopyInto: remove the explicit copying of the fields compiledKey and compiledValue
//  - StepOutput.DeepCopyInto: replace copying of the parsedUrl field with the following:
//      out.parsedUrl = DeepCopyUrl(in.parsedUrl)

// To generate zz_generated.openapi.go:
// openapi-gen  --logtostderr=true -o "" -i ./pkg/apis/bitflow/v1 -O zz_generated.openapi -p ./pkg/apis/bitflow/v1 -r "-" -h /dev/null

// To generate ../../../deploy/crds/*:
// Executed: operator-sdk generate crds
// Currently, there is a bug in the Kubebuilder part of operator-sdk, so a manual change is necessary:
// Set "spec.names.singular" in both yaml files to "bitflow-source" and "bitflow-step", respectively
