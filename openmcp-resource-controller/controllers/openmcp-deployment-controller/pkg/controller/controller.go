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
	"k8s.io/klog"
	"openmcp/openmcp/util/clusterManager"

	syncapis "openmcp/openmcp/openmcp-sync-controller/pkg/apis"

	//"github.com/getlantern/deepcopy"
	"k8s.io/apimachinery/pkg/api/errors"
	"math/rand"
	"reflect"

	"strconv"
	"strings"
	//"reflect"
	"sort"
	"time"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	"openmcp/openmcp/openmcp-resource-controller/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"

	sync "openmcp/openmcp/openmcp-sync-controller/pkg/apis/keti/v1alpha1"
)

var cm *clusterManager.ClusterManager
func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	klog.V(4).Info("[OpenMCP Deployment] Function Called NewController")
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
	if err := syncapis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	//fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPDeployment{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	for _, ghost := range ghosts {
		//fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
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
	klog.V(4).Info("[OpenMCP Deployment] Function Called Reconcile")
	i += 1

	// Fetch the OpenMCPDeployment instance
	instance := &ketiv1alpha1.OpenMCPDeployment{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	klog.V(0).Info("Resource Get => [Name] : "+ instance.Name + " [Namespace]  : " + instance.Namespace)

	if err != nil {

		if errors.IsNotFound(err) {
			klog.V(0).Info("[Delete Detect]")
			// ...TODO: multicluster garbage collector
			// Until then...
			klog.V(0).Info("Delete Deployment of All Cluster")
			err := r.DeleteDeploys(cm, req.NamespacedName.Name, req.NamespacedName.Namespace)

			klog.V(0).Info("Service Notify Send")
			r.ServiceNotifyAll(req.Namespace)

			return reconcile.Result{}, err
		}
		return reconcile.Result{}, err
	}

	if instance.Status.CreateSyncRequestComplete == false {
		klog.V(0).Info("[Create Detect]")
		klog.V(0).Info("Create Deployment Start")
		klog.V(0).Info("SchedulingNeed : ", instance.Status.SchedulingNeed, ", SchedulingComplete : ", instance.Status.SchedulingComplete)
		if instance.Status.SchedulingNeed == false && instance.Status.SchedulingComplete == false {
			klog.V(0).Info("Scheduling 요청 (SchedulingNeed false => true)")
			instance.Status.SchedulingNeed = true

			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				klog.V(0).Info("Failed to update instance status", err)
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, err

			//} else if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false {
		} else if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false && strings.Compare(instance.Spec.Labels["test"], "yes") != 0 {
			//temp
			klog.V(0).Info("Local Scheduling을 시작합니다.(랜덤 스케줄링)")
			klog.V(0).Info("Scheduling Controller와 연계하려면 Labels의 test항목을 no로 변경해주세요")
			replicas := instance.Spec.Replicas

			instance.Status.ClusterMaps = Scheduling(cm, replicas)
			instance.Status.Replicas = replicas

			instance.Status.SchedulingNeed = false
			instance.Status.SchedulingComplete = true
			klog.V(0).Info("Scheduling 완료")
			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				klog.V(0).Info("Failed to update instance status", err)
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, err

		} else if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false && strings.Compare(instance.Spec.Labels["test"], "yes") == 0{
			klog.V(0).Info("Scheduling Wait")
			fmt.Println("testestsetsetseteststsetestsetestttttttttttttttttttttttttttttttt")
			return reconcile.Result{}, err

		} else if instance.Status.SchedulingNeed == false && instance.Status.SchedulingComplete == true {
			klog.V(0).Info("Scheduling 결과를 통해 Deployment의 Sync Resource를 생성합니다.")

			sync_req_name := instance.Status.SyncRequestName

			for _, myCluster := range cm.Cluster_list.Items {
				if instance.Status.ClusterMaps[myCluster.Name] == 0 {
					continue
				}
				found := &appsv1.Deployment{}
				cluster_client := cm.Cluster_genClients[myCluster.Name]

				err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
				if err != nil && errors.IsNotFound(err) {
					// TODO: Today

					replica := instance.Status.ClusterMaps[myCluster.Name]

					dep := r.deploymentForOpenMCPDeployment(req, instance, replica)
					command := "create"
					klog.V(0).Info("SyncResource Create (ClusterName : "+myCluster.Name+", Command : "+ command+", Replicas :", replica, " / ", instance.Status.Replicas, ")")
					sync_req_name, err = r.sendSync(dep, command, myCluster.Name)
					//err = cluster_client.Create(context.Background(), dep)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}
			klog.V(0).Info("Service Notify Send")
			r.ServiceNotify(instance.Spec.Labels, instance.Namespace)

			instance.Status.LastSpec = instance.Spec
			instance.Status.CreateSyncRequestComplete = true
			instance.Status.SyncRequestName = sync_req_name
			klog.V(0).Info("sync_req_name : ", sync_req_name)

			//instance.Status.LastUpdateTime = time.Now().Format(time.RFC3339)
			klog.V(0).Info("Update Status")
			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				klog.V(0).Info("Failed to update instance status", err)
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}

	}

	if !reflect.DeepEqual(instance.Status.LastSpec, instance.Spec) {

		klog.V(0).Info("[Update Detection]")
		sync_req_name := instance.Status.SyncRequestName
		if instance.Status.Replicas != instance.Spec.Replicas {
			klog.V(0).Info("Change Spec Replicas ! ReScheduling Start & Update Deployment")
			cluster_replicas_map := ReScheduling(instance.Spec.Replicas, instance.Status.Replicas, instance.Status.ClusterMaps)

			for _, cluster := range cm.Cluster_list.Items {
				update_replica := cluster_replicas_map[cluster.Name]
				cluster_client := cm.Cluster_genClients[cluster.Name]

				dep := r.deploymentForOpenMCPDeployment(req, instance, update_replica)

				found := &appsv1.Deployment{}
				err := cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
				if err != nil && errors.IsNotFound(err) {
					// Not Exist Deployment.
					if update_replica != 0 {
						// Create !
						command := "create"
						klog.V(0).Info("SyncResource Create (ClusterName : "+cluster.Name+", Command : "+ command+", Replicas :", update_replica, " / ", instance.Status.Replicas, ")")
						sync_req_name, err = r.sendSync(dep, command, cluster.Name)
						//err = cluster_client.Create(context.Background(), dep)
						if err != nil {
							return reconcile.Result{}, err
						}
					}

				} else if err != nil {
					return reconcile.Result{}, err
				} else {
					// Already Exist Deployment.
					if update_replica == 0 {
						// Delete !
						//dep := &appsv1.Deployment{}
						command := "delete"
						klog.V(0).Info("SyncResource Create (ClusterName : "+cluster.Name+", Command : "+ command+", Replicas :", update_replica, " / ", instance.Status.Replicas, ")")
						sync_req_name, err = r.sendSync(dep, command, cluster.Name)

						//err = cluster_client.Delete(context.Background(), dep, req.Namespace, req.Name)

						if err != nil {
							return reconcile.Result{}, err
						}
					} else {
						// Update !
						command := "update"
						klog.V(0).Info("SyncResource Create (ClusterName : "+cluster.Name+", Command : "+ command+", Replicas :", update_replica, " / ", instance.Status.Replicas, ")")
						sync_req_name, err = r.sendSync(dep, command, cluster.Name)
						//err = cluster_client.Update(context.TODO(), dep)
						if err != nil {
							return reconcile.Result{}, err
						}

					}

				}

			}
			r.ServiceNotify(instance.Spec.Labels, instance.Namespace)

			instance.Status.ClusterMaps = cluster_replicas_map
			instance.Status.Replicas = instance.Spec.Replicas
			instance.Status.LastSpec = instance.Spec

			//err := r.live.Status().Update(context.TODO(), instance)
			//if err != nil {
			//	klog.V(0).Info("Failed to update instance status", err)
			//	klog.V(0).Info("check10", err)
			//	return reconcile.Result{}, err
			//}

		}
		if !reflect.DeepEqual(instance.Status.LastSpec.Labels, instance.Spec.Labels) {
			last_label := instance.Status.LastSpec.Labels
			current_label := instance.Spec.Labels
			klog.V(0).Info("Label Changed")
			klog.V(0).Info("Service Notify")
			r.ServiceNotify(last_label, instance.Namespace)
			r.ServiceNotify(current_label, instance.Namespace)
		}

		instance.Status.LastSpec = instance.Spec
		instance.Status.SyncRequestName = sync_req_name
		klog.V(0).Info("sync_req_name : ", sync_req_name)
		klog.V(0).Info("Status Update")
		err := r.live.Status().Update(context.TODO(), instance)
		if err != nil {
			klog.V(0).Info("Failed to update instance status", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, err
	}

	sync_instance := &sync.Sync{}
	nsn := types.NamespacedName{
		"openmcp",
		instance.Status.SyncRequestName,
	}
	err = r.live.Get(context.TODO(), nsn, sync_instance)
	if err == nil {
		// 아직 Sync에서 처리되지 않음
		return reconcile.Result{}, nil
	}

	// Check Deployment in cluster
	klog.V(0).Info("[Member Cluster Check Deployment]")
	sync_req_name := instance.Status.SyncRequestName
	for k, v := range instance.Status.ClusterMaps {
		cluster_name := k
		replica := v

		if replica == 0 {
			continue
		}
		found := &appsv1.Deployment{}
		cluster_client := cm.Cluster_genClients[cluster_name]
		err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)

		if err != nil && errors.IsNotFound(err) {
			// Delete Deployment Detected
			klog.V(0).Info("Cluster '"+cluster_name+"' ReDeployed => ", replica)
			dep := r.deploymentForOpenMCPDeployment(req, instance, replica)
			command := "create"
			klog.V(0).Info("SyncResource Create (ClusterName : "+cluster_name+", Command : "+ command+", Replicas :", replica, " / ", instance.Status.Replicas, ")")
			sync_req_name, err = r.sendSync(dep, command, cluster_name)
			//err = cluster_client.Create(context.Background(), dep)
			if err != nil {
				return reconcile.Result{}, err
			}

		}

	}
	instance.Status.SyncRequestName = sync_req_name
	klog.V(0).Info("sync_req_name : ", sync_req_name)

	err = r.live.Status().Update(context.TODO(), instance)
	if err != nil {
		klog.V(0).Info("Failed to update instance status", err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil // err
}
func (r *reconciler) DeleteDeploys(cm *clusterManager.ClusterManager, name string, namespace string) error {
	klog.V(4).Info("[OpenMCP Deployment] Function Called DeleteDeploys")

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{},
	}
	klog.V(0).Info("Delete Check ", dep.Name, dep.Namespace)
	for _, cluster := range cm.Cluster_list.Items {
		command := "delete"
		_, err := r.sendSync(dep, command, cluster.Name)
		if err != nil {
			return err
		}

	}
	return nil

}
func (r *reconciler) sendSync(dep *appsv1.Deployment, command string, clusterName string) (string, error) {
	klog.V(4).Info("[OpenMCP Deployment] Function Called sendSync")
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
	klog.V(0).Info("Delete Check2 ", s.Spec.Template.(appsv1.Deployment).Name, s.Spec.Template.(appsv1.Deployment).Namespace)

	err := r.live.Create(context.TODO(), s)
	klog.V(0).Info(s.Name)
	return s.Name, err

}
func (r *reconciler) ServiceNotify(label_map map[string]string, namespace string) error {
	klog.V(4).Info("[OpenMCP Deployment] Function Called ServiceNotify")
	klog.V(4).Info("[OpenMCP Deployment] label_map : ", label_map)

	osvc_list := &ketiv1alpha1.OpenMCPServiceList{}
	listOptions := &client.ListOptions{Namespace: namespace}

	klog.V(4).Info("[OpenMCP Deployment] find Notify Target Service")
	r.live.List(context.TODO(), osvc_list, listOptions)
	for _, osvc := range osvc_list.Items {
		for k, v := range osvc.Spec.LabelSelector {
			klog.V(4).Info("[OpenMCP Deployment] find Notify Target Label : ", k, " / ", v)
			if label_map[k] == v {
				klog.V(4).Info("[OpenMCP Deployment] Service '", osvc.Name, "' Will Notify!")
				osvc.Status.ChangeNeed = true
				err := r.live.Status().Update(context.TODO(), &osvc)
				if err != nil {
					return err
				}
				klog.V(4).Info("[OpenMCP Deployment] Service '", osvc.Name, "' Notify Success!")

			}
		}
	}

	return nil
}
func (r *reconciler) ServiceNotifyAll(namespace string) error {
	klog.V(4).Info("[OpenMCP Deployment] Function Called ServiceNotifyAll")

	osvc_list := &ketiv1alpha1.OpenMCPServiceList{}
	listOptions := &client.ListOptions{Namespace: namespace}

	r.live.List(context.TODO(), osvc_list, listOptions)
	for _, osvc := range osvc_list.Items {

		klog.V(0).Info("->", osvc.Name, " notify !")
		osvc.Status.ChangeNeed = true
		err := r.live.Status().Update(context.TODO(), &osvc)
		if err != nil {
			return err
		}
	}

	return nil
}
func (r *reconciler) deploymentForOpenMCPDeployment(req reconcile.Request, m *ketiv1alpha1.OpenMCPDeployment, replica int32) *appsv1.Deployment {
	klog.V(0).Info("[CHECK] deploymentForOpenMCPDeployment")
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

func DeleteDeployments(cm *clusterManager.ClusterManager, nsn types.NamespacedName) error {
	dep := &appsv1.Deployment{}
	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_genClients[cluster.Name]
		err := cluster_client.Get(context.Background(), dep, nsn.Namespace, nsn.Name)
		if err != nil && errors.IsNotFound(err) {
			// all good
			klog.V(0).Info("Not Found")
			continue
		}
		if !isInObject(dep, "OpenMCPDeployment") {
			continue
		}
		klog.V(0).Info(cluster.Name, " Delete Start")
		err = cluster_client.Delete(context.Background(), dep, nsn.Namespace, nsn.Name)
		if err != nil {
			return err
		}
		klog.V(0).Info(cluster.Name, "Delete Complete")
	}
	return nil

}

func Scheduling(cm *clusterManager.ClusterManager, replicas int32) map[string]int32 {
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
				klog.V(0).Info("Scheduling Except Cluster !! Include OpenMCP Label : ", k, v)
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

	klog.V(0).Info("Scheduling Result: ")
	for _, k := range keys {
		v := cluster_replicas_map[k]
		klog.V(0).Info("  ", k, ": ", v)
	}
	return cluster_replicas_map

}
func ReScheduling(spec_replicas int32, status_replicas int32, status_cluster_replicas_map map[string]int32) map[string]int32 {
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
		klog.V(0).Info("cluster_len : ", cluster_len)
		selected_cluster_target_index := rand.Intn(cluster_len)

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

	klog.V(0).Info("ReScheduling Result: ")
	for _, k := range keys {
		v := result_cluster_replicas_map[k]
		prev_v := status_cluster_replicas_map[k]
		klog.V(0).Info("  ", k, ": ", prev_v, " -> ", v)
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
