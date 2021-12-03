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

package serviceInCluster // import "admiralty.io/multicluster-controller/examples/serviceDNS/pkg/controller/serviceDNS"

import (
	"context"
	"fmt"
	"openmcp/openmcp/apis"
	dnsv1alpha1 "openmcp/openmcp/apis/dns/v1alpha1"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-dns-controller/src/controller/serviceDNSRecord"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info(">>> ServiceDNS NewController()")
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
	for _, ghost := range ghosts {
		if err := co.WatchResourceReconcileController(context.TODO(), ghost, &corev1.Service{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
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
	omcplog.V(5).Info("********* [ ServiceInCluster Controller", i, "] *********")
	omcplog.V(5).Info(req.Context, " / ", req.Namespace, " / ", req.Name)

	omcplog.V(2).Info("Service Request")
	osvc := &resourcev1alpha1.OpenMCPService{}
	err := r.live.Get(context.TODO(), req.NamespacedName, osvc)
	if err != nil && errors.IsNotFound(err) {
		omcplog.V(2).Info("Deleted Service. Ignore Request")
		return reconcile.Result{}, nil
	}

	if err == nil {
		omcplog.V(2).Info("Create or Update Service.")
		omcplog.V(2).Info("get ServiceDNSRecord")
		instanceServiceRecord := &dnsv1alpha1.OpenMCPServiceDNSRecord{}
		err2 := r.live.Get(context.TODO(), req.NamespacedName, instanceServiceRecord)
		if err2 != nil && errors.IsNotFound(err2) {
			omcplog.V(0).Info("OpenMCPServiceDNSRecord Not Exist")
			return reconcile.Result{}, nil
		}
		omcplog.V(2).Info("[Get] OpenMCPServiceDNSRecord")

		// Check if a OpenMCPDomain exists
		instanceDomain := &dnsv1alpha1.OpenMCPDomain{}

		domainName := instanceServiceRecord.Spec.DomainRef
		domainNamespace := "kube-federation-system"
		nsn := types.NamespacedName{
			Namespace: domainNamespace,
			Name:      domainName,
		}
		err3 := r.live.Get(context.TODO(), nsn, instanceDomain)
		if err3 != nil && errors.IsNotFound(err3) {
			omcplog.V(0).Info("OpenMCPDomain Not Exist")
			return reconcile.Result{}, nil
		}
		omcplog.V(2).Info("[Get] OpenMCPDomain")

		// Status Update if OpenMCPServiceDNSRecord and OpenMCPDomain exist
		omcplog.V(2).Info("ServiceDNSRecord Status Update")
		serviceDNSRecord.FillStatus(instanceServiceRecord, instanceDomain)

		err4 := r.live.Status().Update(context.TODO(), instanceServiceRecord)
		if err4 != nil {
			omcplog.V(0).Info("[OpenMCP Service DNS Record Controller] : ", err4)
			return reconcile.Result{}, err4
		}

	} else if err != nil && errors.IsNotFound(err) {
		omcplog.V(0).Info("[Service Deleted]: ignore request")
		return reconcile.Result{}, nil
	} else {
		omcplog.V(0).Info("[Service Deleted]: ignore request")
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}
