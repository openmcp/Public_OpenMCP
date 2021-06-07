package snapshot

import (

	// "openmcp/openmcp/migration/pkg/apis"

	"context"
	"encoding/json"
	nanumv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"
	"openmcp/openmcp/openmcp-snapshot/pkg/util/etcd"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	// "openmcp/openmcp/migration/pkg/controller"
)

//volumeSnapshotRun 내에는 PV 만 들어온다고 가정한다.
func etcdSnapshotRun(r *reconciler, snapshotSource *nanumv1alpha1.SnapshotSource, startTime string) error {
	omcplog.V(4).Info(snapshotSource)

	omcplog.V(3).Info("snapshot start")
	snapshotKey := util.MakeSnapshotKeyForSnapshot(startTime, snapshotSource)
	snapshotKeyAllPath := util.MakeSnapshotKeyAllPath(startTime, snapshotKey)

	//Client 로 데이터 가져오기.
	resourceJSONString, err := GetResourceJSON(snapshotSource)
	if err != nil {
		omcplog.V(2).Info("GetResourceJSON for cluster error")
	}

	omcplog.V(2).Info("Input ETCD")
	//ETCD 에 삽입
	etcdCtl := etcd.InitEtcd()
	_ = etcdCtl.Put(snapshotKeyAllPath, resourceJSONString)
	//snapshotSource.SnapshotKey = snapshotKey
	omcplog.V(2).Info("Input ETCD end")
	return nil
}

// GetResourceJSON : https://mingrammer.com/gobyexample/json/ 를 참조하여 작성
func GetResourceJSON(snapshotSource *nanumv1alpha1.SnapshotSource) (string, error) {

	var resourceObj runtime.Object

	client := cm.Cluster_genClients[snapshotSource.ResourceCluster]

	switch snapshotSource.ResourceType {
	case util.DEPLOY:
		resourceObj = &appsv1.Deployment{}
	case util.SERVICE:
		resourceObj = &corev1.Service{}
	case util.PVC:
		resourceObj = &corev1.PersistentVolumeClaim{}
	case util.PV:
		resourceObj = &corev1.PersistentVolume{}
	default:
		omcplog.V(2).Info("Invalid resourceType")
	}
	client.Get(context.TODO(), resourceObj, snapshotSource.ResourceNamespace, snapshotSource.ResourceName)

	omcplog.V(3).Info("resourceType : " + snapshotSource.ResourceType + ", resourceName : " + snapshotSource.ResourceName + ", resourceNamespace: " + snapshotSource.ResourceNamespace)
	ret, err := obj2JsonString(resourceObj)
	if err != nil {
		omcplog.V(2).Info("Json Convert Error")
	}
	return ret, nil
}

// Obj2JsonString : Deployment 등과 같은 interface 를 json string 으로 변환.
func obj2JsonString(obj interface{}) (string, error) {

	json, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	omcplog.V(3).Info("===Obj2JsonString===")
	omcplog.V(3).Info(string(json)[0:40] + "...")

	return string(json), nil
}
