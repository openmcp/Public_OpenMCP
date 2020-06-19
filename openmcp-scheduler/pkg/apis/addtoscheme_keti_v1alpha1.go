package apis

import (
	"openmcpscheduler/pkg/apis/keti/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
	//AddToSchemes = append(AddToSchemes, vpav1beta2.SchemeBuilder.AddToScheme)
	//AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder2.AddToScheme)
}
