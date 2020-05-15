/*
Copyright 2018 The Multicluster-Controller Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller // import "admiralty.io/multicluster-controller/examples/openmcppolicyengine/pkg/controller/openmcppolicyengine"

import (
	"context"
	"fmt"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	//	"admiralty.io/multicluster-controller/pkg/reference"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"resource-controller/apis"
	ketiv1alpha1 "resource-controller/apis/keti/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"

	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
)

var log = logf.Log.WithName("controller_openmcphybridautoscaler")

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string) (*controller.Controller, error) {
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients = append(ghostclients, ghostclient)
	}

	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPPolicyEngine{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

var i int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	i += 1
	fmt.Println("********* [", i, "] *********")
	fmt.Println("Request Context: ", req.Context, " / Request Namespace: ", req.Namespace, " /  Request Name: ", req.Name)
	//cm := NewClusterManager()

	// Fetch the OpenMCPDeployment instance
	instance := &ketiv1alpha1.OpenMCPPolicyEngine{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	//	fmt.Println("instance: ", instance)
	fmt.Println("instance Name: ", instance.Name)
	fmt.Println("instance Namespace: ", instance.Namespace)

	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("Delete Policy Resource")
			return reconcile.Result{}, nil
		}
		fmt.Println("Error1")
		reqLogger.Error(err, "Failed to get hasInstance")
		return reconcile.Result{}, err
	}

	if instance.Spec.PolicyStatus == "Disabled" {
		fmt.Println("Policy Disabled")
	} else if instance.Spec.PolicyStatus == "Enabled" {
		if instance.Spec.RangeOfApplication == "FromNow" {
			fmt.Println("Policy Enabled - FromNow")
		} else if instance.Spec.RangeOfApplication == "All" {
			//fmt.Println("Policy Enabled - All")
			object := instance.Spec.Template.Spec.TargetController.Kind
			if object == "OpenMCPHybridAutoScaler" {
				fmt.Println("Policy Enabled - OpenMCPHybridAutoScaler")
				hpaList := &ketiv1alpha1.OpenMCPHybridAutoScalerList{}
				listOptions := &client.ListOptions{Namespace: ""} //all resources
				r.live.List(context.TODO(), listOptions, hpaList)
				//fmt.Println("List: ", hpaList)
				for _, hpaInstance := range hpaList.Items {
					//fmt.Println("hpastatus: ",hpaInstance.Status.Policies)
					//fmt.Println("policies: ",instance.Spec.Template.Spec.Policies)
					var i = 0
					for index, tmpPolicy := range hpaInstance.Status.Policies { //정책 이름 대조하여 해당 정책만 수정
						if tmpPolicy.Type == instance.Spec.Template.Spec.Policies[0].Type { //같은 정책이 이미 있는 경우
							i++
							hpaInstance.Status.Policies[index].Value = instance.Spec.Template.Spec.Policies[0].Value
							break
						}
					}
					if i == 0 {
						hpaInstance.Status.Policies = append(hpaInstance.Status.Policies, instance.Spec.Template.Spec.Policies...)
					}
					err := r.live.Status().Update(context.TODO(), &hpaInstance)
					if err != nil {
						fmt.Println("OpenMCPHPA Policy Update Error")
						return reconcile.Result{}, err
					} else {
						fmt.Println("OpenMCPHPA Policy UPDATE Success!")
					}
				}
			} else if object == "OpenMCPLoadbalancer" {

			}
		}
	}
	return reconcile.Result{}, nil
}

/*func (cm *ClusterManager) DeleteOpenMCPPolicyEngine(nsn types.NamespacedName) error {
	dep := &appsv1.Deployment{}
	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_clients[cluster.Name]
		fmt.Println(nsn.Namespace, nsn.Name)
		err := cluster_client.Get(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")
		if err != nil && errors.IsNotFound(err) {
			// all good
			fmt.Println("Not Found")
			continue
		}
		fmt.Println(cluster.Name," Delete Start")
		err = cluster_client.Delete(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")
		if err != nil {
			return err
		}
		fmt.Println(cluster.Name, "Delete Complate")
	}
	return nil

}*/

type ClusterManager struct {
	Fed_namespace   string
	Host_config     *rest.Config
	Host_client     genericclient.Client
	Cluster_list    *fedv1b1.KubeFedClusterList
	Cluster_configs map[string]*rest.Config
	Cluster_clients map[string]genericclient.Client
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
func KubeFedClusterClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {

	cluster_clients := make(map[string]genericclient.Client)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		cluster_config := cluster_configs[clusterName]
		cluster_client := genericclient.NewForConfigOrDie(cluster_config)
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
	cluster_clients := KubeFedClusterClients(cluster_list, cluster_configs)

	cm := &ClusterManager{
		Fed_namespace:   fed_namespace,
		Host_config:     host_config,
		Host_client:     host_client,
		Cluster_list:    cluster_list,
		Cluster_configs: cluster_configs,
		Cluster_clients: cluster_clients,
	}
	return cm
}
