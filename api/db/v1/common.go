package v1

import v1 "k8s.io/api/core/v1"

type Storage struct {
	// Storage class to use. If not set default will be used
	StorageClass string `yaml:"storageClass,omitempty" json:"storageClass,omitempty"`
	// Size. Required if persistence is enabled
	Size string `yaml:"size,omitempty" json:"size,omitempty"`
}

type PodResources struct {
	Replicas int `json:"replicas"`
	// +kubebuilder:validation:Optional
	Resources v1.ResourceRequirements `json:"resources"`
}
