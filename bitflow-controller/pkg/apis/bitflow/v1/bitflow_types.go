package v1

// BitflowValidationStatus defines the validation status of a Bitflow object
// +k8s:openapi-gen=true
type BitflowValidationStatus struct {
	ValidationError string `json:"validationError"`
}
