package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SnapshotSpec defines the desired state of Snapshot
type SnapshotSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	SnapshotPolicy  *SnapshotPolicy  `json:"snapshotPolicy,omitempty"`
	SnapshotSources []SnapshotSource `json:"snapshotSources"`
}

// SnapshotSource contains the supported snapshot sources.
type SnapshotSource struct {
	ResourceCluster   string           `json:"resourceCluster,omitempty"`
	ResourceNamespace string           `json:"resourceNamespace,omitempty"`
	ResourceType      string           `json:"resourceType"`
	ResourceName      string           `json:"resourceName"`
	SnapshotKey       string           `json:"SnapshotKey,omitempty"`
	VolumeDataSource  VolumeDataSource `json:"volumeDataSource,omitempty"`
}

// VolumeDataSource contains the supported snapshot sources.
type VolumeDataSource struct {
	VolumeSnapshotClassName  string `json:"volumeSnapshotClassName"`
	VolumeSnapshotSourceKind string `json:"volumeSnapshotSourceKind"`
	VolumeSnapshotSourceName string `json:"volumeSnapshotSourceName"`
	VolumeSnapshotKey        string `json:"volumeSnapshotKey,omitempty"`
}

// SnapshotPolicy defines snapshot policy.
type SnapshotPolicy struct {
	// TimeoutInSecond is the maximal allowed time in second of the entire snapshot process.
	TimeoutInSecond int64 `json:"timeoutInSecond,omitempty"`
	// SnapshotIntervalInSecond is to specify how often operator take snapshot
	// 0 is magic number to indicate one-shot snapshot
	SnapshotIntervalInSecond int64 `json:"snapshotIntervalInSecond,omitempty"`
	// MaxSnapshots is to specify how many snapshots we want to keep
	// 0 is magic number to indicate un-limited snapshots
	MaxSnapshots int `json:"maxSnapshots,omitempty"`
}

// SnapshotStatus defines the observed state of Snapshot
type SnapshotStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Status indicates if the backup has Succeeded.
	Status bool `json:"Status"`
	// Reason indicates the reason for any backup related failures.
	Reason string `json:"Reason,omitempty"`
	// LastSuccessDate indicate the time to get snapshot last time
	LastSuccessDate metav1.Time `json:"lastSuccessDate,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Snapshot is the Schema for the snapshots API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=snapshots,scope=Namespaced
type Snapshot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnapshotSpec   `json:"spec,omitempty"`
	Status SnapshotStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SnapshotList contains a list of Snapshot
type SnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Snapshot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Snapshot{}, &SnapshotList{})
}
