package migration

import (
	"context"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"

	corev1 "k8s.io/api/core/v1"
)

func migpvc(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) error {
	omcplog.V(3).Info("pvc migration")
	targetResource := &corev1.PersistentVolumeClaim{}
	sourceResource := &corev1.PersistentVolumeClaim{}

	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	nameSpace := migSource.nameSpace
	volumePath := migSource.volumePath
	serviceName := migSource.serviceName

	// targetGetErr := targetClient.Get(context.TODO(), targetResource, nameSpace, resource.ResourceName)
	// if targetGetErr != nil {
	// 	omcplog.V(3).Info("get target cluster")
	// }
	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.Error("get source cluster error : ", sourceGetErr)
		return sourceGetErr
	}

	//targetResource = sourceResource
	targetResource = GetLinkSharePvc(sourceResource, volumePath, serviceName)
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.ResourceVersion = ""

	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.Error("target cluster create error : ", targetErr)
		return targetErr
	}

	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.Error("source cluster delete error : ", sourceErr)
		return sourceErr
	}
	return nil
}
