package snapshotrestore

import (
	"context"
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
	if err := co.WatchResourceReconcileObject(live, &nanumv1alpha1.Snapshot{}, controller.WatchOptions{}); err != nil {
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
	omcplog.V(3).Info("Function Called Reconcile")
	omcplog.V(3).Info(time.Now())

	instance := &nanumv1alpha1.SnapshotRestore{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		omcplog.V(0).Info("get instance error")
	}

	//DATE 추출
	snapshotKey := instance.Spec.SnapshotRestoreSource[0].SnapshotKey
	startTime := util.GetStartTimeBySnapshotKey(snapshotKey)

	for idx, snapshotRestoreSource := range instance.Spec.SnapshotRestoreSource {
		omcplog.V(4).Info(snapshotRestoreSource)

		resourceType := snapshotRestoreSource.ResourceType
		omcplog.V(4).Info("\n[" + strconv.Itoa(idx) + "] : Resource : " + resourceType)
		switch resourceType {
		case config.PV:
			volumeSnapshotRestoreRun(r, &snapshotRestoreSource, startTime)
			fallthrough // 이어서 default 실행
		default:
			etcdSnapshotRestoreRun(r, &snapshotRestoreSource, startTime)
		}
	}

	// 작업 후 업데이트
	updateErr := r.live.Update(context.TODO(), instance, &client.UpdateOptions{})
	if updateErr != nil {
		omcplog.V(3).Info("update error : " + string(startTime))
	}
	return reconcile.Result{}, nil

}
