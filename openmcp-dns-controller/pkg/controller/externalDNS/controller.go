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
	"openmcp/openmcp/apis"
	dnsv1alpha1 "openmcp/openmcp/apis/dns/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-dns-controller/pkg/mypdns"
	"openmcp/openmcp/util/clusterManager"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info( ">>> externalDNS NewController()")
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

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &dnsv1alpha1.OpenMCPDNSEndpoint{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
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
	omcplog.V(4).Info( "Function Called Reconcile")
	i += 1
	omcplog.V(5).Info( "********* [ OpenMCP Domain", i, "] *********")
	omcplog.V(5).Info( req.Context, " / ", req.Namespace, " / ", req.Name)

	omcplog.V(2).Info("PdnsNewClient")
	pdnsClient, err := mypdns.PdnsNewClient()
	if err != nil {
		omcplog.V(0).Info(err)
	}


	instance := &dnsv1alpha1.OpenMCPDNSEndpoint{}
	err = r.live.Get(context.TODO(), req.NamespacedName, instance)
	omcplog.V(2).Info("[Get] OpenMCPDNSEndpoint")

	if err != nil && errors.IsNotFound(err){
		omcplog.V(2).Info( "DNSEndpoint Delete Detection, PowerDNS Info Delete")
		err = mypdns.DeleteZone(pdnsClient, r.live)
		if err != nil {
			omcplog.V(0).Info( "[OpenMCP External DNS Controller] : Delete?  ",err)
		}
		return reconcile.Result{}, nil
	}


	for _, domain := range instance.Spec.Domains{
		if domain == ""{
			continue
		}
		omcplog.V(2).Info( "DNSEndpoint Create/Update Detection, PowerDNS Info Refresh")
		err = mypdns.SyncZone(pdnsClient, domain, instance.Spec.Endpoints)
		if err != nil {
			return reconcile.Result{}, err
		}

	}




	return reconcile.Result{}, nil // err
}
