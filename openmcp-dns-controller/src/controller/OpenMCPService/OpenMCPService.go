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

package OpenMCPService // import "admiralty.io/multicluster-controller/examples/serviceDNS/pkg/controller/serviceDNS"

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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	if err := co.WatchResourceReconcileObject(context.TODO(), live, &resourcev1alpha1.OpenMCPService{}, controller.WatchOptions{}); err != nil {
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
	omcplog.V(4).Info("Function Called Reconcile")
	i += 1
	omcplog.V(5).Info("********* [ OpenMCPService Controller", i, "] *********")
	omcplog.V(5).Info(req.Context, " / ", req.Namespace, " / ", req.Name)

	omcplog.V(2).Info("OpenMCPService Request")
	osvc := &resourcev1alpha1.OpenMCPService{}
	err := r.live.Get(context.TODO(), req.NamespacedName, osvc)
	if err != nil && errors.IsNotFound(err) {
		omcplog.V(2).Info("Deleted OpenMCPService. Delete ServiceDNSRecord")
		instance_osvcdnsr := &dnsv1alpha1.OpenMCPServiceDNSRecord{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name + "-by-openmcp",
				Namespace: req.Namespace,
			},
		}
		type ObjectKey = types.NamespacedName
		osdnsr := &dnsv1alpha1.OpenMCPServiceDNSRecord{}
		err_osdnsr := r.live.Get(context.TODO(), ObjectKey{Name: req.Name + "-by-openmcp", Namespace: req.Namespace}, osdnsr)
		if err_osdnsr != nil && errors.IsNotFound(err_osdnsr) {
			err2 := r.live.Delete(context.TODO(), instance_osvcdnsr)
			if err2 != nil && errors.IsNotFound(err2) {
				omcplog.V(2).Info(err2)
				return reconcile.Result{}, nil
			} else if err2 != nil {
				omcplog.V(2).Info("err : ", err2)
				return reconcile.Result{}, err2
			}
		}else if err_osdnsr != nil {
			omcplog.V(2).Info(err_osdnsr)
		}
		return reconcile.Result{}, nil
	}

	if err == nil {
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
				Name:      req.Name + "-by-openmcp",
				Namespace: req.Namespace,
			},
			Spec: dnsv1alpha1.OpenMCPServiceDNSRecordSpec{
				DomainRef: "default-domain",
				RecordTTL: 300,
			},
		}
		nsn := types.NamespacedName{
			Namespace: req.Namespace,
			Name:      req.Name + "-by-openmcp",
		}
		err2 := r.live.Get(context.TODO(), nsn, instance_osvcdnsr)
		if err2 != nil && errors.IsNotFound(err2) {
			omcplog.V(0).Info("[OpenMCPAPIServer] Create OpenMCPServiceDNSReocrd")
			err3 := r.live.Create(context.TODO(), instance_osvcdnsr)
			if err3 != nil && errors.IsAlreadyExists(err3) {
				omcplog.V(0).Info(err3)

			} else if err3 != nil {
				omcplog.V(0).Info(err3)
				return reconcile.Result{}, err3
			}
		}

		err4 := r.live.Get(context.TODO(), nsn, instance_osvcdnsr)
		if err4 != nil && errors.IsNotFound(err4) {
			omcplog.V(0).Info("err:", err4)
			return reconcile.Result{}, err4
		}
		serviceDNSRecord.FillStatus(instance_osvcdnsr, instance_default_domain)
		omcplog.V(0).Info("[OpenMCPAPIServer] Update OpenMCPServiceDNSReocrd")

		err5 := r.live.Status().Update(context.TODO(), instance_osvcdnsr)
		if err5 != nil {
			omcplog.V(0).Info(err5)
			return reconcile.Result{}, err5
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
