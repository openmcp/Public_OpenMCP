/*
Copyright 2018 The Multicluster-Controller Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in ccmpliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"admiralty.io/multicluster-controller/pkg/reference"
	"context"
	"fmt"
	autoscaling "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/types"
	//"k8s.io/klog"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"os"
	"strconv"
	"time"
	"reflect"
	//"sort"
	//"math/rand"
	//"time"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"k8s.io/apimachinery/pkg/api/errors"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"openmcp/openmcp/openmcp-resource-controller/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"

	vpav1beta2 "github.com/kubernetes/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	hpav2beta2 "k8s.io/api/autoscaling/v2beta2"
	//vpav1beta2 "k8s.io/api/autoscaling/v1beta2"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"

	"openmcp/openmcp/openmcp-resource-controller/controllers/openmcp-has-controller/pkg/protobuf"

	syncapis "openmcp/openmcp/openmcp-sync-controller/pkg/apis"
	sync "openmcp/openmcp/openmcp-sync-controller/pkg/apis/keti/v1alpha1"
)

var log = logf.Log.WithName("controller_openmcphybridautoscaler")

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string) (*controller.Controller, error) {
	//fmt.Println("Step 1.	NewController()")
	liveClient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}

	ghostClients := map[string]client.Client{}
	for _, ghost := range ghosts {
		ghostTmp, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostClients[ghost.Name] = ghostTmp
	}

	co := controller.New(&reconciler{live: liveClient, ghosts: ghostClients, ghostNamespace: ghostNamespace}, controller.Options{})

	//live.GetScheme()에 apis scheme 추가
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	//live.GetScheme()에 vpav1beta2 scheme 추가
	if err := vpav1beta2.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	if err := syncapis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: #{err}")
	}

	fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPHybridAutoScaler{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	for _, ghost := range ghosts {
		fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(ghost, &hpav2beta2.HorizontalPodAutoscaler{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}

	for _, ghost := range ghosts {
		fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(ghost, &vpav1beta2.VerticalPodAutoscaler{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}

	return co, nil
}

func (r *reconciler) sendSyncHPA(hpa *hpav2beta2.HorizontalPodAutoscaler, command string, clusterName string) (string, error) {
	syncIndex += 1

	s := &sync.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-hybridautoscaler-hpa-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: sync.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *hpa,
		},
	}
	err := r.live.Create(context.TODO(), s)

	if err != nil {
		omcplog.V(0).Info("syncErr - ", err)
	}

	return s.Name, err

}

func (r *reconciler) sendSyncVPA(vpa *vpav1beta2.VerticalPodAutoscaler, command string, clusterName string) (string, error) {
	syncIndex += 1

	s := &sync.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-hybridautoscaler-vpa-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: sync.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    *vpa,
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
	ghosts         map[string]client.Client
	ghostNamespace string
}

var i = 0
var syncIndex = 0
var lastTimeRebalancing time.Time
var tmpMap = map[string]int32{}

func (r *reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	//klog.V(0).Info("[OpenMCPHAS] Function Called Reconcile")
	//reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	i += 1

	omcplog.V(3).Info("********* [", i, "] *********")
	omcplog.V(3).Info("Namespace : ", request.Namespace, " | Name : ", request.Name, " | Context : ", request.Context)

	cm := clusterManager.NewClusterManager()
	type ObjectKey = types.NamespacedName

	hasInstance := &ketiv1alpha1.OpenMCPHybridAutoScaler{}
	err := r.live.Get(context.TODO(), request.NamespacedName, hasInstance)

	//OpenMCPHAS 리소스 삭제
	if err != nil {
		if errors.IsNotFound(err) {
			//var sync_name string
			for _, cluster := range cm.Cluster_list.Items {
				_ = r.DeleteHPAVPA(cm, cluster.Name, request.Namespace, request.Name)
			}
			return reconcile.Result{}, nil
		}
		omcplog.V(0).Info("!!! Failed to get hasInstance")
		//reqLogger.Error(err, "!!! Failed to get hasInstance")
		return reconcile.Result{}, err
	}

	if !reflect.DeepEqual(hasInstance.Status.LastSpec, hasInstance.Spec) {

		//"has-target-cluster" 정책 적용
		clusterListItems, hasInstance := r.UpdateTargetClusterPolicy(cm, hasInstance)

		target_cluster_policy_err := r.live.Status().Update(context.TODO(), hasInstance)
		if target_cluster_policy_err != nil {
			omcplog.V(0).Info("!!! Policy Status Update Error")
			return reconcile.Result{}, target_cluster_policy_err
		} else {
			omcplog.V(3).Info(">>> Policy Status UPDATE Success")
		}

		//타겟 OpenMCPDeployment Get
		openmcpDep := &ketiv1alpha1.OpenMCPDeployment{}
		openmcpDep_err := r.live.Get(context.TODO(), ObjectKey{Namespace: hasInstance.Namespace, Name: hasInstance.Spec.HpaTemplate.Spec.ScaleTargetRef.Name}, openmcpDep)

		omcplog.V(3).Info(">>> Target OpenMCPDeployment [", openmcpDep.Name, " | ", openmcpDep.Namespace, "]")

		if openmcpDep_err == nil {

			var dep_list_for_hpa []string
			var cluster_dep_replicas map[string]int32
			var cluster_dep_request map[string]bool

			cluster_dep_replicas = make(map[string]int32)
			cluster_dep_request = make(map[string]bool)

			//HPA, VPA가 배포될 클러스터 리스트 저장 (has-target-cluster 정책 적용 및 Deployment 배포 확인)
			for _, cluster := range clusterListItems {
				dep := &appsv1.Deployment{}
				cluster_client := cm.Cluster_genClients[cluster.Name]
				dep_err := cluster_client.Get(context.TODO(), dep, hasInstance.Namespace, hasInstance.Spec.HpaTemplate.Spec.ScaleTargetRef.Name)
				if dep_err == nil {
					//fmt.Println("deploy request: ", dep.Spec.Template.Spec.Containers[0].Resources.Requests) // 컨테이너 여러개일 경우?
					if dep.Spec.Template.Spec.Containers[0].Resources.Requests == nil {
						cluster_dep_request[cluster.Name] = false
					} else {
						cluster_dep_request[cluster.Name] = true
					}
					cluster_dep_replicas[cluster.Name] = *dep.Spec.Replicas
					dep_list_for_hpa = append(dep_list_for_hpa, cluster.Name)
				}
			}

			//리스트에 없는 클러스터들은 HPA, VPA 삭제 (확인 작업)
			var dep_list_except []string
			for _, cluster := range cm.Cluster_list.Items {
				check := 0
				for _, targetCluster := range dep_list_for_hpa {
					if cluster.Name == targetCluster {
						check = 1
						break
					}
				}
				if check == 0 {
					dep_list_except = append(dep_list_except, cluster.Name)
				}
			}
			if dep_list_except != nil {
				//var sync_name string
				for _, cluster := range dep_list_except {
					_ = r.DeleteHPAVPA(cm, cluster, request.Namespace, request.Name)

					err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
					if err_openmcp != nil {
						omcplog.V(0).Info("!!! Failed to update instance status", err_openmcp)
						return reconcile.Result{}, err_openmcp
					} else {
						omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (HPA Create)")
					}
				}

			}

			omcplog.V(3).Info(">>> Target Clusters ", dep_list_for_hpa, " except ", dep_list_except)

			if dep_list_for_hpa != nil {
				// min,max 분배
				cluster_min_map, cluster_max_map, min_max_err := r.UpdateMinMaxDistributionPolicy(hasInstance, cluster_dep_request, dep_list_for_hpa, cluster_dep_replicas)
				if min_max_err != nil {
					omcplog.V(0).Info(min_max_err)
					return reconcile.Result{}, min_max_err
				}
				var sync_name string
				for _, clustername := range dep_list_for_hpa {
					// case 1) HPA 생성
					if hasInstance.Spec.VpaMode == "Never" || (hasInstance.Spec.VpaMode == "Auto" && cluster_dep_request[clustername] == true) {
						foundHPA := &hpav2beta2.HorizontalPodAutoscaler{}
						cluster_client := cm.Cluster_genClients[clustername]
						err = cluster_client.Get(context.TODO(), foundHPA, hasInstance.Namespace, hasInstance.Name)

						// CREATE HPA
						if err != nil && errors.IsNotFound(err) {

							hpa_min := cluster_min_map[clustername]
							hpa_max := cluster_max_map[clustername]

							//HPA 오브젝트 생성
							hpa := r.UpdateHorizontalPodAutoscaler(request, hasInstance, hpa_min, hpa_max)
							//Sync 리소스 생성 (Sync 컨트롤러가 Watch)
							command := "create"
							sync_name, err = r.sendSyncHPA(hpa, command, clustername)

							hasInstance.Status.SyncRequestName = sync_name

							if err != nil {
								omcplog.V(0).Info("!!! Failed to create new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
								//reqLogger.Error(err, "!!! Failed to create new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
								return reconcile.Result{}, err
							}
							omcplog.V(2).Info(">>> "+clustername+" Create HPA [ min:", *hpa.Spec.MinReplicas, " / max:", hpa.Spec.MaxReplicas, " ]")
							//reqLogger.Info("Creating a new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)

							//Status Update
							hasInstance.Status.LastSpec = hasInstance.Spec
							tmpMap[clustername] = 0
							hasInstance.Status.RebalancingCount = tmpMap

							err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
							if err_openmcp != nil {
								omcplog.V(0).Info("!!! Failed to update instance status", err_openmcp)
								return reconcile.Result{}, err_openmcp
							} else {
								omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (HPA Create)")
							}

							// HPA created successfully - return and requeue
							//return reconcile.Result{Requeue: true}, nil
						} else if err != nil {
							omcplog.V(0).Info("!!! Failed to get HPA")
							//reqLogger.Error(err, "!!! Failed to get HPA")
							return reconcile.Result{}, err
						} else {
							if hasInstance.Status.LastSpec.HpaTemplate.Spec.MinReplicas != nil {
								if *hasInstance.Status.LastSpec.HpaTemplate.Spec.MinReplicas != *hasInstance.Spec.HpaTemplate.Spec.MinReplicas || hasInstance.Status.LastSpec.HpaTemplate.Spec.MaxReplicas != hasInstance.Spec.HpaTemplate.Spec.MaxReplicas {

									desired_min_replicas := cluster_min_map[clustername]
									desired_max_replicas := cluster_max_map[clustername]

									if *foundHPA.Spec.MinReplicas != desired_min_replicas || foundHPA.Spec.MaxReplicas != desired_max_replicas {
										foundHPA.TypeMeta.Kind = "HorizontalPodAutoscaler"
										foundHPA.TypeMeta.APIVersion = "autoscaling/v2beta2"

										foundHPA.Spec.MinReplicas = &desired_min_replicas
										foundHPA.Spec.MaxReplicas = desired_max_replicas

										command := "update"
										_, err = r.sendSyncHPA(foundHPA, command, clustername)

										omcplog.V(2).Info(">>> "+clustername+" Update HPA [ min:", *foundHPA.Spec.MinReplicas, " / max:", foundHPA.Spec.MaxReplicas, " ]")

										if err != nil {
											omcplog.V(0).Info("!!! Failed to update HPA", "Hpa.Namespace", foundHPA.Namespace, "Hpa.Name", foundHPA.Name)
											//reqLogger.Error(err, "!!! Failed to update HPA", "Hpa.Namespace", foundHPA.Namespace, "Hpa.Name", foundHPA.Name)
											return reconcile.Result{}, err
										}
									}
								}
							}
						}
					// case 2) HPA + VPA 생성
					} else if hasInstance.Spec.VpaMode == "Always" || (hasInstance.Spec.VpaMode == "Auto" && cluster_dep_request[clustername] == false) {
						// 먼저 HPA 생성
						foundHPA := &hpav2beta2.HorizontalPodAutoscaler{}
						cluster_client := cm.Cluster_genClients[clustername]
						err = cluster_client.Get(context.TODO(), foundHPA, hasInstance.Namespace, hasInstance.Name)
						if err != nil && errors.IsNotFound(err) { //CREATE HPA

							hpa_min := cluster_min_map[clustername]
							hpa_max := cluster_max_map[clustername]

							//HPA 오브젝트 생성
							hpa := r.UpdateHorizontalPodAutoscaler(request, hasInstance, hpa_min, hpa_max)
							//Sync 리소스 생성 (Sync 컨트롤러가 Watch)
							command := "create"
							sync_name, err = r.sendSyncHPA(hpa, command, clustername)

							hasInstance.Status.SyncRequestName = sync_name

							if err != nil {
								omcplog.V(0).Info("!!! Failed to create new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
								//reqLogger.Error(err, "!!! Failed to create new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
								return reconcile.Result{}, err
							}
							omcplog.V(2).Info(">>> "+clustername+" Create HPA [ min:", *hpa.Spec.MinReplicas, " / max:", hpa.Spec.MaxReplicas, " ]")
							//reqLogger.Info("Creating a new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)

							//Status Update
							hasInstance.Status.LastSpec = hasInstance.Spec
							tmpMap[clustername] = 0
							hasInstance.Status.RebalancingCount = tmpMap

							err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
							if err_openmcp != nil {
								omcplog.V(0).Info("!!! Failed to update instance status", err_openmcp)
								return reconcile.Result{}, err_openmcp
							} else {
								omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (HPA Create)")
							}
						} else if err != nil {
							omcplog.V(0).Info("!!! Failed to get HPA")
							//reqLogger.Error(err, "!!! Failed to get HPA")
							return reconcile.Result{}, err
						} else {
							if hasInstance.Status.LastSpec.HpaTemplate.Spec.MinReplicas != nil {
								if *hasInstance.Status.LastSpec.HpaTemplate.Spec.MinReplicas != *hasInstance.Spec.HpaTemplate.Spec.MinReplicas || hasInstance.Status.LastSpec.HpaTemplate.Spec.MaxReplicas != hasInstance.Spec.HpaTemplate.Spec.MaxReplicas {

									desired_min_replicas := cluster_min_map[clustername]
									desired_max_replicas := cluster_max_map[clustername]

									if *foundHPA.Spec.MinReplicas != desired_min_replicas || foundHPA.Spec.MaxReplicas != desired_max_replicas {
										foundHPA.TypeMeta.Kind = "HorizontalPodAutoscaler"
										foundHPA.TypeMeta.APIVersion = "autoscaling/v2beta2"

										foundHPA.Spec.MinReplicas = &desired_min_replicas
										foundHPA.Spec.MaxReplicas = desired_max_replicas

										command := "update"
										_, err = r.sendSyncHPA(foundHPA, command, clustername)

										omcplog.V(2).Info(">>> "+clustername+" Update HPA [ min:", *foundHPA.Spec.MinReplicas, " / max:", foundHPA.Spec.MaxReplicas, " ]")

										if err != nil {
											omcplog.V(0).Info("!!! Failed to update HPA", "Hpa.Namespace", foundHPA.Namespace, "Hpa.Name", foundHPA.Name)
											//reqLogger.Error(err, "!!! Failed to update HPA", "Hpa.Namespace", foundHPA.Namespace, "Hpa.Name", foundHPA.Name)
											return reconcile.Result{}, err
										}
									}
								}
							}
						}
						//그 다음 VPA 생성
						foundVPA := &vpav1beta2.VerticalPodAutoscaler{}
						err = r.ghosts[clustername].Get(context.TODO(), ObjectKey{Namespace: hasInstance.Namespace, Name: hasInstance.Name}, foundVPA)

						if err != nil && errors.IsNotFound(err) { //CREATE VPA

							vpa := r.UpdateVerticalPodAutoscaler(request, hasInstance)

							command := "create"
							_, err = r.sendSyncVPA(vpa, command, clustername)

							hasInstance.Status.SyncRequestName = sync_name

							//reqLogger.Info("Creating a new VPA", "VPA.Namespace", vpa.Namespace, "VPA.Name", vpa.Name)
							if err != nil {
								omcplog.V(0).Info("!!! Failed to create new VPA", "VPA.Namespace", vpa.Namespace, "VPA.Name", vpa.Name)
								//reqLogger.Error(err, "!!! Failed to create new VPA", "VPA.Namespace", vpa.Namespace, "VPA.Name", vpa.Name)
								return reconcile.Result{}, err
							}
							// VPA created successfully - return and requeue
							omcplog.V(2).Info(">>> " + clustername + " Create VPA")

							hasInstance.Status.LastSpec = hasInstance.Spec

							err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
							if err_openmcp != nil {
								omcplog.V(0).Info("!!! Failed to update instance status", err_openmcp)
								return reconcile.Result{}, err_openmcp
							} else {
								omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (VPA Create)")
							}

						} else if err != nil {
							omcplog.V(0).Info("!!! Failed to get VPA")
							//reqLogger.Error(err, "!!! Failed to get VPA")
							return reconcile.Result{}, err
						} else if err == nil { //UPDATE VPA
							//fmt.Println(">>> " + clustername + " Update VPA")
							/*vpa := r.UpdateVerticalPodAutoscaler(request, hasInstance)
							fmt.Println("--- [" + clustername + "] UPDATE VPA SUCCESS\n")
							err = cluster_client.Update(context.TODO(), vpa)
							if err != nil {
								reqLogger.Error(err, "Failed to create new VPA", "VPA.Namespace", vpa.Namespace, "VPA.Name", vpa.Name)
								return reconcile.Result{}, err
							}*/
						}
					}
				}

				//OpenMCPHPA 리소스 변경 여부 확인을 위한 변수 저장
				hasInstance.Status.LastSpec = hasInstance.Spec

				err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
				if err_openmcp != nil {
					omcplog.V(0).Info("!!! Failed to update instance status", err_openmcp)
					return reconcile.Result{}, err_openmcp
				} else {
					omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (End)")
				}
			}
		} else if err != nil && errors.IsNotFound(err) {
			omcplog.V(0).Info("!!! OpenmcpDeployment is not found")
			//reqLogger.Error(err, "!!! OpenMCPDeployment doesn't exist - ", "openmcpDep.Namespace: ", openmcpDep.Namespace, ", openmcpDep.Name: ", openmcpDep.Name)
			return reconcile.Result{}, err
		} else {
			omcplog.V(0).Info("!!! Failed to get OpenMCPDeployment")
			//reqLogger.Error(err, "!!! Failed to get OpenMCPDeployment")
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	/*sync_instance := &sync.Sync{}
	nsn := types.NamespacedName{
		"openmcp",
		hasInstance.Status.SyncRequestName,
	}
	err = r.live.Get(context.TODO(), nsn, sync_instance)
	if err == nil {
		// 아직 Sync에서 처리되지 않음
		return reconcile.Result{}, nil
	}*/

	//타겟 OpenMCPDeployment Get
	openmcpDep := &ketiv1alpha1.OpenMCPDeployment{}
	openmcpDep_err := r.live.Get(context.TODO(), ObjectKey{Namespace: hasInstance.Namespace, Name: hasInstance.Spec.HpaTemplate.Spec.ScaleTargetRef.Name}, openmcpDep)

	var dep_list_for_hpa []string
	var cluster_dep_replicas map[string]int32
	var cluster_dep_request map[string]bool

	if openmcpDep_err == nil {
		cluster_dep_replicas = make(map[string]int32)
		cluster_dep_request = make(map[string]bool)

		//HPA, VPA가 배포될 클러스터 리스트 저장 (has-target-cluster 정책 적용 및 Deployment 배포 확인)
		for _, cluster := range  cm.Cluster_list.Items {
			dep := &appsv1.Deployment{}
			cluster_client := cm.Cluster_genClients[cluster.Name]
			dep_err := cluster_client.Get(context.TODO(), dep, hasInstance.Namespace, hasInstance.Spec.HpaTemplate.Spec.ScaleTargetRef.Name)
			if dep_err == nil {
				//fmt.Println("deploy request: ", dep.Spec.Template.Spec.Containers[0].Resources.Requests) // 컨테이너 여러개일 경우?
				if dep.Spec.Template.Spec.Containers[0].Resources.Requests == nil {
					cluster_dep_request[cluster.Name] = false
				} else {
					cluster_dep_request[cluster.Name] = true
				}
				cluster_dep_replicas[cluster.Name] = *dep.Spec.Replicas
				dep_list_for_hpa = append(dep_list_for_hpa, cluster.Name)
			}
		}
	}

	for _, clustername := range dep_list_for_hpa {

		foundHPA := &hpav2beta2.HorizontalPodAutoscaler{}
		cluster_client := cm.Cluster_genClients[clustername]
		err = cluster_client.Get(context.TODO(), foundHPA, hasInstance.Namespace, hasInstance.Name)

		//UPDATE HPA (Rebalancing 수행 시 or min/max 수정 시)
		if foundHPA.Spec.MaxReplicas == foundHPA.Status.CurrentReplicas && foundHPA.Status.CurrentReplicas == foundHPA.Status.DesiredReplicas {
			if lastTimeRebalancing.IsZero() || (!lastTimeRebalancing.IsZero() && time.Since(lastTimeRebalancing) > time.Second*180) {
				omcplog.V(2).Info(">>> " + clustername + " max rebalancing")

				hasInstance = r.MaxRebalancing(cm, hasInstance, dep_list_for_hpa, clustername, foundHPA)

				err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
				if err_openmcp != nil {
					omcplog.V(0).Info("!!! Failed to update instance status \"RebalancingCount\"", err_openmcp)
					return reconcile.Result{}, err_openmcp
				} else {
					omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (RebalancingCount)")
				}

			}
		} else if *foundHPA.Spec.MinReplicas > 1 && *foundHPA.Spec.MinReplicas == foundHPA.Status.CurrentReplicas && foundHPA.Status.CurrentReplicas == foundHPA.Status.DesiredReplicas {
			if lastTimeRebalancing.IsZero() || (!lastTimeRebalancing.IsZero() && time.Since(lastTimeRebalancing) > time.Second*180) {
				omcplog.V(2).Info(">>> " + clustername + " min rebalancing")

				hasInstance = r.MinRebalancing(cm, hasInstance, dep_list_for_hpa, clustername, foundHPA)

				err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
				if err_openmcp != nil {
					omcplog.V(0).Info("!!! Failed to update instance status \"RebalancingCount\"", err_openmcp)
					return reconcile.Result{}, err_openmcp
				} else {
					omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (RebalancingCount)")
				}

			}
		}
	}


	return reconcile.Result{}, nil

}

