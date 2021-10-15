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
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	// "openmcp/openmcp/migration/pkg/apis"

	"openmcp/openmcp/apis"
	v1alpha1 "openmcp/openmcp/apis/migration/v1alpha1"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
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
	"k8s.io/apimachinery/pkg/types"
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
	if err := co.WatchResourceReconcileObject(context.TODO(), live, &v1alpha1.Migration{}, controller.WatchOptions{}); err != nil {
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
	moveCount      int32
}

// sourceResource : Deploy인 경우 0번에, PV인 경우 pvCheck 를 true로 리턴시키는 함수.
func sortResource(migraionSource v1alpha1.MigrationServiceSource) ([]v1alpha1.MigrationSource, bool) {
	resourceList := migraionSource.MigrationSources
	pvCheck := false
	for i, resource := range resourceList {
		if resource.ResourceType == config.DEPLOY {
			tmp := resourceList[i]
			resourceList[i] = resourceList[0]
			resourceList[0] = tmp
		} else if resource.ResourceType == config.PV {
			pvCheck = true
		}
	}
	return resourceList, pvCheck
}

// setResource : Deploy
func (r *reconciler) setResourceForOnlyDeploy(migraionSource v1alpha1.MigrationServiceSource, migSource *MigrationControllerResource) error {

	sourceClient := migSource.sourceClient
	pvcNames := []string{}
	pvNames := []string{}

	resourceList := migraionSource.MigrationSources
	pvCheck := false
	namespace := migraionSource.NameSpace
	for i, resource := range resourceList {
		resourceName := resource.ResourceName
		if resource.ResourceType == config.DEPLOY {

			// 0. OpenMCPDeployment 의 기능을 잠시 정지시킨다.
			err := r.setBeforeOpenmcpDeploymnet(migraionSource, i)
			if err != nil {
				omcplog.Error("setBeforeOpenmcpDeploymnet error : ", err)
				return err
			}

			// 1. Deployment 데이터를 가져온다.
			omcplog.V(0).Info("- 1. getDeployment -")
			omcplog.V(0).Info(resource.ResourceName)

			sourceDeploy := &appsv1.Deployment{}
			_ = sourceClient.Get(context.TODO(), sourceDeploy, namespace, resourceName)
			omcplog.V(0).Info(sourceDeploy)
			volumeInfo := sourceDeploy.Spec.Template.Spec.DeepCopy().Volumes
			omcplog.V(0).Info(volumeInfo)
			for i, volume := range volumeInfo {
				// 2. PVC 정보가 있는지 체크하고 있으면 기입.
				omcplog.V(0).Info("- 2. check deployment -" + string(rune(i)))
				omcplog.V(0).Info(volume)
				pvcInfo := volume.PersistentVolumeClaim
				if pvcInfo == nil {
					continue
				} else {
					omcplog.V(0).Info("--- pvc bingo ---")
					pvcNames = append(pvcNames, pvcInfo.ClaimName)

					// 3. PVC 정보를 토대로 PV 라벨 정보 추출 후 이름 추출.
					sourcePVC := &corev1.PersistentVolumeClaim{}
					_ = sourceClient.Get(context.TODO(), sourcePVC, namespace, pvcInfo.ClaimName)
					pvc_matchLabel := sourcePVC.Spec.Selector.DeepCopy().MatchLabels
					omcplog.V(0).Info("- 3. check pvc -")
					omcplog.V(0).Info(pvc_matchLabel)
					sourcePVList := &corev1.PersistentVolumeList{}
					_ = sourceClient.List(context.TODO(), sourcePVList, namespace, &client.ListOptions{
						LabelSelector: labels.SelectorFromSet(labels.Set(pvc_matchLabel)),
					})
					for _, pv := range sourcePVList.Items {
						omcplog.V(0).Info("- 4. check pv -")
						omcplog.V(0).Info(pv)
						pvNames = append(pvNames, pv.Name)
					}
				}
			}

			// 5. 저장된 pvName, pvcName 을 토대로 리스트 재작성.
			for _, pvName := range pvNames {
				tmp := v1alpha1.MigrationSource{}
				tmp.ResourceName = pvName
				tmp.ResourceType = config.PV
				resourceList = append(resourceList, tmp)
				pvCheck = true
			}
			for _, pvcName := range pvcNames {
				tmp := v1alpha1.MigrationSource{}
				tmp.ResourceName = pvcName
				tmp.ResourceType = config.PVC
				resourceList = append(resourceList, tmp)
				pvCheck = true
			}
		}
	}
	omcplog.V(0).Info("- resourceList fix -")
	omcplog.V(0).Info(resourceList)

	// 6. 기존에 하던 데이터 보정 실행 (동일 객체 제거, Deployment를 제일 앞으로)
	fixedResourceList := []v1alpha1.MigrationSource{}
	//동일 객체 제거
	for _, resource := range resourceList {
		isConflict := false
		for _, fixedResource := range fixedResourceList {
			if fixedResource.ResourceName == resource.ResourceName && fixedResource.ResourceType == resource.ResourceType {
				// 동일한 경우 리스트에 추가하지 않는다.
				omcplog.V(0).Info(" -- conflict resource ")
				omcplog.V(0).Info(resource)
				isConflict = true
			}
		}
		if !isConflict {
			fixedResourceList = append(fixedResourceList, resource)
		}
	}
	omcplog.V(0).Info(fixedResourceList)
	//Deploy를 가장 앞으로.
	for i, fixedResource := range fixedResourceList {
		if fixedResource.ResourceType == config.DEPLOY {
			tmp := fixedResourceList[i]
			fixedResourceList[i] = fixedResourceList[0]
			fixedResourceList[0] = tmp
		}
	}

	omcplog.V(0).Info("->")
	omcplog.V(0).Info(fixedResourceList)
	omcplog.V(0).Info("--------------------")

	migSource.resourceList = fixedResourceList
	migSource.pvCheck = pvCheck
	return nil
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

func ResourceInit(migraionSource v1alpha1.MigrationServiceSource) (MigrationControllerResource, error) {
	migSource := MigrationControllerResource{}

	migSource.targetCluster = migraionSource.TargetCluster
	migSource.sourceCluster = migraionSource.SourceCluster
	migSource.nameSpace = migraionSource.NameSpace
	migSource.targetClient = cm.Cluster_genClients[migraionSource.TargetCluster]
	migSource.sourceClient = cm.Cluster_genClients[migraionSource.SourceCluster]
	migSource.serviceName = migraionSource.ServiceName
	migSource.resourceList, migSource.pvCheck = sortResource(migraionSource)

	var err error
	err = getVolumePath(&migSource)
	if err != nil {
		omcplog.Error("getVolumePath error : ", err)
		return migSource, err
	}

	return migSource, nil
}

// ResourceInitForOnlyDeployment
func (r *reconciler) ResourceInitForOnlyDeployment(migraionSource v1alpha1.MigrationServiceSource) (MigrationControllerResource, error) {
	migSource := MigrationControllerResource{}
	migSource.targetCluster = migraionSource.TargetCluster
	migSource.sourceCluster = migraionSource.SourceCluster
	migSource.nameSpace = migraionSource.NameSpace
	migSource.targetClient = cm.Cluster_genClients[migraionSource.TargetCluster]
	migSource.sourceClient = cm.Cluster_genClients[migraionSource.SourceCluster]
	migSource.serviceName = migraionSource.ServiceName

	var err error
	//migSource.resourceList, migSource.pvCheck = setResourceForOnlyDeploy(migraionSource)
	err = r.setResourceForOnlyDeploy(migraionSource, &migSource)
	if err != nil {
		omcplog.Error("setResourceForOnlyDeploy error : ", err)
		return migSource, err
	}
	//err, migSource.volumePath, migSource.resourceRequire = getVolumePath(migraionSource)
	err = getVolumePath(&migSource)
	if err != nil {
		omcplog.Error("getVolumePath error : ", err)
		return migSource, err
	}

	return migSource, nil
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(3).Info("Migration Start : Reconcile")
	startDate := time.Now()
	instance := &v1alpha1.Migration{}
	checkResourceName := ""
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		omcplog.Error("get instance error : ", err)
		r.MakeStatus(instance, false, "", err)
		return reconcile.Result{Requeue: false}, nil
	}

	if instance.Status.MigrationStatus == true {
		// 이미 성공한 케이스는 로직을 안탄다.
		omcplog.V(4).Info(instance.Name + " already succeed")
		return reconcile.Result{Requeue: false}, nil
	}
	if instance.Status.MigrationStatus == false && instance.Status.Reason != "" {
		// 이미 실패한 케이스는 로직을 다시 안탄다.
		omcplog.V(4).Info(instance.Name + " already failed")
		return reconcile.Result{Requeue: false}, nil
	}

	for _, migraionSource := range instance.Spec.MigrationServiceSources {
		omcplog.V(4).Info(migraionSource)
		//migSource, initErr := ResourceInit(migraionSource)
		migSource, initErr := r.ResourceInitForOnlyDeployment(migraionSource)
		if initErr != nil {
			omcplog.Error("MigrationControllerResource Init error : ", initErr)
			r.MakeStatusWithMigSource(instance, false, migraionSource, initErr)
			omcplog.V(3).Info("Migration Failed")
			return reconcile.Result{Requeue: false}, nil
		}
		checkNameSpace(migSource.targetClient, migSource.nameSpace)
		var migErr error
		if migSource.pvCheck == true {
			for _, resource := range migSource.resourceList {
				resourceType := resource.ResourceType
				if resourceType == config.DEPLOY {
					migErr = migdeploy(migSource, resource)
					checkResourceName = resource.ResourceName
				} else if resourceType == config.SERVICE {
					migErr = migservice(migSource, resource)
				} else if resourceType == config.PV {
					migErr = migpv(migSource, resource)
				} else if resourceType == config.PVC {
					migErr = migpvc(migSource, resource)
				} else {
					omcplog.Error(fmt.Errorf("Resource Type Error!"))
					r.MakeStatusWithResource(instance, false, migraionSource, resource, fmt.Errorf("Resource Type Error!"))
					omcplog.Error("Migration Failed")
					return reconcile.Result{Requeue: false}, nil
				}
				if migErr != nil {
					omcplog.Error("MigrationControllerResource migration error : ", migErr)
					r.MakeStatusWithResource(instance, false, migraionSource, resource, migErr)
					omcplog.Error("Migration Failed")
					return reconcile.Result{Requeue: false}, nil
				}
			}
		} else {
			for i, resource := range migSource.resourceList {
				resourceType := resource.ResourceType
				if resourceType == config.DEPLOY {
					migErr = migdeployNotVolume(migSource, resource)
					checkResourceName = resource.ResourceName

					err := r.setAfterOpenmcpDeploymnet(migSource, i)
					if err != nil {
						omcplog.Error("setBeforeOpenmcpDeploymnet error : ", err)
						r.MakeStatusWithResource(instance, false, migraionSource, resource, err)
						omcplog.Error("Migration Failed")
						return reconcile.Result{Requeue: false}, nil
					}
				} else if resourceType == config.SERVICE {
					migErr = migserviceNotVolume(migSource, resource)
				} else {
					omcplog.Error(fmt.Errorf("Resource Type Error!"))
					r.MakeStatusWithResource(instance, false, migraionSource, resource, fmt.Errorf("Resource Type Error!"))
					omcplog.Error("Migration Failed")
					return reconcile.Result{Requeue: false}, nil
				}
				if migErr != nil {
					omcplog.Error("MigrationControllerResource migration error : ", migErr)
					r.MakeStatusWithResource(instance, false, migraionSource, resource, migErr)
					omcplog.Error("Migration Failed")
					return reconcile.Result{Requeue: false}, nil
				}
			}
		}

		targetListClient := *cm.Cluster_kubeClients[migSource.targetCluster]
		timeoutcheck := 0

		isMigCompleted := false
		omcplog.V(4).Info("connecting... : " + checkResourceName)
		for !isMigCompleted {
			var regexpErr error
			var errMessage string
			podName := ""
			matchResult := false
			pods, _ := targetListClient.CoreV1().Pods(migSource.nameSpace).List(context.TODO(), metav1.ListOptions{})
			for _, pod := range pods.Items {
				matchResult, regexpErr = regexp.MatchString(checkResourceName, pod.Name)

				if regexpErr == nil && matchResult == true {
					if pod.Name != "" {
						podName = pod.Name
					} else {
						podName = checkResourceName + "-unknown"
					}
					omcplog.V(3).Info("TargetCluster PodName : " + podName)
					if pod.Status.Phase == corev1.PodRunning {
						isMigCompleted = true
						omcplog.V(3).Info(podName + " is Running.")
					} else {
						if pod.Status.Reason == "" {
							errMessage = string(pod.Status.Phase)
						} else {
							errMessage = string(pod.Status.Phase) + "/" + pod.Status.Reason + "/" + pod.Status.Message
						}
						omcplog.V(3).Info(errMessage)
					}
				}
			}

			if timeoutcheck == 30 {
				//시간초과 - 오류 루틴으로 진입
				omcplog.V(3).Info("long time error...")
				omcplog.Error(fmt.Errorf("Target Cluster Pod not loaded... : " + podName))
				omcplog.Error(errMessage)
				r.MakeStatusWithMigSource(instance, false, migraionSource, fmt.Errorf("TargetCluster Pod not loaded... : "+checkResourceName))
				omcplog.Error("Migration Failed")
				return reconcile.Result{Requeue: false}, nil
			}

			timeoutcheck = timeoutcheck + 1
			time.Sleep(time.Second * 1)
			omcplog.V(4).Info("connecting...")
		}
	}
	omcplog.V(3).Info("Migration Complete")
	elapsed := time.Since(startDate)
	r.MakeStatus(instance, true, elapsed.String(), nil)

	return reconcile.Result{Requeue: false}, nil
}

