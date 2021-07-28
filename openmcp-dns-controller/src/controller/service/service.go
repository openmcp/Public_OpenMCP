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

package service // import "admiralty.io/multicluster-controller/examples/serviceDNS/pkg/controller/serviceDNS"

import (
	"context"
	"fmt"
	"openmcp/openmcp/apis"
	dnsv1alpha1 "openmcp/openmcp/apis/dns/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-dns-controller/src/controller/serviceDNSRecord"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info(">>> OpenMCPService NewController()")
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

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &corev1.Service{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	/*
		for _, ghost := range ghosts {
			if err := co.WatchResourceReconcileController(context.TODO(), ghost, &corev1.Service{}, controller.WatchOptions{}); err != nil {
				return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
			}
		}
	*/

	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

var i int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {

	if req.Namespace == "openmcp" && req.Name == "openmcp-apiserver" {
		omcplog.V(4).Info("Function Called Reconcile")
		i += 1
		omcplog.V(5).Info("********* [ OpenMCP OpenMCPService", i, "] *********")
		omcplog.V(5).Info(req.Context, " / ", req.Namespace, " / ", req.Name)

		instance_default_domain := &dnsv1alpha1.OpenMCPDomain{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default-domain",
				Namespace: "kube-federation-system",
			},
			Domain: "openmcp.example.org",
		}

		err := r.live.Create(context.TODO(), instance_default_domain)
		if err != nil && errors.IsAlreadyExists(err) {
			omcplog.V(0).Info(err)
		} else if err != nil {
			omcplog.V(0).Info(err)
		}

		instance_osvcdnsr := &dnsv1alpha1.OpenMCPServiceDNSRecord{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
			Spec: dnsv1alpha1.OpenMCPServiceDNSRecordSpec{
				DomainRef: "default-domain",
				RecordTTL: 300,
			},
		}
		err2 := r.live.Get(context.TODO(), req.NamespacedName, instance_osvcdnsr)
		if err2 != nil && errors.IsNotFound(err2) {
			omcplog.V(0).Info("[OpenMCPAPIServer] Create OpenMCPServiceDNSDomain")
			err3 := r.live.Create(context.TODO(), instance_osvcdnsr)
			if err3 != nil {
				omcplog.V(0).Info(err3)
				return reconcile.Result{}, err3
			}
		}
		serviceDNSRecord.FillStatusInOpenMCP(instance_osvcdnsr, instance_default_domain)
		omcplog.V(0).Info("[OpenMCPAPIServer] Update OpenMCPServiceDNSDomain")
		err4 := r.live.Status().Update(context.TODO(), instance_osvcdnsr)
		if err4 != nil {
			omcplog.V(0).Info(err4)
			return reconcile.Result{}, err4
		}

	}

	return reconcile.Result{}, nil
}
