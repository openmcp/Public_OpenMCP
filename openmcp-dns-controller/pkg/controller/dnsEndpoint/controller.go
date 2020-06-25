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

package dnsEndpoint

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"openmcp/openmcp/openmcp-dns-controller/pkg/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-dns-controller/pkg/apis/keti/v1alpha1"
	resapis "openmcp/openmcp/openmcp-resource-controller/apis"
	resketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
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
	if err := resapis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	// fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPServiceDNSRecord{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPIngressDNSRecord{}, controller.WatchOptions{}); err != nil {
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
	i += 1
	//fmt.Println("********* [ OpenMCP DNS Endpoint", i, "] *********")
	//fmt.Println(req.Context, " / ", req.Namespace, " / ", req.Name)

	// OpenMCPServiceDNSRecord 삭제 요청인 경우 종료
	instanceServiceRecord := &ketiv1alpha1.OpenMCPServiceDNSRecord{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instanceServiceRecord)
	if err == nil {
		//if instanceServiceRecord.Status.DNS == nil {
		//	return reconcile.Result{}, nil
		//}

		nsn := types.NamespacedName{
			Namespace: req.Namespace,
			Name:      "service-"+req.Name,
		}
		instanceEndpoint := &ketiv1alpha1.OpenMCPDNSEndpoint{}
		err := r.live.Get(context.TODO(), nsn, instanceEndpoint)


		if err == nil {
			// Already Exist -> Update
			instanceEndpoint := OpenMCPEndpointUpdateObjectFromServiceDNS(instanceEndpoint, instanceServiceRecord, req.Namespace,  req.Name)
			err = r.live.Update(context.TODO(), instanceEndpoint)
			if err != nil {
				fmt.Println("[OpenMCP DNS Endpoint Controller] : ",err)
				//return reconcile.Result{}, nil
			}
		} else if errors.IsNotFound(err) {
			// Not Exist -> Create
			instanceEndpoint := OpenMCPEndpointCreateObjectFromServiceDNS(instanceServiceRecord, req.Namespace,  req.Name)
			err = r.live.Create(context.TODO(), instanceEndpoint)
			if err != nil {
				fmt.Println("[OpenMCP DNS Endpoint Controller] : ",err)
				//return reconcile.Result{}, nil
			}

		} else {
			// Error !
			fmt.Println("[OpenMCP DNS Endpoint Controller] : ",err)
			//return reconcile.Result{}, nil
		}
	} else if errors.IsNotFound(err) {
		// OpenMCPServiceDNSRecord Deleted -> Delete

		instanceEndpoint := OpenMCPEndpointDeleteObjectFromServiceDNS(req.Namespace,  req.Name)
		err :=r.live.Delete(context.TODO(), instanceEndpoint)
		if err == nil {
			fmt.Println("[OpenMCP DNS Endpoint Controller] : Deleted '", req.Name+"'")
			//return reconcile.Result{}, nil
		}

	}

	// OpenMCPIngressDNSRecord 삭제 요청인 경우 종료
	instanceIngressRecord := &ketiv1alpha1.OpenMCPIngressDNSRecord{}
	err = r.live.Get(context.TODO(), req.NamespacedName, instanceIngressRecord)
	if err == nil {
		//if instanceIngressRecord.Status.DNS == nil {
		//	return reconcile.Result{}, nil
		//}

		// Ingress의 Domains을 구함
		instanceOpenMCPIngress := &resketiv1alpha1.OpenMCPIngress{}
		err = r.live.Get(context.TODO(), req.NamespacedName, instanceOpenMCPIngress)
		if err != nil {
			fmt.Println("[OpenMCP DNS Endpoint Controller] : ",err)
		}

		domains := []string{}
		for _, rule := range instanceOpenMCPIngress.Spec.Template.Spec.Rules {
			domains = append(domains, rule.Host)
		}

		nsn := types.NamespacedName{
			Namespace: req.Namespace,
			Name:      "ingress-"+req.Name,
		}
		instanceEndpoint := &ketiv1alpha1.OpenMCPDNSEndpoint{}
		err := r.live.Get(context.TODO(), nsn, instanceEndpoint)

		if err == nil {
			// Already Exist -> Update
			instanceEndpoint := OpenMCPEndpointUpdateObjectFromIngressDNS(instanceEndpoint, instanceIngressRecord, req.Namespace,  req.Name, domains)
			err = r.live.Update(context.TODO(), instanceEndpoint)
			if err != nil {
				fmt.Println("[OpenMCP DNS Endpoint Controller] : ",err)
				//return reconcile.Result{}, nil
			}
		} else if errors.IsNotFound(err) {
			// Not Exist -> Create
			instanceEndpoint := OpenMCPEndpointCreateObjectFromIngressDNS(instanceIngressRecord, req.Namespace,  req.Name, domains)
			err = r.live.Create(context.TODO(), instanceEndpoint)
			if err != nil {
				fmt.Println("[OpenMCP DNS Endpoint Controller] : ",err)
				//return reconcile.Result{}, nil
			}

		} else {
			// Error !
			fmt.Println("[OpenMCP DNS Endpoint Controller] : ",err)
			//return reconcile.Result{}, nil
		}
	} else if errors.IsNotFound(err) {
		// OpenMCPServiceDNSRecord Deleted -> Delete

		instanceEndpoint := OpenMCPEndpointDeleteObjectFromIngressDNS(req.Namespace,  req.Name)
		err :=r.live.Delete(context.TODO(), instanceEndpoint)
		if err == nil {
			fmt.Println("[OpenMCP DNS Endpoint Controller] : Deleted '", req.Name+"'")
			//return reconcile.Result{}, nil
		}

	}

	return reconcile.Result{}, nil // err
}
func CreateEndpoint(dnsName string, recordTTL ketiv1alpha1.TTL, recordType string, targets []string) *ketiv1alpha1.Endpoint {
	endpoint := &ketiv1alpha1.Endpoint{
		DNSName:    dnsName,
		Targets:    targets,
		RecordType: recordType,
		RecordTTL:  recordTTL,
		Labels:     nil,
	}
	return endpoint
}
func CreateEndpointsFromServiceDNS(instanceServiceRecord *ketiv1alpha1.OpenMCPServiceDNSRecord, namespace, name string) []*ketiv1alpha1.Endpoint {
	endpoints :=  []*ketiv1alpha1.Endpoint{}

	domainRef := instanceServiceRecord.Spec.DomainRef
	recordTTL := instanceServiceRecord.Spec.RecordTTL
	domain := instanceServiceRecord.Status.Domain
	recordType := "A"

	targetsAll := []string{}
	for _, dns := range instanceServiceRecord.Status.DNS {
		if dns.Region == "" || dns.Zones == nil {
			continue
		}
		region := dns.Region
		dnsName := name+"."+namespace+"."+domainRef+".svc." + region + "." + domain

		targets := []string{}
		for _, ingress := range dns.LoadBalancer.Ingress {
			targets = append(targets, ingress.IP)
			targetsAll = append(targetsAll, ingress.IP)
		}


		// Region만 존재하는 DNS
		endpoint := CreateEndpoint(dnsName, recordTTL, recordType, targets)
		endpoints = append(endpoints, endpoint)

		for _, zone := range dns.Zones{
			dnsName := name+"."+namespace+"."+domainRef+".svc." + zone + "." + region + "." + domain

			// Region, Zone 둘다 존재하는 DNS
			endpoint := CreateEndpoint(dnsName, recordTTL, recordType, targets)
			endpoints = append(endpoints, endpoint)
		}
	}

	// Region Zone 둘다 존재하지 않는 DNS
	if domain != ""{
		dnsName := name+"."+namespace+"."+domainRef+".svc." + domain
		endpoint := CreateEndpoint(dnsName, recordTTL, recordType, targetsAll)
		endpoints = append(endpoints, endpoint)
	}
	return endpoints

}
func OpenMCPEndpointCreateObjectFromServiceDNS(instanceServiceRecord *ketiv1alpha1.OpenMCPServiceDNSRecord, namespace, name string) *ketiv1alpha1.OpenMCPDNSEndpoint {

	endpoints := CreateEndpointsFromServiceDNS(instanceServiceRecord, namespace, name)

	instanceEndpoint := &ketiv1alpha1.OpenMCPDNSEndpoint{
		ObjectMeta: v1.ObjectMeta{
			Name: "service-"+name,
			Namespace: namespace,
		},
		Spec: ketiv1alpha1.OpenMCPDNSEndpointSpec{
			Endpoints: endpoints,
			Domains: []string{instanceServiceRecord.Status.Domain},
		},
		Status: ketiv1alpha1.OpenMCPDNSEndpointStatus{},
	}
	return instanceEndpoint
}
func OpenMCPEndpointUpdateObjectFromServiceDNS(instanceEndpoint *ketiv1alpha1.OpenMCPDNSEndpoint, instanceServiceRecord *ketiv1alpha1.OpenMCPServiceDNSRecord, namespace, name string) *ketiv1alpha1.OpenMCPDNSEndpoint {

	endpoints := CreateEndpointsFromServiceDNS(instanceServiceRecord, namespace, name)
	instanceEndpoint.Spec.Endpoints = endpoints
	instanceEndpoint.Spec.Domains = []string{instanceServiceRecord.Status.Domain}
	return instanceEndpoint
}
func OpenMCPEndpointDeleteObjectFromServiceDNS(namespace, name string) *ketiv1alpha1.OpenMCPDNSEndpoint {
	instanceEndpoint := &ketiv1alpha1.OpenMCPDNSEndpoint{
		ObjectMeta: v1.ObjectMeta{
			Name: "service-"+name,
			Namespace: namespace,
		},
		Spec:       ketiv1alpha1.OpenMCPDNSEndpointSpec{},
		Status:     ketiv1alpha1.OpenMCPDNSEndpointStatus{},
	}
	return instanceEndpoint

}





