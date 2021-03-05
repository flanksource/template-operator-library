package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type RedisDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RedisDBSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

type RedisDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisDB `json:"items"`
}

type RedisDBSpec struct {
	Replicas string `json:"replicas,omitempty"`
	CPU      string `json:"cpu,omitempty"`
	Memory   string `json:"memory,omitempty"`

	SentinelReplicas string `json:"sentinelReplicas,omitempty"`
	SentinelCPU      string `json:"sentinelCPU,omitempty"`
	SentinelMemory   string `json:"sentinelMemory,omitempty"`

	Storage RedisVolume `json:"storage,omitempty"`
}

type RedisVolume struct {
	// Storage class to use. If not set default will be used
	StorageClass string `yaml:"storageClass,omitempty" json:"storageClass,omitempty"`
	// Capacity. Required if persistence is enabled
	Size string `yaml:"size,omitempty" json:"size,omitempty"`
}

func init() {
	SchemeBuilder.Register(&RedisDB{}, &RedisDBList{})
}