func (r *reconciler) setBeforeOpenmcpDeploymnet(migraionSource v1alpha1.MigrationServiceSource, idx int) error {
	odeploy := &resourcev1alpha1.OpenMCPDeployment{}
	err := r.live.Get(context.TODO(), types.NamespacedName{Name: migraionSource.MigrationSources[idx].ResourceName, Namespace: migraionSource.NameSpace}, odeploy)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			omcplog.V(4).Info("setBeforeOpenmcpDeploymnet skip : " + migraionSource.MigrationSources[idx].ResourceName)
			return nil
		} else {
			omcplog.Error("setBeforeOpenmcpDeploymnet error : ", err)
		}
		return err
	}
	omcplog.V(4).Info("--- odeploy status ---")
	omcplog.V(4).Info(odeploy.Status)
	moveCount := odeploy.Status.ClusterMaps[migraionSource.SourceCluster]
	omcplog.V(4).Info(moveCount)
	r.moveCount = moveCount

	odeploy.Status.CheckSubResource = true
	err = r.live.Status().Update(context.TODO(), odeploy)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}
	omcplog.V(4).Info("setBeforeOpenmcpDeploymnet end")
	omcplog.V(4).Info("---------------")

	return nil
}

func (r *reconciler) setAfterOpenmcpDeploymnet(micSource MigrationControllerResource, idx int) error {
	odeploy := &resourcev1alpha1.OpenMCPDeployment{}
	err := r.live.Get(context.TODO(), types.NamespacedName{Name: micSource.resourceList[idx].ResourceName, Namespace: micSource.nameSpace}, odeploy)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			omcplog.V(4).Info("setAfterOpenmcpDeploymnet skip : " + micSource.resourceList[idx].ResourceName)
			return nil
		} else {
			omcplog.Error("setAfterOpenmcpDeploymnet error : ", err)
		}
		return err
	}
	omcplog.V(4).Info("--- odeploy ---")
	omcplog.V(4).Info(odeploy.Status)
	omcplog.V(4).Info(r.moveCount)

	odeploy.Status.CheckSubResource = false
	odeploy.Status.ClusterMaps[micSource.sourceCluster] -= r.moveCount
	odeploy.Status.ClusterMaps[micSource.targetCluster] += r.moveCount

	err = r.live.Status().Update(context.TODO(), odeploy)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}
	omcplog.V(4).Info("setAfterOpenmcpDeploymnet end")
	omcplog.V(4).Info("---------------")

	return nil
}

