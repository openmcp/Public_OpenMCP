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
	"os"
	"regexp"

	// "openmcp/openmcp/migration/pkg/apis"
	nanumv1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/migration/pkg/apis"
	config "openmcp/openmcp/migration/pkg/util"
	"openmcp/openmcp/util/clusterManager"

	restclient "k8s.io/client-go/rest"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/client/generic"
)

var cm *clusterManager.ClusterManager

//pod 이름 찾기
func GetPodName(targetClient generic.Client, dpName string, namespace string) string {
	podInfo := &corev1.Pod{}

	listOption := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{
			"name": dpName,
		}),
	}

	targetClient.List(context.TODO(), podInfo, namespace, listOption)

	podName := podInfo.ObjectMeta.Name

	return podName
}

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	cm = myClusterManager
	////omcplog.Info("[OpenMCP] NewController")
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}
	////omcplog.Info("[OpenMCP] 1111111111111")
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients = append(ghostclients, ghostclient)
	}
	////omcplog.Info("[OpenMCP] 2222222222222")
	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}
	////omcplog.Info("[OpenMCP] 33333333333333333")
	if err := co.WatchResourceReconcileObject(live, &nanumv1alpha1.Migration{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	////omcplog.Info("[OpenMCP] 44444444444444444444")
	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

func GetCopyPodName(clusterClient *kubernetes.Clientset, dpName string, namespace string) string {
START:
	pods, _ := clusterClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	podName := ""
	for _, pod := range pods.Items {
		result, err := regexp.MatchString(dpName, pod.Name)
		if err == nil && result == true && pod.Status.Phase == corev1.PodRunning {
			//omcplog.Info("pod:  ", pod)
			podName = pod.Name
			break
		} else if err != nil {
			//omcplog.Info("pod name error ", err)
			continue
		} else {
			continue
		}
	}
	if podName == "" {
		//omcplog.Info("can not found podName")
		goto START
	} else {
		//omcplog.Info("podName :  ", podName)
		return podName
	}
}
func CopyToNfsCMD(oriVolumePath string, serviceName string) (string, string) {
	mkCMD := config.MKDIR_CMD + " " + config.EXTERNAL_NFS_PATH + "/" + serviceName
	copyCMD := config.COPY_CMD + " " + oriVolumePath + " " + config.EXTERNAL_NFS_PATH + "/" + serviceName
	return mkCMD, copyCMD
}

//pod 내부의 볼륨 폴더를 external nfs로 복사
func LinkShareVolume(client kubernetes.Interface, config *restclient.Config, podName string, command string, namespace string) error {

	cmd := []string{
		"sh",
		"-c",
		command,
	}

	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec")
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
		//omcplog.Info("err----------", err)

	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	//omcplog.Info("123123123-----", command)
	if err != nil {
		return fmt.Errorf("error stream: %v", err)
	} else {
		return err
	}
}

func CreateLinkShare(client generic.Client, sourceResource *appsv1.Deployment, volumePath string, mountPath string, serviceName string) (bool, error) {
	//nfs-pvc append
	resourceInfo := &appsv1.Deployment{}
	resourceInfo = sourceResource
	nfsPvc := corev1.Volume{
		Name: config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
				ReadOnly:  false,
			},
		},
	}
	oriPvc := resourceInfo.Spec.Template.Spec.Volumes
	resourceInfo.Spec.Template.Spec.Volumes = append(oriPvc, nfsPvc)

	nfsMount := corev1.VolumeMount{
		Name:      config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
		ReadOnly:  false,
		MountPath: config.EXTERNAL_NFS_PATH,
	}
	oriVolumeMount := resourceInfo.Spec.Template.Spec.Containers[0].VolumeMounts

	resourceInfo.Spec.Template.Spec.Containers[0].VolumeMounts = append(oriVolumeMount, nfsMount)

	nfsNewPv := &corev1.PersistentVolume{}
	nfsNewPvc := &corev1.PersistentVolumeClaim{}

	nfsNewPv.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PV + "-" + serviceName
	nfsNewPv.Kind = config.PV
	nfsNewPv.Labels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PV + "-" + serviceName,
	}
	nfsNewPv.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
		corev1.ReadWriteMany,
	}
	nfsNewPv.Spec.PersistentVolumeReclaimPolicy = corev1.PersistentVolumeReclaimDelete
	nfsNewPv.Spec.NFS = &corev1.NFSVolumeSource{
		Server:   config.EXTERNAL_NFS,
		Path:     config.EXTERNAL_NFS_PATH,
		ReadOnly: false,
	}
	nfsNewPv.Spec.Capacity = corev1.ResourceList{
		corev1.ResourceStorage: resource.MustParse("10Gi"),
	}
	fmt.Println("링크쉐어 pv생성")

	nfsNewPv.ObjectMeta.ResourceVersion = ""
	nfsNewPv.ResourceVersion = ""
	nfsNewPv.Spec.ClaimRef = nil

	nfsNewPvc.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName
	nfsNewPvc.ObjectMeta.Namespace = config.NameSpace
	nfsNewPvc.Kind = config.PVC
	nfsNewPvc.Labels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
	}
	nfsNewPvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
		corev1.ReadWriteMany,
	}
	nfsNewPvc.Spec.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceStorage: resource.MustParse("10Gi"),
		},
	}
	fmt.Println("링크쉐어 pvc생성")
	nfsNewPvc.ObjectMeta.ResourceVersion = ""
	nfsNewPvc.ResourceVersion = ""

	pvErr := client.Create(context.TODO(), nfsNewPv)
	if pvErr != nil {
		//omcplog.Info(pvErr, " --- nfsPv생성 에러")
	}
	pvcErr := client.Create(context.TODO(), nfsNewPvc)
	if pvcErr != nil {
		//omcplog.Info(pvcErr, " ---- nfsPvc생성 에러")
	}
	//client.Update(context.TODO(), resourceInfo)
	//client.Create(context.TODO(), resourceInfo)
	dpErr := client.Update(context.TODO(), resourceInfo)
	if dpErr != nil {
		//omcplog.Info(dpErr, " ---- dp 생성 에러")
	}
	return true, nil

}
func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	//omcplog.Info("Function Called Reconcile")
	instance := &nanumv1alpha1.Migration{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		//omcplog.Info("get instance error")
	}
	//omcplog.Info(len(instance.Spec.MigrationServiceSources))
	for _, migraionSource := range instance.Spec.MigrationServiceSources {
		//omcplog.Info(migraionSource)
		targetCluster := migraionSource.TargetCluster
		sourceCluster := migraionSource.SourceCluster
		nameSpace := migraionSource.NameSpace
		resourceList := migraionSource.MigrationSources
		volumePath := migraionSource.VolumePath
		mountPath := migraionSource.MountPath
		targetClient := cm.Cluster_genClients[targetCluster]
		sourceClient := cm.Cluster_genClients[sourceCluster]
		serviceName := migraionSource.ServiceName

		for _, resource := range resourceList {
			resourceType := resource.ResourceType
			if resourceType == config.DEPLOY {
				//omcplog.Info("deploy")
				targetResource := &appsv1.Deployment{}
				sourceResource := &appsv1.Deployment{}

				sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				if sourceGetErr != nil {
					//omcplog.V(3).Info("get source cluster")
				}

				client := cm.Cluster_kubeClients[sourceCluster]
				restconfig := cm.Cluster_configs[sourceCluster]

				addpvcErr, _ := CreateLinkShare(sourceClient, sourceResource, volumePath, mountPath, serviceName)
				if addpvcErr != true {
					//omcplog.Info("add pvc error!")
				}
				podName := GetCopyPodName(client, sourceResource.Name, nameSpace)
				mkCommand, copyCommand := CopyToNfsCMD(volumePath, serviceName)
				err := LinkShareVolume(client, restconfig, podName, mkCommand, nameSpace)
				if err != nil {
					//omcplog.Info("volume make dir error")
				} else {
					copyErr := LinkShareVolume(client, restconfig, podName, copyCommand, nameSpace)
					if copyErr != nil {
						//omcplog.Info("volume linkshare error")
					} else {
						//omcplog.Info("volume linkshare complete")
					}
				}

				targetResource = sourceResource

				nfsPvc := []corev1.Volume{
					{
						Name: config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
								ReadOnly:  false,
							},
						},
					},
				}

				targetResource.Spec.Template.Spec.Volumes = nfsPvc

				nfsMount := []corev1.VolumeMount{
					{
						Name:      config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
						ReadOnly:  false,
						MountPath: volumePath,
					},
				}
				// targetResource.Spec.Template.Spec.Containers[0].VolumeMounts = append(targetResource.Spec.Template.Spec.Containers[0].VolumeMounts, nfsMount)
				targetResource.Spec.Template.Spec.Containers[0].VolumeMounts = nfsMount
				targetResource.ObjectMeta.ResourceVersion = ""
				targetResource.Spec.Template.ResourceVersion = ""
				targetResource.ResourceVersion = ""
				targetErr := targetClient.Create(context.TODO(), targetResource)
				if targetErr != nil {
					//omcplog.V(3).Info("target cluster create : " + serviceName)
				}

				// sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				// if sourceErr != nil {
				// 	//omcplog.V(3).Info("source cluster delete : " + resource.ResourceName)
				// }

			} else if resourceType == config.SERVICE {
				//omcplog.Info("service")
				targetResource := &corev1.Service{}
				sourceResource := &corev1.Service{}

				sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				if sourceGetErr != nil {
					//omcplog.V(3).Info("get source cluster")
				}
				targetResource = sourceResource
				targetResource.ObjectMeta.ResourceVersion = ""
				targetResource.ResourceVersion = ""
				targetErr := targetClient.Create(context.TODO(), targetResource)
				if targetErr != nil {
					//omcplog.V(3).Info("target cluster create : " + resource.ResourceName)
				}

				// sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				// if sourceErr != nil {
				// 	//omcplog.V(3).Info("source cluster delete : " + resource.ResourceName)
				// }
			} else if resourceType == config.PV {
				//omcplog.Info("pv")
				targetResource := &corev1.PersistentVolume{}
				sourceResource := &corev1.PersistentVolume{}

				sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				if sourceGetErr != nil {
					//omcplog.V(3).Info("get source cluster")
				}

				targetResource = GetLinkSharePv(sourceResource, volumePath, serviceName)
				targetResource.Spec.Capacity = sourceResource.Spec.Capacity
				targetResource.ObjectMeta.ResourceVersion = ""
				targetResource.ResourceVersion = ""
				targetResource.Spec.ClaimRef = nil

				targetErr := targetClient.Create(context.TODO(), targetResource)
				if targetErr != nil {
					//omcplog.V(3).Info("target cluster create : " + resource.ResourceName)
				}

				// sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				// if sourceErr != nil {
				// 	//omcplog.V(3).Info("source cluster delete : " + resource.ResourceName)
				// }
			} else if resourceType == config.PVC {
				//omcplog.Info("pvc")
				targetResource := &corev1.PersistentVolumeClaim{}
				sourceResource := &corev1.PersistentVolumeClaim{}

				// targetGetErr := targetClient.Get(context.TODO(), targetResource, nameSpace, resource.ResourceName)
				// if targetGetErr != nil {
				// 	//omcplog.V(3).Info("get target cluster")
				// }
				sourceGetErr := sourceClient.Get(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				if sourceGetErr != nil {
					//omcplog.V(3).Info("get source cluster")
				}

				//targetResource = sourceResource
				targetResource = GetLinkSharePvc(sourceResource, volumePath, serviceName)
				targetResource.ObjectMeta.ResourceVersion = ""
				targetResource.ResourceVersion = ""

				targetErr := targetClient.Create(context.TODO(), targetResource)
				if targetErr != nil {
					//omcplog.V(3).Info("target cluster create : " + resource.ResourceName)
				}

				// sourceErr := sourceClient.Delete(context.TODO(), sourceResource, nameSpace, resource.ResourceName)
				// if sourceErr != nil {
				// 	//omcplog.V(3).Info("source cluster delete : " + resource.ResourceName)
				// }
			} else {
				//omcplog.V(3).Info("Resource Type Error!")
				return reconcile.Result{}, fmt.Errorf("Resource Type Error!")
			}
		}
	}
	//omcplog.Info("reconcile end")
	return reconcile.Result{}, nil
}

