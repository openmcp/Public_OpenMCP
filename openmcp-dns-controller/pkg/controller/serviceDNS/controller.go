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

package serviceDNS // import "admiralty.io/multicluster-controller/examples/serviceDNS/pkg/controller/serviceDNS"

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"openmcp/openmcp/openmcp-dns-controller/pkg/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-dns-controller/pkg/apis/keti/v1alpha1"
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

	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPServiceDNSRecord{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	for _, ghost := range ghosts {
		if err := co.WatchResourceReconcileController(ghost, &corev1.Service{}, controller.WatchOptions{}); err != nil {
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
	i += 1
	//klog.V(0).Info("********* [ OpenMCP Service DNS Record", i, "] *********")
	//klog.V(0).Info(req.Context, " / ", req.Namespace, " / ", req.Name)
	//cm := clusterManager.NewClusterManager()


	// OpenMCPServiceDNSRecord 삭제 요청인 경우 종료
	instanceServiceRecord := &ketiv1alpha1.OpenMCPServiceDNSRecord{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instanceServiceRecord)
	if err != nil {
		// Delete
		return reconcile.Result{}, nil
	}


	// 도메인이 있는지 체크
	instanceDomain := &ketiv1alpha1.OpenMCPDomain{}

	domainName := instanceServiceRecord.Spec.DomainRef
	domainNamespace := "kube-federation-system"
	nsn := types.NamespacedName{
		Namespace: domainNamespace,
		Name:      domainName,
	}
	err = r.live.Get(context.TODO(), nsn, instanceDomain)

	if err != nil {
		// OpenMCPDomain이 없을경우
		return reconcile.Result{}, nil
	}

	// OpenMCPServiceDNSRecord과 OpenMCPDomain이 존재하는경우
	// Status 업데이트
	FillStatus(instanceServiceRecord, instanceDomain)

	err = r.live.Status().Update(context.TODO(), instanceServiceRecord)
	if err != nil {
		//klog.V(0).Info("[OpenMCP Service DNS Record Controller] : ",err)
		return reconcile.Result{}, nil
	}


	return reconcile.Result{}, nil
}
func ClearStatus(instanceServiceRecord *ketiv1alpha1.OpenMCPServiceDNSRecord) {
	instanceServiceRecord.Status = ketiv1alpha1.OpenMCPServiceDNSRecordStatus{}
}

func FillStatus(instanceServiceRecord *ketiv1alpha1.OpenMCPServiceDNSRecord, instanceDomain *ketiv1alpha1.OpenMCPDomain) error {
	
	instanceServiceRecord.Status = ketiv1alpha1.OpenMCPServiceDNSRecordStatus{}
	
	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_genClients[cluster.Name]

		// 클러스터의 노드정보 (Zone, Region)
		instanceNodeList := &corev1.NodeList{}
		err := cluster_client.List(context.TODO(), instanceNodeList, "default")
		if err != nil {
			klog.V(0).Info("[OpenMCP Service DNS Record Controller] : ",err)
			return nil
		}


		region := ""
		if val, ok := instanceNodeList.Items[0].Labels["topology.kubernetes.io/region"]; ok {
			region = val
		} else if val, ok := instanceNodeList.Items[0].Labels["failure-domain.beta.kubernetes.io/region"]; ok {
			region = val
		}


		zones := []string{}
		zones_dup_map := make(map[string]string) // 중복제거를위한 딕셔너리

		for _, node := range instanceNodeList.Items {
			if val, ok := node.Labels["topology.kubernetes.io/zone"]; ok {
				if _, ok := zones_dup_map[val]; ok {

				} else {
					zones = append(zones, val)
					zones_dup_map[val] = "1"
				}
			} else if val, ok := node.Labels["failure-domain.beta.kubernetes.io/zone"]; ok {
				if _, ok := zones_dup_map[val]; ok {

				} else {
					zones = append(zones, val)
					zones_dup_map[val] = "1"
				}
			}
		}

		// 클러스터의 노드정보 (Zone, Region)
		lb :=  corev1.LoadBalancerStatus{}
		instanceService := &corev1.Service{}
		err = cluster_client.Get(context.TODO(), instanceService,  instanceServiceRecord.Namespace,  instanceServiceRecord.Name)
		if err == nil {
			// 서비스가 존재하면 lb 정보 가져옴
			lb = instanceService.Status.LoadBalancer

		}
		clusterDNS := &ketiv1alpha1.ClusterDNS{
			Cluster:      cluster.Name,
			LoadBalancer: lb, // instanceService.Status.LoadBalancer,
			Zones:        zones,
			Region:       region,
		}

		instanceServiceRecord.Status.DNS = append(instanceServiceRecord.Status.DNS, *clusterDNS)

	}
	instanceServiceRecord.Status.Domain = instanceDomain.Domain
	return nil
}
