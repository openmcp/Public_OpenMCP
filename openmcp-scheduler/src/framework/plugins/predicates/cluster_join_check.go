/*
log 1레벨 결과
2레벨 필터 & 스코어 결과 추가
3레벨 필터 & 스코어 연산 과정 추가
4레벨 연산과정에 대한 모든 로깅
5레벨 디버깅관련 모든 로깅
*/
package predicates

import (
	"openmcp/openmcp/omcplog"
	ketiresource "openmcp/openmcp/openmcp-scheduler/src/resourceinfo"
	"openmcp/openmcp/util/clusterManager"
	"time"
)

/*
 this filter checks status of cluster that it being join or joining
*/
type ClusterJoninCheck struct{}

func (pl *ClusterJoninCheck) Name() string {
	return "ClusterJoninCheck"
}

func (pl *ClusterJoninCheck) Filter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, cm *clusterManager.ClusterManager) bool {
	startTime := time.Now()
	clusterList := clusterInfo.ClusterList
	if clusterList == nil {
		omcplog.V(5).Infof("That instance did not get information from crd cluster.")
	}
	// joinCluster := make(map[string]bool)
	for _, cluster := range clusterList.Items {
		if cluster.Name == "" {
			continue
		}
		if "JOIN" == cluster.Spec.JoinStatus {
			if clusterInfo.ClusterName == cluster.Name {
				//omcplog.V(3).Info("ClusterJoninCheck true ")
				elapsedTime := time.Since(startTime)
				omcplog.V(3).Infof("ClusterJoninCheck Time [%v]", elapsedTime)
				return true
			}
		}
	}
	omcplog.V(4).Info("ClusterJoninCheck false ")
	elapsedTime := time.Since(startTime)
	omcplog.V(3).Infof("ClusterJoninCheck Time [%v]", elapsedTime)
	return false

}
