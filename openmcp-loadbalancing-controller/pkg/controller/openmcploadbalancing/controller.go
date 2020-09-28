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

package openmcploadbalancing // import "admiralty.io/multicluster-controller/examples/openmcploadbalancing/pkg/controller/openmcploadbalancing"

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/apis"
	"openmcp/openmcp/util/clusterManager"
	"os"
	"strings"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/controller/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	resourceapis "openmcp/openmcp/openmcp-resource-controller/apis"
	resourcev1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"

	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/clusterregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/ingressregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/loadbalancingregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/serviceregistry"
)


var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called NewController")
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

	if err := resourceapis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}


	if err := co.WatchResourceReconcileObject(live, &resourcev1alpha1.OpenMCPIngress{}, controller.WatchOptions{}); err != nil {
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
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called Reconcile")

	i += 1
	omcplog.V(3).Info("********* [", i, "] *********")
	omcplog.V(3).Info(req.Context, " / ", req.Namespace, " / ", req.Name)
	instance := &resourcev1alpha1.OpenMCPIngress{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)


	omcplog.V(3).Info("instance Name: ", instance.Name)
	omcplog.V(3).Info("instance Namespace : ", instance.Namespace)

	// delete
	if err != nil && errors.IsNotFound(err) {
		omcplog.V(2).Info("[OpenMCP Loadbalancing Controller] Delete Registry")
		//ingress & ingressName Registry Delete
		if errors.IsNotFound(err) {
			ingressURLs, _ := ingressregistry.Registry.Lookup(loadbalancing.IngressRegistry, req.NamespacedName.Name)
			ingressregistry.Registry.Delete(loadbalancing.IngressRegistry, req.NamespacedName.Name)
			for _, ingressURL := range ingressURLs {
				checkURL, _ := ingressregistry.Registry.CheckURL(loadbalancing.IngressRegistry, ingressURL)
				if checkURL == false {
					s := strings.Split(ingressURL, "/")
					host := s[0]
					var path string
					if len(s) == 1 {
						path = "/"
					} else {
						path = s[1]
					}
					loadbalancingregistry.Registry.IngressDelete(loadbalancing.LoadbalancingRegistry, host, path)
				}
			}
			for _, rule := range instance.Spec.Template.Spec.Rules {
				for _, paths := range rule.HTTP.Paths {
					serviceName := paths.Backend.ServiceName
					serviceregistry.Registry.Delete(loadbalancing.ServiceRegistry, serviceName)
				}
			}
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, nil
	} else { // ingress & ingressName Registry Add or Update
		//add
		omcplog.V(2).Info("[OpenMCP Loadbalancing Controller] Registry Add or Update")

		ingressName := instance.Name
		ingressHosts, err := ingressregistry.Registry.Lookup(loadbalancing.IngressRegistry, ingressName)

		if ingressHosts == nil && errors.IsNotFound(err) {
			for _, ingressHost := range ingressHosts {
				s := strings.Split(ingressHost, "/")
				host := s[0]
				var path string
				if len(s) == 1 {
					path = "/"
				} else {
					path = s[1]
				}
				loadbalancingregistry.Registry.IngressDelete(loadbalancing.LoadbalancingRegistry, host, path)
			}
			ingressregistry.Registry.Delete(loadbalancing.IngressRegistry, ingressName)
		}

		for _, rule := range instance.Spec.Template.Spec.Rules {
			host := rule.Host
			for _, paths := range rule.HTTP.Paths {
				path := paths.Path
				url := host + path
				if len(path) > 1 && path[0] == '/' {
					path = path[1:]
				} else if path == "" {
					path = "/"
				}
				serviceName := paths.Backend.ServiceName
				serviceregistry.Registry.Delete(loadbalancing.ServiceRegistry, serviceName)
				for _, cluster := range cm.Cluster_list.Items {
					cluster_client := cm.Cluster_genClients[cluster.Name]
					fmt.Println(cluster.Name)
					found := &corev1.Service{}
					err := cluster_client.Get(context.TODO(), found, instance.Namespace, serviceName)
					if err != nil && errors.IsNotFound(err) {
						omcplog.V(0).Info(err)
						omcplog.V(0).Info("Service Not Found")
					} else { // Add
						loadbalancingregistry.Registry.Add(loadbalancing.LoadbalancingRegistry, host, path, serviceName)
						serviceregistry.Registry.Add(loadbalancing.ServiceRegistry, serviceName, cluster.Name)
						loadbalancing.RR[host+path] = 0
					}
				}
				ingressregistry.Registry.Add(loadbalancing.IngressRegistry, ingressName, url)
			}
		}
		return reconcile.Result{}, nil // err
	}
	return reconcile.Result{}, nil // err


}

func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called ListKubeFedClusters")

	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		omcplog.V(0).Info("Error retrieving list of federated clusters: %+v", err)
	}
	if len(clusterList.Items) == 0 {
		omcplog.V(0).Info("No federated clusters found")
	}
	return clusterList
}

