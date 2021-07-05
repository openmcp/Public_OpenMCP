package migration

import (
	"context"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"

	corev1 "k8s.io/api/core/v1"
)

func migservice(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) error {
	omcplog.V(3).Info("=== service migration start")
	targetResource := &corev1.Service{}
	sourceResource := &corev1.Service{}

	omcplog.V(3).Info("SourceClient init...")
	omcplog.V(3).Info("TargetClient init...")
	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	omcplog.V(3).Info("SourceClient init complete")
	omcplog.V(3).Info("TargetClient init complete")

	omcplog.V(3).Info("Make deployment resource info...")
	nameSpace := migSource.nameSpace

	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.Error("get source cluster error : ", sourceGetErr)
		return sourceGetErr
	} else {
		omcplog.V(3).Info("get source info : " + resource.ResourceName)
	}
	targetResource = sourceResource
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	omcplog.V(3).Info("Make deployment resource info complete")

	omcplog.V(3).Info("Create for target cluster")
	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.Error("target cluster create error : ", targetErr)
		return targetErr
	}
	omcplog.V(3).Info("Create for target cluster end")

	omcplog.V(3).Info("Delete for source cluster")
	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.Error("source cluster delete rror : ", sourceErr)
		return sourceErr
	}
	omcplog.V(3).Info("Delete for source cluster end")
	return nil
}

func migserviceNotVolume(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) error {
	omcplog.V(3).Info("=== service migration start")
	omcplog.V(3).Info("add Volume")
	targetResource := &corev1.Service{}
	sourceResource := &corev1.Service{}

	omcplog.V(3).Info("SourceClient init...")
	omcplog.V(3).Info("TargetClient init...")
	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	omcplog.V(3).Info("SourceClient init complete")
	omcplog.V(3).Info("TargetClient init complete")

	omcplog.V(3).Info("Make deployment resource info...")
	nameSpace := migSource.nameSpace

	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.Error("get source cluster error : ", sourceGetErr)
		return sourceGetErr
	} else {
		omcplog.V(3).Info("get source info : " + resource.ResourceName)
	}
	targetResource = sourceResource
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	omcplog.V(3).Info("Make deployment resource info complete")

	omcplog.V(3).Info("Create for target cluster")
	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.Error("target cluster create error : ", targetErr)
		return targetErr
	}
	//	time.Sleep(time.Second * 4)
	omcplog.V(3).Info("Create for target cluster end")

	omcplog.V(3).Info("Delete for source cluster")
	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.V(3).Info("source cluster delete error")
		omcplog.Error(sourceErr)
		return sourceErr
	}
	omcplog.V(3).Info("Delete for source cluster end")
	return nil
}
