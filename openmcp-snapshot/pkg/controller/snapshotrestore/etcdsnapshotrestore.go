package snapshotrestore

import (

	// "openmcp/openmcp/migration/pkg/apis"

	"context"
	nanumv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-snapshot/pkg/controller/snapshotrestore/resources"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"
	"openmcp/openmcp/openmcp-snapshot/pkg/util/etcd"

	"k8s.io/apimachinery/pkg/runtime"
	// "openmcp/openmcp/migration/pkg/controller"
)

//volumeSnapshotRun 내에는 PV 만 들어온다고 가정한다.
func etcdSnapshotRestoreRun(r *reconciler, snapshotRestoreSource *nanumv1alpha1.SnapshotRestoreSource, startTime string) error {
	omcplog.V(4).Info(snapshotRestoreSource)

	snapshotKey := snapshotRestoreSource.SnapshotKey

	//ETCD 에서 데이터 가져오기.
	etcdCtl := etcd.InitEtcd()
	resp := etcdCtl.Get(snapshotKey)
	resourceJSONString := string(resp.Kvs[0].Value)

	//Client 로 Deploy
	isSuccess, err := CreateResourceJSON(snapshotRestoreSource, resourceJSONString)
	if err != nil {
		omcplog.V(2).Info("CreateResource for JSON error")
	} else if !isSuccess {
		omcplog.V(2).Info("CreateResource for JSON error")
	}
	return nil
}

// CreateResourceJSON : https://mingrammer.com/gobyexample/json/ 를 참조하여 작성
func CreateResourceJSON(snapshotSource *nanumv1alpha1.SnapshotRestoreSource, resourceJSONString string) (bool, error) {

	var resourceObj runtime.Object

	client := cm.Cluster_genClients[snapshotSource.ResourceCluster]

	switch snapshotSource.ResourceType {
	case util.DEPLOY:
		resourceObj, _ = resources.JSON2Deploy(resourceJSONString)
	case util.SERVICE:
		resourceObj, _ = resources.JSON2Service(resourceJSONString)
	case util.PVC:
		resourceObj, _ = resources.JSON2Pvc(resourceJSONString)
	case util.PV:
		resourceObj, _ = resources.JSON2Pv(resourceJSONString)
	default:
		omcplog.V(2).Info("Invalid resourceType")
	}
	client.Create(context.TODO(), resourceObj)

	return true, nil
}
