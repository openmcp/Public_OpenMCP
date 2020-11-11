package agent

import (
	"errors"

	"openmcp/openmcp/openmcp-globalcache/pkg/utils"
	"openmcp/openmcp/util/clusterManager"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//RegistryNodeManager dd
type RegistryNodeManager struct {
	clientset   *kubernetes.Clientset
	clusterName string
	nodeName    string
	address     string
	port        string
}

var cm *clusterManager.ClusterManager

func (r *RegistryNodeManager) init(clusterName string, nodeName string) error {
	//get Clientset
	var clientset *kubernetes.Clientset
	clientset = cm.Cluster_kubeClients[clusterName]

	//get Address
	result, getErr := clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if getErr != nil {
		return getErr
	}
	var ipAddress string
	//var hostnameAddress string
	for _, address := range result.Status.Addresses {
		/*if address.Type == "Hostname" {
			hostnameAddress = address.Address
		}*/
		if address.Type == "InternalIP" {
			ipAddress = address.Address
		}
	}

	if ipAddress == "" {
		return errors.New("RegistryNodeManager not set ipAddress")
	}

	r.clientset = clientset
	r.clusterName = clusterName
	r.nodeName = nodeName
	r.address = ipAddress
	r.port = string(utils.AgentRegistryManagerDefaultPort)

	return nil
}
