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

package main

import (
	"os"
	"log"
	"fmt"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/controller/openmcploadbalancing"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/controller/service"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing"

	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/reshape"
	"openmcp/openmcp/util/controller/logLevel"
)

//
//type ClusterManager struct {
//	Fed_namespace string
//	Host_config *rest.Config
//	Host_client genericclient.Client
//	Cluster_list *fedv1b1.KubeFedClusterList
//	Cluster_configs map[string]*rest.Config
//	Cluster_clients map[string]genericclient.Client
//}
//
//func NewClusterManager() *ClusterManager {
//	fed_namespace := "kube-federation-system"
//	host_config, _ := rest.InClusterConfig()
//	host_client := genericclient.NewForConfigOrDie(host_config)
//	cluster_list := ListKubeFedClusters(host_client, fed_namespace)
//	cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)
//	cluster_clients := KubeFedClusterClients(cluster_list, cluster_configs)
//
//	cm := &ClusterManager{
//		Fed_namespace: fed_namespace,
//		Host_config: host_config,
//		Host_client: host_client,
//		Cluster_list: cluster_list,
//		Cluster_configs: cluster_configs,
//		Cluster_clients: cluster_clients,
//	}
//	return cm
//}
//
//
//func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
//	clusterList := &fedv1b1.KubeFedClusterList{}
//	err := client.List(context.TODO(), clusterList, namespace)
//	if err != nil {
//		fmt.Println("Error retrieving list of federated clusters: %+v", err)
//	}
//	if len(clusterList.Items) == 0 {
//		fmt.Println("No federated clusters found")
//	}
//	return clusterList
//}
//
//func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, client genericclient.Client, fedNamespace string) map[string]*rest.Config {
//	clusterConfigs := make(map[string]*rest.Config)
//	for _, cluster := range clusterList.Items {
//		config, _ := util.BuildClusterConfig(&cluster, client, fedNamespace)
//		clusterConfigs[cluster.Name] = config
//	}
//	return clusterConfigs
//}
//
//func KubeFedClusterClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {
//
//	cluster_clients := make(map[string]genericclient.Client)
//	for _, cluster := range clusterList.Items {
//		clusterName := cluster.Name
//		cluster_config := cluster_configs[clusterName]
//		cluster_client := genericclient.NewForConfigOrDie(cluster_config)
//		cluster_clients[clusterName] = cluster_client
//	}
//	return cluster_clients
//}
//
//
////var IngressRegistry = ingressregistry.DefaultRegistry{}
////var ClusterRegistry = clusterregistry.DefaultClusterInfo{}
////var CountryRegistry = countryregistry.DefaultCountryInfo{}
//var m = flag.Int("m", traceroute.DEFAULT_MAX_HOPS, `Set the max time-to-live (max number of hops) used in outgoing probe packets (default is 64)`)
//var f = flag.Int("f", traceroute.DEFAULT_FIRST_HOP, `Set the first used time-to-live, e.g. the first hop (default is 1)`)
//var q = flag.Int("q", 1, `Set the number of probes per "ttl" to nqueries (default is one probe).`)
//
//
//func getHopScore(ip string) string {
//	//var m = flag.Int("m", traceroute.DEFAULT_MAX_HOPS, `Set the max time-to-live (max number of hops) used in outgoing probe packets (default is 64)`)
//	//var f = flag.Int("f", traceroute.DEFAULT_FIRST_HOP, `Set the first used time-to-live, e.g. the first hop (default is 1)`)
//	//var q = flag.Int("q", 1, `Set the number of probes per "ttl" to nqueries (default is one probe).`)
//	fmt.Println("hop check Start")
//	fmt.Println(time.Now())
//	flag.Parse()
//	host := "8.8.8.8"
//	options := traceroute.TracerouteOptions{}
//	options.SetRetries(*q - 1)
//	options.SetMaxHops(*m + 1)
//	options.SetFirstHop(*f)
//	fmt.Println(time.Now())
//	_, err := net.ResolveIPAddr("ip", host)
//	if err != nil {
//		return ""
//	}
//	fmt.Println(time.Now())
//	c := make(chan traceroute.TracerouteHop, 0)
//	go func() {
//		for {
//			_, ok := <-c
//			if !ok {
//				return
//			}
//		}
//	}()
//	fmt.Println(time.Now())
//	result, err := traceroute.Traceroute(host, &options, c)
//	if err != nil {
//		fmt.Printf("Error: ", err)
//	}
//	fmt.Println(time.Now())
//	hopLen := len(result.Hops)
//	hopNum := result.Hops[hopLen - 1].TTL
//	fmt.Println(time.Now())
//	fmt.Println("End Hop Check")
//	fmt.Println(hopNum)
//	return strconv.Itoa(hopNum)
//}
//
//
//func updateHopScore() {
//	for {
//		fmt.Println("*****Update Hop Score*****")
//		cm := NewClusterManager()
//		for _, cluster := range cm.Cluster_list.Items {
//			config, _ := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
//			clientset, _ := kubernetes.NewForConfig(config)
//			nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
//			if err != nil {
//				fmt.Println(err)
//			}
//			node := nodes.Items[0]
//			nodeIP := node.Status.Addresses[0].Address
//			fmt.Println(nodeIP)
//			hopScore := getHopScore(nodeIP)
//			lb_proxy.ClusterRegistry[cluster.Name]["HopScore"] = hopScore
//			//ClusterRegistry[cluster.Name]["HopScore"] = hopScore
//		}
//		time.Sleep(time.Second * 3)
//	}
//}
//
//
//func initRegistry() {
//	fmt.Println("*****Init Cluster*****")
//	cm := NewClusterManager()
//
//	for _, cluster := range cm.Cluster_list.Items {
//		lb_proxy.ClusterRegistry[cluster.Name] = map[string]string{}
//		config, _ := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
//		clientset, _ := kubernetes.NewForConfig(config)
//		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
//
//		if err != nil {
//			fmt.Println(err)
//		}
//		node := nodes.Items[0]
//		//label로부터 zone, region 추출
//		country := node.Labels["failure-domain.beta.kubernetes.io/zone"]
//		continent := node.Labels["failure-domain.beta.kubernetes.io/region"]
//
//
//		//iso 적용 전 디버깅용 정의
//		//******************
//		if country == "korea" {
//			country = "KR"
//		} else if country == "" {
//			country = "JP"
//		}
//		if continent == "asia" {
//			continent = "AS"
//		} else if continent == "" {
//			continent = "AS"
//		}
//		//********************
//
//
//		lb_proxy.ClusterRegistry[cluster.Name]["Country"] = country
//		lb_proxy.ClusterRegistry[cluster.Name]["Continent"] = continent
//
//		//label로 부터 위도,경도 추출
//		//lat, lon := convertGeo(node.Labels["latitude"], node.Labels["longitude"])
//		//ClusterRegistry[cluster.Name]["Latitude"] = lat
//		//ClusterRegistry[cluster.Name]["Longitude"] = lon
//
//
//		//member cluster ingress 주소 초기화
//		found := &corev1.Service{}
//		cluster_client := cm.Cluster_clients[cluster.Name]
//		err = cluster_client.Get(context.TODO(), found, "ingress-nginx", "ingress-nginx-controller")
//		if err != nil {
//			fmt.Println(cluster.Name)
//			fmt.Println("Cluster Ingress Controller Not Found")
//		} else {
//			if found.Spec.Type == "LoadBalancer" {
//				fmt.Println("LoadBalancer")
//				lb_proxy.ClusterRegistry[cluster.Name]["IngressIP"] = found.Status.LoadBalancer.Ingress[0].IP
//			} else { //NodrPort
//				fmt.Println("Node Port")
//				port := fmt.Sprint(found.Spec.Ports[0].NodePort)
//				nodeIP := node.Status.Addresses[0].Address
//				fmt.Println(port)
//				fmt.Println(nodeIP)
//				lb_proxy.ClusterRegistry[cluster.Name]["IngressIP"] = nodeIP + ":" + port
//			}
//		}
//	}
//
//
//	//openmcp ingress
//	found := &v1beta1.Ingress{}
//	openmcp_client := cm.Host_client
//	err := openmcp_client.Get(context.TODO(), found, "openmcp", "test-ingress")
//	if err != nil {
//		fmt.Println("OpenMCP Ingress Not found")
//	}
//
//	//found_test := &v1alpha1.OpenMCPIngress{}
//	//err2 := openmcp_client.Get(context.TODO(), found_test, "openmcp", "example-openmcpingress")
//	////list로 다긇어올수있음
//
//	//fmt.Println(err2)
//	//if err2 != nil && errors.IsNotFound(err2){
//	//	fmt.Println(err2)
//	//	fmt.Println("***********OpenMCPIngress Not found************")
//	//}
//	//
//	//fmt.Println(found_test)
//
//	for _,rule := range found.Spec.Rules {
//		host := rule.Host
//		lb_proxy.IngressRegistry[host] = map[string][]string{}
//		//IngressRegistry[host] = map[string][]string{}
//		for _,paths := range rule.HTTP.Paths {
//			path := paths.Path
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
//				if err != nil {
//					fmt.Println("Service Not Found")
//				} else {
//					lb_proxy.IngressRegistry[host][path] = append(lb_proxy.IngressRegistry[host][path], cluster.Name)
//					//IngressRegistry[host][path] = append(IngressRegistry[host][path], cluster.Name)
//				}
//			}
//		}
//	}
//
//	lb_proxy.CountryRegistry["KR"] = "AS"
//	lb_proxy.CountryRegistry["JP"] = "AS"
//	fmt.Println(lb_proxy.CountryRegistry)
//	//fmt.Println(found.Spec.Rules[0].Host)
//	//fmt.Println(found.Spec.Rules[1].Host)
//	//fmt.Println(found.Spec.Rules[0].HTTP.Paths[0].Path)
//	//fmt.Println(found.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName)
//
//	//IngressRegistry["keti.test.com"] = map[string][]string{}
//	//IngressRegistry["keti.test.com"]["service/test"] = append(IngressRegistry["keti.test.com"]["service/test"] ,"cluster1")
//	//IngressRegistry["keti.test.com"]["service/test"] = append(IngressRegistry["keti.test.com"]["service/test"] ,"cluster2")
//	//IngressRegistry["keti.test.com"]["service/test"] = append(IngressRegistry["keti.test.com"]["service/test"] ,"cluster3")
//	//
//	//ClusterRegistry["cluster1"] = map[string]string{}
//	//ClusterRegistry["cluster1"]["Latitude"] = "37.5923"
//	//ClusterRegistry["cluster1"]["Longitude"] = "126.9548"
//	//ClusterRegistry["cluster1"]["IngressIP"] = "10.0.3.201"
//	//ClusterRegistry["cluster1"]["Country"] = "KR"
//	//ClusterRegistry["cluster1"]["Continent"] = "Asia"
//	lb_proxy.ClusterRegistry["cluster1"]["ResourceScore"] = "90"
//	lb_proxy.ClusterRegistry["cluster1"]["HopScore"] = "0"
//	//
//	//ClusterRegistry["cluster2"] = map[string]string{}
//	//ClusterRegistry["cluster2"]["Latitude"] = "39.9035"
//	//ClusterRegistry["cluster2"]["Longitude"] = "116.388"
//	//ClusterRegistry["cluster2"]["IngressIP"] = "10.0.3.211"
//	//ClusterRegistry["cluster2"]["Country"] = "China"
//	//ClusterRegistry["cluster2"]["Continent"] = "Asia"
//	lb_proxy.ClusterRegistry["cluster2"]["ResourceScore"] = "60"
//	lb_proxy.ClusterRegistry["cluster2"]["HopScore"] = "0"
//	//
//	//ClusterRegistry["cluster3"] = map[string]string{}
//	//ClusterRegistry["cluster3"]["Latitude"] = "37.751"
//	//ClusterRegistry["cluster3"]["Longitude"] = "-97.822"
//	//ClusterRegistry["cluster3"]["IngressIP"] = "10.0.3.202"
//	//ClusterRegistry["cluster3"]["Country"] = "US"
//	//ClusterRegistry["cluster3"]["Continent"] = "North America"
//	//ClusterRegistry["cluster3"]["ResourceScore"] = "10"
//
//}
//
//func loadbalancer() {
//
//	initRegistry()
//	go updateHopScore()
//
//	http.HandleFunc("/", lb_proxy.NewMultipleHostReverseProxy(lb_proxy.IngressRegistry, lb_proxy.ClusterRegistry, lb_proxy.CountryRegistry))
//	http.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
//		fmt.Fprintf(writer,"add")
//		//registry.Registry.Add(ServiceRegistry,"service2", "v2", "10.0.3.202:80")
//		ingressregistry.Registry.Add(lb_proxy.IngressRegistry, "keti.host.com" , "service2/v2", "10.0.3.202:80")
//		clusterregistry.Registry.Add(lb_proxy.ClusterRegistry, "cluster4", "15.15", "4.1", "10.0.0.10", "KP", "Asia", "80", "70")
//	})
//	http.HandleFunc("/delete", func(writer http.ResponseWriter, request *http.Request) {
//		fmt.Fprintf(writer,"delete")
//		//registry.Registry.Add(ServiceRegistry,"service2", "v2", "10.0.3.202:80")
//		ingressregistry.Registry.Delete(lb_proxy.IngressRegistry, "keti.test.com" , "test", "cluster1")
//	})
//	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
//		fmt.Fprintf(w, "%v\n", lb_proxy.IngressRegistry)
//		fmt.Fprintf(w, "%v\n", lb_proxy.ClusterRegistry)
//		fmt.Println(req.Host)
//		fmt.Println(req)
//	})
//	println("ready")
//	log.Fatal(http.ListenAndServe(":80", nil))
//}
//

