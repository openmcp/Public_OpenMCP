package util

import (
	"encoding/json"
	"time"

	snapshotv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type etcdInfoParamsMap struct {
	DialTimeout time.Duration
	Endpoints   []string
}

const (

	// RequestTimeout 는  context 생성시의 timeout 기간입니다. 아래와 같이 사용합니다.
	// ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	RequestTimeout = 10 * time.Second

	// CSI Parameters prefixed with csiParameterPrefix are not passed through
	// to the driver on CreateSnapshotRequest calls. Instead they are intended
	// to be used by the CSI external-snapshotter and maybe used to populate
	// fields in subsequent CSI calls or Kubernetes API objects.
	csiParameterPrefix = "csi.storage.k8s.io/"

	// CSISnapshotPrefix 는 스냅샷의 이름에 붙는 header 이며, SnapshotPrefix + "-" + PVC Name 으로 이루어진다.
	CSISnapshotPrefix = "VolumeSnapshot"
	// CRDNameSpace 는 스냅샷 등의 namespace 지정
	CRDNameSpace = "default"
)

// EtcdInfo 는 ETCD 접속정보 입니다.
var EtcdInfo = etcdInfoParamsMap{
	//dial
	DialTimeout: 2 * time.Second,
	Endpoints:   []string{"10.0.0.221:12379"},
}

// GetVolumeSnapshotName 는 GetVolumeSnapshotName 을 구할 때 쓰는 함수. (생성)
func GetVolumeSnapshotName(volumeDataSource snapshotv1alpha1.VolumeDataSource) string {
	return CSISnapshotPrefix + "-" + volumeDataSource.VolumeSnapshotSourceName
}

// convertResourceObj : json String 을 obj 로 변환
func convertResourceObj(resourceInfoJSON string) (*unstructured.Unstructured, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var unstructured *unstructured.Unstructured
	jsonEerr := json.Unmarshal(jsonBytes, &unstructured)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return unstructured, nil
}
