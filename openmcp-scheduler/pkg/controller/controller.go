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

package controller // import "admiralty.io/multicluster-controller/examples/openmcpscheduler/pkg/controller/openmcpscheduler"

import (
	"context"
	"fmt"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/manager"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-resource-controller/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	"openmcp/openmcp/openmcp-scheduler/pkg"
	syncapis "openmcp/openmcp/openmcp-sync-controller/pkg/apis"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"
)

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
	scheduler      *openmcpscheduler.OpenMCPScheduler
}

func NewControllers(cm *clusterManager.ClusterManager, scheduler *openmcpscheduler.OpenMCPScheduler) {
	host_ctx := "openmcp"
	namespace := "openmcp"

	host_cfg := cm.Host_config
	live := cluster.New(host_ctx, host_cfg, cluster.Options{})

	ghosts := []*cluster.Cluster{}

	for _, ghost_cluster := range cm.Cluster_list.Items {
		ghost_ctx := ghost_cluster.Name
		ghost_cfg := cm.Cluster_configs[ghost_ctx]

		ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{})
		ghosts = append(ghosts, ghost)
	}

	sched_cont, err := NewController(live, ghosts, namespace, scheduler)
	if err != nil {
		omcplog.V(0).Info("err New Controller - Scheduler", err)
	}
	reshape_cont, err := reshape.NewController(live, ghosts, namespace)
	if err != nil {
		omcplog.V(0).Info("err New Controller - Reshape", err)
	}
	loglevel_cont, err := logLevel.NewController(live, ghosts, namespace)
	if err != nil {
		omcplog.V(0).Info("err New Controller - logLevel", err)
	}

	m := manager.New()
	m.AddController(sched_cont)
	m.AddController(reshape_cont)
	m.AddController(loglevel_cont)

	stop := reshape.SetupSignalHandler()

	if err := m.Start(stop); err != nil {
		omcplog.V(0).Info(err)
	}
}

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, scheduler *openmcpscheduler.OpenMCPScheduler) (*controller.Controller, error) {
	omcplog.V(4).Info("[OpenMCP Deployment] Function Called NewController")

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

	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace, scheduler: scheduler}, controller.Options{})

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

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	cm := clusterManager.NewClusterManager()

	// Fetch the OpenMCPDeployment instance
	newDeployment := &ketiv1alpha1.OpenMCPDeployment{}
	err := r.live.Get(context.TODO(), req.NamespacedName, newDeployment)
	omcplog.V(0).Infof("Resource Get => [Name] : %v, [Namespace]  : %v", newDeployment.Name, newDeployment.Namespace)

	if err != nil {
		return reconcile.Result{}, err
	}

	// Start scheduling if scheduling is needed
	if newDeployment.Status.CreateSyncRequestComplete == false {
		if newDeployment.Status.SchedulingNeed == true && newDeployment.Status.SchedulingComplete == false {
			cluster_replicas_map, _ := r.scheduler.Scheduling(cm, newDeployment)
			omcplog.V(0).Info("Scheduling Result : ", cluster_replicas_map)

			newDeployment.Status.ClusterMaps = cluster_replicas_map
			newDeployment.Status.Replicas = newDeployment.Spec.Replicas

			newDeployment.Status.SchedulingNeed = false
			newDeployment.Status.SchedulingComplete = true

			// update OpenMCPDeployment to deploy
			err := r.live.Status().Update(context.TODO(), newDeployment)
			if err != nil {
				omcplog.V(0).Infof("Failed to update instance status, %v", err)
				return reconcile.Result{}, err
			}

		}
	}
	return reconcile.Result{}, nil
}
