package analyticEngine

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-analytic-engine/src/influx"
	"openmcp/openmcp/openmcp-analytic-engine/src/protobuf"
	"openmcp/openmcp/util/clusterManager"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	"sigs.k8s.io/kubefed/pkg/controller/util"
)

type AnalyticEngineStruct struct {
	Influx               influx.Influx
	MetricsWeight        map[string]float64            // MetricsWeight["CPU"] = value
	ResourceScore        map[string]float64            // ResourceScore["clusterName"] = value
	ClusterResourceUsage map[string]map[string]float64 // ClusterResourceUsage["clusterName]["CPU"] = value
	//ClusterGeo              map[string]map[string]string
	NetworkInfos            map[string]map[string]*NetworkInfo
	ClusterPodResourceScore map[string]map[string]map[string]float64 // ClusterPodResourceScore[clusterName][NS]["Pod"] = value
	//ClusterSVCResourceScore map[string]map[string]map[string]float64 // ClusterSVCResourceScore[clusterName][NS]["SVC"] = value
	GeoScore []float64 // 0: RegionZoneMatchedScore, 1: OnlyRegionMatcehdScore, 2: NoRegionZoneMatchedScore
	mutex    *sync.Mutex
}

// Network is used to get real-time network information (receive data, transmit data)
// Calculating the difference between previous_data and next_data is needed to get real-time network data
// because the data from Kubelet is cumulative data
type NetworkInfo struct {
	prev_rx int64
	prev_tx int64
	next_rx int64
	next_tx int64
}

func NewAnalyticEngine(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD string) *AnalyticEngineStruct {
	omcplog.V(4).Info("Func NewAnalyticEngine Called")
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println("NewAnalyticEngine")
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	ae := &AnalyticEngineStruct{}
	ae.Influx = *influx.NewInflux(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)
	ae.ResourceScore = make(map[string]float64)
	ae.ClusterResourceUsage = make(map[string]map[string]float64)
	ae.MetricsWeight = make(map[string]float64)
	ae.NetworkInfos = make(map[string]map[string]*NetworkInfo)
	ae.ClusterPodResourceScore = make(map[string]map[string]map[string]float64)
	//ae.ClusterSVCResourceScore = make(map[string]map[string]map[string]float64)
	ae.GeoScore = []float64{-1, -1, -1}
	ae.mutex = &sync.Mutex{}

	return ae
}

func (ae *AnalyticEngineStruct) CalcResourceScore(cm *clusterManager.ClusterManager, quit, quitok chan bool) {
	omcplog.V(4).Info("Func CalcResourceScore Called")
	//cm := clusterManager.NewClusterManager()
	// ae.MetricsWeight = make(map[string]float64)
	// ae.ClusterGeo = map[string]map[string]string{}
	// ae.NetworkInfos = make(map[string]map[string]*NetworkInfo)

	for {
		select {
		case <-quit:
			omcplog.V(2).Info("CalcResourceScore Quit")
			quitok <- true
			return
		default:
			omcplog.V(2).Info("Cluster Metric Score Refresh")
			/*
				Update {
					ae.MetricsWeight
				}
			*/
			err := ae.getMetricsWeight(cm)
			if err != nil {
				omcplog.V(0).Info(err)
			}
			var wg sync.WaitGroup
			wg.Add(len(cm.Cluster_list.Items))

			for _, cluster := range cm.Cluster_list.Items {
				go func(cluster fedv1b1.KubeFedCluster) {
					defer wg.Done()

					if cm.Cluster_genClients[cluster.Name] != nil {
						/*
							Update {
								ae.ClusterResourceUsage
								ae.ResourceScore
							}
						*/
						err = ae.UpdateScore(cluster.Name, cm)
						if err != nil {
							omcplog.V(0).Info(err)
						}
						/*
							Update {
								ae.ClusterPodResourceScore
								ae.GeoScore
							}
						*/
						err = ae.UpdateClusterPodScore(cluster.Name, cm)
						if err != nil {
							omcplog.V(0).Info(err)
						}
						/*
							Update {
								ae.ClusterSVCResourceScore
								ae.GeoScore
							}
						*/
						// err = ae.UpdateClusterSVCScore(cluster.Name, cm)
						// if err != nil {
						// 	omcplog.V(0).Info(err)
						// }

						/*
							Update {
								ae.ClusterGeo
							}
						*/
						// err = ae.updateClusterGeo(cluster, cm)
						// if err != nil {
						// 	omcplog.V(0).Info(err)
						// }

						// Update Network Data from InfluxDB
						/*
							Update {
								ae.NetworkInfos
							}
						*/
						ae.UpdateNetworkData(cluster, cm)

					} else {
						omcplog.V(0).Info("err!! : ", cluster.Name, " cluster client has 'nil'")
					}

				}(cluster)

			}

			wg.Wait()
			time.Sleep(2 * time.Second)
		}

	}
}
func (ae *AnalyticEngineStruct) getMetricsWeight(cm *clusterManager.ClusterManager) error {
	//Get metric-weight Policy----------------------------
	openmcpPolicyInstance, err := cm.Crd_client.OpenMCPPolicy("openmcp").Get("analytic-metrics-weight", metav1.GetOptions{})

	if err != nil {
		return err
	} else {
		policies := openmcpPolicyInstance.Spec.Template.Spec.Policies
		for _, policy := range policies {
			value, _ := strconv.ParseFloat(policy.Value[0], 64)
			ae.MetricsWeight[policy.Type] = value
		}
		omcplog.V(3).Info("metricsWeight : ", ae.MetricsWeight)
	}
	return nil
}
func getNodeList(cluster fedv1b1.KubeFedCluster, cm *clusterManager.ClusterManager) (*corev1.NodeList, error) {

	config, err := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return nodeList, nil
}

// func (ae *AnalyticEngineStruct) updateClusterGeo(cluster fedv1b1.KubeFedCluster, cm *clusterManager.ClusterManager) error {
// 	nodeList, err := getNodeList(cluster, cm)
// 	if len(nodeList.Items) == 0 {
// 		return errors.New("No Exist Nodes in cluster'" + cluster.Name + "'")
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	node := nodeList.Items[0]

// 	//Extract zone, region from Label
// 	ae.ClusterGeo[cluster.Name] = map[string]string{}
// 	ae.ClusterGeo[cluster.Name]["Region"] = node.Labels["topology.kubernetes.io/region"]
// 	ae.ClusterGeo[cluster.Name]["Zone"] = node.Labels["topology.kubernetes.io/zone"]

// 	// ae.ClusterGeo[cluster.Name]["Country"] = node.Labels["failure-domain.beta.kubernetes.io/zone"]
// 	// ae.ClusterGeo[cluster.Name]["Continent"] = node.Labels["failure-domain.beta.kubernetes.io/region"]

