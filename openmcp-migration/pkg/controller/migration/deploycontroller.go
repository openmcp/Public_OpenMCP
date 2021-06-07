package migration

import (
	"context"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"
	"time"

	appsv1 "k8s.io/api/apps/v1"
)

func migdeploy(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) {
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
		omcplog.V(3).Info("get source cluster error : ", sourceGetErr)
	}
	omcplog.V(3).Info("Make resource info complete")

	omcplog.V(3).Info("LinkShare init...")
	targetResource = sourceResource.DeepCopy()
	client := cm.Cluster_kubeClients[sourceCluster]
	restconfig := cm.Cluster_configs[sourceCluster]

	addpvcErr, _, newPv, newPvc := CreateLinkShare(sourceClient, sourceResource, volumePath, serviceName, resourceRequire)
	//addpvcErr, _ := CreateLinkShare(sourceClient, sourceResource, volumePath, serviceName)
	if addpvcErr != true {
		omcplog.V(0).Info("add pvc error : ", addpvcErr)
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
		omcplog.V(0).Info("volume make dir error : ", err)
	} else {
		copyErr := LinkShareVolume(client, restconfig, podName, copyCommand, nameSpace)
		if copyErr != nil {
			omcplog.V(0).Info("volume linkshare error : ", copyErr)
		} else {
			omcplog.V(3).Info("volume linkshare complete")
		}
	}
	omcplog.V(3).Info("LinkShare volume copy end")

	omcplog.V(3).Info("Delete for source cluster")
	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.V(3).Info("source cluster delete error: ", sourceErr)
	}
	pvErr := sourceClient.Delete(context.TODO(), newPv, nameSpace, newPv.Name)
	if pvErr != nil {
		omcplog.V(3).Info("source cluster delete error: ", pvErr)
	}
	pvcError := sourceClient.Delete(context.TODO(), newPvc, nameSpace, newPvc.Name)
	if pvcError != nil {
		omcplog.V(3).Info("source cluster delete error : ", pvcError)
	}
	omcplog.V(3).Info("Delete for source cluster end")

	time.Sleep(time.Second * 1)

	omcplog.V(3).Info("Create for target cluster")
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.Spec.Template.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.V(3).Info("target cluster create error : ", targetErr)
	}
	time.Sleep(time.Second * 2)
	omcplog.V(3).Info("Create for target cluster end")

}

func migdeployNotVolume(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) {
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
		omcplog.V(3).Info("get source cluster error : ", sourceGetErr)
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
		omcplog.V(3).Info("target cluster create error : ", targetErr)
	}
	omcplog.V(3).Info("Create for target cluster end")

	omcplog.V(3).Info("Delete for source cluster")
	sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceErr != nil {
		omcplog.V(3).Info("source cluster delete error : ", sourceErr)
	}
	omcplog.V(3).Info("Delete for source cluster end")
}
