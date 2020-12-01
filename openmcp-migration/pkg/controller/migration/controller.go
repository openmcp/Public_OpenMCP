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
	"time"

	// "openmcp/openmcp/migration/pkg/apis"

	"openmcp/openmcp/apis"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	"openmcp/openmcp/omcplog"
	config "openmcp/openmcp/openmcp-migration/pkg/util"
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

	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/kubefed/pkg/controller/util"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/client/generic"
	// "openmcp/openmcp/migration/pkg/controller"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	cm = myClusterManager
	omcplog.V(4).Info("NewController start")
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		omcplog.V(0).Info("getting delegating client for live cluster: ", err)
		return nil, err
	}
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			omcplog.V(0).Info("getting delegating client for ghost cluster: ", err)
			return nil, err
		}
		ghostclients = append(ghostclients, ghostclient)
	}
	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		omcplog.V(0).Info("adding APIs to live cluster's scheme: ", err)
		return nil, err
	}
	if err := co.WatchResourceReconcileObject(live, &v1alpha1.Migration{}, controller.WatchOptions{}); err != nil {
		omcplog.V(0).Info("setting up Pod watch in live cluster: ", err)
		return nil, err
	}
	omcplog.V(4).Info("NewController end")
	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

func sortResource(migraionSource v1alpha1.MigrationServiceSource) ([]v1alpha1.MigrationSource, bool) {
	resourceList := migraionSource.MigrationSources
	pvCheck := false
	for i, j := range resourceList {
		if j.ResourceType == config.DEPLOY {
			tmp := resourceList[i]
			resourceList[i] = resourceList[0]
			resourceList[0] = tmp
		} else if j.ResourceType == config.PV {
			pvCheck = true
		}
	}
	return resourceList, pvCheck
}

type MigrationControllerResource struct {
	resourceList    []v1alpha1.MigrationSource
	targetCluster   string
	sourceCluster   string
	nameSpace       string
	volumePath      string
	serviceName     string
	pvCheck         bool
	sourceClient    generic.Client
	targetClient    generic.Client
	resourceRequire corev1.ResourceList
}

func MigrationControllerResourceInit(migraionSource v1alpha1.MigrationServiceSource) MigrationControllerResource {
	migSource := MigrationControllerResource{}

	migSource.resourceList, migSource.pvCheck = sortResource(migraionSource)
	migSource.targetCluster = migraionSource.TargetCluster
	migSource.sourceCluster = migraionSource.SourceCluster
	migSource.nameSpace = migraionSource.NameSpace
	migSource.volumePath, migSource.resourceRequire = getVolumePath(migraionSource)
	migSource.targetClient = cm.Cluster_genClients[migraionSource.TargetCluster]
	migSource.sourceClient = cm.Cluster_genClients[migraionSource.SourceCluster]
	migSource.serviceName = migraionSource.ServiceName

	return migSource
}
func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(3).Info("Function Called Reconcile")
	instance := &v1alpha1.Migration{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		omcplog.V(0).Info("get instance error : ", err)
	}
	for _, migraionSource := range instance.Spec.MigrationServiceSources {
		omcplog.V(4).Info(migraionSource)
		migSource := MigrationControllerResourceInit(migraionSource)
		checkNameSpace(migSource.targetClient, migSource.nameSpace)
		if migSource.pvCheck == true {
			for _, resource := range migSource.resourceList {
				resourceType := resource.ResourceType
				if resourceType == config.DEPLOY {
					migdeploy(migSource, resource)
				} else if resourceType == config.SERVICE {
					migservice(migSource, resource)
				} else if resourceType == config.PV {
					migpv(migSource, resource)
				} else if resourceType == config.PVC {
					migpvc(migSource, resource)
				} else {
					omcplog.V(0).Info("Resource Type Error!")
					return reconcile.Result{}, fmt.Errorf("Resource Type Error!")
				}
			}
		} else {
			for _, resource := range migSource.resourceList {
				resourceType := resource.ResourceType
				if resourceType == config.DEPLOY {
					migdeployNotVolume(migSource, resource)
				} else if resourceType == config.SERVICE {
					migserviceNotVolume(migSource, resource)
				}
			}
		}
	}
	return reconcile.Result{}, nil
}

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

