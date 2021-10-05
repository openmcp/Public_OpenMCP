package predicates

import (
	"openmcp/openmcp/omcplog"
	ketiresource "openmcp/openmcp/openmcp-scheduler/src/resourceinfo"
	"openmcp/openmcp/util/clusterManager"
)

/*
 this filter checks status of cluster that it being join or joining
*/
type ClusterJoninCheck struct{}

func (pl *ClusterJoninCheck) Name() string {
	return "ClusterJoninCheck"
}

func (pl *ClusterJoninCheck) Filter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, cm *clusterManager.ClusterManager) bool {
	clusterList := clusterInfo.ClusterList
	if clusterList == nil {
		omcplog.V(0).Infof("That instance did not get information from crd cluster.")
	}
	// joinCluster := make(map[string]bool)
	for _, cluster := range clusterList.Items {
		if cluster.Name == "" {
			continue
		}
		if "JOIN" == cluster.Spec.JoinStatus {
			if clusterInfo.ClusterName == cluster.Name {
				return true
			}
		}
	}
	return false

}
