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

package controller

import (
	"context"
	"fmt"
	"openmcp/openmcp/apis"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	syncv1alpha1 "openmcp/openmcp/apis/sync/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"reflect"
	"strconv"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"admiralty.io/multicluster-controller/pkg/reference"
	"github.com/getlantern/deepcopy"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("Function Called NewController")
	cm = myClusterManager

	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}

	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			omcplog.V(4).Info("Error getting delegating client for ghost cluster [", ghost.Name, "]")
			//return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		} else {
			ghostclients = append(ghostclients, ghostclient)
		}
	}

	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &resourcev1alpha1.OpenMCPConfigMap{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	for _, ghost := range ghosts {
		fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(context.TODO(), ghost, &corev1.ConfigMap{}, controller.WatchOptions{}); err != nil {
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
var syncIndex int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(4).Info("Function Called Reconcile")
	i += 1
	omcplog.V(5).Info("********* [", i, "] *********")
	omcplog.V(3).Info(req.Context, "/", req.Namespace, "/", req.Name)

	instance := &resourcev1alpha1.OpenMCPConfigMap{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	omcplog.V(3).Info("instance Name: ", instance.Name)
	omcplog.V(3).Info("instance Namespace : ", instance.Namespace)

	if err != nil {
		if errors.IsNotFound(err) {
			omcplog.V(3).Info("Delete ConfigMap")

			err := r.DeleteConfigMap(cm, req.NamespacedName.Name, req.NamespacedName.Namespace)
			return reconcile.Result{}, err
		}
		omcplog.V(1).Info(err)
		return reconcile.Result{}, err
	}
	if !reflect.DeepEqual(instance.Status.LastSpec, instance.Spec) {
		omcplog.V(3).Info("Service Update Start")

		err := r.updateConfigMap(req, cm, instance)
		if err != nil {
			omcplog.V(0).Info(err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil

	}
	if instance.Status.ClusterMaps == nil {
		err := r.createConfigMap(req, cm, instance)
		if err != nil {
			omcplog.V(1).Info(err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	if instance.Status.CheckSubResource == true {
		omcplog.V(2).Info("[Member Cluster Check ConfigMap]")
		sync_req_name := ""
		for k, v := range instance.Status.ClusterMaps {
			cluster_name := k
			replica := v

			if v == 0 {
				continue
			}
			// if _, ok := cm.Cluster_genClients[cluster_name]; !ok {
			// 	r.updateConfigMap(req, cm, instance)
			// }
			found := &corev1.ConfigMap{}
			cluster_client := cm.Cluster_genClients[cluster_name]
			err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)

			if err != nil && errors.IsNotFound(err) {
				// Delete ConfigMap Detected
				omcplog.V(2).Info("Cluster '"+cluster_name+"' ReDeployed => ", replica)
				configmap := r.configmapForOpenMCPConfigMap(req, instance)

				command := "create"
				omcplog.V(3).Info("SyncResource Create (ClusterName : "+cluster_name+", Command : "+command+", Replicas :", replica, ")")
				sync_req_name, err = r.sendSync(configmap, command, cluster_name)

				if err != nil {
					return reconcile.Result{}, err
				}

			}

		}
		omcplog.V(3).Info("sync_req_name : ", sync_req_name)

		err = r.live.Status().Update(context.TODO(), instance)
		if err != nil {
			omcplog.V(0).Info("Failed to update instance status", err)
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *reconciler) configmapForOpenMCPConfigMap(req reconcile.Request, m *resourcev1alpha1.OpenMCPConfigMap) *corev1.ConfigMap {
	omcplog.V(4).Info("Function Called configmapForOpenMCPConfigMap")
	configmap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
	}
	deepcopy.Copy(&configmap.Data, &m.Spec.Template.Data)
	reference.SetMulticlusterControllerReference(configmap, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))
	return configmap
}

func (r *reconciler) createConfigMap(req reconcile.Request, cm *clusterManager.ClusterManager, instance *resourcev1alpha1.OpenMCPConfigMap) error {
	omcplog.V(4).Info("Function Called createConfigMap")
	cluster_map := make(map[string]int32)
	for _, cluster := range cm.Cluster_list.Items {
		cluster_map[cluster.Name] = 0
	}
	for _, cluster := range cm.Cluster_list.Items {
		found := &corev1.ConfigMap{}
		cluster_client := cm.Cluster_genClients[cluster.Name]

		err := cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
		if err != nil && errors.IsNotFound(err) {
			omcplog.V(3).Info("Cluster '" + cluster.Name + "' Deployed")
			dep := r.configmapForOpenMCPConfigMap(req, instance)
			command := "create"
			_, err := r.sendSync(dep, command, cluster.Name)
			cluster_map[cluster.Name] = 1
			if err != nil {
				omcplog.V(0).Info(err)
				return err
			}
		}
	}
	instance.Status.ClusterMaps = cluster_map
	instance.Status.LastSpec = instance.Spec
	omcplog.V(3).Info("Update Status")
	err := r.live.Status().Update(context.TODO(), instance)
	return err
}

func (r *reconciler) updateConfigMap(req reconcile.Request, cm *clusterManager.ClusterManager, instance *resourcev1alpha1.OpenMCPConfigMap) error {
	omcplog.V(4).Info("Function Called updateConfigMap")

	for _, cluster := range cm.Cluster_list.Items {

		configmap := r.configmapForOpenMCPConfigMap(req, instance)
		command := "update"
		omcplog.V(3).Info("Update ConfigMap")
		_, err := r.sendSync(configmap, command, cluster.Name)
		if err != nil {
			return err
		}

	}

	instance.Status.LastSpec = instance.Spec

	err := r.live.Status().Update(context.TODO(), instance)
	return err
}

func (r *reconciler) DeleteConfigMap(cm *clusterManager.ClusterManager, name string, namespace string) error {
	omcplog.V(4).Info("Function Called DeleteConfigMap")

	for _, cluster := range cm.Cluster_list.Items {

		omcplog.V(3).Info(cluster.Name, " Delete Start")

		dep := &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
		command := "delete"
		_, err := r.sendSync(dep, command, cluster.Name)

		if err != nil {
			return err
		}
		omcplog.V(3).Info(cluster.Name, "Delete Complete")
	}
	return nil
}

func (r *reconciler) sendSync(configmap *corev1.ConfigMap, command string, clusterName string) (string, error) {
	omcplog.V(4).Info("[OpenMCP ConfigMap] Function Called sendSync")
	syncIndex += 1

	s := &syncv1alpha1.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-configmap-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: syncv1alpha1.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *configmap,
		},
	}
	omcplog.V(5).Info("Delete Check ", s.Spec.Template.(corev1.ConfigMap).Name, s.Spec.Template.(corev1.ConfigMap).Namespace)

	err := r.live.Create(context.TODO(), s)

	if err != nil {
		omcplog.V(0).Info(err)
	}

	omcplog.V(0).Info(s.Name)
	return s.Name, err
}