func GetCopyPodName(clusterClient *kubernetes.Clientset, dpName string, namespace string) string {
START:
	pods, _ := clusterClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	podName := ""
	for _, pod := range pods.Items {
		result, err := regexp.MatchString(dpName, pod.Name)
		if err == nil && result == true && pod.Labels["updateCheck"] == "true" {
			//omcplog.V(4).Info("pod STATUS:  ", pod.Status.ContainerStatuses)
			podName = pod.Name
			break
		} else if err != nil {
			//omcplog.V(0).Info("pod name error", err)
			continue
		} else {
			continue
		}
	}
	if podName == "" {
		omcplog.V(4).Info("can not found podName")
		time.Sleep(time.Second * 1)
		goto START
	} else {
		omcplog.V(4).Info("podName :  ", podName)
		return podName
	}
}
func CopyToNfsCMD(oriVolumePath string, serviceName string) (string, string) {
	mkCMD := config.MKDIR_CMD + " " + config.EXTERNAL_NFS_PATH + "/" + serviceName
	copyCMD := config.COPY_CMD + " " + oriVolumePath + " " + config.EXTERNAL_NFS_PATH + "/" + serviceName
	omcplog.V(4).Info("mkdir cmd  :  ", mkCMD)
	omcplog.V(4).Info("copy cmd  :  ", copyCMD)
	return mkCMD, copyCMD
}

//pod 내부의 볼륨 폴더를 external nfs로 복사
func LinkShareVolume(client kubernetes.Interface, config *restclient.Config, podName string, command string, namespace string) error {
	cmd := []string{
		"sh",
		"-c",
		command,
	}
	timeoutcheck := 0
START:
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
		omcplog.V(0).Info("NewSPDYExecutor err: ", err)

	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		timeoutcheck = timeoutcheck + 1
		if timeoutcheck == 10 {
			omcplog.V(0).Info("error stream: ", err)
			return err
		}
		omcplog.V(4).Info("connecting: ", podName)
		time.Sleep(time.Second * 1)
		goto START
	} else {
		return err
	}
}

func CreateLinkShare(client generic.Client, sourceResource *appsv1.Deployment, volumePath string, serviceName string, resourceRequire corev1.ResourceList) (bool, error, *corev1.PersistentVolume, *corev1.PersistentVolumeClaim) {
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
	resourceInfo.Spec.Template.Labels["updateCheck"] = "true"
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
	nfsNewPv.Spec.Capacity = resourceRequire

	nfsNewPv.ObjectMeta.ResourceVersion = ""
	nfsNewPv.ResourceVersion = ""
	nfsNewPv.Spec.ClaimRef = nil

	nfsNewPvc.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName
	nfsNewPvc.ObjectMeta.Namespace = sourceResource.Namespace
	nfsNewPvc.Kind = config.PVC
	nfsNewPvc.Labels = map[string]string{
		"name": config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
	}
	nfsNewPvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
		corev1.ReadWriteMany,
	}
	nfsNewPvc.Spec.Resources = corev1.ResourceRequirements{
		Requests: resourceRequire,
	}

	nfsNewPvc.ObjectMeta.ResourceVersion = ""
	nfsNewPvc.ResourceVersion = ""

	pvErr := client.Create(context.TODO(), nfsNewPv)
	if pvErr != nil {
		omcplog.V(0).Info("nfsPv 생성 에러: ", pvErr)
	}
	omcplog.V(3).Info("nfsPv 생성 완료")
	pvcErr := client.Create(context.TODO(), nfsNewPvc)
	if pvcErr != nil {
		omcplog.V(0).Info("nfsPvc 생성 에러: ", pvcErr)
	}
	omcplog.V(3).Info("nfsPvc 생성 완료")
	//client.Update(context.TODO(), resourceInfo)
	//client.Create(context.TODO(), resourceInfo)
	dpErr := client.Update(context.TODO(), resourceInfo)
	if dpErr != nil {
		omcplog.V(0).Info("dp 생성 에러: ", dpErr)
	}
	omcplog.V(3).Info("dp 생성 완료")
	return true, nil, nfsNewPv, nfsNewPvc
	// return true, nil
}
func getVolumePath(migraionSource v1alpha1.MigrationServiceSource) (string, corev1.ResourceList) {
	pvcName := ""
	dpResource := &appsv1.Deployment{}
	pvcResource := &corev1.PersistentVolumeClaim{}
	sourceCluster := migraionSource.SourceCluster
	sourceClient := cm.Cluster_genClients[sourceCluster]
	volumeName := ""
	mountPath := ""
	for _, resource := range migraionSource.MigrationSources {
		if resource.ResourceType == config.PVC {
			pvcName = resource.ResourceName
			sourceClient.Get(context.TODO(), pvcResource, migraionSource.NameSpace, resource.ResourceName)
		} else if resource.ResourceType == config.DEPLOY {
			sourceClient.Get(context.TODO(), dpResource, migraionSource.NameSpace, resource.ResourceName)
		} else {
			continue
		}
	}
	if pvcName == "" {
		omcplog.V(0).Info("error pvc name nil")
	}

	volumeList := dpResource.Spec.Template.Spec.Volumes
	for _, volume := range volumeList {
		if volume.PersistentVolumeClaim.ClaimName == pvcName {
			volumeName = volume.Name
		}
	}
	if volumeName == "" {
		omcplog.V(0).Info("error volume name nil")
	}

	containerList := dpResource.Spec.Template.Spec.Containers
	for _, container := range containerList {
		for _, volumeMount := range container.VolumeMounts {
			if volumeMount.Name == volumeName {
				mountPath = volumeMount.MountPath
			}
		}
	}
	if mountPath == "" {
		omcplog.V(0).Info("error mount path nil")
	}
	storageSize := corev1.ResourceList{}
	if pvcResource.Spec.Resources.Requests.Storage() != nil {
		storageSize = pvcResource.Spec.Resources.Requests
	} else {
		storageSize = corev1.ResourceList{
			corev1.ResourceStorage: resource.MustParse(config.DEFAULT_VOLUME_SIZE),
		}
	}
	return mountPath, storageSize
}

