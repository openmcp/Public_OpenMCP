// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the keti v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=keti.example.com
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	//	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	GroupName = "keti.example.com"
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
		&OpenMCPPolicyEngine{},
		&OpenMCPPolicyEngineList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}