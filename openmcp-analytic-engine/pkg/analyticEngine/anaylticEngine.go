package analyticEngine

import (
	"context"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"net"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-analytic-engine/pkg/Geo"
	"openmcp/openmcp/openmcp-analytic-engine/pkg/influx"
	"openmcp/openmcp/openmcp-analytic-engine/pkg/protobuf"
	"openmcp/openmcp/util/clusterManager"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	"sort"
	"strconv"
	"time"
)

type AnalyticEngineStruct struct {
	Influx        influx.Influx
	MetricsWeight map[string]float64
	ResourceScore map[string]float64
	ClusterGeo    map[string]map[string]string
	NetworkInfos  map[string]map[string]*NetworkInfo
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
	ae := &AnalyticEngineStruct{}
	ae.Influx = *influx.NewInflux(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)
	ae.ResourceScore = make(map[string]float64)
	return ae
}

func (ae *AnalyticEngineStruct) CalcResourceScore() {
	omcplog.V(4).Info("Func CalcResourceScore Called")
	cm := clusterManager.NewClusterManager()
	ae.MetricsWeight = make(map[string]float64)
	ae.ClusterGeo = map[string]map[string]string{}
	ae.NetworkInfos = make(map[string]map[string]*NetworkInfo)

	//정책 엔진 - 메트릭 가중치 읽어오기----------------------------
	openmcpPolicyInstance, target_cluster_policy_err := cm.Crd_client.OpenMCPPolicy("openmcp").Get("analytic-metrics-weight", metav1.GetOptions{})

	if target_cluster_policy_err != nil {
		omcplog.V(0).Info(target_cluster_policy_err)
	} else {
		a := openmcpPolicyInstance.Spec.Template.Spec.Policies
		for _, b := range a {
			value, _ := strconv.ParseFloat(b.Value[0], 64)
			ae.MetricsWeight[b.Type] = value
		}
		omcplog.V(3).Info("metricsWeight : ", ae.MetricsWeight)
	}
	//-------------------------------------------------------
	for {
		omcplog.V(2).Info("Cluster들의 Metric Score를 갱신합니다.")
		for _, cluster := range cm.Cluster_list.Items {
			ae.ResourceScore[cluster.Name] = ae.UpdateScore(cluster.Name)
			config, _ := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
			clientset, _ := kubernetes.NewForConfig(config)
			nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
			if err != nil {
				fmt.Println(err)
			}

			if len(nodes.Items) != 0 {
				node := nodes.Items[0]

				//label로부터 zone, region 추출
				ae.ClusterGeo[cluster.Name] = map[string]string{}
				ae.ClusterGeo[cluster.Name]["Country"] = node.Labels["failure-domain.beta.kubernetes.io/zone"]
				ae.ClusterGeo[cluster.Name]["Continent"] = node.Labels["failure-domain.beta.kubernetes.io/region"]
			}

			// Update Network Data from InfluxDB
			ae.UpdateNetworkData(cluster.Name, nodes)
		}
		time.Sleep(5 * time.Second)
	}
}

// Update Network data from InfluxDB
func (ae *AnalyticEngineStruct) UpdateNetworkData(clusterName string, nodeList *corev1.NodeList) {
	omcplog.V(4).Info("Func UpdateNetworkData Called")
	// Initialize cluster's network data
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
}

func (ae *AnalyticEngineStruct) UpdateScore(clusterName string) float64 {
	omcplog.V(4).Info("Func UpdateScore Called")
	var score float64 = 0
	result := ae.Influx.GetClusterMetricsData(clusterName)

	MetricsMap := make(map[string]float64)
	prevMetricsMap := make(map[string]float64)
	var totalCpuCore int64 = 0

	cm := clusterManager.NewClusterManager()

	for _, ser := range result[0].Series {

		nodeCapacity := &corev1.Node{}
		err := cm.Cluster_genClients[ser.Tags["cluster"]].Get(context.TODO(), nodeCapacity, "", ser.Tags["node"])
		if err != nil {
			omcplog.V(0).Info("nodelist err : ", err)
		} else {
			omcplog.V(2).Info("[CPU Capacity] ", ser.Tags["cluster"], "/", ser.Tags["node"], "/", nodeCapacity.Status.Capacity.Cpu().Value())
		}

		totalCpuCore = totalCpuCore + nodeCapacity.Status.Capacity.Cpu().Value()
		//fmt.Println(totalCpuCore)
		//fmt.Println(ser)
		for c, colName := range ser.Columns {
			for r, _ := range ser.Values {

				Strval := fmt.Sprintf("%v", ser.Values[r][c])
				QuanVal, _ := resource.ParseQuantity(Strval)

				if r == 0 {
					if _, ok := MetricsMap[colName]; ok {
						MetricsMap[colName] = MetricsMap[colName] + float64(QuanVal.Value())
					} else {
						MetricsMap[colName] = float64(QuanVal.Value())
					}
				} else if r == 1 {
					if _, ok := prevMetricsMap[colName]; ok {
						prevMetricsMap[colName] = prevMetricsMap[colName] + float64(QuanVal.Value())
					} else {
						prevMetricsMap[colName] = float64(QuanVal.Value())
					}
				}

			}
		}
		//fmt.Println(totalCpuCore)
	}

	cpuScore := (float64(totalCpuCore) - MetricsMap["CPUUsageNanoCores"]) / float64(totalCpuCore) * 100 // 확인필요
	memScore := MetricsMap["MemoryAvailableBytes"] / (MetricsMap["MemoryUsageBytes"] + MetricsMap["MemoryAvailableBytes"]) * 100
	//netScore := (MetricsMap["NetworkRxBytes"] - prevMetricsMap["NetworkRxBytes"]) + (MetricsMap["NetworkTxBytes"] - prevMetricsMap["NetworkTxBytes"]) / float64(totalNet) * 100// 확인필요
	diskScore := MetricsMap["FsAvailableBytes"] / MetricsMap["FsCapacityBytes"] * 100

	score = cpuScore*ae.MetricsWeight["CPU"] + memScore*ae.MetricsWeight["Memory"] + diskScore*ae.MetricsWeight["FS"]

	omcplog.V(2).Info("--------------------------------------------------------------------------------------------")
	omcplog.V(2).Info("totalScore : ", score)
	omcplog.V(2).Info("--------------------------------------------------------------------------------------------")

	return score
}

