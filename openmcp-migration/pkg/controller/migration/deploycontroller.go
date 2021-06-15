package migration

import (
	"context"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"
	"time"

	appsv1 "k8s.io/api/apps/v1"
)

func migdeploy(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) error {
	omcplog.V(3).Info("=== deploy migration start")
	omcplog.V(3).Info("add Volume")

	targetResource := &appsv1.Deployment{}
	sourceResource := &appsv1.Deployment{}

	omcplog.V(3).Info("SourceClient init...")
	omcplog.V(3).Info("TargetClient init...")
	time.Sleep(time.Second * 2)
	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	omcplog.V(3).Info("SourceClient init complete")
	omcplog.V(3).Info("TargetClient init complete")

	omcplog.V(3).Info("Make resource info...")
	time.Sleep(time.Second * 2)
	nameSpace := migSource.nameSpace
	sourceCluster := migSource.sourceCluster
	volumePath := migSource.volumePath
	serviceName := migSource.serviceName
	resourceRequire := migSource.resourceRequire

	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.Error("get source cluster error : ", sourceGetErr)
		return sourceGetErr
	}
	omcplog.V(3).Info("Make resource info complete")

	omcplog.V(3).Info("LinkShare init...")
	targetResource = sourceResource.DeepCopy()
	client := cm.Cluster_kubeClients[sourceCluster]
	restconfig := cm.Cluster_configs[sourceCluster]

	_, addpvcErr, newPv, newPvc := CreateLinkShare(sourceClient, sourceResource, volumePath, serviceName, resourceRequire)
	//addpvcErr, _ := CreateLinkShare(sourceClient, sourceResource, volumePath, serviceName)
	if addpvcErr != nil {
		omcplog.Error("add pvc error : ", addpvcErr)
		return addpvcErr
	}
	omcplog.V(3).Info(" volumePath : " + volumePath)
	omcplog.V(3).Info(" serviceName : " + serviceName)
	omcplog.V(3).Info("LinkShare init end")

	// time.Sleep(time.Second * 3)

	omcplog.V(3).Info("LinkShare volume copy")
	podName := GetCopyPodName(client, sourceResource.Name, nameSpace)
	mkCommand, copyCommand := CopyToNfsCMD(volumePath, serviceName)

	err := LinkShareVolume(client, restconfig, podName, mkCommand, nameSpace)
	if err != nil {
		omcplog.Error("volume make dir error : ", err)
		return err
	} else {
		copyErr := LinkShareVolume(client, restconfig, podName, copyCommand, nameSpace)
		if copyErr != nil {
			omcplog.Error("volume linkshare error : ", copyErr)
			return copyErr
		} else {
			omcplog.V(3).Info("volume linkshare complete")
		}
	}
	omcplog.V(3).Info("LinkShare volume copy end")

	omcplog.V(3).Info("Delete for source cluster")
	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.Error("source cluster deploy delete error : ", sourceErr)
		return sourceErr
	}
	pvcError := sourceClient.Delete(context.TODO(), newPvc, nameSpace, newPvc.Name)
	if pvcError != nil {
		omcplog.Error("source cluster pvc delete error : ", pvcError)
		return pvcError
	}
	pvErr := sourceClient.Delete(context.TODO(), newPv, nameSpace, newPv.Name)
	if pvErr != nil {
		omcplog.Error("source cluster pv delete error: ", pvErr)
		return pvErr
	}
	omcplog.V(3).Info("Delete for source cluster end")

	time.Sleep(time.Second * 1)

	omcplog.V(3).Info("Create for target cluster")
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.Spec.Template.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.Error("target cluster deploy create error : ", targetErr)
		return targetErr
	}
	time.Sleep(time.Second * 2)
	omcplog.V(3).Info("Create for target cluster end")
	return nil
}

func migdeployNotVolume(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) error {
	omcplog.V(3).Info("=== deploy migration start : " + resource.ResourceName)
	omcplog.V(3).Info("not Volume")

	targetResource := &appsv1.Deployment{}
	sourceResource := &appsv1.Deployment{}

	omcplog.V(3).Info("SourceClient init...")
	omcplog.V(3).Info("TargetClient init...")
	//	time.Sleep(time.Second * 3)
	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	omcplog.V(3).Info("SourceClient init complete")
	omcplog.V(3).Info("TargetClient init complete")

	omcplog.V(3).Info("Make resource info...")
	//	time.Sleep(time.Second * 3)
	nameSpace := migSource.nameSpace

	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.Error("get source cluster error : ", sourceGetErr)
		return sourceGetErr
	} else {
		omcplog.V(3).Info("get source info : " + resource.ResourceName)
	}
	omcplog.V(3).Info("Make deployment resource info complete")

	omcplog.V(3).Info("Create for target cluster")
	targetResource = sourceResource.DeepCopy()
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.Spec.Template.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.Error("target cluster create error : ", targetErr)
		return targetErr
	}
	omcplog.V(3).Info("Create for target cluster end")

	omcplog.V(3).Info("Delete for source cluster")
	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.Error("source cluster delete error : ", sourceErr)
		return sourceErr
	}
	omcplog.V(3).Info("Delete for source cluster end")
	return nil
}
