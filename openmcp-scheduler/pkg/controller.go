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

package openmcpscheduler // import "admiralty.io/multicluster-controller/examples/openmcpscheduler/pkg/controller/openmcpscheduler"

import (
	"context"
	"fmt"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/openmcp-resource-controller/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	syncapis "openmcp/openmcp/openmcp-sync-controller/pkg/apis"
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/klog"
	appsv1 "k8s.io/api/apps/v1"
	"openmcp/openmcp/openmcp-scheduler/pkg/protobuf"
)

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, grpcServer protobuf.RequestAnalysisClient) (*controller.Controller, error) {
	klog.V(4).Info("[OpenMCP Deployment] Function Called NewController")
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

	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace, grpcServer: grpcServer}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}
	if err := syncapis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPDeployment{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	for _, ghost := range ghosts {
		if err := co.WatchResourceReconcileController(ghost, &appsv1.Deployment{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}
	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
	grpcServer		protobuf.RequestAnalysisClient
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	klog.Infof("*********** Reconcile ***********")
	cm := clusterManager.NewClusterManager()

	// Fetch the OpenMCPDeployment instance
	instance := &ketiv1alpha1.OpenMCPDeployment{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	klog.V(0).Infof("Resource Get => [Name] : %v, [Namespace]  : %v", instance.Name, instance.Namespace)

	if err != nil {
		return reconcile.Result{}, err
	}

	if instance.Status.CreateSyncRequestComplete == false {
		if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false {	
			
			cluster_replicas_map := Scheduling(cm, instance, r.grpcServer)
			klog.Infof("Check cluster's replicas_map :  %v", cluster_replicas_map)

			instance.Status.ClusterMaps = cluster_replicas_map
			instance.Status.Replicas = instance.Spec.Replicas

			instance.Status.SchedulingNeed = false
			instance.Status.SchedulingComplete = true

			// update OpenMCPDeployment to deploy
			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				klog.V(0).Infof("Failed to update instance status, %v", err)
				return reconcile.Result{}, err
			}

		} 
	}

	return reconcile.Result{}, nil
}
