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
	"fmt"
	"openmcp/openmcp/apis"
	syncv1alpha1 "openmcp/openmcp/apis/sync/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	vpav1beta2 "github.com/kubernetes/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	appsv1 "k8s.io/api/apps/v1"
	hpav2beta2 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &syncv1alpha1.Sync{}, controller.WatchOptions{}); err != nil {
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
	ghosts         map[string]client.Client
	ghostNamespace string
}

var i int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(4).Info("[OpenMCP Sync] Function Called Reconcile")
	i += 1

	// Fetch the Sync instance
	instance := &syncv1alpha1.Sync{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		return reconcile.Result{}, nil
	}
	omcplog.V(5).Info("Resource Get => [Name] : " + instance.Name + " [Namespace]  : " + instance.Namespace)

	// Instance Delete
	err = r.live.Delete(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	omcplog.V(2).Info("Resource Extract from SyncResource")
	obj, clusterName, command := r.resourceForSync(instance)

	jsonbody, err := json.Marshal(obj)
	if err != nil {
		// do error check
		return reconcile.Result{}, err
	}
	//clusterClient := r.ghosts[clusterName]
	clusterClient := cm.Cluster_genClients[clusterName]
	clusterClient2 := r.ghosts[clusterName]

	if obj.GetKind() == "Deployment" {
		subInstance := &appsv1.Deployment{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create Deployment : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete Deployment : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update Deployment : ", err)
			}
		}

	} else if obj.GetKind() == "Service" {
		subInstance := &corev1.Service{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create Service : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete Service : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update Service : ", err)
			}
		}

	} else if obj.GetKind() == "Ingress" {
		subInstance := &extv1b1.Ingress{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create Ingress : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete Ingress : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update Ingress : ", err)
			}
		}

	} else if obj.GetKind() == "HorizontalPodAutoscaler" {
		subInstance := &hpav2beta2.HorizontalPodAutoscaler{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create HorizontalPodAutoscaler : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete HorizontalPodAutoscaler : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update HorizontalPodAutoscaler : ", err)
			}
		}

	} else if obj.GetKind() == "VerticalPodAutoscaler" {
		subInstance := &vpav1beta2.VerticalPodAutoscaler{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}

		if command == "create" {
			//err = clusterClient.Create(context.TODO(), subInstance)
			err = clusterClient2.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create VerticalPodAutoscaler : ", err)
			}
		} else if command == "delete" {
			//err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			err = clusterClient2.Delete(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete VerticalPodAutoscaler : ", err)
			}
		} else if command == "update" {
			//err = clusterClient.Update(context.TODO(), subInstance)
			err = clusterClient2.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update VerticalPodAutoscaler : ", err)
			}
		}

	} else if obj.GetKind() == "ConfigMap" {
		subInstance := &corev1.ConfigMap{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create ConfigMap : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete ConfigMap : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update ConfigMap : ", err)
			}
		}

	} else if obj.GetKind() == "Secret" {
		subInstance := &corev1.Secret{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create Secret : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete Secret : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update Secret : ", err)
			}
		}

	} else if obj.GetKind() == "PersistentVolume" {
		subInstance := &corev1.PersistentVolume{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create PersistentVolume : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete PersistentVolume : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update PersistentVolume : ", err)
			}
		}

	} else if obj.GetKind() == "PersistentVolumeClaim" {
		subInstance := &corev1.PersistentVolumeClaim{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create PersistentVolumeClaim : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete PersistentVolumeClaim : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update PersistentVolumeClaim : ", err)
			}
		}

	} else if obj.GetKind() == "StatefulSet" {
		subInstance := &appsv1.StatefulSet{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create StatefulSet : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete StatefulSet : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update StatefulSet : ", err)
			}
		}

	} else if obj.GetKind() == "DaemonSet" {
		subInstance := &appsv1.DaemonSet{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create DaemonSet : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete DaemonSet : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update DaemonSet : ", err)
			}
		}

	} else if obj.GetKind() == "Job" {
		subInstance := &batchv1.Job{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create Job : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete Job : ", err)
			}
		} else if command == "update" {
			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update Job : ", err)
			}
		}
	} else if obj.GetKind() == "Namespace" {
		subInstance := &corev1.Namespace{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return reconcile.Result{}, err
		}
		if command == "create" {
			subInstance.ResourceVersion = ""
			err = clusterClient.Create(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				if !errors.IsNotFound(err) {
					return reconcile.Result{}, err // err
				}
			} else {
				omcplog.V(0).Info("[Error] Cannot Create Namespace : ", err)
			}
		} else if command == "delete" {
			err = clusterClient.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
			if err == nil {
				omcplog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Delete Namespace : ", err)
			}
		} else if command == "update" {

			err = clusterClient.Update(context.TODO(), subInstance)
			if err == nil {
				omcplog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
			} else {
				omcplog.V(0).Info("[Error] Cannot Update Namespace : ", err)
			}
		}
	}

	return reconcile.Result{}, nil // err
}

func (r *reconciler) resourceForSync(instance *syncv1alpha1.Sync) (*unstructured.Unstructured, string, string) {
	omcplog.V(4).Info("[OpenMCP Sync] Function Called resourceForSync")
	clusterName := instance.Spec.ClusterName
	command := instance.Spec.Command

	u := &unstructured.Unstructured{}

	omcplog.V(2).Info("[Parsing Sync] ClusterName : ", clusterName, ", command : ", command)
	var err error
	u.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&instance.Spec.Template)
	if err != nil {
		omcplog.V(0).Info(err)
	}
	omcplog.V(4).Info(u.GetName(), " / ", u.GetNamespace())

	return u, clusterName, command
}
