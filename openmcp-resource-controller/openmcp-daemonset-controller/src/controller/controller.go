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

package openmcpdaemonset

import (
	"context"
	"fmt"
	syncv1alpha1 "openmcp/openmcp/apis/sync/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"strconv"

	"admiralty.io/multicluster-controller/pkg/reference"
	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/api/errors"

	"openmcp/openmcp/apis"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	cm = myClusterManager

	liveClient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}

	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			omcplog.V(4).Info("Error getting delegating client for ghost cluster [", ghost.Name, "]")
			//return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		} else {
			ghostclients = append(ghostclients, ghostclient)
		}
	}

	co := controller.New(&reconciler{live: liveClient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})

	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(context.TODO(), live, &resourcev1alpha1.OpenMCPDaemonSet{}, controller.WatchOptions{}); err != nil {

		fmt.Println("err: ", err)
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	for _, ghost := range ghosts {
		fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(context.TODO(), ghost, &appsv1.DaemonSet{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}

	return co, nil
}
func (r *reconciler) sendSync(ds *appsv1.DaemonSet, command string, clusterName string) (string, error) {
	syncIndex += 1

	s := &syncv1alpha1.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-daemonset-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: syncv1alpha1.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *ds,
		},
	}
	err := r.live.Create(context.TODO(), s)

	if err != nil {
		omcplog.V(0).Info("syncErr - ", err)
	}

	return s.Name, err

}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

var i = 0
var syncIndex = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	i += 1
	omcplog.V(2).Info("********* [", i, "] *********")
	omcplog.V(2).Info("Namespace : ", req.Namespace, " | Name : ", req.Name, " | Context : ", req.Context)

	ods_instance := &resourcev1alpha1.OpenMCPDaemonSet{}
	err := r.live.Get(context.TODO(), req.NamespacedName, ods_instance)

	if err != nil && errors.IsNotFound(err) {

		for _, cluster := range cm.Cluster_list.Items {
			ds := &appsv1.DaemonSet{}

			cluster_client := cm.Cluster_genClients[cluster.Name]
			err = cluster_client.Get(context.TODO(), ds, req.Namespace, req.Name)
			//delete
			if err == nil {
				odsinstance := &appsv1.DaemonSet{
					TypeMeta: metav1.TypeMeta{
						Kind:       "DaemonSet",
						APIVersion: "apps/v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      req.Name,
						Namespace: req.Namespace,
					},
				}
				command := "delete"
				_, err_sync := r.sendSync(odsinstance, command, cluster.Name)

				if err_sync != nil {
					omcplog.V(0).Info("err_sync : ", err_sync)
					return reconcile.Result{}, err_sync
				} else {
					omcplog.V(2).Info("Success to Delete DaemonSet in " + cluster.Name)
				}
			}
		}

		return reconcile.Result{}, nil

	} else if err != nil {
		omcplog.V(0).Info(err)
		return reconcile.Result{}, err
	}

	if ods_instance.Status.ClusterMaps == nil {

		ds := r.setDSResourceStruct(req, ods_instance)
		cluster_map := make(map[string]int32)

		if len(ods_instance.Spec.Clusters) == 0 {
			omcplog.V(2).Info("Deploy DaemonSet Resource On All Clusters ...")

			for _, cluster := range cm.Cluster_list.Items {
				foundds := &appsv1.DaemonSet{}
				cluster_client := cm.Cluster_genClients[cluster.Name]

				err = cluster_client.Get(context.TODO(), foundds, ods_instance.Namespace, ods_instance.Name)
				if err != nil && errors.IsNotFound(err) {
					//create
					command := "create"
					_, err_sync := r.sendSync(ds, command, cluster.Name)
					cluster_map[cluster.Name] = 1
					if err_sync != nil {
						return reconcile.Result{}, err_sync
					}

					omcplog.V(2).Info("Success to Create DaemonSet in " + cluster.Name)
				}
			}
		} else {
			omcplog.V(2).Info("Deploy DaemonSet Resource On Specified Clusters ...")

			for _, clustername := range ods_instance.Spec.Clusters {
				foundds := &appsv1.DaemonSet{}
				cluster_client := cm.Cluster_genClients[clustername]

				err = cluster_client.Get(context.TODO(), foundds, ods_instance.Namespace, ods_instance.Name)
				if err != nil && errors.IsNotFound(err) {
					//create
					command := "create"
					_, err_sync := r.sendSync(ds, command, clustername)
					cluster_map[clustername] = 1
					if err_sync != nil {
						return reconcile.Result{}, err_sync
					}

					omcplog.V(2).Info("Success to Create DaemonSet in " + clustername)
				}
			}
		}
		ods_instance.Status.ClusterMaps = cluster_map
		ods_instance.Status.LastSpec = ods_instance.Spec
		ods_instance.Status.CheckSubResource = true

		err_status_update := r.live.Status().Update(context.TODO(), ods_instance)
		if err_status_update != nil {
			omcplog.V(0).Info("Failed to update instance status", err_status_update)
			return reconcile.Result{}, err_status_update
		}

	}

	return reconcile.Result{}, nil
}

func (r *reconciler) setDSResourceStruct(req reconcile.Request, m *resourcev1alpha1.OpenMCPDaemonSet) *appsv1.DaemonSet {
	omcplog.V(4).Info("setDSResourceStruct() Function Called")

	ls := LabelsForDS(m.Name)

	ds := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},

		Spec: m.Spec.Template.Spec,
	}

	reference.SetMulticlusterControllerReference(ds, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return ds
}

func LabelsForDS(name string) map[string]string {
	return map[string]string{"app": "openmcpdaemonset", "openmcpdaemonset_cr": name}
}
