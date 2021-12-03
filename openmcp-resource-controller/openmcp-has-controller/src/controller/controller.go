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

package openmcphas

import (
	"context"
	"fmt"
	"openmcp/openmcp/apis"
	policyv1alpha1 "openmcp/openmcp/apis/policy/v1alpha1"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	syncv1alpha1 "openmcp/openmcp/apis/sync/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-analytic-engine/src/protobuf"
	"openmcp/openmcp/openmcp-resource-controller/openmcp-has-controller/src/analyticResource"
	"openmcp/openmcp/util/clusterManager"
	"os"
	"reflect"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"admiralty.io/multicluster-controller/pkg/reference"
	vpav1beta2 "github.com/kubernetes/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	appsv1 "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v1"
	hpav2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"strconv"
	"time"

	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	cm = myClusterManager

	liveClient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}

	ghostClients := map[string]client.Client{}
	for _, ghost := range ghosts {
		ghostTmp, err := ghost.GetDelegatingClient()
		if err != nil {
			omcplog.V(4).Info("Error getting delegating client for ghost cluster [", ghost.Name, "]")
			//return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		} else {
			ghostClients[ghost.Name] = ghostTmp
		}
	}

	co := controller.New(&reconciler{live: liveClient, ghosts: ghostClients, ghostNamespace: ghostNamespace}, controller.Options{})

	//live.GetScheme() - apis scheme ADD
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	//live.GetScheme() - vpav1beta2 scheme ADD
	if err := vpav1beta2.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	//fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(context.TODO(), live, &resourcev1alpha1.OpenMCPHybridAutoScaler{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	for _, ghost := range ghosts {
		//fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(context.TODO(), ghost, &hpav2beta2.HorizontalPodAutoscaler{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}

	for _, ghost := range ghosts {
		//fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(context.TODO(), ghost, &vpav1beta2.VerticalPodAutoscaler{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}

	return co, nil
}

