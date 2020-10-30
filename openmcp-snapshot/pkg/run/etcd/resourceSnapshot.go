package etcd

import (
	"fmt"

	snapshotv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
)

//func GetResource(resourceType string, resourceName string) {
//}

// CreateSnapshot : 스냅샷 생성
func CreateSnapshot(snapshotSource snapshotv1alpha1.SnapshotSource) (bool, error) {
	clientset := InitKube()

	// 실패 케이스1
	resourceType := snapshotSource.ResourceType
	resourceName := snapshotSource.ResourceName
	resourceNamespace := snapshotSource.ResourceNamespace

	val, err := GetResourceJSON(clientset, resourceType, resourceName, resourceNamespace)

	if err != nil {
		return false, err
	}

	isSuccess, etcdErr := InsertEtcd(resourceType+"."+resourceName, val)
	if etcdErr != nil {
		return false, etcdErr
	}
	if !isSuccess {
		return false, fmt.Errorf("Insert Etcd Fail")
	}

	return true, nil
}

// RestoreSnapshot : 스냅샷 복구
func RestoreSnapshot(snapshotKey string, resourceType string) (bool, error) {

	// ETCD에서 json 가져오기
	jsonStr, etcdErr := GetEtcd(snapshotKey)
	if etcdErr != nil {
		return false, etcdErr
	}

	clientset := InitKube()
	// defer 는 fatal 전에
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Create Resource ERROR", r)
		}
	}()
	isSuccess, err := CreateResourceJSON(clientset, resourceType, jsonStr)
	if err != nil {

		return false, err // 에러 발생
	} else if !isSuccess {
		return false, fmt.Errorf("error") // 에러 발생
	}

	return true, nil
}
