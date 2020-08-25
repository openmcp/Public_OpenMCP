package main

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"openmcp/openmcp/testRestClient/rest"
	"openmcp/openmcp/testRestClient/v1alpha1"
)


func main() {
	host_config, _ := rest.InClusterConfig()
	crdClient, err := v1alpha1.NewForConfig(host_config)
	if err != nil {
		fmt.Println(err)
	}
	openmcpPolicyInstance, err  := crdClient.OpenMCPPolicy("openmcp").Get("metric-collector-period", metav1.GetOptions{})

	fmt.Println(openmcpPolicyInstance.Name)
}