// 	return nil

// }

// Update Network data from InfluxDB
func (ae *AnalyticEngineStruct) UpdateNetworkData(cluster fedv1b1.KubeFedCluster, cm *clusterManager.ClusterManager) error {
	omcplog.V(4).Info("Func UpdateNetworkData Called")
	// Initialize cluster's network data
	clusterName := cluster.Name
	nodeList, err := getNodeList(cluster, cm)
	if err != nil {
		return err
	}
	_, exists := ae.NetworkInfos[clusterName]
	if !exists {
		newClusterInfo := make(map[string]*NetworkInfo)
		ae.NetworkInfos[clusterName] = newClusterInfo
	}

	// Update Node's network data
	for _, node := range nodeList.Items {

		// Initialize Node's network data
		_, exists := ae.NetworkInfos[clusterName][node.Name]
		if !exists {
			ae.NetworkInfos[clusterName][node.Name] = &NetworkInfo{}
		}

		// Get cumulative network data from InfluxDB
		result := ae.Influx.GetNetworkData(clusterName, node.Name)

		// If data is not stored, cannot calculate real-time network data
		if len(result) == 0 {
			continue
		}

		for _, ser := range result[0].Series {

			// Second row is previous data because the result is ordered by desc
			prev_rx, _ := strconv.ParseInt(fmt.Sprintf("%s", ser.Values[1][1]), 10, 64)
			prev_tx, _ := strconv.ParseInt(fmt.Sprintf("%s", ser.Values[1][2]), 10, 64)

			// First row is next data
			next_rx, _ := strconv.ParseInt(fmt.Sprintf("%s", ser.Values[0][1]), 10, 64)
			next_tx, _ := strconv.ParseInt(fmt.Sprintf("%s", ser.Values[0][2]), 10, 64)

			// Update network data on NetworkInfo structure
			ae.NetworkInfos[clusterName][node.Name].prev_rx = prev_rx
			ae.NetworkInfos[clusterName][node.Name].prev_tx = prev_tx
			ae.NetworkInfos[clusterName][node.Name].next_rx = next_rx
			ae.NetworkInfos[clusterName][node.Name].next_tx = next_tx
		}
	}
	return nil
}

