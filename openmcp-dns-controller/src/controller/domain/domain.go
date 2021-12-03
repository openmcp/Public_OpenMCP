package domain

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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info(">>> Domain NewController()")
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

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &dnsv1alpha1.OpenMCPDomain{}, controller.WatchOptions{}); err != nil {
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

func (r *reconciler) UpdateStatusServiceDNSRecordFromDelete() error {
	// Update Blank Status to delete Service DNS Record information when OpenMCPDomain Delete

	instanceServiceRecordList := &dnsv1alpha1.OpenMCPServiceDNSRecordList{}
	err := r.live.List(context.TODO(), instanceServiceRecordList, &client.ListOptions{})
	omcplog.V(2).Info("[List] OpenMCPServiceDNSRecordList")
	if err != nil {
		return err
	}
	instanceDomainList := &dnsv1alpha1.OpenMCPDomainList{}
	err = r.live.List(context.TODO(), instanceDomainList, &client.ListOptions{})
	omcplog.V(2).Info("[List] OpenMCPDomainList")
	if err != nil {
		return err
	}

	deleted_index := -1
	for i, instanceServiceRecord := range instanceServiceRecordList.Items {
		find := false
		deleted_index = i
		for _, inDomain := range instanceDomainList.Items {
			if instanceServiceRecord.Status.Domain == inDomain.Domain {
				find = true
				break
			}
		}

		if !find {
			omcplog.V(2).Info("[OpenMCP Domain Controller] Service DNS Record Delete :", instanceServiceRecordList.Items[deleted_index].Name)

			omcplog.V(2).Info("ServiceRecord Clear Status")
			serviceDNSRecord.ClearStatus(&instanceServiceRecordList.Items[deleted_index])
			err = r.live.Status().Update(context.TODO(), &instanceServiceRecordList.Items[deleted_index])
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
func (r *reconciler) UpdateStatusServiceDNSRecordFromCreate(instanceDomain *dnsv1alpha1.OpenMCPDomain) error {
	//  OpenMCPServiceDNSRecord update when OpenMCPDomain Create
	instanceServiceRecordList := &dnsv1alpha1.OpenMCPServiceDNSRecordList{}
	err := r.live.List(context.TODO(), instanceServiceRecordList, &client.ListOptions{})
	omcplog.V(2).Info("[List] OpenMCPServiceDNSRecordList")
	if err != nil {
		return err
	}

	for _, instanceServiceRecord := range instanceServiceRecordList.Items {
		if instanceServiceRecord.Spec.DomainRef == instanceDomain.Name {

			omcplog.V(2).Info("ServiceRecord Fill Status")
			serviceDNSRecord.FillStatus(&instanceServiceRecord, instanceDomain)

			err = r.live.Status().Update(context.TODO(), &instanceServiceRecord)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(4).Info("Function Called Reconcile")
	i += 1
	omcplog.V(5).Info("********* [ OpenMCP Domain", i, "] *********")
	omcplog.V(5).Info(req.Context, " / ", req.Namespace, " / ", req.Name)

	instanceDomain := &dnsv1alpha1.OpenMCPDomain{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instanceDomain)
	omcplog.V(2).Info("[List] OpenMCPDomain")
	if err != nil {
		omcplog.V(2).Info("Domain Delete Detection")
		err = r.UpdateStatusServiceDNSRecordFromDelete()
		if err != nil {
			omcplog.V(0).Info("[OpenMCP Domain Controller] : ", err)
		}

		return reconcile.Result{}, nil
	}
	omcplog.V(2).Info("Domain Create Detection")
	omcplog.V(2).Info("ServiceDNSRecord Status Update")
	err = r.UpdateStatusServiceDNSRecordFromCreate(instanceDomain)
	if err != nil {
		omcplog.V(0).Info("[OpenMCP Domain Controller] : ", err)
	}

	return reconcile.Result{}, nil
}
