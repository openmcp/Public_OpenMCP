package analyticResource

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"openmcp/openmcp/openmcp-analytic-engine/src/protobuf"
	"openmcp/openmcp/util/clusterManager"
	"os"
	"strconv"
	"strings"
	"time"

	"admiralty.io/multicluster-controller/pkg/cluster"
	appsv1 "k8s.io/api/apps/v1"
)

type CPAValue struct {
	OmcpdeployInfo *protobuf.CPADeployInfo
	//	InitReplicas         int32
	ReplicasAfterScaling int32
	CpaMin               int32
	CpaMax               int32
	AutoscalingTime      map[string]time.Time
}

type CPAKey struct {
	HASName      string
	HASNamespace string
}

var CPAInfoList = make(map[CPAKey]CPAValue)

func CalcPodMetrics(cm *clusterManager.ClusterManager, live *cluster.Cluster) {
	fmt.Println("Start CPA ... CPAList_ClusterNum -> [", len(CPAInfoList), "]")
	SERVER_IP := os.Getenv("GRPC_SERVER")
	SERVER_PORT := os.Getenv("GRPC_PORT")
Exit:
	for {
		if len(CPAInfoList) > 0 {
			fmt.Println("*** CPAInfoList => ", CPAInfoList)
			fmt.Println("")

			var CPADeployList = protobuf.CPADeployList{}
			for _, v := range CPAInfoList {
				v.OmcpdeployInfo.ReplicasNum = v.ReplicasAfterScaling
				CPADeployList.CPADeployInfo = append(CPADeployList.CPADeployInfo, v.OmcpdeployInfo)
			}

			grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)
			grpcResponse, gRPC_err := grpcClient.SendCPAAnalysis(context.TODO(), &CPADeployList)

			if gRPC_err == nil {
				result := grpcResponse.ResponseCPADeploy
				if len(result) > 0 {
					//Get cpa period policy
					var period_int64 int64
					foundPolicy, period_policy_err := cm.Crd_client.OpenMCPPolicy("openmcp").Get("cpa-period", metav1.GetOptions{})

					if period_policy_err != nil {
						period_int64 = 300
					} else {
						cpa_policy := foundPolicy.Spec.Template.Spec.Policies
						period := cpa_policy[0].Value[0]
						period_int64, _ = strconv.ParseInt(period, 10, 64)
					}

					fmt.Println("===============================================================")
					fmt.Println("Receive CPA Result from Analytic Engine (" + strconv.FormatInt(period_int64, 10) + "s)")

					for _, deploy := range result {
						fmt.Println("")
						cpaKey := CPAKey{}
						cpaKey.HASName = deploy.CPAName
						cpaKey.HASNamespace = deploy.Namespace
						lastTime := CPAInfoList[cpaKey].AutoscalingTime[deploy.TargetCluster]

						if lastTime.IsZero() || (!lastTime.IsZero() && time.Since(lastTime) > time.Second*time.Duration(period_int64)) {
							if strings.Contains(deploy.PodState, "Warning") {
								todo := deploy.Action
								targetCluster := deploy.TargetCluster

								//beforeScaling := CPAInfoList[cpaKey].ReplicasAfterScaling
								var beforeScaling int32
								beforeScaling = 0
								for _, cluster := range CPAInfoList[cpaKey].OmcpdeployInfo.Clusters {
									dep := &appsv1.Deployment{}
									cluster_client := cm.Cluster_genClients[cluster]
									dep_err := cluster_client.Get(context.TODO(), dep, deploy.Namespace, deploy.Name)
									if dep_err == nil {
										beforeScaling = beforeScaling + *dep.Spec.Replicas
									} else {
										fmt.Println("!! [", deploy.Name, "] Fail to Get Deployment Replicas")
									}
								}

								dep := &appsv1.Deployment{}
								cluster_client := cm.Cluster_genClients[targetCluster]
								dep_err := cluster_client.Get(context.TODO(), dep, deploy.Namespace, deploy.Name)

								if dep_err != nil {
									fmt.Println("getTargetDeploy err : ", dep_err)
									continue Exit
								} else {
									if todo == "Scale-in" {
										//replicas 축소 -1
										beforeScaling -= 1
										if beforeScaling >= CPAInfoList[cpaKey].CpaMin && *dep.Spec.Replicas > 1 {
											*dep.Spec.Replicas -= 1
											depupdate_err := cluster_client.Update(context.TODO(), dep)
											if depupdate_err != nil {
												fmt.Println("UpdateDeployReplicas err : ", depupdate_err)
												continue Exit
											} else {
												//fmt.Println("[" + deploy.Name + "] Success to Scale-in in ", targetCluster)
												fmt.Println("[", deploy.TargetCluster, "] Success to Scale-in Deployment '"+deploy.Name+"'")
												tmptime := CPAInfoList[cpaKey]
												if tmptime.AutoscalingTime == nil {
													scalingTime := make(map[string]time.Time)
													tmptime.AutoscalingTime = scalingTime
												}
												tmptime.AutoscalingTime[deploy.TargetCluster] = time.Now()
												CPAInfoList[cpaKey] = tmptime
											}

										} else {
											if beforeScaling < CPAInfoList[cpaKey].CpaMin {
												fmt.Println("!! [", deploy.TargetCluster, "] Fail to Scale-in (CurrentReplicas = MinReplicas)")
											} else if *dep.Spec.Replicas <= 1 {
												fmt.Println("!! [", deploy.TargetCluster, "] Fail to Scale-in (CurrentReplicas = 1)")
											}
										}
									} else if todo == "Scale-out" {
										//replicas 증가 +1
										beforeScaling += 1
										if beforeScaling <= CPAInfoList[cpaKey].CpaMax {
											*dep.Spec.Replicas += 1
											depupdate_err := cluster_client.Update(context.TODO(), dep)
											if depupdate_err != nil {
												fmt.Println("UpdateDeployReplicas err : ", depupdate_err)
												continue Exit
											} else {
												//fmt.Println("[" + deploy.Name + "] Success to Scale-out in ", targetCluster)
												fmt.Println("[", deploy.TargetCluster, "] Success to Scale-out Deployment '"+deploy.Name+"'")
												tmptime := CPAInfoList[cpaKey]
												if tmptime.AutoscalingTime == nil {
													scalingTime := make(map[string]time.Time)
													tmptime.AutoscalingTime = scalingTime
												}
												tmptime.AutoscalingTime[deploy.TargetCluster] = time.Now()
												CPAInfoList[cpaKey] = tmptime
											}

										} else if beforeScaling == CPAInfoList[cpaKey].CpaMax+1 {
											r_cluster_list := deploy.RestCluster

											if len(r_cluster_list) == 0 {
												fmt.Println("!! [", deploy.TargetCluster, "] Fail to Scale-out. There is no rest cluster.")
												//continue Exit
											}

											for _, r_cluster := range r_cluster_list {
												if r_cluster != targetCluster {

													//pod 감소
													r_dep := &appsv1.Deployment{}
													r_cluster_client := cm.Cluster_genClients[r_cluster]
													r_dep_err := r_cluster_client.Get(context.TODO(), r_dep, deploy.Namespace, deploy.Name)

													if *r_dep.Spec.Replicas > 1 {
														fmt.Println("*** CurrentReplicas == MaxReplicas. Get Replicas from other clusters.")

														if r_dep_err != nil {
															fmt.Println("getTargetDeploy err : ", r_dep_err)
															continue Exit
														} else {
															*r_dep.Spec.Replicas -= 1
															r_depupdate_err := r_cluster_client.Update(context.TODO(), r_dep)
															if r_depupdate_err != nil {
																fmt.Println("UpdateDeployReplicas err : ", r_depupdate_err)
																continue Exit
															} else {
																//fmt.Println("[" + deploy.Name + "] Success to Scale-in in ", r_cluster)
																fmt.Println("[", r_cluster, "] Success to Scale-in Deployment '"+deploy.Name+"'")
															}
														}

														//pod 증가
														*dep.Spec.Replicas += 1
														depupdate_err := cluster_client.Update(context.TODO(), dep)
														if depupdate_err != nil {
															fmt.Println("UpdateDeployReplicas err : ", depupdate_err)
															*r_dep.Spec.Replicas += 1
															r_depupdate_err2 := r_cluster_client.Update(context.TODO(), r_dep)
															if r_depupdate_err2 != nil {
																fmt.Println("UpdateDeployReplicas err : ", r_depupdate_err2)
																continue Exit
															}
															continue Exit
														} else {
															//fmt.Println("[" + deploy.Name + "] Success to Scale-in in ", targetCluster)
															fmt.Println("[", deploy.TargetCluster, "] Success to Scale-out Deployment '"+deploy.Name+"'")
															tmptime := CPAInfoList[cpaKey]
															if tmptime.AutoscalingTime == nil {
																scalingTime := make(map[string]time.Time)
																tmptime.AutoscalingTime = scalingTime
															}
															tmptime.AutoscalingTime[deploy.TargetCluster] = time.Now()
															CPAInfoList[cpaKey] = tmptime
														}

														fmt.Println("Success to send replica from", r_cluster, "to", targetCluster)
														break
													}
												}
											}
										}
									}
								}

								//deploy 돌면서 총 Replicas 다시 검사
								var totalReplicas int32
								totalReplicas = 0
								for _, cluster := range CPAInfoList[cpaKey].OmcpdeployInfo.Clusters {
									dep := &appsv1.Deployment{}
									cluster_client := cm.Cluster_genClients[cluster]
									dep_err := cluster_client.Get(context.TODO(), dep, deploy.Namespace, deploy.Name)
									if dep_err == nil {
										totalReplicas = totalReplicas + *dep.Spec.Replicas
									} else {
										fmt.Println("!! [", deploy.Name, "] Fail to Get Deployment Replicas")
										continue Exit
									}
								}
								/*if totalReplicas != beforeScaling {
									fmt.Println("!![error] Fail to Update Replicas")
								} else {*/
								fmt.Println("[", deploy.Name, "] Update Replicas (", totalReplicas, "/", CPAInfoList[cpaKey].CpaMax, ")")
								tmp := CPAInfoList[cpaKey]
								tmp.ReplicasAfterScaling = totalReplicas
								CPAInfoList[cpaKey] = tmp
								//}
							}
						} else {
							fmt.Println("[", deploy.TargetCluster, "] Scaling Already Done. Please Wait ...")
						}
					}
					fmt.Println("===============================================================")
					fmt.Println("")
				}
			} else {
				fmt.Println(gRPC_err)
			}
		}

		time.Sleep(time.Second * 5)
	}
}
