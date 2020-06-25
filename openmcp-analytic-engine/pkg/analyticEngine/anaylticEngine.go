package analyticEngine

import (
	"github.com/oschwald/geoip2-golang"
	hpav2beta1 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"

	//"admiralty.io/multicluster-controller/pkg/cluster"
	//"k8s.io/apimachinery/pkg/types"
	"fmt"
	//"k8s.io/client-go/rest"

	//	appsv1 "k8s.io/api/apps/v1"
	//hpav2beta1 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/api/resource"
	//"k8s.io/apimachinery/pkg/types"

	//"k8s.io/apimachinery/pkg/types"
	"openmcp/openmcp/openmcp-analytic-engine/pkg/protobuf"
	"sort"
	"time"

	"context"
	//"github.com/influxdata/influxdb/client/v2"
	//"github.com/influxdata/influxdb/models"
	"google.golang.org/grpc"
	"log"
	//"math"
	"net"
	"openmcp/openmcp/openmcp-analytic-engine/pkg/Geo"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/openmcp-analytic-engine/pkg/influx"

	//fedapis "sigs.k8s.io/kubefed/pkg/apis"
	//ketiapis "resource-controller/apis"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/kubefed/pkg/controller/util"
)

type AnalyticEngineStruct struct {
	Influx        influx.Influx
	MetricsWeight map[string]float64
	ResourceScore map[string]float64
	ClusterGeo    map[string]map[string]string
}

func NewAnalyticEngine(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD string) *AnalyticEngineStruct {
	ae := &AnalyticEngineStruct{}
	ae.Influx = *influx.NewInflux(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)
	ae.ResourceScore = make(map[string]float64)
	return ae
}

func (ae *AnalyticEngineStruct) CalcResourceScore() {
	cm := clusterManager.NewClusterManager()
	ae.MetricsWeight = make(map[string]float64)
	ae.ClusterGeo = map[string]map[string]string{}
	//정책 엔진 - 메트릭 가중치 읽어오기----------------------------
	openmcpPolicyInstance, target_cluster_policy_err := cm.Crd_client.OpenMCPPolicyEngine("openmcp").Get("analytic-metrics-weight", metav1.GetOptions{})

	if target_cluster_policy_err != nil {
		fmt.Println(target_cluster_policy_err)
	} else {
		a := openmcpPolicyInstance.Spec.Template.Spec.Policies
		for _, b := range a {
			value, _ := strconv.ParseFloat(b.Value[0], 64)
			ae.MetricsWeight[b.Type] = value
		}
		fmt.Println("metricsWeight : ", ae.MetricsWeight)
	}
	//-------------------------------------------------------
	for {
		fmt.Println("Cluster들의 Metric Score를 갱신합니다.")
		for _, cluster := range cm.Cluster_list.Items {
			ae.ResourceScore[cluster.Name] = ae.UpdateScore(cluster.Name)
			config, _ := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
			clientset, _ := kubernetes.NewForConfig(config)
			nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
			if err != nil {
				fmt.Println(err)
			}
			node := nodes.Items[0]
			//label로부터 zone, region 추출
			ae.ClusterGeo[cluster.Name] = map[string]string{}
			ae.ClusterGeo[cluster.Name]["Country"] = node.Labels["failure-domain.beta.kubernetes.io/zone"]
			ae.ClusterGeo[cluster.Name]["Continent"] = node.Labels["failure-domain.beta.kubernetes.io/region"]
		}
		time.Sleep(5 * time.Second)
	}
}

