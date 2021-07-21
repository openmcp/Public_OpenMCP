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

package openmcpservice // import "admiralty.io/multicluster-controller/examples/serviceDNS/pkg/controller/serviceDNS"

import (
	"context"
	"fmt"
	"openmcp/openmcp/apis"
	dnsv1alpha1 "openmcp/openmcp/apis/dns/v1alpha1"
	ketiv1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
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

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &ketiv1alpha1.OpenMCPService{}, controller.WatchOptions{}); err != nil {
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
	omcplog.V(4).Info("Function Called Reconcile")
	i += 1
	omcplog.V(5).Info("********* [ OpenMCP OpenMCPService", i, "] *********")
	omcplog.V(5).Info(req.Context, " / ", req.Namespace, " / ", req.Name)

	// Return for OpenMCPService deletion request
	omcplog.V(2).Info("ServiceDNSRecord Request")
	osvc := &ketiv1alpha1.OpenMCPService{}
	err := r.live.Get(context.TODO(), req.NamespacedName, osvc)
	omcplog.V(2).Info("[Get] OpenMCPService")
	if err != nil {
		if errors.IsNotFound(err) {
			// Delete
			omcplog.V(0).Info("OpenMCPService has been deleted. Delete Default OpenMCPServiceDNSRecord")
			instance_osvcdnsr := &dnsv1alpha1.OpenMCPServiceDNSRecord{
				ObjectMeta: metav1.ObjectMeta{
					Name:      req.Name,
					Namespace: req.Namespace,
				},
				Spec: dnsv1alpha1.OpenMCPServiceDNSRecordSpec{
					DomainRef: "openmcp-default-domain",
					RecordTTL: 300,
				},
			}
			err = r.live.Delete(context.TODO(), instance_osvcdnsr)
			if err != nil && !errors.IsAlreadyExists(err) {
				omcplog.V(0).Info(err)
			}

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, err
	}
	omcplog.V(2).Info("OpenMCPService Create Detection")

	instance_default_domain := &dnsv1alpha1.OpenMCPDomain{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-default-domain",
			Namespace: "kube-federation-system",
		},
		Domain: "openmcp.example.org",
	}

	instance_osvcdnsr := &dnsv1alpha1.OpenMCPServiceDNSRecord{
		ObjectMeta: metav1.ObjectMeta{
			Name:      osvc.Name,
			Namespace: osvc.Namespace,
		},
		Spec: dnsv1alpha1.OpenMCPServiceDNSRecordSpec{
			DomainRef: "openmcp-default-domain",
			RecordTTL: 300,
		},
	}
	err = r.live.Create(context.TODO(), instance_default_domain)
	if err != nil && !errors.IsAlreadyExists(err) {
		omcplog.V(0).Info(err)
	}
	err = r.live.Create(context.TODO(), instance_osvcdnsr)
	if err != nil && !errors.IsAlreadyExists(err) {
		omcplog.V(0).Info(err)
	}

	/*

		// Check if a OpenMCPDomain exists
		instanceDomain := &dnsv1alpha1.OpenMCPDomain{}

		domainName := instanceServiceRecord.Spec.DomainRef
		domainNamespace := "kube-federation-system"
		nsn := types.NamespacedName{
			Namespace: domainNamespace,
			Name:      domainName,
		}
		err = r.live.Get(context.TODO(), nsn, instanceDomain)
		omcplog.V(2).Info("[Get] OpenMCPDomain")

		if err != nil {
			// OpenMCPDomain is not Exist
			omcplog.V(0).Info("OpenMCPDomain Not Exist")
			return reconcile.Result{}, nil
		}

		// Status Update if OpenMCPServiceDNSRecord and OpenMCPDomain exist
		omcplog.V(2).Info("ServiceDNSRecord Status Update")
		FillStatus(instanceServiceRecord, instanceDomain)

		err = r.live.Status().Update(context.TODO(), instanceServiceRecord)
		if err != nil {
			omcplog.V(0).Info("[OpenMCP Service DNS Record Controller] : ", err)
			return reconcile.Result{}, nil
		}
	*/

	return reconcile.Result{}, nil
}
