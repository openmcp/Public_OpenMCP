package domain

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"context"
	"fmt"
	"k8s.io/klog"
	"openmcp/openmcp/openmcp-dns-controller/pkg/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-dns-controller/pkg/apis/keti/v1alpha1"
	"openmcp/openmcp/openmcp-dns-controller/pkg/controller/serviceDNS"
	"openmcp/openmcp/util/clusterManager"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string) (*controller.Controller, error) {
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

	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPDomain{}, controller.WatchOptions{}); err != nil {
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
	// OpenMCPDomain Delete시 Service DNS Record 정보 삭제 하도록 빈 Status 업데이트

	instanceServiceRecordList := &ketiv1alpha1.OpenMCPServiceDNSRecordList{}
	err := r.live.List(context.TODO(), instanceServiceRecordList, &client.ListOptions{})
	if err != nil {
		return err
	}
	instanceDomainList := &ketiv1alpha1.OpenMCPDomainList{}
	err = r.live.List(context.TODO(), instanceDomainList, &client.ListOptions{})
	if err != nil {
		klog.V(0).Info( "[OpenMCP Domain Controller] : ", err)
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
			klog.V(0).Info( "[OpenMCP Domain Controller] Service DNS Record Delete :", instanceServiceRecordList.Items[deleted_index].Name)

			serviceDNS.ClearStatus(&instanceServiceRecordList.Items[deleted_index])
			err = r.live.Status().Update(context.TODO(), &instanceServiceRecordList.Items[deleted_index])
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
func (r *reconciler) UpdateStatusServiceDNSRecordFromCreate(instanceDomain *ketiv1alpha1.OpenMCPDomain) error {
	// OpenMCPDomain Create시 OpenMCPServiceDNSRecord 업데이트
	instanceServiceRecordList := &ketiv1alpha1.OpenMCPServiceDNSRecordList{}
	err := r.live.List(context.TODO(), instanceServiceRecordList, &client.ListOptions{})
	if err != nil {
		return err
	}

	for _, instanceServiceRecord := range instanceServiceRecordList.Items {
		if instanceServiceRecord.Spec.DomainRef == instanceDomain.Name {

			serviceDNS.FillStatus(&instanceServiceRecord, instanceDomain)

			err = r.live.Status().Update(context.TODO(), &instanceServiceRecord)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	i += 1
	//klog.V(0).Info( "********* [ OpenMCP Domain", i, "] *********")
	//klog.V(0).Info( req.Context, " / ", req.Namespace, " / ", req.Name)
	//cm := clusterManager.NewClusterManager()

	instanceDomain := &ketiv1alpha1.OpenMCPDomain{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instanceDomain)
	if err != nil {
		err = r.UpdateStatusServiceDNSRecordFromDelete()
		if err != nil {
			klog.V(0).Info( "[OpenMCP Domain Controller] : ", err)
		}

		return reconcile.Result{}, nil
	}

	err = r.UpdateStatusServiceDNSRecordFromCreate(instanceDomain)
	if err != nil {
		klog.V(0).Info( "[OpenMCP Domain Controller] : ", err)
	}

	return reconcile.Result{}, nil
}
