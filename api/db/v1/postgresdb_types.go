package v1

import (
	kommonsv1 "github.com/flanksource/kommons/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type PostgresqlDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PostgresqlDBSpec   `json:"spec,omitempty"`
	Status            PostgresqlDBStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type PostgresqlDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresqlDB `json:"items"`
}

type PostgresqlDBSpec struct {
	PodResources `json:",inline"`
	Storage      Storage `json:"storage,omitempty"`
	// +kubebuilder:validation:Optional
	Backup     PostgresqlBackup  `json:"backup,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
	Replicas   int               `yaml:"replicas,omitempty" json:"replicas,omitempty"`
	// Deprecated
	CPU string `json:"cpu,omitempty"`
	// Deprecated
	Memory string `json:"memory,omitempty"`
}

type PostgresqlDBStatus struct {
	Conditions kommonsv1.ConditionList `json:"conditions"`
}

type PostgresqlBackup struct {
	Bucket    string          `json:"bucket,omitempty"`
	Schedule  string          `json:"schedule,omitempty"`
	Retention BackupRetention `yaml:"retention,omitempty" json:"retention,omitempty"`
}

type BackupRetention struct {
	KeepLast    int `yaml:"keepLast,omitempty" json:"keepLast,omitempty"`
	KeepHourly  int `yaml:"keepHourly,omitempty" json:"keepHourly,omitempty"`
	KeepDaily   int `yaml:"keepDaily,omitempty" json:"keepDaily,omitempty"`
	KeepWeekly  int `yaml:"keepWeekly,omitempty" json:"keepWeekly,omitempty"`
	KeepMonthly int `yaml:"keepMonthly,omitempty" json:"keepMonthly,omitempty"`
	KeepYearly  int `yaml:"keepYearly,omitempty" json:"keepYearly,omitempty"`
}

func init() {
	SchemeBuilder.Register(&PostgresqlDB{}, &PostgresqlDBList{})
}