func (ae *AnalyticEngineStruct) SendLBAnalysis(ctx context.Context, in *protobuf.LBInfo) (*protobuf.ResponseLB, error) {
	omcplog.V(4).Info("Func SendLBAnalysis Called")

	//clusterNameList := in.ClusterNameList
	clusterScoreMap := make(map[string]float64)

	//for _, clusterName := range clusterNameList {
	//	clusterScoreMap[clusterName] = ae.ResourceScore[clusterName]
	//}
	//
	//omcplog.V(5).Info(clusterScoreMap)
	//
	//omcplog.V(2).Info("Geo Score")
	//clientIP := in.ClientIP
	//country := ae.getCountry(clientIP)
	//continent := ae.getContinent(country)
	//
	//score := ae.geoScore(clusterNameList, country, continent)
	//
	//for _, clusterName := range clusterNameList {
	//	omcplog.V(5).Info(clusterScoreMap[clusterName])
	//	omcplog.V(5).Info(score[clusterName])
	//	clusterScoreMap[clusterName] = clusterScoreMap[clusterName] + score[clusterName]
	//}
	clusterScoreMap = ae.ResourceScore

	omcplog.V(2).Info("LB Response")
	omcplog.V(3).Info(clusterScoreMap)
	return &protobuf.ResponseLB{ScoreMap: clusterScoreMap}, nil
}

func (ae *AnalyticEngineStruct) SelectHPACluster(data *protobuf.HASInfo) []string {
	omcplog.V(4).Info("Func SelectHPACluster Called")

	scoreMap := map[float64]string{}
	score := []float64{}
	omcplog.V(5).Info(ae.ResourceScore)
	for key, value := range ae.ResourceScore {
		if key != data.ClusterName && value > 0 {
			scoreMap[value] = key
			score = append(score, value)
		}
	}
	sort.Float64s(score)

	filteringCluster := []string{}

	for i := 0; i < len(score); i++ {
		filteringCluster = append(filteringCluster, scoreMap[score[i]])
	}

	/*filteringCluster = append(filteringCluster, "cluster2")
	filteringCluster = append(filteringCluster, "cluster3")
	fmt.Println(filteringCluster)*/

	return filteringCluster
}

