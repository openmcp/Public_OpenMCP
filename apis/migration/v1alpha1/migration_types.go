/*
Copyright 2018 The Multicluster-Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MigrationSpec defines the desired state of Migration
type MigrationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	MigrationServiceSources []MigrationServiceSource `json:"MigrationServiceSource"`
}

type MigrationServiceSource struct {
	// container service spec
	MigrationSources []MigrationSource `json:"MigrationSource"`
	ServiceName      string            `json:"ServiceName"`
	TargetCluster    string            `json:"TargetCluster"`
	SourceCluster    string            `json:"SourceCluster"`
	NameSpace        string            `json:"NameSpace"`
}
type MigrationSource struct {
	// Migration source
	ResourceType string `json:"ResourceType"`
	ResourceName string `json:"ResourceName"`
}

// MigrationStatus defines the observed state of Migration
type MigrationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	MigrationStatus bool   `json:"MigrationStatus"`
	Reason          string `json:"Reason"`
	ElapsedTime     string `json:"ElapsedTime,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Migration is the Schema for the Migrations API
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Created time stamp"
// +kubebuilder:printcolumn:name="IsSuccess",type="boolean",JSONPath=".status.MigrationStatus",description="-"
// +kubebuilder:printcolumn:name="SourceCluster",type="string",JSONPath=".spec.MigrationServiceSource[*].SourceCluster",description="-"
// +kubebuilder:printcolumn:name="TargetCluster",type="string",JSONPath=".spec.MigrationServiceSource[*].TargetCluster",description="-"
// +kubebuilder:printcolumn:name="ServiceName",type="string",JSONPath=".spec.MigrationServiceSource[*].ServiceName",description="-"
// +kubebuilder:printcolumn:name="NameSpace",type="string",JSONPath=".spec.MigrationServiceSource[*].NameSpace",description="-"
// +kubebuilder:printcolumn:name="REASON",type="string",JSONPath=".status.Reason",description="-"
// +kubebuilder:printcolumn:name="ElapsedTime",type="string",JSONPath=".status.ElapsedTime",description="ElapsedTime"
type Migration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MigrationSpec   `json:"spec,omitempty"`
	Status MigrationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MigrationList contains a list of Migration
type MigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Migration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Migration{}, &MigrationList{})
}
