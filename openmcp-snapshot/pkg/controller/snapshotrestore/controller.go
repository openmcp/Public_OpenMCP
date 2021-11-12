package snapshotrestore

import (
	"context"
	"strconv"
	"time"

	nanumv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"
	config "openmcp/openmcp/openmcp-snapshot/pkg/util"
	"openmcp/openmcp/openmcp-snapshot/pkg/util/etcd"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"

	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/apis"
	"sigs.k8s.io/kubefed/pkg/client/generic"
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
	if err := co.WatchResourceReconcileObject(context.TODO(), live, &nanumv1alpha1.SnapshotRestore{}, controller.WatchOptions{}); err != nil {
		omcplog.V(0).Info("setting up Pod watch in live cluster: ", err)
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
	instance := &nanumv1alpha1.SnapshotRestore{}
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

	instance.Status.IsVolumeSnapshot = false

	//groupSnapshotKey 추출
	groupSnapshotKey := instance.Spec.GroupSnapshotKey
	if instance.Spec.IsGroupSnapshot {
		setGroupErr, _ := setGroupSnapshotRestoreRun(instance, groupSnapshotKey)
		if setGroupErr != nil {
			omcplog.Error("1. get SnapshotInfo for ETCD error : ", setGroupErr)
			r.makeStatusRun(instance, "Running", "1. get SnapshotInfo for ETCD Failed", "", setGroupErr)
			omcplog.V(3).Info("1. get SnapshotInfo for ETCD Failed")
			return reconcile.Result{Requeue: false}, nil
		}
	} else {
		instance.Status.SnapshotRestoreSource = instance.Spec.DeepCopy().SnapshotRestoreSource

		//progress++
		r.progressCurrent++ //getResource 하기 전이라 추가함.
		r.progressMax++
		omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
		omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))
		omcplog.V(4).Info("1. get SnapshotInfo for ETCD ... !Update :" + strconv.Itoa(r.progressCurrent))
		r.makeStatusRun(instance, "Running", "1. get SnapshotInfo for ETCD success", "", nil)
		omcplog.V(4).Info("1. get SnapshotInfo for ETCD ... !Update End")
	}

	// 스냅샷 복구 정보들을가공하는 부분 -> r.progressMax 추가 산출 함수.
	initErr := r.setResource(instance.Status.SnapshotRestoreSource)
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
		omcplog.V(3).Info("SnapshotRestore Failed")
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
	for idx, snapshotRestoreSource := range instance.Status.SnapshotRestoreSource {
		desc := ""
		resourceType := snapshotRestoreSource.ResourceType
		omcplog.V(4).Info("\n[" + strconv.Itoa(idx) + "] : Resource : " + resourceType)

		if resourceType == config.PV {
			instance.Status.IsVolumeSnapshot = true
			volumeSnapshotKey := util.MakeVolumeProcessResourceKey(groupSnapshotKey, strconv.Itoa(int(startDate.Unix())), strconv.Itoa(pvIdx))
			instance.Status.SnapshotRestoreSource[idx].VolumeSnapshotKey = volumeSnapshotKey
			volumeSnapshotRestoreErr, errDetail := volumeSnapshotRestoreRun(r, snapshotRestoreSource.ResourceCluster, snapshotRestoreSource.ResourceSnapshotKey, groupSnapshotKey, volumeSnapshotKey, pvIdx)
			if volumeSnapshotRestoreErr != nil {
				if errDetail == nil {
					omcplog.Error("3. volumeSnapshotRestoreRun error : ", volumeSnapshotRestoreErr)
					r.makeStatusRun(instance, corev1.ConditionFalse, "3. volumeSnapshotRestoreRun error", "", volumeSnapshotRestoreErr)
					omcplog.V(3).Info("SnapshotRestore Failed")
					return reconcile.Result{Requeue: false}, nil

				} else if errDetail.Error() == "Command Error : TargetFile is empty!" {
					// 타겟이 없는 경우 성공처리
					omcplog.V(3).Info("3. snapshot volume is zero.")
					instance.Status.SnapshotRestoreSource[idx].VolumeSnapshotKey += "-Empty"
				} else {
					omcplog.Error("3. volumeSnapshotRestoreRun error(detail) : ", errDetail)
					r.makeStatusRun(instance, corev1.ConditionFalse, "3. volumeSnapshotRestoreRun error(detail)", "", errDetail)
					omcplog.V(3).Info("SnapshotRestore Failed")
					return reconcile.Result{Requeue: false}, nil
				}
			} else {
				r.progressCurrent++
				omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
				omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))
				omcplog.V(4).Info("3. volumeSnapshotRestoreRun... !Update :" + strconv.Itoa(r.progressCurrent))
				r.makeStatusRun(instance, "Running", "3. volumeSnapshotRestoreRun success : "+desc, "", nil)
				omcplog.V(4).Info("3. volumeSnapshotRestoreRun ... !Update End")
			}
			pvIdx++
		}
		etcdSnapshotRestoreErr := etcdSnapshotRestoreRun(r, &snapshotRestoreSource, groupSnapshotKey)
		if etcdSnapshotRestoreErr != nil {
			omcplog.Error("etcdSnapshotRestoreRun error : ", etcdSnapshotRestoreErr)
			r.makeStatusRun(instance, corev1.ConditionFalse, "3. etcdSnapshotRestoreRun error", "", etcdSnapshotRestoreErr)
			omcplog.V(3).Info("SnapshotRestore Failed")
			return reconcile.Result{Requeue: false}, nil
		} else {
			r.progressCurrent++
			omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
			omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))
			omcplog.V(4).Info("3. etcdSnapshotRestoreRun... !Update :" + strconv.Itoa(r.progressCurrent))
			r.makeStatusRun(instance, "Running", "3. etcdSnapshotRestoreRun success : "+desc, "", nil)
			omcplog.V(4).Info("3. etcdSnapshotRestoreRun ... !Update End")
		}
	}

	r.progressCurrent++
	omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
	omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))
	omcplog.V(3).Info("Snapshot Restore complete")
	elapsed := time.Since(startDate)
	r.makeStatusRun(instance, corev1.ConditionTrue, "Snapshot succeed", elapsed.String(), nil)

	return reconcile.Result{Requeue: false}, nil
}

