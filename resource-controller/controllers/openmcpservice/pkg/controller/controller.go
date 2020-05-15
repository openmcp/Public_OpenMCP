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

package openmcpservice // import "admiralty.io/multicluster-controller/examples/openmcpservice/pkg/controller/openmcpservice"

import (
	"admiralty.io/multicluster-controller/pkg/reference"
	"context"
	"encoding/json"
	"fmt"
	"github.com/getlantern/deepcopy"
	//"github.com/jinzhu/copier"
	"k8s.io/apimachinery/pkg/api/errors"
	//"math/rand"
	"reflect"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	//"reflect"
	//"sort"
	//"time"
	//"github.com/jinzhu/copier"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"resource-controller/apis"
	ketiv1alpha1 "resource-controller/apis/keti/v1alpha1"
	//corev1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
)

type ClusterManager struct {
	Fed_namespace   string
	Host_config     *rest.Config
	Host_client     genericclient.Client
	Cluster_list    *fedv1b1.KubeFedClusterList
	Cluster_configs map[string]*rest.Config
	Cluster_clients map[string]genericclient.Client
}

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
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPService{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	for _, ghost := range ghosts {
		fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(ghost, &corev1.Service{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}

	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

var i int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	i += 1
	fmt.Println("********* [", i, "] *********")
	fmt.Println(req.Context, " / ", req.Namespace, " / ", req.Name)
	cm := NewClusterManager()

	// Fetch the OpenMCPService instance
	instance := &ketiv1alpha1.OpenMCPService{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	fmt.Println("instance Name: ", instance.Name)
	fmt.Println("instance Namespace : ", instance.Namespace)

	if err != nil {
		if errors.IsNotFound(err) {
			// ...TODO: multicluster garbage collector
			// Until then...
			fmt.Println("Delete Services ..Cluster")
			err := cm.DeleteServices(req.NamespacedName)
			return reconcile.Result{}, err
		}
		fmt.Println("Error1")
		return reconcile.Result{}, err
	}
	if instance.Status.ClusterMaps == nil {
		fmt.Println("Service Create Start")
		err := r.createService(req, cm, instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil

	}
	if !reflect.DeepEqual(instance.Status.LastSpec, instance.Spec) {
		fmt.Println("Service Update Start")
		err := r.updateService(req, cm, instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil

	}
	if instance.Status.ChangeNeed {
		fmt.Println("Receive notify from OpenMCP Deployment ")

		instance.Status.ChangeNeed = false
		r.updateService(req, cm, instance)
	}

	//odeploy := &ketiv1alpha1.OpenMCPDeployment{}
	//cm.Host_client.Get(context.TODO(), odeploy, req.Namespace, req.Name)
	//
	//for _, cluster := range cm.Cluster_list.Items {
	//	if odeploy.Status.ClusterMaps[cluster.Name] == 0 && instance.Status.ClusterMaps[cluster.Name] == 1 {
	//
	//	} else if odeploy.Status.ClusterMaps[cluster.Name] != 0 && instance.Status.ClusterMaps[cluster.Name] == 0 {
	//
	//	}
	//}
	//// Check Service in cluster
	//for k, _ := range instance.Status.ClusterMaps {
	//	cluster_name := k
	//
	//	found := &corev1.Service{}
	//	cluster_client := cm.Cluster_clients[cluster_name]
	//	err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-service")
	//	if err != nil && errors.IsNotFound(err) {
	//		// Delete Serivce Detected
	//		fmt.Println("Cluster '" + cluster_name  + "' ReDeployed")
	//		svc := r.serviceForOpenMCPService(req, instance)
	//		err = cluster_client.Create(context.Background(), svc)
	//
	//		if err != nil {
	//			return reconcile.Result{}, err
	//		}
	//
	//	}
	//
	//}

	return reconcile.Result{}, nil // err
}

func (r *reconciler) serviceForOpenMCPService(req reconcile.Request, m *ketiv1alpha1.OpenMCPService) *corev1.Service {

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-service",
			Namespace: m.Namespace,
		},
		//Spec: m.Spec.Template.Spec,
	}
	//svc.Spec.Selector = m.Spec.OpenMCPLabelSelector

	deepcopy.Copy(&svc.Spec, &m.Spec.Template.Spec)
	deepcopy.Copy(&svc.Spec.Selector, &m.Spec.LabelSelector)

	reference.SetMulticlusterControllerReference(svc, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return svc
}
func (cm *ClusterManager) DeleteServices(nsn types.NamespacedName) error {
	svc := &corev1.Service{}
	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_clients[cluster.Name]
		fmt.Println(nsn.Namespace, nsn.Name)
		err := cluster_client.Get(context.Background(), svc, nsn.Namespace, nsn.Name+"-service")
		if err != nil && errors.IsNotFound(err) {
			// all good
			fmt.Println("Not Found")
			continue
		}
		if !isInObject(svc, "OpenMCPService") {
			continue
		}
		fmt.Println(cluster.Name, " Delete Start")
		err = cluster_client.Delete(context.Background(), svc, nsn.Namespace, nsn.Name+"-service")
		if err != nil {
			return err
		}
		fmt.Println(cluster.Name, "Delete Complate")
	}
	return nil

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

func isInObject(child *corev1.Service, parent string) bool {
	refKind_str := child.ObjectMeta.Annotations["multicluster.admiralty.io/controller-reference"]
	refKind_map := make(map[string]interface{})
	err := json.Unmarshal([]byte(refKind_str), &refKind_map)
	if err != nil {
		panic(err)
	}
	if refKind_map["kind"] == parent {
		return true
	}
	return false
}
func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (r *reconciler) getClusterIncludeLabel(label_map map[string]string, namespace string) []string {
	result_cluster_list := []string{}

	odeploy_list := &ketiv1alpha1.OpenMCPDeploymentList{}
	listOptions := &client.ListOptions{Namespace: namespace}

	r.live.List(context.TODO(), listOptions, odeploy_list)
	for _, odeploy := range odeploy_list.Items {
		fmt.Println("odeploy Name : ", odeploy.Name)
		for k, v := range odeploy.Spec.Labels {
			fmt.Println("label : ", k, " / ", v)
			if label_map[k] == v {
				fmt.Println("Match!")
				for cluster_name, replica := range odeploy.Status.ClusterMaps {
					fmt.Println(cluster_name, replica)
					if replica != 0 {
						result_cluster_list = append(result_cluster_list, cluster_name)
					}
				}
			}
		}
	}
	result_cluster_list = unique(result_cluster_list)

	fmt.Println(result_cluster_list)
	return result_cluster_list
}
func (r *reconciler) createService(req reconcile.Request, cm *ClusterManager, instance *ketiv1alpha1.OpenMCPService) error {
	cluster_map := make(map[string]int32)

	for _, cluster := range cm.Cluster_list.Items {
		cluster_map[cluster.Name] = 0
	}

	label_include_cluster_list := r.getClusterIncludeLabel(instance.Spec.LabelSelector, instance.Namespace)
	svc := r.serviceForOpenMCPService(req, instance)

	for _, cluster_name := range label_include_cluster_list {
		found := &corev1.Service{}
		cluster_client := cm.Cluster_clients[cluster_name]

		err := cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-service")

		if err != nil && errors.IsNotFound(err) {
			err = cluster_client.Create(context.TODO(), svc)
			cluster_map[cluster_name] = 1
			if err != nil {
				return err
			}
		}
	}

	instance.Status.LastSpec = instance.Spec
	instance.Status.ClusterMaps = cluster_map

	err := r.live.Status().Update(context.TODO(), instance)
	return err

}
func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
func (r *reconciler) updateService(req reconcile.Request, cm *ClusterManager, instance *ketiv1alpha1.OpenMCPService) error {
	cluster_map := make(map[string]int32)

	for _, cluster := range cm.Cluster_list.Items {
		cluster_map[cluster.Name] = 0
	}
	label_include_cluster_list := r.getClusterIncludeLabel(instance.Spec.LabelSelector, instance.Namespace)

	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_clients[cluster.Name]

		found := &corev1.Service{}
		err := cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-service")

		if contains(label_include_cluster_list, cluster.Name) {
			if err != nil && errors.IsNotFound(err) {
				svc := r.serviceForOpenMCPService(req, instance)
				cluster_map[cluster.Name] = 1

				err = cluster_client.Create(context.TODO(), svc)
				if err != nil {
					return err
				}
			} else if err == nil {
				svc := r.serviceForOpenMCPService(req, instance)

				svc.Spec.ClusterIP = found.Spec.ClusterIP
				svc.ResourceVersion = found.ResourceVersion

				cluster_map[cluster.Name] = 1

				err = cluster_client.Update(context.TODO(), svc)
				if err != nil {
					return err
				}
			}
		} else {
			if err != nil && errors.IsNotFound(err) {
				continue
			} else if err == nil {
				err = cluster_client.Delete(context.TODO(), found, instance.Namespace, instance.Name+"-service")
				if err != nil {
					return err
				}
			}
		}
	}
	instance.Status.LastSpec = instance.Spec
	err := r.live.Status().Update(context.TODO(), instance)
	return err

}
