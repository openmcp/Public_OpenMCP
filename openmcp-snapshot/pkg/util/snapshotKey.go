package util

/**
startTime, snapshotSource.ResourceCluster, snapshotSource.ResourceType, snapshotSource.ResourceName 를 '-' 로 조합하여 만든다.
volume 스냅샷의 경우 resourceType 이 volume 이다.
*/

import (
	"fmt"
	nanumv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"strconv"
	"strings"
)

const TypeVolumeSnapshot = "volume"

//func MakeVolumeSnapshotKey(startTime string, clusterName string, resourceName string) string {
//	ret := MakeSnapshotKey(startTime, clusterName, TypeVolumeSnapshot, resourceName)
//	return ret
//}

// ex) 1624008391-volume-1  -> sns-1624008391-volume-1-job, ...
func MakeVolumeProcessResourceKey(startTime string, middleName string, idx string) string {
	head := startTime
	ret := strings.Join([]string{head, middleName, idx}, "-")
	return ret
}

// ex) 1624008391/cluster1/Deployment/iot-gateway
func MakeSnapshotKeyForSnapshot(greoupSnapshotKey string, snapshotSource *nanumv1alpha1.SnapshotSource) string {
	ret := MakeSnapshotKey(greoupSnapshotKey, snapshotSource.ResourceCluster, snapshotSource.ResourceType, snapshotSource.ResourceName)
	return ret
}

// ex) 1624008391/cluster1/Deployment/iot-gateway
func MakeSnapshotKeyForSnapshotRestore(greoupSnapshotKey string, snapshotSourceRestore *nanumv1alpha1.SnapshotSource) string {
	ret := MakeSnapshotKey(greoupSnapshotKey, snapshotSourceRestore.ResourceCluster, snapshotSourceRestore.ResourceType, snapshotSourceRestore.ResourceName)
	return ret
}

//
func MakeSnapshotKey(greoupSnapshotKey string, resourceCluster string, resourceType string, resourceName string) string {
	ret := strings.Join([]string{MakePrifix(greoupSnapshotKey), resourceCluster, resourceType, resourceName}, "/")
	return ret
}

func MakePrifix(greoupSnapshotKey string) string {
	head := ""
	if ETCDROOT != "" {
		head = strings.Join([]string{ETCDROOT, greoupSnapshotKey}, "/")
	} else {
		head = strings.Join([]string{greoupSnapshotKey}, "/")
	}
	return head
}

func GetGroupSnapshotKeyBySnapshotKey(snapshotKey string) (string, error) {
	return GetCommonNameBySnapshotKey(snapshotKey, 0)
}

func GetClusterNameBySnapshotKey(snapshotKey string) (string, error) {
	return GetCommonNameBySnapshotKey(snapshotKey, 1)
}

func GetResourceTypeBySnapshotKey(snapshotKey string) (string, error) {
	return GetCommonNameBySnapshotKey(snapshotKey, 2)
}

func GetResourceNameBySnapshotKey(snapshotKey string) (string, error) {
	return GetCommonNameBySnapshotKey(snapshotKey, 3)
}

func GetCommonNameBySnapshotKey(snapshotKey string, idx int) (string, error) {
	retVal := ""
	tmp := strings.Split(snapshotKey, "/")
	if len(tmp) < 1 {
		return "", fmt.Errorf("GetCommonNameBySnapshotKey : snapshotKey is not valid [" + snapshotKey + ", " + strconv.Itoa(idx) + "]")
	}
	if ETCDROOT != "" {
		retVal = tmp[idx+len(strings.Split(ETCDROOT, "/"))]
	} else {
		retVal = tmp[idx]
	}
	if retVal == "" {
		return "", fmt.Errorf("GetCommonNameBySnapshotKey : groupSnapshotKey is not valid [" + snapshotKey + ", " + strconv.Itoa(idx) + "]")
	}
	return retVal, nil
}