//setGroupSnapshotRestoreRun 실행
func setGroupSnapshotRestoreRun(instance *nanumv1alpha1.SnapshotRestore, groupSnapshotKey string) (error, error) {
	omcplog.V(3).Info(instance)
	omcplog.V(3).Info("groupSnapshot is true")

	//ETCD 에서 데이터 가져오기.
	etcdCtl, etcdInitErr := etcd.InitEtcd()
	if etcdInitErr != nil {
		return etcdInitErr, nil
	}
	resp, etcdGetErr := etcdCtl.GetEtcdGroupSnapshot(instance.Spec.GroupSnapshotKey)
	if etcdGetErr != nil {
		return etcdGetErr, nil
	}

	omcplog.V(3).Info("1) set SnapshotRestoreList")
	snapshotRestoreSources := []nanumv1alpha1.SnapshotRestoreSource{}

	for idx, kv := range resp.Kvs {
		//snapshot Restore Source 에 내용을 추가한다. namespace 의 경우 알 수 없지만, 알 수 없어도 동작하게끔 구현되어 있다.

		snapshotRestoreSource := nanumv1alpha1.SnapshotRestoreSource{}
		ResourceCluster, _ := util.GetClusterNameBySnapshotKey(string(kv.Key))
		ResourceSnapshotKey, _ := util.GetGroupSnapshotKeyBySnapshotKey(string(kv.Key))
		ResourceType, _ := util.GetResourceTypeBySnapshotKey(string(kv.Key))
		snapshotRestoreSource.ResourceCluster = ResourceCluster
		//snapshotRestoreSource.ResourceNamespace = ResourceCluster
		snapshotRestoreSource.ResourceSnapshotKey = ResourceSnapshotKey
		snapshotRestoreSource.ResourceType = ResourceType
		snapshotRestoreSource.ResourceSnapshotKey = string(kv.Key)

		omcplog.V(3).Info("idx : " + strconv.Itoa(idx))
		omcplog.V(3).Info(snapshotRestoreSource)
		snapshotRestoreSources = append(snapshotRestoreSources, snapshotRestoreSource)
	}
	omcplog.V(3).Info("  set SnapshotRestoreList end")

	instance.Status.SnapshotRestoreSource = snapshotRestoreSources
	return nil, nil
}

// setResource : Deploy
func (r *reconciler) setResource(resourceList []nanumv1alpha1.SnapshotRestoreSource) error {

	for _, resource := range resourceList {
		if resource.ResourceType == config.PV {
			r.progressMax++
		}
		r.progressMax++
	}

	omcplog.V(4).Info("+ ProgressMax Setting end")
	omcplog.V(4).Info("++ progressCurrent add :" + strconv.Itoa(r.progressCurrent))
	omcplog.V(4).Info("++ progressMax add :" + strconv.Itoa(r.progressMax))
	omcplog.V(4).Info("++ progress Count : " + strconv.Itoa(r.progressCurrent) + "/" + strconv.Itoa(r.progressMax))

	return nil
}
