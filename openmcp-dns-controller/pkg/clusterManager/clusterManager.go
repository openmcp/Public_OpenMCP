package clusterManager

import (
	"context"
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"sigs.k8s.io/kubefed/pkg/controller/util"
)

type ClusterManager struct {
	Fed_namespace       string
	Host_config         *rest.Config
	Host_client         genericclient.Client
	Cluster_list        *fedv1b1.KubeFedClusterList
	Cluster_configs     map[string]*rest.Config
	Cluster_genClients  map[string]genericclient.Client
	Cluster_kubeClients map[string]*kubernetes.Clientset
	Cluster_dynClients  map[string]dynamic.Interface
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

func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, client genericclient.Client, fedNamespace string) map[string]*rest.Config {
	clusterConfigs := make(map[string]*rest.Config)
	for _, cluster := range clusterList.Items {
		config, _ := util.BuildClusterConfig(&cluster, client, fedNamespace)
		clusterConfigs[cluster.Name] = config
	}
	return clusterConfigs
}
func KubeFedClusterGenClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {

	cluster_clients := make(map[string]genericclient.Client)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		cluster_config := cluster_configs[clusterName]
		cluster_client := genericclient.NewForConfigOrDie(cluster_config)
		cluster_clients[clusterName] = cluster_client
	}
	return cluster_clients
}
func KubeFedClusterKubeClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]*kubernetes.Clientset {

	cluster_clients := make(map[string]*kubernetes.Clientset)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		cluster_config := cluster_configs[clusterName]
		cluster_client := kubernetes.NewForConfigOrDie(cluster_config)
		cluster_clients[clusterName] = cluster_client
	}
	return cluster_clients
}
func KubeFedClusterDynClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]dynamic.Interface {

	cluster_clients := make(map[string]dynamic.Interface)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		cluster_config := cluster_configs[clusterName]
		cluster_client := dynamic.NewForConfigOrDie(cluster_config)
		cluster_clients[clusterName] = cluster_client
	}
	return cluster_clients
}
func NewClusterManager() *ClusterManager {
	fed_namespace := "kube-federation-system"

	host_config, _ := rest.InClusterConfig()
	host_client := genericclient.NewForConfigOrDie(host_config)
	cluster_list := ListKubeFedClusters(host_client, fed_namespace)
	cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)
	cluster_gen_clients := KubeFedClusterGenClients(cluster_list, cluster_configs)
	cluster_kube_clients := KubeFedClusterKubeClients(cluster_list, cluster_configs)
	cluster_dyn_clients := KubeFedClusterDynClients(cluster_list, cluster_configs)

	cm := &ClusterManager{
		Fed_namespace:       fed_namespace,
		Host_config:         host_config,
		Host_client:         host_client,
		Cluster_list:        cluster_list,
		Cluster_configs:     cluster_configs,
		Cluster_genClients:  cluster_gen_clients,
		Cluster_kubeClients: cluster_kube_clients,
		Cluster_dynClients:  cluster_dyn_clients,
	}
	return cm
}
