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
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/restmapper"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-sync-controller/pkg/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-sync-controller/pkg/apis/keti/v1alpha1"
	"openmcp/openmcp/util/clusterManager"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	cm = myClusterManager

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
	omcplog.V(4).Info("[OpenMCP Sync] Function Called Reconcile")
	i += 1

	// Fetch the Sync instance
	instance := &ketiv1alpha1.Sync{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		return reconcile.Result{}, nil
	}
	omcplog.V(5).Info("Resource Get => [Name] : "+ instance.Name + " [Namespace]  : " + instance.Namespace)


	// Instance 삭제
	err = r.live.Delete(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	omcplog.V(2).Info("Resource Extract from SyncResource")
	obj, clusterName, command := r.resourceForSync(instance)

	if command == "create" {
		omcplog.V(2).Info("Create Resource Start")
		err := CreateObj(cm, obj, clusterName)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				omcplog.V(0).Info(err)
			} else {
				return reconcile.Result{}, err
			}
		} else {
			omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "' in Cluster'" + clusterName + "'")
			omcplog.V(3).Info("  Name : " + obj.GetName())
			omcplog.V(3).Info("  Namespace : " + obj.GetNamespace())
			omcplog.V(2).Info()
		}
	} else if command == "update" {
		omcplog.V(2).Info("Update Resource Start")
		err := UpdateObj(cm, obj, clusterName)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				omcplog.V(0).Info(err)
			} else {
				return reconcile.Result{}, err
			}
		} else {
			omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "' in Cluster'" + clusterName + "'")
			omcplog.V(3).Info("  Name : " + obj.GetName())
			omcplog.V(3).Info("  Namespace : " + obj.GetNamespace())
			omcplog.V(2).Info()
		}
	} else if command == "delete" { // Delete
		omcplog.V(2).Info("Delete Resource Start")
		err := DeleteObj(cm, obj, clusterName)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				omcplog.V(0).Info(err)
			} else {
				return reconcile.Result{}, err
			}
		} else {
			omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "' in Cluster'" + clusterName + "'")
			omcplog.V(3).Info("  Name : " + obj.GetName())
			omcplog.V(3).Info("  Namespace : " + obj.GetNamespace())
			omcplog.V(2).Info()
		}
	} else {
		omcplog.V(1).Info("Command '" + command + "' is not a valid command.")
	}

	return reconcile.Result{}, nil // err
}

func (r *reconciler) resourceForSync(instance *ketiv1alpha1.Sync) (*unstructured.Unstructured, string, string) {
	omcplog.V(4).Info("[OpenMCP Sync] Function Called resourceForSync")
	clusterName := instance.Spec.ClusterName
	command := instance.Spec.Command

	u := &unstructured.Unstructured{}

	omcplog.V(2).Info("[Parsing Sync] ClusterName : ", clusterName, "command : ",command)
	var err error
	u.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&instance.Spec.Template)
	if err != nil {
		omcplog.V(0).Info(err)
	}
	omcplog.V(4).Info(u.GetName(), u.GetNamespace())

	return u, clusterName, command
}
func CreateObj(cm *clusterManager.ClusterManager, obj *unstructured.Unstructured, clusterName string) error {
	omcplog.V(4).Info("[OpenMCP Sync] Function Called CreateObj")
	//deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	gvk := obj.GroupVersionKind()

	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}

	clientset := cm.Cluster_kubeClients[clusterName]
	groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	mapping, err := rm.RESTMapping(gk, gvk.Version)

	if err != nil {
		omcplog.V(0).Info(err)
	}
	//omcplog.V(0).Info(mapping.Resource.Group, mapping.Resource.Version, mapping.Resource.Resource)
	_, err = cm.Cluster_dynClients[clusterName].Resource(mapping.Resource).Namespace(obj.GetNamespace()).Create(obj, metav1.CreateOptions{})
	return err

}
func UpdateObj(cm *clusterManager.ClusterManager, obj *unstructured.Unstructured, clusterName string) error {
	omcplog.V(4).Info("[OpenMCP Sync] Function Called UpdateObj")
	//deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	gvk := obj.GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	clientset := cm.Cluster_kubeClients[clusterName]
	groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)
	mapping, err := rm.RESTMapping(gk, gvk.Version)

	// omcplog.V(0).Info(mapping.Resource.Group, mapping.Resource.Version, mapping.Resource.Resource)
	_, err = cm.Cluster_dynClients[clusterName].Resource(mapping.Resource).Namespace(obj.GetNamespace()).Update(obj, metav1.UpdateOptions{})
	return err

}
func DeleteObj(cm *clusterManager.ClusterManager, obj *unstructured.Unstructured, clusterName string) error {
	omcplog.V(4).Info("[OpenMCP Sync] Function Called DeleteObj")
	//deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	gvk := obj.GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}

	omcplog.V(5).Info(gvk.Kind, gvk.Group, gvk.Version)
	omcplog.V(5).Info(obj.GetName(), obj.GetNamespace())
	clientset := cm.Cluster_kubeClients[clusterName]
	groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	mapping, err := rm.RESTMapping(gk, gvk.Version)
	found, err := cm.Cluster_dynClients[clusterName].Resource(mapping.Resource).Namespace(obj.GetNamespace()).Get(obj.GetName(), metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		// all good
		omcplog.V(0).Info("Not Found")
		return nil
	}
	if !isInObject(found, "OpenMCP") {
		return nil
	}

	// omcplog.V(0).Info(mapping.Resource.Group, mapping.Resource.Version, mapping.Resource.Resource)
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
