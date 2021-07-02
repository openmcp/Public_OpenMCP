package snapshot

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	nanumv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"
	config "openmcp/openmcp/openmcp-snapshot/pkg/util"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"

	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/kubefed/pkg/controller/util"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/apis"
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
	omcplog.V(4).Info("-----------")
	omcplog.V(4).Info(live)
	omcplog.V(4).Info("-----------")
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
	pvIdx := 1
	for idx, snapshotSources := range instance.Spec.SnapshotSources {
		resourceType := snapshotSources.ResourceType
		omcplog.V(4).Info("\n[" + strconv.Itoa(idx) + "] : Resource : " + resourceType)

		if resourceType == config.PV {
			instance.Status.IsVolumeSnapshot = true
			volumeSnapshotKey := util.MakeVolumeProcessResourceKey(groupSnapshotKey, strconv.Itoa(int(startDate.Unix())), strconv.Itoa(pvIdx))
			volumeDataSource := &nanumv1alpha1.VolumeDataSource{VolumeSnapshotKey: volumeSnapshotKey}
			instance.Spec.SnapshotSources[idx].VolumeDataSource = volumeDataSource
			omcplog.V(4).Info("volumeSnapshotKey : " + volumeSnapshotKey)
			volumeSnapshotErr, errDetail := volumeSnapshotRun(r, &snapshotSources, groupSnapshotKey, volumeSnapshotKey, pvIdx)
			if volumeSnapshotErr != nil {
				if errDetail.Error() == "Command Error : TargetFile is empty!" {
					// 타겟이 없은 경우 성공처리
					omcplog.V(3).Info("snapshot volume is zero.")
					instance.Spec.SnapshotSources[idx].VolumeDataSource.VolumeSnapshotKey += " Empty"
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
		instance.Spec.SnapshotSources[idx].ResourceSnapshotKey = etcdSnapshotKeyAllPath
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