func (r *reconciler) sendSyncHPA(hpa *hpav2beta2.HorizontalPodAutoscaler, command string, clusterName string) (string, error) {
	syncIndex += 1

	s := &syncv1alpha1.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-hybridautoscaler-hpa-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: syncv1alpha1.SyncSpec{
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

	s := &syncv1alpha1.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openmcp-hybridautoscaler-vpa-sync-" + strconv.Itoa(syncIndex),
			Namespace: "openmcp",
		},
		Spec: syncv1alpha1.SyncSpec{
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

	i += 1

	omcplog.V(3).Info("********* [", i, "] *********")
	omcplog.V(3).Info("Namespace : ", request.Namespace, " | Name : ", request.Name, " | Context : ", request.Context)

	totalHASTimeStart1 := time.Now()

	type ObjectKey = types.NamespacedName

	hasInstance := &resourcev1alpha1.OpenMCPHybridAutoScaler{}
	err := r.live.Get(context.TODO(), request.NamespacedName, hasInstance)

	//Delete OpenMCPHAS resource
	if err != nil {
		if errors.IsNotFound(err) {

			r.DeleteAllHPAVPA(cm, request.Namespace, request.Name)

			hasList := &resourcev1alpha1.OpenMCPHybridAutoScalerList{}
			err_haslist := r.live.List(context.TODO(), hasList)

			if err_haslist == nil {
				for k, _ := range analyticResource.CPAInfoList {
					check := 0
					for _, has := range hasList.Items {
						cpaKey := analyticResource.CPAKey{}

						cpaKey.HASName = has.Name
						cpaKey.HASNamespace = has.Namespace

						if cpaKey == k {
							check = 1
						}
					}
					if check == 0 {
						delete(analyticResource.CPAInfoList, k)
						fmt.Println("Success to Delete CPA ", k)
					}
				}

			} else {
				fmt.Println(err_haslist)
			}
			//return 맞나? cluster 여러개 조인되어 있을 때 hpa, vpa 생성/삭제 다시 테스트
			return reconcile.Result{}, nil
		} else {
			return reconcile.Result{}, err
		}
	}

	if !reflect.DeepEqual(hasInstance.Status.LastSpec, hasInstance.Spec) || hasInstance.Status.ChangeNeed == true {
		omcplog.V(3).Info("hasInstance.Status.ChangeNeed : ", hasInstance.Status.ChangeNeed)
		//Apply "has-target-cluster" Policy
		clusterListItems, hasInstance := r.UpdateTargetClusterPolicy(cm, hasInstance)

		target_cluster_policy_err := r.live.Status().Update(context.TODO(), hasInstance)
		if target_cluster_policy_err != nil {
			omcplog.V(0).Info("!![error] Fail to Update Policy Status")
			return reconcile.Result{}, target_cluster_policy_err
		} else {
			//omcplog.V(3).Info(">>> Success to Update Policy Status")
		}

		//totalHASTimeEnd1 := time.Since(totalHASTimeStart1)
		//omcplog.V(4).Info("------ Check HAS Time (1) : ", totalHASTimeEnd1)
		//totalHASTimeStart2 := time.Now()

		if hasInstance.Spec.MainController == "OpenMCP" {
			cpaKey := analyticResource.CPAKey{}

			cpaKey.HASName = hasInstance.Name
			cpaKey.HASNamespace = hasInstance.Namespace

			_, prs := analyticResource.CPAInfoList[cpaKey]
			//fmt.Println("CPA Exists? ", prs, "(",cpaKey,")")
			//cpa 리스트에 없는 경우
			if prs == false {
				//타겟으로 하는 openmcpdeployment 배포되어있는지 검사
				openmcpDep := &resourcev1alpha1.OpenMCPDeployment{}
				openmcpDep_err := r.live.Get(
					context.TODO(),
					ObjectKey{
						Namespace: hasInstance.Namespace,
						Name:      hasInstance.Spec.ScalingOptions.CpaTemplate.ScaleTargetRef.Name},
					openmcpDep)

				var totalCpuRequest int64
				var totalMemRequest int64

				totalCpuRequest = 0
				totalMemRequest = 0

				for _, container := range openmcpDep.Spec.Template.Spec.Template.Spec.Containers {

					cpuInt64 := container.Resources.Requests.Cpu().MilliValue()
					memInt64 := container.Resources.Requests.Memory().MilliValue()

					if cpuInt64 > 0 {
						totalCpuRequest += cpuInt64
					}
					if memInt64 > 0 {
						totalMemRequest += memInt64
					}
				}

				if totalCpuRequest == 0 || totalMemRequest == 0 {
					omcplog.V(0).Info("!![error] Fail to CPA - There is no request on Deployment")
				} else {
					if openmcpDep_err == nil {
						cpaCluster := make([]string, 0)
						var totalReplicas int32
						totalReplicas = 0

						//어느 클러스터에 배포되어있는지 확인
						for _, cluster := range cm.Cluster_list.Items {
							dep := &appsv1.Deployment{}
							cluster_client := cm.Cluster_genClients[cluster.Name]
							dep_err := cluster_client.Get(context.TODO(), dep, hasInstance.Namespace, hasInstance.Spec.ScalingOptions.CpaTemplate.ScaleTargetRef.Name)
							if dep_err == nil {
								totalReplicas = totalReplicas + *dep.Spec.Replicas
								cpaCluster = append(cpaCluster, cluster.Name)
							}
						}
						deployInfo := &protobuf.CPADeployInfo{
							Name:        openmcpDep.Name,
							Namespace:   openmcpDep.Namespace,
							ReplicasNum: totalReplicas,
							CPAName:     hasInstance.Name,
							Clusters:    cpaCluster,
							//containers total Request 계산해서 put
							CpuRequest: totalCpuRequest,
							MemRequest: totalMemRequest,
						}

						cpaValue := analyticResource.CPAValue{}

						cpaValue.OmcpdeployInfo = deployInfo
						//cpaValue.InitReplicas = openmcpDep.Spec.Replicas
						cpaValue.ReplicasAfterScaling = totalReplicas
						cpaValue.CpaMin = hasInstance.Spec.ScalingOptions.CpaTemplate.MinReplicas
						cpaValue.CpaMax = hasInstance.Spec.ScalingOptions.CpaTemplate.MaxReplicas

						analyticResource.CPAInfoList[cpaKey] = cpaValue

						fmt.Println("Success to Put CPA ", cpaKey)

					} else {
						fmt.Println(openmcpDep_err)
					}
				}
			}

			return reconcile.Result{}, nil

		} else if hasInstance.Spec.MainController == "MemberCluster" {

			//Get Target OpenMCPDeployment
			openmcpDep := &resourcev1alpha1.OpenMCPDeployment{}
			openmcpDep_err := r.live.Get(context.TODO(), ObjectKey{Namespace: hasInstance.Namespace, Name: hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.ScaleTargetRef.Name}, openmcpDep)

			omcplog.V(3).Info(">>> Target OpenMCPDeployment [", openmcpDep.Name, " | ", openmcpDep.Namespace, "]")

			//totalHASTimeEnd2 := time.Since(totalHASTimeStart2)
			//omcplog.V(4).Info("------ Check HAS Time (2) : ", totalHASTimeEnd2)
			//totalHASTimeStart3 := time.Now()

			if openmcpDep_err == nil {

				var dep_list_for_hpa []string
				var cluster_dep_replicas map[string]int32
				var cluster_dep_request map[string]bool

				cluster_dep_replicas = make(map[string]int32)
				cluster_dep_request = make(map[string]bool)

				//Select Clusters to Deploy HPA, VPA (Policy + Check Deployment)
				for _, cluster := range clusterListItems {
					dep := &appsv1.Deployment{}
					cluster_client := cm.Cluster_genClients[cluster.Name]
					dep_err := cluster_client.Get(context.TODO(), dep, hasInstance.Namespace, hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.ScaleTargetRef.Name)
					if dep_err == nil {
						for i, _ := range dep.Spec.Template.Spec.Containers {
							if dep.Spec.Template.Spec.Containers[i].Resources.Requests == nil {
								cluster_dep_request[cluster.Name] = false
								break
							} else {
								cluster_dep_request[cluster.Name] = true
							}
						}
						cluster_dep_replicas[cluster.Name] = *dep.Spec.Replicas
						dep_list_for_hpa = append(dep_list_for_hpa, cluster.Name)
					}
				}
				//Delete HPA, VPA on except_cluster list
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
					for _, cluster := range dep_list_except {
						_ = r.DeleteHPAVPA(cm, cluster, request.Namespace, request.Name)

						err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
						if err_openmcp != nil {
							omcplog.V(0).Info("!! Fail to Update instance status", err_openmcp)
							return reconcile.Result{}, err_openmcp
						} else {
							omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (HPA Create)")
						}
					}

				}

				omcplog.V(3).Info(">>> Target Clusters ", dep_list_for_hpa, " except ", dep_list_except)

				//totalHASTimeEnd3 := time.Since(totalHASTimeStart3)
				//omcplog.V(4).Info("------ Check HAS Time (3) : ", totalHASTimeEnd3)
				//totalHASTimeStart4 := time.Now()

				if dep_list_for_hpa != nil {
					// Distribute min,max
					cluster_min_map, cluster_max_map, min_max_err := r.UpdateMinMaxDistributionPolicy(hasInstance, cluster_dep_request, dep_list_for_hpa, cluster_dep_replicas)
					if min_max_err != nil {
						omcplog.V(0).Info(min_max_err)
						return reconcile.Result{}, min_max_err
					}

					//totalHASTimeEnd4 := time.Since(totalHASTimeStart4)
					//omcplog.V(4).Info("------ Check HAS Time (4) : ", totalHASTimeEnd4)
					//totalHASTimeStart5 := time.Now()

					var sync_name string
					for _, clustername := range dep_list_for_hpa {

						// case 1) HPA
						if hasInstance.Spec.ScalingOptions.VpaTemplate == "Never" || (hasInstance.Spec.ScalingOptions.VpaTemplate == "Auto" && cluster_dep_request[clustername] == true) {
							foundHPA := &hpav2beta2.HorizontalPodAutoscaler{}
							cluster_client := cm.Cluster_genClients[clustername]
							err = cluster_client.Get(context.TODO(), foundHPA, hasInstance.Namespace, hasInstance.Name)
							fmt.Println("b")
							// CREATE HPA
							if err != nil && errors.IsNotFound(err) {

								hpa_min := cluster_min_map[clustername]
								hpa_max := cluster_max_map[clustername]

								//Create HPA Object
								hpa := r.UpdateHorizontalPodAutoscaler(request, hasInstance, hpa_min, hpa_max)
								//Create Sync resource (Sync Controller Watching)
								command := "create"
								sync_name, err = r.sendSyncHPA(hpa, command, clustername)

								hasInstance.Status.SyncRequestName = sync_name

								if err != nil {
									omcplog.V(0).Info("!![error] Fail to Create new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
									return reconcile.Result{}, err
								}
								omcplog.V(2).Info(">>> "+clustername+" Create HPA [ min:", *hpa.Spec.MinReplicas, " / max:", hpa.Spec.MaxReplicas, " ]")
								fmt.Println("d")
								//Status Update
								hasInstance.Status.LastSpec = hasInstance.Spec
								tmpMap[clustername] = 0
								hasInstance.Status.RebalancingCount = tmpMap
								hasInstance.Status.ChangeNeed = false
								hasInstance.Status.CheckSubResource = true
								hasInstance.Status.ClusterMinMaps = cluster_min_map
								hasInstance.Status.ClusterMaxMaps = cluster_max_map

								err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
								if err_openmcp != nil {
									omcplog.V(0).Info("!![error] Fail to update instance status", err_openmcp)
									return reconcile.Result{}, err_openmcp
								} else {
									omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (HPA Create)")
								}

								// HPA created successfully - return and requeue
							} else if err != nil {
								omcplog.V(0).Info("!![error] Fail to Get HPA")
								return reconcile.Result{}, err
							} else {
								if hasInstance.Status.LastSpec.ScalingOptions.HpaTemplate.Spec.MinReplicas != nil {
									if *hasInstance.Status.LastSpec.ScalingOptions.HpaTemplate.Spec.MinReplicas != *hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas || hasInstance.Status.LastSpec.ScalingOptions.HpaTemplate.Spec.MaxReplicas != hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas {

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
												omcplog.V(0).Info("!![error] Fail to Update HPA", "Hpa.Namespace", foundHPA.Namespace, "Hpa.Name", foundHPA.Name)
												return reconcile.Result{}, err
											}
										}
									}
								}
							}
							// case 2) HPA + VPA
						} else if hasInstance.Spec.ScalingOptions.VpaTemplate == "Always" || (hasInstance.Spec.ScalingOptions.VpaTemplate == "Auto" && cluster_dep_request[clustername] == false) {
							// First, Create HPA

							//Status Update
							foundHPA := &hpav2beta2.HorizontalPodAutoscaler{}
							cluster_client := cm.Cluster_genClients[clustername]
							err = cluster_client.Get(context.TODO(), foundHPA, hasInstance.Namespace, hasInstance.Name)
							if err != nil && errors.IsNotFound(err) { //CREATE HPA

								hpa_min := cluster_min_map[clustername]
								hpa_max := cluster_max_map[clustername]

								//Create HPA Object
								hpa := r.UpdateHorizontalPodAutoscaler(request, hasInstance, hpa_min, hpa_max)
								//Create Sync resource (Sync Controller Watching)
								command := "create"
								sync_name, err = r.sendSyncHPA(hpa, command, clustername)

								hasInstance.Status.SyncRequestName = sync_name

								if err != nil {
									omcplog.V(0).Info("!![error] Fail to Create new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
									return reconcile.Result{}, err
								}
								omcplog.V(2).Info(">>> "+clustername+" Create HPA [ min:", *hpa.Spec.MinReplicas, " / max:", hpa.Spec.MaxReplicas, " ]")

								//Status Update
								hasInstance.Status.LastSpec = hasInstance.Spec
								tmpMap[clustername] = 0
								hasInstance.Status.RebalancingCount = tmpMap
								hasInstance.Status.ChangeNeed = false
								hasInstance.Status.CheckSubResource = true
								hasInstance.Status.ClusterMinMaps = cluster_min_map
								hasInstance.Status.ClusterMaxMaps = cluster_max_map

								err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
								if err_openmcp != nil {
									omcplog.V(0).Info("!![error] Fail to Update instance status", err_openmcp)
									return reconcile.Result{}, err_openmcp
								} else {
									omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (HPA Create)")
								}
							} else if err != nil {
								omcplog.V(0).Info("!![error] Fail to Get HPA")
								return reconcile.Result{}, err
							} else {

								if hasInstance.Status.LastSpec.ScalingOptions.HpaTemplate.Spec.MinReplicas != nil {
									if *hasInstance.Status.LastSpec.ScalingOptions.HpaTemplate.Spec.MinReplicas != *hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas || hasInstance.Status.LastSpec.ScalingOptions.HpaTemplate.Spec.MaxReplicas != hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas {

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
												omcplog.V(0).Info("!![error] Fail to Update HPA", "Hpa.Namespace", foundHPA.Namespace, "Hpa.Name", foundHPA.Name)
												return reconcile.Result{}, err
											}
										}
									}
								}
							}
							//And then, Create VPA
							foundVPA := &vpav1beta2.VerticalPodAutoscaler{}
							err = r.ghosts[clustername].Get(context.TODO(), ObjectKey{Namespace: hasInstance.Namespace, Name: hasInstance.Name}, foundVPA)

							if err != nil && errors.IsNotFound(err) { //CREATE VPA

								vpa := r.UpdateVerticalPodAutoscaler(request, hasInstance)

								command := "create"
								_, err = r.sendSyncVPA(vpa, command, clustername)

								hasInstance.Status.SyncRequestName = sync_name

								if err != nil {
									omcplog.V(0).Info("!![error] Fail to Create new VPA", "VPA.Namespace", vpa.Namespace, "VPA.Name", vpa.Name)
									return reconcile.Result{}, err
								}
								// VPA created successfully - return and requeue
								omcplog.V(2).Info(">>> " + clustername + " Create VPA")

								hasInstance.Status.LastSpec = hasInstance.Spec

								err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
								if err_openmcp != nil {
									omcplog.V(0).Info("!![error] Fail to Update instance status", err_openmcp)
									return reconcile.Result{}, err_openmcp
								} else {
									omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (VPA Create)")
								}

							} else if err != nil {
								omcplog.V(0).Info("!![error] Fail to Get VPA")
								return reconcile.Result{}, err
							} else if err == nil { //UPDATE VPA
							}
						}
					}

					//totalHASTimeEnd5 := time.Since(totalHASTimeStart5)
					//omcplog.V(4).Info("------ Check HAS Time (5) : ", totalHASTimeEnd5)

					//Check for Changes to OpenMCPHPA resource
					hasInstance.Status.LastSpec = hasInstance.Spec
					hasInstance.Status.ChangeNeed = false

					err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
					if err_openmcp != nil {
						omcplog.V(0).Info("!![error] Fail to Update instance status", err_openmcp)
						return reconcile.Result{}, err_openmcp
					} else {
						omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (End)")
					}
				}

			} else if openmcpDep_err != nil && errors.IsNotFound(openmcpDep_err) {
				omcplog.V(0).Info("!![error] Can't find OpenmcpDeployment")
				for _, cluster := range clusterListItems {
					_ = r.DeleteHPAVPA(cm, cluster.Name, request.Namespace, request.Name)

					err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
					if err_openmcp != nil {
						omcplog.V(0).Info("!! Fail to Update instance status", err_openmcp)
						return reconcile.Result{}, err_openmcp
					} else {
						omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (HPA Create)")
					}
				}
				return reconcile.Result{}, openmcpDep_err
			} else {
				omcplog.V(0).Info("!![error] Fail to Get OpenMCPDeployment - ", openmcpDep_err)

				return reconcile.Result{}, openmcpDep_err
			}
		}

		totalHASTimeEnd0 := time.Since(totalHASTimeStart1)

		omcplog.V(4).Info("------ ToTal HAS Time : ", totalHASTimeEnd0, " ------")

		return reconcile.Result{}, nil
	}

	rebalancingTimeStart := time.Now()

	//Get Target  OpenMCPDeployment
	openmcpDep := &resourcev1alpha1.OpenMCPDeployment{}
	openmcpDep_err := r.live.Get(context.TODO(), ObjectKey{Namespace: hasInstance.Namespace, Name: hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.ScaleTargetRef.Name}, openmcpDep)

	var dep_list_for_hpa []string
	//var cluster_dep_replicas map[string]int32
	var cluster_dep_request map[string]bool

	if openmcpDep_err == nil {
		//cluster_dep_replicas = make(map[string]int32)
		cluster_dep_request = make(map[string]bool)

		for _, cluster := range cm.Cluster_list.Items {
			dep := &appsv1.Deployment{}
			cluster_client := cm.Cluster_genClients[cluster.Name]
			dep_err := cluster_client.Get(context.TODO(), dep, hasInstance.Namespace, hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.ScaleTargetRef.Name)
			if dep_err == nil {

				for i, _ := range dep.Spec.Template.Spec.Containers {
					if dep.Spec.Template.Spec.Containers[i].Resources.Requests == nil {
						cluster_dep_request[cluster.Name] = false
						break
					} else {
						cluster_dep_request[cluster.Name] = true
					}
				}

				//cluster_dep_replicas[cluster.Name] = *dep.Spec.Replicas
				dep_list_for_hpa = append(dep_list_for_hpa, cluster.Name)
			}
		}
	}

	for _, clustername := range dep_list_for_hpa {

		foundHPA := &hpav2beta2.HorizontalPodAutoscaler{}
		cluster_client := cm.Cluster_genClients[clustername]
		err = cluster_client.Get(context.TODO(), foundHPA, hasInstance.Namespace, hasInstance.Name)

		if foundHPA.Spec.MinReplicas != nil {

			//UPDATE HPA (Rebalancing or Update min/max value)
			if foundHPA.Spec.MaxReplicas == foundHPA.Status.CurrentReplicas && foundHPA.Status.CurrentReplicas == foundHPA.Status.DesiredReplicas {
				if lastTimeRebalancing.IsZero() || (!lastTimeRebalancing.IsZero() && time.Since(lastTimeRebalancing) > time.Second*100) {

					hasInstance, check := r.MaxRebalancing(cm, hasInstance, dep_list_for_hpa, clustername, foundHPA)

					if check == "Success" {
						err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
						if err_openmcp != nil {
							omcplog.V(0).Info("!![error] Fail to update instance status \"RebalancingCount\"", err_openmcp)
							return reconcile.Result{}, err_openmcp
						} else {
							omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (RebalancingCount)")
						}

						rebalancingTimeEnd := time.Since(rebalancingTimeStart)
						omcplog.V(4).Info("------ Check Rebalancing Time : ", rebalancingTimeEnd, " ------")
					}

				}
			} else if *foundHPA.Spec.MinReplicas > 1 && *foundHPA.Spec.MinReplicas == foundHPA.Status.CurrentReplicas && foundHPA.Status.CurrentReplicas == foundHPA.Status.DesiredReplicas {
				if lastTimeRebalancing.IsZero() || (!lastTimeRebalancing.IsZero() && time.Since(lastTimeRebalancing) > time.Second*100) {

					hasInstance, check := r.MinRebalancing(cm, hasInstance, dep_list_for_hpa, clustername, foundHPA)

					if check == "Success" {

						err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
						if err_openmcp != nil {
							omcplog.V(0).Info("!![error] Fail to Update instance status \"RebalancingCount\"", err_openmcp)
							return reconcile.Result{}, err_openmcp
						} else {
							omcplog.V(3).Info(">>> OpenMCPHPA LastSpec Update (RebalancingCount)")
						}

						rebalancingTimeEnd := time.Since(rebalancingTimeStart)
						omcplog.V(4).Info("------ Check Rebalancing Time : ", rebalancingTimeEnd, " ------")
					}
				}
			}

		}
	}
	//멤버 클러스터의 hpa, vpa가 삭제되었을 때 다시 재생성
	if hasInstance.Status.CheckSubResource == true {
		omcplog.V(2).Info("[Member Cluster Check Service]")
		sync_req_name := ""
		for k, v := range hasInstance.Status.ClusterMinMaps {
			cluster_name := k
			min_replicas := v
			max_replicas := hasInstance.Status.ClusterMaxMaps[cluster_name]

			if v == 0 {
				continue
			}

			found := &hpav2beta2.HorizontalPodAutoscaler{}
			cluster_client := cm.Cluster_genClients[cluster_name]
			err = cluster_client.Get(context.TODO(), found, hasInstance.Namespace, hasInstance.Name)

			if err != nil && errors.IsNotFound(err) {
				omcplog.V(2).Info("HPA resource ReDeployed in [", cluster_name, "]")

				//Re-Create HPA Object
				hpa := r.UpdateHorizontalPodAutoscaler(request, hasInstance, min_replicas, max_replicas)
				//Create Sync resource (Sync Controller Watching)
				command := "create"
				sync_name, err := r.sendSyncHPA(hpa, command, cluster_name)

				hasInstance.Status.SyncRequestName = sync_name

				if err != nil {
					omcplog.V(0).Info("!![error] Fail to ReDeploy new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
					return reconcile.Result{}, err
				}
				omcplog.V(2).Info(">>> "+cluster_name+" ReDeploy HPA [ min:", *hpa.Spec.MinReplicas, " / max:", hpa.Spec.MaxReplicas, " ]")

			}

		}
		omcplog.V(3).Info("sync_req_name : ", sync_req_name)
	}

	return reconcile.Result{}, nil

}

