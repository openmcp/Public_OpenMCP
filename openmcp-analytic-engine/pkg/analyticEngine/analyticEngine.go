package analyticEngine

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"net"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-analytic-engine/pkg/influx"
	"openmcp/openmcp/openmcp-analytic-engine/pkg/protobuf"
	"openmcp/openmcp/util/clusterManager"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	"sort"
	"strconv"
	"time"
)

type AnalyticEngineStruct struct {
	Influx               influx.Influx
	MetricsWeight        map[string]float64
	ResourceScore        map[string]float64
	ClusterResourceUsage map[string]map[string]float64
	ClusterGeo           map[string]map[string]string
	NetworkInfos         map[string]map[string]*NetworkInfo
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
	ae.ClusterResourceUsage = make(map[string]map[string]float64)
	return ae
}

func (ae *AnalyticEngineStruct) CalcResourceScore(cm *clusterManager.ClusterManager, quit, quitok chan bool) {
	omcplog.V(4).Info("Func CalcResourceScore Called")
	//cm := clusterManager.NewClusterManager()
	ae.MetricsWeight = make(map[string]float64)
	ae.ClusterGeo = map[string]map[string]string{}
	ae.NetworkInfos = make(map[string]map[string]*NetworkInfo)

	//Get metric-weight Policy----------------------------
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
	for {
		select {
		case <-quit:
			omcplog.V(2).Info("CalcResourceScore Quit")
			quitok <- true
			return
		default:
			omcplog.V(2).Info("Cluster Metric Score Refresh")
			for _, cluster := range cm.Cluster_list.Items {
				ae.ResourceScore[cluster.Name], ae.ClusterResourceUsage[cluster.Name] = ae.UpdateScore(cluster.Name, cm)
				config, err := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
				if err != nil {
					omcplog.V(0).Info(err)
				}
				clientset, err := kubernetes.NewForConfig(config)
				if err != nil {
					omcplog.V(0).Info(err)
				}
				nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
				if err != nil {
					omcplog.V(0).Info(err)
				}

				if len(nodes.Items) != 0 {
					node := nodes.Items[0]

					//Extract zone, region from Label
					ae.ClusterGeo[cluster.Name] = map[string]string{}
					ae.ClusterGeo[cluster.Name]["Country"] = node.Labels["failure-domain.beta.kubernetes.io/zone"]
					ae.ClusterGeo[cluster.Name]["Continent"] = node.Labels["failure-domain.beta.kubernetes.io/region"]
				}

				// Update Network Data from InfluxDB
				ae.UpdateNetworkData(cluster.Name, nodes)
			}
			time.Sleep(7 * time.Second)
		}

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

func (ae *AnalyticEngineStruct) UpdateScore(clusterName string, cm *clusterManager.ClusterManager) (float64, map[string]float64) {
	omcplog.V(4).Info("Func UpdateScore Called")
	var score float64 = 0
	scoreMap := make(map[string]float64)

	result := ae.Influx.GetClusterMetricsData(clusterName)

	MetricsMap := make(map[string]float64)
	prevMetricsMap := make(map[string]float64)
	var totalCpuCore int64 = 0

	//cm := clusterManager.NewClusterManager()
	if len(result) == 0 {
		return 0, nil
	}
	time_t := ""
	for _, ser := range result[0].Series {
		nodeCapacity := &corev1.Node{}
		err := cm.Cluster_genClients[ser.Tags["cluster"]].Get(context.TODO(), nodeCapacity, "", ser.Tags["node"])
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
							MetricsMap[colName] = MetricsMap[colName] + float64(QuanVal.Value())
						} else {
							MetricsMap[colName] = float64(QuanVal.Value())
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

	}

	var netScore float64

	cpuScore := MetricsMap["CPUUsageNanoCores"] / float64(totalCpuCore) * 100
	memScore := (MetricsMap["MemoryUsageBytes"] / (MetricsMap["MemoryUsageBytes"] + MetricsMap["MemoryAvailableBytes"])) * 100
	netScore = ((MetricsMap["NetworkRxBytes"] - prevMetricsMap["NetworkRxBytes"]) + (MetricsMap["NetworkTxBytes"] - prevMetricsMap["NetworkTxBytes"])) / 1000000
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

	if score > 0 && clusterName == "eks-cluster1" {
		omcplog.V(2).Info("--------------------------------------------------------")
		omcplog.V(2).Info(" rx -> ", MetricsMap["NetworkRxBytes"]-prevMetricsMap["NetworkRxBytes"])
		omcplog.V(2).Info(" tx -> ", MetricsMap["NetworkTxBytes"]-prevMetricsMap["NetworkTxBytes"])
		omcplog.V(2).Info(" cpuscore : ", cpuScore, MetricsMap["CPUUsageNanoCores"], float64(totalCpuCore))
		omcplog.V(2).Info(" memScore : ", memScore)
		omcplog.V(2).Info(" netScore : ", netScore)
		omcplog.V(2).Info(" diskScore : ", diskScore)
		omcplog.V(2).Info(" latencyScore : ", latencyScore)
		omcplog.V(2).Info("\"", clusterName, "\" totalScore : ", score)
		omcplog.V(2).Info("--------------------------------------------------------")
	}

	return score, scoreMap
}

func (ae *AnalyticEngineStruct) SendLBAnalysis(ctx context.Context, in *protobuf.LBInfo) (*protobuf.ResponseLB, error) {
	omcplog.V(4).Info("Func SendLBAnalysis Called")

	clusterScoreMap := make(map[string]float64)

	clusterScoreMap = ae.ResourceScore

	omcplog.V(2).Info("LB Response")
	omcplog.V(3).Info(clusterScoreMap)
	return &protobuf.ResponseLB{ScoreMap: clusterScoreMap}, nil
}

func (ae *AnalyticEngineStruct) SelectHPACluster(data *protobuf.HASInfo) []string {
	omcplog.V(4).Info("Func SelectHPACluster Called")

	scoreMap := map[float64]string{}
	scoreT := []float64{}

	for key, value := range ae.ResourceScore {
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

	fmt.Println(deploy.CPADeployInfo)
	deployList := protobuf.ResponseCPADeployList{}

	for _, cDeploy := range deploy.CPADeployInfo {
		target, state, todo := ae.AnalyzeCPADeployment(cDeploy)
		if todo == "Scale-out" || todo == "Scale-in" {
			d := &protobuf.ResponseCPADeploy{
				Name:          cDeploy.Name,
				Namespace:     cDeploy.Namespace,
				CPAName:       cDeploy.CPAName,
				PodState:      state,
				Action:        todo,
				TargetCluster: target,
			}
			deployList.ResponseCPADeploy = append(deployList.ResponseCPADeploy, d)
		}
	}

	return &deployList, nil
}

func (ae *AnalyticEngineStruct) AnalyzeCPADeployment(cDeploy *protobuf.CPADeployInfo) (string, string, string) {
	omcplog.V(4).Info("Func AnalyzeCPADeployment Called")
	//** request 없으면 CPA 불가

	//cpu, memory, fs 문제면 cluster 내
	//cpuusage/cpurequest > 80 or memusage/memrequest > 80 이면 Warning
	//노드 용량이 넉넉하면 'scale-out'
	//cpuusage/cpurequest < 10 or memusage/memrequest < 10 이면 Warning
	//'scale-in'

	//네트워크 문제면 (서비스 문제) cluster 간
	//networkLatency > x 인 경우 Warning
	//클러스터 ResourceScore를 비교해서 점수가 낮은 클러스터에 'scale-out' (낮을수록 좋음)

	cpuRequestInt64 := cDeploy.CpuRequest
	memRequestInt64 := cDeploy.MemRequest

	for _, cluster := range cDeploy.Clusters {
		fmt.Println("Start Analysis ...")
		fmt.Println("podNum : ", strconv.FormatInt(int64(cDeploy.ReplicasNum), 10))
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
			fmt.Println("Empty")
		}

		num := float64(cDeploy.ReplicasNum)

		//cpu
		cpuUsage = cputotal / num
		fmt.Println("[", cDeploy.Name, "] CPU 사용률 ", cpuUsage/float64(cpuRequestInt64)*100, "%")

		if cpuUsage/float64(cpuRequestInt64)*100 > 80 {
			if ae.ClusterResourceUsage[cluster]["cpu"] < 80 {
				fmt.Println("CPU Warning! -> Scale-out")
				return cluster, "Warning-cpu", "Scale-out"
			} else {
				fmt.Println("CPU Warning! -> Can't Scale-out (No Capacity)")
			}
		} else if cpuUsage/float64(cpuRequestInt64)*100 < 0.1 {
			fmt.Println("CPU Warning! -> Scale-in")
			return cluster, "Warning-cpu", "Scale-in"
		}

		//memory
		memUsage = memtotal / num
		fmt.Println("[", cDeploy.Name, "] MEM 사용률 ", memUsage/float64(memRequestInt64)*100, "%")

		if memUsage/float64(memRequestInt64)*100 > 80 {
			if ae.ClusterResourceUsage[cluster]["memory"] < 80 {
				fmt.Println("Memory Warning! -> Scale-out")
				return cluster, "Warning-memory", "Scale-out"
			} else {
				fmt.Println("Memory Warning! -> Can't Scale-out (No Capacity)")
			}
		} else if memUsage/float64(memRequestInt64)*100 < 1 {
			fmt.Println("Memory Warning! -> Scale-in")
			return cluster, "Warning-memory", "Scale-in"
		}

		/*
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
				fmt.Println("Network Warning! -> Scale-out")
				return resultCluster, "Warning-network", "scale-out"
			}
		*/

	}

	return "", "", ""

}