func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, client genericclient.Client, fedNamespace string) map[string]*rest.Config {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called KubeFedClusterConfigs")
	clusterConfigs := make(map[string]*rest.Config)
	for _, cluster := range clusterList.Items {
		config, _ := util.BuildClusterConfig(&cluster, client, fedNamespace)
		clusterConfigs[cluster.Name] = config
	}
	return clusterConfigs
}
func KubeFedClusterClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called KubeFedClusterClients")

	cluster_clients := make(map[string]genericclient.Client)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		cluster_config := cluster_configs[clusterName]
		cluster_client := genericclient.NewForConfigOrDie(cluster_config)
		cluster_clients[clusterName] = cluster_client
	}
	return cluster_clients
}


var OPENMCP_IP = ""

func initRegistry() {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called initRegistry")


	for _, cluster := range cm.Cluster_list.Items {
		loadbalancing.ClusterRegistry[cluster.Name] = map[string]string{}

		config, _ := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
		clientset, _ := kubernetes.NewForConfig(config)
		nodes, err := clientset.CoreV1().Nodes().List( metav1.ListOptions{})

		if err != nil {
			omcplog.V(0).Info(err)
		}
		node := nodes.Items[0]
		country := node.Labels["failure-domain.beta.kubernetes.io/zone"]
		continent := node.Labels["failure-domain.beta.kubernetes.io/region"]


		loadbalancing.ClusterRegistry[cluster.Name]["Country"] = country
		loadbalancing.ClusterRegistry[cluster.Name]["Continent"] = continent

		found := &corev1.Service{}
		cluster_client := cm.Cluster_genClients[cluster.Name]
		err = cluster_client.Get(context.TODO(), found, "ingress-nginx", "ingress-nginx")
		if err != nil {
			omcplog.V(0).Info("Cluster Ingress Controller Not Found")
		} else {
			if found.Spec.Type == "LoadBalancer" {
				omcplog.V(5).Info("[OpenMCP Loadbalancing Controller] Service Type LoadBalancer")
				if len (found.Status.LoadBalancer.Ingress) > 0 {
					loadbalancing.ClusterRegistry[cluster.Name]["IngressIP"] = found.Status.LoadBalancer.Ingress[0].IP
				}
			} else {
				omcplog.V(5).Info("[OpenMCP Loadbalancing Controller] Service Type NodePort")
				port := fmt.Sprint(found.Spec.Ports[0].NodePort)
				nodeIP := node.Status.Addresses[0].Address
				omcplog.V(5).Info("[OpenMCP Loadbalancing Controller] NodePort :" + port)
				omcplog.V(5).Info("[OpenMCP Loadbalancing Controller] NodeIP :" + nodeIP)
				loadbalancing.ClusterRegistry[cluster.Name]["IngressIP"] = nodeIP + ":" + port
			}
		}
	}

}



func Loadbalancer(openmcpIP string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called Reconcile")

	initRegistry()

	lb := os.Getenv("LB")

	if lb == "RR" {
		http.HandleFunc("/", loadbalancing.NewMultipleHostReverseProxyRR(loadbalancing.LoadbalancingRegistry, loadbalancing.ClusterRegistry, loadbalancing.CountryRegistry, loadbalancing.ServiceRegistry, openmcpIP))

	} else {
		http.HandleFunc("/", loadbalancing.NewMultipleHostReverseProxy(loadbalancing.LoadbalancingRegistry, loadbalancing.ClusterRegistry, loadbalancing.CountryRegistry, loadbalancing.ServiceRegistry, openmcpIP))
	}
	http.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "add")
	})
	http.HandleFunc("/delete", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "delete")
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%v\n", loadbalancing.LoadbalancingRegistry)
		fmt.Fprintf(w, "")
		fmt.Fprintf(w, "%v\n", loadbalancing.IngressRegistry)
		fmt.Fprintf(w, "")
		fmt.Fprintf(w, "%v\n", loadbalancing.ServiceRegistry)
		fmt.Fprintf(w, "")
		fmt.Fprintf(w, "%v\n", loadbalancing.ClusterRegistry)
	})
	omcplog.V(3).Info("ready")
	log.Fatal(http.ListenAndServe(":80", nil))
}
