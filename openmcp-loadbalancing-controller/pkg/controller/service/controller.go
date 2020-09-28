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
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/serviceregistry"
	resourceapis "openmcp/openmcp/openmcp-resource-controller/apis"
	resourcev1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)



var cm *clusterManager.ClusterManager
func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Function Called NewController")
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

	instance := &resourcev1alpha1.OpenMCPService{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	omcplog.V(3).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] instance Name: ", instance.Name)
	omcplog.V(3).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] instance Namespace : ", instance.Namespace)

	serviceName := instance.Name

	// delete
	if err != nil && errors.IsNotFound(err) {
		if errors.IsNotFound(err) {
			omcplog.V(2).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Delete Service Registry")
			serviceregistry.Registry.Delete(loadbalancing.ServiceRegistry, serviceName)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, nil
	} else { 
		target, err := serviceregistry.Registry.Lookup(loadbalancing.ServiceRegistry, serviceName)

		if target != nil && !errors.IsNotFound(err) {
			omcplog.V(2).Info("[OpenMCP Loadbalancing Controller(Service Watch Controller)] Update Registry")
			serviceregistry.Registry.Delete(loadbalancing.ServiceRegistry, serviceName)
			for _, cluster := range cm.Cluster_list.Items {
				cluster_client := cm.Cluster_genClients[cluster.Name]
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