func (r *reconciler) UpdateMinMaxDistributionPolicy(hasInstance *ketiv1alpha1.OpenMCPHybridAutoScaler, cluster_dep_request map[string]bool, dep_list_for_hpa []string, cluster_dep_replicas map[string]int32) (map[string]int32, map[string]int32, error){
	type ObjectKey = types.NamespacedName

	foundPolicy := &ketiv1alpha1.OpenMCPPolicy{}
	minmax_policy_err := r.live.Get(context.TODO(), ObjectKey{Namespace: "openmcp", Name: "hpa-minmax-distribution-mode"}, foundPolicy)

	cluster_min_map, cluster_max_map, hasInstance, min_max_err := r.CreateMinMaxMap(hasInstance, cluster_dep_request, foundPolicy, minmax_policy_err, dep_list_for_hpa, cluster_dep_replicas)

	return cluster_min_map, cluster_max_map, min_max_err
}

func (r *reconciler) UpdateTargetClusterPolicy(cm *clusterManager.ClusterManager, hasInstance *ketiv1alpha1.OpenMCPHybridAutoScaler) ([]fedv1b1.KubeFedCluster, *ketiv1alpha1.OpenMCPHybridAutoScaler){
	checkPolicy := 0
	clusterListItems := make([]fedv1b1.KubeFedCluster, 0)

	type ObjectKey = types.NamespacedName

	if hasInstance.Status.Policies != nil { // 기존에 설정된 정책이 있는지 확인
		for _, tmp := range hasInstance.Status.Policies {
			if tmp.Type == "Target" { // Target 이미 설정되어 있으면
				//fmt.Println(">>> Policy \"Cluster Target\" Existed")
				if tmp.Value[0] == "Default" {
					clusterListItems = cm.Cluster_list.Items
				} else {
					for _, cluster := range cm.Cluster_list.Items {
						for _, value := range tmp.Value {
							if cluster.Name == value {
								clusterListItems = append(clusterListItems, cluster)
							}
						}
					}
				}
				checkPolicy = 1
				break
			}
		}
	}
	if checkPolicy == 0 { // Target 설정 X
		foundPolicy := &ketiv1alpha1.OpenMCPPolicy{}
		target_cluster_policy_err := r.live.Get(context.TODO(), ObjectKey{Namespace: "openmcp", Name: "has-target-cluster"}, foundPolicy)
		if target_cluster_policy_err == nil {
			if foundPolicy.Spec.PolicyStatus == "Enabled" {
				omcplog.V(3).Info(">>> Policy \"Cluster Target\" Apply (Enabled)")
				hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				for _, cluster := range cm.Cluster_list.Items {
					for _, value := range foundPolicy.Spec.Template.Spec.Policies[0].Value {
						if cluster.Name == value {
							clusterListItems = append(clusterListItems, cluster)
						}
					}
				}
			} else {
				omcplog.V(3).Info(">>> Policy \"Cluster Target\" Apply (Disabled - set default)")
				omp := make([]ketiv1alpha1.OpenMCPPolicies, 1)
				omp[0].Type = "Target"
				omp_value := make([]string, 0)
				omp_value = append(omp_value, "Default")
				omp[0].Value = omp_value
				hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
				clusterListItems = cm.Cluster_list.Items
			}
		} else {
			omcplog.V(1).Info("!!! Fail to get policy \"Cluster Target\" (set default)")
			//klog.V(0).Info("policy_err : ", target_cluster_policy_err)
			omp := make([]ketiv1alpha1.OpenMCPPolicies, 1)
			omp[0].Type = "Target"
			omp_value := make([]string, 0)
			omp_value = append(omp_value, "Default")
			omp[0].Value = omp_value
			hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
			clusterListItems = cm.Cluster_list.Items
		}
	}

	return clusterListItems, hasInstance
}