func (ae *AnalyticEngineStruct) UpdateScore(clusterName string) float64 {
	var score float64 = 0
	result := ae.Influx.GetClusterMetricsData(clusterName)

	MetricsMap := make(map[string]float64)
	prevMetricsMap := make(map[string]float64)
	var totalCpuCore int64 = 0

	cm := clusterManager.NewClusterManager()

	for _, ser := range result[0].Series {

		nodeCapacity := &corev1.Node{}
		//err := cm.Cluster_genClients["cluster4"].List(context.TODO(), nodeCapacity, "")
		err := cm.Cluster_genClients[ser.Tags["cluster"]].Get(context.TODO(), nodeCapacity, "", ser.Tags["node"])
		if err != nil {
			fmt.Println("nodelist err : ", err)
		} else {
			//	fmt.Println(nodeCapacity)
			//	fmt.Println(nodeCapacity.Items[0].Status.Capacity.Cpu().MilliValue())
			fmt.Println("[CPU Capacity] ", ser.Tags["cluster"], "/", ser.Tags["node"], "/", nodeCapacity.Status.Capacity.Cpu().Value())
			//fmt.Println(ser.Tags["node"])
		}

		totalCpuCore = totalCpuCore + nodeCapacity.Status.Capacity.Cpu().Value()

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
	}

	cpuScore := (float64(totalCpuCore) - MetricsMap["CPUUsageNanoCores"]) / float64(totalCpuCore) * 100 // 확인필요
	memScore := MetricsMap["MemoryAvailableBytes"] / (MetricsMap["MemoryUsageBytes"] + MetricsMap["MemoryAvailableBytes"]) * 100
	//netScore := (MetricsMap["NetworkRxBytes"] - prevMetricsMap["NetworkRxBytes"]) + (MetricsMap["NetworkTxBytes"] - prevMetricsMap["NetworkTxBytes"]) / float64(totalNet) * 100// 확인필요
	diskScore := MetricsMap["FsAvailableBytes"] / MetricsMap["FsCapacityBytes"] * 100

	score = cpuScore*ae.MetricsWeight["CPU"] + memScore*ae.MetricsWeight["Memory"] + diskScore*ae.MetricsWeight["FS"]
	/*
		cpuUsed := resource.NewQuantity(int64(totalCpuCore) - int64(MetricsMap["CPUUsageNanoCores"]),resource.BinarySI).String()
		cpuTotal := resource.NewQuantity(int64(totalCpuCore), resource.BinarySI).String()

		memUsed := resource.NewQuantity(int64(MetricsMap["MemoryAvailableBytes"]),resource.BinarySI).String()
		memTotal := resource.NewQuantity(int64(MetricsMap["MemoryUsageBytes"] + MetricsMap["MemoryAvailableBytes"]), resource.BinarySI).String()

		//netUsed := resource.NewQuantity(int64(MetricsMap["NetworkRxBytes"]) - int64(prevMetricsMap["NetworkRxBytes"]) +  int64(MetricsMap["NetworkTxBytes"]) - int64(prevMetricsMap["NetworkTxBytes"]) ,resource.BinarySI).String()
		//netTotal := resource.NewQuantity(int64(totalNet), resource.BinarySI).String()

		diskUsed := resource.NewQuantity(int64(MetricsMap["FsAvailableBytes"]),resource.BinarySI).String()
		diskTotal := resource.NewQuantity(int64(MetricsMap["FsCapacityBytes"]), resource.BinarySI).String()



		fmt.Println("--------------------------------------------------------------------------------------------")
		fmt.Println("Cluster : ", clusterName)
		fmt.Println("--------------------------------------------------------------------------------------------")
		fmt.Println("cpuScore : ", cpuScore, "(",cpuUsed, " / ",  cpuTotal, ") , weight : ", ae.MetricsWeight["CPU"])
		fmt.Println("memScore : ", memScore, "(", memUsed, " / ", memTotal, ") , weight : ", ae.MetricsWeight["Memory"])
		//fmt.Println("netScore : ", netScore, "(", netUsed, " / ", netTotal, ")")
		fmt.Println("diskScore : ", diskScore, "(", diskUsed, "/", diskTotal, ") , weight : ", ae.MetricsWeight["FS"])
	*/
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("totalScore : ", score)
	fmt.Println("--------------------------------------------------------------------------------------------")

	return score
}

