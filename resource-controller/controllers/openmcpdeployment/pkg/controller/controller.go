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

package openmcpdeployment // import "admiralty.io/multicluster-controller/examples/openmcpdeployment/pkg/controller/openmcpdeployment"

import (
	"admiralty.io/multicluster-controller/pkg/reference"
	"context"
	"encoding/json"
	"fmt"
	syncapis "sync-controller/pkg/apis"

	//"github.com/getlantern/deepcopy"
	"k8s.io/apimachinery/pkg/api/errors"
	"math/rand"
	"reflect"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	"strconv"
	"strings"
	//"reflect"
	"sort"
	"time"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	"resource-controller/apis"
	ketiv1alpha1 "resource-controller/apis/keti/v1alpha1"
	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	sync "sync-controller/pkg/apis/keti/v1alpha1"
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
	if err := syncapis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPDeployment{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	for _, ghost := range ghosts {
		fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(ghost, &appsv1.Deployment{}, controller.WatchOptions{}); err != nil {
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
var syncIndex int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	i += 1
	fmt.Println("********* [", i, "] *********")
	fmt.Println(req.Context, " / ", req.Namespace, " / ", req.Name)
	cm := NewClusterManager()
	fmt.Println("cm check!!!", cm)

	// Fetch the OpenMCPDeployment instance
	instance := &ketiv1alpha1.OpenMCPDeployment{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	fmt.Println("instance Name: ", instance.Name)
	fmt.Println("instance Namespace : ", instance.Namespace)

	if err != nil {
		fmt.Println("1. Delete Detection")
		if errors.IsNotFound(err) {
			// ...TODO: multicluster garbage collector
			// Until then...
			fmt.Println("Delete Deployments ..Cluster")
			//err := cm.DeleteDeployments(req.NamespacedName)
			err := r.DeleteDeploys(cm, instance.Name, instance.Namespace)

			r.ServiceNotifyAll(req.Namespace)

			return reconcile.Result{}, err
		}
		fmt.Println("check2", err)
		return reconcile.Result{}, err
	}
	if instance.Status.ClusterMaps == nil {
		fmt.Println("2. Create Detection")
		if instance.Status.SchedulingNeed == false && instance.Status.SchedulingComplete == false {
			instance.Status.SchedulingNeed = true

			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				fmt.Println("Failed to update instance status", err)
				fmt.Println("check16", err)
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, err

			//} else if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false {
		} else if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false && strings.Compare(instance.Spec.Labels["test"], "yes") != 0 {
			//temp
			fmt.Println("Scheduling Start")
			replicas := instance.Spec.Replicas

			instance.Status.ClusterMaps = cm.Scheduling(replicas)
			instance.Status.Replicas = replicas

			instance.Status.SchedulingNeed = false
			instance.Status.SchedulingComplete = true

			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				fmt.Println("Failed to update instance status", err)
				fmt.Println("check4", err)
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, err

		} else if instance.Status.SchedulingNeed == false && instance.Status.SchedulingComplete == true {
			fmt.Println("Create Deployments")

			for _, cluster := range cm.Cluster_list.Items {

				if instance.Status.ClusterMaps[cluster.Name] == 0 {
					continue
				}
				found := &appsv1.Deployment{}
				cluster_client := cm.Cluster_clients[cluster.Name]

				err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
				if err != nil && errors.IsNotFound(err) {
					// TODO: Today

					replica := instance.Status.ClusterMaps[cluster.Name]
					fmt.Println("Cluster '"+cluster.Name+"' Deployed (", replica, " / ", instance.Status.Replicas, ")")
					dep := r.deploymentForOpenMCPDeployment(req, instance, replica)
					command := "create"
					err = r.sendSync(cm, dep, command, cluster.Name)
					//err = cluster_client.Create(context.Background(), dep)
					if err != nil {
						fmt.Println("check3", err)
						return reconcile.Result{}, err
					}
				}
			}
		}
		r.ServiceNotify(instance.Spec.Labels, instance.Namespace)

		instance.Status.LastSpec = instance.Spec
		//instance.Status.LastUpdateTime = time.Now().Format(time.RFC3339)

		err := r.live.Status().Update(context.TODO(), instance)
		if err != nil {
			fmt.Println("Failed to update instance status", err)
			fmt.Println("check4", err)
			return reconcile.Result{}, err
		}

		fmt.Println("check5", err)
		return reconcile.Result{}, nil

	}
	if !reflect.DeepEqual(instance.Status.LastSpec, instance.Spec) {

		fmt.Println("3. Update Detection")
		if instance.Status.Replicas != instance.Spec.Replicas {
			fmt.Println("Change Spec Replicas ! ReScheduling Start & Update Deployment")
			cluster_replicas_map := cm.ReScheduling(instance.Spec.Replicas, instance.Status.Replicas, instance.Status.ClusterMaps)

			for _, cluster := range cm.Cluster_list.Items {
				update_replica := cluster_replicas_map[cluster.Name]
				cluster_client := cm.Cluster_clients[cluster.Name]

				dep := r.deploymentForOpenMCPDeployment(req, instance, update_replica)

				found := &appsv1.Deployment{}
				err := cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
				if err != nil && errors.IsNotFound(err) {
					// Not Exist Deployment.
					if update_replica != 0 {
						// Create !
						command := "create"
						err = r.sendSync(cm, dep, command, cluster.Name)
						//err = cluster_client.Create(context.Background(), dep)
						if err != nil {
							fmt.Println("check6", err)
							return reconcile.Result{}, err
						}
					}

				} else if err != nil {
					fmt.Println("check7", err)
					return reconcile.Result{}, err
				} else {
					// Already Exist Deployment.
					if update_replica == 0 {
						// Delete !
						//dep := &appsv1.Deployment{}
						command := "delete"
						err = r.sendSync(cm, dep, command, cluster.Name)
						//err = cluster_client.Delete(context.Background(), dep, req.Namespace, req.Name)

						if err != nil {
							fmt.Println("check8", err)
							return reconcile.Result{}, err
						}
					} else {
						// Update !
						command := "update"
						err = r.sendSync(cm, dep, command, cluster.Name)
						//err = cluster_client.Update(context.TODO(), dep)
						if err != nil {
							fmt.Println("check9", err)
							return reconcile.Result{}, err
						}

					}

				}

			}
			r.ServiceNotify(instance.Spec.Labels, instance.Namespace)

			instance.Status.ClusterMaps = cluster_replicas_map
			instance.Status.Replicas = instance.Spec.Replicas
			instance.Status.LastSpec = instance.Spec
			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				fmt.Println("Failed to update instance status", err)
				fmt.Println("check10", err)
				return reconcile.Result{}, err
			}

		}
		if !reflect.DeepEqual(instance.Status.LastSpec.Labels, instance.Spec.Labels) {
			last_label := instance.Status.LastSpec.Labels
			current_label := instance.Spec.Labels

			r.ServiceNotify(last_label, instance.Namespace)
			r.ServiceNotify(current_label, instance.Namespace)
		}

		instance.Status.LastSpec = instance.Spec
		err := r.live.Status().Update(context.TODO(), instance)
		if err != nil {
			fmt.Println("Failed to update instance status", err)
			fmt.Println("check10", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, err
	}

	// Check Deployment in cluster
	for k, v := range instance.Status.ClusterMaps {
		fmt.Println("4. Check Clusters")
		cluster_name := k
		replica := v

		if replica == 0 {
			continue
		}
		found := &appsv1.Deployment{}
		cluster_client := cm.Cluster_clients[cluster_name]
		err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)

		if err != nil && errors.IsNotFound(err) {
			// Delete Deployment Detected
			fmt.Println("Cluster '"+cluster_name+"' ReDeployed => ", replica)
			dep := r.deploymentForOpenMCPDeployment(req, instance, replica)
			command := "create"
			err = r.sendSync(cm, dep, command, cluster_name)
			//err = cluster_client.Create(context.Background(), dep)
			if err != nil {
				fmt.Println("check11", err)
				return reconcile.Result{}, err
			}

		}

	}

	fmt.Println("check12", err)
	return reconcile.Result{}, nil // err
}
func (r *reconciler) DeleteDeploys(cm *ClusterManager, name string, namespace string) error {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	for _, cluster := range cm.Cluster_list.Items {
		command := "delete"
		err := r.sendSync(cm, dep, command, cluster.Name)
		if err != nil {
			return err
		}

	}
	return nil

}
func (r *reconciler) sendSync(cm *ClusterManager, dep *appsv1.Deployment, command string, clusterName string) error {
	syncIndex += 1

	s := &sync.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-deployment-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: sync.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *dep,
		},
	}

	err := r.live.Create(context.TODO(), s)

	return err

}
func (r *reconciler) ServiceNotify(label_map map[string]string, namespace string) error {
	fmt.Println("ServiceNotify Called", label_map)

	osvc_list := &ketiv1alpha1.OpenMCPServiceList{}
	listOptions := &client.ListOptions{Namespace: namespace}

	//labelSelector := metav1.LabelSelector{MatchLabels: label_map}
	//label_selector, err := metav1.LabelSelectorAsSelector(&labelSelector)
	//if err != nil {
	//	return err
	//}
	//listOptions := &client.ListOptions{
	//	Namespace: namespace,
	//	LabelSelector: label_selector,
	//}

	r.live.List(context.TODO(), listOptions, osvc_list)
	for _, osvc := range osvc_list.Items {
		for k, v := range osvc.Spec.LabelSelector {
			fmt.Println(k, " / ", v)
			if label_map[k] == v {
				fmt.Println("->", osvc.Name, " notify !")
				osvc.Status.ChangeNeed = true
				err := r.live.Status().Update(context.TODO(), &osvc)
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}
func (r *reconciler) ServiceNotifyAll(namespace string) error {
	fmt.Println("ServiceNotifyAll Called")

	osvc_list := &ketiv1alpha1.OpenMCPServiceList{}
	listOptions := &client.ListOptions{Namespace: namespace}

	//labelSelector := metav1.LabelSelector{MatchLabels: label_map}
	//label_selector, err := metav1.LabelSelectorAsSelector(&labelSelector)
	//if err != nil {
	//	return err
	//}
	//listOptions := &client.ListOptions{
	//	Namespace: namespace,
	//	LabelSelector: label_selector,
	//}

	r.live.List(context.TODO(), listOptions, osvc_list)
	for _, osvc := range osvc_list.Items {

		fmt.Println("->", osvc.Name, " notify !")
		osvc.Status.ChangeNeed = true
		err := r.live.Status().Update(context.TODO(), &osvc)
		if err != nil {
			return err
		}
	}

	return nil
}
func (r *reconciler) deploymentForOpenMCPDeployment(req reconcile.Request, m *ketiv1alpha1.OpenMCPDeployment, replica int32) *appsv1.Deployment {
	fmt.Println("[CHECK] deploymentForOpenMCPDeployment")
	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec:   openmcpDeploymentTemplateSpecToDeploymentSpec(m.Spec.Template.Spec),
		Status: appsv1.DeploymentStatus{},
	}

	if dep.Spec.Selector == nil {
		dep.Spec.Selector = &metav1.LabelSelector{}
	}

	dep.Spec.Selector.MatchLabels = m.Spec.Labels
	dep.Spec.Template.ObjectMeta.Labels = m.Spec.Labels
	dep.Spec.Replicas = &replica

	reference.SetMulticlusterControllerReference(dep, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return dep
}
func (cm *ClusterManager) DeleteDeployments(nsn types.NamespacedName) error {
	dep := &appsv1.Deployment{}
	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_clients[cluster.Name]
		err := cluster_client.Get(context.Background(), dep, nsn.Namespace, nsn.Name)
		if err != nil && errors.IsNotFound(err) {
			// all good
			fmt.Println("Not Found")
			continue
		}
		if !isInObject(dep, "OpenMCPDeployment") {
			continue
		}
		fmt.Println(cluster.Name, " Delete Start")
		err = cluster_client.Delete(context.Background(), dep, nsn.Namespace, nsn.Name)
		if err != nil {
			return err
		}
		fmt.Println(cluster.Name, "Delete Complete")
	}
	return nil

}

func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		fmt.Println("Error retrieving list of federated clusters: ", err)
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
	namespace := "kube-federation-system"
	cluster_len := len(cm.Cluster_list.Items)
	for i, cluster := range cm.Cluster_list.Items {
		except := false
		joined_cluster := &fedv1b1.KubeFedCluster{}
		err := cm.Host_client.Get(context.TODO(), joined_cluster, namespace, cluster.Name)
		if err != nil {
			return nil
		}
		for k, v := range joined_cluster.Labels {
			if k == "openmcp" && v == "true" {
				fmt.Println("Scheduling Except Cluster !! Include OpenMCP Label : ", k, v)
				except = true
				break
			}
		}
		if except {
			continue
		}

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

func isInObject(child *appsv1.Deployment, parent string) bool {
	refKind_str := child.ObjectMeta.Annotations["multicluster.admiralty.io/controller-reference"]
	refKind_map := make(map[string]interface{})
	err := json.Unmarshal([]byte(refKind_str), &refKind_map)
	if err != nil {
		panic(err)
	}
	if refKind_map["kind"] == parent {
		return true
	}
	return false
}

func openmcpContainersToContainers(containers []ketiv1alpha1.OpenMCPContainer) []corev1.Container {
	var newContainers []corev1.Container

	for _, container := range containers {
		newContainer := corev1.Container{
			Name:       container.Name,
			Image:      container.Image,
			Command:    container.Command,
			Args:       container.Args,
			WorkingDir: container.WorkingDir,
			Ports:      container.Ports,
			EnvFrom:    container.EnvFrom,
			Env:        container.Env,
			Resources: corev1.ResourceRequirements{
				Limits:   container.Resources.Limits,
				Requests: container.Resources.Requests,
			},
			VolumeMounts:             container.VolumeMounts,
			VolumeDevices:            container.VolumeDevices,
			LivenessProbe:            container.LivenessProbe,
			ReadinessProbe:           container.ReadinessProbe,
			Lifecycle:                container.Lifecycle,
			TerminationMessagePath:   container.TerminationMessagePath,
			TerminationMessagePolicy: container.TerminationMessagePolicy,
			ImagePullPolicy:          container.ImagePullPolicy,
			SecurityContext:          container.SecurityContext,
			Stdin:                    container.Stdin,
			StdinOnce:                container.StdinOnce,
			TTY:                      container.TTY,
		}
		newContainers = append(newContainers, newContainer)
	}

	return newContainers
}

func openmcpPodSpecToPodSpec(spec ketiv1alpha1.OpenMCPPodSpec) corev1.PodSpec {
	return corev1.PodSpec{
		Volumes:                       spec.Volumes,
		InitContainers:                openmcpContainersToContainers(spec.InitContainers),
		Containers:                    openmcpContainersToContainers(spec.Containers),
		RestartPolicy:                 spec.RestartPolicy,
		TerminationGracePeriodSeconds: spec.TerminationGracePeriodSeconds,
		ActiveDeadlineSeconds:         spec.ActiveDeadlineSeconds,
		DNSPolicy:                     spec.DNSPolicy,
		NodeSelector:                  spec.NodeSelector,
		ServiceAccountName:            spec.ServiceAccountName,
		DeprecatedServiceAccount:      spec.DeprecatedServiceAccount,
		AutomountServiceAccountToken:  spec.AutomountServiceAccountToken,
		NodeName:                      spec.NodeName,
		HostNetwork:                   spec.HostNetwork,
		HostPID:                       spec.HostPID,
		HostIPC:                       spec.HostIPC,
		ShareProcessNamespace:         spec.ShareProcessNamespace,
		SecurityContext:               spec.SecurityContext,
		ImagePullSecrets:              spec.ImagePullSecrets,
		Hostname:                      spec.Hostname,
		Subdomain:                     spec.Subdomain,
		Affinity:                      spec.Affinity,
		SchedulerName:                 spec.SchedulerName,
		Tolerations:                   spec.Tolerations,
		HostAliases:                   spec.HostAliases,
		PriorityClassName:             spec.PriorityClassName,
		Priority:                      spec.Priority,
		DNSConfig:                     spec.DNSConfig,
		ReadinessGates:                spec.ReadinessGates,
		RuntimeClassName:              spec.RuntimeClassName,
		EnableServiceLinks:            spec.EnableServiceLinks,
	}
}

func openmcpPodTemplateSpecToPodTemplateSpec(template ketiv1alpha1.OpenMCPPodTemplateSpec) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: template.ObjectMeta,
		Spec:       openmcpPodSpecToPodSpec(template.Spec),
	}
}

func openmcpDeploymentTemplateSpecToDeploymentSpec(templateSpec ketiv1alpha1.OpenMCPDeploymentTemplateSpec) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas:                templateSpec.Replicas,
		Selector:                templateSpec.Selector,
		Template:                openmcpPodTemplateSpecToPodTemplateSpec(templateSpec.Template),
		Strategy:                templateSpec.Strategy,
		MinReadySeconds:         templateSpec.MinReadySeconds,
		RevisionHistoryLimit:    templateSpec.RevisionHistoryLimit,
		Paused:                  templateSpec.Paused,
		ProgressDeadlineSeconds: templateSpec.ProgressDeadlineSeconds,
	}
}
