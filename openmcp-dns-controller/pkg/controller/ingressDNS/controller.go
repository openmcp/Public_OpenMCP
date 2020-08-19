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

package ingressDNS

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-dns-controller/pkg/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-dns-controller/pkg/apis/keti/v1alpha1"
	"openmcp/openmcp/util/clusterManager"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info( ">>> IngressDNS NewController()")
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

	// fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPIngressDNSRecord{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	if err := co.WatchResourceReconcileObject(live, &extv1b1.Ingress{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	for _, ghost := range ghosts {
		if err := co.WatchResourceReconcileController(ghost, &extv1b1.Ingress{}, controller.WatchOptions{}); err != nil {
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
	omcplog.V(4).Info( "Function Called Reconcile")
	i += 1
	omcplog.V(5).Info( "********* [ OpenMCP Domain", i, "] *********")
	omcplog.V(5).Info( req.Context, " / ", req.Namespace, " / ", req.Name)
	//cm := clusterManager.NewClusterManager()

	// Fetch the Sync instance
	instanceIngressRecord := &ketiv1alpha1.OpenMCPIngressDNSRecord{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instanceIngressRecord)
	if err != nil {
		// Delete
		omcplog.V(2).Info("IngressDNSRecord 삭제 감지")
		return reconcile.Result{}, nil
	}
	omcplog.V(2).Info("IngressDNSRecord or Ingress 생성 감지")



	omcplog.V(2).Info("IngressDNSRecord Status 업데이트")
	FillStatus(instanceIngressRecord)
	err = r.live.Status().Update(context.TODO(), instanceIngressRecord)
	if err != nil {
		omcplog.V(0).Info( "[OpenMCP Ingress DNS Record Controller] : ",err)
		return reconcile.Result{}, nil
	}



	return reconcile.Result{}, nil // err
}

func ClearStatus(instanceIngressRecord *ketiv1alpha1.OpenMCPIngressDNSRecord) {
	instanceIngressRecord.Status = ketiv1alpha1.OpenMCPIngressDNSRecordStatus{}
}

func FillStatus(instanceIngressRecord *ketiv1alpha1.OpenMCPIngressDNSRecord) error {
	instanceIngressRecord.Status = ketiv1alpha1.OpenMCPIngressDNSRecordStatus{}

	// OpenMCP Ingress도 기록
	lb :=  corev1.LoadBalancerStatus{}
	instanceIngress := &extv1b1.Ingress{}
	err := cm.Host_client.Get(context.TODO(), instanceIngress,  instanceIngressRecord.Namespace,  instanceIngressRecord.Name)
	if err == nil {
		// 서비스가 존재하면 lb 정보 가져옴
		lb = instanceIngress.Status.LoadBalancer
	}

	hosts := []string{}
	for i := 0; i <len(instanceIngress.Spec.Rules); i++ {
		hosts = append(hosts, instanceIngress.Spec.Rules[i].Host)
	}
	clusterIngressDNS := &ketiv1alpha1.ClusterIngressDNS{
		Cluster:      "openmcp",
		LoadBalancer: lb,
		Hosts: hosts,
	}
	omcplog.V(3).Info("Cluster : ", clusterIngressDNS.Cluster)
	omcplog.V(3).Info("LoadBalancer Ingress : ", clusterIngressDNS.LoadBalancer.Ingress)
	omcplog.V(3).Info("Hosts : ", clusterIngressDNS.Hosts)
	instanceIngressRecord.Status.DNS = append(instanceIngressRecord.Status.DNS, *clusterIngressDNS)

	// 각 클러스터의 Ingress 정보 기록
	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_genClients[cluster.Name]

		lb :=  corev1.LoadBalancerStatus{}
		instanceIngress := &extv1b1.Ingress{}
		err := cluster_client.Get(context.TODO(), instanceIngress,  instanceIngressRecord.Namespace,  instanceIngressRecord.Name)
		if err == nil {
			// 서비스가 존재하면 lb 정보 가져옴
			lb = instanceIngress.Status.LoadBalancer
		}

		hosts := []string{}
		for i := 0; i <len(instanceIngress.Spec.Rules); i++ {
			hosts = append(hosts, instanceIngress.Spec.Rules[i].Host)
		}
		clusterIngressDNS := &ketiv1alpha1.ClusterIngressDNS{
			Cluster:      cluster.Name,
			LoadBalancer: lb,
			Hosts: hosts,
		}
		omcplog.V(3).Info("Cluster : ", clusterIngressDNS.Cluster)
		omcplog.V(3).Info("LoadBalancer Ingress : ", clusterIngressDNS.LoadBalancer.Ingress)
		omcplog.V(3).Info("Hosts : ", clusterIngressDNS.Hosts)
		instanceIngressRecord.Status.DNS = append(instanceIngressRecord.Status.DNS, *clusterIngressDNS)
	}


	return nil


}