func CreateEndpointsFromIngressDNS(instanceIngressRecord *ketiv1alpha1.OpenMCPIngressDNSRecord, namespace, name string) []*ketiv1alpha1.Endpoint {
	endpoints :=  []*ketiv1alpha1.Endpoint{}

	recordTTL := instanceIngressRecord.Spec.RecordTTL
	recordType := "A"

	for _, dns := range instanceIngressRecord.Status.DNS {
		for _, host := range dns.Hosts {
			dnsName := host
			targets := []string{}
			for _, ingress := range dns.LoadBalancer.Ingress {
				targets = append(targets, ingress.IP)
			}
			endpoint := CreateEndpoint(dnsName, recordTTL, recordType, targets)
			endpoints = append(endpoints, endpoint)
		}

	}

	return endpoints

}
func OpenMCPEndpointCreateObjectFromIngressDNS(instanceIngressRecord *ketiv1alpha1.OpenMCPIngressDNSRecord, namespace, name string, domains []string) *ketiv1alpha1.OpenMCPDNSEndpoint {

	endpoints := CreateEndpointsFromIngressDNS(instanceIngressRecord, namespace, name)

	instanceEndpoint := &ketiv1alpha1.OpenMCPDNSEndpoint{
		ObjectMeta: v1.ObjectMeta{
			Name: "ingress-"+name,
			Namespace: namespace,
		},
		Spec: ketiv1alpha1.OpenMCPDNSEndpointSpec{
			Endpoints: endpoints,
			Domains: domains,
		},
		Status: ketiv1alpha1.OpenMCPDNSEndpointStatus{},
	}
	return instanceEndpoint
}
func OpenMCPEndpointUpdateObjectFromIngressDNS(instanceEndpoint *ketiv1alpha1.OpenMCPDNSEndpoint, instanceIngressRecord *ketiv1alpha1.OpenMCPIngressDNSRecord, namespace, name string, domains []string) *ketiv1alpha1.OpenMCPDNSEndpoint {

	endpoints := CreateEndpointsFromIngressDNS(instanceIngressRecord, namespace, name)
	instanceEndpoint.Spec.Endpoints = endpoints
	instanceEndpoint.Spec.Domains = domains
	return instanceEndpoint
}
func OpenMCPEndpointDeleteObjectFromIngressDNS(namespace, name string) *ketiv1alpha1.OpenMCPDNSEndpoint {
	instanceEndpoint := &ketiv1alpha1.OpenMCPDNSEndpoint{
		ObjectMeta: v1.ObjectMeta{
			Name: "ingress-"+name,
			Namespace: namespace,
		},
		Spec:       ketiv1alpha1.OpenMCPDNSEndpointSpec{},
		Status:     ketiv1alpha1.OpenMCPDNSEndpointStatus{},
	}
	return instanceEndpoint

}