package snapshot

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"openmcp/openmcp/apis"
	nanumv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"
	config "openmcp/openmcp/openmcp-snapshot/pkg/util"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"

	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/kubefed/pkg/controller/util"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/client/generic"
	// "openmcp/openmcp/snapshot/pkg/controller"
)

var cm *clusterManager.ClusterManager
var namespacedName types.NamespacedName

//pod 이름 찾기
func GetPodName(targetClient generic.Client, dpName string, namespace string) (string, error) {
	podInfo := &corev1.Pod{}

	listOption := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{
			"name": dpName,
		}),
	}

	err := targetClient.List(context.TODO(), podInfo, namespace, listOption)
	if err != nil {
		omcplog.Error("----------- : ", err)
		return "", err
	}

	podName := podInfo.ObjectMeta.Name

	return podName, nil
}

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	cm = myClusterManager
	omcplog.V(4).Info("NewController start")
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		omcplog.Error("getting delegating client for live cluster: ", err)
		return nil, err
	}
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			omcplog.Error("getting delegating client for ghost cluster: ", err)
			return nil, err
		}
		ghostclients = append(ghostclients, ghostclient)
	}
	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		omcplog.Error("adding APIs to live cluster's scheme: ", err)
		return nil, err
	}
	if err := co.WatchResourceReconcileObject(context.TODO(), live, &nanumv1alpha1.Snapshot{}, controller.WatchOptions{}); err != nil {
		omcplog.Error("setting up Pod watch in live cluster: ", err)
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

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(3).Info("Snapshot Start : Reconcile")
	startDate := time.Now()
	omcplog.V(3).Info(startDate)

	instance := &nanumv1alpha1.Snapshot{}
	namespacedName = req.NamespacedName
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		omcplog.Error("get instance error : ", err)
		r.MakeStatus(instance, false, "", err)
		return reconcile.Result{Requeue: false}, nil
	}

	if instance.Status.Status == true {
		// 이미 성공한 케이스는 로직을 안탄다.
		omcplog.V(4).Info(instance.Name + " already succeed")
		return reconcile.Result{Requeue: false}, nil
	}
	if instance.Status.Status == false && instance.Status.Reason != "" {
		// 이미 실패한 케이스는 로직을 다시 안탄다.
		omcplog.V(4).Info(instance.Name + " already failed")
		return reconcile.Result{Requeue: false}, nil
	}

	//groupSnapshot 키값에 사용될 시간 추출
	//startTime := strconv.Itoa(int(time.Now().Unix()))
	groupSnapshotKey := strconv.Itoa(int(time.Now().Unix()))

	omcplog.V(3).Info(time.Now())
	omcplog.V(4).Info("[Reconcile] startTime : " + groupSnapshotKey)
	instance.Status.IsVolumeSnapshot = false
	instance.Spec.GroupSnapshotKey = groupSnapshotKey
	//instance.Status.SnapshotKey = startTime

	// 스냅샷 Resource 정보들을 가공하는 부분
	instance.Status.SnapshotSources = instance.Spec.DeepCopy().SnapshotSources
	setResourceForOnlyDeploy(&instance.Status)

	pvIdx := 1
	for idx, snapshotSources := range instance.Status.SnapshotSources {
		resourceType := snapshotSources.ResourceType
		omcplog.V(4).Info("\n[" + strconv.Itoa(idx) + "] : Resource : " + resourceType)

		if resourceType == config.PV {
			instance.Status.IsVolumeSnapshot = true
			volumeSnapshotKey := util.MakeVolumeProcessResourceKey(groupSnapshotKey, strconv.Itoa(int(startDate.Unix())), strconv.Itoa(pvIdx))
			volumeDataSource := &nanumv1alpha1.VolumeDataSource{VolumeSnapshotKey: volumeSnapshotKey}
			instance.Status.SnapshotSources[idx].VolumeDataSource = volumeDataSource
			omcplog.V(4).Info("volumeSnapshotKey : " + volumeSnapshotKey)
			volumeSnapshotErr, errDetail := volumeSnapshotRun(r, &snapshotSources, groupSnapshotKey, volumeSnapshotKey, pvIdx)
			if volumeSnapshotErr != nil {
				if errDetail == nil {
					omcplog.Error("volumeSnapshotRun error : ", volumeSnapshotErr)
					r.MakeStatusWithSource(instance, false, snapshotSources, volumeSnapshotErr, errDetail)
					omcplog.V(3).Info("Snapshot Failed")
					return reconcile.Result{Requeue: false}, nil
				} else if errDetail.Error() == "Command Error : TargetFile is empty!" {
					// 타겟이 없은 경우 성공처리
					omcplog.V(3).Info("snapshot volume is zero.")
					instance.Status.SnapshotSources[idx].VolumeDataSource.VolumeSnapshotKey += " Empty"
				} else {
					omcplog.Error("volumeSnapshotRun error : ", errDetail)
					r.MakeStatusWithSource(instance, false, snapshotSources, volumeSnapshotErr, errDetail)
					omcplog.V(3).Info("Snapshot Failed")
					return reconcile.Result{Requeue: false}, nil
				}
			}
			pvIdx++
		}
		etcdSnapshotKeyAllPath, etcdSnapshotErr := etcdSnapshotRun(r, &snapshotSources, groupSnapshotKey)
		instance.Status.SnapshotSources[idx].ResourceSnapshotKey = etcdSnapshotKeyAllPath
		if etcdSnapshotErr != nil {
			omcplog.Error("etcdSnapshotRun error : ", etcdSnapshotErr)
			r.MakeStatusWithSource(instance, false, snapshotSources, etcdSnapshotErr, nil)
			omcplog.V(3).Info("Snapshot Failed")
			return reconcile.Result{Requeue: false}, nil
		}
	}
	omcplog.V(3).Info("Snapshot Complete")
	elapsed := time.Since(startDate)
	r.MakeStatus(instance, true, elapsed.String(), nil)

	return reconcile.Result{Requeue: false}, nil

}

