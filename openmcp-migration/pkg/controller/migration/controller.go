/*
Copyright 2018 The Multicluster-Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package migration

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	// "openmcp/openmcp/migration/pkg/apis"
	"openmcp/openmcp/apis"
	migrationv1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"
	config "openmcp/openmcp/openmcp-migration/pkg/util"
	"openmcp/openmcp/util/clusterManager"

	restclient "k8s.io/client-go/rest"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/kubefed/pkg/controller/util"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/client/generic"
	// "openmcp/openmcp/migration/pkg/controller"
)

var cm *clusterManager.ClusterManager

//pod 이름 찾기
func GetPodName(targetClient generic.Client, dpName string, namespace string) string {
	podInfo := &corev1.Pod{}
	fmt.Println(" pod name 찾기 11111111111111")
	listOption := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{
			"name": dpName,
		}),
	}
	fmt.Println(" pod name 찾기 2222222222222222222")
	targetClient.List(context.TODO(), podInfo, namespace, listOption)
	fmt.Println(" pod name 찾기 3333333333333333")
	podName := podInfo.ObjectMeta.Name
	fmt.Println(" pod name 찾기 44444444444444: " + podName)
	// for _, pod := range podInfo {
	// 	result, err := regexp.MatchString(dpName, pod.Name)
	// 	if err == nil && result == true {
	// 		podName = pod.Name
	// 		break
	// 	} else if err != nil {
	// 		log.Println("pod name error ", err)
	// 		continue
	// 	} else {
	// 		continue
	// 	}
	// }
	return podName
}

//pod 내부의 볼륨 폴더를 external nfs로 복사
// func LinkShareVolume(client kubernetes.Interface, config *restclient.Config, podName string,
// 	command string, namespace string) error {

// 	cmd := []string{
// 		"sh",
// 		"-c",
// 		command,
// 	}
// 	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec")
// 	option := &v1.PodExecOptions{
// 		Command: cmd,
// 		Stdin:   true,
// 		Stdout:  true,
// 		Stderr:  true,
// 		TTY:     false,
// 	}

// 	req.VersionedParams(
// 		option,
// 		scheme.ParameterCodec,
// 	)
// 	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
// 	if err != nil {
// 		return err
// 	}
// 	err = exec.Stream(remotecommand.StreamOptions{
// 		Stdin:  os.Stdin,
// 		Stdout: os.Stdout,
// 		Stderr: os.Stderr,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("error stream: %v", err)

// 	} else {
// 		return err
// 	}
// }
// func CreateLinkShare(deployInfo *appsv1.Deployment, targetClient generic.Client, volumePath string) bool {
// 	nfsPvc := corev1.Volume{
// 		Name: config.EXTERNAL_NFS_NAME_PVC,
// 		VolumeSource: corev1.VolumeSource{
// 			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
// 				ClaimName: config.EXTERNAL_NFS_NAME_PVC,
// 				ReadOnly:  false,
// 			},
// 		},
// 	}
// 	fmt.Println(" 링크쉐어 걸기 11111111111111")
// 	oriPvc := deployInfo.Spec.Template.Spec.Volumes
// 	deployInfo.Spec.Template.Spec.Volumes = append(oriPvc, nfsPvc)
// 	fmt.Println(" 링크쉐어 걸기 222222222222222222")
// 	nfsMount := corev1.VolumeMount{
// 		Name:      config.EXTERNAL_NFS_NAME_PVC ,
// 		ReadOnly:  false,
// 		MountPath: config.EXTERNAL_NFS_PATH ,
// 	}
// 	fmt.Println(" 링크쉐어 걸기 3333333333333333")
// 	deployInfo.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployInfo.Spec.Template.Spec.Containers[0].VolumeMounts, nfsMount)
// 	deleteErr := targetClient.Delete(context.TODO(), deployInfo, deployInfo.Namespace, deployInfo.Name)
// 	if deleteErr != nil {
// 		omcplog.Info("delete deploy false")
// 		return false
// 	}
// 	fmt.Println(" 링크쉐어 걸기 44444444444")
// 	deployInfo.ObjectMeta.ResourceVersion = ""
// 	deployInfo.Spec.Template.ResourceVersion = ""
// 	deployInfo.ResourceVersion = ""
// 	createErr := targetClient.Create(context.TODO(), deployInfo)
// 	fmt.Println(" 링크쉐어 걸기 555555555555555555")
// 	if createErr != nil {
// 		omcplog.Info("linkshare false")
// 		return false
// 	} else {
// 		omcplog.Info("linkshare success")
// 		return true
// 	}

// }
func AddLinkShareVolume(deployInfo *appsv1.Deployment, volumePath string) *appsv1.Deployment {
	nfsPvc := corev1.Volume{
		Name: config.EXTERNAL_NFS_NAME_PVC,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: config.EXTERNAL_NFS_NAME_PVC,
				ReadOnly:  false,
			},
		},
	}
	fmt.Println(" 타겟 클러스터 링크쉐어 111111111111111111111")
	oriPvc := deployInfo.Spec.Template.Spec.Volumes
	deployInfo.Spec.Template.Spec.Volumes = append(oriPvc, nfsPvc)
	fmt.Println(" 타겟 클러스터 링크쉐어 222222222222222222")
	deployInfo.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name = config.EXTERNAL_NFS_NAME_PVC
	deployInfo.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath = config.EXTERNAL_NFS_PATH
	deployInfo.ObjectMeta.ResourceVersion = ""
	deployInfo.Spec.Template.ResourceVersion = ""
	deployInfo.ResourceVersion = ""
	fmt.Println(" 타겟 클러스터 링크쉐어 333333333333333")
	return deployInfo
}

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	cm = myClusterManager
	omcplog.Info("[OpenMCP] NewController")
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}
	omcplog.Info("[OpenMCP] 1111111111111")
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients = append(ghostclients, ghostclient)
	}
	omcplog.Info("[OpenMCP] 2222222222222")
	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}
	omcplog.Info("[OpenMCP] 33333333333333333")
	if err := co.WatchResourceReconcileObject(live, &migrationv1alpha1.Migration{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	omcplog.Info("[OpenMCP] 44444444444444444444")
	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

// func linkShareVolumeCheck(client generic.Client, nameSpace string, volumePath string) (bool, error) {
// 	linkShareVolumePvc := config.EXTERNAL_NFS_NAME_PVC
// 	obj := &corev1.PersistentVolumeClaim{}
// 	err := client.Get(context.TODO(), obj, nameSpace, linkShareVolumePvc)
// 	fmt.Println("링크쉐어 볼륨 체크")
// 	if err == nil {
// 		return true, nil
// 	} else {
// 		client.Create(context.TODO(), GetLinkSharePvc())
// 		client.Create(context.TODO(), GetLinkSharePv(volumePath))
// 		omcplog.Info("create linkshare volume")
// 		return false, nil
// 	}
// }

//해당 deploy에서 생성된 pod명 Get
// func GetPodName(clusterClient *kubernetes.Clientset, dpName string, namespace string) string {
// 	pods, _ := clusterClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
// 	podName := ""
// 	for _, pod := range pods.Items {
// 		result, err := regexp.MatchString(dpName, pod.Name)
// 		if err == nil && result == true {
// 			podName = pod.Name
// 			break
// 		} else if err != nil {
// 			log.Println("pod name error ", err)
// 			continue
// 		} else {
// 			continue
// 		}
// 	}
// 	return podName
// }
func GetCopyPodName(clusterClient *kubernetes.Clientset, dpName string, namespace string) string {
	pods, _ := clusterClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	podName := ""
	for _, pod := range pods.Items {
		result, err := regexp.MatchString(dpName, pod.Name)
		if err == nil && result == true {
			podName = pod.Name
			break
		} else if err != nil {
			log.Println("pod name error ", err)
			continue
		} else {
			continue
		}
	}
	return podName
}
func CopyToNfsCMD(oriVolumePath string) string {
	copyCMD := config.COPY_CMD + " " + oriVolumePath + " " + config.EXTERNAL_NFS_PATH
	return copyCMD
}

//pod 내부의 볼륨 폴더를 external nfs로 복사
func LinkShareVolume(client kubernetes.Interface, config *restclient.Config, podName string, command string, namespace string) error {

	cmd := []string{
		"sh",
		"-c",
		command,
	}
	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec")
	omcplog.Info(podName, "---------", namespace, "------------", command)
	option := &corev1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		return fmt.Errorf("error stream: %v", err)

	} else {
		return err
	}
}

func CreateLinkShare(client generic.Client, sourceResource *appsv1.Deployment, volumePath string) (bool, error) {
	//nfs-pvc append
	resourceInfo := &appsv1.Deployment{}
	resourceInfo = sourceResource
	nfsPvc := corev1.Volume{
		Name: config.EXTERNAL_NFS_NAME_PVC,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: config.EXTERNAL_NFS_NAME_PVC,
				ReadOnly:  false,
			},
		},
	}
	oriPvc := resourceInfo.Spec.Template.Spec.Volumes
	resourceInfo.Spec.Template.Spec.Volumes = append(oriPvc, nfsPvc)

	//resourceInfo.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath = "/nfsdir2/" + lsPvcName
	nfsMount := corev1.VolumeMount{
		Name:      config.EXTERNAL_NFS_NAME_PVC,
		ReadOnly:  false,
		MountPath: config.EXTERNAL_NFS_PATH,
	} //jihoon
	oriVolumeMount := resourceInfo.Spec.Template.Spec.Containers[0].VolumeMounts
	resourceInfo.Spec.Template.Spec.Containers[0].VolumeMounts = append(oriVolumeMount, nfsMount)

	nfsNewPv := &corev1.PersistentVolume{}
	nfsNewPvc := &corev1.PersistentVolumeClaim{}

	nfsNewPv.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PV
	nfsNewPv.Kind = config.PV
	nfsNewPv.Labels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PV,
	}
	nfsNewPv.Spec.PersistentVolumeReclaimPolicy = corev1.PersistentVolumeReclaimDelete
	nfsNewPv.Spec.NFS = &corev1.NFSVolumeSource{
		Server:   config.EXTERNAL_NFS,
		Path:     config.EXTERNAL_NFS_PATH,
		ReadOnly: false,
	}
	fmt.Println("링크쉐어 pv생성")

	nfsNewPv.ObjectMeta.ResourceVersion = ""
	nfsNewPv.ResourceVersion = ""
	nfsNewPv.Spec.ClaimRef = nil

	nfsNewPvc.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PVC
	nfsNewPvc.ObjectMeta.Namespace = config.NameSpace
	nfsNewPvc.Kind = config.PVC
	nfsNewPvc.Labels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PVC,
	}
	fmt.Println("링크쉐어 pvc생성")
	nfsNewPvc.ObjectMeta.ResourceVersion = ""
	nfsNewPvc.ResourceVersion = ""
	client.Delete(context.TODO(), resourceInfo, sourceResource.Namespace, sourceResource.Name)
	// resourceInfo.ObjectMeta.ResourceVersion = ""
	// resourceInfo.Spec.Template.ResourceVersion = ""
	// resourceInfo.ResourceVersion = ""

	client.Create(context.TODO(), nfsNewPv)
	client.Create(context.TODO(), nfsNewPvc)
	//client.Update(context.TODO(), resourceInfo)
	client.Create(context.TODO(), resourceInfo)

	return true, nil

}
func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.Info("Function Called Reconcile")
	instance := &migrationv1alpha1.Migration{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		omcplog.Info("get instance error")
	}
	omcplog.Info(len(instance.Spec.MigrationServiceSources))
	for _, migraionSource := range instance.Spec.MigrationServiceSources {
		omcplog.Info(migraionSource)
		targetCluster := migraionSource.TargetCluster
		sourceCluster := migraionSource.SourceCluster
		nameSpace := migraionSource.NameSpace
		resourceList := migraionSource.MigrationSources
		volumePath := migraionSource.VolumePath
		omcplog.Info("Reconcile 11111111111")
		targetClient := cm.Cluster_genClients[targetCluster]
		sourceClient := cm.Cluster_genClients[sourceCluster]

		for _, resource := range resourceList {
			resourceType := resource.ResourceType
			if resourceType == config.DEPLOY {
				omcplog.Info("deploy")
				targetResource := &appsv1.Deployment{}
				sourceResource := &appsv1.Deployment{}

				sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				if sourceGetErr != nil {
					omcplog.V(3).Info("get source cluster")
				}

				client := cm.Cluster_kubeClients[sourceCluster]
				restconfig := cm.Cluster_configs[sourceCluster]
				podName := GetCopyPodName(client, sourceResource.Name, nameSpace)
				addpvcErr, _ := CreateLinkShare(sourceClient, sourceResource, volumePath)
				if addpvcErr != true {
					omcplog.Info("add pvc error!")
				}

				command := CopyToNfsCMD(volumePath)
				err := LinkShareVolume(client, restconfig, podName, command, nameSpace)
				if err != nil {
					omcplog.Info("volume linkshare error")
				} else {
					omcplog.Info("volume linkshare complete")
				}
				// targetResource = AddLinkShareVolume(sourceResource, volumePath)
				// result := CreateLinkShare(sourceResource, sourceClient)
				// if result == true {
				// 	fmt.Println("링크쉐어 완료")
				// }
				// targetGetErr := targetClient.Get(context.TODO(), targetResource, nameSpace, resource.ResourceName)
				// if targetGetErr != nil {
				// 	omcplog.V(3).Info("get target cluster")
				// }
				targetResource = sourceResource

				nfsPvc := []corev1.Volume{
					{
						Name: config.EXTERNAL_NFS_NAME_PVC,
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: config.EXTERNAL_NFS_NAME_PVC,
								ReadOnly:  false,
							},
						},
					},
				}
				//oriPvc := sourceResource.Spec.Template.Spec.Volumes
				// targetResource.Spec.Template.Spec.Volumes = append(oriPvc, nfsPvc)
				targetResource.Spec.Template.Spec.Volumes = nfsPvc
				omcplog.Info("deploy 1111111111111")
				//resourceInfo.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath = "/nfsdir2/" + lsPvcName
				nfsMount := []corev1.VolumeMount{
					{
						Name:      config.EXTERNAL_NFS_NAME_PVC,
						ReadOnly:  false,
						MountPath: config.EXTERNAL_NFS_PATH,
					},
				}
				omcplog.Info("deploy 22222222222222222")
				// targetResource.Spec.Template.Spec.Containers[0].VolumeMounts = append(targetResource.Spec.Template.Spec.Containers[0].VolumeMounts, nfsMount)
				targetResource.Spec.Template.Spec.Containers[0].VolumeMounts = nfsMount
				omcplog.Info("deploy 33333333333")
				targetResource.ObjectMeta.ResourceVersion = ""
				targetResource.Spec.Template.ResourceVersion = ""
				targetResource.ResourceVersion = ""
				targetErr := targetClient.Create(context.TODO(), targetResource)
				if targetErr != nil {
					omcplog.V(3).Info("target cluster create : " + resource.ResourceName)
				}

				// sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				// if sourceErr != nil {
				// 	omcplog.V(3).Info("source cluster delete : " + resource.ResourceName)
				// }

			} else if resourceType == config.SERVICE {
				omcplog.Info("service")
				targetResource := &corev1.Service{}
				sourceResource := &corev1.Service{}

				sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				if sourceGetErr != nil {
					omcplog.V(3).Info("get source cluster")
				}
				targetResource = sourceResource
				targetResource.ObjectMeta.ResourceVersion = ""
				targetResource.ResourceVersion = ""
				targetErr := targetClient.Create(context.TODO(), targetResource)
				if targetErr != nil {
					omcplog.V(3).Info("target cluster create : " + resource.ResourceName)
				}

				// sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				// if sourceErr != nil {
				// 	omcplog.V(3).Info("source cluster delete : " + resource.ResourceName)
				// }
			} else if resourceType == config.PV {
				omcplog.Info("pv")
				targetResource := &corev1.PersistentVolume{}
				sourceResource := &corev1.PersistentVolume{}

				sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				if sourceGetErr != nil {
					omcplog.V(3).Info("get source cluster")
				}
				targetResource = sourceResource
				targetResource = GetLinkSharePv(sourceResource, volumePath)
				targetResource.ObjectMeta.ResourceVersion = ""
				targetResource.ResourceVersion = ""
				targetResource.Spec.ClaimRef = nil
				// targetGetErr := targetClient.Get(context.TODO(), targetResource, nameSpace, resource.ResourceName)
				// if targetGetErr != nil {
				// 	omcplog.V(3).Info("get target cluster")
				// }

				targetErr := targetClient.Create(context.TODO(), targetResource)
				if targetErr != nil {
					omcplog.V(3).Info("target cluster create : " + resource.ResourceName)
				}

				// sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				// if sourceErr != nil {
				// 	omcplog.V(3).Info("source cluster delete : " + resource.ResourceName)
				// }
			} else if resourceType == config.PVC {
				omcplog.Info("pvc")
				targetResource := &corev1.PersistentVolumeClaim{}
				sourceResource := &corev1.PersistentVolumeClaim{}

				// targetGetErr := targetClient.Get(context.TODO(), targetResource, nameSpace, resource.ResourceName)
				// if targetGetErr != nil {
				// 	omcplog.V(3).Info("get target cluster")
				// }
				sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				if sourceGetErr != nil {
					omcplog.V(3).Info("get source cluster")
				}

				//targetResource = sourceResource
				targetResource = GetLinkSharePvc(sourceResource, volumePath)
				targetResource.ObjectMeta.ResourceVersion = ""
				targetResource.ResourceVersion = ""

				targetErr := targetClient.Create(context.TODO(), targetResource)
				if targetErr != nil {
					omcplog.V(3).Info("target cluster create : " + resource.ResourceName)
				}

				// sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				// if sourceErr != nil {
				// 	omcplog.V(3).Info("source cluster delete : " + resource.ResourceName)
				// }
			} else {
				omcplog.V(3).Info("Resource Type Error!")
				return reconcile.Result{}, fmt.Errorf("Resource Type Error!")
			}
		}
	}
	omcplog.Info("reconcile end")
	return reconcile.Result{}, nil
}

func GetLinkSharePvc(sourceResource *corev1.PersistentVolumeClaim, volumePath string) *corev1.PersistentVolumeClaim {
	linkSharePvc := &corev1.PersistentVolumeClaim{}
	linkSharePvc = sourceResource
	linkSharePvc.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PVC
	linkSharePvc.ObjectMeta.Namespace = config.NameSpace
	linkSharePvc.Kind = config.PVC
	linkSharePvc.Labels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PVC,
	}
	omcplog.Info("111111111 ")
	// linkSharePvc.Spec.Selector = &metav1.LabelSelector{
	// 	MatchLabels: map[string]string{
	// 		"name": config.EXTERNAL_NFS_NAME_PV + sourceResource.ObjectMeta.Labels["name"],
	// 	},
	// }
	fmt.Println("링크쉐어 pvc생성")
	linkSharePvc.ObjectMeta.ResourceVersion = ""
	linkSharePvc.ResourceVersion = ""
	return linkSharePvc
}

func GetLinkSharePv(sourceResource *corev1.PersistentVolume, volumePath string) *corev1.PersistentVolume {
	linkSharePv := &corev1.PersistentVolume{}
	linkSharePv.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PV
	linkSharePv.Kind = config.PV
	linkSharePv.Labels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PV,
	}
	linkSharePv.Spec.PersistentVolumeReclaimPolicy = corev1.PersistentVolumeReclaimDelete
	linkSharePv.Spec.NFS = &corev1.NFSVolumeSource{
		Server:   config.EXTERNAL_NFS,
		Path:     config.EXTERNAL_NFS_PATH,
		ReadOnly: false,
	}
	fmt.Println("링크쉐어 pv생성")

	linkSharePv.ObjectMeta.ResourceVersion = ""
	linkSharePv.ResourceVersion = ""
	linkSharePv.Spec.ClaimRef = nil
	// linkSharePv.Spec = sourceResource.Spec
	// 현재 nfs 진행 (해당 클러스터 pv정보 필요)
	// 	var volumeInfo *corev1.NFSVolumeSource
	// if sourceResource.Spec.NFS != nil {
	// 	json.Unmarshal([]byte(volumeString), volumeInfo)
	// 	linkSharePv.Spec.NFS = volumeInfo
	// 	linkSharePv.Spec.NFS.Path = VolumePath
	// } else if sourceResource.Spec.HostPath != nil {
	// 	var volumeInfo *corev1.HostPathVolumeSource
	// 	json.Unmarshal([]byte(volumeString), &volumeInfo)
	// 	resourceInfo.Spec.HostPath = volumeInfo
	// 	resourceInfo.Spec.HostPath.Path = VolumePath
	// } else if sourceResource.Spec.ISCSI != nil {
	// 	var volumeInfo *corev1.ISCSIPersistentVolumeSource
	// 	json.Unmarshal([]byte(volumeString), &volumeInfo)
	// 	resourceInfo.Spec.ISCSI = volumeInfo
	// 	//resourceInfo.Spec.ISCSI. = VolumePath
	// } else if sourceResource.Spec.Glusterfs != nil {
	// 	var volumeInfo *corev1.GlusterfsPersistentVolumeSource
	// 	json.Unmarshal([]byte(volumeString), &volumeInfo)
	// 	resourceInfo.Spec.Glusterfs = volumeInfo
	// 	resourceInfo.Spec.Glusterfs.Path = VolumePath
	// }

	return linkSharePv
}
