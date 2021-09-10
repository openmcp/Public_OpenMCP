package analyticResource

import (
	"context"
	"fmt"
	"openmcp/openmcp/openmcp-analytic-engine/src/protobuf"
	"openmcp/openmcp/util/clusterManager"
	"os"
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
	autoscalingTime      time.Time
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
			var CPADeployList = protobuf.CPADeployList{}
			for _, v := range CPAInfoList {
				v.OmcpdeployInfo.ReplicasNum = v.ReplicasAfterScaling
				CPADeployList.CPADeployInfo = append(CPADeployList.CPADeployInfo, v.OmcpdeployInfo)

				//구현
				/*if v.ReplicasAfterScaling > v.CpaMax {
					//pod 줄이기
				}else if v.ReplicasAfterScaling < v.CpaMin {
					//pod 늘리기
				}*/
			}
			grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)
			grpcResponse, gRPC_err := grpcClient.SendCPAAnalysis(context.TODO(), &CPADeployList)

			if gRPC_err == nil {
				result := grpcResponse.ResponseCPADeploy
				if len(result) > 0 {
					fmt.Println("===============================================================")
					fmt.Println("CPAInfoList : ", CPAInfoList)
					for _, deploy := range result {
						cpaKey := CPAKey{}
						cpaKey.HASName = deploy.CPAName
						cpaKey.HASNamespace = deploy.Namespace
						lastTime := CPAInfoList[cpaKey].autoscalingTime
						if lastTime.IsZero() || (!lastTime.IsZero() && time.Since(lastTime) > time.Second*300) {
							if strings.Contains(deploy.PodState, "Warning") {
								todo := deploy.Action
								targetCluster := deploy.TargetCluster
								beforeScaling := CPAInfoList[cpaKey].ReplicasAfterScaling

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
										if beforeScaling >= CPAInfoList[cpaKey].CpaMin {
											*dep.Spec.Replicas -= 1
											depupdate_err := cluster_client.Update(context.TODO(), dep)
											if depupdate_err != nil {
												fmt.Println("UpdateDeployReplicas err : ", depupdate_err)
												continue Exit
											} else {
												fmt.Println("[" + deploy.Name + "] Success to Scale-in")
												tmptime := CPAInfoList[cpaKey]
												tmptime.autoscalingTime = time.Now()
												CPAInfoList[cpaKey] = tmptime
											}

										} else {
											fmt.Println("!![error] Fail to Scale-in (CurrentReplicas = MinReplicas)")
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
												fmt.Println("[" + deploy.Name + "] Success to Scale-out")
												tmptime := CPAInfoList[cpaKey]
												tmptime.autoscalingTime = time.Now()
												CPAInfoList[cpaKey] = tmptime
											}

										} else {
											fmt.Println("!![error] Fail to Scale-out (CurrentReplicas = MaxReplicas)")
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
										fmt.Println("!![error] Fail to Get Deployment Replicas")
										continue Exit
									}
								}
								if totalReplicas != beforeScaling {
									fmt.Println("!![error] Fail to Update Replicas")
								} else {
									fmt.Println("[" + deploy.Name + "] Success to Update Replicas")
									tmp := CPAInfoList[cpaKey]
									tmp.ReplicasAfterScaling = totalReplicas
									CPAInfoList[cpaKey] = tmp
								}
							}
						} else {
							fmt.Println("Already Done. Please Wait ...")
						}
					}
					fmt.Println("===============================================================")
				}
			} else {
				fmt.Println(gRPC_err)
			}
		}

		time.Sleep(time.Second * 5)
	}
}