func main() {
	logLevel.KetiLogInit()

	for {
		cm := clusterManager.NewClusterManager()

		host_ctx := "openmcp"
		namespace := "openmcp"

		host_cfg := cm.Host_config
		live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})

		ghosts := []*cluster.Cluster{}

		SERVER_IP := os.Getenv("GRPC_SERVER")

		for _, ghost_cluster := range cm.Cluster_list.Items {
			ghost_ctx := ghost_cluster.Name
			ghost_cfg := cm.Cluster_configs[ghost_ctx]

			ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})

			ghosts = append(ghosts, ghost)
		}
		for _, ghost := range ghosts {
			fmt.Println(ghost.Name)
		}
		co, _ := openmcploadbalancing.NewController(live, ghosts, namespace)
		serviceWatch, _ := service.NewController(live, ghosts, namespace)
		reshape_cont, _ := reshape.NewController(live, ghosts, namespace)
		loglevel_cont, _ := logLevel.NewController(live, ghosts, namespace)

		m := manager.New()
		m.AddController(co)
		m.AddController(serviceWatch)
		m.AddController(reshape_cont)
		m.AddController(loglevel_cont)
		go openmcploadbalancing.Loadbalancer(SERVER_IP)
		tempCluster := []string {}
		go loadbalancing.RequestResourceScore(tempCluster,SERVER_IP)
		go loadbalancing.GetPolicy()

		stop := reshape.SetupSignalHandler()

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}
	}
}
