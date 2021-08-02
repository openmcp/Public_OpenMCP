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

package controller // import "admiralty.io/multicluster-controller/examples/openmcpdeployment/pkg/controller/openmcpdeployment"

import (
	"context"
	"fmt"
	"openmcp/openmcp/apis"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	syncv1alpha1 "openmcp/openmcp/apis/sync/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"reflect"
	"strconv"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"admiralty.io/multicluster-controller/pkg/reference"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("[OpenMCP Deployment] Function Called NewController")
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

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &resourcev1alpha1.OpenMCPDeployment{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	for _, ghost := range ghosts {
		if err := co.WatchResourceReconcileController(context.TODO(), ghost, &appsv1.Deployment{}, controller.WatchOptions{}); err != nil {
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
	/*org_cm.Mutex.Lock()
	var cm *clusterManager.ClusterManager
	copier.Copy(cm, org_cm)
	org_cm.Mutex.Unlock()

	omcplog.V(4).Info("CM : ", len(cm.Cluster_list.Items))
	*/

	omcplog.V(4).Info("[OpenMCP Deployment] Function Called Reconcile")

	i += 1

	// Fetch the OpenMCPDeployment instance
	instance := &resourcev1alpha1.OpenMCPDeployment{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	omcplog.V(5).Info("Resource Get => [Name] : " + instance.Name + " [Namespace]  : " + instance.Namespace)

	if err != nil {

		if errors.IsNotFound(err) {
			omcplog.V(2).Info("[Delete Detect]")
			omcplog.V(2).Info("Delete Deployment of All Cluster")
			err := r.DeleteDeploys(cm, req.NamespacedName.Name, req.NamespacedName.Namespace)

			omcplog.V(2).Info("Service Notify Send")
			//r.ServiceNotifyAll(req.Namespace)

			return reconcile.Result{}, err
		}
		return reconcile.Result{}, err
	}

	if instance.Status.CreateSyncRequestComplete == false {
		omcplog.V(2).Info("[Create Detect]")
		omcplog.V(2).Info("Create Deployment Start")
		omcplog.V(3).Info("SchedulingNeed : ", instance.Status.SchedulingNeed, ", SchedulingComplete : ", instance.Status.SchedulingComplete)

		if instance.Status.SchedulingNeed == false && instance.Status.SchedulingComplete == false {
			omcplog.V(3).Info("Scheduling 요청 (SchedulingNeed false => true)")
			instance.Status.SchedulingNeed = true

			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				omcplog.V(0).Info("Failed to update instance status", err)
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, err

			//} else if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false {
		} else if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false {
			omcplog.V(2).Info("Scheduling Wait")
			return reconcile.Result{}, nil

		} else if instance.Status.SchedulingNeed == false && instance.Status.SchedulingComplete == true {

			omcplog.V(2).Info("Create a Sync Resource for Deployment with Scheduling results.")

			sync_req_name := ""

			omcplog.V(5).Info("Cluster Count : ", len(cm.Cluster_list.Items))
			for _, myCluster := range cm.Cluster_list.Items {
				replica := instance.Status.ClusterMaps[myCluster.Name]
				cluster_client := cm.Cluster_genClients[myCluster.Name]

				dep := r.deploymentForOpenMCPDeployment(req, instance, replica, myCluster.Name)

				found := &appsv1.Deployment{}
				err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
				omcplog.V(2).Info("/// cluster_client : ", cluster_client)

				if err != nil && errors.IsNotFound(err) {
					// Not Exist Deployment.
					if replica != 0 {
						// Create !
						command := "create"
						omcplog.V(2).Info("SyncResource Create (ClusterName : "+myCluster.Name+", Command : "+command+", Replicas :", replica, " / ", instance.Status.Replicas, ")")
						sync_req_name, err = r.sendSync(dep, command, myCluster.Name)

						if err != nil {
							return reconcile.Result{}, err
						}
					}

				} else if err != nil {
					omcplog.V(2).Info("err : ", err)
					return reconcile.Result{}, err
				} else {
					// Already Exist Deployment.
					if replica == 0 {
						// Delete !
						command := "delete"
						omcplog.V(2).Info("SyncResource Create (ClusterName : "+myCluster.Name+", Command : "+command+", Replicas :", replica, " / ", instance.Status.Replicas, ")")
						sync_req_name, err = r.sendSync(dep, command, myCluster.Name)

						if err != nil {
							return reconcile.Result{}, err
						}
					} else {
						// Update !
						command := "update"
						omcplog.V(2).Info("SyncResource Create (ClusterName : "+myCluster.Name+", Command : "+command+", Replicas :", replica, " / ", instance.Status.Replicas, ")")
						sync_req_name, err = r.sendSync(dep, command, myCluster.Name)
						if err != nil {
							return reconcile.Result{}, err
						}

					}

				}
			}
			// omcplog.V(2).Info("Service Notify Send")
			// r.ServiceNotify(instance.Spec.Labels, instance.Namespace)

			instance.Status.LastSpec = instance.Spec
			instance.Status.CreateSyncRequestComplete = true
			instance.Status.SyncRequestName = sync_req_name
			omcplog.V(3).Info("sync_req_name : ", sync_req_name)
			omcplog.V(2).Info("Update Status")
			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				omcplog.V(0).Info("Failed to update instance status", err)
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}

	}

	if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false {
		return reconcile.Result{}, nil
	}
	if !reflect.DeepEqual(instance.Status.LastSpec, instance.Spec) {

		omcplog.V(2).Info("[Update Detection]")
		sync_req_name := instance.Status.SyncRequestName
		if instance.Status.Replicas != instance.Spec.Replicas {
			omcplog.V(2).Info("Change Spec Replicas ! ReScheduling Start & Update Deployment")
			instance.Status.CreateSyncRequestComplete = false
			instance.Status.SchedulingNeed = true
			instance.Status.SchedulingComplete = false
			err = r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				omcplog.V(0).Info("Failed to update instance status", err)
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil

		}
		if !reflect.DeepEqual(instance.Status.LastSpec.Labels, instance.Spec.Labels) {
			// last_label := instance.Status.LastSpec.Labels
			// current_label := instance.Spec.Labels
			// omcplog.V(2).Info("Label Changed")
			// omcplog.V(2).Info("Service Notify")
			// r.ServiceNotify(last_label, instance.Namespace)
			// r.ServiceNotify(current_label, instance.Namespace)
		}

		instance.Status.LastSpec = instance.Spec
		instance.Status.SyncRequestName = sync_req_name
		omcplog.V(2).Info("sync_req_name : ", sync_req_name)
		omcplog.V(2).Info("Status Update")
		err := r.live.Status().Update(context.TODO(), instance)
		if err != nil {
			omcplog.V(2).Info("Failed to update instance status", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	sync_instance := &syncv1alpha1.Sync{}
	nsn := types.NamespacedName{
		Namespace: "openmcp",
		Name:      instance.Status.SyncRequestName,
	}
	err = r.live.Get(context.TODO(), nsn, sync_instance)
	if err == nil {
		// 아직 Sync에서 처리되지 않음
		return reconcile.Result{}, nil
	}

	// Check Deployment in cluster
	if instance.Status.BlockSubResource == false {
		omcplog.V(2).Info("[Member Cluster Check Deployment]")
		sync_req_name := ""
		for k, v := range instance.Status.ClusterMaps {
			cluster_name := k
			replica := v

			if replica == 0 {
				continue
			}

			// if _, ok := cm.Cluster_genClients[cluster_name]; !ok {
			// 	instance.Status.CreateSyncRequestComplete = false
			// 	instance.Status.SchedulingNeed = true
			// 	instance.Status.SchedulingComplete = false
			// 	err = r.live.Status().Update(context.TODO(), instance)
			// 	if err != nil {
			// 		omcplog.V(0).Info("Failed to update instance status", err)
			// 		return reconcile.Result{}, err
			// 	}
			// 	return reconcile.Result{}, nil

			// }
			found := &appsv1.Deployment{}
			cluster_client := cm.Cluster_genClients[cluster_name]
			err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)

			if err != nil && errors.IsNotFound(err) {
				// Delete Deployment Detected
				omcplog.V(2).Info("Cluster '"+cluster_name+"' ReDeployed => ", replica)
				dep := r.deploymentForOpenMCPDeployment(req, instance, replica, cluster_name)

				command := "create"
				omcplog.V(3).Info("SyncResource Create (ClusterName : "+cluster_name+", Command : "+command+", Replicas :", replica, " / ", instance.Status.Replicas, ")")
				sync_req_name, err = r.sendSync(dep, command, cluster_name)

				if err != nil {
					return reconcile.Result{}, err
				}

			}

		}
		instance.Status.SyncRequestName = sync_req_name
		omcplog.V(3).Info("sync_req_name : ", sync_req_name)

		err = r.live.Status().Update(context.TODO(), instance)
		if err != nil {
			omcplog.V(0).Info("Failed to update instance status", err)
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil // err
}
func (r *reconciler) DeleteDeploys(cm *clusterManager.ClusterManager, name string, namespace string) error {
	omcplog.V(4).Info("[OpenMCP Deployment] Function Called DeleteDeploys")

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
	omcplog.V(2).Info("Delete Check ", dep.Name, dep.Namespace)
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
	omcplog.V(4).Info("[OpenMCP Deployment] Function Called sendSync")
	syncIndex += 1

	s := &syncv1alpha1.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-deployment-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: syncv1alpha1.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *dep,
		},
	}

	err := r.live.Create(context.TODO(), s)
	return s.Name, err

}
func (r *reconciler) ServiceNotify(label_map map[string]string, namespace string) error {
	omcplog.V(4).Info("[OpenMCP Deployment] Function Called ServiceNotify")
	omcplog.V(4).Info("[OpenMCP Deployment] label_map : ", label_map)

	osvc_list := &resourcev1alpha1.OpenMCPServiceList{}
	listOptions := &client.ListOptions{Namespace: namespace}

	omcplog.V(2).Info("[OpenMCP Deployment] find Notify Target Service")
	r.live.List(context.TODO(), osvc_list, listOptions)
	for _, osvc := range osvc_list.Items {
		for k, v := range osvc.Spec.LabelSelector {
			omcplog.V(4).Info("[OpenMCP Deployment] find Notify Target Label : ", k, " / ", v)
			if label_map[k] == v {
				omcplog.V(2).Info("[OpenMCP Deployment] Service '", osvc.Name, "' Will Notify!")
				osvc.Status.ChangeNeed = true
				err := r.live.Status().Update(context.TODO(), &osvc)
				if err != nil {
					return err
				}
				omcplog.V(2).Info("[OpenMCP Deployment] Service '", osvc.Name, "' Notify Success!")

			}
		}
	}

	return nil
}
func (r *reconciler) ServiceNotifyAll(namespace string) error {
	omcplog.V(4).Info("[OpenMCP Deployment] Function Called ServiceNotifyAll")

	osvc_list := &resourcev1alpha1.OpenMCPServiceList{}
	listOptions := &client.ListOptions{Namespace: namespace}

	r.live.List(context.TODO(), osvc_list, listOptions)
	for _, osvc := range osvc_list.Items {

		omcplog.V(0).Info("->", osvc.Name, " notify !")
		osvc.Status.ChangeNeed = true
		err := r.live.Status().Update(context.TODO(), &osvc)
		if err != nil {
			return err
		}
	}

	return nil
}
func (r *reconciler) deploymentForOpenMCPDeployment(req reconcile.Request, m *resourcev1alpha1.OpenMCPDeployment, replica int32, clusterName string) *appsv1.Deployment {
	omcplog.V(4).Info("[OpenMCP Deployment] Function Called deploymentForOpenMCPDeployment")

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

	r.ApplyClusterLabel(dep, clusterName)

	reference.SetMulticlusterControllerReference(dep, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return dep
}
func (r *reconciler) ApplyClusterLabel(dep *appsv1.Deployment, clusterName string) {
	newLabel := dep.Spec.Template.ObjectMeta.Labels
	if newLabel == nil {
		newLabel = make(map[string]string)
	}
	newLabel["cluster"] = clusterName

	dep.Spec.Template.ObjectMeta.Labels = newLabel
	// dep.Spec.Selector.MatchLabels = newLabel

}

func openmcpContainersToContainers(containers []resourcev1alpha1.OpenMCPContainer) []corev1.Container {
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

func openmcpPodSpecToPodSpec(spec resourcev1alpha1.OpenMCPPodSpec) corev1.PodSpec {
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

func openmcpPodTemplateSpecToPodTemplateSpec(template resourcev1alpha1.OpenMCPPodTemplateSpec) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: template.ObjectMeta,
		Spec:       openmcpPodSpecToPodSpec(template.Spec),
	}
}

func openmcpDeploymentTemplateSpecToDeploymentSpec(templateSpec resourcev1alpha1.OpenMCPDeploymentTemplateSpec) appsv1.DeploymentSpec {
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