func (ae *AnalyticEngineStruct) CompareHPAMaxInfo(clusterList []string, data *protobuf.HASInfo) string {
	omcplog.V(4).Info("Func CompareHPAMaxInfo Called")
	replicasGap := map[string]int32{}
	rebalancingCount := map[string]int32{}

	for _, cluster := range clusterList {
		omcplog.V(3).Info(cluster, " hpa MaxReplicas : ", data.HPAMinORMaxReplicas[cluster], " / CurrentReplicas : ", data.HPACurrentReplicas[cluster])
		//calc := hpaInstance.Spec.MaxReplicas - hpaInstance.Status.CurrentReplicas
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

func (ae *AnalyticEngineStruct) SendHASMaxAnalysis(ctx context.Context, data *protobuf.HASInfo) (*protobuf.ResponseHAS, error) {
	omcplog.V(4).Info("Func SendHASMaxAnalysis Called")

	//fmt.Println(data)
	timeStart_analysis := time.Now()
	filteringCluster := ae.SelectHPACluster(data)
	timeEnd_analysis := time.Since(timeStart_analysis)
	omcplog.V(2).Info("[1] SelectCandidateCluster \t", timeEnd_analysis)

	var result string
	/*	if len(filteringCluster) == 1 {
		result = filteringCluster[0]

	}else {*/
	result = ae.CompareHPAMaxInfo(filteringCluster, data)
	//	}

	timeEnd_analysis4 := time.Since(timeStart_analysis)
	omcplog.V(2).Info("-----------------------------------------")
	omcplog.V(2).Info("==> Total Analysis time \t", timeEnd_analysis4)
	omcplog.V(2).Info("ResultCluster\t[", result, "]")

	omcplog.V(2).Info("*******  [End] HAS Rebalancing Analysis  ******* \n")
	//fmt.Println("---------HAS Response End---------")

	return &protobuf.ResponseHAS{TargetCluster: result}, nil
}

func (ae *AnalyticEngineStruct) SendHASMinAnalysis(ctx context.Context, data *protobuf.HASInfo) (*protobuf.ResponseHAS, error) {
	omcplog.V(4).Info("Func SendHASMinAnalysis Called")
	//fmt.Println("\n******* [Start] HAS Rebalancing Analysis *******")

	timeStart_analysis := time.Now()
	filteringCluster := ae.SelectHPACluster(data)
	timeEnd_analysis := time.Since(timeStart_analysis)
	omcplog.V(2).Info("[1] SelectCandidateCluster \t", timeEnd_analysis)

	var result string
	//if len(filteringCluster) == 1 {
	//	result = filteringCluster[0]
	//}else {
	//timeStart_analysis2 := time.Now()
	result = ae.CompareHPAMinInfo(filteringCluster, data)
	//timeEnd_analysis2 := time.Since(timeStart_analysis2)
	//fmt.Println("[2] CompareHPAMinInfo() \t", timeEnd_analysis2)
	//}

	timeEnd_analysis4 := time.Since(timeStart_analysis)
	omcplog.V(2).Info("-----------------------------------------")
	omcplog.V(2).Info("==> Total Analysis time \t", timeEnd_analysis4)
	omcplog.V(2).Info("ResultCluster\t[", result, "]")

	omcplog.V(2).Info("*******  [End] HAS Rebalancing Analysis  ******* \n")

	return &protobuf.ResponseHAS{TargetCluster: result}, nil
}

func (ae *AnalyticEngineStruct) SendNetworkAnalysis(ctx context.Context, data *protobuf.NodeInfo) (*protobuf.ReponseNetwork, error) {
	omcplog.V(4).Info("Func SendNetworkAnalysis Called")
	//klog.Info("***** [Start] Network Analysis *****")
	startTime := time.Now()

	// calculate difference between previous data and next data
	diff_rx := ae.NetworkInfos[data.ClusterName][data.NodeName].next_rx - ae.NetworkInfos[data.ClusterName][data.NodeName].prev_rx
	diff_tx := ae.NetworkInfos[data.ClusterName][data.NodeName].next_tx - ae.NetworkInfos[data.ClusterName][data.NodeName].prev_tx

	elapsedTime := time.Since(startTime)
	omcplog.V(2).Info("%-30s [",elapsedTime,"]", "=> Total Anlysis time")
	omcplog.V(2).Info("***** [End] Network Analysis *****")

	return &protobuf.ReponseNetwork{RX: diff_rx, TX: diff_tx}, nil
}

func (ae *AnalyticEngineStruct) StartGRPC(GRPC_PORT string) {
	omcplog.V(4).Info("Func StartGRPC Called")
	log.Printf("Grpc Server Start at Port %s\n", GRPC_PORT)

	l, err := net.Listen("tcp", ":"+GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	protobuf.RegisterRequestAnalysisServer(grpcServer, ae)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("fail to serve: %v", err)
	}

}

//LoadBalancing
func (ae *AnalyticEngineStruct) geoScore(clusters []string, clientCountry, clientContinent string) map[string]float64 {
	omcplog.V(4).Info("Func geoScore Called")

//	fmt.Println("*****Geo Score*****")

	midScore := 100.0
	policy := ae.MetricsWeight["GeoRate"] * 100.0
	//policy := 70.0
	omcplog.V(5).Info(ae.ClusterGeo)

	score := map[string]float64{}
	for _, cluster := range clusters {
		clustercountry := ae.ClusterGeo[cluster]["Country"]
		clustercontinent := ae.ClusterGeo[cluster]["Continent"]

		if clientCountry == clustercountry {
			score[cluster] = midScore + (midScore * policy / 100.0)
		} else if clientContinent == clustercontinent {
			score[cluster] = midScore
		} else {
			score[cluster] = midScore - (midScore * policy / 100.0)
		}
	}
	omcplog.V(5).Info(score)
	return score
}

func (ae *AnalyticEngineStruct) getCountry(clientip string) string {
	omcplog.V(4).Info("Func getCountry Called")
	//fmt.Println("*****Extract Country*****")
	omcplog.V(2).Info("'/root/GeoLite2-City.mmdb' Open")
	db, err := geoip2.Open("/root/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//ip := net.ParseIP(clientip)
	ip := net.ParseIP("14.102.132.0")

	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	return record.Country.IsoCode
}

func (ae *AnalyticEngineStruct) getContinent(country string) string {
	return Geo.Geo[country]
}
