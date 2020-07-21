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

package externalDNS

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog"
	"openmcp/openmcp/openmcp-dns-controller/pkg/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-dns-controller/pkg/apis/keti/v1alpha1"
	"openmcp/openmcp/openmcp-dns-controller/pkg/mypdns"
	"openmcp/openmcp/util/clusterManager"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string) (*controller.Controller, error) {
	klog.V(0).Info( ">>> externalDNS NewController()")
	cm = clusterManager.NewClusterManager()

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

	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPDNSEndpoint{}, controller.WatchOptions{}); err != nil {
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
	klog.V(0).Info( "********* [ OpenMCP External DNS", i, "] *********")
	klog.V(0).Info( req.Context, " / ", req.Namespace, " / ", req.Name)
	//	cm := clusterManager.NewClusterManager()

	pdnsClient, err := mypdns.PdnsNewClient()
	if err != nil {
		klog.V(0).Info( err)
	}

	// Fetch the Sync instance
	instance := &ketiv1alpha1.OpenMCPDNSEndpoint{}
	err = r.live.Get(context.TODO(), req.NamespacedName, instance)


	if err != nil && errors.IsNotFound(err){
		err = mypdns.DeleteZone(pdnsClient, r.live)
		if err != nil {
			klog.V(0).Info( "[OpenMCP External DNS Controller] : Delete?  ",err)
		}
		return reconcile.Result{}, nil
	}


	for _, domain := range instance.Spec.Domains{
		if domain == ""{
			continue
		}
		err = mypdns.SyncZone(pdnsClient, domain, instance.Spec.Endpoints)
		if err != nil {
			return reconcile.Result{}, err
		}

		//zone, err := mypdns.GetZone(pdnsClient, domain)
		//if err != nil {
		//	klog.V(0).Info( "[OpenMCP External DNS Controller] : '", domain + "' Not Found")
		//}
		//if err == nil {
		//	// Already Exist
		//	err = mypdns.UpdateZoneWithRecords(pdnsClient, *zone, domain, instance.Spec.Endpoints)
		//	if err != nil {
		//		klog.V(0).Info( "[OpenMCP External DNS Controller] : UpdateZone?  ",err)
		//	}
		//} else {
		//	err = mypdns.CreateZoneWithRecords(pdnsClient, domain, instance.Spec.Endpoints)
		//	if err != nil {
		//		klog.V(0).Info( "[OpenMCP External DNS Controller] : CreateZone? ",err)
		//	}
		//}
	}




	return reconcile.Result{}, nil // err
}
