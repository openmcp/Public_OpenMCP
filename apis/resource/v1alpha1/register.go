// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the keti v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=openmcp.k8s.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	GroupName    = "openmcp.k8s.io"
	GroupVersion = "v1alpha1"
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}
)
var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	SchemeBuilder.Register(addKnownTypes)
}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&OpenMCPDeployment{},
		&OpenMCPDeploymentList{},
		&OpenMCPService{},
		&OpenMCPServiceList{},
		&OpenMCPHybridAutoScaler{},
		&OpenMCPHybridAutoScalerList{},
		&OpenMCPIngress{},
		&OpenMCPIngressList{},
		&OpenMCPConfigMap{},
		&OpenMCPConfigMapList{},
		&OpenMCPSecret{},
		&OpenMCPSecretList{},
		&OpenMCPJob{},
		&OpenMCPJobList{},
		&OpenMCPNamespace{},
		&OpenMCPNamespaceList{},
		&OpenMCPVirtualService{},
		&OpenMCPVirtualServiceList{},
		&OpenMCPPersistentVolumeClaim{},
		&OpenMCPPersistentVolumeClaimList{},
		&OpenMCPPersistentVolume{},
		&OpenMCPPersistentVolumeList{},
		&OpenMCPStatefulSet{},
		&OpenMCPStatefulSetList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
