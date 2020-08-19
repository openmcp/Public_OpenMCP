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

package service

import (
	"context"
	"fmt"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/apis"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/controller/util"

	corev1 "k8s.io/api/core/v1"
	resourceapis "openmcp/openmcp/openmcp-resource-controller/apis"
	resourcev1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"

	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/serviceregistry"
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
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Function Called NewController")
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

	if err := resourceapis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}


	if err := co.WatchResourceReconcileObject(live, &resourcev1alpha1.OpenMCPService{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	for _, ghost := range ghosts {
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
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Function Called Reconcile")
	i += 1
	omcplog.V(3).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)]********* [", i, "] *********")
	omcplog.V(3).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)]" + req.Context, " / ", req.Namespace, " / ", req.Name)
	cm := NewClusterManager()

	//instance := &corev1.Service{}
	instance := &resourcev1alpha1.OpenMCPService{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	omcplog.V(3).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] instance Name: ", instance.Name)
	omcplog.V(3).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] instance Namespace : ", instance.Namespace)

	serviceName := instance.Name

	// delete
	if err != nil && errors.IsNotFound(err) {
		//해당 instance가 없을 경우 ingress & ingressName Registry 삭제
		if errors.IsNotFound(err) {
			omcplog.V(2).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Delete Service Registry")
			serviceregistry.Registry.Delete(loadbalancing.ServiceRegistry, serviceName)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, nil
	} else { // 해당 instance가 있을 경우 ingress & ingressName Registry Add or Update
	//	add
		target, err := serviceregistry.Registry.Lookup(loadbalancing.ServiceRegistry, serviceName)

		if target != nil && !errors.IsNotFound(err) {
			omcplog.V(2).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Update Registry")
			serviceregistry.Registry.Delete(loadbalancing.ServiceRegistry, serviceName)
			for _, cluster := range cm.Cluster_list.Items {
				cluster_client := cm.Cluster_clients[cluster.Name]
				found := &corev1.Service{}
				err := cluster_client.Get(context.TODO(), found, "openmcp", serviceName)
				if err != nil && errors.IsNotFound(err) {
					omcplog.V(0).Info(err)
					omcplog.V(0).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Service Not Found")
				} else { // Add
					omcplog.V(2).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Service Registry Add")
					serviceregistry.Registry.Add(loadbalancing.ServiceRegistry, serviceName, cluster.Name)
				}
			}
		} else {
			omcplog.V(0).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Not Ingress Service")
		}
	}
	return reconcile.Result{}, nil
}


func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Function Called ListKubeFedClusters")
	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		omcplog.V(0).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Error retrieving list of federated clusters: %+v", err)
	}
	if len(clusterList.Items) == 0 {
		omcplog.V(0).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] No federated clusters found")
	}
	return clusterList
}

func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, client genericclient.Client, fedNamespace string) map[string]*rest.Config {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Function Called KubeFedClusterConfigs")
	clusterConfigs := make(map[string]*rest.Config)
	for _, cluster := range clusterList.Items {
		config, _ := util.BuildClusterConfig(&cluster, client, fedNamespace)
		clusterConfigs[cluster.Name] = config
	}
	return clusterConfigs
}
func KubeFedClusterClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Function Called KubeFedClusterClients")

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
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Function Called NewClusterManager")
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
