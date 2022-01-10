package syncWeight

import (
	//"fmt"

	"os"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	//"admiralty.io/multicluster-controller/pkg/cluster"
	"context"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-analytic-engine/src/protobuf"
	"openmcp/openmcp/util/clusterManager"
	"strings"
	"time"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var DistributeList map[types.NamespacedName]drInfo
var NodeList map[string]nodeInfo // key : CLUSTERNAME/NODENAME

type drInfo struct {
	DRScore map[string][]DRWeight
}

type DRWeight struct {
	ToRegionZone    string
	InitScore       int
	ConvertToWeight int
}

type nodeInfo struct {
	node_region string
	node_zone   string
}

//모든 서비스의 DestinationRule 상시 업데이트
func SyncDRWeight(myClusterManager *clusterManager.ClusterManager, quit, quitok chan bool) {
	cm := myClusterManager

	DistributeList = map[types.NamespacedName]drInfo{}

	SERVER_IP := os.Getenv("GRPC_SERVER")
	SERVER_PORT := os.Getenv("GRPC_PORT")

	grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)


	for {

		select {
		case <-quit:
			omcplog.V(2).Info("SyncWeight Quit")
			quitok <- true
			return
		default:
			omcplog.V(2).Info("Analyzing DestinationRule Weight ...")

			//모든 노드들의 region_zone 리스트
			ocList, err := cm.Crd_client.OpenMCPCluster("openmcp").List(v1.ListOptions{})
			NodeList = map[string]nodeInfo{}

			if err != nil {
				omcplog.V(0).Info("OpenMCPClusterList err : ", err)
			}
			for _, oc := range ocList.Items {
				clustername := oc.Name
				if oc.Spec.JoinStatus == "JOIN" {
					region := oc.Spec.NodeInfo.Region
					zones := strings.Split(oc.Spec.NodeInfo.Zone, ",")

					for i := 0; i < len(zones); i++ {
						tmp := nodeInfo{
							node_region: region,
							node_zone:   zones[i],
						}

						//cluster_node := clustername + "/" + node.Name

						NodeList[clustername] = tmp
					}
				}
			}

			/*//모든 노드들의 region_zone 리스트
			for _, memberCluster := range cm.Cluster_list.Items {

				clusterName := memberCluster.Name

				nodeList := &corev1.NodeList{}
				err_node := cm.Cluster_genClients[clusterName].List(context.TODO(), nodeList, "default")

				if err_node != nil {
					omcplog.V(2).Info("err_node : ", err_node)
				} else {
					for _, node := range nodeList.Items {
						region := node.Labels["topology.kubernetes.io/region"]
						zone := node.Labels["topology.kubernetes.io/zone"]

						tmp := nodeInfo{
							node_region: region,
							node_zone:   zone,
						}

						cluster_node := clusterName + "/" + node.Name

						NodeList[cluster_node] = tmp
					}
				}
			}*/

			node_list_all := map[string]int{}

			for _, node := range NodeList {
				region_zone := node.node_region + "/" + node.node_zone + "/*"
				node_list_all[region_zone] = 1
			}

			distributeListTm := map[types.NamespacedName]drInfo{}

			//OpenMCPService 조회 ->  Pod Selector 조회 -> 배포된 노드 정보 가져오기
			//osvcList := &resourcev1alpha1.OpenMCPServiceList{}
			//err_osvc := liveclient.List(context.TODO(), osvcList)
			osvcList, err_osvc := cm.Crd_client.OpenMCPService(corev1.NamespaceAll).List(v1.ListOptions{})

			if err_osvc == nil {
				omcplog.V(2).Info("osvcList : ", osvcList)
			} else {
				omcplog.V(2).Info(err_osvc)
			}

			for _, osvc := range osvcList.Items {
				target_node_list := map[string][]string{}

				for _, mcluster := range cm.Cluster_list.Items {
					podNodeName := ""
					tmp_pod_list := []string{}
					cluster_client := cm.Cluster_genClients[mcluster.Name]

					//omcplog.V(2).Info("***********",osvc.Name," selector :", osvc.Spec.Template.Spec.Selector)

					if cm.Cluster_genClients[mcluster.Name] != nil && len(osvc.Spec.Template.Spec.Selector) != 0 {

						listOption := &client.ListOptions{
							LabelSelector: labels.SelectorFromSet(
								osvc.Spec.Template.Spec.Selector,
							),
						}

						podList := &corev1.PodList{}
						err_pod := cluster_client.List(context.TODO(), podList, osvc.Namespace, listOption)

						if err_pod == nil {
							for _, pod := range podList.Items {
								tmp_pod_list = append(tmp_pod_list, pod.Name)
								podNodeName = pod.Spec.NodeName
							}
						} else {
							omcplog.V(2).Info(err_pod)
						}

						if podNodeName != "" {
							cluster_node := mcluster.Name // + "/" + podNodeName
							region_zone := NodeList[cluster_node].node_region + "/" + NodeList[cluster_node].node_zone + "/" + mcluster.Name

							target_node_list[region_zone] = append(target_node_list[region_zone], tmp_pod_list...)
						}
					} else {
						omcplog.V(0).Info("err!! : ", mcluster.Name, " cluster client has 'nil'")
					}
				}

				omcplog.V(2).Info("[",osvc.Name,"] target_node_list : ", target_node_list)
				if len(target_node_list) > 0 {
					tmp_rz := map[string][]DRWeight{}
					for rz, _ := range node_list_all { //from
						tmp_score := []DRWeight{}
						for pn, podlist := range target_node_list { //to
							average_pod_score := 0
							tmp_length := len(podlist)
							for _, podname := range podlist {
								s := analyzeScore(podname, osvc.Namespace, rz, pn, grpcClient)
								average_pod_score += s

								if s == 0 {
									tmp_length -= 1
								}
							}
							initScore := 0
							if tmp_length != 0 {
								initScore = average_pod_score / tmp_length

							}
							d := DRWeight{
								ToRegionZone:    pn,
								InitScore:       initScore,
								ConvertToWeight: 0,
							}
							tmp_score = append(tmp_score, d)
						}
						//score -> weight 변환
						var totalscore int
						totalscore = 0
						for _, target := range tmp_score {
							totalscore += target.InitScore
						}
						//if totalscore > 0 {

						for i, target := range tmp_score {
							var f float32
							if totalscore == 0 {
								tmp_score[i].ConvertToWeight = 0
							} else {
								f = float32(target.InitScore) / float32(totalscore)
								tmp_score[i].ConvertToWeight = int(f * 100)
							}

							//fmt.Println("[before] init : ", target.InitScore, "total : ", totalscore, "weight : ", tmp_score[i].ConvertToWeight)
						}

						totalweight := 0

						maxscore := 0
						maxindex := 0
						for i, target := range tmp_score {
							if target.ConvertToWeight == 0 {
								tmp_score[i].ConvertToWeight = 1
							}
							if maxscore <= target.ConvertToWeight {
								maxscore = target.ConvertToWeight
								maxindex = i
							}
							totalweight += tmp_score[i].ConvertToWeight
						}

						if totalweight > 0 && totalweight < 100 {
							for i := 0; i < 100-totalweight; i++ {
								a := i % len(tmp_score)
								tmp_score[a].ConvertToWeight += 1
							}
						} else if totalweight > 100 {
							tmp_score[maxindex].ConvertToWeight -= totalweight - 100
						}
						//}

						tmp_rz[rz] = tmp_score
					}
					osvc_n_ns := types.NamespacedName{
						Namespace: osvc.Namespace,
						Name:      osvc.Name,
					}
					dl := drInfo{
						DRScore: tmp_rz,
					}

					//distributeListTm 갱신
					distributeListTm[osvc_n_ns] = dl

					//Update DestinationRules

					drweight := map[string][]DRWeight{}
					drweight = distributeListTm[osvc_n_ns].DRScore

					distributeTarget := map[string]map[string]uint32{}

					for key, value := range drweight {
						tmp := map[string]uint32{}
						for _, w := range value {
							tmp[w.ToRegionZone] = uint32(w.ConvertToWeight)
						}
						distributeTarget[key] = tmp
					}

					distribute := []*networkingv1alpha3.LocalityLoadBalancerSetting_Distribute{}

					for k, v := range distributeTarget {
						tmp_dis := &networkingv1alpha3.LocalityLoadBalancerSetting_Distribute{
							From: k,
							To:   v,
						}
						distribute = append(distribute, tmp_dis)
					}

					obj_dr, err_get_dr := cm.Crd_istio_client.DestinationRule(osvc.Namespace).Get(osvc.Name, v1.GetOptions{})
					//obj_dr := &v1alpha3.DestinationRule{}
					//err_get_dr := liveclient.Get(context.TODO(), osvc_n_ns, obj_dr)

					if err_get_dr == nil {

						tmp_dr := &v1alpha3.DestinationRule{
							TypeMeta: v1.TypeMeta{
								Kind:       "DestinationRule",
								APIVersion: "networking.istio.io/v1alpha3",
							},
							ObjectMeta: v1.ObjectMeta{
								Name:      obj_dr.Name,
								Namespace: obj_dr.Namespace,
							},
							Spec: obj_dr.Spec,
						}
						tmp_dr.ResourceVersion = obj_dr.ResourceVersion

						obj_dr.Spec.TrafficPolicy.LoadBalancer.LocalityLbSetting.Distribute = distribute

						_, err_update_dr := cm.Crd_istio_client.DestinationRule(osvc.Namespace).Update(tmp_dr)
						//err_update_dr := liveclient.Update(context.TODO(), obj_dr)

						if err_update_dr != nil {
							omcplog.V(2).Info(err_update_dr)
						} else {
							omcplog.V(2).Info("update dr - ", osvc_n_ns)
						}
					} else {
						omcplog.V(2).Info(err_get_dr)
					}
				} else {
					osvc_n_ns := types.NamespacedName{
						Namespace: osvc.Namespace,
						Name:      osvc.Name,
					}

					obj_dr, err_get_dr := cm.Crd_istio_client.DestinationRule(osvc.Namespace).Get(osvc.Name, v1.GetOptions{})

					if err_get_dr == nil {

						tmp_dr := &v1alpha3.DestinationRule{
							TypeMeta: v1.TypeMeta{
								Kind:       "DestinationRule",
								APIVersion: "networking.istio.io/v1alpha3",
							},
							ObjectMeta: v1.ObjectMeta{
								Name:      obj_dr.Name,
								Namespace: obj_dr.Namespace,
							},
							Spec: obj_dr.Spec,
						}
						tmp_dr.ResourceVersion = obj_dr.ResourceVersion

						obj_dr.Spec.TrafficPolicy.LoadBalancer.LocalityLbSetting.Distribute = nil

						_, err_update_dr := cm.Crd_istio_client.DestinationRule(osvc.Namespace).Update(tmp_dr)
						//err_update_dr := liveclient.Update(context.TODO(), obj_dr)

						if err_update_dr != nil {
							omcplog.V(2).Info(err_update_dr)
						} else {
							omcplog.V(2).Info("delete distribute in dr - ", osvc_n_ns)
						}
					} else {
						omcplog.V(2).Info(err_get_dr)
					}
				}

			}

			DistributeList = distributeListTm

			omcplog.V(2).Info(">>> Update All DestinationRule Resources ")
			//fmt.Println(DistributeList)

			time.Sleep(time.Second * 10)

		}

	}
}

func analyzeScore(podname string, namespace string, from string, to string, grpcClient protobuf.RequestAnalysisClient) int {
	slice_from := strings.Split(from, "/")
	slice_to := strings.Split(to, "/")

	rzinfo := &protobuf.RegionZoneInfo{
		FromRegion:    slice_from[0],
		FromZone:      slice_from[1],
		ToRegion:      slice_to[0],
		ToZone:        slice_to[1],
		ToClusterName: slice_to[2],
		ToNamespace:   namespace,
		ToPodName:     podname,
	}

	result, err_grpc := grpcClient.SendRegionZoneInfo(context.TODO(), rzinfo)

	tmp_result := 0

	if err_grpc != nil {
		omcplog.V(2).Info("[ ",podname," err ] ",err_grpc)
	} else {
		tmp_result = int(result.Weight)
	}

	return tmp_result
}