func (r *reconciler) CreateMinMaxMap(hasInstance *ketiv1alpha1.OpenMCPHybridAutoScaler, cluster_dep_request map[string]bool,foundPolicy *ketiv1alpha1.OpenMCPPolicy, minmax_policy_err error, dep_list_for_hpa []string, cluster_dep_replicas map[string]int32) (map[string]int32, map[string]int32, *ketiv1alpha1.OpenMCPHybridAutoScaler, error) {
	timeStart_mixmaxdist := time.Now()

	cluster_min_map := make(map[string]int32)
	cluster_max_map := make(map[string]int32)
	var err error

	type ObjectKey = types.NamespacedName

	checkPolicy := 0
	if hasInstance.Status.Policies != nil { // 정책이 하나라도 있는지 확인
		for n, tmp := range hasInstance.Status.Policies {
			if tmp.Type == "Mode" { // Mode 이미 설정되어 있으면
				//klog.V(0).Info(">>> Policy \"Min Max Distribution\" Existed")
				if hasInstance.Status.Policies[n].Value[0] == "Equal" {
					//fmt.Println("Policy - min max Equal")
					cluster_min_map = HpaMinScheduling_equal(dep_list_for_hpa, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling_equal(dep_list_for_hpa, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
				} else {
					//fmt.Println("Policy - min max Unequal")
					cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)

					for _, cluster := range dep_list_for_hpa {
						if hasInstance.Spec.VpaMode == "Always" || (hasInstance.Spec.VpaMode == "Auto" && cluster_dep_request[cluster] == false){
							if cluster_min_map[cluster] < 2 {
								cluster_min_map[cluster] = 2
							}
						}
					}
				}
				checkPolicy = 1
				break
			}
		}
	}
	if checkPolicy == 0 { // Mode 설정 X

		if minmax_policy_err == nil {
			if foundPolicy.Spec.PolicyStatus == "Enabled" {
				omcplog.V(3).Info(">>> Policy \"Min Max Distribution\" Apply (Enabled)")
				if foundPolicy.Spec.Template.Spec.Policies[0].Value[0] == "Equal" {
					//fmt.Println("New Policy - min max Equal")
					cluster_min_map = HpaMinScheduling_equal(dep_list_for_hpa, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling_equal(dep_list_for_hpa, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
					hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				} else if foundPolicy.Spec.Template.Spec.Policies[0].Value[0] == "Unequal" {
					//fmt.Println("New Policy - min max Unequal")
					cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
					hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				} else {
					//fmt.Println("Default - min max Unequal")
					cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
					foundPolicy.Spec.Template.Spec.Policies[0].Value[0] = "Default"
					hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				}

			} else {
				omcplog.V(3).Info(">>> Policy \"Min Max Distribution\" Apply (Disabled - set default)")
				cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
				cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
				omp := make([]ketiv1alpha1.OpenMCPPolicies, 1)
				omp_value := make([]string, 1)
				omp_value[0] = "Default"
				omp[0].Type = "Mode"
				omp[0].Value = omp_value
				hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
			}
		} else {
			//klog.V(0).Info("policy_err : ", minmax_policy_err)
			omcplog.V(1).Info("Fail to get policy \"Min Max Distribution\" (set default)")
			cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
			cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
			omp := make([]ketiv1alpha1.OpenMCPPolicies, 1)
			omp_value := make([]string, 1)
			omp_value[0] = "Default"
			omp[0].Type = "Mode"
			omp[0].Value = omp_value
			hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
		}

		for _, cluster := range dep_list_for_hpa {
			if hasInstance.Spec.VpaMode == "Always" || (hasInstance.Spec.VpaMode == "Auto" && cluster_dep_request[cluster] == false){
				if cluster_min_map[cluster] < 2 {
					cluster_min_map[cluster] = 2
				}
			}
		}

		omcplog.V(4).Info("******* [Start] HAS Min/Max Distribution *******")
		timeEnd_mixmaxdist := time.Since(timeStart_mixmaxdist)
		for _, cluster := range dep_list_for_hpa {
			omcplog.V(4).Info("[",cluster,"] min:",cluster_min_map[cluster]," / max:",cluster_max_map[cluster],"")
		}
		omcplog.V(4).Info("=> Total Min/Max Distribution Time [", timeEnd_mixmaxdist,"]")
		omcplog.V(4).Info("*******  [End] HAS Min/Max Distribution  *******")


		//정책 업데이트
		err = r.live.Status().Update(context.TODO(), hasInstance)
		if err != nil {
			omcplog.V(0).Info("Policy Status Update Error")
		} else {
			omcplog.V(3).Info(">>> Policy Status UPDATE Success")
		}

	}

	return cluster_min_map, cluster_max_map, hasInstance, err
}

func (r *reconciler) MinRebalancing(cm *clusterManager.ClusterManager, hasInstance *ketiv1alpha1.OpenMCPHybridAutoScaler, dep_list_for_hpa []string, clustername string, foundHPA *hpav2beta2.HorizontalPodAutoscaler) *ketiv1alpha1.OpenMCPHybridAutoScaler {

	var dep_list_for_analysis []string

	minReplicasMap := map[string]int32{}
	currentReplicasMap := map[string]int32{}

	for _, cn := range dep_list_for_hpa {
		analysisHPA := &hpav2beta2.HorizontalPodAutoscaler{}
		err := cm.Cluster_genClients[cn].Get(context.TODO(), analysisHPA, hasInstance.Namespace, hasInstance.Name)
		if err == nil {
			//fmt.Println("min : ", *analysisHPA.Spec.MinReplicas)
			if *analysisHPA.Spec.MinReplicas < analysisHPA.Status.CurrentReplicas-1 {
				dep_list_for_analysis = append(dep_list_for_analysis, cn)
				minReplicasMap[cn] = *analysisHPA.Spec.MinReplicas
				currentReplicasMap[cn] = analysisHPA.Status.CurrentReplicas
			}
		} else {
			omcplog.V(0).Info("!!! Failed to get hpa info : ", err)
		}
	}

	//후보 클러스터가 없을 때
	if len(dep_list_for_analysis) == 0 {
		omcplog.V(0).Info("!!! Failed Rebalancing : There is no candidate cluster")
	} else {
		//timeStart_mixmaxrebal := time.Now()
		//분석 엔진을 통해 얻은 결과 (gRPC 통신)
		//---------------------------------------------------------------------------------------------------
		SERVER_IP := os.Getenv("GRPC_SERVER")
		SERVER_PORT := os.Getenv("GRPC_PORT")

		grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

		hi := &protobuf.HASInfo{HPAName: hasInstance.Name, HPANamespace: hasInstance.Namespace, ClusterName: clustername, HPACurrentReplicas: currentReplicasMap, HPAMinORMaxReplicas: minReplicasMap, HASRebalancingCount: hasInstance.Status.RebalancingCount}
		timeStart_mixmaxrebal := time.Now()
		result, gRPCerr := grpcClient.SendHASMinAnalysis(context.TODO(), hi)
		timeEnd_mixmaxrebal := time.Since(timeStart_mixmaxrebal)
		if gRPCerr != nil || len(result.TargetCluster) == 0 {
			if gRPCerr != nil {
				omcplog.V(0).Info("could not connect : %v", gRPCerr)
			}
			omcplog.V(0).Info("!!! Failed Rebalacing : Failed to get Analysis Result :(")
		} else {
			//---------------------------------------------------------------------------------------------------
			qosCluster := result.TargetCluster
			omcplog.V(3).Info("     => Anlysis Result [", qosCluster, "]")


			omcplog.V(4).Info("     => Total Rebalancing Time [", timeEnd_mixmaxrebal, "] (Analysis + gRPC)")
			//fmt.Println("-------------------- Min rebalancing result 분석 시간 : ", timeEnd_mixmaxrebal,"------------------")

			foundHPA.TypeMeta.Kind = "HorizontalPodAutoscaler"
			foundHPA.TypeMeta.APIVersion = "autoscaling/v2beta2"

			updateHPA := foundHPA
			*updateHPA.Spec.MinReplicas -= 1

			qos_cluster_client := cm.Cluster_genClients[qosCluster]
			foundQosHPA := &hpav2beta2.HorizontalPodAutoscaler{}
			err := qos_cluster_client.Get(context.TODO(), foundQosHPA, hasInstance.Namespace, hasInstance.Name)

			if err != nil {
				omcplog.V(0).Info("err: ", err)
			}

			foundQosHPA.TypeMeta.Kind = "HorizontalPodAutoscaler"
			foundQosHPA.TypeMeta.APIVersion = "autoscaling/v2beta2"

			updateQosHPA := foundQosHPA
			*updateQosHPA.Spec.MinReplicas += 1

			command := "update"
			_, current_err := r.sendSyncHPA(updateHPA, command, clustername)
			sync_name, qos_err := r.sendSyncHPA(updateQosHPA, command, qosCluster)

			hasInstance.Status.SyncRequestName = sync_name

			if current_err != nil || qos_err != nil { //둘 중 하나라도 에러가 날 경우 롤백
				omcplog.V(0).Info("current_err : ", current_err)
				omcplog.V(0).Info("qos_err : ", qos_err)

				command := "update"
				r.sendSyncHPA(foundHPA, command, clustername)
				r.sendSyncHPA(foundQosHPA, command, qosCluster)

				omcplog.V(3).Info(">>> "+clustername+" Rollback HPA [ min:", *foundHPA.Spec.MinReplicas, " / max:", foundHPA.Spec.MaxReplicas, " ]")
				omcplog.V(3).Info(">>> "+qosCluster+" Rollback HPA [ min:", *foundQosHPA.Spec.MinReplicas, " / max:", foundQosHPA.Spec.MaxReplicas, " ]")
			} else if current_err == nil && qos_err == nil {
				omcplog.V(3).Info(">>> "+clustername+" Update HPA [ min:", *updateHPA.Spec.MinReplicas, " / max:", updateHPA.Spec.MaxReplicas, " ]")
				omcplog.V(3).Info(">>> "+qosCluster+" Update HPA [ min:", *updateQosHPA.Spec.MinReplicas, " / max:", updateQosHPA.Spec.MaxReplicas, " ]")

				//Status Update
				hasInstance.Status.RebalancingCount[clustername] += 1

				lastTimeRebalancing = time.Now()
				omcplog.V(4).Info("     => RebalancingTime : ", lastTimeRebalancing)
			}
		}
	}

	return hasInstance
}

func (r *reconciler) MaxRebalancing(cm *clusterManager.ClusterManager, hasInstance *ketiv1alpha1.OpenMCPHybridAutoScaler, dep_list_for_hpa []string, clustername string, foundHPA *hpav2beta2.HorizontalPodAutoscaler) *ketiv1alpha1.OpenMCPHybridAutoScaler {

	var dep_list_for_analysis []string

	maxReplicasMap := map[string]int32{}
	currentReplicasMap := map[string]int32{}

	for _, cn := range dep_list_for_hpa {
		analysisHPA := &hpav2beta2.HorizontalPodAutoscaler{}
		err := cm.Cluster_genClients[cn].Get(context.TODO(), analysisHPA, hasInstance.Namespace, hasInstance.Name)

		if err == nil {
			if analysisHPA.Spec.MaxReplicas > analysisHPA.Status.CurrentReplicas+1 {
				dep_list_for_analysis = append(dep_list_for_analysis, cn)
				maxReplicasMap[cn] = analysisHPA.Spec.MaxReplicas
				currentReplicasMap[cn] = analysisHPA.Status.CurrentReplicas
			}
		} else {
			omcplog.V(0).Info("!!! Failed to get hpa info : ", err)
		}
	}

	//후보 클러스터가 없을 때
	if len(dep_list_for_analysis) == 0 {
		omcplog.V(0).Info("!!! Failed Rebalancing : There is no candidate cluster")
	} else {

		//분석 엔진을 통해 얻은 결과 (gRPC 통신)
		//---------------------------------------------------------------------------------------------------
		SERVER_IP := os.Getenv("GRPC_SERVER")
		SERVER_PORT := os.Getenv("GRPC_PORT")

		grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

		hi := &protobuf.HASInfo{HPAName: hasInstance.Name, HPANamespace: hasInstance.Namespace, ClusterName: clustername, HPACurrentReplicas: currentReplicasMap, HPAMinORMaxReplicas: maxReplicasMap, HASRebalancingCount: hasInstance.Status.RebalancingCount}
		timeStart_mixmaxrebal := time.Now()
		result, gRPCerr := grpcClient.SendHASMaxAnalysis(context.TODO(), hi)
		timeEnd_mixmaxrebal := time.Since(timeStart_mixmaxrebal)
		if gRPCerr != nil || len(result.TargetCluster) == 0 {
			if gRPCerr != nil {
				omcplog.V(0).Info("could not connect : %v", gRPCerr)
			}
			omcplog.V(0).Info("!!! Failed Rebalacing : Failed to get Analysis Result :(")
		} else {
			//---------------------------------------------------------------------------------------------------
			qosCluster := result.TargetCluster
			omcplog.V(3).Info("     => Anlysis Result [", qosCluster, "]")


			omcplog.V(4).Info("     => Total Rebalancing Time [", timeEnd_mixmaxrebal, "] (Analysis + gRPC)")

			foundHPA.TypeMeta.Kind = "HorizontalPodAutoscaler"
			foundHPA.TypeMeta.APIVersion = "autoscaling/v2beta2"

			updateHPA := foundHPA
			updateHPA.Spec.MaxReplicas += 1

			qos_cluster_client := cm.Cluster_genClients[qosCluster]
			foundQosHPA := &hpav2beta2.HorizontalPodAutoscaler{}
			err := qos_cluster_client.Get(context.TODO(), foundQosHPA, hasInstance.Namespace, hasInstance.Name)

			if err != nil {
				omcplog.V(0).Info("err: ", err)
			}

			foundQosHPA.TypeMeta.Kind = "HorizontalPodAutoscaler"
			foundQosHPA.TypeMeta.APIVersion = "autoscaling/v2beta2"

			updateQosHPA := foundQosHPA
			updateQosHPA.Spec.MaxReplicas -= 1

			command := "update"
			_, current_err := r.sendSyncHPA(updateHPA, command, clustername)
			sync_name, qos_err := r.sendSyncHPA(updateQosHPA, command, qosCluster)

			hasInstance.Status.SyncRequestName = sync_name

			if current_err != nil || qos_err != nil { //둘 중 하나라도 에러가 날 경우 롤백
				omcplog.V(0).Info("current_err : ", current_err)
				omcplog.V(0).Info("qos_err : ", qos_err)

				command := "update"
				r.sendSyncHPA(foundHPA, command, clustername)
				r.sendSyncHPA(foundQosHPA, command, qosCluster)

				omcplog.V(3).Info(">>> "+clustername+" Rollback HPA [ min:", *foundHPA.Spec.MinReplicas, " / max:", foundHPA.Spec.MaxReplicas, " ]")
				omcplog.V(3).Info(">>> "+qosCluster+" Rollback HPA [ min:", *foundQosHPA.Spec.MinReplicas, " / max:", foundQosHPA.Spec.MaxReplicas, " ]")
			} else if current_err == nil && qos_err == nil {
				omcplog.V(3).Info(">>> "+clustername+" Update HPA [ min:", *updateHPA.Spec.MinReplicas, " / max:", updateHPA.Spec.MaxReplicas, " ]")
				omcplog.V(3).Info(">>> "+qosCluster+" Update HPA [ min:", *updateQosHPA.Spec.MinReplicas, " / max:", updateQosHPA.Spec.MaxReplicas, " ]")

				//Status Update
				hasInstance.Status.RebalancingCount[clustername] += 1

				lastTimeRebalancing = time.Now()
				omcplog.V(4).Info("     => RebalancingTime : ", lastTimeRebalancing)

			}
		}
	}

	return hasInstance

}

func (r *reconciler) UpdateHorizontalPodAutoscaler(req reconcile.Request, m *ketiv1alpha1.OpenMCPHybridAutoScaler, min int32, max int32) *hpav2beta2.HorizontalPodAutoscaler {
	ls := LabelsForHpa(m.Name)

	hpa := &hpav2beta2.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2beta2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},

		Spec:   m.Spec.HpaTemplate.Spec,
		Status: hpav2beta2.HorizontalPodAutoscalerStatus{},
	}

	hpa.Spec.MinReplicas = &min
	hpa.Spec.MaxReplicas = max
	hpa.Spec.ScaleTargetRef.APIVersion = "apps/v1"
	hpa.Spec.ScaleTargetRef.Kind = "Deployment"
	hpa.Spec.ScaleTargetRef.Name = hpa.Spec.ScaleTargetRef.Name

	reference.SetMulticlusterControllerReference(hpa, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return hpa
}

func (r *reconciler) UpdateVerticalPodAutoscaler(req reconcile.Request, m *ketiv1alpha1.OpenMCPHybridAutoScaler) *vpav1beta2.VerticalPodAutoscaler {
	ls := LabelsForHpa(m.Name)
	vpaUpdateMode := vpav1beta2.UpdateModeAuto

	vpa := &vpav1beta2.VerticalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VerticalPodAutoscaler",
			APIVersion: "autoscaling.k8s.io/v1beta2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: vpav1beta2.VerticalPodAutoscalerSpec{
			TargetRef: &autoscaling.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       m.Spec.HpaTemplate.Spec.ScaleTargetRef.Name,
			},
			UpdatePolicy: &vpav1beta2.PodUpdatePolicy{
				UpdateMode: &vpaUpdateMode,
			},
		},
		Status: vpav1beta2.VerticalPodAutoscalerStatus{},
	}
	reference.SetMulticlusterControllerReference(vpa, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))
	// Set Memcached instance as the owner and controller
	//controllerutil.SetControllerReference(m, vpa, r.scheme)
	return vpa
}

