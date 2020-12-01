package migration

import (
	"context"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"

	corev1 "k8s.io/api/core/v1"
)

func migservice(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) {
	omcplog.V(3).Info("service migration")

	targetResource := &corev1.Service{}
	sourceResource := &corev1.Service{}

	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	nameSpace := migSource.nameSpace

	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.V(3).Info("get source cluster error : ", sourceGetErr)
	}
	targetResource = sourceResource
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.V(3).Info("target cluster create error : ", targetErr)
	}

	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.V(3).Info("source cluster delete error : ", sourceErr)
	}
}

func migserviceNotVolume(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) {
	omcplog.V(3).Info("service migration start : " + resource.ResourceName)
	targetResource := &corev1.Service{}
	sourceResource := &corev1.Service{}

	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	nameSpace := migSource.nameSpace

	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.V(3).Info("get source cluster error : ", sourceGetErr)
	} else {
		omcplog.V(3).Info("get source info : " + resource.ResourceName)
	}
	targetResource = sourceResource
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.V(3).Info("target cluster create error : ", targetErr)
	}

	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.V(3).Info("source cluster delete rror : ", sourceErr)
	}
}
