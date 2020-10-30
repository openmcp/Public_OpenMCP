package apis

import (
	dnsv1alpha1 "openmcp/openmcp/apis/dns/v1alpha1"
	loadbalancingv1alpha1 "openmcp/openmcp/apis/loadbalancing/v1alpha1"
	migrationv1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	policyv1alpha1 "openmcp/openmcp/apis/policy/v1alpha1"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	snapshotv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	syncv1alpha1 "openmcp/openmcp/apis/sync/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, resourcev1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, dnsv1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, loadbalancingv1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, migrationv1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, syncv1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, policyv1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, snapshotv1alpha1.SchemeBuilder.AddToScheme)
}
