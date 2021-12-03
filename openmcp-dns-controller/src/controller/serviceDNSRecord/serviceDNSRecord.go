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

package serviceDNSRecord

import (
	"context"
	"fmt"
	"net"
	"openmcp/openmcp/apis"
	dnsv1alpha1 "openmcp/openmcp/apis/dns/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"strings"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
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

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &dnsv1alpha1.OpenMCPServiceDNSRecord{}, controller.WatchOptions{}); err != nil {
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
	omcplog.V(5).Info("********* [ OpenMCP ServiceDNSRecord", i, "] *********")
	omcplog.V(5).Info(req.Context, " / ", req.Namespace, " / ", req.Name)

	// Return for OpenMCPServiceDNSRecord deletion request
	instanceServiceRecord := &dnsv1alpha1.OpenMCPServiceDNSRecord{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instanceServiceRecord)
	omcplog.V(2).Info("[Get] OpenMCPServiceDNSRecord")

	if err == nil {

		if len(instanceServiceRecord.Status.DNS) == 0 {

			instanceDomain := &dnsv1alpha1.OpenMCPDomain{}
			domainName := instanceServiceRecord.Spec.DomainRef
			domainNamespace := "kube-federation-system"
			nsn := types.NamespacedName{
				Namespace: domainNamespace,
				Name:      domainName,
			}
			err2 := r.live.Get(context.TODO(), nsn, instanceDomain)
			if err2 != nil && errors.IsNotFound(err2) {
				return reconcile.Result{}, nil
			}

			FillStatus(instanceServiceRecord, instanceDomain)

			err3 := r.live.Status().Update(context.TODO(), instanceServiceRecord)
			if err3 != nil {
				return reconcile.Result{}, nil
			}

		}

		nsn := types.NamespacedName{
			Namespace: req.Namespace,
			Name:      "service-" + req.Name,
		}
		instanceEndpoint := &dnsv1alpha1.OpenMCPDNSEndpoint{}
		err := r.live.Get(context.TODO(), nsn, instanceEndpoint)
		omcplog.V(2).Info("[Get] OpenMCPDNSEndpoint")

		if err == nil {
			// Already Exist -> Update
			omcplog.V(2).Info("Try OpenMCPDNSEndpoint Update From OpenMCPServiceDNS")
			instanceEndpoint := OpenMCPEndpointUpdateObjectFromServiceDNS(instanceEndpoint, instanceServiceRecord, req.Namespace, req.Name)
			err = r.live.Update(context.TODO(), instanceEndpoint)
			if err != nil {
				omcplog.V(0).Info("[OpenMCP DNS Endpoint Controller] : ", err)

			}
		} else if errors.IsNotFound(err) {
			// Not Exist -> Create
			omcplog.V(2).Info("Try OpenMCPDNSEndpoint Create From OpenMCPServiceDNS")
			instanceEndpoint := OpenMCPEndpointCreateObjectFromServiceDNS(instanceServiceRecord, req.Namespace, req.Name)
			err = r.live.Create(context.TODO(), instanceEndpoint)
			if err != nil {
				omcplog.V(0).Info("[OpenMCP DNS Endpoint Controller] : ", err)

			}

		} else {
			// Error !
			omcplog.V(0).Info("[OpenMCP DNS Endpoint Controller] : ", err)

		}
	} else if errors.IsNotFound(err) {
		// OpenMCPServiceDNSRecord Deleted -> Delete
		omcplog.V(2).Info("Try OpenMCPDNSEndpoint Delete From OpenMCPServiceDNS")
		instanceEndpoint := OpenMCPEndpointDeleteObjectFromServiceDNS(req.Namespace, req.Name)
		err := r.live.Delete(context.TODO(), instanceEndpoint)
		if err == nil {
			omcplog.V(2).Info("[OpenMCP DNS Endpoint Controller] : Deleted '", req.Name+"'")

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
func CreateEndpointsFromServiceDNS(instanceServiceRecord *dnsv1alpha1.OpenMCPServiceDNSRecord, namespace, name string) []*dnsv1alpha1.Endpoint {
	endpoints := []*dnsv1alpha1.Endpoint{}

	domainRef := instanceServiceRecord.Spec.DomainRef
	recordTTL := instanceServiceRecord.Spec.RecordTTL
	domain := instanceServiceRecord.Status.Domain
	recordType := "A"

	targetsAll := []string{}
	for _, dns := range instanceServiceRecord.Status.DNS {
		targets := []string{}
		for _, ingress := range dns.LoadBalancer.Ingress {

			// If there is no ip (EKS uses domain instead of ip)
			// Check and insert the ip corresponding to domain.
			if ingress.IP == "" {

				addrs, err := net.LookupIP(ingress.Hostname)
				if err != nil {
					fmt.Println("Unknown host: ", ingress.Hostname)
					continue
				} else {
					fmt.Println(addrs)
					for _, addr := range addrs {
						target := addr.String()
						targets = append(targets, target)
						targetsAll = append(targetsAll, target)
					}
				}
			} else { // use ip as it is if there is an address ip (not EKS)
				target := ingress.IP
				targets = append(targets, target)
				targetsAll = append(targetsAll, target)
			}

		}

		if dns.Region == "" || dns.Zones == nil {
			continue
		}
		region := dns.Region
		dnsName := name + "." + namespace + "." + domainRef + ".svc." + region + "." + domain

		// DNS where only Region exists
		endpoint := CreateEndpoint(dnsName, recordTTL, recordType, targets)
		endpoints = append(endpoints, endpoint)

		for _, zone := range dns.Zones {
			dnsName := name + "." + namespace + "." + domainRef + ".svc." + zone + "." + region + "." + domain

			// DNS where both Region and Zone exist
			endpoint := CreateEndpoint(dnsName, recordTTL, recordType, targets)
			omcplog.V(3).Info("DNSName : ", endpoint.DNSName)
			omcplog.V(3).Info("RecordTTL : ", endpoint.RecordTTL)
			omcplog.V(3).Info("RecordType : ", endpoint.RecordType)
			omcplog.V(3).Info("Targets : ", endpoint.Targets)
			endpoints = append(endpoints, endpoint)
		}
	}

	// DNS that neither Region exists

	if domain != "" {
		dnsName := name + "." + namespace + "." + domainRef + ".svc." + domain
		endpoint := CreateEndpoint(dnsName, recordTTL, recordType, targetsAll)
		omcplog.V(3).Info("DNSName : ", endpoint.DNSName)
		omcplog.V(3).Info("RecordTTL : ", endpoint.RecordTTL)
		omcplog.V(3).Info("RecordType : ", endpoint.RecordType)
		omcplog.V(3).Info("Targets : ", endpoint.Targets)

		endpoints = append(endpoints, endpoint)
	}

	return endpoints

}
func OpenMCPEndpointCreateObjectFromServiceDNS(instanceServiceRecord *dnsv1alpha1.OpenMCPServiceDNSRecord, namespace, name string) *dnsv1alpha1.OpenMCPDNSEndpoint {

	endpoints := CreateEndpointsFromServiceDNS(instanceServiceRecord, namespace, name)

	instanceEndpoint := &dnsv1alpha1.OpenMCPDNSEndpoint{
		ObjectMeta: v1.ObjectMeta{
			Name:      "service-" + name,
			Namespace: namespace,
		},
		Spec: dnsv1alpha1.OpenMCPDNSEndpointSpec{
			Endpoints: endpoints,
			Domains:   []string{instanceServiceRecord.Status.Domain},
		},
		Status: dnsv1alpha1.OpenMCPDNSEndpointStatus{},
	}
	return instanceEndpoint
}
func OpenMCPEndpointUpdateObjectFromServiceDNS(instanceEndpoint *dnsv1alpha1.OpenMCPDNSEndpoint, instanceServiceRecord *dnsv1alpha1.OpenMCPServiceDNSRecord, namespace, name string) *dnsv1alpha1.OpenMCPDNSEndpoint {

	endpoints := CreateEndpointsFromServiceDNS(instanceServiceRecord, namespace, name)
	instanceEndpoint.Spec.Endpoints = endpoints
	instanceEndpoint.Spec.Domains = []string{instanceServiceRecord.Status.Domain}
	return instanceEndpoint
}
func OpenMCPEndpointDeleteObjectFromServiceDNS(namespace, name string) *dnsv1alpha1.OpenMCPDNSEndpoint {
	instanceEndpoint := &dnsv1alpha1.OpenMCPDNSEndpoint{
		ObjectMeta: v1.ObjectMeta{
			Name:      "service-" + name,
			Namespace: namespace,
		},
		Spec:   dnsv1alpha1.OpenMCPDNSEndpointSpec{},
		Status: dnsv1alpha1.OpenMCPDNSEndpointStatus{},
	}
	return instanceEndpoint

}
func ClearStatus(instanceServiceRecord *dnsv1alpha1.OpenMCPServiceDNSRecord) {
	instanceServiceRecord.Status = dnsv1alpha1.OpenMCPServiceDNSRecordStatus{}
}

func FillStatus(instanceServiceRecord *dnsv1alpha1.OpenMCPServiceDNSRecord, instanceDomain *dnsv1alpha1.OpenMCPDomain) error {

	instanceServiceRecord.Status = dnsv1alpha1.OpenMCPServiceDNSRecordStatus{}

	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_genClients[cluster.Name]

		// Node Info of Cluster (Zone, Region)
		instanceNodeList := &corev1.NodeList{}
		err := cluster_client.List(context.TODO(), instanceNodeList, "default")
		if err != nil {
			omcplog.V(0).Info("[OpenMCP Service DNS Record Controller] : ", err)
			return nil
		}

		region := ""
		if len(instanceNodeList.Items) >= 1 {
			if val, ok := instanceNodeList.Items[0].Labels["topology.kubernetes.io/region"]; ok {
				region = val
			} else if val, ok := instanceNodeList.Items[0].Labels["failure-domain.beta.kubernetes.io/region"]; ok {
				region = val
			}
		}

		zones := []string{}
		zones_dup_map := make(map[string]string) // Map for deduplication

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

		// Node Info of Cluster (Zone, Region)
		lb := corev1.LoadBalancerStatus{}
		instanceService := &corev1.Service{}

		// if default domain use
		// OpenmcpService made default dns with SERVICENAME +"-by-openmcp"
		name := instanceServiceRecord.Name
		last11 := name[len(name)-11:]

		if last11 == "-by-openmcp" {
			name = name[:len(name)-11]

		}
		err = cluster_client.Get(context.TODO(), instanceService, instanceServiceRecord.Namespace, name)
		if err == nil {
			// Get lb information if service exists
			lb = instanceService.Status.LoadBalancer

		}
		clusterDNS := &dnsv1alpha1.ClusterDNS{
			Cluster:      cluster.Name,
			LoadBalancer: lb, // instanceService.Status.LoadBalancer,
			Zones:        zones,
			Region:       region,
		}

		omcplog.V(3).Info("Cluster : ", clusterDNS.Cluster)
		omcplog.V(3).Info("LoadBalancer Ingress : ", clusterDNS.LoadBalancer.Ingress)
		omcplog.V(3).Info("Zones : ", clusterDNS.Zones)
		omcplog.V(3).Info("Region : ", clusterDNS.Region)

		instanceServiceRecord.Status.DNS = append(instanceServiceRecord.Status.DNS, *clusterDNS)

	}
	instanceServiceRecord.Status.Domain = instanceDomain.Domain
	return nil
}
func FillStatusInOpenMCP(instanceServiceRecord *dnsv1alpha1.OpenMCPServiceDNSRecord, instanceDomain *dnsv1alpha1.OpenMCPDomain) error {

	instanceServiceRecord.Status = dnsv1alpha1.OpenMCPServiceDNSRecordStatus{}

	host_client := cm.Host_client

	// Node Info of Cluster (Zone, Region)
	instanceNodeList := &corev1.NodeList{}
	err := host_client.List(context.TODO(), instanceNodeList, "default")
	if err != nil {
		omcplog.V(0).Info("[OpenMCP Service DNS Record Controller] : ", err)
		return nil
	}

	region := ""
	if len(instanceNodeList.Items) >= 1 {
		if val, ok := instanceNodeList.Items[0].Labels["topology.kubernetes.io/region"]; ok {
			region = val
		} else if val, ok := instanceNodeList.Items[0].Labels["failure-domain.beta.kubernetes.io/region"]; ok {
			region = val
		}
	}

	zones := []string{}
	zones_dup_map := make(map[string]string) // Map for deduplication

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

	// Node Info of Cluster (Zone, Region)
	lb := corev1.LoadBalancerStatus{}
	instanceService := &corev1.Service{}
	err = host_client.Get(context.TODO(), instanceService, instanceServiceRecord.Namespace, instanceServiceRecord.Name)
	if err == nil {
		// Get lb information if service exists
		lb = instanceService.Status.LoadBalancer

	}
	clusterDNS := &dnsv1alpha1.ClusterDNS{
		Cluster:      "openmcp",
		LoadBalancer: lb, // instanceService.Status.LoadBalancer,
		Zones:        zones,
		Region:       region,
	}

	omcplog.V(3).Info("Cluster : ", clusterDNS.Cluster)
	omcplog.V(3).Info("LoadBalancer Ingress : ", clusterDNS.LoadBalancer.Ingress)
	omcplog.V(3).Info("Zones : ", clusterDNS.Zones)
	omcplog.V(3).Info("Region : ", clusterDNS.Region)

	instanceServiceRecord.Status.DNS = append(instanceServiceRecord.Status.DNS, *clusterDNS)
	instanceServiceRecord.Status.Domain = instanceDomain.Domain
	return nil
}
