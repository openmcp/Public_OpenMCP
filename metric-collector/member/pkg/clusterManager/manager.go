package clusterManager

import (
	"cluster-metric-collector/pkg/kubeletClient"

	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ClusterManager struct {
	Host_config    *rest.Config
	Host_client    *kubernetes.Clientset
	Node_list      *corev1.NodeList
	Kubelet_client *kubeletClient.KubeletClient
}

func NewClusterManager() *ClusterManager {
	config, err := rest.InClusterConfig()
	if err != nil {

	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {

	}
	node_list, err := GetNodeList(clientSet)
	if err != nil {

	}
	kubeletClient, err := kubeletClient.NewKubeletClient()
	if err != nil {

	}
	cm := &ClusterManager{
		Host_config:    config,
		Host_client:    clientSet,
		Node_list:      node_list,
		Kubelet_client: kubeletClient,
	}
	return cm
}

func GetNodeList(clientSet *kubernetes.Clientset) (*corev1.NodeList, error) {

	nodeList, err := clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		fmt.Println("Error retrieving list of Node: %+v", err)
		return nodeList, err
	}
	if len(nodeList.Items) == 0 {
		fmt.Println("No Nodes found")
		return nodeList, err
	}
	return nodeList, nil
}