func (r *reconciler) MakeStatusWithResource(instance *v1alpha1.Migration, migrationStatus bool, migraionSource v1alpha1.MigrationServiceSource, resource v1alpha1.MigrationSource, err error) {
	r.makeStatusRun(instance, migrationStatus, migraionSource, resource, "", err)
}

func (r *reconciler) MakeStatusWithMigSource(instance *v1alpha1.Migration, migrationStatus bool, migraionSource v1alpha1.MigrationServiceSource, err error) {
	r.makeStatusRun(instance, migrationStatus, migraionSource, v1alpha1.MigrationSource{}, "", err)
}

func (r *reconciler) MakeStatus(instance *v1alpha1.Migration, migrationStatus bool, elapsed string, err error) {
	r.makeStatusRun(instance, migrationStatus, v1alpha1.MigrationServiceSource{}, v1alpha1.MigrationSource{}, elapsed, err)
}

func (r *reconciler) makeStatusRun(instance *v1alpha1.Migration, migrationStatus bool, migraionSource v1alpha1.MigrationServiceSource, resource v1alpha1.MigrationSource, elapsedTime string, err error) {
	instance.Status.MigrationStatus = migrationStatus

	if elapsedTime == "" {
		elapsedTime = "0"
	}
	instance.Status.ElapsedTime = elapsedTime
	omcplog.V(3).Info("elapsedTime : ", elapsedTime)

	if !migrationStatus {
		tmp := make(map[string]interface{})
		tmp["SourceCluster"] = migraionSource.SourceCluster
		tmp["TargetCluster"] = migraionSource.TargetCluster
		tmp["ServiceName"] = migraionSource.ServiceName
		tmp["NameSpace"] = migraionSource.NameSpace
		tmp["Reason"] = err.Error()

		jsonTmp, err := json.Marshal(tmp)
		if err != nil {
			omcplog.V(3).Info(err, "-----------")
		}
		instance.Status.Reason = string(jsonTmp)
	}

	//r.live.Update(context.TODO(), instance)
	//r.live.Status().Patch(context.TODO(), instance)
	//r.live.Status().Update(context.TODO(), instance)
	err = r.live.Status().Update(context.TODO(), instance)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}
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
	pods, _ := clusterClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
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
		omcplog.V(4).Info(" Pod hasn't created yet")
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
		omcplog.Error("NewSPDYExecutor err : ", err)
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		timeoutcheck = timeoutcheck + 1
		if timeoutcheck == 10 {
			omcplog.Error("error stream", err)
			return err
		}
		omcplog.V(4).Info("connecting: ", podName)
		time.Sleep(time.Second * 1)
		goto START
	} else {
		//success
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
		omcplog.Error("nfsPv 생성 에러: ", pvErr)
		return false, pvErr, nfsNewPv, nfsNewPvc
	}
	omcplog.V(3).Info("nfsPv 생성 완료")
	pvcErr := client.Create(context.TODO(), nfsNewPvc)
	if pvcErr != nil {
		omcplog.Error("nfsPvc 생성 에러: ", pvcErr)
		return false, pvcErr, nfsNewPv, nfsNewPvc
	}
	omcplog.V(3).Info("nfsPvc 생성 완료")
	//client.Update(context.TODO(), resourceInfo)
	//client.Create(context.TODO(), resourceInfo)
	dpErr := client.Update(context.TODO(), resourceInfo)
	if dpErr != nil {
		omcplog.Error("dp 생성 에러: ", dpErr)
		return false, dpErr, nfsNewPv, nfsNewPvc
	}
	omcplog.V(3).Info("dp 생성 완료")
	return true, nil, nfsNewPv, nfsNewPvc
	// return true, nil
}
func getVolumePath(migSource *MigrationControllerResource) error {
	omcplog.V(0).Info("--- get Volume Path ---")
	//func getVolumePath(migraionSource v1alpha1.MigrationServiceSource, migSource *MigrationControllerResource) error {
	pvcName := ""
	dpResource := &appsv1.Deployment{}
	pvcResource := &corev1.PersistentVolumeClaim{}
	sourceClient := migSource.sourceClient
	volumeName := ""
	mountPath := ""
	for _, resource := range migSource.resourceList {
		if resource.ResourceType == config.PVC {
			pvcName = resource.ResourceName
			err := sourceClient.Get(context.TODO(), pvcResource, migSource.nameSpace, resource.ResourceName)
			if err != nil {
				omcplog.Error("failed get Source PVC :  ", err)
				return err
			}
		} else if resource.ResourceType == config.DEPLOY {
			err := sourceClient.Get(context.TODO(), dpResource, migSource.nameSpace, resource.ResourceName)
			if err != nil {
				omcplog.Error("failed get Source Deploy : ", err)
				return err
			}
		} else {
			continue
		}
	}
	if pvcName == "" {
		omcplog.V(3).Info("MigrationSpec pvcName none")
		return nil
	}
	if dpResource.Name == "" {
		omcplog.V(3).Info("Source dpResourceName none")
		return nil
	}
	if pvcResource.Name == "" {
		omcplog.V(3).Info("Source pvcResourceName none")
		return nil
	}

	volumeList := dpResource.Spec.Template.Spec.Volumes
	for _, volume := range volumeList {
		if volume.PersistentVolumeClaim.ClaimName == pvcName {
			volumeName = volume.Name
		}
	}
	if volumeName == "" {
		omcplog.Error(fmt.Errorf("error volume name nil : dpResource / " + dpResource.String()))
		return fmt.Errorf("error volume name nil : dpResource / " + dpResource.String())
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
		omcplog.Error(fmt.Errorf("error mount path nil : dpResource / " + dpResource.String()))
		return fmt.Errorf("error mount path nil : dpResource / " + dpResource.String())
	}
	storageSize := corev1.ResourceList{}
	if pvcResource.Spec.Resources.Requests.Storage() != nil {
		storageSize = pvcResource.Spec.Resources.Requests
	} else {
		storageSize = corev1.ResourceList{
			corev1.ResourceStorage: resource.MustParse(config.DEFAULT_VOLUME_SIZE),
		}
	}
	omcplog.V(0).Info("mountPath : ")
	omcplog.V(0).Info(mountPath)
	omcplog.V(0).Info("storageSize : ")
	omcplog.V(0).Info(storageSize)
	migSource.volumePath = mountPath
	migSource.resourceRequire = storageSize
	return nil
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
