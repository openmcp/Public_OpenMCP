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

package openmcpingress // import "admiralty.io/multicluster-controller/examples/openmcpingress/pkg/controller/openmcpingress"

import (
	"context"
	"fmt"
	"github.com/getlantern/deepcopy"
	"github.com/mittwald/go-powerdns/apis/zones"
	corev1 "k8s.io/api/core/v1"
	"openmcp/openmcp/openmcp-dns-controller/pkg/mypdns"

	//"reflect"
	"admiralty.io/multicluster-controller/pkg/reference"
	"k8s.io/apimachinery/pkg/api/errors"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"openmcp/openmcp/openmcp-resource-controller/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	extv1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

	fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPIngress{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	if err := co.WatchResourceReconcileController(live, &extv1b1.Ingress{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
	}
	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	for _, ghost := range ghosts {
		fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(ghost, &extv1b1.Ingress{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
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
	fmt.Println("********* [", i, "] *********")
	fmt.Println(req.Context, " / ", req.Namespace, " / ", req.Name)
	cm := clusterManager.NewClusterManager()

	// Fetch the OpenMCPDeployment instance
	instance := &ketiv1alpha1.OpenMCPIngress{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	fmt.Println("instance Name: ", instance.Name)
	fmt.Println("instance Namespace : ", instance.Namespace)

	if err != nil {
		if errors.IsNotFound(err) {
			// ...TODO: multicluster garbage collector
			// Until then...
			fmt.Println("Delete Deployments ..Cluster")
			err := DeleteIngress(cm, req.NamespacedName)
			return reconcile.Result{}, err
		}
		fmt.Println("Error1")
		return reconcile.Result{}, err
	}
	if instance.Status.ClusterMaps == nil || instance.Status.ChangeNeed == true {
		fmt.Println("Ingress Create Start")
		r.createIngress(req, cm, instance)

		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil

	} else {
		// Check Ingress In Openmcp
		fmt.Println("check Ingress In Openmcp")
		foundIngress := &extv1b1.Ingress{}
		err = cm.Host_client.Get(context.TODO(), foundIngress, instance.Namespace, instance.Name)

		if err == nil {
			foundService := &corev1.Service{}
			nsn := types.NamespacedName{
				Namespace: "openmcp",
				Name:      "loadbalancing-controller",
			}
			err = r.live.Get(context.TODO(), nsn, foundService)
			if err != nil && errors.IsNotFound(err) {
				fmt.Println("LoadBalancing-controller Service Not Found")
				return reconcile.Result{}, err
			} else {
				fmt.Println("Update Ingress Status")
				foundIngress.Status.LoadBalancer = foundService.Status.LoadBalancer
				err = cm.Host_client.UpdateStatus(context.Background(), foundIngress)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		} else if errors.IsNotFound(err) {
			fmt.Println("Create Ingress")
			host_ing, _ := r.ingressForOpenMCPIngress(req, instance)
			err = cm.Host_client.Create(context.Background(), host_ing)
			if err != nil {
				return reconcile.Result{}, err
			}
		}

		// Check Ingress in cluster
		for k, _ := range instance.Status.ClusterMaps {
			cluster_name := k
			//isExist := v
			found := &extv1b1.Ingress{}
			cluster_client := cm.Cluster_genClients[cluster_name]
			err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
			if err != nil && errors.IsNotFound(err) {
				// Delete Ingress Detected
				fmt.Println("Cluster '" + cluster_name + "' ReDeployed")
				_, ing := r.ingressForOpenMCPIngress(req, instance)
				err = cluster_client.Create(context.Background(), ing)
				if err != nil {
					return reconcile.Result{}, err
				}

			}
		}
	}

	return reconcile.Result{}, nil // err
}

func (r *reconciler) registerPdnsServer(ingress *extv1b1.Ingress) error {
	pdnsClient, err := mypdns.PdnsNewClient()
	if err != nil {
		fmt.Println(err)
		return err
	}

	found := &corev1.Service{}
	nsn := types.NamespacedName{
		Namespace: "openmcp",
		Name:      "loadbalancing-controller",
	}
	err = r.live.Get(context.TODO(), nsn, found)
	if err != nil && errors.IsNotFound(err) {
		fmt.Println("LoadBalancing-controller Service Not Found")
		return err
	} else {
		ip := found.Status.LoadBalancer.Ingress[0].IP

		for _, rule := range ingress.Spec.Rules {
			resourceRecordSet := zones.ResourceRecordSet{
				Name:       ".",
				Type:       "A",
				TTL:        300,
				ChangeType: 0,
				Records:    []zones.Record{{Content: ip, Disabled: false, SetPTR: false}},
				Comments:   nil,
			}
			pdnsClient.Zones().AddRecordSetToZone(context.TODO(), "localhost", rule.Host+".", resourceRecordSet)

		}

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil

}

//
//func insertDB(cm *ClusterManager, ingress *extv1b1.Ingress) {
//	fmt.Println("insertDB")
//	db, err := sql.Open("mysql", "root:ketilinux@tcp(10.0.3.12:3306)/powerdns")
//	if err != nil {
//		fmt.Println(err)
//	}
//	var domainID string
//	err = db.QueryRow("SELECT id FROM domain WHERE name=\"openmcp.org\"").Scan(&domainID)
//	if err != nil {
//		fmt.Println(err)
//	}
//	s := rand.NewSource(time.Now().UnixNano())
//	r := rand.New(s)
//
//	for _, rule := range ingress.Spec.Rules {
//		host := rule.Host
//		var exists int
//		err := db.QueryRow("SELECT id FROM records WHERE name='" + host + "'").Scan(&exists)
//
//		if err == nil && exists !=0{
//			fmt.Println("Host Exist")
//		} else if err != nil && exists == 0{
//			idCheck := false
//			//var domainID string
//			var id string
//			for idCheck == false {
//				randNum := r.Intn(9999999)
//				fmt.Println("*****Create Records Table ID*****")
//				id = strconv.Itoa(randNum)
//				err = db.QueryRow("SELECT id FROM records WHERE id=" + id).Scan(&exists)
//				if err == nil && exists != 0 {
//					fmt.Println("ID Exist")
//				} else if err != nil && exists == 0  {
//					idCheck = true
//				}
//			}
//			found := &corev1.Service{}
//			openmcp := cm.Host_client
//			err = openmcp.Get(context.TODO(), found, "openmcp", "loadbalancing-controller")
//			if err != nil && errors.IsNotFound(err) {
//				fmt.Println("LoadBalancing-controller Service Not Found")
//			} else {
//				ip := found.Status.LoadBalancer.Ingress[0].IP
//				queryValue := "(" + id + "," + domainID + ",'" + host + "','A','" + ip + "',300,0 ,NULL ,0,NULL,1)"
//				_, err := db.Exec("INSERT INTO records VALUES " + queryValue)
//				if err != nil {
//					fmt.Println(err)
//				}
//			}
//		}
//	}
//	defer db.Close()
//}

func (r *reconciler) createIngress(req reconcile.Request, cm *clusterManager.ClusterManager, instance *ketiv1alpha1.OpenMCPIngress) error {
	host_ing, ing := r.ingressForOpenMCPIngress(req, instance)

	found := &extv1b1.Ingress{}
	err := cm.Host_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
	if err != nil && errors.IsNotFound(err) {
		//insertDB(cm, host_ing)
		//r.registerPdnsServer(host_ing)
		err = cm.Host_client.Create(context.Background(), host_ing)
		if err != nil {
			return err
		}
	}

	serviceFound := &corev1.Service{}
	cluster_map := make(map[string]int32)
	for _, cluster := range cm.Cluster_list.Items {
		cluster_ing := &extv1b1.Ingress{}
		deepcopy.Copy(cluster_ing, &ing)

		isService := true
		found := &extv1b1.Ingress{}
		cluster_client := cm.Cluster_genClients[cluster.Name]
		err := cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)

		fmt.Println("**********************")
		fmt.Println(err)

		if err != nil && errors.IsNotFound(err) {
			for i, rule := range cluster_ing.Spec.Rules {
				for _, paths := range rule.HTTP.Paths {
					serviceName := paths.Backend.ServiceName
					fmt.Println("service name")
					fmt.Println(serviceName)
					serviceErr := cluster_client.Get(context.TODO(), serviceFound, instance.Namespace, serviceName)
					if serviceErr != nil && errors.IsNotFound(serviceErr) {
						fmt.Println("service not found")
						isService = false
					}
				}

				cluster_ing.Spec.Rules[i].Host = cluster.Name + "." + rule.Host

			}
			if isService == true {
				fmt.Println("Create Ingress Resource - ", cluster.Name)
				err = cluster_client.Create(context.Background(), cluster_ing)
				cluster_map[cluster.Name] = 1
				if err != nil {
					fmt.Println(cluster.Name, " - ", err)
				}
			}
		}
	}
	instance.Status.ClusterMaps = cluster_map

	err = r.live.Status().Update(context.TODO(), instance)
	return err

}
func (r *reconciler) ingressForOpenMCPIngress(req reconcile.Request, m *ketiv1alpha1.OpenMCPIngress) (*extv1b1.Ingress, *extv1b1.Ingress) {
	host_ing := &extv1b1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		// Spec: m.Spec.Template.Spec,
	}
	deepcopy.Copy(&host_ing.Spec, &m.Spec.Template.Spec)

	host_ing.Namespace = "openmcp"

	for i, _ := range host_ing.Spec.Rules {
		for j, _ := range host_ing.Spec.Rules[i].HTTP.Paths {
			host_ing.Spec.Rules[i].HTTP.Paths[j].Backend.ServiceName = "loadbalancing-controller"
			host_ing.Spec.Rules[i].HTTP.Paths[j].Backend.ServicePort.IntVal = 80
		}
	}

	reference.SetMulticlusterControllerReference(host_ing, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	ing := &extv1b1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		// Spec: m.Spec.Template.Spec,
	}
	deepcopy.Copy(&ing.Spec, &m.Spec.Template.Spec)

	reference.SetMulticlusterControllerReference(ing, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return host_ing, ing
}
func DeleteIngress(cm *clusterManager.ClusterManager, nsn types.NamespacedName) error {
	ing := &extv1b1.Ingress{}
	err := cm.Host_client.Get(context.Background(), ing, nsn.Namespace, nsn.Name)
	if err != nil && errors.IsNotFound(err) {
		// all good
		fmt.Println("Not Found")
	} else if err != nil && !errors.IsNotFound(err) {
		return err
	}
	fmt.Println("OpenMCP Delete Start")
	err = cm.Host_client.Delete(context.Background(), ing, nsn.Namespace, nsn.Name)
	if err != nil {
		return err
	}
	fmt.Println("OpenMCP Delete Complate")

	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_genClients[cluster.Name]
		fmt.Println(nsn.Namespace, nsn.Name)
		err := cluster_client.Get(context.Background(), ing, nsn.Namespace, nsn.Name)
		if err != nil && errors.IsNotFound(err) {
			// all good
			fmt.Println("Not Found")
			continue
		}
		fmt.Println(cluster.Name, " Delete Start")
		err = cluster_client.Delete(context.Background(), ing, nsn.Namespace, nsn.Name)
		if err != nil {
			return err
		}
		fmt.Println(cluster.Name, "Delete Complate")
	}
	return nil

}

/*
func(cm *ClusterManager) Scheduling(replicas int32) map[string]int32{
        rand.Seed(time.Now().UTC().UnixNano())

        cluster_replicas_map := make(map[string]int32)

        remain_rep := replicas
        rep := 0
        cluster_len := len(cm.Cluster_list.Items)
        for i, cluster := range cm.Cluster_list.Items {
                if i == cluster_len-1 {
                        rep  = int(remain_rep)
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
*/
/*
func(cm *ClusterManager) ReScheduling(spec_replicas int32, status_replicas int32, status_cluster_replicas_map map[string]int32) map[string]int32{
	rand.Seed(time.Now().UTC().UnixNano())

	result_cluster_replicas_map := make(map[string]int32)
	for k,v := range status_cluster_replicas_map {
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
*/
/*
func keyOf(my_map map[string]int32, target_index int) string {
	index := 0
	for k, _ := range my_map {
		if index == target_index{
			return k
		}
		index += 1
        }
	return ""

}
*/
