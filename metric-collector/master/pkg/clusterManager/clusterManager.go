package clusterManager

import (
	"context"
	"fmt"
	"k8s.io/client-go/rest"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
)

type ClusterManager struct {
	HostConfig  *rest.Config
	HostClient  genericclient.Client
	ClusterList *fedv1b1.KubeFedClusterList
}

func NewClusterManager() *ClusterManager {
	fedNamespace := "kube-federation-system"
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	hostClient, err := genericclient.New(config)
	if err != nil {
		fmt.Println(err)
	}
	clusterList := ListKubeFedClusters(hostClient, fedNamespace)

	cm := &ClusterManager{
		HostConfig:  config,
		HostClient:  hostClient,
		ClusterList: clusterList,
	}
	return cm
}
func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		fmt.Println("Error retrieving list of federated clusters: %+v", err)
	}
	if len(clusterList.Items) == 0 {
		fmt.Println("No federated clusters found")
	}
	return clusterList
}
