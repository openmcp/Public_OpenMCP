package syncResource

import (
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SyncResource(cm *clusterManager.ClusterManager) map[string]int32 {
	resultClusterMaps := make(map[string]int32)

	oclusterList, err := cm.Crd_client.OpenMCPCluster("openmcp").List(metav1.ListOptions{})
	if err != nil {
		omcplog.V(0).Info(err)
	}

	for _, ocluster := range oclusterList.Items {
		if ocluster.Spec.JoinStatus == "JOIN" || ocluster.Spec.JoinStatus == "JOINING" {
			resultClusterMaps[ocluster.Name] = 1
		} else {
			resultClusterMaps[ocluster.Name] = 0
		}

	}

	return resultClusterMaps
}
