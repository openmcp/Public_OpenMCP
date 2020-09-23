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
	vpav1beta2 "github.com/kubernetes/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	appsv1 "k8s.io/api/apps/v1"
	hpav2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
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
	ghostclients := map[string]client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients[ghost.Name] = ghostclient
	}
	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})

	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}
	if err := vpav1beta2.AddToScheme(live.GetScheme()); err != nil {
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
	ghosts         map[string]client.Client
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

	// do(cm, obj, clusterName, command)
	jsonbody, err := json.Marshal(obj)
	if err != nil {
		// do error check
		fmt.Println(err)
		return reconcile.Result{}, err
	}
	//clusterClient := cm.Cluster_clusterClients[clusterName]
	clusterClient := r.ghosts[clusterName]



	if obj.GetKind() == "Deployment"{
		subInstance := &appsv1.Deployment{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create"{
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Create Deployment : ", err)
			}
		} else if command == "delete"{
			err = clusterClient.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Delete Deployment : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Update Deployment : ", err)
			}
		}


	} else if obj.GetKind() == "Service"{
		subInstance := &corev1.Service{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create"{
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Create Service : ", err)
			}
		} else if command == "delete"{
			err = clusterClient.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Delete Service : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Update Service : ", err)
			}
		}

	} else if obj.GetKind() == "Ingress"{
		subInstance := &extv1b1.Ingress{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create"{
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Create Ingress : ", err)
			}
		} else if command == "delete"{
			err = clusterClient.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Delete Ingress : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Update Ingress : ", err)
			}
		}

	} else if obj.GetKind() == "HorizontalPodAutoscaler"{
		subInstance := &hpav2beta2.HorizontalPodAutoscaler{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create"{
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Create HorizontalPodAutoscaler : ", err)
			}
		} else if command == "delete"{
			err = clusterClient.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Delete HorizontalPodAutoscaler : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Update HorizontalPodAutoscaler : ", err)
			}
		}

	} else if obj.GetKind() == "VerticalPodAutoscaler"{
		subInstance := &vpav1beta2.VerticalPodAutoscaler{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}

		if command == "create"{
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Create VerticalPodAutoscaler : ", err)
			}
		} else if command == "delete"{
			err = clusterClient.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Delete VerticalPodAutoscaler : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Update VerticalPodAutoscaler : ", err)
			}
		}

	} else if obj.GetKind() == "ConfigMap"{
		subInstance := &corev1.ConfigMap{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create"{
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Create ConfigMap : ", err)
			}
		} else if command == "delete"{
			err = clusterClient.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Delete ConfigMap : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Update ConfigMap : ", err)
			}
		}

	} else if obj.GetKind() == "Secret"{
		subInstance := &corev1.Secret{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create"{
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Create Secret : ", err)
			}
		} else if command == "delete"{
			err = clusterClient.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Delete Secret : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Update Secret : ", err)
			}
		}

	}else if obj.GetKind() == "PersistentVolume"{
		subInstance := &corev1.PersistentVolume{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create"{
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Create PersistentVolume : ", err)
			}
		} else if command == "delete"{
			err = clusterClient.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Delete PersistentVolume : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Update PersistentVolume : ", err)
			}
		}

	} else if obj.GetKind() == "PersistentVolumeClaim"{
		subInstance := &corev1.PersistentVolumeClaim{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create"{
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Create PersistentVolumeClaim : ", err)
			}
		} else if command == "delete"{
			err = clusterClient.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Delete PersistentVolumeClaim : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() +  "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() +"', in Cluster'" + clusterName + "'")
			}else {
				omcplog.V(0).Info("[Error] Cannot Update PersistentVolumeClaim : ", err)
			}
		}

	}


	return reconcile.Result{}, nil // err
}

func (r *reconciler) resourceForSync(instance *ketiv1alpha1.Sync) (*unstructured.Unstructured, string, string) {
	omcplog.V(4).Info("[OpenMCP Sync] Function Called resourceForSync")
	clusterName := instance.Spec.ClusterName
	command := instance.Spec.Command

	u := &unstructured.Unstructured{}


	omcplog.V(2).Info("[Parsing Sync] ClusterName : ", clusterName, ", command : ",command)
	var err error
	u.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&instance.Spec.Template)
	if err != nil {
		omcplog.V(0).Info(err)
	}
	omcplog.V(4).Info(u.GetName()," / ", u.GetNamespace())

	return u, clusterName, command
}
func CreateObj(cm *clusterManager.ClusterManager, obj *unstructured.Unstructured, clusterName string) error {
	omcplog.V(4).Info("[OpenMCP Sync] Function Called CreateObj")
	//deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	fmt.Println("check1")
	gvk := obj.GroupVersionKind()
	fmt.Println("check2")
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	fmt.Println("check3")
	clientset := cm.Cluster_kubeClients[clusterName]
	groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	fmt.Println("check4")
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)
	fmt.Println("check5")

	mapping, err := rm.RESTMapping(gk, gvk.Version)
	fmt.Println("check6")

	if err != nil {
		omcplog.V(0).Info(err)
	}
	//omcplog.V(0).Info(mapping.Resource.Group, mapping.Resource.Version, mapping.Resource.Resource)
	_, err = cm.Cluster_dynClients[clusterName].Resource(mapping.Resource).Namespace(obj.GetNamespace()).Create(obj, metav1.CreateOptions{})
	fmt.Println("check7")
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