func (r *reconciler) DeleteHPAVPA(cm *clusterManager.ClusterManager, cluster string, namespace string, name string) string{
	hpa := &hpav2beta2.HorizontalPodAutoscaler{}

	vpa := &vpav1beta2.VerticalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VerticalPodAutoscaler",
			APIVersion: "autoscaling.k8s.io/v1beta2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	type ObjectKey = types.NamespacedName
	var sync_name string

	cluster_client := cm.Cluster_genClients[cluster]

	err1 := cluster_client.Get(context.TODO(), hpa, namespace, name)

	if err1 != nil && errors.IsNotFound(err1) {
		//fmt.Println("Fail to Delete HPA - ", err1)
	} else if err1 != nil {
		omcplog.V(0).Info("!!! Fail to Delete HPA - ", err1)
	} else if err1 == nil {
		var sync_err1 error

		hpa.TypeMeta.Kind = "HorizontalPodAutoscaler"
		hpa.TypeMeta.APIVersion = "autoscaling/v2beta2"

		command := "delete"
		sync_name, sync_err1 = r.sendSyncHPA(hpa, command, cluster)

		if sync_err1 != nil && errors.IsNotFound(sync_err1) {
			//klog.V(0).Info("Fail to Delete VPA - ", err3)
		} else if sync_err1 != nil {
			omcplog.V(0).Info("!!! Fail to Delete VPA - ", sync_err1)
		} else if sync_err1 == nil {
			omcplog.V(2).Info(">>> " + cluster + " Delete HPA")
		}
	}

	err2 := r.ghosts[cluster].Get(context.TODO(), ObjectKey{Namespace: namespace, Name: name}, vpa)
	if err2 != nil && errors.IsNotFound(err2) {
		//fmt.Println("Fail to Get VPA - ", err2)
	} else if err2 != nil {
		omcplog.V(0).Info("!!! Fail to Get VPA - ", err2)
	} else if err2 == nil {
		var sync_err2 error

		command := "delete"
		sync_name, sync_err2 = r.sendSyncVPA(vpa, command, cluster)

		if sync_err2 != nil && errors.IsNotFound(sync_err2) {
			//klog.V(0).Info("Fail to Delete VPA - ", err3)
		} else if sync_err2 != nil {
			omcplog.V(0).Info("!!! Fail to Delete VPA - ", sync_err2)
		} else if sync_err2 == nil {
			omcplog.V(2).Info(">>> " + cluster + " Delete VPA")
		}
	}

	return sync_name

}