func (ae *AnalyticEngineStruct) UpdateScore(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(4).Info("Func UpdateScore Called")
	var score float64 = 0
	scoreMap := make(map[string]float64)

	result, err := ae.Influx.GetClusterMetricsData(clusterName)
	if err != nil {
		return err
	}
	MetricsMap := make(map[string]float64)
	prevMetricsMap := make(map[string]float64)
	var totalCpuCore int64 = 0

	time_t := ""

	if len(result) > 0 {
		for _, ser := range result[0].Series {
			if cm.Cluster_genClients[clusterName] != nil {
				nodeCapacity := &corev1.Node{}
				if clusterGenClient, ok := cm.Cluster_genClients[ser.Tags["cluster"]]; ok {
					err := clusterGenClient.Get(context.TODO(), nodeCapacity, "", ser.Tags["node"])

					if err != nil {
						omcplog.V(0).Info("nodelist err!  : ", err)
						continue
					} else {
						omcplog.V(2).Info("[CPU Capacity] ", ser.Tags["cluster"], "/", ser.Tags["node"], "/", nodeCapacity.Status.Capacity.Cpu().Value())
					}

					totalCpuCore = totalCpuCore + nodeCapacity.Status.Capacity.Cpu().Value()

					for c, colName := range ser.Columns {
						for r, _ := range ser.Values {

							Strval := fmt.Sprintf("%v", ser.Values[r][c])
							QuanVal, _ := resource.ParseQuantity(Strval)

							if r == 0 {
								if colName == "NetworkLatency" {
									MetricsMap[colName], _ = strconv.ParseFloat(Strval, 64)
								} else if colName == "time" {
									time_t = Strval
								} else {
									if _, ok := MetricsMap[colName]; ok {
										if colName == "CPUUsageNanoCores" {
											MetricsMap[colName] = MetricsMap[colName] + float64(QuanVal.AsApproximateFloat64())
										} else {
											MetricsMap[colName] = MetricsMap[colName] + float64(QuanVal.Value())
										}

									} else {
										if colName == "CPUUsageNanoCores" {
											MetricsMap[colName] = MetricsMap[colName] + float64(QuanVal.AsApproximateFloat64())
										} else {
											MetricsMap[colName] = float64(QuanVal.Value())
										}
									}
								}

							} else if r == 4 {
								if _, ok := prevMetricsMap[colName]; ok {
									prevMetricsMap[colName] = prevMetricsMap[colName] + float64(QuanVal.Value())
								} else {
									prevMetricsMap[colName] = float64(QuanVal.Value())
								}
							}
						}
					}
				} else {
					return nil
				}

			} else {
				omcplog.V(0).Info("err!! : ", clusterName, " cluster client has 'nil'")
			}
		}

		cpuScore := MetricsMap["CPUUsageNanoCores"] / float64(totalCpuCore) * 100
		memScore := (MetricsMap["MemoryUsageBytes"] / (MetricsMap["MemoryUsageBytes"] + MetricsMap["MemoryAvailableBytes"])) * 100
		netScore := ((MetricsMap["NetworkRxBytes"] - prevMetricsMap["NetworkRxBytes"]) + (MetricsMap["NetworkTxBytes"] - prevMetricsMap["NetworkTxBytes"])) / 1000000
		diskScore := (MetricsMap["FsUsedBytes"] / MetricsMap["FsCapacityBytes"]) * 100
		latencyScore := MetricsMap["NetworkLatency"] * 1000

		if len(result[0].Series) != 0 {
			ae.Influx.InsertClusterStatus(clusterName, time_t, cpuScore, memScore, netScore, diskScore, latencyScore)
		}

		score = cpuScore*ae.MetricsWeight["CPU"] + memScore*ae.MetricsWeight["Memory"] + diskScore*ae.MetricsWeight["FS"] + netScore*ae.MetricsWeight["NET"] + latencyScore*ae.MetricsWeight["LATENCY"]

		scoreMap["cpu"] = cpuScore
		scoreMap["memory"] = memScore
		scoreMap["network"] = netScore
		scoreMap["disk"] = diskScore
		scoreMap["latency"] = latencyScore

		// if score > 0 && clusterName == "eks-cluster1" {
		// 	omcplog.V(2).Info("--------------------------------------------------------")
		// 	omcplog.V(2).Info(" rx -> ", MetricsMap["NetworkRxBytes"]-prevMetricsMap["NetworkRxBytes"])
		// 	omcplog.V(2).Info(" tx -> ", MetricsMap["NetworkTxBytes"]-prevMetricsMap["NetworkTxBytes"])
		// 	omcplog.V(2).Info(" cpuscore : ", cpuScore, MetricsMap["CPUUsageNanoCores"], float64(totalCpuCore))
		// 	omcplog.V(2).Info(" memScore : ", memScore)
		// 	omcplog.V(2).Info(" netScore : ", netScore)
		// 	omcplog.V(2).Info(" diskScore : ", diskScore)
		// 	omcplog.V(2).Info(" latencyScore : ", latencyScore)
		// 	omcplog.V(2).Info("\"", clusterName, "\" totalScore : ", score)
		// 	omcplog.V(2).Info("--------------------------------------------------------")
		// }

		ae.mutex.Lock()
		ae.ResourceScore[clusterName] = score
		ae.mutex.Unlock()
		ae.ClusterResourceUsage[clusterName] = scoreMap
	}

	return nil
}
func (ae *AnalyticEngineStruct) UpdateClusterPodScore(clusterName string, cm *clusterManager.ClusterManager) error {

	var geoRate float64
	var resourceRate float64
	var RegionZoneMatchedScore float64
	var OnlyRegionMatchedScore float64
	var NoRegionZoneMatchedScore float64

	ClusterPodResourceScore := make(map[string]map[string]map[string]float64)

	// ae.GeoScore Update
	openmcpPolicyInstance, err := cm.Crd_client.OpenMCPPolicy("openmcp").Get("lb-scoring-weight", metav1.GetOptions{})
	if err != nil {
		omcplog.V(0).Info(err)
		return err
	}

	policies := openmcpPolicyInstance.Spec.Template.Spec.Policies

	for _, policy := range policies {
		value, _ := strconv.ParseFloat(policy.Value[0], 64)
		if policy.Type == "GeoRate" {
			geoRate = value
			resourceRate = 1 - geoRate
		} else if policy.Type == "RegionZoneMatchedScore" {
			RegionZoneMatchedScore = value
		} else if policy.Type == "OnlyRegionMatchedScore" {
			OnlyRegionMatchedScore = value
		} else if policy.Type == "NoRegionZoneMatchedScore" {
			NoRegionZoneMatchedScore = value
		}

	}

	ae.GeoScore[0] = geoRate * RegionZoneMatchedScore
	ae.GeoScore[1] = geoRate * OnlyRegionMatchedScore
	ae.GeoScore[2] = geoRate * NoRegionZoneMatchedScore

	// ae.ClusterPodResourceScore Update
	_, exists := ClusterPodResourceScore[clusterName]
	if !exists {
		tmp := make(map[string]map[string]float64)
		ClusterPodResourceScore[clusterName] = tmp
	}

	openmcpPolicyInstance, err = cm.Crd_client.OpenMCPPolicy("openmcp").Get("analytic-metrics-weight", metav1.GetOptions{})
	if err != nil {
		omcplog.V(0).Info(err)
		return err
	}

	policies = openmcpPolicyInstance.Spec.Template.Spec.Policies

	var CpuWeight float64
	var MemWeight float64
	for _, policy := range policies {
		value, _ := strconv.ParseFloat(policy.Value[0], 64)
		if policy.Type == "CPU" {
			CpuWeight = value
		} else if policy.Type == "Memory" {
			MemWeight = value
		}

	}

	nodeCpuCapacity := make(map[string]map[string]int64)
	nodeMemCapacity := make(map[string]map[string]int64)

	_, exists = nodeCpuCapacity[clusterName]
	if !exists {
		tempCpuCapacityMap := make(map[string]int64)
		nodeCpuCapacity[clusterName] = tempCpuCapacityMap
	}
	_, exists = nodeMemCapacity[clusterName]
	if !exists {
		tempMemoryCapacityMap := make(map[string]int64)
		nodeMemCapacity[clusterName] = tempMemoryCapacityMap
	}

	nodeList := &corev1.NodeList{}
	if clusterGenClient, ok := cm.Cluster_genClients[clusterName]; ok {
		err = clusterGenClient.List(context.TODO(), nodeList, "default")
		if err != nil {
			omcplog.V(0).Info(err)
			return err
		}
		for _, node := range nodeList.Items {
			nodeCpuCapacity[clusterName][node.Name] = node.Status.Capacity.Cpu().Value()
			nodeMemCapacity[clusterName][node.Name] = node.Status.Capacity.Memory().Value()
		}

		svcList := &corev1.ServiceList{}
		if clusterGenClient, ok := cm.Cluster_genClients[clusterName]; ok {
			err = clusterGenClient.List(context.TODO(), svcList, corev1.NamespaceAll)
			if err != nil {
				omcplog.V(0).Info(err)
				return err
			}
			for _, svc := range svcList.Items {

				ns := svc.Namespace

				_, exists = ClusterPodResourceScore[clusterName][ns]
				if !exists {
					tmp2 := make(map[string]float64)
					ClusterPodResourceScore[clusterName][ns] = tmp2
				}

				matchLabels := client.MatchingLabels{}

				for key, value := range svc.Spec.Selector {
					matchLabels[key] = value
				}
				podList := &corev1.PodList{}
				if clusterGenClient, ok := cm.Cluster_genClients[clusterName]; ok {
					err = clusterGenClient.List(context.TODO(), podList, ns, matchLabels)
					if err != nil {
						omcplog.V(0).Info(err)
						return err
					}

					for _, pod := range podList.Items {
						result, err := ae.Influx.GetClusterPodsData(clusterName, pod.Name)
						if err != nil {
							omcplog.V(0).Info(err)
							return err
						}
						if len(result) == 0 {
							return errors.New("Error : Influx DB Data length is 0")
						}
						for _, ser := range result[0].Series {
							for r, _ := range ser.Values {
								var nodeName string
								var podCpuUsage float64
								var podMemUsage float64
								for c, colName := range ser.Columns {
									Strval := fmt.Sprintf("%v", ser.Values[r][c])
									QuanVal, _ := resource.ParseQuantity(Strval)

									if colName == "CPUUsageNanoCores" {
										//MetricsMap[colName], _ = strconv.ParseFloat(Strval, 64)
										podCpuUsage = QuanVal.AsApproximateFloat64()

										// fmt.Println("check1 - CPUUsageNanoCores: ", Strval, QuanVal, QuanVal.AsApproximateFloat64())
									} else if colName == "MemoryUsageBytes" {
										//time_t = Strval
										// fmt.Println("check2 - MemoryUsageBytes: ", Strval, QuanVal, QuanVal.Value())
										podMemUsage = float64(QuanVal.Value())
									} else if colName == "node" {
										nodeName = Strval
									}
								}
								podCpuIdle := 100 - ((podCpuUsage)/float64(nodeCpuCapacity[clusterName][nodeName]))*100
								podMemIdle := 100 - ((podMemUsage)/float64(nodeMemCapacity[clusterName][nodeName]))*100

								// podCpuIdleTotal = podCpuIdleTotal + podCpuIdle
								// podMemIdleTotal = podMemIdlenTotal + podMemIdle

								// 상수값을 빼서 작은 cpu/mem 변화에도 가중치 차이를 둠
								podCpuIdleRescale := podCpuIdle // - 99
								podMemIdleRescale := podMemIdle // - 99

								// if podCpuIdleRescale < 0 {
								// 	podCpuIdleRescale = 0
								// }
								// if podMemIdleRescale < 0 {
								// 	podMemIdleRescale = 0
								// }
								score := (podCpuIdleRescale*CpuWeight + podMemIdleRescale*MemWeight) * resourceRate
								ClusterPodResourceScore[clusterName][ns][pod.Name] = score

								// fmt.Println("Cluster: ", clusterName)
								// fmt.Println("Node: ", nodeName)
								// fmt.Println("Pod: ", pod.Name)
								// fmt.Println("podCpuIdle: ", podCpuIdle)
								// fmt.Println("podMemIdle: ", podMemIdle)

								if clusterName == "cluster04" || clusterName == "cluster09" || clusterName == "cluster17" {
									if strings.Contains(pod.Name, "productpage") {
										omcplog.V(2).Info("clusterName:", clusterName)
										omcplog.V(2).Info("nodeName:", nodeName)
										omcplog.V(2).Info("nodeCpuCapacity[clusterName][nodeName]:", nodeCpuCapacity[clusterName][nodeName])
										omcplog.V(2).Info("podCpuUsage:", podCpuUsage)
										omcplog.V(2).Info("podCpuIdle:", podCpuIdle)
										omcplog.V(2).Info("nodeMemCapacity[clusterName][nodeName]:", nodeMemCapacity[clusterName][nodeName])
										omcplog.V(2).Info("podMemUsage:", podMemUsage)
										omcplog.V(2).Info("podMemIdle:", podMemIdle)
										omcplog.V(2).Info("score:", score)
									}

								}

							}
						}
					}
				} else {
					return nil
				}

			}

			_, exists = ae.ClusterPodResourceScore[clusterName]
			if !exists {
				tmp := make(map[string]map[string]float64)
				ae.ClusterPodResourceScore[clusterName] = tmp
			}
			ae.ClusterPodResourceScore[clusterName] = ClusterPodResourceScore[clusterName]

			omcplog.V(2).Info("ae.GeoScore:", ae.GeoScore)
			omcplog.V(2).Info("["+clusterName+"] ae.ClusterPodResourceScore:", ae.ClusterPodResourceScore[clusterName])
		}

	}

	return nil
}