func (r *reconciler) UpdateMinMaxDistributionPolicy(hasInstance *resourcev1alpha1.OpenMCPHybridAutoScaler, cluster_dep_request map[string]bool, dep_list_for_hpa []string, cluster_dep_replicas map[string]int32) (map[string]int32, map[string]int32, error) {
	type ObjectKey = types.NamespacedName

	foundPolicy := &policyv1alpha1.OpenMCPPolicy{}
	minmax_policy_err := r.live.Get(context.TODO(), ObjectKey{Namespace: "openmcp", Name: "hpa-minmax-distribution-mode"}, foundPolicy)

	cluster_min_map, cluster_max_map, hasInstance, min_max_err := r.CreateMinMaxMap(hasInstance, cluster_dep_request, foundPolicy, minmax_policy_err, dep_list_for_hpa, cluster_dep_replicas)

	return cluster_min_map, cluster_max_map, min_max_err
}

func (r *reconciler) UpdateTargetClusterPolicy(cm *clusterManager.ClusterManager, hasInstance *resourcev1alpha1.OpenMCPHybridAutoScaler) ([]fedv1b1.KubeFedCluster, *resourcev1alpha1.OpenMCPHybridAutoScaler) {
	checkPolicy := 0
	clusterListItems := make([]fedv1b1.KubeFedCluster, 0)

	type ObjectKey = types.NamespacedName

	if hasInstance.Status.Policies != nil { // Check policy
		for _, tmp := range hasInstance.Status.Policies {
			if tmp.Type == "Target" { // if policy exists
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
	if checkPolicy == 0 { // if policy doesn't exist
		foundPolicy := &policyv1alpha1.OpenMCPPolicy{}
		target_cluster_policy_err := r.live.Get(context.TODO(), ObjectKey{Namespace: "openmcp", Name: "has-target-cluster"}, foundPolicy)
		if target_cluster_policy_err == nil {
			if foundPolicy.Spec.PolicyStatus == "Enabled" {
				omcplog.V(3).Info(">>> Apply \"Cluster Target\" Policy (Enabled)")
				hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				for _, cluster := range cm.Cluster_list.Items {
					for _, value := range foundPolicy.Spec.Template.Spec.Policies[0].Value {
						if cluster.Name == value {
							clusterListItems = append(clusterListItems, cluster)
						}
					}
				}
			} else {
				omcplog.V(3).Info(">>> Apply \"Cluster Target\" Policy (Disabled - set default)")
				omp := make([]policyv1alpha1.OpenMCPPolicies, 1)
				omp[0].Type = "Target"
				omp_value := make([]string, 0)
				omp_value = append(omp_value, "Default")
				omp[0].Value = omp_value
				hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
				clusterListItems = cm.Cluster_list.Items
			}
		} else {
			omcplog.V(1).Info("!![error] Fail to get Policy \"Cluster Target\" (set default)")
			omp := make([]policyv1alpha1.OpenMCPPolicies, 1)
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

func (r *reconciler) CreateMinMaxMap(hasInstance *resourcev1alpha1.OpenMCPHybridAutoScaler, cluster_dep_request map[string]bool, foundPolicy *policyv1alpha1.OpenMCPPolicy, minmax_policy_err error, dep_list_for_hpa []string, cluster_dep_replicas map[string]int32) (map[string]int32, map[string]int32, *resourcev1alpha1.OpenMCPHybridAutoScaler, error) {
	timeStart_mixmaxdist := time.Now()

	cluster_min_map := make(map[string]int32)
	cluster_max_map := make(map[string]int32)
	var err error

	type ObjectKey = types.NamespacedName

	checkPolicy := 0
	if hasInstance.Status.Policies != nil { // Check minmaxdist policy
		for n, tmp := range hasInstance.Status.Policies {
			if tmp.Type == "Mode" { // if policy exists
				if hasInstance.Status.Policies[n].Value[0] == "Equal" {
					cluster_min_map = HpaMinScheduling_equal(dep_list_for_hpa, *hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling_equal(dep_list_for_hpa, hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas)
				} else {
					cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas)

					for _, cluster := range dep_list_for_hpa {
						if hasInstance.Spec.ScalingOptions.VpaTemplate == "Always" || (hasInstance.Spec.ScalingOptions.VpaTemplate == "Auto" && cluster_dep_request[cluster] == false) {
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
	if checkPolicy == 0 { // if policy doesn't exist

		if minmax_policy_err == nil {
			if foundPolicy.Spec.PolicyStatus == "Enabled" {
				omcplog.V(3).Info(">>> Apply \"Min Max Distribution\" Policy (Enabled)")
				if foundPolicy.Spec.Template.Spec.Policies[0].Value[0] == "Equal" {
					cluster_min_map = HpaMinScheduling_equal(dep_list_for_hpa, *hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling_equal(dep_list_for_hpa, hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas)
					hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				} else if foundPolicy.Spec.Template.Spec.Policies[0].Value[0] == "Unequal" {
					cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas)
					hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				} else {
					cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas)
					foundPolicy.Spec.Template.Spec.Policies[0].Value[0] = "Default"
					hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				}

			} else {
				omcplog.V(3).Info(">>> Apply \"Min Max Distribution\" Policy (Disabled - set default)")
				cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas)
				cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas)
				omp := make([]policyv1alpha1.OpenMCPPolicies, 1)
				omp_value := make([]string, 1)
				omp_value[0] = "Default"
				omp[0].Type = "Mode"
				omp[0].Value = omp_value
				hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
			}
		} else {
			omcplog.V(1).Info("!![error] Fail to get Policy \"Min Max Distribution\" (set default)")
			cluster_min_map = HpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas)
			cluster_max_map = HpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas)
			omp := make([]policyv1alpha1.OpenMCPPolicies, 1)
			omp_value := make([]string, 1)
			omp_value[0] = "Default"
			omp[0].Type = "Mode"
			omp[0].Value = omp_value
			hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
		}

		for _, cluster := range dep_list_for_hpa {
			if hasInstance.Spec.ScalingOptions.VpaTemplate == "Always" || (hasInstance.Spec.ScalingOptions.VpaTemplate == "Auto" && cluster_dep_request[cluster] == false) {
				if cluster_min_map[cluster] < 2 {
					cluster_min_map[cluster] = 2
				}
			}
		}

		omcplog.V(4).Info("******* [Start] HAS Min/Max Distribution *******")
		timeEnd_mixmaxdist := time.Since(timeStart_mixmaxdist)
		for _, cluster := range dep_list_for_hpa {
			omcplog.V(4).Info("[", cluster, "] min:", cluster_min_map[cluster], " / max:", cluster_max_map[cluster], "")
		}
		omcplog.V(4).Info("=> Total Min/Max Distribution Time [", timeEnd_mixmaxdist, "]")
		omcplog.V(4).Info("*******  [End] HAS Min/Max Distribution  *******")

		//Update OpenMCPHAS Policy Option
		err = r.live.Status().Update(context.TODO(), hasInstance)
		if err != nil {
			omcplog.V(0).Info("!![error] Fail to Update Policy Status")
		} else {
			//omcplog.V(3).Info(">>> Success to Update Policy Status")
		}

	}

	return cluster_min_map, cluster_max_map, hasInstance, err
}

func (r *reconciler) MinRebalancing(cm *clusterManager.ClusterManager, hasInstance *resourcev1alpha1.OpenMCPHybridAutoScaler, dep_list_for_hpa []string, clustername string, foundHPA *hpav2beta2.HorizontalPodAutoscaler) (*resourcev1alpha1.OpenMCPHybridAutoScaler, string) {

	var dep_list_for_analysis []string

	minReplicasMap := map[string]int32{}
	currentReplicasMap := map[string]int32{}

	for _, cn := range dep_list_for_hpa {
		analysisHPA := &hpav2beta2.HorizontalPodAutoscaler{}
		err := cm.Cluster_genClients[cn].Get(context.TODO(), analysisHPA, hasInstance.Namespace, hasInstance.Name)
		if err == nil {
			if *analysisHPA.Spec.MinReplicas < analysisHPA.Status.CurrentReplicas-1 {
				dep_list_for_analysis = append(dep_list_for_analysis, cn)
				minReplicasMap[cn] = *analysisHPA.Spec.MinReplicas
				currentReplicasMap[cn] = analysisHPA.Status.CurrentReplicas
			}
		} else {
			omcplog.V(0).Info("!![error] Fail to Get hpa info : ", err)
		}
	}

	//when there is no candidate cluster
	if len(dep_list_for_analysis) == 0 {
		omcplog.V(4).Info("!![error] Fail Rebalancing : There is no candidate cluster")

		return hasInstance, "Fail"
	} else {
		//Result from Analytic Engine(gRPC)
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
			omcplog.V(0).Info("!![error] Fail Min Rebalacing : Failed to get Analysis Result :(")
		} else {
			qosCluster := result.TargetCluster
			omcplog.V(2).Info(">>> " + clustername + " min rebalancing")

			omcplog.V(3).Info("     => Anlysis Result [", qosCluster, "]")

			omcplog.V(4).Info("     => Total Rebalancing Time [", timeEnd_mixmaxrebal, "] (Analysis + gRPC)")

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

			if current_err != nil || qos_err != nil { //if errors
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

				//Update Status
				hasInstance.Status.RebalancingCount[clustername] += 1

				lastTimeRebalancing = time.Now()
				omcplog.V(4).Info("     => RebalancingTime : ", lastTimeRebalancing)
			}
		}
	}

	return hasInstance, "Success"
}

func (r *reconciler) MaxRebalancing(cm *clusterManager.ClusterManager, hasInstance *resourcev1alpha1.OpenMCPHybridAutoScaler, dep_list_for_hpa []string, clustername string, foundHPA *hpav2beta2.HorizontalPodAutoscaler) (*resourcev1alpha1.OpenMCPHybridAutoScaler, string) {

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
			omcplog.V(0).Info("!![error] Fail to Get hpa info : ", err)
		}
	}

	//when there is no candidate cluster
	if len(dep_list_for_analysis) == 0 {
		omcplog.V(4).Info("!![error] Fail Max Rebalancing : There is no candidate cluster")
		return hasInstance, "Fail"
	} else {

		//Result from Analytic Engine (gRPC)
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
			omcplog.V(0).Info("!![error] Fail Rebalacing : Failed to get Analysis Result :(")
		} else {
			qosCluster := result.TargetCluster
			omcplog.V(2).Info(">>> " + clustername + " max rebalancing")
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

			if current_err != nil || qos_err != nil { //if errors
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

				//Update Status
				hasInstance.Status.RebalancingCount[clustername] += 1

				lastTimeRebalancing = time.Now()
				omcplog.V(4).Info("     => RebalancingTime : ", lastTimeRebalancing)

			}
		}
	}

	return hasInstance, "Success"

}

func (r *reconciler) UpdateHorizontalPodAutoscaler(req reconcile.Request, m *resourcev1alpha1.OpenMCPHybridAutoScaler, min int32, max int32) *hpav2beta2.HorizontalPodAutoscaler {
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

		Spec:   m.Spec.ScalingOptions.HpaTemplate.Spec,
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

func (r *reconciler) UpdateVerticalPodAutoscaler(req reconcile.Request, m *resourcev1alpha1.OpenMCPHybridAutoScaler) *vpav1beta2.VerticalPodAutoscaler {
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
				Name:       m.Spec.ScalingOptions.HpaTemplate.Spec.ScaleTargetRef.Name,
			},
			UpdatePolicy: &vpav1beta2.PodUpdatePolicy{
				UpdateMode: &vpaUpdateMode,
			},
		},
		Status: vpav1beta2.VerticalPodAutoscalerStatus{},
	}
	reference.SetMulticlusterControllerReference(vpa, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))
	return vpa
}
func (r *reconciler) DeleteHPAVPA(cm *clusterManager.ClusterManager, cluster string, namespace string, name string) string {
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
	} else if err1 != nil {
		omcplog.V(0).Info("!![error] Fail to Delete HPA - ", err1)
	} else if err1 == nil {
		var sync_err1 error

		hpa.TypeMeta.Kind = "HorizontalPodAutoscaler"
		hpa.TypeMeta.APIVersion = "autoscaling/v2beta2"

		command := "delete"
		sync_name, sync_err1 = r.sendSyncHPA(hpa, command, cluster)

		if sync_err1 != nil && errors.IsNotFound(sync_err1) {
		} else if sync_err1 != nil {
			omcplog.V(0).Info("!![error] Fail to Delete VPA - ", sync_err1)
		} else if sync_err1 == nil {
			omcplog.V(2).Info(">>> " + cluster + " Delete HPA")
		}
	}

	err2 := r.ghosts[cluster].Get(context.TODO(), ObjectKey{Namespace: namespace, Name: name}, vpa)
	if err2 != nil && errors.IsNotFound(err2) {
	} else if err2 != nil {
		omcplog.V(0).Info("!![error] Fail to Get VPA - ", err2)
	} else if err2 == nil {
		var sync_err2 error

		command := "delete"
		sync_name, sync_err2 = r.sendSyncVPA(vpa, command, cluster)

		if sync_err2 != nil && errors.IsNotFound(sync_err2) {
		} else if sync_err2 != nil {
			omcplog.V(0).Info("!![error] Fail to Delete VPA - ", sync_err2)
		} else if sync_err2 == nil {
			omcplog.V(2).Info(">>> " + cluster + " Delete VPA")
		}
	}

	return sync_name

}
func (r *reconciler) DeleteAllHPAVPA(cm *clusterManager.ClusterManager, namespace string, name string) {
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

	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_genClients[cluster.Name]

		err1 := cluster_client.Get(context.TODO(), hpa, namespace, name)

		if err1 != nil && errors.IsNotFound(err1) {
		} else if err1 != nil {
			omcplog.V(0).Info("!![error] Fail to Delete HPA - ", err1)
		} else if err1 == nil {
			var sync_err1 error

			hpa.TypeMeta.Kind = "HorizontalPodAutoscaler"
			hpa.TypeMeta.APIVersion = "autoscaling/v2beta2"

			command := "delete"
			_, sync_err1 = r.sendSyncHPA(hpa, command, cluster.Name)

			if sync_err1 != nil && errors.IsNotFound(sync_err1) {
			} else if sync_err1 != nil {
				omcplog.V(0).Info("!![error] Fail to Delete HPA - ", sync_err1)
			} else if sync_err1 == nil {
				omcplog.V(2).Info(">>> " + cluster.Name + " Delete HPA")
			}
		}

		err2 := r.ghosts[cluster.Name].Get(context.TODO(), ObjectKey{Namespace: namespace, Name: name}, vpa)
		if err2 != nil && errors.IsNotFound(err2) {
		} else if err2 != nil {
			omcplog.V(0).Info("!![error] Fail to Get VPA - ", err2)
		} else if err2 == nil {
			var sync_err2 error

			command := "delete"
			_, sync_err2 = r.sendSyncVPA(vpa, command, cluster.Name)

			if sync_err2 != nil && errors.IsNotFound(sync_err2) {
			} else if sync_err2 != nil {
				omcplog.V(0).Info("!![error] Fail to Delete VPA - ", sync_err2)
			} else if sync_err2 == nil {
				omcplog.V(2).Info(">>> " + cluster.Name + " Delete VPA")
			}
		}
	}

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
	return cluster_maxreplicas_map
}

func HpaMinScheduling_equal(dep_list_for_hpa []string, min int32) map[string]int32 {

	cluster_minreplicas_map := make(map[string]int32)

	for _, clustername := range dep_list_for_hpa {
		cluster_minreplicas_map[clustername] = 0
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
	return cluster_maxreplicas_map
}

func LabelsForHpa(name string) map[string]string {
	return map[string]string{"app": "openmcphybridautoscaler", "openmcphybridautoscaler_cr": name}
}
