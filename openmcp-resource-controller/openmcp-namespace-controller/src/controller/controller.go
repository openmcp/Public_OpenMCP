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

package openmcpnamespace

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
	"time"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"admiralty.io/multicluster-controller/pkg/reference"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("Function Called NewController")
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
	r := &reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}

	co := controller.New(r, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &resourcev1alpha1.OpenMCPNamespace{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	for _, ghost := range ghosts {
		if err := co.WatchResourceReconcileController(context.TODO(), ghost, &corev1.Namespace{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}

	//r.newClusterDeployNamespace()
	return co, nil
}
func (r *reconciler) newClusterDeployNamespace() error {
	omcplog.V(4).Info("Function Called newClusterDeployNamespace")

	ns := &corev1.Namespace{}
	onList, err := cm.Crd_client.OpenMCPNamespace("default").List(metav1.ListOptions{})

	if err != nil && errors.IsNotFound(err) {
		omcplog.V(2).Info("Not Exist OpenMCPNamespaceList Resource")
		return err
	} else if err != nil {
		return err
	}
	omcplog.V(2).Info("Exist OpenMCPNamespaceList Resource ", len(onList.Items))
	for _, ons := range onList.Items {
		for _, cl := range cm.Cluster_list.Items {

			err = cm.Cluster_genClients[cl.Name].Get(context.TODO(), ns, metav1.NamespaceDefault, ons.Name)
			if err != nil && errors.IsNotFound(err) {
				err = r.live.Status().Update(context.TODO(), &ons)
				if err != nil {
					omcplog.V(1).Info("Failed to update instance status", err)
					return err
				}
				return nil

			} else if err != nil {
				return err
			}
		}
	}
	return nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

var i int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(4).Info("Function Called Reconcile")
	i += 1
	omcplog.V(5).Info("********* [", i, "] *********")
	omcplog.V(3).Info(req.Context, " / ", req.Namespace, " / ", req.Name)

	// Fetch the OpenMCPDeployment instance
	instance := &resourcev1alpha1.OpenMCPNamespace{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	omcplog.V(4).Info("instance Name: ", instance.Name)
	omcplog.V(4).Info("instance Namespace : ", instance.Namespace)

	if err != nil {
		if errors.IsNotFound(err) {
			// ...TODO: multicluster garbage collector
			// Until then...
			omcplog.V(3).Info("Delete OpenMCPNamespace")
			err := r.DeleteNamespace(cm, req.NamespacedName.Name, req.NamespacedName.Namespace)
			return reconcile.Result{}, err
		}
		omcplog.V(1).Info(err)
		return reconcile.Result{}, err
	}

	if instance.Status.ClusterMaps == nil {
		err := r.createNamespace(req, cm, instance)
		if err != nil {
			omcplog.V(1).Info(err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	if !reflect.DeepEqual(instance.Status.LastSpec, instance.Spec) {
		err := r.updateNamespace(req, cm, instance)
		if err != nil {
			omcplog.V(1).Info(err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	//err = r.live.Status().Update(context.TODO(), instance)
	//if err != nil {
	//	omcplog.V(1).Info("Failed to update instance status", err)
	//	return reconcile.Result{}, err
	//}
	//
	//return reconcile.Result{}, nil // err
	// Check Job in cluster
	if instance.Status.CheckSubResource == true {
		omcplog.V(2).Info("[Member Cluster Check Job]")
		for k, v := range instance.Status.ClusterMaps {
			cluster_name := k
			replica := v

			if v == 0 {
				continue
			}
			found := &corev1.Namespace{}
			cluster_client := cm.Cluster_genClients[cluster_name]
			fmt.Println(cm.Cluster_genClients[cluster_name], cluster_name)
			err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)

			if err != nil && errors.IsNotFound(err) {
				// Delete Service Detected
				omcplog.V(2).Info("Cluster '"+cluster_name+"' ReDeployed => ", replica)
				svc := r.namespaceForOpenMCPNamespace(req, instance)

				command := "create"
				omcplog.V(3).Info("SyncResource Create (ClusterName : "+cluster_name+", Command : "+command+", Replicas :", replica, ")")
				_, err = r.sendSync(svc, command, cluster_name)

				if err != nil {
					return reconcile.Result{}, err
				}

			}

		}

	}
	return reconcile.Result{}, nil // err
}

func (r *reconciler) namespaceForOpenMCPNamespace(req reconcile.Request, m *resourcev1alpha1.OpenMCPNamespace) *corev1.Namespace {
	omcplog.V(4).Info("Function Called namespaceForOpenMCPNamespace")

	newLabel := m.Labels
	if newLabel == nil {
		newLabel = make(map[string]string)
	}
	newLabel["istio-injection"] = "enabled"

	dep := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: m.Name,
			//Namespace: m.Namespace,
			Labels: newLabel,
		},
	}

	//deepcopy.Copy(&dep.Spec, &m.Spec.Template.Spec)

	reference.SetMulticlusterControllerReference(dep, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return dep
}

var syncIndex int = 0

func (r *reconciler) sendSync(secret *corev1.Namespace, command string, clusterName string) (string, error) {
	omcplog.V(4).Info("Function Called sendSync")

	syncIndex += 1

	s := &syncv1alpha1.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-namespace-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: syncv1alpha1.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *secret,
		},
	}
	// omcplog.V(5).Info("Delete Check ", s.Spec.Template.(corev1.Namespace).Name, s.Spec.Template.(corev1.Namespace).Namespace)

	err := r.live.Create(context.TODO(), s)

	if err != nil {
		omcplog.V(0).Info(err)
	}

	omcplog.V(0).Info(s.Name)
	return s.Name, err
}

func (r *reconciler) createNamespace(req reconcile.Request, cm *clusterManager.ClusterManager, instance *resourcev1alpha1.OpenMCPNamespace) error {
	omcplog.V(4).Info("Function Called createNamespace")
	dep := r.namespaceForOpenMCPNamespace(req, instance)

	// err := cm.Host_client.Create(context.TODO(), dep)
	// if err != nil {
	// 	return err
	// }
	cluster_map := make(map[string]int32)
	for _, cluster := range cm.Cluster_list.Items {
		found := &corev1.Namespace{}
		cluster_client := cm.Cluster_genClients[cluster.Name]

		err := cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name)

		if err != nil && errors.IsNotFound(err) {
			omcplog.V(3).Info("Cluster '" + cluster.Name + "' Deployed")

			command := "create"
			_, err := r.sendSync(dep, command, cluster.Name)
			cluster_map[cluster.Name] = 1
			if err != nil {
				return err
			}
		}

	}
	instance.Status.ClusterMaps = cluster_map
	instance.Status.LastSpec = instance.Spec
	err := r.live.Status().Update(context.TODO(), instance)
	return err
}

func (r *reconciler) DeleteNamespace(cm *clusterManager.ClusterManager, name string, namespace string) error {
	omcplog.V(4).Info("Function Called DeleteNamespace")

	dep := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			//Namespace: namespace,
		},
	}
	err := cm.Host_client.Delete(context.TODO(), dep, dep.Namespace, dep.Name)
	if err != nil && errors.IsNotFound(err) {
		fmt.Println(err)
	} else if err != nil {
		return err
	}

	for _, cluster := range cm.Cluster_list.Items {

		omcplog.V(3).Info(cluster.Name, " Delete Start")

		command := "delete"
		_, err := r.sendSync(dep, command, cluster.Name)

		if err != nil {
			return err
		}
		omcplog.V(3).Info(cluster.Name, "Delete Complete")
	}
	return nil
}

func (r *reconciler) updateNamespace(req reconcile.Request, cm *clusterManager.ClusterManager, instance *resourcev1alpha1.OpenMCPNamespace) error {
	omcplog.V(4).Info("Function Called updateNamespace")

	for _, cluster := range cm.Cluster_list.Items {

		omcplog.V(3).Info("Cluster '" + cluster.Name + "' updated")
		dep := r.namespaceForOpenMCPNamespace(req, instance)
		command := "update"
		_, err := r.sendSync(dep, command, cluster.Name)
		if err != nil {
			return err
		}

	}
	instance.Status.LastSpec = instance.Spec
	err := r.live.Status().Update(context.TODO(), instance)
	return err

}
func CheckClusterNamespaceStatus(myClusterManager *clusterManager.ClusterManager, quit, quitok chan bool) {
	cm = myClusterManager
	for {
		select {
		case <-quit:
			omcplog.V(2).Info("CheckClusterNamespaceStatus Quit")
			quitok <- true
			return
		default:
			omcplog.V(2).Info("CheckClusterNamespaceStatus Start")
			onsList, err := cm.Crd_client.OpenMCPNamespace("default").List(metav1.ListOptions{})
			if err != nil {
				omcplog.V(0).Info(err)
			}
			for _, ons := range onsList.Items {
				nsStatusActiveAll := true

				for k, v := range ons.Status.ClusterMaps {
					clusterName := k

					if v == 0 {
						continue
					}

					if cm.Cluster_genClients[clusterName] != nil {
						ns := &corev1.Namespace{}
						err = cm.Cluster_genClients[clusterName].Get(context.TODO(), ns, corev1.NamespaceDefault, ons.Name)
						if err != nil && errors.IsNotFound(err) {
							omcplog.V(2).Info(err, ": Create Namespace In Cluster '", clusterName, "'")
							newLabel := ons.Labels
							if newLabel == nil {
								newLabel = make(map[string]string)
							}
							newLabel["istio-injection"] = "enabled"

							CreateNSInstance := &corev1.Namespace{
								TypeMeta: metav1.TypeMeta{
									Kind:       "Namespace",
									APIVersion: "v1",
								},
								ObjectMeta: metav1.ObjectMeta{
									Name:   ons.Name,
									Labels: newLabel,
								},
							}
							reference.SetMulticlusterControllerReference(CreateNSInstance, reference.NewMulticlusterOwnerReference(&ons, ons.GroupVersionKind(), "openmcp"))

							err2 := cm.Cluster_genClients[clusterName].Create(context.TODO(), CreateNSInstance)
							if err2 != nil {
								omcplog.V(0).Info(err2)
							}
							nsStatusActiveAll = false

						} else if err != nil && errors.IsAlreadyExists(err) {
							if ns.Status.Phase != "Active" {
								nsStatusActiveAll = false
								break
							}

						} else if err != nil {
							omcplog.V(0).Info(err)
							nsStatusActiveAll = false
							break
						} else {
							if ns.Status.Phase != "Active" {
								nsStatusActiveAll = false
								break
							}
						}
					}

				}
				if nsStatusActiveAll {
					ns := CreateNSFromONS(ons)

					checkNS := &corev1.Namespace{}
					err := cm.Host_client.Get(context.TODO(), checkNS, corev1.NamespaceDefault, ons.Name)
					if err != nil && errors.IsNotFound(err) {

						err = cm.Host_client.Create(context.TODO(), ns)
						if err != nil {
							omcplog.V(0).Info(err)
						}
					} else if err == nil {
						err = cm.Host_client.Update(context.TODO(), ns)
						if err != nil {
							omcplog.V(0).Info(err)
						}
					} else {
						omcplog.V(0).Info(err)
					}

				}

			}
			time.Sleep(1 * time.Second)

		}

	}
}

func CreateNSFromONS(ons resourcev1alpha1.OpenMCPNamespace) *corev1.Namespace {
	newLabel := ons.Labels
	if newLabel == nil {
		newLabel = make(map[string]string)
	}
	newLabel["istio-injection"] = "enabled"

	ns := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   ons.Name,
			Labels: newLabel,
		},
	}
	reference.SetMulticlusterControllerReference(ns, reference.NewMulticlusterOwnerReference(&ons, ons.GroupVersionKind(), "openmcp"))
	return ns
}
