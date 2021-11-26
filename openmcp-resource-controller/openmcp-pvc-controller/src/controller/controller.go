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

package openmcppvc

import (
	"admiralty.io/multicluster-controller/pkg/reference"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	syncv1alpha1 "openmcp/openmcp/apis/sync/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"strconv"

	"k8s.io/apimachinery/pkg/api/errors"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"openmcp/openmcp/apis"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

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
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients = append(ghostclients, ghostclient)
	}

	co := controller.New(&reconciler{live: liveClient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})

	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(context.TODO(), live, &resourcev1alpha1.OpenMCPPersistentVolumeClaim{}, controller.WatchOptions{}); err != nil {

		fmt.Println("err: ", err)
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	for _, ghost := range ghosts {
		fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(context.TODO(), ghost, &v1.PersistentVolumeClaim{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}

	return co, nil
}
func (r *reconciler) sendSync(pvc *v1.PersistentVolumeClaim, command string, clusterName string) (string, error) {
	syncIndex += 1

	s := &syncv1alpha1.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-pvc-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: syncv1alpha1.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *pvc,
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
	fmt.Println("********* [", i, "] *********")
	omcplog.V(3).Info("Namespace : ", req.Namespace, " | Name : ", req.Name, " | Context : ", req.Context)

	opvc_instance := &resourcev1alpha1.OpenMCPPersistentVolumeClaim{}
	err := r.live.Get(context.TODO(), req.NamespacedName, opvc_instance)

	if err != nil && errors.IsNotFound(err) {
		omcplog.V(3).Info("Delete PersistentVolumeClaim")

		for _, cluster := range cm.Cluster_list.Items {
			pvc := &v1.PersistentVolumeClaim{}

			cluster_client := cm.Cluster_genClients[cluster.Name]
			err = cluster_client.Get(context.TODO(), pvc, req.Namespace, req.Name)
			//delete
			if err == nil {
				pvcinstance := &v1.PersistentVolumeClaim{
					TypeMeta: metav1.TypeMeta{
						Kind:       "PersistentVolumeClaim",
						APIVersion: "v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      req.Name,
						Namespace: req.Namespace,
					},
				}
				command := "delete"
				_, err_sync := r.sendSync(pvcinstance, command, cluster.Name)

				if err_sync != nil {
					omcplog.V(3).Info("err_sync : ", err_sync)
					return reconcile.Result{}, err_sync
				} else {
					omcplog.V(3).Info("Success to Delete PVC in ", cluster.Name)
				}
			}
		}

		return reconcile.Result{}, nil

	} else if err != nil {
		omcplog.V(3).Info(err)
		return reconcile.Result{}, err
	}

	if opvc_instance.Status.ClusterMaps == nil {

		pvc := r.setPVCResourceStruct(req, opvc_instance)
		cluster_map := make(map[string]int32)

		for _, clustername := range opvc_instance.Spec.Clusters {
			foundpvc := &v1.PersistentVolumeClaim{}
			cluster_client := cm.Cluster_genClients[clustername]

			err = cluster_client.Get(context.TODO(), foundpvc, opvc_instance.Namespace, opvc_instance.Name)
			if err != nil && errors.IsNotFound(err) {
				//create
				command := "create"
				_, err_sync := r.sendSync(pvc, command, clustername)
				cluster_map[clustername] = 1
				if err_sync != nil {
					return reconcile.Result{}, err_sync
				}

				fmt.Println("Success to Create PVC in ", clustername)
			}
		}

		opvc_instance.Status.ClusterMaps = cluster_map
		opvc_instance.Status.LastSpec = opvc_instance.Spec

		err_status_update := r.live.Status().Update(context.TODO(), opvc_instance)
		if err_status_update != nil {
			fmt.Println("Failed to update instance status", err_status_update)
			return reconcile.Result{}, err_status_update
		}

	}

	return reconcile.Result{}, nil
}

func (r *reconciler) setPVCResourceStruct(req reconcile.Request, m *resourcev1alpha1.OpenMCPPersistentVolumeClaim) *v1.PersistentVolumeClaim {
	omcplog.V(4).Info("setPVCResourceStruct() Function Called")

	ls := LabelsForPVC(m.Name)

	pvc := &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			CreationTimestamp: m.CreationTimestamp,
			Name:              m.Name,
			Namespace:         m.Namespace,
			Labels:            ls,
		},

		Spec: m.Spec.Template.Spec,
	}

	reference.SetMulticlusterControllerReference(pvc, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return pvc
}

func LabelsForPVC(name string) map[string]string {
	return map[string]string{"app": "openmcppersistentvolumeclaim", "openmcppersistentvolumeclaim_cr": name}
}
