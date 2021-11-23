package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
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
	VolumeInfos         []VolumeInfo      `json:"VolumeInfo,omitempty"`
}

// VolumeDataSource contains the supported snapshot sources.
type VolumeDataSource struct {
	VolumeSnapshotClassName  string `json:"volumeSnapshotClassName,omitempty"`
	VolumeSnapshotSourceKind string `json:"volumeSnapshotSourceKind,omitempty"`
	VolumeSnapshotSourceName string `json:"volumeSnapshotSourceName,omitempty"`
	VolumeSnapshotKey        string `json:"volumeSnapshotKey"`
}

type VolumeInfo struct {
	VolumeSnapshotKey  string `json:"volumeSnapshotKey"`
	VolumeSnapshotDate string `json:"volumeSnapshotDate,omitempty"`
	VolumeSnapshotSize string `json:"volumeSnapshotSize,omitempty"`
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

	// LastSuccessDate indicate the time to get snapshot last time
	//LastSuccessDate metav1.Time `json:"lastSuccessDate,omitempty"`
	ElapsedTime string `json:"elapsedTime,omitempty"`
	// isVolumeSnapshot
	IsVolumeSnapshot bool `json:"isVolumeSnapshot,omitempty"`

	SnapshotSources []SnapshotSource `json:"snapshotSource,omitempty"`
	// 현재 진행도
	CurrentCount int32 `json:"currentCount,omitempty"`
	// 최대 진행도
	MaxCount int32 `json:"maxCount,omitempty"`
	// 최대 컨디션 카운트. currentCount/maxCount
	ConditionProgress string `json:"progress,omitempty"`

	// Condition 정보 리스트
	//Conditions []MigrationCondition `json:"conditions,omitempty"`
	// Status of the condition, one of True, False, Unknown.
	Status      v1.ConditionStatus `json:"status,omitempty" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	Description string             `json:"description"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Snapshot is the Schema for the snapshots API
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Created time stamp"
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.snapshotSources[*].resourceCluster",description="-"
// +kubebuilder:printcolumn:name="NameSpace",type="string",JSONPath=".spec.snapshotSources[*].resourceNamespace",description="-"
// +kubebuilder:printcolumn:name="GroupSnapshotKey",type="string",JSONPath=".spec.groupSnapshotKey",description="-"
// +kubebuilder:printcolumn:name="IsVolumeSnapshot",type="boolean",JSONPath=".status.isVolumeSnapshot",description="-"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="-"
// +kubebuilder:printcolumn:name="Description",type="string",JSONPath=".status.description",description="-"
// +kubebuilder:printcolumn:name="Progress",type="string",JSONPath=".status.progress",description="-"
// +kubebuilder:printcolumn:name="ElapsedTime",type="string",JSONPath=".status.elapsedTime",description="ElapsedTime"
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
