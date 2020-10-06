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

package openmcpingress

import (
	"context"
	"fmt"
	"github.com/getlantern/deepcopy"
	"github.com/mittwald/go-powerdns/apis/zones"
	corev1 "k8s.io/api/core/v1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-dns-controller/pkg/mypdns"
	"strconv"

	"admiralty.io/multicluster-controller/pkg/reference"
	"k8s.io/apimachinery/pkg/api/errors"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"k8s.io/apimachinery/pkg/types"
	"openmcp/openmcp/apis"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	extv1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	syncv1alpha1 "openmcp/openmcp/apis/sync/v1alpha1"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("[OpenMCP Ingress Controller] Function Called NewController")
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

	if err := co.WatchResourceReconcileObject(live, &resourcev1alpha1.OpenMCPIngress{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	if err := co.WatchResourceReconcileController(live, &extv1b1.Ingress{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
	}

	for _, ghost := range ghosts {
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
	omcplog.V(4).Info("[OpenMCP Ingress Controller] Function Called Reconcile")
	i += 1
	omcplog.V(5).Info("********* [", i, "] *********")
	omcplog.V(3).Info(req.Context, " / ", req.Namespace, " / ", req.Name)

	instance := &resourcev1alpha1.OpenMCPIngress{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	omcplog.V(3).Info("instance Name: ", instance.Name)
	omcplog.V(3).Info("instance Namespace : ", instance.Namespace)

	if err != nil {
		if errors.IsNotFound(err) {
			omcplog.V(3).Info("Delete Deployments ..Cluster")
			err := r.DeleteIngress(cm, req.NamespacedName.Name, req.NamespacedName.Namespace)
			return reconcile.Result{}, err
		}
		omcplog.V(1).Info(err)
		return reconcile.Result{}, err
	}
	if instance.Status.ClusterMaps == nil || instance.Status.ChangeNeed == true {
		omcplog.V(3).Info("Ingress Create Start")
		r.createIngress(req, cm, instance)

		if err != nil {
			omcplog.V(1).Info(err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil

	} else {
		// Check Ingress In Openmcp
		omcplog.V(3).Info("Check Ingress In Openmcp")
		foundIngress := &extv1b1.Ingress{}
		err = cm.Host_client.Get(context.TODO(), foundIngress, instance.Namespace, instance.Name)

		if err == nil {
			foundService := &corev1.Service{}
			nsn := types.NamespacedName{
				Namespace: "openmcp",
				Name:      "openmcp-loadbalancing-controller",
			}
			err = r.live.Get(context.TODO(), nsn, foundService)
			if err != nil && errors.IsNotFound(err) {
				omcplog.V(0).Info("LoadBalancing-controller Service Not Found")
				return reconcile.Result{}, err
			} else {
				omcplog.V(3).Info("Update Ingress Status")
				foundIngress.Status.LoadBalancer = foundService.Status.LoadBalancer
				err = cm.Host_client.UpdateStatus(context.Background(), foundIngress)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		} else if errors.IsNotFound(err) {
			omcplog.V(3).Info("Create Ingress")
			host_ing, _ := r.ingressForOpenMCPIngress(req, instance)
			err = cm.Host_client.Create(context.Background(), host_ing)
			if err != nil {
				return reconcile.Result{}, err
			}
		}

		// Check Ingress in cluster
		for k, _ := range instance.Status.ClusterMaps {
			cluster_name := k
			found := &extv1b1.Ingress{}
			cluster_client := cm.Cluster_genClients[cluster_name]
			err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
			if err != nil && errors.IsNotFound(err) {
				// Delete Ingress Detected
				omcplog.V(3).Info("Cluster '" + cluster_name + "' ReDeployed")
				_, ing := r.ingressForOpenMCPIngress(req, instance)

				command := "create"
				_, err = r.sendSync(ing, command, cluster_name)

				if err != nil {
					return reconcile.Result{}, err
				}

			}
		}
	}

	return reconcile.Result{}, nil // err
}

func (r *reconciler) registerPdnsServer(ingress *extv1b1.Ingress) error {
	omcplog.V(4).Info("[OpenMCP Ingress Controller] Function Called registerPdnsServer")
	pdnsClient, err := mypdns.PdnsNewClient()
	if err != nil {
		omcplog.V(0).Info(err)
		return err
	}

	found := &corev1.Service{}
	nsn := types.NamespacedName{
		Namespace: "openmcp",
		Name:      "openmcp-loadbalancing-controller",
	}
	err = r.live.Get(context.TODO(), nsn, found)
	if err != nil && errors.IsNotFound(err) {
		omcplog.V(0).Info("LoadBalancing-controller Service Not Found")
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
			omcplog.V(0).Info(err)
			return err
		}
	}

	return nil

}

func (r *reconciler) createIngress(req reconcile.Request, cm *clusterManager.ClusterManager, instance *resourcev1alpha1.OpenMCPIngress) error {
	omcplog.V(4).Info("[OpenMCP Ingress Controller] Function Called createIngress")
	host_ing, ing := r.ingressForOpenMCPIngress(req, instance)

	found := &extv1b1.Ingress{}
	err := cm.Host_client.Get(context.TODO(), found, instance.Namespace, instance.Name)
	if err != nil && errors.IsNotFound(err) {
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

		if err != nil && errors.IsNotFound(err) {
			for i, rule := range cluster_ing.Spec.Rules {
				for _, paths := range rule.HTTP.Paths {
					serviceName := paths.Backend.ServiceName
					omcplog.V(5).Info("service name")
					omcplog.V(5).Info(serviceName)
					serviceErr := cluster_client.Get(context.TODO(), serviceFound, instance.Namespace, serviceName)
					if serviceErr != nil && errors.IsNotFound(serviceErr) {
						omcplog.V(0).Info("service not found")
						isService = false
					}
				}

				cluster_ing.Spec.Rules[i].Host = cluster.Name + "." + rule.Host

			}
			if isService == true {
				omcplog.V(3).Info("Create Ingress Resource - ", cluster.Name)
				command := "create"
				_, err = r.sendSync(cluster_ing, command, cluster.Name)
				cluster_map[cluster.Name] = 1
				if err != nil {
					omcplog.V(0).Info(cluster.Name, " - ", err)
				}
			}
		}
	}
	instance.Status.ClusterMaps = cluster_map

	err = r.live.Status().Update(context.TODO(), instance)
	return err
}

func (r *reconciler) ingressForOpenMCPIngress(req reconcile.Request, m *resourcev1alpha1.OpenMCPIngress) (*extv1b1.Ingress, *extv1b1.Ingress) {
	omcplog.V(4).Info("[OpenMCP Ingress Controller] Function Called ingressForOpenMCPIngress")
	host_ing := &extv1b1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
	}
	deepcopy.Copy(&host_ing.Spec, &m.Spec.Template.Spec)

	host_ing.Namespace = "openmcp"

	for i, _ := range host_ing.Spec.Rules {
		for j, _ := range host_ing.Spec.Rules[i].HTTP.Paths {
			host_ing.Spec.Rules[i].HTTP.Paths[j].Backend.ServiceName = "openmcp-loadbalancing-controller"
			host_ing.Spec.Rules[i].HTTP.Paths[j].Backend.ServicePort.IntVal = 80
		}
	}

	reference.SetMulticlusterControllerReference(host_ing, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	ing := &extv1b1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
	}
	deepcopy.Copy(&ing.Spec, &m.Spec.Template.Spec)

	reference.SetMulticlusterControllerReference(ing, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return host_ing, ing
}

func (r *reconciler) DeleteIngress(cm *clusterManager.ClusterManager, name string, namespace string) error {
	omcplog.V(4).Info("[OpenMCP Ingress Controller] Function Called DeleteIngress")
	ing := &extv1b1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	err := cm.Host_client.Get(context.Background(), ing, namespace, name)
	if err != nil && errors.IsNotFound(err) {
		omcplog.V(0).Info("Not Found")
	} else if err != nil && !errors.IsNotFound(err) {
		return err
	}
	omcplog.V(3).Info("OpenMCP Delete Start")
	err = cm.Host_client.Delete(context.Background(), ing, namespace, name)

	if err != nil {
		return err
	}
	omcplog.V(3).Info("OpenMCP Delete Complete")

	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_genClients[cluster.Name]
		omcplog.V(0).Info(namespace, name)
		err := cluster_client.Get(context.Background(), ing, namespace, name)
		if err != nil && errors.IsNotFound(err) {
			omcplog.V(0).Info("Not Found")
			continue
		}
		omcplog.V(3).Info(cluster.Name, " Delete Start")
		command := "delete"
		omcplog.V(0).Info(name)
		omcplog.V(0).Info(namespace)

		ing = &extv1b1.Ingress{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Ingress",
				APIVersion: "networking.k8s.io/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
		_, err = r.sendSync(ing, command, cluster.Name)

		if err != nil {
			return err
		}
		omcplog.V(3).Info(cluster.Name, "Delete Complete")
	}
	return nil

}

var syncIndex int = 0

func (r *reconciler) sendSync(ingress *extv1b1.Ingress, command string, clusterName string) (string, error) {
	omcplog.V(4).Info("[OpenMCP Ingress Controller] Function Called sendSync")
	syncIndex += 1

	s := &syncv1alpha1.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-ingress-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: syncv1alpha1.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *ingress,
		},
	}
	omcplog.V(5).Info("Delete Check ", s.Spec.Template.(extv1b1.Ingress).Name, s.Spec.Template.(extv1b1.Ingress).Namespace)

	err := r.live.Create(context.TODO(), s)

	if err != nil {
		omcplog.V(1).Info(err)
	}

	return s.Name, err
}
