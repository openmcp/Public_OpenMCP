package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SnapshotRestoreSpec defines the desired state of SnapshotRestore
type SnapshotRestoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	SnapshotRestoreSource []SnapshotRestoreSource `json:"snapshotRestoreSource,omitempty"`
	GroupSnapshotKey      string                  `json:"groupSnapshotKey"`
	IsGroupSnapshot       bool                    `json:"isGroupSnapshot,omitempty"`
}

// SnapshotRestoreSource contains the supported SnapshotRestore sources.
type SnapshotRestoreSource struct {
	ResourceCluster     string `json:"resourceCluster"`
	ResourceNamespace   string `json:"resourceNamespace"`
	ResourceType        string `json:"resourceType"`
	ResourceSnapshotKey string `json:"resourceSnapshotKey"`
	VolumeSnapshotKey   string `json:"volumeSnapshotKey,omitempty"`
}

// SnapshotRestoreStatus defines the observed state of SnapshotRestore
type SnapshotRestoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Succeeded indicates if the backup has Succeeded.
	// LastSuccessDate indicate the time to get snapshot last time
	//LastSuccessDate metav1.Time `json:"lastSuccessDate,omitempty"`
	ElapsedTime string `json:"ElapsedTime,omitempty"`
	// isVolumeSnapshot
	IsVolumeSnapshot      bool                    `json:"isVolumeSnapshot,omitempty"`
	SnapshotRestoreSource []SnapshotRestoreSource `json:"snapshotRestoreSource,omitempty"`

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

// SnapshotRestore is the Schema for the openmcpsnapshotrestores API
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=snapshotrestores,scope=Namespaced
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Created time stamp"
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.snapshotRestoreSources[*].resourceCluster",description="-"
// +kubebuilder:printcolumn:name="NameSpace",type="string",JSONPath=".spec.snapshotRestoreSources[*].resourceNamespace",description="-"
// +kubebuilder:printcolumn:name="SnapshotKey",type="string",JSONPath=".spec.groupSnapshotKey",description="-"
// +kubebuilder:printcolumn:name="IsGroupSnapshot",type="boolean",JSONPath=".spec.isGroupSnapshot",description="-"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="-"
// +kubebuilder:printcolumn:name="Description",type="string",JSONPath=".status.description",description="-"
// +kubebuilder:printcolumn:name="Progress",type="string",JSONPath=".status.progress",description="-"
// +kubebuilder:printcolumn:name="ElapsedTime",type="string",JSONPath=".status.elapsedTime",description="ElapsedTime"
type SnapshotRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnapshotRestoreSpec   `json:"spec,omitempty"`
	Status SnapshotRestoreStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SnapshotRestoreList contains a list of SnapshotRestore
type SnapshotRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SnapshotRestore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SnapshotRestore{}, &SnapshotRestoreList{})
}
