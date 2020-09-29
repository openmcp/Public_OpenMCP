package migration

import (
	"fmt"
	"openmcp/openmcp/util/clusterManager"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
)

func TestReconcile(t *testing.T) {
	fmt.Println("test reconcile")

	cm = clusterManager.NewClusterManager()
	fmt.Println("test reconcile111")
	cluster_client := cm.Cluster_genClients["cluster1"]
	newCm := cm.Cluster_list
	fmt.Println(newCm)
	fmt.Println(cm.Cluster_list.Items[0].ClusterName)

	for _, myCluster := range cm.Cluster_list.Items {
		found := &appsv1.Deployment{}
		cluster_client = cm.Cluster_genClients[myCluster.Name]
		fmt.Println(cluster_client, found)
	}
}