// func (ae *AnalyticEngineStruct) UpdateClusterSVCScore(clusterName string, cm *clusterManager.ClusterManager) error {

// 	var geoRate float64
// 	var resourceRate float64
// 	var RegionZoneMatchedScore float64
// 	var OnlyRegionMatchedScore float64
// 	var NoRegionZoneMatchedScore float64

// 	// ae.GeoScore Update
// 	openmcpPolicyInstance, err := cm.Crd_client.OpenMCPPolicy("openmcp").Get("lb-scoring-weight", metav1.GetOptions{})
// 	if err != nil {
// 		return err
// 	}

// 	policies := openmcpPolicyInstance.Spec.Template.Spec.Policies

// 	for _, policy := range policies {
// 		value, _ := strconv.ParseFloat(policy.Value[0], 64)
// 		if policy.Type == "GeoRate" {
// 			geoRate = value
// 			resourceRate = 1 - geoRate
// 		} else if policy.Type == "RegionZoneMatchedScore" {
// 			RegionZoneMatchedScore = value
// 		} else if policy.Type == "OnlyRegionMatchedScore" {
// 			OnlyRegionMatchedScore = value
// 		} else if policy.Type == "NoRegionZoneMatchedScore" {
// 			NoRegionZoneMatchedScore = value
// 		}

// 	}

// 	ae.GeoScore[0] = geoRate * RegionZoneMatchedScore
// 	ae.GeoScore[1] = geoRate * OnlyRegionMatchedScore
// 	ae.GeoScore[2] = geoRate * NoRegionZoneMatchedScore

// 	// ae.ClusterSVCResourceScore Update
// 	_, exists := ae.ClusterSVCResourceScore[clusterName]
// 	if !exists {
// 		tempClusterSVCResourceScore := make(map[string]map[string]float64)
// 		ae.ClusterSVCResourceScore[clusterName] = tempClusterSVCResourceScore
// 	}

// 	openmcpPolicyInstance, err = cm.Crd_client.OpenMCPPolicy("openmcp").Get("analytic-metrics-weight", metav1.GetOptions{})
// 	if err != nil {
// 		return err
// 	}

// 	policies = openmcpPolicyInstance.Spec.Template.Spec.Policies

// 	var CpuWeight float64
// 	var MemWeight float64
// 	for _, policy := range policies {
// 		value, _ := strconv.ParseFloat(policy.Value[0], 64)
// 		if policy.Type == "CPU" {
// 			CpuWeight = value
// 		} else if policy.Type == "Memory" {
// 			MemWeight = value
// 		}

// 	}

// 	nodeCpuCapacity := make(map[string]map[string]int64)
// 	nodeMemCapacity := make(map[string]map[string]int64)

// 	_, exists = nodeCpuCapacity[clusterName]
// 	if !exists {
// 		tempCpuCapacityMap := make(map[string]int64)
// 		nodeCpuCapacity[clusterName] = tempCpuCapacityMap
// 	}
// 	_, exists = nodeMemCapacity[clusterName]
// 	if !exists {
// 		tempMemoryCapacityMap := make(map[string]int64)
// 		nodeMemCapacity[clusterName] = tempMemoryCapacityMap
// 	}

// 	nodeList := &corev1.NodeList{}
// 	err = cm.Cluster_genClients[clusterName].List(context.TODO(), nodeList, "default")
// 	if err != nil {
// 		return err
// 	}
// 	for _, node := range nodeList.Items {
// 		nodeCpuCapacity[clusterName][node.Name] = node.Status.Capacity.Cpu().Value()
// 		nodeMemCapacity[clusterName][node.Name] = node.Status.Capacity.Memory().Value()
// 	}