func (ae *AnalyticEngineStruct) SendLBAnalysis(ctx context.Context, in *protobuf.LBInfo) (*protobuf.ResponseLB, error) {
	fmt.Println("LB Requested")
	clusterNameList := in.ClusterNameList
	clusterScoreMap := make(map[string]float64)

	for _, clusterName := range clusterNameList {
		//*****************************
		//얘가 값을 제대로 못받아옴
		clusterScoreMap[clusterName] = ae.ResourceScore[clusterName]
		//clusterScoreMap[clusterName] = 100.0
	}

	fmt.Println(clusterScoreMap)

	fmt.Println("Geo Score")
	clientIP := in.ClientIP
	country := ae.getCountry(clientIP)
	continent := ae.getContinent(country)

	//clusters := []string{"cluster4", "cluster5", "cluster6"}
	score := ae.geoScore(clusterNameList, country, continent)

	for _, clusterName := range clusterNameList {
		fmt.Println(clusterScoreMap[clusterName])
		fmt.Println(score[clusterName])
		clusterScoreMap[clusterName] = clusterScoreMap[clusterName] + score[clusterName]
	}

	fmt.Println("LB Response")
	fmt.Println(clusterScoreMap)
	return &protobuf.ResponseLB{ScoreMap: clusterScoreMap}, nil
}

func (ae *AnalyticEngineStruct) SelectHPACluster(data *protobuf.HASInfo) []string {

	scoreMap := map[float64]string{}
	score := []float64{}

	for key, value := range ae.ResourceScore {
		if key != data.ClusterName && value > 0 {
			scoreMap[value] = key
			score = append(score, value)
		}
	}
	sort.Float64s(score)
	fmt.Println("scoreMap : ", scoreMap)

	filteringCluster := []string{}

	if len(score) > 1 {
		filteringCluster = append(filteringCluster, scoreMap[score[0]])
		filteringCluster = append(filteringCluster, scoreMap[score[1]])
	}

	fmt.Println("filteringCluster : ", filteringCluster)

	//result := ae.CompareHPAInfo(filteringCluster, data.HPAName, data.HPANamespace)
	//fmt.Println("result : ", result)

	return filteringCluster
}

func (ae *AnalyticEngineStruct) CompareHPAMaxInfo(clusterList []string, hpaName string, hpaNamespace string) string {
	replicasGap := map[string]int32{}
	rebalancingCount := map[string]int32{}

	cm := clusterManager.NewClusterManager()

	hpaInstance := &hpav2beta1.HorizontalPodAutoscaler{}
	for _, cluster := range clusterList {
		err := cm.Cluster_genClients[cluster].Get(context.TODO(), hpaInstance, hpaNamespace, hpaName)
		if err == nil {
			fmt.Println(cluster, " hpa : ", hpaInstance.Spec.MaxReplicas, " / ", hpaInstance.Status.CurrentReplicas)
			calc := hpaInstance.Spec.MaxReplicas - hpaInstance.Status.CurrentReplicas
			if calc > 0 {
				replicasGap[cluster] = calc
			}
		} else {
			fmt.Println(err)
		}
	}

	openmcphasInstance, err := cm.Crd_client.OpenMCPHybridAutoScaler(hpaNamespace).Get(hpaName, metav1.GetOptions{})

	if err == nil {
		//fmt.Println("success: ",openmcphasInstance)
		for cluster, _ := range replicasGap {
			rebalancingCount[cluster] = openmcphasInstance.Status.RebalancingCount[cluster]
		}
	} else {
		fmt.Println(err)
	}

	fmt.Println("desiredReplicas : ", replicasGap)
	fmt.Println("countRebalancing : ", rebalancingCount)

	result := ""

	for cluster, _ := range rebalancingCount {
		if result == "" {
			result = cluster
		} else {
			if (replicasGap[result] - rebalancingCount[result]) < (replicasGap[cluster] - rebalancingCount[cluster]) {
				result = cluster
			}
		}
		//fmt.Println(result)
	}

	return result
}

