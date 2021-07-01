package snapshotrestore

import (
	"context"
	"encoding/json"
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
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(3).Info("Snapshot Start : Reconcile")
	startDate := time.Now()
	omcplog.V(3).Info(startDate)

	instance := &nanumv1alpha1.SnapshotRestore{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		omcplog.Error("get instance error : ", err)
		r.MakeStatus(instance, false, "", err)
		return reconcile.Result{Requeue: false}, nil
	}
	// if len(instance.Spec.SnapshotRestoreSource) < 1 {
	// 	omcplog.V(0).Info("========= SnapshotRestoreSource size 0")
	// 	return reconcile.Result{}, nil
	// }

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

	instance.Status.IsVolumeSnapshot = false

	//groupSnapshotKey 추출
	groupSnapshotKey := instance.Spec.GroupSnapshotKey
	if instance.Spec.IsGroupSnapshot {
		setGroupErr, _ := setGroupSnapshotRestoreRun(instance, groupSnapshotKey)
		if setGroupErr != nil {
			omcplog.Error("setGroupSnapshotRestoreRun error : ", setGroupErr)
			r.MakeStatus(instance, false, "", setGroupErr)
			omcplog.V(3).Info("SnapshotRestore Failed")
			return reconcile.Result{Requeue: false}, nil
		}
	} else {
		instance.Status.SnapshotRestoreSource = instance.Spec.DeepCopy().SnapshotRestoreSource
	}

	err = r.live.Update(context.TODO(), instance)
	//err = r.live.Status().Update(context.TODO(), instance)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}

	pvIdx := 1
	for idx, snapshotRestoreSource := range instance.Status.SnapshotRestoreSource {
		resourceType := snapshotRestoreSource.ResourceType
		omcplog.V(4).Info("\n[" + strconv.Itoa(idx) + "] : Resource : " + resourceType)

		if resourceType == config.PV {
			instance.Status.IsVolumeSnapshot = true
			volumeSnapshotKey := util.MakeVolumeProcessResourceKey(groupSnapshotKey, strconv.Itoa(int(startDate.Unix())), strconv.Itoa(pvIdx))
			instance.Status.SnapshotRestoreSource[idx].VolumeSnapshotKey = volumeSnapshotKey
			volumeSnapshotRestoreErr, errDetail := volumeSnapshotRestoreRun(r, snapshotRestoreSource.ResourceCluster, snapshotRestoreSource.ResourceSnapshotKey, groupSnapshotKey, volumeSnapshotKey, pvIdx)
			if volumeSnapshotRestoreErr != nil {
				if errDetail.Error() == "Command Error : TargetFile is empty!" {
					// 타겟이 없는 경우 성공처리
					omcplog.V(3).Info("snapshot volume is zero.")
					instance.Status.SnapshotRestoreSource[idx].VolumeSnapshotKey += "-Empty"
				} else {
					omcplog.Error("volumeSnapshotRestoreRun error : ", volumeSnapshotRestoreErr)
					r.MakeStatusWithSource(instance, false, snapshotRestoreSource, volumeSnapshotRestoreErr, errDetail)
					omcplog.V(3).Info("SnapshotRestore Failed")
					return reconcile.Result{Requeue: false}, nil
				}
			}
			pvIdx++
		}
		etcdSnapshotRestoreErr := etcdSnapshotRestoreRun(r, &snapshotRestoreSource, groupSnapshotKey)
		if etcdSnapshotRestoreErr != nil {
			omcplog.Error("etcdSnapshotRestoreRun error : ", etcdSnapshotRestoreErr)
			r.MakeStatusWithSource(instance, false, snapshotRestoreSource, etcdSnapshotRestoreErr, nil)
			omcplog.V(3).Info("SnapshotRestore Failed")
			return reconcile.Result{Requeue: false}, nil
		}
	}

	omcplog.V(3).Info("Snapshot Restore complete")
	elapsed := time.Since(startDate)
	r.MakeStatus(instance, true, elapsed.String(), nil)

	return reconcile.Result{Requeue: false}, nil
}

func (r *reconciler) MakeStatusWithSource(instance *nanumv1alpha1.SnapshotRestore, status bool, snapshotRestoreSource nanumv1alpha1.SnapshotRestoreSource, err error, detailErr error) {
	r.makeStatusRun(instance, status, snapshotRestoreSource, "", err, detailErr)
}

func (r *reconciler) MakeStatus(instance *nanumv1alpha1.SnapshotRestore, status bool, elapsed string, err error) {
	r.makeStatusRun(instance, status, nanumv1alpha1.SnapshotRestoreSource{}, elapsed, err, nil)
}

func (r *reconciler) makeStatusRun(instance *nanumv1alpha1.SnapshotRestore, status bool, snapshotRestoreSource nanumv1alpha1.SnapshotRestoreSource, elapsedTime string, err error, detailErr error) {
	instance.Status.Status = status

	if elapsedTime == "" {
		elapsedTime = "0"
	}
	instance.Status.ElapsedTime = elapsedTime
	omcplog.V(3).Info("[Exit]")
	omcplog.V(3).Info("snapshotStatus : ", status)

	if !status {
		omcplog.V(3).Info("err : ", err.Error())
		tmp := make(map[string]interface{})
		tmp["Cluster"] = snapshotRestoreSource.ResourceCluster
		//tmp["NameSpace"] = snapshotRestoreSource.ResourceNamespace
		tmp["SnapshotKey"] = snapshotRestoreSource.ResourceSnapshotKey
		tmp["ResourceType"] = snapshotRestoreSource.ResourceType
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
