package main

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	ketiapis "resource-controller/apis"
	ketiv1alpha1 "resource-controller/apis/keti/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type A struct{
	live client.Client
}
func main() {
	host_config, _ := rest.InClusterConfig()
	live := cluster.New("openmcp", host_config, cluster.Options{CacheOptions: cluster.CacheOptions{}})
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		fmt.Printf("getting delegating client for live cluster: %v", err)
	}

	if err := ketiapis.AddToScheme(live.GetScheme()); err != nil {
		fmt.Println("adding APIs to live cluster's scheme: %v", err)
	}

	r := &A{
		live: liveclient,
	}

	openmcphasInstance := &ketiv1alpha1.OpenMCPHybridAutoScaler{}

	err = r.live.Get(context.TODO(), types.NamespacedName{Namespace:"openmcp", Name:"openmcp-hpa"},openmcphasInstance)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(openmcphasInstance.Name)
	}


	openmcpDeployInstance := &appsv1.Deployment{}
	err = r.live.Get(context.TODO(), types.NamespacedName{Namespace:"openmcp", Name:"sync-controller"},openmcpDeployInstance)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(openmcpDeployInstance.Name+"!")
	}



}
