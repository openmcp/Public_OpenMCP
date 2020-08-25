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

package openmcpsecret // import "admiralty.io/multicluster-controller/examples/openmcpsecret/pkg/controller/openmcpsecret"

import (
	"context"
	"fmt"
	//"k8s.io/klog"
	"openmcp/openmcp/omcplog"
	"strconv"

	"admiralty.io/multicluster-controller/pkg/reference"
	//"reflect"
	//"sort"
	//"math/rand"
	//"time"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/kubefed/pkg/controller/util"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"openmcp/openmcp/openmcp-resource-controller/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	//corev1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	//appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/getlantern/deepcopy"
	"k8s.io/client-go/rest"
	syncapis "openmcp/openmcp/openmcp-sync-controller/pkg/apis"
	sync "openmcp/openmcp/openmcp-sync-controller/pkg/apis/keti/v1alpha1"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
)
type ClusterManager struct {
        Fed_namespace string
        Host_config *rest.Config
        Host_client genericclient.Client
        Cluster_list *fedv1b1.KubeFedClusterList
        Cluster_configs map[string]*rest.Config
        Cluster_clients map[string]genericclient.Client
}



func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string) (*controller.Controller, error) {
	omcplog.V(4).Info("Function Called NewController")
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
	if err := syncapis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPSecret{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	for _, ghost := range ghosts {
		if err := co.WatchResourceReconcileController(ghost, &corev1.Secret{}, controller.WatchOptions{}); err != nil {
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
	omcplog.V(4).Info("Function Called Reconcile")
	i += 1
	omcplog.V(5).Info("********* [",i,"] *********")
	omcplog.V(3).Info(req.Context," / ", req.Namespace," / ", req.Name)
	cm := NewClusterManager()

	// Fetch the OpenMCPDeployment instance
        instance := &ketiv1alpha1.OpenMCPSecret{}
        err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	omcplog.V(3).Info("instance Name: ", instance.Name)
	omcplog.V(3).Info("instance Namespace : ", instance.Namespace)

	if err != nil {
		if errors.IsNotFound(err) {
			// ...TODO: multicluster garbage collector
			// Until then...
			omcplog.V(3).Info("Delete Deployments ..Cluster")
			err := r.DeleteSecret(cm, req.NamespacedName.Name, req.NamespacedName.Namespace)
			return reconcile.Result{}, err
		}
		omcplog.V(1).Info(err)
		return reconcile.Result{}, err
	}

	if instance.Status.ClusterMaps == nil {
		err := r.createSecret(req, cm, instance)
		if err != nil {
			omcplog.V(1).Info(err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	} else {
		err := r.updateSecret(req, cm, instance)
		if err != nil {
			omcplog.V(1).Info(err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	err = r.live.Status().Update(context.TODO(), instance)
	if err != nil {
		omcplog.V(1).Info("Failed to update instance status", err)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *reconciler) secretForOpenMCPSecret(req reconcile.Request, m *ketiv1alpha1.OpenMCPSecret) *corev1.Secret {
	omcplog.V(4).Info("Function Called secretForOpenMCPSecret")

        dep := &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Name,
				Namespace: m.Namespace,
			},
                //Spec: m.Spec.Template.Spec,
        }
        deepcopy.Copy(&dep.Data, &m.Spec.Template.Data)
        //dep.Spec.Replicas = &replica

	reference.SetMulticlusterControllerReference(dep, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return dep
}



func (r *reconciler) createSecret(req reconcile.Request, cm *ClusterManager, instance *ketiv1alpha1.OpenMCPSecret) error {
	omcplog.V(4).Info("Function Called createSecret")
	cluster_map := make(map[string]int32)
	for _, cluster := range cm.Cluster_list.Items {

		omcplog.V(3).Info("Cluster '" + cluster.Name + "' Deployed")
		dep := r.secretForOpenMCPSecret(req, instance)
		command := "create"
		_, err := r.sendSync(dep, command, cluster.Name)
		cluster_map[cluster.Name] = 1
		if err != nil {
			return err
		}
	}
	instance.Status.ClusterMaps = cluster_map
	err := r.live.Status().Update(context.TODO(), instance)
	return err
}


func (r *reconciler) updateSecret(req reconcile.Request, cm *ClusterManager, instance *ketiv1alpha1.OpenMCPSecret) error {
	omcplog.V(4).Info("Function Called updateSecret")

	for _, cluster := range cm.Cluster_list.Items {

		omcplog.V(3).Info("Cluster '" + cluster.Name + "' Deployed")
		dep := r.secretForOpenMCPSecret(req, instance)
		command := "update"
		_, err := r.sendSync(dep, command, cluster.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *reconciler) DeleteSecret(cm *ClusterManager, name string, namespace string) error {
	omcplog.V(4).Info("Function Called DeleteSecret")

	for _, cluster := range cm.Cluster_list.Items {

		omcplog.V(3).Info(cluster.Name," Delete Start")

		dep := &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
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


func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	omcplog.V(4).Info("Function Called ListKubeFedClusters")
        clusterList := &fedv1b1.KubeFedClusterList{}
        err := client.List(context.TODO(), clusterList, namespace)
        if err != nil {
			omcplog.V(1).Info("Error retrieving list of federated clusters: %+v", err)
        }
        if len(clusterList.Items) == 0 {
			omcplog.V(1).Info("No federated clusters found")
        }
        return clusterList
}

func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, client genericclient.Client, fedNamespace string) map[string]*rest.Config {
	omcplog.V(4).Info("Function Called KubeFedClusterConfigs")
        clusterConfigs := make(map[string]*rest.Config)
        for _, cluster := range clusterList.Items {
                config, _ := util.BuildClusterConfig(&cluster, client, fedNamespace)
                clusterConfigs[cluster.Name] = config
        }
        return clusterConfigs
}
func KubeFedClusterClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {
	omcplog.V(4).Info("Function Called KubeFedClusterClients")

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
	omcplog.V(4).Info("Function Called NewClusterManager")
        fed_namespace := "kube-federation-system"
        host_config, _ := rest.InClusterConfig()
        host_client := genericclient.NewForConfigOrDie(host_config)
        cluster_list := ListKubeFedClusters(host_client, fed_namespace)
        cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)
        cluster_clients := KubeFedClusterClients(cluster_list, cluster_configs)

        cm := &ClusterManager{
                Fed_namespace: fed_namespace,
                Host_config: host_config,
                Host_client: host_client,
                Cluster_list: cluster_list,
                Cluster_configs: cluster_configs,
                Cluster_clients: cluster_clients,
        }
        return cm
}


var syncIndex int = 0
func (r *reconciler) sendSync(secret *corev1.Secret, command string, clusterName string) (string, error) {
	omcplog.V(4).Info("Function Called sendSync")

	syncIndex += 1

	s := &sync.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-secret-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: sync.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *secret,
		},
	}
	omcplog.V(5).Info("Delete Check ", s.Spec.Template.(corev1.Secret).Name, s.Spec.Template.(corev1.Secret).Namespace)

	err := r.live.Create(context.TODO(), s)

	if err != nil {
		omcplog.V(0).Info(err)
	}

	omcplog.V(0).Info(s.Name)
	return s.Name, err
}