// 	svcList := &corev1.ServiceList{}
// 	err = cm.Cluster_genClients[clusterName].List(context.TODO(), svcList, corev1.NamespaceAll)
// 	if err != nil {
// 		return err
// 	}

// 	for _, svc := range svcList.Items {
// 		var podCpuUtilizationTotal float64
// 		var podMemUtilizationTotal float64

// 		ns := svc.Namespace

// 		_, exists = ae.ClusterSVCResourceScore[clusterName][ns]
// 		if !exists {
// 			tempClusterSVCResourceScore2 := make(map[string]float64)
// 			ae.ClusterSVCResourceScore[clusterName][ns] = tempClusterSVCResourceScore2
// 		}

// 		matchLabels := client.MatchingLabels{}

// 		for key, value := range svc.Spec.Selector {
// 			matchLabels[key] = value
// 		}
// 		podList := &corev1.PodList{}
// 		err = cm.Cluster_genClients[clusterName].List(context.TODO(), podList, ns, matchLabels)
// 		if err != nil {
// 			return err
// 		}

// 		for _, pod := range podList.Items {
// 			result, err := ae.Influx.GetClusterPodsData(clusterName, pod.Name)
// 			if err != nil {
// 				return err
// 			}
// 			if len(result) == 0 {
// 				return errors.New("Error : Influx DB Data length is 0")
// 			}
// 			for _, ser := range result[0].Series {
// 				for r, _ := range ser.Values {
// 					var nodeName string
// 					var podCpuUsage float64
// 					var podMemUsage float64
// 					for c, colName := range ser.Columns {
// 						Strval := fmt.Sprintf("%v", ser.Values[r][c])
// 						QuanVal, _ := resource.ParseQuantity(Strval)

// 						if colName == "CPUUsageNanoCores" {
// 							//MetricsMap[colName], _ = strconv.ParseFloat(Strval, 64)
// 							podCpuUsage = QuanVal.AsApproximateFloat64()

// 							// fmt.Println("check1 - CPUUsageNanoCores: ", Strval, QuanVal, QuanVal.AsApproximateFloat64())
// 						} else if colName == "MemoryUsageBytes" {
// 							//time_t = Strval
// 							// fmt.Println("check2 - MemoryUsageBytes: ", Strval, QuanVal, QuanVal.Value())
// 							podMemUsage = float64(QuanVal.Value())
// 						} else if colName == "node" {
// 							nodeName = Strval
// 						}
// 					}
// 					podCpuUtilization := 100 - (podCpuUsage/float64(nodeCpuCapacity[clusterName][nodeName]))*100
// 					podMemUtilization := 100 - (podMemUsage/float64(nodeMemCapacity[clusterName][nodeName]))*100

// 					podCpuUtilizationTotal = podCpuUtilizationTotal + podCpuUtilization
// 					podMemUtilizationTotal = podMemUtilizationTotal + podMemUtilization

// 					fmt.Println("Cluster: ", clusterName)
// 					fmt.Println("Node: ", nodeName)
// 					fmt.Println("Pod: ", pod.Name)
// 					fmt.Println("podCpuUtilization: ", podCpuUtilization)
// 					fmt.Println("podMemUtilization: ", podMemUtilization)

// 				}
// 			}
// 		}
// 		podCpuUtilizationAvg := podCpuUtilizationTotal / float64(len(podList.Items))
// 		podMemUtilizationAvg := podMemUtilizationTotal / float64(len(podList.Items))

// 		score := (podCpuUtilizationAvg*CpuWeight + podMemUtilizationAvg*MemWeight) * resourceRate
// 		ae.ClusterSVCResourceScore[clusterName][ns][svc.Name] = score
// 	}

// 	fmt.Println("ae.GeoScore:", ae.GeoScore)
// 	fmt.Println("["+clusterName+"] ae.ClusterSVCResourceScore:", ae.ClusterSVCResourceScore[clusterName])

// 	return nil
// }
func (ae *AnalyticEngineStruct) SendLBAnalysis(ctx context.Context, in *protobuf.LBInfo) (*protobuf.ResponseLB, error) {
	omcplog.V(4).Info("Func SendLBAnalysis Called")

	// clusterScoreMap := ae.ResourceScore
	clusterScoreMap := make(map[string]float64)
	ae.mutex.Lock()
	copier.Copy(clusterScoreMap, ae.ResourceScore)
	ae.mutex.Unlock()

	omcplog.V(2).Info("LB Response")
	omcplog.V(3).Info(clusterScoreMap)
	return &protobuf.ResponseLB{ScoreMap: clusterScoreMap}, nil
}

func (ae *AnalyticEngineStruct) SelectHPACluster(data *protobuf.HASInfo) []string {
	omcplog.V(4).Info("Func SelectHPACluster Called")

	scoreMap := map[float64]string{}
	scoreT := []float64{}

	// clusterScoreMap := ae.ResourceScore
	clusterScoreMap := make(map[string]float64)
	ae.mutex.Lock()
	copier.Copy(clusterScoreMap, ae.ResourceScore)
	ae.mutex.Unlock()

	for key, value := range clusterScoreMap {
		if key != data.ClusterName && value > 0 {
			scoreMap[value] = key
			scoreT = append(scoreT, value)
		}
	}

	sort.Float64s(scoreT)

	filteringCluster := []string{}

	for i := 0; i < len(scoreT); i++ {
		if i == 2 {
			break
		}
		filteringCluster = append(filteringCluster, scoreMap[scoreT[i]])
	}

	omcplog.V(5).Info(filteringCluster)

	return filteringCluster
}

func (ae *AnalyticEngineStruct) CompareHPAMaxInfo(clusterList []string, data *protobuf.HASInfo) string {
	omcplog.V(4).Info("Func CompareHPAMaxInfo Called")
	replicasGap := map[string]int32{}
	rebalancingCount := map[string]int32{}

	for _, cluster := range clusterList {
		omcplog.V(3).Info(cluster, " hpa MaxReplicas : ", data.HPAMinORMaxReplicas[cluster], " / CurrentReplicas : ", data.HPACurrentReplicas[cluster])
		calc := data.HPAMinORMaxReplicas[cluster] - data.HPACurrentReplicas[cluster]
		if calc > 0 {
			replicasGap[cluster] = calc
		}
	}

	for cluster, _ := range replicasGap {
		rebalancingCount[cluster] = data.HASRebalancingCount[cluster]
	}

	omcplog.V(3).Info("desiredReplicas : ", replicasGap)
	omcplog.V(3).Info("countRebalancing : ", rebalancingCount)

	result := ""

	for cluster, _ := range rebalancingCount {
		if result == "" {
			result = cluster
		} else {
			if (replicasGap[result] - rebalancingCount[result]) < (replicasGap[cluster] - rebalancingCount[cluster]) {
				result = cluster
			}
		}
	}

	return result
}

