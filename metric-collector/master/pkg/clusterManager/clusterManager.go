package clusterManager

import (
	"context"
	"fmt"
	"k8s.io/client-go/rest"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"

	clientV1alpha1 "resource-controller/clientset/v1alpha1"
)

type ClusterManager struct {
	HostConfig  *rest.Config
	HostClient  genericclient.Client
	ClusterList *fedv1b1.KubeFedClusterList

	Crd_client			*clientV1alpha1.ExampleV1Alpha1Client
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

	crd_client, _ := clientV1alpha1.NewForConfig(config)

	cm := &ClusterManager{
		HostConfig:  config,
		HostClient:  hostClient,
		ClusterList: clusterList,
		Crd_client:	 crd_client,
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
