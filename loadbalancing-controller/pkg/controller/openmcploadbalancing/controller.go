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
	"math/rand"
	"net/http"
	"sort"
	"time"

	"loadbalancing-controller/pkg/apis"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"admiralty.io/multicluster-controller/pkg/reference"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/controller/util"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ketiv1alpha1 "loadbalancing-controller/pkg/apis/keti/v1alpha1"
	resourceapis "resource-controller/apis"
	v1alpha1 "resource-controller/apis/keti/v1alpha1"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"

	"github.com/HanJaeseung/LoadBalancing"
	"github.com/HanJaeseung/LoadBalancing/clusterregistry"
	"github.com/HanJaeseung/LoadBalancing/ingressregistry"
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
	fmt.Println("api check")

	fmt.Printf("%T, %s\n", live, live.GetClusterName())
	//if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPLoadbalancing{}, controller.WatchOptions{}); err != nil {
	//	return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	//}

	if err := co.WatchResourceReconcileObject(live, &corev1.Service{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	//if err := co.WatchResourceReconcileObject(live, &v1beta1.Ingress{}, controller.WatchOptions{}); err != nil {
	//	return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	//}
	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	//for _, ghost := range ghosts {
	//	fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
	//	if err := co.WatchResourceReconcileController(ghost, &appsv1.Deployment{}, controller.WatchOptions{}); err != nil {
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

	i += 1
	fmt.Println("********* [", i, "] *********")
	fmt.Println(req.Context, " / ", req.Namespace, " / ", req.Name)
	cm := NewClusterManager()

	found_test := &v1alpha1.OpenMCPIngress{}
	err2 := cm.Host_client.Get(context.TODO(), found_test, "openmcp", "example-openmcpingress")

	fmt.Println(err2)
	if err2 != nil && errors.IsNotFound(err2) {
		fmt.Println(err2)
		fmt.Println("***********OpenMCPIngress Not found************")
	}

	//lbinstance :=  &ketiv1alpha1.OpenMCPLoadbalancing{}
	//ingressinstance := &v1beta1.Ingress{}

	//test_instance := &corev1.Service{}
	//
	//err2 := cm.Host_client.Get(context.TODO(), test_instance, "openmcp", "openmcp-test")
	//

	//
	//for _, cluster := range cm.Cluster_list.Items {
	//	cluster_client := cm.Cluster_clients[cluster.Name]
	//	err := cluster_client.Get(context.TODO(), test_instance, "openmcp", "openmcp-test")
	//	if err != nil {
	//		fmt.Println(cluster.Name)
	//		fmt.Println("Test Service Not Found")
	//	} else {
	//		fmt.Println("found Service")
	//		fmt.Println(time.Now())
	//	}
	//}
	//fmt.Println(time.Now())

	//if err := r.live.Get(context.TODO(), req.NamespacedName, lbinstance); err != nil {
	//	err = r.live.Get(context.TODO(), req.NamespacedName, ingressinstance)
	//}
	//fmt.Println(r.live.Status().GetAPIVersion())
	// Fetch the OpenMCPDeployment instance
	instance := &ketiv1alpha1.OpenMCPLoadbalancing{}
	//err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	fmt.Println("instance Name: ", instance.Name)
	fmt.Println("instance Namespace : ", instance.Namespace)

	//if err != nil {
	//	fmt.Println(err)
	//	if errors.IsNotFound(err) {
	//		// ...TODO: multicluster garbage collector
	//		// Until then...
	//		fmt.Println("Delete Deployments ..Cluster")
	//		//err := cm.DeleteDeployments(req.NamespacedName)
	//		fmt.Println(cm)
	//		return reconcile.Result{}, err
	//	}
	//	fmt.Println("Error1")
	//	return reconcile.Result{}, err
	//}
	//if instance.Status.ClusterMaps == nil {
	//	if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false && strings.Compare(instance.Spec.Labels["test"], "yes") != 0 {
	//		fmt.Println("Scheduling Start & Create Deployments")
	//		replicas := instance.Spec.Replicas
	//		cluster_replicas_map := cm.Scheduling(replicas)
	//
	//		for _, cluster := range cm.Cluster_list.Items {
	//
	//			if cluster_replicas_map[cluster.Name] == 0 {
	//				continue
	//			}
	//			found := &appsv1.Deployment{}
	//			cluster_client := cm.Cluster_clients[cluster.Name]
	//
	//			err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-deploy")
	//			if err != nil && errors.IsNotFound(err) {
	//				replica := cluster_replicas_map[cluster.Name]
	//				fmt.Println("Cluster '"+cluster.Name+"' Deployed (", replica, " / ", replicas, ")")
	//				dep := r.deploymentForOpenMCPLoadbalancing(req, instance, replica)
	//				err = cluster_client.Create(context.Background(), dep)
	//				if err != nil {
	//					return reconcile.Result{}, err
	//				}
	//			}
	//		}
	//		instance.Status.ClusterMaps = cluster_replicas_map
	//		instance.Status.Replicas = replicas
	//
	//		err := r.live.Status().Update(context.TODO(), instance)
	//		if err != nil {
	//			fmt.Println("Failed to update instance status", err)
	//			return reconcile.Result{}, err
	//		}
	//
	//		return reconcile.Result{}, nil
	//	}
	//} else if instance.Status.Replicas != instance.Spec.Replicas{
	//	fmt.Println("Change Spec Replicas ! ReScheduling Start & Update Deployment")
	//	cluster_replicas_map := cm.ReScheduling(instance.Spec.Replicas, instance.Status.Replicas, instance.Status.ClusterMaps)
	//
	//            for _, cluster := range cm.Cluster_list.Items {
	//		update_replica := cluster_replicas_map[cluster.Name]
	//		cluster_client := cm.Cluster_clients[cluster.Name]
	//
	//		dep := r.deploymentForOpenMCPLoadbalancing(req, instance, update_replica)
	//
	//		found := &appsv1.Deployment{}
	//                    err := cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-deploy")
	//                    if err != nil && errors.IsNotFound(err) {
	//                            // Not Exist Deployment.
	//			if update_replica != 0{
	//				// Create !
	//                                err = cluster_client.Create(context.Background(), dep)
	//	                        if err != nil {
	//		                        return reconcile.Result{}, err
	//			        }
	//			}
	//
	//                    } else if err != nil {
	//			return reconcile.Result{}, err
	//		} else {
	//			// Already Exist Deployment.
	//			if update_replica == 0 {
	//				// Delete !
	//				dep := &appsv1.Deployment{}
	//				err = cluster_client.Delete(context.Background(), dep, req.Namespace, req.Name+"-deploy")
	//
	//                                    if err != nil {
	//                                            return reconcile.Result{}, err
	//                                    }
	//			} else {
	//				// Update !
	//				err = cluster_client.Update(context.TODO(), dep)
	//                                if err != nil {
	//	                                return reconcile.Result{}, err
	//                                }
	//
	//			}
	//
	//		}
	//
	//	}
	//	instance.Status.ClusterMaps = cluster_replicas_map
	//        instance.Status.Replicas = instance.Spec.Replicas
	//	err := r.live.Status().Update(context.TODO(), instance)
	//        if err != nil {
	//                 fmt.Println("Failed to update instance status", err)
	//                 return reconcile.Result{}, err
	//         }
	//
	//} else {
	//	// Check Deployment in cluster
	//	for k, v := range instance.Status.ClusterMaps {
	//		cluster_name := k
	//		replica := v
	//
	//		if replica == 0 {
	//			continue
	//		}
	//		found := &appsv1.Deployment{}
	//		cluster_client := cm.Cluster_clients[cluster_name]
	//		err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-deploy")
	//		if err != nil && errors.IsNotFound(err) {
	//			// Delete Deployment Detected
	//                            fmt.Println("Cluster '" + cluster_name  + "' ReDeployed => ", replica)
	//                            dep := r.deploymentForOpenMCPLoadbalancing(req, instance, replica)
	//                            err = cluster_client.Create(context.Background(), dep)
	//			if err != nil {
	//                                    return reconcile.Result{}, err
	//                            }
	//
	//		}
	//
	//	}
	//
	//
	//}

	return reconcile.Result{}, nil // err
}

func (r *reconciler) deploymentForOpenMCPLoadbalancing(req reconcile.Request, m *ketiv1alpha1.OpenMCPLoadbalancing, replica int32) *appsv1.Deployment {

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-deploy",
			Namespace: m.Namespace,
		},
		Spec: m.Spec.Template.Spec,
	}
	dep.Spec.Replicas = &replica

	reference.SetMulticlusterControllerReference(dep, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return dep
}
func (cm *ClusterManager) DeleteDeployments(nsn types.NamespacedName) error {
	dep := &appsv1.Deployment{}
	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_clients[cluster.Name]
		fmt.Println(nsn.Namespace, nsn.Name)
		err := cluster_client.Get(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")
		if err != nil && errors.IsNotFound(err) {
			// all good
			fmt.Println("Not Found")
			continue
		}
		fmt.Println(cluster.Name, " Delete Start")
		err = cluster_client.Delete(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")
		if err != nil {
			return err
		}
		fmt.Println(cluster.Name, "Delete Complate")
	}
	return nil

}

func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		fmt.Println("Error retrieving list of federated clusters: %+v", err)
	}
	if len(clusterList.Items) == 0 {
		fmt.Println("No federated clusters found")
	}
	return clusterList
}

