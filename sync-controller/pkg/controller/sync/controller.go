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

package sync // import "admiralty.io/multicluster-controller/examples/serviceDNS/pkg/controller/serviceDNS"

import (
	"context"
	"encoding/json"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"

	//"encoding/json"
	"fmt"
	"reflect"

	//"github.com/jinzhu/copier"
	//corev1 "k8s.io/api/core/v1"

	//corev1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/runtime"

	//"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"

	//corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	//appsv1 "k8s.io/api/apps/v1"2
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	//"reflect"

	"sigs.k8s.io/kubefed/pkg/controller/util"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sync-controller/pkg/apis"
	ketiv1alpha1 "sync-controller/pkg/apis/keti/v1alpha1"
	//corev1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/rest"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
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

	// fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.Sync{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	//for _, ghost := range ghosts {
	//	fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
	//	if err := co.WatchResourceReconcileController(ghost, &appsv1.Deployment{}, controller.WatchOptions{}); err != nil {
	//		return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
	//	}
	//}
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

	// Fetch the Sync instance
	instance := &ketiv1alpha1.Sync{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		return reconcile.Result{}, nil
	}

	fmt.Println("instance Name: ", instance.Name)
	fmt.Println("instance Namespace : ", instance.Namespace)

	// Instance 삭제
	err = r.live.Delete(context.TODO(), instance)
	if err != nil {
		fmt.Println("Delete Err", err)
	}

	obj, clusterName, command := r.resourceForSync(instance)

	if command == "create" {
		err := CreateObj(cm, obj, clusterName)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				fmt.Println(err)
			} else {
				return reconcile.Result{}, err
			}
		} else {
			fmt.Println("Created Resource '" + obj.GetKind() + "' in Cluster'" + clusterName + "'")
			fmt.Println("  Name : " + obj.GetName())
			fmt.Println("  Namespace : " + obj.GetNamespace())
			fmt.Println()
		}
	} else if command == "update" {
		err := UpdateObj(cm, obj, clusterName)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				fmt.Println(err)
			} else {
				return reconcile.Result{}, err
			}
		} else {
			fmt.Println("Updated Resource '" + obj.GetKind() + "' in Cluster'" + clusterName + "'")
			fmt.Println("  Name : " + obj.GetName())
			fmt.Println("  Namespace : " + obj.GetNamespace())
			fmt.Println()
		}
	} else if command == "delete" { // Delete
		err := DeleteObj(cm, obj, clusterName)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				fmt.Println(err)
			} else {
				return reconcile.Result{}, err
			}
		} else {
			fmt.Println("Deleted Resource '" + obj.GetKind() + "' in Cluster'" + clusterName + "'")
			fmt.Println("  Name : " + obj.GetName())
			fmt.Println("  Namespace : " + obj.GetNamespace())
			fmt.Println()
		}
	} else {
		fmt.Println("Command '" + command + "' is not a valid command.")
	}

	return reconcile.Result{}, nil // err
}

func (r *reconciler) resourceForSync(instance *ketiv1alpha1.Sync) (*unstructured.Unstructured, string, string) {

	clusterName := instance.Spec.ClusterName
	command := instance.Spec.Command

	u := &unstructured.Unstructured{}

	fmt.Println(instance.Spec.ClusterName)
	var err error
	u.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&instance.Spec.Template)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(u.GetName(), u.GetNamespace())

	return u, clusterName, command
}
func CreateObj(cm *ClusterManager, obj *unstructured.Unstructured, clusterName string) error {
	//deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	gvk := obj.GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	clientset := cm.Cluster_kubeClients[clusterName]
	groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)
	mapping, err := rm.RESTMapping(gk, gvk.Version)

	//fmt.Println(mapping.Resource.Group, mapping.Resource.Version, mapping.Resource.Resource)
	_, err = cm.Cluster_dynClients[clusterName].Resource(mapping.Resource).Namespace(obj.GetNamespace()).Create(obj, metav1.CreateOptions{})
	return err

}
func UpdateObj(cm *ClusterManager, obj *unstructured.Unstructured, clusterName string) error {
	//deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	gvk := obj.GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	clientset := cm.Cluster_kubeClients[clusterName]
	groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)
	mapping, err := rm.RESTMapping(gk, gvk.Version)

	// fmt.Println(mapping.Resource.Group, mapping.Resource.Version, mapping.Resource.Resource)
	_, err = cm.Cluster_dynClients[clusterName].Resource(mapping.Resource).Namespace(obj.GetNamespace()).Update(obj, metav1.UpdateOptions{})
	return err

}
func DeleteObj(cm *ClusterManager, obj *unstructured.Unstructured, clusterName string) error {
	//deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	gvk := obj.GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}

	fmt.Println(gvk.Kind, gvk.Group, gvk.Version)
	fmt.Println(obj.GetName(), obj.GetNamespace())
	clientset := cm.Cluster_kubeClients[clusterName]
	groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)
	mapping, err := rm.RESTMapping(gk, gvk.Version)
	found, err := cm.Cluster_dynClients[clusterName].Resource(mapping.Resource).Namespace(obj.GetNamespace()).Get(obj.GetName(), metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		// all good
		fmt.Println("Not Found")
		return nil
	}
	if !isInObject(found, "OpenMCP") {
		return nil
	}

	// fmt.Println(mapping.Resource.Group, mapping.Resource.Version, mapping.Resource.Resource)
	err = cm.Cluster_dynClients[clusterName].Resource(mapping.Resource).Namespace(obj.GetNamespace()).Delete(obj.GetName(), &metav1.DeleteOptions{})
	return err

}
func isInObject(obj *unstructured.Unstructured, subString string) bool {
	refKind_str := obj.GetAnnotations()["multicluster.admiralty.io/controller-reference"]
	refKind_map := make(map[string]interface{})
	err := json.Unmarshal([]byte(refKind_str), &refKind_map)
	if err != nil {
		panic(err)
	}
	if strings.Contains(fmt.Sprintf("%v", refKind_map["kind"]), subString) {
		return true
	}
	return false
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
func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
