package clusterManager

import (
	"context"
	"fmt"
	clientV1alpha1 "openmcp/openmcp/clientset/v1alpha1"
	"openmcp/openmcp/omcplog"

	"k8s.io/client-go/util/flowcontrol"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"sigs.k8s.io/kubefed/pkg/controller/util"
)

type ClusterManager struct {
	Fed_namespace       string
	Host_config         *rest.Config
	Host_client         genericclient.Client
	Host_kubeClient     *kubernetes.Clientset
	Crd_client          *clientV1alpha1.ExampleV1Alpha1Client
	Crd_cluster_client  *clientV1alpha1.ExampleV1Alpha1Client
	Crd_istio_client    *clientV1alpha1.ExampleV1Alpha1Client
	Cluster_list        *fedv1b1.KubeFedClusterList
	Node_list           *corev1.NodeList
	Cluster_configs     map[string]*rest.Config
	Cluster_genClients  map[string]genericclient.Client
	Cluster_kubeClients map[string]*kubernetes.Clientset
	Cluster_dynClients  map[string]dynamic.Interface
	//Mutex	*sync.Mutex
}

func ListKubeFedClusters(genClient genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	tempClusterList := &fedv1b1.KubeFedClusterList{}
	clusterList := &fedv1b1.KubeFedClusterList{}

	err := genClient.List(context.TODO(), tempClusterList, namespace, &client.ListOptions{})

	if err != nil {
		fmt.Println("Error retrieving list of federated clusters: %+v", err)
	}
	if len(tempClusterList.Items) == 0 {
		fmt.Println("No federated clusters found")
	}

	// Status Check
	for _, cluster := range tempClusterList.Items {
		status := true

		for _, cond := range cluster.Status.Conditions {
			if cond.Type == "Offline" {
				status = false
				break
			}
		}
		if status {
			clusterList.Items = append(clusterList.Items, cluster)
		}
	}

	return clusterList
}

func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, genClient genericclient.Client, fedNamespace string) map[string]*rest.Config {
	clusterConfigs := make(map[string]*rest.Config)
	for _, cluster := range clusterList.Items {
		config, _ := util.BuildClusterConfig(&cluster, genClient, fedNamespace)
		clusterConfigs[cluster.Name] = config
	}
	return clusterConfigs
}
func KubeFedClusterGenClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {

	cluster_clients := make(map[string]genericclient.Client)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		cluster_config := cluster_configs[clusterName]
		//cluster_client := genericclient.NewForConfigOrDie(cluster_config)
		cluster_config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(100, 100) // default 5 10
		cluster_client, _ := genericclient.New(cluster_config)
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
	omcplog.V(4).Info("NewClusterManager Called")
	//mutex := &sync.Mutex{}
	fed_namespace := "kube-federation-system"

	//컨트롤러 모듈 생성 시 정의한 sa, secret 기반 config
	host_config, _ := rest.InClusterConfig()

	//컨트롤러 배포시 생성된 secret 기반 kube client, crd client 생성
	host_kubeClient := kubernetes.NewForConfigOrDie(host_config)             // k8sscheme 기본 kube 리소스
	crd_client, _ := clientV1alpha1.NewForConfig(host_config)                //resourcev1alpha1
	crd_cluster_client, _ := clientV1alpha1.NewClusterForConfig(host_config) //clusterv1alpha1
	crd_istio_client, _ := clientV1alpha1.NewIstioForConfig(host_config)     // v1alpha3 networking.istio.io

	//federation list 받아오기 위한 컨트롤러 SA secret 기반 generic client 생성
	host_client := genericclient.NewForConfigOrDie(host_config) //fedapis, k8sscheme
	cluster_list := ListKubeFedClusters(host_client, fed_namespace)

	//federation 조인시 member cluster에 생성한 sa, secret 기반 config
	cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)

	//kubefed secret 기반 generic, kube, dyn client 생성
	cluster_gen_clients := KubeFedClusterGenClients(cluster_list, cluster_configs)   //fedapis, k8sscheme //throttling log
	cluster_kube_clients := KubeFedClusterKubeClients(cluster_list, cluster_configs) // k8sscheme
	cluster_dyn_clients := KubeFedClusterDynClients(cluster_list, cluster_configs)

	node_list, _ := GetNodeList(host_kubeClient)

	cm := &ClusterManager{
		Fed_namespace:       fed_namespace,
		Host_config:         host_config,
		Host_client:         host_client,
		Host_kubeClient:     host_kubeClient,
		Crd_client:          crd_client,
		Crd_cluster_client:  crd_cluster_client,
		Crd_istio_client:    crd_istio_client,
		Cluster_list:        cluster_list,
		Node_list:           node_list,
		Cluster_configs:     cluster_configs,
		Cluster_genClients:  cluster_gen_clients,
		Cluster_kubeClients: cluster_kube_clients,
		Cluster_dynClients:  cluster_dyn_clients,
		//Mutex:	mutex,
	}
	return cm
}
func GetNodeList(clientSet *kubernetes.Clientset) (*corev1.NodeList, error) {

	nodeList := &corev1.NodeList{}
	nodeList, err := clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nodeList, err
	}
	if len(nodeList.Items) == 0 {
		fmt.Println("No Nodes found")
		return nodeList, err
	}
	return nodeList, nil
}
