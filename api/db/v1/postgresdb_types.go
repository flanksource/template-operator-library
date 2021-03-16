package v1

import (
	kommonsv1 "github.com/flanksource/kommons/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type PostgresDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PostgresDBSpec   `json:"spec,omitempty"`
	Status            PostgresDBStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type PostgresDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresDB `json:"items"`
}

type PostgresDBSpec struct {
	PodResources `json:",inline"`
	Storage      Storage `json:"storage,omitempty"`
	// +kubebuilder:validation:Optional
	Backup     PostgresBackup    `json:"backup,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

type PostgresDBStatus struct {
	Conditions kommonsv1.ConditionList `json:"conditions"`
}

type PostgresBackup struct {
	Bucket   string `json:"bucket,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}

func init() {
	SchemeBuilder.Register(&PostgresDB{}, &PostgresDBList{})
}
