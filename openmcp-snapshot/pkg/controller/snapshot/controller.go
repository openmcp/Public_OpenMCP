package snapshot

import (
	"context"
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

	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/kubefed/pkg/controller/util"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/client/generic"
	// "openmcp/openmcp/snapshot/pkg/controller"
)

var cm *clusterManager.ClusterManager

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
	live            client.Client
	ghosts          []client.Client
	ghostNamespace  string
	progressMax     int
	progressCurrent int
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(3).Info("Snapshot Start : Reconcile")
	startDate := time.Now()
	omcplog.V(3).Info(startDate)
	instance := &nanumv1alpha1.Snapshot{}
	r.progressMax = 1
	r.progressCurrent = 0
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		omcplog.Error("get instance error : ", err)
		r.makeStatusRun(instance, corev1.ConditionFalse, "0. get instance error", "", err)
		return reconcile.Result{Requeue: false}, nil

	}

	if instance.Status.Status == corev1.ConditionTrue {
		// 이미 성공한 케이스는 로직을 안탄다.
		omcplog.V(4).Info(instance.Name + " already succeed")
		return reconcile.Result{Requeue: false}, nil
	}
	if instance.Status.Status == corev1.ConditionFalse {
		// 이미 실패한 케이스는 로직을 다시 안탄다.
		omcplog.V(4).Info(instance.Name + " already failed")
		return reconcile.Result{Requeue: false}, nil
	}
	omcplog.V(4).Info("0. get instance ... !Update :" + strconv.Itoa(r.progressCurrent))
	r.makeStatusRun(instance, "Running", "0. get resource instance success", "", nil)
	omcplog.V(4).Info("0. get instance ... !Update End")

	//groupSnapshot 키값에 사용될 시간 추출
	//startTime := strconv.Itoa(int(time.Now().Unix()))
	groupSnapshotKey := strconv.Itoa(int(time.Now().Unix()))

	omcplog.V(3).Info(time.Now())
	omcplog.V(4).Info("[Reconcile] startTime : " + groupSnapshotKey)
	instance.Status.IsVolumeSnapshot = false
	instance.Spec.GroupSnapshotKey = groupSnapshotKey
	//instance.Status.SnapshotKey = startTime

	//progress++
	r.progressCurrent++ //getResource 하기 전이라 추가함.
	r.progressMax++
	omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
	omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))
	omcplog.V(4).Info("1. get SnapshotKey ... !Update :" + strconv.Itoa(r.progressCurrent))
	r.makeStatusRun(instance, "Running", "1. get SnapshotKey success", "", nil)
	omcplog.V(4).Info("1. get SnapshotKey ... !Update End")

	// 스냅샷 Resource 정보들을 가공하는 부분
	var initErr error
	instance.Status.SnapshotSources = instance.Spec.DeepCopy().SnapshotSources
	instance.Status.SnapshotSources, initErr = r.setResource(&instance.Status)
	// r.progressMax++ //1번은 최초 reconcile 초기화때 1을 가지고 시작함.
	r.progressMax++ //2번
	//r.progressMax++ //3번은  setResource에서 etcd, volume 개수만큼 진행
	//r.progressMax++ //4번  미구현
	omcplog.V(4).Info("+ [Both] Progress Setting end")
	omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
	omcplog.V(4).Info("++ progressMax add :" + strconv.Itoa(r.progressMax))
	omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))

	if initErr != nil {
		omcplog.Error("2. Init error : ", initErr)
		r.makeStatusRun(instance, corev1.ConditionFalse, "2. Init error", "", initErr)
		omcplog.V(3).Info("Snapshot Failed")
		return reconcile.Result{Requeue: false}, nil
	} else {
		r.progressCurrent++
		omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
		omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))
		omcplog.V(4).Info("2. Init ... !Update :" + strconv.Itoa(r.progressCurrent))
		r.makeStatusRun(instance, "Running", "2. Init success", "", nil)
		omcplog.V(4).Info("2. Init ... !Update End")
	}

	pvIdx := 1
	for idx, snapshotSources := range instance.Status.SnapshotSources {
		desc := ""
		resourceType := snapshotSources.ResourceType
		omcplog.V(4).Info("\n[" + strconv.Itoa(idx) + "] : Resource : " + resourceType)
		omcplog.V(4).Info(snapshotSources)

		// resource Type PV일때 볼륨 스냅샷 진행.
		if resourceType == config.PV {
			instance.Status.IsVolumeSnapshot = true
			volumeSnapshotKey := util.MakeVolumeProcessResourceKey(groupSnapshotKey, strconv.Itoa(int(startDate.Unix())), strconv.Itoa(pvIdx))
			volumeDataSource := &nanumv1alpha1.VolumeDataSource{VolumeSnapshotKey: volumeSnapshotKey}
			instance.Status.SnapshotSources[idx].VolumeDataSource = volumeDataSource
			omcplog.V(4).Info("volumeSnapshotKey : " + volumeSnapshotKey)
			volumeSnapshotErr, errDetail := volumeSnapshotRun(r, &snapshotSources, groupSnapshotKey, volumeSnapshotKey, pvIdx)
			if volumeSnapshotErr != nil {
				if errDetail == nil {
					omcplog.Error("3. volumeSnapshotRun error : ", volumeSnapshotErr)
					r.makeStatusRun(instance, corev1.ConditionFalse, "3. volumeSnapshotRun error", "", volumeSnapshotErr)
					omcplog.V(3).Info("Snapshot Failed")
					return reconcile.Result{Requeue: false}, nil

				} else if errDetail.Error() == "Command Error : TargetFile is empty!" {
					// 타겟이 없은 경우 성공처리
					omcplog.V(3).Info("3. snapshot volume is zero.")
					instance.Status.SnapshotSources[idx].VolumeDataSource.VolumeSnapshotKey += " Empty"
				} else {
					omcplog.Error("3. volumeSnapshotRun error(detail) : ", errDetail)
					r.makeStatusRun(instance, corev1.ConditionFalse, "3. volumeSnapshotRun error(detail)", "", errDetail)
					omcplog.V(3).Info("Snapshot Failed")
					return reconcile.Result{Requeue: false}, nil
				}
			} else {
				r.progressCurrent++
				omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
				omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))
				omcplog.V(4).Info("3. volumeSnapshotRun... !Update :" + strconv.Itoa(r.progressCurrent))
				r.makeStatusRun(instance, "Running", "3. volumeSnapshotRun success : "+desc, "", nil)
				omcplog.V(4).Info("3. volumedSnapshotRun ... !Update End")
			}
			pvIdx++
		}
		etcdSnapshotKeyAllPath, etcdSnapshotErr := etcdSnapshotRun(r, &snapshotSources, groupSnapshotKey)
		instance.Status.SnapshotSources[idx].ResourceSnapshotKey = etcdSnapshotKeyAllPath
		if etcdSnapshotErr != nil {
			omcplog.Error("3. etcdSnapshotRun error : ", etcdSnapshotErr)
			r.makeStatusRun(instance, corev1.ConditionFalse, "3. etcdSnapshotRun error", "", etcdSnapshotErr)
			omcplog.V(3).Info("Snapshot Failed")
			return reconcile.Result{Requeue: false}, nil
		} else {
			r.progressCurrent++
			omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
			omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))
			omcplog.V(4).Info("3. etcdSnapshotRun... !Update :" + strconv.Itoa(r.progressCurrent))
			r.makeStatusRun(instance, "Running", "3. etcdSnapshotRun success : "+desc, "", nil)
			omcplog.V(4).Info("3. etcdSnapshotRun ... !Update End")
		}
	}

	// TODO get Snapshot 정보 추출

	// TODO get Snapshot List 정보 추출.
	r.progressCurrent++
	omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
	omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))
	omcplog.V(3).Info("Snapshot Complete")
	elapsed := time.Since(startDate)
	r.makeStatusRun(instance, corev1.ConditionTrue, "Snapshot succeed", elapsed.String(), nil)

	return reconcile.Result{Requeue: false}, nil

}

