package v1

import (
	kommonsv1 "github.com/flanksource/kommons/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:storageversion

type MongoDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MongoDBSpec   `json:"spec,omitempty"`
	Status            MongoDBStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type MongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDB `json:"items"`
}

type MongoDBSpec struct {
	PodResources `json:",inline"`
	Storage      Storage `json:"storage,omitempty"`
	Database     string  `json:"database,omitempty"`
	User         string  `json:"user,omitempty"`
	// Deprecated
	CPU string `json:"cpu,omitempty"`
	// Deprecated
	Memory string `json:"memory,omitempty"`
}

type MongoDBStatus struct {
	Conditions kommonsv1.ConditionList `json:"conditions"`
}