func (ae *AnalyticEngineStruct) CompareHPAMinInfo(clusterList []string, data *protobuf.HASInfo) string {
	omcplog.V(4).Info("Func CompareHPAMinInfo Called")
	replicasGap := map[string]int32{}
	rebalancingCount := map[string]int32{}

	timeStart_analysis := time.Now()

	for _, cluster := range clusterList {
		omcplog.V(3).Info(cluster, " hpa MinReplicas : ", data.HPAMinORMaxReplicas[cluster], " / CurrentReplicas : ", data.HPACurrentReplicas[cluster])
		calc := data.HPACurrentReplicas[cluster] - data.HPAMinORMaxReplicas[cluster]
		if calc > 0 {
			replicasGap[cluster] = calc
		}
	}

	timeEnd_analysis := time.Since(timeStart_analysis)

	omcplog.V(3).Info("[2] GetHPAInfo \t\t\t", timeEnd_analysis)

	timeStart_analysis2 := time.Now()
	for cluster, _ := range replicasGap {
		rebalancingCount[cluster] = data.HASRebalancingCount[cluster]
	}

	omcplog.V(3).Info("desiredReplicas : ", replicasGap)
	omcplog.V(3).Info("countRebalancing : ", rebalancingCount)

	timeEnd_analysis2 := time.Since(timeStart_analysis2)
	omcplog.V(3).Info("[3] GetHASInfo \t\t\t", timeEnd_analysis2)

	result := ""
	timeStart_analysis3 := time.Now()
	for cluster, _ := range rebalancingCount {
		if result == "" {
			result = cluster
		} else {
			if (replicasGap[result] - rebalancingCount[result]) < (replicasGap[cluster] - rebalancingCount[cluster]) {
				result = cluster
			}
		}
	}

	timeEnd_analysis3 := time.Since(timeStart_analysis3)
	omcplog.V(3).Info("[4] CompareQoSScore \t\t", timeEnd_analysis3)

	return result
}
func (ae *AnalyticEngineStruct) SendRegionZoneInfo(ctx context.Context, data *protobuf.RegionZoneInfo) (*protobuf.ResponseWeight, error) {
	omcplog.V(4).Info("Func SendRegionZoneInfo Called")

	var tempGeoScore float64
	var tempClusterPodResourceScore float64

	if data.FromRegion == data.ToRegion && data.FromZone == data.ToZone {
		if ae.GeoScore[0] == -1 {
			return nil, errors.New("GeoScore Initiailzing. Retry again")
		}
		tempGeoScore = ae.GeoScore[0]
	} else if data.FromRegion == data.ToRegion {
		if ae.GeoScore[1] == -1 {
			return nil, errors.New("GeoScore Initiailzing. Retry again")
		}
		tempGeoScore = ae.GeoScore[1]
	} else {
		if ae.GeoScore[2] == -1 {
			return nil, errors.New("GeoScore Initiailzing. Retry again")
		}
		tempGeoScore = ae.GeoScore[2]
	}

	//fmt.Println(ae.ClusterPodResourceScore)
	_, exists := ae.ClusterPodResourceScore[data.ToClusterName]
	if !exists {
		return nil, errors.New("Not Exist Cluster '" + data.ToClusterName + "' - ClusterPodResourceScore is Initializing. Try again")
	}

	_, exists = ae.ClusterPodResourceScore[data.ToClusterName][data.ToNamespace]
	if !exists {
		return nil, errors.New("Not Exist Namespace '" + data.ToNamespace + "' - ClusterPodResourceScore is Initializing. Try again")
	}

	_, exists = ae.ClusterPodResourceScore[data.ToClusterName][data.ToNamespace][data.ToPodName]
	if !exists {
		return nil, errors.New("Not Exist Pod '" + data.ToPodName + "' - ClusterPodResourceScore is Initializing or Pod has Pending State. Try again")
	}

	tempClusterPodResourceScore = ae.ClusterPodResourceScore[data.ToClusterName][data.ToNamespace][data.ToPodName]

	weight := float32(tempGeoScore + tempClusterPodResourceScore)

	omcplog.V(4).Info("FromRegion: ", data.FromRegion)
	omcplog.V(4).Info("ToRegion: ", data.ToRegion)
	omcplog.V(4).Info("FromZone: ", data.FromZone)
	omcplog.V(4).Info("ToZone: ", data.ToZone)
	omcplog.V(4).Info("ToClusterName: ", data.ToClusterName)
	omcplog.V(4).Info("ToNamespace: ", data.ToNamespace)
	omcplog.V(4).Info("ToPodName: ", data.ToPodName)

	omcplog.V(4).Info("GeoScore: ", tempGeoScore)
	omcplog.V(4).Info("ClusterPodResourceScore: ", tempClusterPodResourceScore)
	omcplog.V(4).Info("weight: ", weight)
	omcplog.V(4).Info("")

	return &protobuf.ResponseWeight{Weight: weight}, nil

}

// func (ae *AnalyticEngineStruct) SendRegionZoneInfo(ctx context.Context, data *protobuf.RegionZoneInfo) (*protobuf.ResponseWeight, error) {
// 	omcplog.V(4).Info("Func SendRegionZoneInfo Called")

// 	var tempGeoScore float64
// 	var tempClusterSVCResourceScore float64

// 	if data.FromRegion == data.ToRegion && data.FromZone == data.ToZone {
// 		if ae.GeoScore[0] == -1 {
// 			return nil, errors.New("GeoScore Initiailzing. Retry again")
// 		}
// 		tempGeoScore = ae.GeoScore[0]
// 	} else if data.FromRegion == data.ToRegion {
// 		if ae.GeoScore[1] == -1 {
// 			return nil, errors.New("GeoScore Initiailzing. Retry again")
// 		}
// 		tempGeoScore = ae.GeoScore[1]
// 	} else {
// 		if ae.GeoScore[2] == -1 {
// 			return nil, errors.New("GeoScore Initiailzing. Retry again")
// 		}
// 		tempGeoScore = ae.GeoScore[2]
// 	}

// 	_, exists := ae.ClusterSVCResourceScore[data.ToClusterName]
// 	if !exists {
// 		return nil, errors.New("ClusterSVCResourceScore Initiailzing. Retry again")
// 	}
// 	_, exists = ae.ClusterSVCResourceScore[data.ToClusterName][data.ToNamespace]
// 	if !exists {
// 		return nil, errors.New("ClusterSVCResourceScore Initiailzing. Retry again")
// 	}
// 	_, exists = ae.ClusterSVCResourceScore[data.ToClusterName][data.ToNamespace][data.ToServiceName]
// 	if !exists {
// 		return nil, errors.New("ClusterSVCResourceScore Initiailzing. Retry again")
// 	}

