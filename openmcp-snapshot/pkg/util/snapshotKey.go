package util

/**
startTime, snapshotSource.ResourceCluster, snapshotSource.ResourceType, snapshotSource.ResourceName 를 '-' 로 조합하여 만든다.
volume 스냅샷의 경우 resourceType 이 volume 이다.
*/

import (
	nanumv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"strings"
)

const TypeVolumeSnapshot = "volume"

func MakeVolumeSnapshotKey(startTime string, clusterName string, resourceName string) string {
	ret := MakeSnapshotKey(startTime, clusterName, TypeVolumeSnapshot, resourceName)
	return ret
}

func MakeSnapshotKeyForSnapshot(startTime string, snapshotSource *nanumv1alpha1.SnapshotSource) string {
	ret := MakeSnapshotKey(startTime, snapshotSource.ResourceCluster, snapshotSource.ResourceType, snapshotSource.ResourceName)
	return ret
}

func MakeSnapshotKeyForSnapshotRestore(startTime string, snapshotSourceRestore *nanumv1alpha1.SnapshotSource) string {
	ret := MakeSnapshotKey(startTime, snapshotSourceRestore.ResourceCluster, snapshotSourceRestore.ResourceType, snapshotSourceRestore.ResourceName)
	return ret
}

func MakeSnapshotKey(startTime string, clusterName string, SnapshotType string, resourceName string) string {
	ret := strings.Join([]string{startTime, clusterName, SnapshotType, resourceName}, "-")
	return ret
}

func GetStartTimeBySnapshotKey(snapshotKey string) string {
	tmp := strings.Split(snapshotKey, "-")
	startTime := tmp[0]
	return startTime
}

func GetResourceNameBySnapshotKey(snapshotKey string) string {
	tmp := strings.Split(snapshotKey, "-")
	resourceName := tmp[len(tmp)-1]
	return resourceName
}
