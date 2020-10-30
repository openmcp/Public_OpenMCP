package util

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//GetClientset kubernetes의 clientset 생성
func GetClientset(clusterInfo string) (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	con, err := clientcmd.NewClientConfigFromBytes([]byte(clusterInfo))
	if err != nil {
		return nil, err
	}
	clientconf, err := con.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err = kubernetes.NewForConfig(clientconf)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
