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
	//"github.com/HanJaeseung/LoadBalancing"
	//"github.com/HanJaeseung/LoadBalancing/clusterregistry"
	//"github.com/HanJaeseung/LoadBalancing/ingressregistry"
	//"github.com/HanJaeseung/LoadBalancing/ingressnameregistry"
)



type ClusterManager struct {
	Fed_namespace   string
	Host_config     *rest.Config
	Host_client     genericclient.Client
	Cluster_list    *fedv1b1.KubeFedClusterList
	Cluster_configs map[string]*rest.Config
	Cluster_clients map[string]genericclient.Client
}

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string) (*controller.Controller, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called NewController")
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


	//if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPLoadbalancing{}, controller.WatchOptions{}); err != nil {
	//	return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	//}

	/*if err := co.WatchResourceReconcileObject(live, &v1beta1.Ingress{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}*/

	if err := co.WatchResourceReconcileObject(live, &resourcev1alpha1.OpenMCPIngress{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	//for _, ghost := range ghosts {
	//	fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
	//	if err := co.WatchResourceReconcileController(ghost, &v1beta1.Ingress{}, controller.WatchOptions{}); err != nil {
	//		return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
	//	}
	//}
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
	cm := NewClusterManager()
	//instance := &v1beta1.Ingress{}
	//err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	//type ObjectKey = types.NamespacedName
	//
	//openmcpingressinstance := &resourcev1alpha1.OpenMCPIngress{}
	//err1 := r.live.Get(context.TODO(), ObjectKey{Namespace: "openmcp", Name: "openmcp-ingress-example"}, openmcpingressinstance)
	//fmt.Println("openmcpingressinstance",openmcpingressinstance)

	instance := &resourcev1alpha1.OpenMCPIngress{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	//if err1 != nil {
	//	fmt.Println(err1, "Fail to get openmcpingress")
	//}

	omcplog.V(3).Info("instance Name: ", instance.Name)
	omcplog.V(3).Info("instance Namespace : ", instance.Namespace)

	// delete
	if err != nil && errors.IsNotFound(err) {
		omcplog.V(2).Info("[OpenMCP Loadbalancing Controller] Delete Registry")
		//해당 instance가 없을 경우 ingress & ingressName Registry 삭제
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
	} else { // 해당 instance가 있을 경우 ingress & ingressName Registry Add or Update
		//	add
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
				//queue := loadbalancing.Queue{}
				for _, cluster := range cm.Cluster_list.Items {
					cluster_client := cm.Cluster_clients[cluster.Name]
					fmt.Println(cluster.Name)
					found := &corev1.Service{}
					err := cluster_client.Get(context.TODO(), found, instance.Namespace, serviceName)
					if err != nil && errors.IsNotFound(err) {
						omcplog.V(0).Info(err)
						omcplog.V(0).Info("Service Not Found")
					} else { // Add
						loadbalancingregistry.Registry.Add(loadbalancing.LoadbalancingRegistry, host, path, serviceName)
						serviceregistry.Registry.Add(loadbalancing.ServiceRegistry, serviceName, cluster.Name)
						//queue.Set(cluster.Name)
						//loadbalancing.PV(loadbalancing.RR[host+path]).Set(cluster.Name)
						//loadbalancing.RR[host+path] = append(loadbalancing.RR[host+path], cluster.Name)
						loadbalancing.RR[host+path] = 0
					}
				}
				//loadbalancing.RR[host + path] = queue
				ingressregistry.Registry.Add(loadbalancing.IngressRegistry, ingressName, url)
			}
		}
		return reconcile.Result{}, nil // err
	}
	return reconcile.Result{}, nil // err

	//else { // 해당 instance가 있을 경우 ingress & ingressName Registry Add or Update
	//	//	add
	//	fmt.Println("Registry Add or Update")
	//	ingressName := instance.Name
	//	for _, rule := range instance.Spec.Rules {
	//		host := rule.Host
	//		for _, paths := range rule.HTTP.Paths {
	//			path := paths.Path
	//			if len(path) > 1 && path[0] == '/' {
	//				path = path[1:]
	//			} else if path == "" {
	//				path = "/"
	//			}
	//			serviceName := paths.Backend.ServiceName
	//			serviceregistry.Registry.Delete(loadbalancing.ServiceRegistry, serviceName)
	//			for _, cluster := range cm.Cluster_list.Items {
	//				cluster_client := cm.Cluster_clients[cluster.Name]
	//				found := &corev1.Service{}
	//				err := cluster_client.Get(context.TODO(), found, "openmcp", serviceName)
	//				if err != nil && errors.IsNotFound(err) {
	//					fmt.Println("Service Not Found")
	//				} else {
	//					//기존 host path에 endpoint가 중복되서 들어가는거 방지 test-ingress1, test-ingress2 생성하는데 둘다 내용이 같을 경우
	//					//checkIngress := loadbalancingregistry.Registry.IngressLookup(loadbalancing.LoadbalancingRegistry, host, path, cluster.Name)
	//					//checkEndpoint := serviceregistry.Registry.EndpointCheck(loadbalancing.ServiceRegistry, serviceName, cluster.Name)
	//					//if checkEndpoint == true  {
	//					loadbalancingregistry.Registry.Add(loadbalancing.LoadbalancingRegistry, host, path, serviceName)
	//					serviceregistry.Registry.Add(loadbalancing.ServiceRegistry, serviceName, cluster.Name)
	//					//}
	//				}
	//			}
	//			ingressHosts, err := ingressregistry.Registry.Lookup(loadbalancing.IngressRegistry, ingressName)
	//			if ingressHosts != nil && err == nil {
	//				for _, ingressHost := range ingressHosts {
	//					s := strings.Split(ingressHost, "/")
	//					host := s[0]
	//					path := s[1]
	//					loadbalancingregistry.Registry.IngressDelete(loadbalancing.LoadbalancingRegistry, host, path)
	//				}
	//				ingressregistry.Registry.Delete(loadbalancing.IngressRegistry, ingressName)
	//			}
	//			ingressregistry.Registry.Add(loadbalancing.IngressRegistry, ingressName, host + paths.Path)
	//		}
	//	}
	//}

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

func NewClusterManager() *ClusterManager {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called NewClusterManager")
	fed_namespace := "kube-federation-system"
	host_config, _ := rest.InClusterConfig()
	host_client := genericclient.NewForConfigOrDie(host_config)
	cluster_list := ListKubeFedClusters(host_client, fed_namespace)
	cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)
	cluster_clients := KubeFedClusterClients(cluster_list, cluster_configs)

	cm := &ClusterManager{
		Fed_namespace:   fed_namespace,
		Host_config:     host_config,
		Host_client:     host_client,
		Cluster_list:    cluster_list,
		Cluster_configs: cluster_configs,
		Cluster_clients: cluster_clients,
	}
	return cm
}

var OPENMCP_IP = ""

//*************************************************************************************************************************************************************

func initRegistry() {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Function Called initRegistry")

	cm := NewClusterManager()

	for _, cluster := range cm.Cluster_list.Items {
		loadbalancing.ClusterRegistry[cluster.Name] = map[string]string{}

		config, _ := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
		clientset, _ := kubernetes.NewForConfig(config)
		nodes, err := clientset.CoreV1().Nodes().List( metav1.ListOptions{})

		if err != nil {
			omcplog.V(0).Info(err)
		}
		node := nodes.Items[0]
		//label로부터 zone, region 추출
		country := node.Labels["failure-domain.beta.kubernetes.io/zone"]
		continent := node.Labels["failure-domain.beta.kubernetes.io/region"]

		////iso 적용 전 디버깅용 정의
		////******************
		//if country == "korea" {
		//	country = "KR"
		//} else if country == "" {
		//	country = "JP"
		//}
		//if continent == "asia" {
		//	continent = "AS"
		//} else if continent == "" {
		//	continent = "AS"
		//}
		//********************

		loadbalancing.ClusterRegistry[cluster.Name]["Country"] = country
		loadbalancing.ClusterRegistry[cluster.Name]["Continent"] = continent

		found := &corev1.Service{}
		cluster_client := cm.Cluster_clients[cluster.Name]
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

	////openmcp ingress
	//openmcp_client := cm.Host_client
	//ingressList := &v1beta1.IngressList{}
	//ingressListErr := openmcp_client.List(context.TODO(), ingressList, "openmcp")
	//
	//if ingressListErr != nil && errors.IsNotFound(ingressListErr) {
	//	fmt.Println(ingressListErr)
	//}
	//
	//for _, ingress := range ingressList.Items {
	//	ingressName := ingress.Name
	//	//lb_proxy.IngressNameRegistry[ingressName] = []string{}
	//	for _, rule := range ingress.Spec.Rules {
	//		host := rule.Host
	//		lb_proxy.IngressRegistry[host] = map[string][]string{}
	//		for _, paths := range rule.HTTP.Paths {
	//			path := paths.Path
	//			//lb_proxy.IngressNameRegistry[ingressName] = append(lb_proxy.IngressNameRegistry[ingressName], host + path)
	//			ingressnameregistry.Registry.Add(lb_proxy.IngressNameRegistry, ingressName, host+path)
	//			if len(path) > 1 && path[0] == '/' {
	//				path = path[1:]
	//			} else if path == "" {
	//				path = "/"
	//			}
	//			serviceName := paths.Backend.ServiceName
	//			for _, cluster := range cm.Cluster_list.Items {
	//				cluster_client := cm.Cluster_clients[cluster.Name]
	//				found := &corev1.Service{}
	//				err := cluster_client.Get(context.TODO(), found, "openmcp", serviceName)
	//				if err != nil && errors.IsNotFound(err) {
	//					fmt.Println("Service Not Found")
	//				} else {
	//					ingressregistry.Registry.Add(lb_proxy.IngressRegistry, host, path, cluster.Name)
	//					//lb_proxy.IngressRegistry[host][path] = append(lb_proxy.IngressRegistry[host][path], cluster.Name)
	//				}
	//			}
	//		}
	//	}
	//}
	//loadbalancing.CountryRegistry["KR"] = "AS"
	//loadbalancing.CountryRegistry["JP"] = "AS"

	//fmt.Println(found.Spec.Rules[0].Host)
	//fmt.Println(found.Spec.Rules[1].Host)
	//fmt.Println(found.Spec.Rules[0].HTTP.Paths[0].Path)
	//fmt.Println(found.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName)

	//IngressRegistry["keti.test.com"] = map[string][]string{}
	//IngressRegistry["keti.test.com"]["service/test"] = append(IngressRegistry["keti.test.com"]["service/test"] ,"cluster1")
	//IngressRegistry["keti.test.com"]["service/test"] = append(IngressRegistry["keti.test.com"]["service/test"] ,"cluster2")
	//IngressRegistry["keti.test.com"]["service/test"] = append(IngressRegistry["keti.test.com"]["service/test"] ,"cluster3")
	//
	//ClusterRegistry["cluster1"] = map[string]string{}
	//ClusterRegistry["cluster1"]["Latitude"] = "37.5923"
	//ClusterRegistry["cluster1"]["Longitude"] = "126.9548"
	//ClusterRegistry["cluster1"]["IngressIP"] = "10.0.3.201"
	//ClusterRegistry["cluster1"]["Country"] = "KR"
	//ClusterRegistry["cluster1"]["Continent"] = "Asia"
	//loadbalancing.ClusterRegistry["cluster1"]["ResourceScore"] = "90"
	//loadbalancing.ClusterRegistry["cluster1"]["HopScore"] = "0"
	////
	////ClusterRegistry["cluster2"] = map[string]string{}
	////ClusterRegistry["cluster2"]["Latitude"] = "39.9035"
	////ClusterRegistry["cluster2"]["Longitude"] = "116.388"
	////ClusterRegistry["cluster2"]["IngressIP"] = "10.0.3.211"
	////ClusterRegistry["cluster2"]["Country"] = "China"
	////ClusterRegistry["cluster2"]["Continent"] = "Asia"
	//loadbalancing.ClusterRegistry["cluster2"]["ResourceScore"] = "60"
	//loadbalancing.ClusterRegistry["cluster2"]["HopScore"] = "0"

	//
	//ClusterRegistry["cluster3"] = map[string]string{}
	//ClusterRegistry["cluster3"]["Latitude"] = "37.751"
	//ClusterRegistry["cluster3"]["Longitude"] = "-97.822"
	//ClusterRegistry["cluster3"]["IngressIP"] = "10.0.3.202"
	//ClusterRegistry["cluster3"]["Country"] = "US"
	//ClusterRegistry["cluster3"]["Continent"] = "North America"
	//ClusterRegistry["cluster3"]["ResourceScore"] = "10"

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
		loadbalancingregistry.Registry.Add(loadbalancing.LoadbalancingRegistry, "keti.host.com", "service2/v2", "10.0.3.202:80")
		clusterregistry.Registry.Add(loadbalancing.ClusterRegistry, "cluster4", "15.15", "4.1", "10.0.0.10", "KP", "Asia", "80", "70")
	})
	http.HandleFunc("/delete", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "delete")
		//loadbalancingregistry.Registry.Delete(loadbalancing.LoadbalancingRegistry, "keti.test.com", "test", "cluster1")
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%v\n", loadbalancing.LoadbalancingRegistry)
		//fmt.Fprintf(w, "%v\n", loadbalancing.ClusterRegistry)
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