// 	tempClusterSVCResourceScore = ae.ClusterSVCResourceScore[data.ToClusterName][data.ToNamespace][data.ToServiceName]

// 	weight := float32(tempGeoScore + tempClusterSVCResourceScore)

// 	omcplog.V(4).Info("FromRegion: ", data.FromRegion)
// 	omcplog.V(4).Info("ToRegion: ", data.ToRegion)
// 	omcplog.V(4).Info("FromZone: ", data.FromZone)
// 	omcplog.V(4).Info("ToZone: ", data.ToZone)
// 	omcplog.V(4).Info("ToClusterName: ", data.ToClusterName)
// 	omcplog.V(4).Info("ToNamespace: ", data.ToNamespace)
// 	omcplog.V(4).Info("ToServiceName: ", data.ToServiceName)

// 	omcplog.V(4).Info("GeoScore: ", tempGeoScore)
// 	omcplog.V(4).Info("ClusterSVCResourceScore: ", tempClusterSVCResourceScore)
// 	omcplog.V(4).Info("weight: ", weight)
// 	omcplog.V(4).Info("")

// 	return &protobuf.ResponseWeight{Weight: weight}, nil

// }
func (ae *AnalyticEngineStruct) SendHASMaxAnalysis(ctx context.Context, data *protobuf.HASInfo) (*protobuf.ResponseHAS, error) {
	omcplog.V(4).Info("Func SendHASMaxAnalysis Called")

	timeStart_analysis := time.Now()
	filteringCluster := ae.SelectHPACluster(data)
	timeEnd_analysis := time.Since(timeStart_analysis)
	omcplog.V(2).Info("[1] SelectCandidateCluster \t", timeEnd_analysis)

	var result string
	result = ae.CompareHPAMaxInfo(filteringCluster, data)

	timeEnd_analysis4 := time.Since(timeStart_analysis)
	omcplog.V(2).Info("-----------------------------------------")
	omcplog.V(2).Info("==> Total Analysis time \t", timeEnd_analysis4)
	omcplog.V(2).Info("ResultCluster\t[", result, "]")

	omcplog.V(2).Info("*******  [End] HAS Rebalancing Analysis  ******* \n")

	return &protobuf.ResponseHAS{TargetCluster: result}, nil
}

func (ae *AnalyticEngineStruct) SendHASMinAnalysis(ctx context.Context, data *protobuf.HASInfo) (*protobuf.ResponseHAS, error) {
	omcplog.V(4).Info("Func SendHASMinAnalysis Called")

	timeStart_analysis := time.Now()
	filteringCluster := ae.SelectHPACluster(data)
	timeEnd_analysis := time.Since(timeStart_analysis)
	omcplog.V(2).Info("[1] SelectCandidateCluster \t", timeEnd_analysis)

	var result string
	result = ae.CompareHPAMinInfo(filteringCluster, data)

	timeEnd_analysis4 := time.Since(timeStart_analysis)
	omcplog.V(2).Info("-----------------------------------------")
	omcplog.V(2).Info("==> Total Analysis time \t", timeEnd_analysis4)
	omcplog.V(2).Info("ResultCluster\t[", result, "]")

	omcplog.V(2).Info("*******  [End] HAS Rebalancing Analysis  ******* \n")

	return &protobuf.ResponseHAS{TargetCluster: result}, nil
}

func (ae *AnalyticEngineStruct) SendNetworkAnalysis(ctx context.Context, data *protobuf.NodeInfo) (*protobuf.ResponseNetwork, error) {
	omcplog.V(4).Info("Func SendNetworkAnalysis Called")
	startTime := time.Now()

	// calculate difference between previous data and next data
	var diff_rx int64 = -1
	var diff_tx int64 = -1
	if _, ok := ae.NetworkInfos[data.ClusterName][data.NodeName]; ok {
		diff_rx = ae.NetworkInfos[data.ClusterName][data.NodeName].next_rx - ae.NetworkInfos[data.ClusterName][data.NodeName].prev_rx
		diff_tx = ae.NetworkInfos[data.ClusterName][data.NodeName].next_tx - ae.NetworkInfos[data.ClusterName][data.NodeName].prev_tx
	}

	elapsedTime := time.Since(startTime)
	omcplog.V(2).Info("%-30s [", elapsedTime, "]", "=> Total Anlysis time")
	omcplog.V(2).Info("***** [End] Network Analysis *****")

	return &protobuf.ResponseNetwork{RX: diff_rx, TX: diff_tx}, nil
}

var grpcServer *grpc.Server

func (ae *AnalyticEngineStruct) StartGRPC(GRPC_PORT string) {
	omcplog.V(4).Info("Func StartGRPC Called")
	log.Printf("Grpc Server Start at Port %s\n", GRPC_PORT)

	l, err := net.Listen("tcp", ":"+GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer = grpc.NewServer()

	protobuf.RegisterRequestAnalysisServer(grpcServer, ae)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("fail to serve: %v", err)
	}

}
func (ae *AnalyticEngineStruct) StopGRPC() {
	omcplog.V(4).Info("Func StopGRPC Called")
	grpcServer.Stop()

}

func (ae *AnalyticEngineStruct) SendCPAAnalysis(ctx context.Context, deploy *protobuf.CPADeployList) (*protobuf.ResponseCPADeployList, error) {
	omcplog.V(4).Info("Func SendCPAAnalysis Called")

	omcplog.V(4).Info("CPADeployInfo : ", deploy.CPADeployInfo)
	deployList := protobuf.ResponseCPADeployList{}

	for _, cDeploy := range deploy.CPADeployInfo {
		tmp_deployList := protobuf.ResponseCPADeployList{}
		tmp_cluster_list := []string{}

		for _, cluster := range cDeploy.Clusters {
			state, todo, check_rest := ae.AnalyzeCPADeployment(cDeploy, cluster)
			if todo == "Scale-out" || todo == "Scale-in" {
				d := &protobuf.ResponseCPADeploy{
					Name:          cDeploy.Name,
					Namespace:     cDeploy.Namespace,
					CPAName:       cDeploy.CPAName,
					PodState:      state,
					Action:        todo,
					TargetCluster: cluster,
				}
				tmp_deployList.ResponseCPADeploy = append(tmp_deployList.ResponseCPADeploy, d)
			}
			if check_rest == 1 {
				tmp_cluster_list = append(tmp_cluster_list, cluster)
			}
		}

		for i, _ := range tmp_deployList.ResponseCPADeploy {
			tmp_deployList.ResponseCPADeploy[i].RestCluster = tmp_cluster_list
		}

		deployList.ResponseCPADeploy = append(deployList.ResponseCPADeploy, tmp_deployList.ResponseCPADeploy...)
	}
	omcplog.V(4).Info("*** deployList : ", deployList)
	return &deployList, nil
}