// setResource : Deploy
func (r *reconciler) setResource(status *nanumv1alpha1.SnapshotStatus) ([]nanumv1alpha1.SnapshotSource, error) {
	resourceList := status.SnapshotSources

	for _, resource := range resourceList {
		namespace := resource.ResourceNamespace
		resourceName := resource.ResourceName
		sourceClient := cm.Cluster_genClients[resource.ResourceCluster]

		if resource.ResourceType == config.DEPLOY {
			pvcNames := []string{}
			pvNames := []string{}
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
				omcplog.V(0).Info("- 2. check deployment -" + strconv.Itoa(i))

				omcplog.V(0).Info(volume)
				pvcInfo := volume.PersistentVolumeClaim
				if pvcInfo == nil {
					continue
				} else {
					omcplog.V(0).Info("--- pvc contain ---")
					omcplog.V(0).Info(pvcInfo.ClaimName)
					pvcNames = append(pvcNames, pvcInfo.ClaimName)

					// 3. PVC 정보를 토대로 PV 라벨 정보 추출
					omcplog.V(0).Info("- 3. check pvc -")
					sourcePVC := &corev1.PersistentVolumeClaim{}
					err := sourceClient.Get(context.TODO(), sourcePVC, namespace, pvcInfo.ClaimName)
					if err != nil {
						omcplog.V(0).Info(pvcInfo.ClaimName + " is not exist!!!")
						continue
					}
					omcplog.V(0).Info(sourcePVC.Name)
					pvc_matchLabel := sourcePVC.Spec.Selector.DeepCopy().MatchLabels

					omcplog.V(0).Info("pvc_matchLabel :")
					omcplog.V(0).Info(pvc_matchLabel)

					// 4. PVC에서 추출된 라벨 정보로 PV 추출
					sourcePVList := &corev1.PersistentVolumeList{}
					err = sourceClient.List(context.TODO(), sourcePVList, namespace, &client.ListOptions{
						LabelSelector: labels.SelectorFromSet(labels.Set(pvc_matchLabel)),
					})
					if err != nil {
						omcplog.V(0).Info(pvc_matchLabel)
						omcplog.V(0).Info(" this label not exist!!!")
						continue
					}
					for _, pv := range sourcePVList.Items {
						omcplog.V(0).Info("- 4. check pv -")
						omcplog.V(0).Info(pv.Name)
						pvNames = append(pvNames, pv.Name)
					}
				}
			}

			// 5. 저장된 pvName, pvcName 을 토대로 리스트 재작성.
			for _, pvName := range pvNames {
				tmp := nanumv1alpha1.SnapshotSource{}
				tmp.ResourceName = pvName
				tmp.ResourceType = config.PV
				tmp.ResourceNamespace = "default"
				tmp.ResourceCluster = resource.ResourceCluster
				resourceList = append(resourceList, tmp)
			}
			for _, pvcName := range pvcNames {
				tmp := nanumv1alpha1.SnapshotSource{}
				tmp.ResourceName = pvcName
				tmp.ResourceType = config.PVC
				tmp.ResourceNamespace = resource.ResourceNamespace
				tmp.ResourceCluster = resource.ResourceCluster
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
	//Deploy를 가장 뒤로.
	for i, fixedResource := range fixedResourceList {
		if fixedResource.ResourceType == config.DEPLOY {
			tmp := fixedResourceList[i]
			fixedResourceList[i] = fixedResourceList[len(fixedResourceList)-1]
			fixedResourceList[len(fixedResourceList)-1] = tmp
		}
	}

	omcplog.V(0).Info("->")
	omcplog.V(0).Info(fixedResourceList)
	omcplog.V(0).Info("--------------------")

	omcplog.V(0).Info("----check resource........----")
	//deploy의 볼륨개수, service pv,pvc 등의 개수에 따라 수정.
	for i, fixedResource := range fixedResourceList {
		if fixedResource.ResourceType == config.DEPLOY {
			tmp := fixedResourceList[i]
			fixedResourceList[i] = fixedResourceList[len(fixedResourceList)-1]
			fixedResourceList[len(fixedResourceList)-1] = tmp
		} else if fixedResource.ResourceType == config.PV {
			//PV 가 있으면 한번 더. 3. snapshot 진행시 VolumeSnapshot 에 대한 내용을 한번 더함.
			r.progressMax++
		}
		//etcd 스냅샷에 대한 내용.
		r.progressMax++
	}
	omcplog.V(4).Info("+ ProgressMax Setting end")
	omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
	omcplog.V(4).Info("++ progressMax add :" + strconv.Itoa(r.progressMax))
	omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))

	resourceList = fixedResourceList
	return resourceList, nil
}
