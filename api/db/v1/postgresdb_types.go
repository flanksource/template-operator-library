package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type PostgresDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PostgresDBSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

type PostgresDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresDB `json:"items"`
}

type PostgresDBSpec struct {
	Replicas int            `json:"replicas,omitempty"`
	Storage  PostgresVolume `json:"storage,omitempty"`
	// +kubebuilder:validation:Optional
	Backup     PostgresBackup    `json:"backup,omitempty"`
	CPU        string            `json:"cpu,omitempty"`
	Memory     string            `json:"memory,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

type PostgresVolume struct {
	// Storage class to use. If not set default will be used
	StorageClass string `yaml:"storageClass,omitempty" json:"storageClass,omitempty"`
	// Size. Required if persistence is enabled
	Size string `yaml:"size,omitempty" json:"size,omitempty"`
}

type PostgresBackup struct {
	Bucket   string `json:"bucket,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}

func init() {
	SchemeBuilder.Register(&PostgresDB{}, &PostgresDBList{})
}