func checkNameSpace(client generic.Client, namespace string) bool {
	ns := &corev1.NamespaceList{}
	client.List(context.TODO(), ns, "")
	for _, nss := range ns.Items {
		if nss.Name == namespace {
			omcplog.V(4).Info("Already exists namespace: ", namespace)
			return true
		}
	}

	nameSpaceObj := &corev1.Namespace{}
	nameSpaceObj.Name = namespace
	client.Create(context.TODO(), nameSpaceObj)
	omcplog.V(4).Info("Create namespace: ", namespace)
	return true
}
func GetLinkSharePvc(sourceResource *corev1.PersistentVolumeClaim, volumePath string, serviceName string) *corev1.PersistentVolumeClaim {
	linkSharePvc := &corev1.PersistentVolumeClaim{}
	linkSharePvc = sourceResource.DeepCopy()
	// linkSharePvc.ObjectMeta.Name = config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName
	// linkSharePvc.ObjectMeta.Namespace = sourceResource.NameSpace
	// linkSharePvc.Kind = config.PVC
	// linkSharePvc.Spec.VolumeName = config.EXTERNAL_NFS_NAME_PV + "-" + serviceName
	// linkSharePvc.Labels = map[string]string{
	// 	"name": config.EXTERNAL_NFS_NAME_PVC + "-" + serviceName,
	// }
	// linkSharePvc.Spec.Selector.MatchLabels = map[string]string{
	// 	"name": config.EXTERNAL_NFS_NAME_PV + "-" + serviceName,
	// }
	// linkSharePvc.Spec.Selector = &metav1.LabelSelector{
	// 	MatchLabels: map[string]string{
	// 		"name": config.EXTERNAL_NFS_NAME_PV + sourceResource.ObjectMeta.Labels["name"],
	// 	},
	// }
	linkSharePvc.ObjectMeta.ResourceVersion = ""
	linkSharePvc.ResourceVersion = ""
	return linkSharePvc
}

func GetLinkSharePv(sourceResource *corev1.PersistentVolume, volumePath string, serviceName string) *corev1.PersistentVolume {
	linkSharePv := &corev1.PersistentVolume{}
	linkSharePv.ObjectMeta.Name = sourceResource.ObjectMeta.Name
	linkSharePv.Kind = config.PV
	// linkSharePv.Labels = map[string]string{
	// 	"name": config.EXTERNAL_NFS_NAME_PV + "-" + serviceName,
	// }
	linkSharePv.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{
		corev1.ReadWriteMany,
	}
	linkSharePv.Spec.PersistentVolumeReclaimPolicy = corev1.PersistentVolumeReclaimDelete
	linkSharePv.Spec.NFS = &corev1.NFSVolumeSource{
		Server:   config.EXTERNAL_NFS,
		Path:     config.EXTERNAL_NFS_PATH + "/" + serviceName,
		ReadOnly: false,
	}

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