func HpaMinScheduling(dep_list_for_hpa []string, cluster_dep_replicas map[string]int32, min int32) map[string]int32 {

	cluster_minreplicas_map := make(map[string]int32)
	var replicas_sum int32

	for _, clustername := range dep_list_for_hpa {
		cluster_minreplicas_map[clustername] = 0

		replicas_sum += cluster_dep_replicas[clustername]
	}

	cluster_len := int32(len(dep_list_for_hpa))
	remain_min := min

	if min < cluster_len {
		for _, clustername := range dep_list_for_hpa {
			cluster_minreplicas_map[clustername] = 1
		}
	} else {
		for replicas_sum != 0 && remain_min != 0 {
			for _, clustername := range dep_list_for_hpa {
				if cluster_dep_replicas[clustername] > 0 {
					cluster_minreplicas_map[clustername] += 1
					cluster_dep_replicas[clustername] -= 1

					replicas_sum -= 1
					remain_min -= 1

					if remain_min == 0 {
						break
					}
				}
			}
		}
		for remain_min != 0 {
			for _, clustername := range dep_list_for_hpa {
				cluster_minreplicas_map[clustername] += 1
				remain_min -= 1

				if remain_min == 0 {
					break
				}
			}
		}

	}
	//fmt.Println("min_map: ", cluster_minreplicas_map)
	return cluster_minreplicas_map
}

