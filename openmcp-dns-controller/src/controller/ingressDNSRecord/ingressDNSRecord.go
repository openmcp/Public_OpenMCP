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

package ingressDNSRecord

import (
	"context"
	"fmt"
	"net"
	"openmcp/openmcp/apis"
	dnsv1alpha1 "openmcp/openmcp/apis/dns/v1alpha1"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"strings"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info(">>> DNSEndpoint NewController()")
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

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &dnsv1alpha1.OpenMCPIngressDNSRecord{}, controller.WatchOptions{}); err != nil {
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
	omcplog.V(5).Info("********* [ OpenMCP IngressDNSRecord", i, "] *********")
	omcplog.V(5).Info(req.Context, " / ", req.Namespace, " / ", req.Name)

	// Return for OpenMCPIngressDNSRecord deletion request
	instanceIngressRecord := &dnsv1alpha1.OpenMCPIngressDNSRecord{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instanceIngressRecord)
	omcplog.V(2).Info("[Get] OpenMCPIngressDNSRecord")
	if err == nil {

		// Get Ingress Domains
		instanceOpenMCPIngress := &resourcev1alpha1.OpenMCPIngress{}
		err = r.live.Get(context.TODO(), req.NamespacedName, instanceOpenMCPIngress)
		omcplog.V(2).Info("[Get] OpenMCPIngress")
		if err != nil {
			omcplog.V(0).Info("[OpenMCP DNS Endpoint Controller] : ", err)
		}

		domains := []string{}
		for _, rule := range instanceOpenMCPIngress.Spec.Template.Spec.Rules {
			domains = append(domains, rule.Host)
		}

		nsn := types.NamespacedName{
			Namespace: req.Namespace,
			Name:      "ingress-" + req.Name,
		}
		instanceEndpoint := &dnsv1alpha1.OpenMCPDNSEndpoint{}
		err := r.live.Get(context.TODO(), nsn, instanceEndpoint)
		omcplog.V(2).Info("[Get] OpenMCPDNSEndpoint")

		if err == nil {
			// Already Exist -> Update
			omcplog.V(2).Info("Try OpenMCPDNSEndpoint Update From OpenMCPIngressDNS")
			instanceEndpoint := OpenMCPEndpointUpdateObjectFromIngressDNS(instanceEndpoint, instanceIngressRecord, req.Namespace, req.Name, domains)
			err = r.live.Update(context.TODO(), instanceEndpoint)
			if err != nil {
				omcplog.V(0).Info("[OpenMCP DNS Endpoint Controller] : ", err)

			}
		} else if errors.IsNotFound(err) {
			// Not Exist -> Create
			omcplog.V(2).Info("Try OpenMCPDNSEndpoint Create From OpenMCPIngressDNS")
			instanceEndpoint := OpenMCPEndpointCreateObjectFromIngressDNS(instanceIngressRecord, req.Namespace, req.Name, domains)
			err = r.live.Create(context.TODO(), instanceEndpoint)
			if err != nil {
				omcplog.V(0).Info("[OpenMCP DNS Endpoint Controller] : ", err)

			}

		} else {
			// Error !
			omcplog.V(0).Info("[OpenMCP DNS Endpoint Controller] : ", err)

		}
	} else if errors.IsNotFound(err) {
		// OpenMCPIngressDNSRecord Deleted -> Delete
		omcplog.V(2).Info("Try OpenMCPDNSEndpoint Delete From OpenMCPIngressDNS")
		instanceEndpoint := OpenMCPEndpointDeleteObjectFromIngressDNS(req.Namespace, req.Name)
		err := r.live.Delete(context.TODO(), instanceEndpoint)
		if err == nil {
			omcplog.V(0).Info("[OpenMCP DNS Endpoint Controller] : Deleted '", req.Name+"'")

		}

	}

	return reconcile.Result{}, nil // err
}
func CreateEndpoint(dnsName string, recordTTL dnsv1alpha1.TTL, recordType string, targets []string) *dnsv1alpha1.Endpoint {
	endpoint := &dnsv1alpha1.Endpoint{
		DNSName:    strings.ToLower(dnsName),
		Targets:    targets,
		RecordType: recordType,
		RecordTTL:  recordTTL,
		Labels:     nil,
	}
	return endpoint
}