func (ae *AnalyticEngineStruct) CompareHPAMinInfo(clusterList []string, hpaName string, hpaNamespace string) string {
	replicasGap := map[string]int32{}
	rebalancingCount := map[string]int32{}

	cm := clusterManager.NewClusterManager()

	hpaInstance := &hpav2beta1.HorizontalPodAutoscaler{}
	for _, cluster := range clusterList {
		err := cm.Cluster_genClients[cluster].Get(context.TODO(), hpaInstance, hpaNamespace, hpaName)
		if err == nil {
			fmt.Println(cluster, " hpa : ", *hpaInstance.Spec.MinReplicas, " / ", hpaInstance.Status.CurrentReplicas)
			calc := hpaInstance.Status.CurrentReplicas - *hpaInstance.Spec.MinReplicas
			if calc > 0 {
				replicasGap[cluster] = calc
			}
		} else {
			fmt.Println(err)
		}
	}

	openmcphasInstance, err := cm.Crd_client.OpenMCPHybridAutoScaler(hpaNamespace).Get(hpaName, metav1.GetOptions{})

	if err == nil {
		//fmt.Println("success: ",openmcphasInstance)
		for cluster, _ := range replicasGap {
			rebalancingCount[cluster] = openmcphasInstance.Status.RebalancingCount[cluster]
		}
	} else {
		fmt.Println(err)
	}

	fmt.Println("desiredReplicas : ", replicasGap)
	fmt.Println("countRebalancing : ", rebalancingCount)

	result := ""

	for cluster, _ := range rebalancingCount {
		if result == "" {
			result = cluster
		} else {
			if (replicasGap[result] - rebalancingCount[result]) < (replicasGap[cluster] - rebalancingCount[cluster]) {
				result = cluster
			}
		}
		//fmt.Println(result)
	}

	return result
}

func (ae *AnalyticEngineStruct) SendHASMaxAnalysis(ctx context.Context, data *protobuf.HASInfo) (*protobuf.ResponseHAS, error) {
	fmt.Println("---------HAS Request Start---------")

	filteringCluster := ae.SelectHPACluster(data)
	fmt.Println("(Max)filteringCluster : " ,filteringCluster)
	result := ae.CompareHPAMaxInfo(filteringCluster, data.HPAName, data.HPANamespace)

	fmt.Println("---------HAS Response End---------")

	return &protobuf.ResponseHAS{TargetCluster: result}, nil
}

func (ae *AnalyticEngineStruct) SendHASMinAnalysis(ctx context.Context, data *protobuf.HASInfo) (*protobuf.ResponseHAS, error) {
	fmt.Println("---------HAS Request Start---------")

	filteringCluster := ae.SelectHPACluster(data)
	fmt.Println("(Min)filteringCluster : " ,filteringCluster)
	result := ae.CompareHPAMinInfo(filteringCluster, data.HPAName, data.HPANamespace)

	fmt.Println("---------HAS Response End---------")

	return &protobuf.ResponseHAS{TargetCluster: result}, nil
}

//
//func (ae *AnalyticEngineStruct) SendHASAnalysis(ctx context.Context, data *protobuf.HASInfo) (*protobuf.ResponseHAS, error) {
//	fmt.Println(">>> Get Request From HAS Controller [",data.ClusterName,"]")
//	//요청한 클러스터들의 매트릭 값 받아오기
//	influxData := ae.Influx.SelectMetricsData()
//	//매트릭 값을 기반으로 QoS 분석하기
//	result := ae.QoSAnalysisResult(data.HpaInfo, influxData)
//
//	fmt.Println(">>> Send Response(Result) To HAS Controller")
//	//결과값을 리턴값으로 넘겨주기
//	return &protobuf.ResponseHAS{
//		TargetCluster:          result,
//	}, nil
//}

