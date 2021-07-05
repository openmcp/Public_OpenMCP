package v1alpha1

import (
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
	Status bool `json:"Status"`
	// Reason indicates the reason for any backup related failures.
	Reason       string `json:"Reason"`
	ReasonDetail string `json:"ReasonDetail"`
	// LastSuccessDate indicate the time to get snapshot last time
	//LastSuccessDate metav1.Time `json:"lastSuccessDate,omitempty"`
	ElapsedTime string `json:"ElapsedTime,omitempty"`
	// isVolumeSnapshot
	IsVolumeSnapshot      bool                    `json:"isVolumeSnapshot,omitempty"`
	SnapshotRestoreSource []SnapshotRestoreSource `json:"snapshotRestoreSource"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SnapshotRestore is the Schema for the openmcpsnapshotrestores API
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=snapshotrestores,scope=Namespaced
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Created time stamp"
// +kubebuilder:printcolumn:name="IsSuccess",type="boolean",JSONPath=".status.Status",description="-"
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.snapshotRestoreSources[*].resourceCluster",description="-"
// +kubebuilder:printcolumn:name="NameSpace",type="string",JSONPath=".spec.snapshotRestoreSources[*].resourceNamespace",description="-"
// +kubebuilder:printcolumn:name="SnapshotKey",type="string",JSONPath=".spec.groupSnapshotKey",description="-"
// +kubebuilder:printcolumn:name="IsGroupSnapshot",type="boolean",JSONPath=".spec.isGroupSnapshot",description="-"
// +kubebuilder:printcolumn:name="REASON",type="string",JSONPath=".status.Reason",description="-"
// +kubebuilder:printcolumn:name="ElapsedTime",type="string",JSONPath=".status.ElapsedTime",description="ElapsedTime"
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
