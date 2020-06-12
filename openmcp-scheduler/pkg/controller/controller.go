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

package openmcpscheduler // import "admiralty.io/multicluster-controller/examples/openmcpscheduler/pkg/controller/openmcpscheduler"

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	"admiralty.io/multicluster-controller/pkg/reference"

	"openmcpscheduler/pkg/apis"
    ketiv1alpha1 "openmcpscheduler/pkg/apis/keti/v1alpha1"
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/klog"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

type ClusterManager struct {
	Fed_namespace string
    Host_config *rest.Config
    Host_client genericclient.Client
    Cluster_list *fedv1b1.KubeFedClusterList
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

	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPDeployment{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	for _, ghost := range ghosts {
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

// controller use events to eventually trigger reconcile requests.
// reconcile use clients to access API objects.
func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	klog.Infof("*********** Reconcile ***********")
	cm := NewClusterManager()

	// get OpenMCPDeployment
	instance := &ketiv1alpha1.OpenMCPDeployment{}
    err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			err := cm.DeleteDeployments(req.NamespacedName)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, err
	}

	if instance.Status.ClusterMaps == nil {

		if instance.Status.SchedulingNeed == false && instance.Status.SchedulingComplete == false {
			klog.Info("[YENA] Detect Creation...")	
			instance.Status.SchedulingNeed = true

			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
					klog.Infof("Failed to update instance status %v", err)
					return reconcile.Result{}, err
			}
			return reconcile.Result{}, err

		} else if instance.Status.SchedulingNeed == true && instance.Status.SchedulingComplete == false && strings.Compare(instance.Spec.Labels["test"], "yes") == 0 {	
			klog.Info("[SCHEDULING] Need Scheduling...")		
			
			cluster_replicas_map := cm.Scheduling(instance)
			klog.Infof("[YENA] Check cluster_replicas_map %v", cluster_replicas_map)

			instance.Status.ClusterMaps = cluster_replicas_map
			instance.Status.Replicas = instance.Spec.Replicas

			instance.Status.SchedulingNeed = false
			instance.Status.SchedulingComplete = true

			// update OpenMCPDeployment to deploy
			err := r.live.Status().Update(context.TODO(), instance)
			if err != nil {
				klog.Infof("Failed to update instance status, %v", err)
				return reconcile.Result{}, err
			}

		} else if instance.Status.SchedulingNeed == false && instance.Status.SchedulingComplete == true{
			klog.Info("[SCHEDULING] Need Deployment...")

			for _, cluster := range cm.Cluster_list.Items {

				if instance.Status.ClusterMaps[cluster.Name] == 0{
					continue
				}
				found := &appsv1.Deployment{}
				cluster_client := cm.Cluster_clients[cluster.Name]

				err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-deploy")
				if err != nil && errors.IsNotFound(err) {
					replica := instance.Status.ClusterMaps[cluster.Name]
					klog.Infof("Cluster %v Deployed (%v / %v)", cluster.Name, replica, instance.Status.Replicas)
					dep := r.deploymentForOpenMCPDeployment(req, instance, replica)
					err = cluster_client.Create(context.Background(), dep)

					if err != nil {
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
			klog.Infof("Failed to update instance status, %v", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (r *reconciler) ServiceNotify(label_map map[string]string, namespace string) error{
	fmt.Println("ServiceNotify Called", label_map)

	osvc_list := &ketiv1alpha1.OpenMCPServiceList{}
	listOptions := &client.ListOptions{Namespace: namespace}

	r.live.List(context.TODO(), listOptions, osvc_list)
	for _, osvc := range osvc_list.Items {
		for k, v := range osvc.Spec.LabelSelector{
			fmt.Println(k, " / ",v)
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

func (r *reconciler) deploymentForOpenMCPDeployment(req reconcile.Request, m *ketiv1alpha1.OpenMCPDeployment, replica int32) *appsv1.Deployment {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
		   Name:      m.Name+"-deploy",
           Namespace: m.Namespace,
		},
		Spec: openmcpDeploymentTemplateSpecToDeploymentSpec(m.Spec.Template.Spec),
    }

	if dep.Spec.Selector == nil{
		dep.Spec.Selector = &metav1.LabelSelector{}
	}

	dep.Spec.Selector.MatchLabels = m.Spec.Labels
	dep.Spec.Template.ObjectMeta.Labels = m.Spec.Labels
	dep.Spec.Replicas = &replica

	reference.SetMulticlusterControllerReference(dep, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return dep
}

func isInObject(child *appsv1.Deployment, parent string) bool {
	refKind_str := child.ObjectMeta.Annotations["multicluster.admiralty.io/controller-reference"]
	refKind_map := make(map[string]interface{})
	err := json.Unmarshal([]byte(refKind_str), &refKind_map)
	if err!= nil{
		panic(err)
	}
	if refKind_map["kind"] == parent{
		return true
	}
	return false
}

func (cm *ClusterManager) DeleteDeployments(nsn types.NamespacedName) error {
	dep := &appsv1.Deployment{}

	for _, cluster := range cm.Cluster_list.Items {
	    cluster_client := cm.Cluster_clients[cluster.Name]
		err := cluster_client.Get(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")

		if err != nil && errors.IsNotFound(err) {
			fmt.Println("Not Found")
			continue
		}

		if !isInObject(dep, "OpenMCPDeployment"){
			continue
		}

		err = cluster_client.Delete(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")
		if err != nil {
			return err
		}
    }
	return nil
}

func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		klog.Infof("Error retrieving list of federated clusters: %v", err)
	}
	if len(clusterList.Items) == 0 {
		klog.Info("No federated clusters found")
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
	// return a config object which uses the service account kubernetes gives to pods
	// *rest.Config is for talking to a Kubernetes apiserver.
	host_config, _ := rest.InClusterConfig()
	host_client := genericclient.NewForConfigOrDie(host_config)
	cluster_list := ListKubeFedClusters(host_client, fed_namespace)
	cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)
	cluster_clients := KubeFedClusterClients(cluster_list, cluster_configs)

	cm := &ClusterManager{
			Fed_namespace: fed_namespace,
			Host_config: host_config,
			Host_client: host_client,
			Cluster_list: cluster_list,
			Cluster_configs: cluster_configs,
			Cluster_clients: cluster_clients,
	}
	return cm
}

func openmcpPodTemplateSpecToPodTemplateSpec(template ketiv1alpha1.OpenMCPPodTemplateSpec) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: template.ObjectMeta,
		Spec: openmcpPodSpecToPodSpec(template.Spec),
	}
}

func openmcpDeploymentTemplateSpecToDeploymentSpec(templateSpec ketiv1alpha1.OpenMCPDeploymentTemplateSpec) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas: templateSpec.Replicas,
		Selector: templateSpec.Selector,
		Template: openmcpPodTemplateSpecToPodTemplateSpec(templateSpec.Template),
		Strategy: templateSpec.Strategy,
		MinReadySeconds: templateSpec.MinReadySeconds,
		RevisionHistoryLimit: templateSpec.RevisionHistoryLimit,
		Paused:templateSpec.Paused,
		ProgressDeadlineSeconds: templateSpec.ProgressDeadlineSeconds,
	}
}

func openmcpContainersToContainers(containers []ketiv1alpha1.OpenMCPContainer) []corev1.Container {
	var newContainers []corev1.Container

	for _, container := range containers {
		newContainer := corev1.Container{
			Name: container.Name,
			Image: container.Image,
			Command: container.Command,
			Args: container.Args,
			WorkingDir: container.WorkingDir,
			Ports: container.Ports,
			EnvFrom: container.EnvFrom,
			Env: container.Env,
			Resources: corev1.ResourceRequirements{
				Limits: container.Resources.Limits,
				Requests: container.Resources.Requests,
			},
			VolumeMounts: container.VolumeMounts,
			VolumeDevices: container.VolumeDevices,
			LivenessProbe: container.LivenessProbe,
			ReadinessProbe: container.ReadinessProbe,
			Lifecycle: container.Lifecycle,
			TerminationMessagePath: container.TerminationMessagePath,
			TerminationMessagePolicy: container.TerminationMessagePolicy,
			ImagePullPolicy: container.ImagePullPolicy,
			SecurityContext: container.SecurityContext,
			Stdin: container.Stdin,
			StdinOnce: container.StdinOnce,
			TTY: container.TTY,
		}
		newContainers = append(newContainers, newContainer)
	}

	return newContainers
}

func openmcpPodSpecToPodSpec(spec ketiv1alpha1.OpenMCPPodSpec) corev1.PodSpec {
	return corev1.PodSpec{
		Volumes: spec.Volumes,
		InitContainers: openmcpContainersToContainers(spec.InitContainers),
		Containers: openmcpContainersToContainers(spec.Containers),
		RestartPolicy: spec.RestartPolicy,
		TerminationGracePeriodSeconds: spec.TerminationGracePeriodSeconds,
		ActiveDeadlineSeconds: spec.ActiveDeadlineSeconds,
		DNSPolicy: spec.DNSPolicy,
		NodeSelector: spec.NodeSelector,
		ServiceAccountName: spec.ServiceAccountName,
		DeprecatedServiceAccount: spec.DeprecatedServiceAccount,
		AutomountServiceAccountToken: spec.AutomountServiceAccountToken,
		NodeName: spec.NodeName,
		HostNetwork: spec.HostNetwork,
		HostPID: spec.HostPID,
		HostIPC: spec.HostIPC,
		ShareProcessNamespace: spec.ShareProcessNamespace,
		SecurityContext: spec.SecurityContext,
		ImagePullSecrets: spec.ImagePullSecrets,
		Hostname: spec.Hostname,
		Subdomain: spec.Subdomain,
		Affinity: spec.Affinity,
		SchedulerName: spec.SchedulerName,
		Tolerations: spec.Tolerations,
		HostAliases: spec.HostAliases,
		PriorityClassName: spec.PriorityClassName,
		Priority: spec.Priority,
		DNSConfig: spec.DNSConfig,
		ReadinessGates: spec.ReadinessGates,
		RuntimeClassName: spec.RuntimeClassName,
		EnableServiceLinks: spec.EnableServiceLinks,
	}
}