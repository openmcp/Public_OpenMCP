package migration

import (
	"context"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"

	appsv1 "k8s.io/api/apps/v1"
)

func migdeploy(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) {
	omcplog.V(3).Info("deploy migration start")
	targetResource := &appsv1.Deployment{}
	sourceResource := &appsv1.Deployment{}

	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	nameSpace := migSource.nameSpace
	sourceCluster := migSource.sourceCluster
	volumePath := migSource.volumePath
	serviceName := migSource.serviceName
	resourceRequire := migSource.resourceRequire

	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.V(3).Info("get source cluster error : ", sourceGetErr)
	}
	targetResource = sourceResource.DeepCopy()
	client := cm.Cluster_kubeClients[sourceCluster]
	restconfig := cm.Cluster_configs[sourceCluster]

	addpvcErr, _, newPv, newPvc := CreateLinkShare(sourceClient, sourceResource, volumePath, serviceName, resourceRequire)
	//addpvcErr, _ := CreateLinkShare(sourceClient, sourceResource, volumePath, serviceName)
	if addpvcErr != true {
		omcplog.V(0).Info("add pvc error : ", addpvcErr)
	}
	// time.Sleep(time.Second * 3)
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
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.Spec.Template.ResourceVersion = ""
	targetResource.ResourceVersion = ""
	targetErr := targetClient.Create(context.TODO(), targetResource)
	if targetErr != nil {
		omcplog.V(3).Info("target cluster create error : ", targetErr)
	}
	// nfsPvc := []corev1.Volume{
	// 	{
	// 		Name: config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
	// 		VolumeSource: corev1.VolumeSource{
	// 			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
	// 				ClaimName: config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
	// 				ReadOnly:  false,
	// 			},
	// 		},
	// 	},
	// }

	// targetResource.Spec.Template.Spec.Volumes = nfsPvc

	// nfsMount := []corev1.VolumeMount{
	// 	{
	// 		Name:      config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
	// 		ReadOnly:  false,
	// 		MountPath: volumePath,
	// 	},
	// }
	// // targetResource.Spec.Template.Spec.Containers[0].VolumeMounts = append(targetResource.Spec.Template.Spec.Containers[0].VolumeMounts, nfsMount)
	// targetResource.Spec.Template.Spec.Containers[0].VolumeMounts = nfsMount

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
}

func migdeployNotVolume(migSource MigrationControllerResource, resource v1alpha1.MigrationSource) {
	omcplog.V(3).Info("deploy migration start : " + resource.ResourceName)
	targetResource := &appsv1.Deployment{}
	sourceResource := &appsv1.Deployment{}

	sourceClient := migSource.sourceClient
	targetClient := migSource.targetClient
	nameSpace := migSource.nameSpace

	sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
	if sourceGetErr != nil {
		omcplog.V(3).Info("get source cluster error : ", sourceGetErr)
	} else {
		omcplog.V(3).Info("get source info : " + resource.ResourceName)
	}

	targetResource = sourceResource.DeepCopy()
	targetResource.ObjectMeta.ResourceVersion = ""
	targetResource.Spec.Template.ResourceVersion = ""
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
