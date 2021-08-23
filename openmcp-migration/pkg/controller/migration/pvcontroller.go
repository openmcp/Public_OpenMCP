package migration

import (
	"context"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func migpv(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) error {
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
		omcplog.Error("get source cluster error : ", sourceGetErr)
		return sourceGetErr
	}

	targetResource = GetLinkSharePv(sourceResource, volumePath, serviceName)
	targetResource.Spec.Capacity = sourceResource.Spec.Capacity
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	targetResource.Spec.ClaimRef = nil

	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		if strings.Contains(targetErr.Error(), "already exists") {
			omcplog.V(3).Info("target cluster create error : ", targetErr)
			omcplog.V(3).Info("continue...")
		} else {
			omcplog.Error("target cluster create error : ", targetErr)
			return targetErr
		}
	}

	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.Error("source cluster delete error : ", sourceErr)
		return sourceErr
	}
	return nil
}
