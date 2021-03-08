package v1

type Storage struct {
	// Storage class to use. If not set default will be used
	StorageClass string `yaml:"storageClass,omitempty" json:"storageClass,omitempty"`
	// Capacity. Required if persistence is enabled
	Size string `yaml:"size,omitempty" json:"size,omitempty"`
}