/*func (ae *AnalyticEngineStruct) SearchInfluxDB(){

}*/
//
//type metricInfo struct {
//	clusterName string
//	metricValue interface{}
//}
//
//func (ae *AnalyticEngineStruct) QoSAnalysisResult(hpainfo []*protobuf.HASInfo, result []client.Result) string{
//	fmt.Println(">>> Analysis Rebalancing Metric")
//
//	clusterNum := len(hpainfo[0].CL)
//	filteringValue :=  make([]models.Row, 0)
//	resultMap := make(map[string]metricInfo)
//	countCluster := make(map[string]int)
//
//	tmp := 0
//	targetCluster := ""
//
//
//	resultMap["cpu_usage"] = metricInfo{clusterName: "", metricValue: float64(100000000000000000)}
//	resultMap["fs_usage"] = metricInfo{clusterName: "", metricValue: float64(100000000000000000)}
//	resultMap["memory_usage"] = metricInfo{clusterName: "", metricValue: float64(100000000000000000)}
//	resultMap["network_rx_usage"] = metricInfo{clusterName: "", metricValue: float64(100000000000000000)}
//	resultMap["network_tx_usage"] = metricInfo{clusterName: "", metricValue: float64(100000000000000000)}
//
//	if clusterNum == 1 {
//		targetCluster = result[0].Series[0].Tags["cluster"]
//	}else if clusterNum > 1{
//		for i := 0; i < clusterNum; i++ {
//			for _, data := range result[0].Series {
//				if hpainfo[0].CL[i].ClusterName == data.Tags["cluster"] {
//					filteringValue = append(filteringValue, data)
//				}
//			}
//		}
//
//		for _, data := range filteringValue {
//			//fmt.Println("===>", data.Tags["node"])
//			for i := 0; i < len(data.Columns); i++ {
//				if i != 0 && i != 1 && i != 4 {
//					if data.Values[0][i] != nil {
//						//	fmt.Println(data.Values[0][i].(string))
//						//	fmt.Println(bestResult[data.Columns[i]].metricValue)
//
//						var aa float64
//						if i == 2 {
//							aa, _ = strconv.ParseFloat(data.Values[0][i].(string)[:len(data.Values[0][i].(string))-1], 64)
//						} else if i == 3 || i == 5 {
//							aa, _ = strconv.ParseFloat(data.Values[0][i].(string)[:len(data.Values[0][i].(string))-2], 64)
//						} else {
//							aa, _ = strconv.ParseFloat(data.Values[0][i].(string), 64)
//						}
//
//						bb := resultMap[data.Columns[i]].metricValue.(float64)
//
//						min := math.Min(aa, bb)
//						resultMap[data.Columns[i]] = metricInfo{clusterName: data.Tags["cluster"], metricValue: min}
//						countCluster[data.Tags["cluster"]] = 0
//						//fmt.Println("[", data.Columns[i], "] compare : ", aa, " -- ", bb, " // min : ", min)
//					}
//				}
//			}
//			//fmt.Println()
//		}
//
//		for _, value := range resultMap {
//			countCluster[value.clusterName] += 1
//		}
//
//		for key, value := range countCluster {
//			if tmp < value {
//				targetCluster = key
//				tmp = value
//			}
//		}
//
//		//fmt.Println("*********compare Result : ", resultMap)
//	}
//
//
//	//fmt.Println("*********targetCluster : ", targetCluster)
//	fmt.Println("     => Anlysis Result [", targetCluster, "]")
//
//	return targetCluster
//}

func (ae *AnalyticEngineStruct) StartGRPC(GRPC_PORT string) {
	log.Printf("Grpc Server Start at Port %s\n", GRPC_PORT)

	//manager = NewClusterManager()
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

//*****************************************
//LoadBalancing

func (ae *AnalyticEngineStruct) geoScore(clusters []string, clientCountry, clientContinent string) map[string]float64 {
	fmt.Println("*****Geo Score*****")

	midScore := 100.0
	policy := ae.MetricsWeight["GeoRate"] * 100.0
	//policy := 70.0
	fmt.Println(ae.ClusterGeo)

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
	fmt.Println(score)
	return score
}

func (ae *AnalyticEngineStruct) getCountry(clientip string) string {
	fmt.Println("*****Extract Country*****")
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
