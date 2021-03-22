package v1

import (
	kommonsv1 "github.com/flanksource/kommons/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type RedisDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RedisDBSpec   `json:"spec,omitempty"`
	Status            RedisDBStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type RedisDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisDB `json:"items"`
}

type RedisDBSpec struct {
	Redis    RedisSpec    `json:"redis"`
	Sentinel SentinelSpec `json:"sentinel"`
	Storage  Storage      `json:"storage,omitempty"`
}

type RedisSpec struct {
	PodResources `json:",inline"`
}

type RedisDBStatus struct {
	Conditions kommonsv1.ConditionList `json:"conditions"`
}

type SentinelSpec struct {
	PodResources `json:",inline"`
}

func init() {
	SchemeBuilder.Register(&RedisDB{}, &RedisDBList{})
}
