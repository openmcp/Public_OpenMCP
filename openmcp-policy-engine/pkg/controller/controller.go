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

package controller // import "admiralty.io/multicluster-controller/examples/openmcppolicyengine/pkg/controller/openmcppolicyengine"

import (
	"context"
	"fmt"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"k8s.io/apimachinery/pkg/api/errors"
	//	"admiralty.io/multicluster-controller/pkg/reference"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"openmcp/openmcp/openmcp-resource-controller/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"

	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)

var log = logf.Log.WithName("controller_openmcphybridautoscaler")

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("[OpenMCP Policy Engine] Function Called NewController")
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

	omcplog.V(4).Info(live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPPolicy{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

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
	omcplog.V(4).Info("********* [", i, "] *********")
	omcplog.V(5).Info("Request Context: ", req.Context, " / Request Namespace: ", req.Namespace, " /  Request Name: ", req.Name)
	//cm := NewClusterManager()

	// Fetch the OpenMCPDeployment instance
	instance := &ketiv1alpha1.OpenMCPPolicy{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	//	omcplog.V(4).Info("instance: ", instance)
	omcplog.V(5).Info("instance Name: ", instance.Name)
	omcplog.V(5).Info("instance Namespace: ", instance.Namespace)

	if err != nil {
		if errors.IsNotFound(err) {
			omcplog.V(2).Info("Delete Policy Resource")
			return reconcile.Result{}, nil
		}
		omcplog.V(0).Info("Error: ", err)

		return reconcile.Result{}, err
	}

	if instance.Spec.PolicyStatus == "Disabled" {
		omcplog.V(2).Info("Policy Disabled")
	} else if instance.Spec.PolicyStatus == "Enabled" {
		if instance.Spec.RangeOfApplication == "FromNow" {
			omcplog.V(2).Info("Policy Enabled - FromNow")
		} else if instance.Spec.RangeOfApplication == "All" {
			//omcplog.V(4).Info("Policy Enabled - All")
			object := instance.Spec.Template.Spec.TargetController.Kind
			if object == "OpenMCPHybridAutoScaler" {
				omcplog.V(2).Info("Policy Enabled - OpenMCPHybridAutoScaler")
				hpaList := &ketiv1alpha1.OpenMCPHybridAutoScalerList{}
				listOptions := &client.ListOptions{Namespace: ""} //all resources
				r.live.List(context.TODO(), hpaList, listOptions)
				//omcplog.V(4).Info("List: ", hpaList)
				for _, hpaInstance := range hpaList.Items {
					//omcplog.V(4).Info("hpastatus: ",hpaInstance.Status.Policies)
					//omcplog.V(4).Info("policies: ",instance.Spec.Template.Spec.Policies)
					var i = 0
					for index, tmpPolicy := range hpaInstance.Status.Policies { //정책 이름 대조하여 해당 정책만 수정
						if tmpPolicy.Type == instance.Spec.Template.Spec.Policies[0].Type { //같은 정책이 이미 있는 경우
							i++
							hpaInstance.Status.Policies[index].Value = instance.Spec.Template.Spec.Policies[0].Value
							break
						}
					}
					if i == 0 {
						hpaInstance.Status.Policies = append(hpaInstance.Status.Policies, instance.Spec.Template.Spec.Policies...)
					}
					err := r.live.Status().Update(context.TODO(), &hpaInstance)
					if err != nil {
						omcplog.V(0).Info("OpenMCPHPA Policy Update Error")
						return reconcile.Result{}, err
					} else {
						omcplog.V(2).Info("OpenMCPHPA Policy UPDATE Success!")
					}
				}
			} else if object == "OpenMCPLoadbalancer" {

			}
		}
	}
	return reconcile.Result{}, nil
}

/*func (cm *ClusterManager) DeleteOpenMCPPolicyEngine(nsn types.NamespacedName) error {
	dep := &appsv1.Deployment{}
	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_clients[cluster.Name]
		omcplog.V(4).Info(nsn.Namespace, nsn.Name)
		err := cluster_client.Get(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")
		if err != nil && errors.IsNotFound(err) {
			// all good
			omcplog.V(4).Info("Not Found")
			continue
		}
		omcplog.V(4).Info(cluster.Name," Delete Start")
		err = cluster_client.Delete(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")
		if err != nil {
			return err
		}
		omcplog.V(4).Info(cluster.Name, "Delete Complate")
	}
	return nil

}*/