func HpaMaxScheduling(dep_list_for_hpa []string, cluster_dep_replicas map[string]int32, max int32) map[string]int32 {

	cluster_maxreplicas_map := make(map[string]int32)
	var replicas_sum int32

	for _, clustername := range dep_list_for_hpa {
		cluster_maxreplicas_map[clustername] = 0

		replicas_sum += cluster_dep_replicas[clustername]
	}

	cluster_len := int32(len(dep_list_for_hpa))

	remain_max := max

	if max < cluster_len {
		for _, clustername := range dep_list_for_hpa {
			cluster_maxreplicas_map[clustername] = 2
		}
	} else {
		for replicas_sum != 0 && remain_max != 0 {
			for _, clustername := range dep_list_for_hpa {
				if cluster_dep_replicas[clustername] > 0 {
					cluster_maxreplicas_map[clustername] += 1
					cluster_dep_replicas[clustername] -= 1

					replicas_sum -= 1
					remain_max -= 1

					if remain_max == 0 {
						break
					}
				}
			}
		}

		for remain_max != 0 {
			for _, clustername := range dep_list_for_hpa {
				cluster_maxreplicas_map[clustername] += 1
				remain_max -= 1

				if remain_max == 0 {
					break
				}
			}
		}
	}
	//fmt.Println("max_map: ", cluster_maxreplicas_map)
	return cluster_maxreplicas_map
}