func CreateEndpointsFromIngressDNS(instanceIngressRecord *dnsv1alpha1.OpenMCPIngressDNSRecord, namespace, name string) []*dnsv1alpha1.Endpoint {
	endpoints := []*dnsv1alpha1.Endpoint{}

	recordTTL := instanceIngressRecord.Spec.RecordTTL
	recordType := "A"

	for _, dns := range instanceIngressRecord.Status.DNS {
		for _, host := range dns.Hosts {
			dnsName := host
			targets := []string{}
			for _, ingress := range dns.LoadBalancer.Ingress {
				target := ingress.IP
				if ingress.IP == "" {
					addrs, err := net.LookupIP(ingress.Hostname)
					if err != nil {
						fmt.Println("Unknown host: ", ingress.Hostname)
					} else {
						for _, addr := range addrs {
							target := addr.String()
							fmt.Println("CHECK : ", target)
							targets = append(targets, target)
						}
					}
				} else {
					targets = append(targets, target)
				}

			}
			endpoint := CreateEndpoint(dnsName, recordTTL, recordType, targets)
			omcplog.V(3).Info("DNSName : ", endpoint.DNSName)
			omcplog.V(3).Info("RecordTTL : ", endpoint.RecordTTL)
			omcplog.V(3).Info("RecordType : ", endpoint.RecordType)
			omcplog.V(3).Info("Targets : ", endpoint.Targets)
			endpoints = append(endpoints, endpoint)
		}

	}

	return endpoints

}
func OpenMCPEndpointCreateObjectFromIngressDNS(instanceIngressRecord *dnsv1alpha1.OpenMCPIngressDNSRecord, namespace, name string, domains []string) *dnsv1alpha1.OpenMCPDNSEndpoint {

	endpoints := CreateEndpointsFromIngressDNS(instanceIngressRecord, namespace, name)

	instanceEndpoint := &dnsv1alpha1.OpenMCPDNSEndpoint{
		ObjectMeta: v1.ObjectMeta{
			Name:      "ingress-" + name,
			Namespace: namespace,
		},
		Spec: dnsv1alpha1.OpenMCPDNSEndpointSpec{
			Endpoints: endpoints,
			Domains:   domains,
		},
		Status: dnsv1alpha1.OpenMCPDNSEndpointStatus{},
	}
	return instanceEndpoint
}
func OpenMCPEndpointUpdateObjectFromIngressDNS(instanceEndpoint *dnsv1alpha1.OpenMCPDNSEndpoint, instanceIngressRecord *dnsv1alpha1.OpenMCPIngressDNSRecord, namespace, name string, domains []string) *dnsv1alpha1.OpenMCPDNSEndpoint {

	endpoints := CreateEndpointsFromIngressDNS(instanceIngressRecord, namespace, name)
	instanceEndpoint.Spec.Endpoints = endpoints
	instanceEndpoint.Spec.Domains = domains
	return instanceEndpoint
}
func OpenMCPEndpointDeleteObjectFromIngressDNS(namespace, name string) *dnsv1alpha1.OpenMCPDNSEndpoint {
	instanceEndpoint := &dnsv1alpha1.OpenMCPDNSEndpoint{
		ObjectMeta: v1.ObjectMeta{
			Name:      "ingress-" + name,
			Namespace: namespace,
		},
		Spec:   dnsv1alpha1.OpenMCPDNSEndpointSpec{},
		Status: dnsv1alpha1.OpenMCPDNSEndpointStatus{},
	}
	return instanceEndpoint

}
func ClearStatus(instanceIngressRecord *dnsv1alpha1.OpenMCPIngressDNSRecord) {
	instanceIngressRecord.Status = dnsv1alpha1.OpenMCPIngressDNSRecordStatus{}
}

func FillStatus(instanceIngressRecord *dnsv1alpha1.OpenMCPIngressDNSRecord) error {
	instanceIngressRecord.Status = dnsv1alpha1.OpenMCPIngressDNSRecordStatus{}

	lb := corev1.LoadBalancerStatus{}
	instanceIngress := &extv1b1.Ingress{}
	err := cm.Host_client.Get(context.TODO(), instanceIngress, instanceIngressRecord.Namespace, instanceIngressRecord.Name)
	if err == nil {
		// Get lb information if service exists
		lb = instanceIngress.Status.LoadBalancer
	}

	hosts := []string{}
	for i := 0; i < len(instanceIngress.Spec.Rules); i++ {
		hosts = append(hosts, instanceIngress.Spec.Rules[i].Host)
	}
	clusterIngressDNS := &dnsv1alpha1.ClusterIngressDNS{
		Cluster:      "openmcp",
		LoadBalancer: lb,
		Hosts:        hosts,
	}
	omcplog.V(3).Info("Cluster : ", clusterIngressDNS.Cluster)
	omcplog.V(3).Info("LoadBalancer Ingress : ", clusterIngressDNS.LoadBalancer.Ingress)
	omcplog.V(3).Info("Hosts : ", clusterIngressDNS.Hosts)
	instanceIngressRecord.Status.DNS = append(instanceIngressRecord.Status.DNS, *clusterIngressDNS)

	// Record Ingress information for each cluster
	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_genClients[cluster.Name]

		lb := corev1.LoadBalancerStatus{}
		instanceIngress := &extv1b1.Ingress{}
		err := cluster_client.Get(context.TODO(), instanceIngress, instanceIngressRecord.Namespace, instanceIngressRecord.Name)
		if err == nil {
			// Get lb information if service exists
			lb = instanceIngress.Status.LoadBalancer
		}

		hosts := []string{}
		for i := 0; i < len(instanceIngress.Spec.Rules); i++ {
			hosts = append(hosts, instanceIngress.Spec.Rules[i].Host)
		}
		clusterIngressDNS := &dnsv1alpha1.ClusterIngressDNS{
			Cluster:      cluster.Name,
			LoadBalancer: lb,
			Hosts:        hosts,
		}
		omcplog.V(3).Info("Cluster : ", clusterIngressDNS.Cluster)
		omcplog.V(3).Info("LoadBalancer Ingress : ", clusterIngressDNS.LoadBalancer.Ingress)
		omcplog.V(3).Info("Hosts : ", clusterIngressDNS.Hosts)
		instanceIngressRecord.Status.DNS = append(instanceIngressRecord.Status.DNS, *clusterIngressDNS)
	}

	return nil

}
