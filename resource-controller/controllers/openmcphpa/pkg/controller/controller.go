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
	"os"

	//"sort"
	//"math/rand"
	//"time"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	//corev1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/types"
	//"admiralty.io/multicluster-controller/pkg/reference"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"resource-controller/apis"
	ketiv1alpha1 "resource-controller/apis/keti/v1alpha1"

	vpav1beta2 "github.com/kubernetes/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	hpav2beta2 "k8s.io/api/autoscaling/v2beta2"
	//vpav1beta2 "k8s.io/api/autoscaling/v1beta2"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"

	"resource-controller/controllers/openmcphpa/pkg/protobuf"
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
		//ghostclients = append(ghostclients, ghostclient)
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

type reconciler struct {
	live           client.Client
	ghosts         map[string]client.Client
	ghostNamespace string
}

var i int = 0

func (r *reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	//fmt.Println("Step 7.	r.Reconcile()")
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	i += 1
	fmt.Println("********* [", i, "] *********")
	fmt.Println("Request Namespace: ", request.Namespace, " /  Request Name: ", request.Name, " / Request Context: ", request.Context)
	cm := NewClusterManager()
	type ObjectKey = types.NamespacedName

	// Fetch the OpenMCPHybridAutoScaler instance
	hasInstance := &ketiv1alpha1.OpenMCPHybridAutoScaler{}
	err := r.live.Get(context.TODO(), request.NamespacedName, hasInstance)
	//OpenMCPHPA Delete
	if err != nil {
		if errors.IsNotFound(err) {

			r.DeleteOpenMCPHPA(cm, request.Namespace, request.Name+"-hpa", request.Name+"-vpa")

			reqLogger.Info("openmcphybridautoscaler resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get hasInstance")
		return reconcile.Result{}, err
	}

	//----------------------------타겟 클러스터 선정 정책------------------------------------------------------------------------------------------------------------
	checkPolicy := 0
	clusterListItems := make([]fedv1b1.KubeFedCluster, 0)

	if hasInstance.Status.Policies != nil { // 기존에 설정된 정책이 있는지 확인
		for _, tmp := range hasInstance.Status.Policies {
			if tmp.Type == "Target" { // Target 이미 설정되어 있으면
				fmt.Println(">>> Policy \"Cluster Target\" Existed")
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
		foundPolicy := &ketiv1alpha1.OpenMCPPolicyEngine{}
		target_cluster_policy_err := r.live.Get(context.TODO(), ObjectKey{Namespace: "openmcp", Name: "hpa-target-cluster"}, foundPolicy)
		if target_cluster_policy_err == nil {
			if foundPolicy.Spec.PolicyStatus == "Enabled" {
				fmt.Println(">>> Policy \"Cluster Target\" Apply (Enabled)")
				hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				for _, cluster := range cm.Cluster_list.Items {
					for _, value := range foundPolicy.Spec.Template.Spec.Policies[0].Value {
						if cluster.Name == value {
							clusterListItems = append(clusterListItems, cluster)
						}
					}
				}
			} else {
				fmt.Println(">>> Policy \"Cluster Target\" Apply (Disabled - set default)")
				omp := make([]ketiv1alpha1.OpenMCPPolicies, 1)
				omp[0].Type = "Target"
				omp_value := make([]string, 0)
				omp_value = append(omp_value, "Default")
				omp[0].Value = omp_value
				hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
				clusterListItems = cm.Cluster_list.Items
			}
		} else {
			fmt.Println("policy_err : ", target_cluster_policy_err)
			fmt.Println("Fail to get policy \"Cluster Target\" (set default)")
			omp := make([]ketiv1alpha1.OpenMCPPolicies, 1)
			omp[0].Type = "Target"
			omp_value := make([]string, 0)
			omp_value = append(omp_value, "Default")
			omp[0].Value = omp_value
			hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
			clusterListItems = cm.Cluster_list.Items
		}

		err := r.live.Status().Update(context.TODO(), hasInstance)
		if err != nil {
			fmt.Println("Policy Status Update Error")
			return reconcile.Result{}, err
		} else {
			fmt.Println(">>> Policy Status UPDATE Success")
		}
	}
	//-------------------------------------------------------------------------------------------------------------------------------------------------------

	openmcpDep := &ketiv1alpha1.OpenMCPDeployment{}
	openmcpDep_err := r.live.Get(context.TODO(), ObjectKey{Namespace: hasInstance.Namespace, Name: hasInstance.Spec.HpaTemplate.Spec.ScaleTargetRef.Name}, openmcpDep)

	fmt.Println(">>> Target OpenMCPDeployment [", openmcpDep.Name, ":", openmcpDep.Namespace, "]")

	if openmcpDep_err == nil {

		var dep_list_for_hpa []string
		var cluster_dep_replicas map[string]int32
		var cluster_dep_request map[string]bool

		cluster_dep_replicas = make(map[string]int32)
		cluster_dep_request = make(map[string]bool)

		//*** cluster list policy
		//해당 cluster에만 hpa 생성
		for _, cluster := range clusterListItems {
			dep := &appsv1.Deployment{}
			cluster_client := cm.Cluster_clients[cluster.Name]
			dep_err := cluster_client.Get(context.TODO(), dep, hasInstance.Namespace, hasInstance.Spec.HpaTemplate.Spec.ScaleTargetRef.Name+"-deploy")
			if dep_err == nil {
				//fmt.Println("deploy request: ", dep.Spec.Template.Spec.Containers[0].Resources.Requests) // 컨테이너 여러개일 경우?
				if dep.Spec.Template.Spec.Containers[0].Resources.Requests == nil {
					//fmt.Println("no request!")
					cluster_dep_request[cluster.Name] = false
				} else {
					//fmt.Println("request exists!")
					cluster_dep_request[cluster.Name] = true
				}
				cluster_dep_replicas[cluster.Name] = *dep.Spec.Replicas
				dep_list_for_hpa = append(dep_list_for_hpa, cluster.Name)
			}
		}

		//-----------타켓 클러스터 외의 클러스터에 배포된 hpa, vpa 삭제-------------------------------------------
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
		//fmt.Println("dep_list_except:", dep_list_except)
		if dep_list_except != nil {
			r.DeleteHPAVPA(cm, dep_list_except, request.Namespace, request.Name+"-hpa", request.Name+"-vpa")
		}

		fmt.Println(">>> Target Clusters ", dep_list_for_hpa, " /// except ", dep_list_except)
		//----------------------------------------------------------------------------------------------

		if dep_list_for_hpa != nil {
			//min/max 분배 정책
			cluster_min_map, cluster_max_map, hasInstance, min_max_err := CreateMinMaxMap(r, cm, hasInstance, dep_list_for_hpa, cluster_dep_replicas)

			if min_max_err != nil {
				fmt.Println(min_max_err)
				return reconcile.Result{}, min_max_err
			}

			for _, clustername := range dep_list_for_hpa {
				// VPA 생성하지 않는 경우
				if hasInstance.Spec.VpaMode == "Never" || (hasInstance.Spec.VpaMode == "Auto" && cluster_dep_request[clustername] == true) {
					//fmt.Println("[ Only HPA ]")
					// Check if the HPA already exists, if not create a new one
					foundHPA := &hpav2beta2.HorizontalPodAutoscaler{}
					cluster_client := cm.Cluster_clients[clustername]
					err = cluster_client.Get(context.TODO(), foundHPA, hasInstance.Namespace, hasInstance.Name+"-hpa")

					if err != nil && errors.IsNotFound(err) { //CREATE HPA
						// Define a new HPA
						hpa_min := cluster_min_map[clustername]
						hpa_max := cluster_max_map[clustername]
						//fmt.Println(clustername, " - min : ", hpa_min, ", max : ", hpa_max)
						hpa := r.updateHorizontalPodAutoscaler(request, hasInstance, hpa_min, hpa_max)

						reqLogger.Info("Creating a new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
						err = cluster_client.Create(context.TODO(), hpa)
						if err != nil {
							reqLogger.Error(err, "Failed to create new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
							return reconcile.Result{}, err
						}
						fmt.Println(">>> "+clustername+" Create HPA [ min:", *hpa.Spec.MinReplicas, "/ max:", hpa.Spec.MaxReplicas, "]")
						// HPA created successfully - return and requeue

						//Status Update
						hasInstance.Status.LastSpec = hasInstance.Spec

						err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
						if err_openmcp != nil {
							fmt.Println("Failed to update instance status", err)
							return reconcile.Result{}, err
						} else {
							fmt.Println(">>> OpenMCPHPA LastSpec Update (HPA Create)")
						}

						return reconcile.Result{Requeue: true}, nil
					} else if err != nil {
						reqLogger.Error(err, "Failed to get HPA")
						return reconcile.Result{}, err
					} else if err == nil { //UPDATE HPA - Rebalancing 조건을 만족했을 때 또는 min/max값이 수정되었을때
						if foundHPA.Spec.MaxReplicas == foundHPA.Status.CurrentReplicas && foundHPA.Status.CurrentReplicas == foundHPA.Status.DesiredReplicas { //||  (*foundHPA.Spec.MinReplicas == foundHPA.Status.CurrentReplicas && foundHPA.Status.CurrentReplicas == foundHPA.Status.DesiredReplicas) {
							fmt.Println(">>> " + clustername + " min/max rebalancing")
							//분석 엔진을 통해 얻은 결과
							//---------------------------------------------------------------------------------------------------
							SERVER_IP := os.Getenv("GRPC_SERVER")
							SERVER_PORT := os.Getenv("GRPC_PORT")

							grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

							hi := &protobuf.OpenMCPHPAInfo{DeploymentName: "test-pod", DeploymentNamespace: "test", CL: nil}
							hi.CL = make([]*protobuf.ClusterList, 0)
							hi.CL = append(hi.CL, &protobuf.ClusterList{ClusterName: "cluster2"})
							hi.CL = append(hi.CL, &protobuf.ClusterList{ClusterName: "cluster3"})

							rv := &protobuf.RequestValue{}
							rv.HpaInfo = make([]*protobuf.OpenMCPHPAInfo, 0)
							rv.HpaInfo = append(rv.HpaInfo, hi)

							r, gRPCerr := grpcClient.SendAnalysisResult(context.TODO(), rv)
							if gRPCerr != nil {
								fmt.Printf("could not connect : %v", gRPCerr)
							}
							fmt.Println(">>> Anlysis Result [", r, "]")
							//---------------------------------------------------------------------------------------------------
							qosCluster := r.TargetCluster

							updateHPA := foundHPA
							updateHPA.Spec.MaxReplicas += 1

							qos_cluster_client := cm.Cluster_clients[qosCluster]
							foundQosHPA := &hpav2beta2.HorizontalPodAutoscaler{}
							err := qos_cluster_client.Get(context.TODO(), foundQosHPA, hasInstance.Namespace, hasInstance.Name+"-hpa")

							if err != nil {
								fmt.Println("err: ", err)
							}

							updateQosHPA := foundQosHPA
							updateQosHPA.Spec.MaxReplicas -= 1

							current_err := cluster_client.Update(context.TODO(), updateHPA)
							fmt.Println(">>> "+clustername+" Update HPA [ min:", *updateHPA.Spec.MinReplicas, "/ max:", updateHPA.Spec.MaxReplicas, "]")
							qos_err := qos_cluster_client.Update(context.TODO(), updateQosHPA)
							fmt.Println(">>> "+qosCluster+" Update HPA [ min:", *updateQosHPA.Spec.MinReplicas, "/ max:", updateQosHPA.Spec.MaxReplicas, "]")

							if current_err != nil || qos_err != nil { //둘 중 하나라도 에러가 날 경우 롤백
								fmt.Println("current_err : ", current_err)
								fmt.Println("qos_err : ", qos_err)
								cluster_client.Update(context.TODO(), foundHPA)
								fmt.Println(">>> "+clustername+" Rollback HPA [ min:", *foundHPA.Spec.MinReplicas, "/ max:", foundHPA.Spec.MaxReplicas, "]")
								qos_cluster_client.Update(context.TODO(), foundQosHPA)
								fmt.Println(">>> "+qosCluster+" Rollback HPA [ min:", *foundQosHPA.Spec.MinReplicas, "/ max:", foundQosHPA.Spec.MaxReplicas, "]")
							}
						} else {
							//fmt.Println("OpenMCPHPA / HPA Changes")
							if hasInstance.Status.LastSpec.HpaTemplate.Spec.MinReplicas != nil {
								if *hasInstance.Status.LastSpec.HpaTemplate.Spec.MinReplicas != *hasInstance.Spec.HpaTemplate.Spec.MinReplicas || hasInstance.Status.LastSpec.HpaTemplate.Spec.MaxReplicas != hasInstance.Spec.HpaTemplate.Spec.MaxReplicas {
									desired_min_replicas := cluster_min_map[clustername]
									desired_max_replicas := cluster_max_map[clustername]

									if *foundHPA.Spec.MinReplicas != desired_min_replicas || foundHPA.Spec.MaxReplicas != desired_max_replicas {
										updateHPA := foundHPA
										updateHPA.Spec.MinReplicas = &desired_min_replicas
										updateHPA.Spec.MaxReplicas = desired_max_replicas

										err = cluster_client.Update(context.TODO(), updateHPA)
										//fmt.Println("min: ", desired_min_replicas, "max: ", desired_max_replicas)
										fmt.Println(">>> "+clustername+" Update HPA [ min:", *updateHPA.Spec.MinReplicas, "/ max:", updateHPA.Spec.MaxReplicas, "]")

										if err != nil {
											reqLogger.Error(err, "Failed to update HPA", "Hpa.Namespace", updateHPA.Namespace, "Hpa.Name", updateHPA.Name)
											return reconcile.Result{}, err
										}
									}
								}
							}
						}
					}
					// VPA 생성
				} else if hasInstance.Spec.VpaMode == "Always" || (hasInstance.Spec.VpaMode == "Auto" && cluster_dep_request[clustername] == false) {
					//fmt.Println("[ HPA VPA ]")
					// Check if the HPA already exists, if not create a new one
					//----------------------------HPA-------------------------------
					foundHPA := &hpav2beta2.HorizontalPodAutoscaler{}
					cluster_client := cm.Cluster_clients[clustername]
					err = cluster_client.Get(context.TODO(), foundHPA, hasInstance.Namespace, hasInstance.Name+"-hpa")

					if err != nil && errors.IsNotFound(err) { //CREATE HPA
						// Define a new HPA
						hpa_min := cluster_min_map[clustername]
						hpa_max := cluster_max_map[clustername]
						//fmt.Println(clustername, " - min : ", hpa_min, ", max : ", hpa_max)
						hpa := r.updateHorizontalPodAutoscalerWithVpa(request, hasInstance, hpa_min, hpa_max)

						//fmt.Println("hpa: ", hpa)
						reqLogger.Info("Creating a new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
						err = cluster_client.Create(context.TODO(), hpa)

						if err != nil {
							reqLogger.Error(err, "Failed to create new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
							return reconcile.Result{}, err
						}

						fmt.Println(">>> "+clustername+" Create HPA [ min:", *hpa.Spec.MinReplicas, "/ max:", hpa.Spec.MaxReplicas, "]")
						// HPA created successfully - return and requeue

						//Status Update
						hasInstance.Status.LastSpec = hasInstance.Spec

						err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
						if err_openmcp != nil {
							fmt.Println("Failed to update instance status", err)
							return reconcile.Result{}, err
						} else {
							fmt.Println(">>> OpenMCPHPA LastSpec Update (HPA Create)")
						}

						return reconcile.Result{Requeue: true}, nil
					} else if err != nil {
						reqLogger.Error(err, "Failed to get HPA")
						return reconcile.Result{}, err
					} else if err == nil { //UPDATE HPA - Rebalancing 조건을 만족했을 때 또는 min/max값이 수정되었을때
						if foundHPA.Spec.MaxReplicas == foundHPA.Status.CurrentReplicas && foundHPA.Status.CurrentReplicas == foundHPA.Status.DesiredReplicas { //||  (*foundHPA.Spec.MinReplicas == foundHPA.Status.CurrentReplicas && foundHPA.Status.CurrentReplicas == foundHPA.Status.DesiredReplicas) {
							fmt.Println(">>> " + clustername + " min/max rebalancing")

							//분석 엔진을 통해 얻은 결과
							//---------------------------------------------------------------------------------------------------
							SERVER_IP := os.Getenv("GRPC_SERVER")
							SERVER_PORT := os.Getenv("GRPC_PORT")

							grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

							hi := &protobuf.OpenMCPHPAInfo{DeploymentName: "test-pod", DeploymentNamespace: "test", CL: nil}
							hi.CL = make([]*protobuf.ClusterList, 0)
							//hi.CL = append(hi.CL, &protobuf.ClusterList{ClusterName:"cluster1"})
							hi.CL = append(hi.CL, &protobuf.ClusterList{ClusterName: "cluster2"})
							hi.CL = append(hi.CL, &protobuf.ClusterList{ClusterName: "cluster3"})

							rv := &protobuf.RequestValue{}
							rv.HpaInfo = make([]*protobuf.OpenMCPHPAInfo, 0)
							rv.HpaInfo = append(rv.HpaInfo, hi)

							r, gRPCerr := grpcClient.SendAnalysisResult(context.TODO(), rv)
							if gRPCerr != nil {
								fmt.Printf("could not connect : %v", gRPCerr)
							}
							fmt.Println("Anlysis Result:", r)
							//---------------------------------------------------------------------------------------------------
							qosCluster := r.TargetCluster
							//qosCluster := "cluster1"
							updateHPA := foundHPA
							updateHPA.Spec.MaxReplicas += 1

							qos_cluster_client := cm.Cluster_clients[qosCluster]
							foundQosHPA := &hpav2beta2.HorizontalPodAutoscaler{}
							err := qos_cluster_client.Get(context.TODO(), foundQosHPA, hasInstance.Namespace, hasInstance.Name+"-hpa")

							if err != nil {
								fmt.Println("err: ", err)
							}

							updateQosHPA := foundQosHPA
							updateQosHPA.Spec.MaxReplicas -= 1

							current_err := cluster_client.Update(context.TODO(), updateHPA)
							fmt.Println(">>> "+clustername+" Update HPA [ min:", *updateHPA.Spec.MinReplicas, "/ max:", updateHPA.Spec.MaxReplicas, "]")
							qos_err := qos_cluster_client.Update(context.TODO(), updateQosHPA)
							fmt.Println(">>> "+qosCluster+" Update HPA [ min:", *updateQosHPA.Spec.MinReplicas, "/ max:", updateQosHPA.Spec.MaxReplicas, "]")

							if current_err != nil || qos_err != nil { //둘 중 하나라도 에러가 날 경우 롤백
								fmt.Println("current_err : ", current_err)
								fmt.Println("qos_err : ", qos_err)
								cluster_client.Update(context.TODO(), foundHPA)
								qos_cluster_client.Update(context.TODO(), foundQosHPA)
							}

						} else {
							if hasInstance.Status.LastSpec.HpaTemplate.Spec.MinReplicas != nil {
								if *hasInstance.Status.LastSpec.HpaTemplate.Spec.MinReplicas != *hasInstance.Spec.HpaTemplate.Spec.MinReplicas || hasInstance.Status.LastSpec.HpaTemplate.Spec.MaxReplicas != hasInstance.Spec.HpaTemplate.Spec.MaxReplicas {
									desired_min_replicas := cluster_min_map[clustername]
									desired_max_replicas := cluster_max_map[clustername]

									if *foundHPA.Spec.MinReplicas != desired_min_replicas || foundHPA.Spec.MaxReplicas != desired_max_replicas {
										updateHPA := foundHPA
										updateHPA.Spec.MinReplicas = &desired_min_replicas
										updateHPA.Spec.MaxReplicas = desired_max_replicas

										err = cluster_client.Update(context.TODO(), updateHPA)
										fmt.Println(">>> "+clustername+" Update HPA [ min:", *updateHPA.Spec.MinReplicas, "/ max:", updateHPA.Spec.MaxReplicas, "]")

										if err != nil {
											reqLogger.Error(err, "Failed to update HPA", "Hpa.Namespace", updateHPA.Namespace, "Hpa.Name", updateHPA.Name)
											return reconcile.Result{}, err
										}
									}
								}
							}
						}
					}

					//----------------------------VPA-------------------------------
					foundVPA := &vpav1beta2.VerticalPodAutoscaler{}
					//cluster_client = cm.Cluster_clients[clustername]
					err = r.ghosts[clustername].Get(context.TODO(), ObjectKey{Namespace: hasInstance.Namespace, Name: hasInstance.Name + "-vpa"}, foundVPA)
					// err = cluster_client.Get(context.TODO(), foundVPA, hasInstance.Namespace, hasInstance.Name+"-vpa")

					if err != nil && errors.IsNotFound(err) { //CREATE VPA
						// Define a new VPA
						vpa := r.updateVerticalPodAutoscaler(request, hasInstance)

						reqLogger.Info("Creating a new VPA", "VPA.Namespace", vpa.Namespace, "VPA.Name", vpa.Name)
						err = r.ghosts[clustername].Create(context.TODO(), vpa)
						//err = cluster_client.Create(context.TODO(), vpa)

						if err != nil {
							reqLogger.Error(err, "Failed to create new VPA", "VPA.Namespace", vpa.Namespace, "VPA.Name", vpa.Name)
							return reconcile.Result{}, err
						}
						// VPA created successfully - return and requeue
						fmt.Println(">>> " + clustername + " Create VPA")

						hasInstance.Status.LastSpec = hasInstance.Spec

						err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
						if err_openmcp != nil {
							fmt.Println("Failed to update instance status", err)
							return reconcile.Result{}, err
						} else {
							fmt.Println(">>> OpenMCPHPA LastSpec Update (VPA Create)")
						}

						return reconcile.Result{Requeue: true}, nil
					} else if err != nil {
						fmt.Println("Failed to get VPA")
						reqLogger.Error(err, "Failed to get VPA")
						return reconcile.Result{}, err
					} else if err == nil { //UPDATE VPA
						//fmt.Println(">>> " + clustername + " Update VPA")
						/*vpa := r.updateVerticalPodAutoscaler(request, hasInstance)

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
				fmt.Println("Failed to update instance status", err)
				return reconcile.Result{}, err
			} else {
				fmt.Println(">>> OpenMCPHPA LastSpec Update (End)")
			}

		}
	} else if err != nil && errors.IsNotFound(err) {
		fmt.Println("OpenmcpDeployment is not found")
		reqLogger.Error(err, "OpenMCPDeployment doesn't exist - ", "openmcpDep.Namespace: ", openmcpDep.Namespace, ", openmcpDep.Name: ", openmcpDep.Name)
		return reconcile.Result{}, err
	} else {
		fmt.Println("Failed to get OpenMCPDeployment")
		reqLogger.Error(err, "Failed to get OpenMCPDeployment")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *reconciler) UpdateStatusLastspec(hasInstance *ketiv1alpha1.OpenMCPHybridAutoScaler) {
	hasInstance.Status.LastSpec = hasInstance.Spec

	err_openmcp := r.live.Status().Update(context.TODO(), hasInstance)
	if err_openmcp != nil {
		fmt.Println("Failed to update instance status", err_openmcp)
	} else {
		fmt.Println("Update OpenMCPHPA LastSpec : ", hasInstance.Status.LastSpec.HpaTemplate.Spec.MaxReplicas)
	}
}

func CreateMinMaxMap(r *reconciler, cm *ClusterManager, hasInstance *ketiv1alpha1.OpenMCPHybridAutoScaler, dep_list_for_hpa []string, cluster_dep_replicas map[string]int32) (map[string]int32, map[string]int32, *ketiv1alpha1.OpenMCPHybridAutoScaler, error) {
	cluster_min_map := make(map[string]int32)
	cluster_max_map := make(map[string]int32)
	var err error

	type ObjectKey = types.NamespacedName

	checkPolicy := 0
	if hasInstance.Status.Policies != nil { // 정책이 하나라도 있는지 확인
		for n, tmp := range hasInstance.Status.Policies {
			if tmp.Type == "Mode" { // Mode 이미 설정되어 있으면
				fmt.Println(">>> Policy \"Min Max Distribution\" Existed")
				if hasInstance.Status.Policies[n].Value[0] == "Equal" {
					//fmt.Println("Policy - min max Equal")
					cluster_min_map = cm.hpaMinScheduling_equal(dep_list_for_hpa, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = cm.hpaMaxScheduling_equal(dep_list_for_hpa, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
				} else {
					//fmt.Println("Policy - min max Unequal")
					cluster_min_map = cm.hpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = cm.hpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
				}
				checkPolicy = 1
				break
			}
		}
	}
	if checkPolicy == 0 { // Mode 설정 X
		foundPolicy := &ketiv1alpha1.OpenMCPPolicyEngine{}
		minmax_policy_err := r.live.Get(context.TODO(), ObjectKey{Namespace: "openmcp", Name: "hpa-minmax-distribution-mode"}, foundPolicy)
		if minmax_policy_err == nil {
			if foundPolicy.Spec.PolicyStatus == "Enabled" {
				fmt.Println(">>> Policy \"Min Max Distribution\" Apply (Enabled)")
				if foundPolicy.Spec.Template.Spec.Policies[0].Value[0] == "Equal" {
					//fmt.Println("New Policy - min max Equal")
					cluster_min_map = cm.hpaMinScheduling_equal(dep_list_for_hpa, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = cm.hpaMaxScheduling_equal(dep_list_for_hpa, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
					hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				} else if foundPolicy.Spec.Template.Spec.Policies[0].Value[0] == "Unequal" {
					//fmt.Println("New Policy - min max Unequal")
					cluster_min_map = cm.hpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = cm.hpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
					hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				} else {
					//fmt.Println("Default - min max Unequal")
					cluster_min_map = cm.hpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
					cluster_max_map = cm.hpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
					foundPolicy.Spec.Template.Spec.Policies[0].Value[0] = "Default"
					hasInstance.Status.Policies = append(hasInstance.Status.Policies, foundPolicy.Spec.Template.Spec.Policies...)
				}

			} else {
				fmt.Println(">>> Policy \"Min Max Distribution\" Apply (Disabled - set default)")
				cluster_min_map = cm.hpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
				cluster_max_map = cm.hpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
				omp := make([]ketiv1alpha1.OpenMCPPolicies, 1)
				omp_value := make([]string, 1)
				omp_value[0] = "Default"
				omp[0].Type = "Mode"
				omp[0].Value = omp_value
				hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
			}
		} else {
			fmt.Println("policy_err : ", minmax_policy_err)
			fmt.Println("Fail to get policy \"Min Max Distribution\" (set default)")
			cluster_min_map = cm.hpaMinScheduling(dep_list_for_hpa, cluster_dep_replicas, *hasInstance.Spec.HpaTemplate.Spec.MinReplicas)
			cluster_max_map = cm.hpaMaxScheduling(dep_list_for_hpa, cluster_dep_replicas, hasInstance.Spec.HpaTemplate.Spec.MaxReplicas)
			omp := make([]ketiv1alpha1.OpenMCPPolicies, 1)
			omp_value := make([]string, 1)
			omp_value[0] = "Default"
			omp[0].Type = "Mode"
			omp[0].Value = omp_value
			hasInstance.Status.Policies = append(hasInstance.Status.Policies, omp...)
		}

		//정책 업데이트
		err = r.live.Status().Update(context.TODO(), hasInstance)
		if err != nil {
			fmt.Println("Policy Status Update Error")
		} else {
			fmt.Println("Policy Status UPDATE Success")
		}

	}

	return cluster_min_map, cluster_max_map, hasInstance, err
}

func (r *reconciler) DeleteOpenMCPHPA(cm *ClusterManager, namespace string, hpaName string, vpaName string) {
	hpa := &hpav2beta2.HorizontalPodAutoscaler{}
	vpa := &vpav1beta2.VerticalPodAutoscaler{}
	type ObjectKey = types.NamespacedName

	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_clients[cluster.Name]

		err1 := cluster_client.Delete(context.Background(), hpa, namespace, hpaName)
		if err1 != nil && errors.IsNotFound(err1) {
			//fmt.Println("Fail to Delete HPA - ", err1)
		} else if err1 != nil {
			fmt.Println("Fail to Delete HPA - ", err1)
		} else if err1 == nil {
			fmt.Println(">>> " + cluster.Name + " Delete HPA")
		}

		err2 := r.ghosts[cluster.Name].Get(context.TODO(), ObjectKey{Namespace: namespace, Name: vpaName}, vpa)
		if err2 != nil && errors.IsNotFound(err2) {
			//fmt.Println("Fail to Get VPA - ", err2)
		} else if err2 != nil {
			fmt.Println("Fail to Get VPA - ", err2)
		} else if err2 == nil {
			err3 := r.ghosts[cluster.Name].Delete(context.Background(), vpa)
			if err3 != nil && errors.IsNotFound(err3) {
				//fmt.Println("Fail to Delete VPA - ", err3)
			} else if err3 != nil {
				fmt.Println("Fail to Delete VPA - ", err3)
			} else if err3 == nil {
				fmt.Println(">>> " + cluster.Name + " Delete VPA")
			}
		}

	}
}

func (r *reconciler) DeleteHPAVPA(cm *ClusterManager, clustername []string, namespace string, hpaName string, vpaName string) {
	hpa := &hpav2beta2.HorizontalPodAutoscaler{}
	vpa := &vpav1beta2.VerticalPodAutoscaler{}
	type ObjectKey = types.NamespacedName

	for _, cluster := range clustername {
		cluster_client := cm.Cluster_clients[cluster]

		err1 := cluster_client.Delete(context.Background(), hpa, namespace, hpaName)
		if err1 != nil && errors.IsNotFound(err1) {
			//fmt.Println("Fail to Delete HPA - ", err1)
		} else if err1 != nil {
			fmt.Println("Fail to Delete HPA - ", err1)
		} else if err1 == nil {
			fmt.Println(">>> " + cluster + " Delete HPA")
		}

		err2 := r.ghosts[cluster].Get(context.TODO(), ObjectKey{Namespace: namespace, Name: vpaName}, vpa)
		if err2 != nil && errors.IsNotFound(err2) {
			//fmt.Println("Fail to Get VPA - ", err2)
		} else if err2 != nil {
			fmt.Println("Fail to Get VPA - ", err2)
		} else if err2 == nil {
			err3 := r.ghosts[cluster].Delete(context.Background(), vpa)
			if err3 != nil && errors.IsNotFound(err3) {
				//fmt.Println("Fail to Delete VPA - ", err3)
			} else if err3 != nil {
				fmt.Println("Fail to Delete VPA - ", err3)
			} else if err3 == nil {
				fmt.Println(">>> " + cluster + " Delete VPA")
			}
		}

	}
}

func (cm *ClusterManager) hpaMinScheduling(dep_list_for_hpa []string, cluster_dep_replicas map[string]int32, min int32) map[string]int32 {

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

func (cm *ClusterManager) hpaMaxScheduling(dep_list_for_hpa []string, cluster_dep_replicas map[string]int32, max int32) map[string]int32 {

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

func (cm *ClusterManager) hpaMinScheduling_equal(dep_list_for_hpa []string, min int32) map[string]int32 {

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

func (cm *ClusterManager) hpaMaxScheduling_equal(dep_list_for_hpa []string, max int32) map[string]int32 {

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

func (r *reconciler) updateHorizontalPodAutoscaler(req reconcile.Request, m *ketiv1alpha1.OpenMCPHybridAutoScaler, min int32, max int32) *hpav2beta2.HorizontalPodAutoscaler {
	ls := labelsForHpa(m.Name)

	hpa := &hpav2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-hpa",
			Namespace: m.Namespace,
			Labels:    ls,
		},

		Spec: m.Spec.HpaTemplate.Spec,
	}

	hpa.Spec.MinReplicas = &min
	hpa.Spec.MaxReplicas = max
	hpa.Spec.ScaleTargetRef.APIVersion = "apps/v1"
	hpa.Spec.ScaleTargetRef.Kind = "Deployment"
	hpa.Spec.ScaleTargetRef.Name = hpa.Spec.ScaleTargetRef.Name + "-deploy"

	reference.SetMulticlusterControllerReference(hpa, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return hpa
}

func (r *reconciler) updateHorizontalPodAutoscalerWithVpa(req reconcile.Request, m *ketiv1alpha1.OpenMCPHybridAutoScaler, min int32, max int32) *hpav2beta2.HorizontalPodAutoscaler {
	ls := labelsForHpa(m.Name)

	hpa := &hpav2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-hpa",
			Namespace: m.Namespace,
			Labels:    ls,
		},

		Spec: m.Spec.HpaTemplate.Spec,
	}

	if min < 2 {
		*hpa.Spec.MinReplicas = 2
	} else {
		hpa.Spec.MinReplicas = &min
	}
	hpa.Spec.MaxReplicas = max
	hpa.Spec.ScaleTargetRef.APIVersion = "apps/v1"
	hpa.Spec.ScaleTargetRef.Kind = "Deployment"
	hpa.Spec.ScaleTargetRef.Name = hpa.Spec.ScaleTargetRef.Name + "-deploy"

	reference.SetMulticlusterControllerReference(hpa, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return hpa
}

func (r *reconciler) updateVerticalPodAutoscaler(req reconcile.Request, m *ketiv1alpha1.OpenMCPHybridAutoScaler) *vpav1beta2.VerticalPodAutoscaler {
	ls := labelsForHpa(m.Name)
	vpaUpdateMode := vpav1beta2.UpdateModeAuto

	vpa := &vpav1beta2.VerticalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-vpa",
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: vpav1beta2.VerticalPodAutoscalerSpec{
			TargetRef: &autoscaling.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       m.Spec.HpaTemplate.Spec.ScaleTargetRef.Name + "-deploy",
			},
			UpdatePolicy: &vpav1beta2.PodUpdatePolicy{
				UpdateMode: &vpaUpdateMode,
			},
		},
	}
	reference.SetMulticlusterControllerReference(vpa, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))
	// Set Memcached instance as the owner and controller
	//controllerutil.SetControllerReference(m, vpa, r.scheme)
	return vpa
}

func labelsForHpa(name string) map[string]string {
	return map[string]string{"app": "openmcphybridautoscaler", "openmcphybridautoscaler_cr": name}
}

type ClusterManager struct {
	Fed_namespace   string
	Host_config     *rest.Config
	Host_client     genericclient.Client
	Cluster_list    *fedv1b1.KubeFedClusterList
	Cluster_configs map[string]*rest.Config
	Cluster_clients map[string]genericclient.Client
}

func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		fmt.Println("Error retrieving list of federated clusters: %+v", err)
	}
	if len(clusterList.Items) == 0 {
		fmt.Println("No federated clusters found")
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
	//fmt.Println("Step 0.	NewClusterManager()")
	fed_namespace := "kube-federation-system"
	host_config, _ := rest.InClusterConfig()
	host_client := genericclient.NewForConfigOrDie(host_config)
	cluster_list := ListKubeFedClusters(host_client, fed_namespace)
	cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)
	cluster_clients := KubeFedClusterClients(cluster_list, cluster_configs)

	//fmt.Println("cluster_list : ", cluster_list)

	cm := &ClusterManager{
		Fed_namespace:   fed_namespace,
		Host_config:     host_config,
		Host_client:     host_client,
		Cluster_list:    cluster_list,
		Cluster_configs: cluster_configs,
		Cluster_clients: cluster_clients,
	}
	return cm
}
