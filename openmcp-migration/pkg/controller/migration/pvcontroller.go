package migration

import (
	"context"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"

	corev1 "k8s.io/api/core/v1"
)

func migpv(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) {
	omcplog.V(3).Info("pv migration")
	targetResource := &corev1.PersistentVolume{}
	sourceResource := &corev1.PersistentVolume{}

	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	nameSpace := migSource.nameSpace
	volumePath := migSource.volumePath
	serviceName := migSource.serviceName

	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.V(3).Info("get source cluster error : ", sourceGetErr)
	}

	targetResource = GetLinkSharePv(sourceResource, volumePath, serviceName)
	targetResource.Spec.Capacity = sourceResource.Spec.Capacity
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	targetResource.Spec.ClaimRef = nil

	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.V(3).Info("target cluster create error : ", targetErr)
	}

	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.V(3).Info("source cluster delete error : ", sourceErr)
	}
}