func HpaMinScheduling_equal(dep_list_for_hpa []string, min int32) map[string]int32 {

	cluster_minreplicas_map := make(map[string]int32)

	for _, clustername := range dep_list_for_hpa {
		cluster_minreplicas_map[clustername] = 0

		//replicas_sum += cluster_dep_replicas[clustername]
	}

	cluster_len := int32(len(dep_list_for_hpa))
	remain_min := min

	if min < cluster_len {
		for _, clustername := range dep_list_for_hpa {
			cluster_minreplicas_map[clustername] = 1
		}
	} else {
		for remain_min != 0 {
			for _, clustername := range dep_list_for_hpa {
				cluster_minreplicas_map[clustername] += 1
				remain_min -= 1

				if remain_min == 0 {
					break
				}
			}
		}

	}
	//fmt.Println("min_map: ", cluster_minreplicas_map)
	return cluster_minreplicas_map
}

func HpaMaxScheduling_equal(dep_list_for_hpa []string, max int32) map[string]int32 {

	cluster_maxreplicas_map := make(map[string]int32)

	for _, clustername := range dep_list_for_hpa {
		cluster_maxreplicas_map[clustername] = 0
	}

	cluster_len := int32(len(dep_list_for_hpa))

	remain_max := max

	if max < cluster_len {
		for _, clustername := range dep_list_for_hpa {
			cluster_maxreplicas_map[clustername] = 2
		}
	} else {
		for remain_max != 0 {
			for _, clustername := range dep_list_for_hpa {
				cluster_maxreplicas_map[clustername] += 1
				remain_max -= 1

				if remain_max == 0 {
					break
				}
			}
		}
	}
	//fmt.Println("max_map: ", cluster_maxreplicas_map)
	return cluster_maxreplicas_map
}

func LabelsForHpa(name string) map[string]string {
	return map[string]string{"app": "openmcphybridautoscaler", "openmcphybridautoscaler_cr": name}
}