func GetLinkSharePvc(sourceResource *corev1.PersistentVolumeClaim, volumePath string, serviceName string) *corev1.PersistentVolumeClaim {
	linkSharePvc := &corev1.PersistentVolumeClaim{}
	linkSharePvc = sourceResource
	linkSharePvc.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName
	linkSharePvc.ObjectMeta.Namespace = config.NameSpace
	linkSharePvc.Kind = config.PVC
	linkSharePvc.Spec.VolumeName = config.EXTERNAL_NFS_NAME_PV + "-" + serviceName
	linkSharePvc.Labels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
	}
	linkSharePvc.Spec.Selector.MatchLabels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PV + "-" + serviceName,
	}
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

func GetLinkSharePv(sourceResource *corev1.PersistentVolume, volumePath string, serviceName string) *corev1.PersistentVolume {
	linkSharePv := &corev1.PersistentVolume{}
	linkSharePv.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PV + "-" + serviceName
	linkSharePv.Kind = config.PV
	linkSharePv.Labels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PV + "-" + serviceName,
	}
	linkSharePv.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
		corev1.ReadWriteMany,
	}
	linkSharePv.Spec.PersistentVolumeReclaimPolicy = corev1.PersistentVolumeReclaimDelete
	linkSharePv.Spec.NFS = &corev1.NFSVolumeSource{
		Server:   config.EXTERNAL_NFS,
		Path:     config.EXTERNAL_NFS_PATH + "/" + serviceName,
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