func (ae *AnalyticEngineStruct) AnalyzeCPADeployment(cDeploy *protobuf.CPADeployInfo, cluster string) (string, string, int) {
	omcplog.V(4).Info("Func AnalyzeCPADeployment Called")
	//** request 없으면 CPA 불가

	//cpu, memory, fs 문제면 cluster 내
	//cpuusage/cpurequest > 60 or memusage/memrequest > 60 이면 Warning
	//노드 용량이 넉넉하면 'scale-out'
	//cpuusage/cpurequest < 1 or memusage/memrequest < 1 이면 Warning
	//'scale-in'

	//네트워크 문제면 (서비스 문제) cluster 간
	//networkLatency > x 인 경우 Warning
	//클러스터 ResourceScore를 비교해서 점수가 낮은 클러스터에 'scale-out' (낮을수록 좋음)

	cpuRequestInt64 := cDeploy.CpuRequest
	memRequestInt64 := cDeploy.MemRequest

	omcplog.V(3).Info("Start Analysis ...")
	omcplog.V(3).Info("podNum : ", strconv.FormatInt(int64(cDeploy.ReplicasNum), 10))
	result := ae.Influx.GetCPAMetricsData(cluster, cDeploy.Namespace, cDeploy.Name, strconv.FormatInt(int64(cDeploy.ReplicasNum), 10))

	var cpuUsage float64
	var cputotal float64
	var memUsage float64
	var memtotal float64

	cputotal = 0
	memtotal = 0

	if result != nil {
		for _, row := range result[0].Series[0].Values {
			cpu := row[1]
			cpuString := fmt.Sprintf("%v", cpu)
			cpuString = cpuString[:len(cpuString)-1]
			cpuFloat64, _ := strconv.ParseFloat(cpuString, 64)
			cpuFloat64 = cpuFloat64 * 1e-06
			cputotal += cpuFloat64

			mem := row[2]
			memString := fmt.Sprintf("%v", mem)
			memString = memString[:len(memString)-2]
			memFloat64, _ := strconv.ParseFloat(memString, 64)
			memFloat64 = memFloat64 * 1e+6
			memtotal += memFloat64
		}
	} else {
		omcplog.V(1).Info("Empty")
	}

	num := float64(cDeploy.ReplicasNum)

	//cpu,memory

	cpuUsage = cputotal / num
	memUsage = memtotal / num
	omcplog.V(1).Info("*** ", cluster, " ***")
	omcplog.V(1).Info("========================================================")
	omcplog.V(1).Info("[", cDeploy.Name, "] CPU 사용률 ", cpuUsage/float64(cpuRequestInt64)*100, "%")
	omcplog.V(1).Info("[", cDeploy.Name, "] MEM 사용률 ", memUsage/float64(memRequestInt64)*100, "%")
	omcplog.V(2).Info("--------------------------------------------------------")
	omcplog.V(2).Info("ClusterResourceUsage[CPU] : ", ae.ClusterResourceUsage[cluster]["cpu"])
	omcplog.V(2).Info("ClusterResourceUsage[MEM] : ", ae.ClusterResourceUsage[cluster]["memory"])
	omcplog.V(1).Info("========================================================")
	omcplog.V(1).Info("")

	if cpuUsage/float64(cpuRequestInt64)*100 > 60 {
		if ae.ClusterResourceUsage[cluster]["cpu"] < 80 {
			omcplog.V(1).Info("CPU Warning! -> Scale-out")
			return "Warning-cpu", "Scale-out", 0
		} else {
			omcplog.V(1).Info("CPU Warning! -> Can't Scale-out (No Capacity)")
		}
	} else if memUsage/float64(memRequestInt64)*100 > 60 {
		if ae.ClusterResourceUsage[cluster]["memory"] < 80 {
			omcplog.V(1).Info("Memory Warning! -> Scale-out")
			return "Warning-memory", "Scale-out", 0
		} else {
			omcplog.V(1).Info("Memory Warning! -> Can't Scale-out (No Capacity)")
		}
	} else if cpuUsage/float64(cpuRequestInt64)*100 < 5 && memUsage/float64(memRequestInt64)*100 < 5 {
		omcplog.V(1).Info("CPU/Memory Warning! -> Scale-in")
		return "Warning-cpu/memory", "Scale-in", 0
	}

	if cpuUsage/float64(cpuRequestInt64)*100 < 20 && memUsage/float64(memRequestInt64)*100 < 20 {
		return "", "", 1
	}

	/*
		//cpu
		cpuUsage = cputotal / num
		memUsage = memtotal / num

		fmt.Println("[", cDeploy.Name, "] CPU 사용률 ", cpuUsage/float64(cpuRequestInt64)*100, "%")
		fmt.Println("[", cDeploy.Name, "] MEM 사용률 ", memUsage/float64(memRequestInt64)*100, "%")

		if cpuUsage/float64(cpuRequestInt64)*100 > 80 {
			if ae.ClusterResourceUsage[cluster]["cpu"] < 80 {
				fmt.Println("CPU Warning! -> Scale-out")
				return cluster, "Warning-cpu", "Scale-out"
			} else {
				fmt.Println("CPU Warning! -> Can't Scale-out (No Capacity)")
			}
		} else if memUsage/float64(memRequestInt64)*100 > 80 {
			if ae.ClusterResourceUsage[cluster]["memory"] < 80 {
				fmt.Println("Memory Warning! -> Scale-out")
				return cluster, "Warning-memory", "Scale-out"
			} else {
				fmt.Println("Memory Warning! -> Can't Scale-out (No Capacity)")
			}
		} else if cpuUsage/float64(cpuRequestInt64)*100 < 0.1 && memUsage/float64(memRequestInt64)*100 < 1 {
			fmt.Println("CPU/Memory Warning! -> Scale-in")
			return cluster, "Warning-cpu/memory", "Scale-in"
		}


			//network
			netLatency := 0  //influxdb

			if netLatency > 1000 {
				var minScore float64
				minScore = 1000000
				resultCluster := ""
				for _, clustername := range cDeploy.Clusters {
					s := ae.ResourceScore[clustername]
					if s < minScore {
						minScore = s
						resultCluster = clustername
					}
				}
			} else if memUsage/float64(memRequestInt64)*100 < 1 {
				fmt.Println("Memory Warning! -> Scale-in")
				return cluster, "Warning-memory", "Scale-in"
			}
	*/

	return "", "", 0

}
