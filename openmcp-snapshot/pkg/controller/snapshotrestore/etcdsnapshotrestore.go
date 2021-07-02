package snapshotrestore

import (

	// "openmcp/openmcp/migration/pkg/apis"
	"context"
	"fmt"
	nanumv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-snapshot/pkg/controller/snapshotrestore/resources"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"
	"openmcp/openmcp/openmcp-snapshot/pkg/util/etcd"
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
	// "openmcp/openmcp/migration/pkg/controller"
)

//volumeSnapshotRun 내에는 PV 만 들어온다고 가정한다.
func etcdSnapshotRestoreRun(r *reconciler, snapshotRestoreSource *nanumv1alpha1.SnapshotRestoreSource, groupSnapshotKey string) error {
	omcplog.V(3).Info(snapshotRestoreSource)

	omcplog.V(3).Info("etcd snapshot restore start")
	resourceJSONString, getErr := getETCDResource(groupSnapshotKey, snapshotRestoreSource.ResourceSnapshotKey)
	if getErr != nil {
		omcplog.Error("getETCDResource")
		return getErr
	}
	omcplog.V(3).Info("resource json size : " + strconv.Itoa(len(resourceJSONString)))
	//Client 로 Deploy
	err := CreateResourceJSON(snapshotRestoreSource, resourceJSONString)
	if err != nil {
		omcplog.Error("etcdsnapshotrestore.go : GetResourceJSON for cluster error")
		return err
	}
	return nil
}

func getETCDResource(startTime string, resourceSnapshotKey string) (string, error) {
	omcplog.V(4).Info("get Etcd key AllPath : " + resourceSnapshotKey)
	//ETCD 에서 데이터 가져오기.
	etcdCtl, etcdInitErr := etcd.InitEtcd()
	if etcdInitErr != nil {
		omcplog.Error("etcdsnapshotresource.go : Etcd Init Err")
		return "", etcdInitErr
	}
	resp, etcdGetErr := etcdCtl.Get(resourceSnapshotKey)
	if etcdGetErr != nil {
		omcplog.Error("etcdsnapshotresource.go : Etcd Get Err")
		return "", etcdGetErr
	}
	resourceJSONString := string(resp.Kvs[0].Value)
	return resourceJSONString, nil
}

// CreateResourceJSON : https://mingrammer.com/gobyexample/json/ 를 참조하여 작성
func CreateResourceJSON(snapshotSource *nanumv1alpha1.SnapshotRestoreSource, resourceJSONString string) error {

	var resourceObj runtime.Object

	client := cm.Cluster_genClients[snapshotSource.ResourceCluster]
	//clientconfig := cm.Cluster_configs
	//omcplog.V(3).Info("----------", clientconfig)

	switch snapshotSource.ResourceType {
	case util.DEPLOY:
		resourceObj, _ = resources.JSON2Deploy(resourceJSONString)
	case util.SERVICE:
		resourceObj, _ = resources.JSON2Service(resourceJSONString, client)
	case util.PVC:
		resourceObj, _ = resources.JSON2Pvc(resourceJSONString)
	case util.PV:
		resourceObj, _ = resources.JSON2Pv(resourceJSONString)
	default:
		omcplog.Error("Invalid resourceType")
		return fmt.Errorf("Invalid resourceType : " + snapshotSource.ResourceType)
	}

	err := client.Update(context.TODO(), resourceObj)
	if err != nil {
		omcplog.Error("UpdateResource for JSON error", err)
		return err
	}
	omcplog.V(3).Info("Update " + snapshotSource.ResourceType + "resource complete!")
	return nil
}
