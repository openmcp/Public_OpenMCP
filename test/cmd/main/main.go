package main

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	clientV1alpha1 "clientset/v1alpha1"
)

var kubeconfig string
func main() {
	config, err := rest.InClusterConfig()

	clientSet, err := clientV1alpha1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	openmcpdeployments, err := clientSet.OpenMCPDeployment("openmcp").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("openmcpdeployments found: %+v\n", openmcpdeployments)
}