func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, client genericclient.Client, fedNamespace string) map[string]*rest.Config {
	clusterConfigs := make(map[string]*rest.Config)
	for _, cluster := range clusterList.Items {
		config, _ := util.BuildClusterConfig(&cluster, client, fedNamespace)
		clusterConfigs[cluster.Name] = config
	}
	return clusterConfigs
}
func KubeFedClusterClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {

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
func (cm *ClusterManager) Scheduling(replicas int32) map[string]int32 {
	rand.Seed(time.Now().UTC().UnixNano())

	cluster_replicas_map := make(map[string]int32)

	remain_rep := replicas
	rep := 0
	cluster_len := len(cm.Cluster_list.Items)
	for i, cluster := range cm.Cluster_list.Items {
		if i == cluster_len-1 {
			rep = int(remain_rep)
		} else {
			rep = rand.Intn(int(remain_rep + 1))
		}
		remain_rep = remain_rep - int32(rep)
		cluster_replicas_map[cluster.Name] = int32(rep)

	}
	keys := make([]string, 0)
	for k, _ := range cluster_replicas_map {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("Scheduling Result: ")
	for _, k := range keys {
		v := cluster_replicas_map[k]
		fmt.Println("  ", k, ": ", v)
	}
	return cluster_replicas_map

}
func (cm *ClusterManager) ReScheduling(spec_replicas int32, status_replicas int32, status_cluster_replicas_map map[string]int32) map[string]int32 {
	rand.Seed(time.Now().UTC().UnixNano())

	result_cluster_replicas_map := make(map[string]int32)
	for k, v := range status_cluster_replicas_map {
		result_cluster_replicas_map[k] = v
	}

	action := "dec"
	replica_rate := spec_replicas - status_replicas
	if replica_rate > 0 {
		action = "inc"
	}

	remain_replica := replica_rate

	for remain_replica != 0 {
		cluster_len := len(result_cluster_replicas_map)
		selected_cluster_target_index := rand.Intn(int(cluster_len))

		target_key := keyOf(result_cluster_replicas_map, selected_cluster_target_index)
		if action == "inc" {
			result_cluster_replicas_map[target_key] += 1
			remain_replica -= 1
		} else {
			if result_cluster_replicas_map[target_key] >= 1 {
				result_cluster_replicas_map[target_key] -= 1
				remain_replica += 1
			}
		}
	}
	keys := make([]string, 0)
	for k, _ := range result_cluster_replicas_map {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("ReScheduling Result: ")
	for _, k := range keys {
		v := result_cluster_replicas_map[k]
		prev_v := status_cluster_replicas_map[k]
		fmt.Println("  ", k, ": ", prev_v, " -> ", v)
	}

	return result_cluster_replicas_map

}
func keyOf(my_map map[string]int32, target_index int) string {
	index := 0
	for k, _ := range my_map {
		if index == target_index {
			return k
		}
		index += 1
	}
	return ""

}

//*************************************************************************************************************************************************************

func initRegistry() {
	fmt.Println("*****Init Cluster*****")
	cm := NewClusterManager()

	for _, cluster := range cm.Cluster_list.Items {
		lb_proxy.ClusterRegistry[cluster.Name] = map[string]string{}
		config, _ := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
		clientset, _ := kubernetes.NewForConfig(config)
		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})

		if err != nil {
			fmt.Println(err)
		}
		node := nodes.Items[0]
		//label로부터 zone, region 추출
		country := node.Labels["failure-domain.beta.kubernetes.io/zone"]
		continent := node.Labels["failure-domain.beta.kubernetes.io/region"]

		//iso 적용 전 디버깅용 정의
		//******************
		if country == "korea" {
			country = "KR"
		} else if country == "" {
			country = "JP"
		}
		if continent == "asia" {
			continent = "AS"
		} else if continent == "" {
			continent = "AS"
		}
		//********************

		lb_proxy.ClusterRegistry[cluster.Name]["Country"] = country
		lb_proxy.ClusterRegistry[cluster.Name]["Continent"] = continent

		//label로 부터 위도,경도 추출
		//lat, lon := convertGeo(node.Labels["latitude"], node.Labels["longitude"])
		//ClusterRegistry[cluster.Name]["Latitude"] = lat
		//ClusterRegistry[cluster.Name]["Longitude"] = lon

		//member cluster ingress 주소 초기화
		found := &corev1.Service{}
		cluster_client := cm.Cluster_clients[cluster.Name]
		err = cluster_client.Get(context.TODO(), found, "ingress-nginx", "ingress-nginx-controller")
		if err != nil {
			fmt.Println(cluster.Name)
			fmt.Println("Cluster Ingress Controller Not Found")
		} else {
			if found.Spec.Type == "LoadBalancer" {
				fmt.Println("LoadBalancer")
				lb_proxy.ClusterRegistry[cluster.Name]["IngressIP"] = found.Status.LoadBalancer.Ingress[0].IP
			} else { //NodrPort
				fmt.Println("Node Port")
				port := fmt.Sprint(found.Spec.Ports[0].NodePort)
				nodeIP := node.Status.Addresses[0].Address
				fmt.Println(port)
				fmt.Println(nodeIP)
				lb_proxy.ClusterRegistry[cluster.Name]["IngressIP"] = nodeIP + ":" + port
			}
		}
	}

	//openmcp ingress
	openmcp_client := cm.Host_client
	ingressList := &v1beta1.IngressList{}
	ingressListErr := openmcp_client.List(context.TODO(), ingressList, "openmcp")

	if ingressListErr != nil && errors.IsNotFound(ingressListErr) {
		fmt.Println(ingressListErr)
	}

	for _, ingress := range ingressList.Items {
		for _, rule := range ingress.Spec.Rules {
			host := rule.Host
			lb_proxy.IngressRegistry[host] = map[string][]string{}
			for _, paths := range rule.HTTP.Paths {
				path := paths.Path
				if len(path) > 1 && path[0] == '/' {
					path = path[1:]
				} else if path == "" {
					path = "/"
				}
				serviceName := paths.Backend.ServiceName
				for _, cluster := range cm.Cluster_list.Items {
					cluster_client := cm.Cluster_clients[cluster.Name]
					found := &corev1.Service{}
					err := cluster_client.Get(context.TODO(), found, "openmcp", serviceName)
					if err != nil && errors.IsNotFound(err) {
						fmt.Println("Service Not Found")
					} else {
						lb_proxy.IngressRegistry[host][path] = append(lb_proxy.IngressRegistry[host][path], cluster.Name)
					}
				}
			}
		}
	}
	lb_proxy.CountryRegistry["KR"] = "AS"
	lb_proxy.CountryRegistry["JP"] = "AS"
	fmt.Println(lb_proxy.CountryRegistry)
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
	lb_proxy.ClusterRegistry["cluster1"]["ResourceScore"] = "90"
	lb_proxy.ClusterRegistry["cluster1"]["HopScore"] = "0"
	//
	//ClusterRegistry["cluster2"] = map[string]string{}
	//ClusterRegistry["cluster2"]["Latitude"] = "39.9035"
	//ClusterRegistry["cluster2"]["Longitude"] = "116.388"
	//ClusterRegistry["cluster2"]["IngressIP"] = "10.0.3.211"
	//ClusterRegistry["cluster2"]["Country"] = "China"
	//ClusterRegistry["cluster2"]["Continent"] = "Asia"
	lb_proxy.ClusterRegistry["cluster2"]["ResourceScore"] = "60"
	lb_proxy.ClusterRegistry["cluster2"]["HopScore"] = "0"
	//
	//ClusterRegistry["cluster3"] = map[string]string{}
	//ClusterRegistry["cluster3"]["Latitude"] = "37.751"
	//ClusterRegistry["cluster3"]["Longitude"] = "-97.822"
	//ClusterRegistry["cluster3"]["IngressIP"] = "10.0.3.202"
	//ClusterRegistry["cluster3"]["Country"] = "US"
	//ClusterRegistry["cluster3"]["Continent"] = "North America"
	//ClusterRegistry["cluster3"]["ResourceScore"] = "10"

}

func Loadbalancer() {

	initRegistry()

	http.HandleFunc("/", lb_proxy.NewMultipleHostReverseProxy(lb_proxy.IngressRegistry, lb_proxy.ClusterRegistry, lb_proxy.CountryRegistry))
	http.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "add")
		ingressregistry.Registry.Add(lb_proxy.IngressRegistry, "keti.host.com", "service2/v2", "10.0.3.202:80")
		clusterregistry.Registry.Add(lb_proxy.ClusterRegistry, "cluster4", "15.15", "4.1", "10.0.0.10", "KP", "Asia", "80", "70")
	})
	http.HandleFunc("/delete", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "delete")
		ingressregistry.Registry.Delete(lb_proxy.IngressRegistry, "keti.test.com", "test", "cluster1")
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%v\n", lb_proxy.IngressRegistry)
		fmt.Fprintf(w, "%v\n", lb_proxy.ClusterRegistry)
	})
	println("ready")
	log.Fatal(http.ListenAndServe(":80", nil))
}