func (r *reconciler) MakeStatusWithSource(instance *nanumv1alpha1.Snapshot, snapshotStatus bool, snapshotSource nanumv1alpha1.SnapshotSource, err error, detailErr error) {
	r.makeStatusRun(instance, snapshotStatus, snapshotSource, "", err, detailErr)
}

func (r *reconciler) MakeStatus(instance *nanumv1alpha1.Snapshot, snapshotStatus bool, elapsed string, err error) {
	r.makeStatusRun(instance, snapshotStatus, nanumv1alpha1.SnapshotSource{}, elapsed, err, nil)
}

func (r *reconciler) makeStatusRun(instance *nanumv1alpha1.Snapshot, snapshotStatus bool, snapshotSource nanumv1alpha1.SnapshotSource, elapsedTime string, err error, detailErr error) {
	instance.Status.Status = snapshotStatus

	if elapsedTime == "" {
		elapsedTime = "0"
	}
	instance.Status.ElapsedTime = elapsedTime
	omcplog.V(3).Info("[Exit]")
	omcplog.V(3).Info("snapshotStatus : ", snapshotStatus)

	if !snapshotStatus {
		omcplog.V(3).Info("err : ", err.Error())
		tmp := make(map[string]interface{})
		tmp["Cluster"] = snapshotSource.ResourceCluster
		tmp["NameSpace"] = snapshotSource.ResourceNamespace
		//tmp["ResourceType"] = snapshotSource.ResourceType
		//tmp["ResourceName"] = snapshotSource.ResourceName
		tmp["GroupSnapshotKey"] = instance.Spec.GroupSnapshotKey
		//tmp["VolumeSnapshotClassName"] = snapshotSource.VolumeDataSource.VolumeSnapshotClassName
		//tmp["VolumeSnapshotSourceKind"] = snapshotSource.VolumeDataSource.VolumeSnapshotSourceKind
		//tmp["VolumeSnapshotSourceName"] = snapshotSource.VolumeDataSource.VolumeSnapshotSourceName
		//tmp["VolumeSnapshotKey"] = instance.Status.VolumeDataSource.VolumeSnapshotKey
		tmp["Reason"] = err.Error()
		//tmp["ReasonDetail"] = detailErr.Error()

		jsonTmp, err := json.Marshal(tmp)
		if err != nil {
			omcplog.V(3).Info(err, "-----------")
		}
		instance.Status.Reason = string(jsonTmp)
		if detailErr != nil {
			instance.Status.ReasonDetail = detailErr.Error()
		}
	}

	//r.live.Update(context.TODO(), instance)
	//r.live.Status().Patch(context.TODO(), instance)
	//r.live.Status().Update(context.TODO(), instance)
	//err = r.live.Status().Update(context.TODO(), instance)
	omcplog.V(3).Info("live update")
	err = r.live.Update(context.Background(), instance)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}
	err = r.live.Status().Update(context.Background(), instance)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}
	time.Sleep(5 * time.Second)

	omcplog.V(3).Info("live update end")
}

// setResource : Deploy
func setResourceForOnlyDeploy(status *nanumv1alpha1.SnapshotStatus) error {
	resourceList := status.SnapshotSources

	for _, resource := range resourceList {
		namespace := resource.ResourceNamespace
		resourceName := resource.ResourceName
		sourceClient := cm.Cluster_genClients[resource.ResourceCluster]
		pvcNames := []string{}
		pvNames := []string{}

		if resource.ResourceType == config.DEPLOY {
			// 1. Deployment 데이터를 가져온다.
			omcplog.V(0).Info("- 1. getDeployment -")
			omcplog.V(0).Info(resourceName)

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
				tmp := nanumv1alpha1.SnapshotSource{}
				tmp.ResourceName = pvName
				tmp.ResourceType = config.PV
				resourceList = append(resourceList, tmp)
			}
			for _, pvcName := range pvcNames {
				tmp := nanumv1alpha1.SnapshotSource{}
				tmp.ResourceName = pvcName
				tmp.ResourceType = config.PVC
				resourceList = append(resourceList, tmp)
			}
		}
	}

	omcplog.V(0).Info("- resourceList fix -")
	omcplog.V(0).Info(resourceList)

	// 6. 기존에 하던 데이터 보정 실행 (동일 객체 제거, Deployment를 제일 앞으로)
	fixedResourceList := []nanumv1alpha1.SnapshotSource{}
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

	resourceList = fixedResourceList
	return nil
}
