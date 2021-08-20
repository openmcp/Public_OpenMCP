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
	SnapshotPolicy   *SnapshotPolicy  `json:"snapshotPolicy,omitempty"`
	SnapshotSources  []SnapshotSource `json:"snapshotSources"`
	GroupSnapshotKey string           `json:"groupSnapshotKey,omitempty"`
}

// SnapshotSource contains the supported snapshot sources.
type SnapshotSource struct {
	ResourceCluster     string            `json:"resourceCluster"`
	ResourceNamespace   string            `json:"resourceNamespace"`
	ResourceType        string            `json:"resourceType"`
	ResourceName        string            `json:"resourceName"`
	ResourceSnapshotKey string            `json:"resourceSnapshotKey,omitempty"`
	VolumeDataSource    *VolumeDataSource `json:"volumeDataSource,omitempty"`
}

// VolumeDataSource contains the supported snapshot sources.
type VolumeDataSource struct {
	VolumeSnapshotClassName  string `json:"volumeSnapshotClassName,omitempty"`
	VolumeSnapshotSourceKind string `json:"volumeSnapshotSourceKind,omitempty"`
	VolumeSnapshotSourceName string `json:"volumeSnapshotSourceName,omitempty"`
	VolumeSnapshotKey        string `json:"volumeSnapshotKey"`
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
	Reason       string `json:"Reason"`
	ReasonDetail string `json:"ReasonDetail,omitempty"`
	// LastSuccessDate indicate the time to get snapshot last time
	//LastSuccessDate metav1.Time `json:"lastSuccessDate,omitempty"`
	ElapsedTime string `json:"ElapsedTime,omitempty"`
	// isVolumeSnapshot
	IsVolumeSnapshot bool             `json:"isVolumeSnapshot,omitempty"`
	SnapshotSources  []SnapshotSource `json:"snapshotSource"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Snapshot is the Schema for the snapshots API
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Created time stamp"
// +kubebuilder:printcolumn:name="IsSuccess",type="boolean",JSONPath=".status.Status",description="-"
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.snapshotSources[*].resourceCluster",description="-"
// +kubebuilder:printcolumn:name="NameSpace",type="string",JSONPath=".spec.snapshotSources[*].resourceNamespace",description="-"
// +kubebuilder:printcolumn:name="GroupSnapshotKey",type="string",JSONPath=".spec.groupSnapshotKey",description="-"
// +kubebuilder:printcolumn:name="IsVolumeSnapshot",type="boolean",JSONPath=".status.isVolumeSnapshot",description="-"
// +kubebuilder:printcolumn:name="REASON",type="string",JSONPath=".status.Reason",description="-"
// +kubebuilder:printcolumn:name="ElapsedTime",type="string",JSONPath=".status.ElapsedTime",description="ElapsedTime"
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
