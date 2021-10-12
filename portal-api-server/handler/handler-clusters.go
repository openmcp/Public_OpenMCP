package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetJoinedClusters(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	// clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
	clusterurl := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters?clustername=openmcp"
	go CallAPI(token, clusterurl, ch)
	clusters := <-ch
	clusterData := clusters.data

	resCluster := ClustersRes{}

	//get clusters Information
	clusterNames := []string{}

	for _, element := range clusterData["items"].([]interface{}) {
		joinStatus := GetStringElement(element, []string{"spec", "joinStatus"})

		if joinStatus == "JOIN" {
			clusterName := GetStringElement(element, []string{"metadata", "name"})
			provider := GetStringElement(element, []string{"spec", "clusterPlatformType"})

			clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/namespaces/kube-federation-system/kubefedclusters/" + clusterName + "?clustername=openmcp"
			go CallAPI(token, clusterurl, ch)
			clusters := <-ch
			clusterData := clusters.data

			cluster := ClusterInfo{}
			clusterType := GetStringElement(clusterData["status"], []string{"conditions", "type"})
			if clusterType == "Ready" {
				clusterNames = append(clusterNames, clusterName)
				cluster.Name = clusterName
				cluster.Provider = provider

				// // if 조건으로 테스트용 Provider 입력해보자
				// if clusterName == "cluster1" {
				// 	provider = "On-Premis" //>>> KVM이라고할꺼임
				// } else if clusterName == "cluster2" {
				// 	provider = "GKE"
				// } else if clusterName == "openmcp" {
				// 	provider = "AKS"
				// } else {
				// 	provider = "EKS"
				// }
				resCluster.Clusters = append(resCluster.Clusters, cluster)
			}
		}
	}

	for i, cluster := range resCluster.Clusters {

		cluster.Status = "Healthy"

		// get node names, cpu(capacity)
		nodeURL := "https://" + openmcpURL + "/api/v1/nodes?clustername=" + cluster.Name
		go CallAPI(token, nodeURL, ch)
		nodeResult := <-ch
		nodeData := nodeResult.data
		nodeItems := nodeData["items"].([]interface{})

		cpuCapSum := 0
		memoryCapSum := 0
		fsCapSum := 0
		cpuUseSum := 0
		memoryUseSum := 0
		fsUseSum := 0
		networkSum := 0

		// get nodename, cpu capacity Information
		for _, element := range nodeItems {
			isMaster := GetStringElement(element, []string{"metadata", "labels", "node-role.kubernetes.io/master"})

			if isMaster != "-" && isMaster == "" {
				zone := GetStringElement(element, []string{"metadata", "labels", "topology.kubernetes.io/region"})
				region := GetStringElement(element, []string{"metadata", "labels", "topology.kubernetes.io/zone"})
				resCluster.Clusters[i].Zones = zone
				resCluster.Clusters[i].Region = region
			}

			nodeName := GetStringElement(element, []string{"metadata", "name"})
			// fmt.Println(nodeName)
			status := ""
			statusInfo := element.(map[string]interface{})["status"]
			//GetStringElement(element, []string{"status"})

			var healthCheck = make(map[string]string)
			for _, elem := range statusInfo.(map[string]interface{})["conditions"].([]interface{}) {
				conType := GetStringElement(elem, []string{"type"})
				// elem.(map[string]interface{})["type"].(string)
				tf := GetStringElement(elem, []string{"status"})
				// elem.(map[string]interface{})["status"].(string)
				healthCheck[conType] = tf
			}

			if healthCheck["Ready"] == "True" && (healthCheck["NetworkUnavailable"] == "" || healthCheck["NetworkUnavailable"] == "False") && healthCheck["MemoryPressure"] == "False" && healthCheck["DiskPressure"] == "False" && healthCheck["PIDPressure"] == "False" {
				// healthyNodeCnt++
				status = "Healthy"
			} else {
				if healthCheck["Ready"] == "Unknown" || (healthCheck["NetworkUnavailable"] == "" || healthCheck["NetworkUnavailable"] == "Unknown") || healthCheck["MemoryPressure"] == "Unknown" || healthCheck["DiskPressure"] == "Unknown" || healthCheck["PIDPressure"] == "Unknown" {
					status = "Unknown"
				} else {
					status = "Unhealthy"
				}
			}
			if status == "Healthy" {
				resCluster.Clusters[i].Status = "Healthy"
			} else {
				resCluster.Clusters[i].Status = "Unhealthy"
			}

			cpuCapacity := GetStringElement(element, []string{"status", "capacity", "cpu"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["capacity"].(map[string]interface{})["cpu"].(string)
			cpuCapInt, _ := strconv.Atoi(cpuCapacity)
			memoryCapacity := GetStringElement(element, []string{"status", "capacity", "memory"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["capacity"].(map[string]interface{})["memory"].(string)
			memoryCapacity = strings.Split(memoryCapacity, "Ki")[0]
			memoryCapInt, _ := strconv.Atoi(memoryCapacity)

			cpuCapSum += cpuCapInt
			memoryCapSum += memoryCapInt

			clMetricURL := "https://" + openmcpURL + "/metrics/nodes/" + nodeName + "?clustername=" + cluster.Name
			// fmt.Println("check usl ::: http://" + openmcpURL + "/metrics/nodes/" + nodeName + "?clustername=" + cluster.Name)
			go CallAPI(token, clMetricURL, ch)
			clMetricResult := <-ch
			clMetricData := clMetricResult.data

			cpuUse := "0n"
			memoryUse := "0Ki"
			fsUse := "0Ki"
			fsCap := "0Ki"
			ntRx := "0"
			ntTx := "0"
			//  cluster CPU Usage, Memroy Usage 확인
			if clMetricData["nodemetrics"] != nil {
				for _, element := range clMetricData["nodemetrics"].([]interface{}) {

					cpuUseCheck := GetInterfaceElement(element, []string{"cpu", "CPUUsageNanoCores"})
					// element.(map[string]interface{})["cpu"].(map[string]interface{})["CPUUsageNanoCores"]
					if cpuUseCheck == nil {
						cpuUse = "0n"
					} else {
						cpuUse = cpuUseCheck.(string)
					}
					cpuUse = strings.Split(cpuUse, "n")[0]
					cpuUseInt, _ := strconv.Atoi(cpuUse)

					memoryUseCheck := GetInterfaceElement(element, []string{"memory", "MemoryUsageBytes"})
					// element.(map[string]interface{})["memory"].(map[string]interface{})["MemoryUsageBytes"]
					if memoryUseCheck == nil {
						memoryUse = "0Ki"
					} else {
						memoryUse = memoryUseCheck.(string)
					}
					memoryUse = strings.Split(memoryUse, "Ki")[0]
					memoryUseInt, _ := strconv.Atoi(memoryUse)

					cpuUseSum += cpuUseInt
					memoryUseSum += memoryUseInt

					fsCapCheck := GetInterfaceElement(element, []string{"fs", "FsCapacityBytes"})
					// element.(map[string]interface{})["fs"].(map[string]interface{})["FsCapacityBytes"]
					if fsCapCheck == nil {
						fsCap = "0Ki"
					} else {
						fsCap = fsCapCheck.(string)
					}
					fsCap = strings.Split(fsCap, "Ki")[0]
					fsCapInt, _ := strconv.Atoi(fsCap)
					fsCapSum += fsCapInt

					fsUseCheck := GetInterfaceElement(element, []string{"fs", "FsUsedBytes"})
					// element.(map[string]interface{})["fs"].(map[string]interface{})["FsUsedBytes"]
					if fsUseCheck == nil {
						fsUse = "0Ki"
					} else {
						fsUse = fsUseCheck.(string)
					}
					fsUse = strings.Split(fsUse, "Ki")[0]
					fsUseInt, _ := strconv.Atoi(fsUse)
					fsUseSum += fsUseInt

					ntRxCheck := GetInterfaceElement(element, []string{"network", "NetworkRxBytes"})
					// element.(map[string]interface{})["network"].(map[string]interface{})["NetworkRxBytes"]
					if ntRxCheck == nil {
						ntRx = "0"
					} else {
						ntRx = ntRxCheck.(string)
					}
					ntTxCheck := GetInterfaceElement(element, []string{"network", "NetworkTxBytes"})
					// element.(map[string]interface{})["network"].(map[string]interface{})["NetworkTxBytes"]
					if ntTxCheck == nil {
						ntTx = "0"
					} else {
						ntTx = ntTxCheck.(string)
					}
					ntTxUseInt, _ := strconv.Atoi(ntTx)
					ntRxUseInt, _ := strconv.Atoi(ntRx)
					rTxSum := ntRxUseInt + ntTxUseInt

					networkSum += rTxSum

				}
			}
		}

		//calculate cpu, memory unit
		cpuUseSumF := float64(cpuUseSum) / 1000 / 1000 / 1000
		cpuUseSumS := fmt.Sprintf("%.1f", cpuUseSumF)
		memoryUseSumF := float64(memoryUseSum) / 1000 / 1000
		memoryUseSumS := fmt.Sprintf("%.1f", memoryUseSumF)
		memoryCapSumF := float64(memoryCapSum) / 1000 / 1000
		memoryCapSumS := fmt.Sprintf("%.1f", memoryCapSumF)

		fsUseSumF := float64(fsUseSum) / 1000 / 1000
		fsUseSumS := fmt.Sprintf("%.1f", fsUseSumF)
		fsCapSumF := float64(fsCapSum) / 1000 / 1000
		fsCapSumS := fmt.Sprintf("%.1f", fsCapSumF)
		networkCapSumF := float64(networkSum) / 1000 / 1000 / 1000 / 1000
		networkCapSumS := fmt.Sprintf("%.1f", networkCapSumF)
		// networkSumS := strconv.Itoa(networkSum)

		// fmt.Println(fsUseSumS, fsCapSumS)

		var cpuStatus []NameVal
		var memStatus []NameVal
		var fsStatus []NameVal

		cpuStatus = append(cpuStatus, NameVal{"Used", cpuUseSumF})
		cpuStatus = append(cpuStatus, NameVal{"Total", float64(cpuCapSum)})
		cpuUnit := Unit{"core", cpuStatus}

		memStatus = append(memStatus, NameVal{"Used", memoryUseSumF})
		memStatus = append(memStatus, NameVal{"Total", memoryCapSumF})
		memUnit := Unit{"Gi", memStatus}

		fsStatus = append(fsStatus, NameVal{"Used", fsUseSumF})
		if fsCapSumF == 0 {
			fsCapSumF = 100.0
		}
		fsStatus = append(fsStatus, NameVal{"Total", fsCapSumF})
		fsUnit := Unit{"Gi", fsStatus}

		resUsage := ClusterResourceUsage{cpuUnit, memUnit, fsUnit}

		resCluster.Clusters[i].Nodes = len(nodeItems)
		resCluster.Clusters[i].Cpu = cpuUseSumS + "/" + strconv.Itoa(cpuCapSum) + " Core"
		resCluster.Clusters[i].Ram = memoryUseSumS + "/" + memoryCapSumS + " Gi"
		resCluster.Clusters[i].Disk = PercentUseString(fsUseSumS, fsCapSumS) + "%"
		resCluster.Clusters[i].Network = networkCapSumS + " byte/s"
		resCluster.Clusters[i].ResourceUsage = resUsage
	}
	// fmt.Println(resCluster.Clusters)
	json.NewEncoder(w).Encode(resCluster.Clusters)
}

func GetJoinableClusters(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusOK)

	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	// url := "https://" + openmcpURL + "/joinable" //못쓴다 이제... 전체 cluster검색후 metadata : joinStatus: UNJOIN인것들만 모아서 뿌려야함.

	clusterurl := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters?clustername=openmcp"

	go CallAPI(token, clusterurl, ch)
	clusters := <-ch
	clusterData := clusters.data

	type joinable struct {
		Name     string `json:"name"`
		Endpoint string `json:"endpoint"`
		Provider string `json:"provider"`
		Region   string `json:"region"`
		Zone     string `json:"zone"`
	}

	var joinableLists []joinable

	if clusterData["items"] != nil {
		for _, element := range clusterData["items"].([]interface{}) {
			joinStatus := GetStringElement(element, []string{"spec", "joinStatus"})
			if joinStatus == "UNJOIN" {
				clusterName := GetStringElement(element, []string{"metadata", "name"})
				// endpoint := element.(map[string]interface{})["endpoint"].(string)
				provider := GetStringElement(element, []string{"spec", "clusterPlatformType"})
				region := "-"
				zone := "-"

				clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
				go CallAPI(token, clusterurl, ch)
				clusters2 := <-ch
				clusterData2 := clusters2.data

				endpoint := ""
				for _, element := range clusterData2["items"].([]interface{}) {
					targetClusterName := GetStringElement(element, []string{"metadata", "name"})

					if clusterName == targetClusterName {
						endpoint = GetStringElement(element, []string{"spec", "apiEndpoint"})
						endpoint = strings.Replace(endpoint, "https://", "", -1)
						endpoint = strings.Replace(endpoint, "http://", "", -1)
						address := strings.Split(endpoint, ":")
						endpoint = address[0]
					}
				}

				nodeURL := "https://" + openmcpURL + "/api/v1/nodes?clustername=" + clusterName
				go CallAPI(token, nodeURL, ch)
				nodeResult := <-ch
				nodeData := nodeResult.data
				nodeItems := nodeData["items"].([]interface{})

				for _, element := range nodeItems {
					isMaster := GetStringElement(element, []string{"metadata", "labels", "node-role.kubernetes.io/master"})

					if isMaster != "-" && isMaster == "" {
						zone = GetStringElement(element, []string{"metadata", "labels", "topology.kubernetes.io/region"})
						region = GetStringElement(element, []string{"metadata", "labels", "topology.kubernetes.io/zone"})
					}
				}

				res := joinable{clusterName, endpoint, provider, region, zone}
				joinableLists = append(joinableLists, res)
			}
		}
		json.NewEncoder(w).Encode(joinableLists)
	} else {
		json.NewEncoder(w).Encode(joinableLists)
	}
}

func ClusterOverview(w http.ResponseWriter, r *http.Request) {
	clusterNm := r.URL.Query().Get("clustername")
	if clusterNm == "" {
		errorMSG := jsonErr{500, "failed", "need some params"}
		json.NewEncoder(w).Encode(errorMSG)
	} else {
		ch := make(chan Resultmap)
		token := GetOpenMCPToken()

		clusterurl := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/" + clusterNm + "?clustername=openmcp"
		// https://192.168.0.152:30000/apis/core.kubefed.io/v1beta1/namespaces/kube-federation-system/kubefedclusters/cluster1?clustername=openmcp

		// https://192.168.0.152:30000/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/cluster1?clustername=openmcp

		go CallAPI(token, clusterurl, ch)
		clusters := <-ch
		clusterData := clusters.data

		//get clusters Information
		// clusterNames := []string{}
		//set master cluster info
		// clusterNames = append(clusterNames, "openmcp")
		region := "-"
		zone := "-"
		provider := GetStringElement(clusterData["spec"], []string{"clusterPlatformType"})

		// region = GetStringElement(clusterData["status"], []string{"region"})
		// zone = GetStringElement(clusterData["status"], []string{"zones"})
		// if clusterNm == "cluster1" {
		// 	provider = "eks"
		// } else if clusterNm == "cluster2" {
		// 	provider = "kvm"
		// } else if clusterNm == "openmcp" {
		// 	provider = "aks"
		// } else {
		// 	provider = "-"
		// }

		// GetInfluxNodesMetric()
		InitInfluxConfig()
		inf := NewInflux(InfluxConfig.Influx.Ip, InfluxConfig.Influx.Port, InfluxConfig.Influx.Username, InfluxConfig.Influx.Username)
		results := GetInfluxPodsMetric(clusterNm, inf)

		type podUsage struct {
			namespace string
			cpuUsage  int
			memUsage  int
		}

		namspaceUsage := make(map[string]podUsage)
		cpuSum := make(map[string]int)
		memSum := make(map[string]int)

		for _, result := range results {
			if result.Series != nil {
				for _, ser := range result.Series {
					namespace := ser.Tags["namespace"]
					// pod := ser.Tags["pod"]
					cpu := 0
					mem := 0
					if ser.Values[0][1] != nil {
						cpuUsage := ser.Values[0][1].(string)
						if cpuUsage != "0" {
							cpuUsage = cpuUsage[:len(cpuUsage)-1]
							cpu, _ = strconv.Atoi(cpuUsage)
						}
					}
					if ser.Values[0][6] != nil {
						memUsage := ser.Values[0][6].(string)
						if memUsage != "0" {
							memUsage = memUsage[:len(memUsage)-2]
							mem, _ = strconv.Atoi(memUsage)
						}
					}

					cpuSum[namespace] += cpu
					memSum[namespace] += mem
					namspaceUsage[namespace] = podUsage{namespace, cpuSum[namespace], memSum[namespace]}

				}
			}
		}

		nodeURL := "https://" + openmcpURL + "/api/v1/nodes?clustername=" + clusterNm
		go CallAPI(token, nodeURL, ch)

		nodeResult := <-ch
		nodeData := nodeResult.data
		nodeItems := nodeData["items"].([]interface{})

		clusterCPUCapSum := 0
		clusterMemoryCapSum := 0

		nodeResCPU := make(map[string]float64)
		nodeResMem := make(map[string]float64)
		healthyNodeCnt := 0
		unknownNodeCnt := 0
		unhealthyNodeCnt := 0
		var nodeResCPUSum float64
		var nodeResMemSum float64
		var nodeResFSSum int
		var nodeResFSCapaSum int

		var nodeNameList []string
		var kubeVersion string
		for _, element := range nodeItems {
			isMaster := GetStringElement(element, []string{"metadata", "labels", "node-role.kubernetes.io/master"})
			if isMaster != "-" && isMaster == "" {
				zone = GetStringElement(element, []string{"metadata", "labels", "topology.kubernetes.io/region"})
				region = GetStringElement(element, []string{"metadata", "labels", "topology.kubernetes.io/zone"})
			}

			status := element.(map[string]interface{})["status"]
			var healthCheck = make(map[string]string)
			kubeVersion = status.(map[string]interface{})["nodeInfo"].(map[string]interface{})["kubeletVersion"].(string)
			for _, elem := range status.(map[string]interface{})["conditions"].([]interface{}) {
				conType := elem.(map[string]interface{})["type"].(string)
				tf := elem.(map[string]interface{})["status"].(string)
				healthCheck[conType] = tf

			}

			if healthCheck["Ready"] == "True" && (healthCheck["NetworkUnavailable"] == "" || healthCheck["NetworkUnavailable"] == "False") && healthCheck["MemoryPressure"] == "False" && healthCheck["DiskPressure"] == "False" && healthCheck["PIDPressure"] == "False" {
				healthyNodeCnt++
			} else {
				if healthCheck["Ready"] == "Unknown" || (healthCheck["NetworkUnavailable"] == "" || healthCheck["NetworkUnavailable"] == "Unknown") || healthCheck["MemoryPressure"] == "Unknown" || healthCheck["DiskPressure"] == "Unknown" || healthCheck["PIDPressure"] == "Unknown" {
					unknownNodeCnt++
				} else {
					unhealthyNodeCnt++
				}
			}

			nodeName := element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
			cpuCapacity := element.(map[string]interface{})["status"].(map[string]interface{})["capacity"].(map[string]interface{})["cpu"].(string)
			cpuCapInt, _ := strconv.Atoi(cpuCapacity)
			nodeNameList = append(nodeNameList, nodeName)
			memoryCapacity := element.(map[string]interface{})["status"].(map[string]interface{})["capacity"].(map[string]interface{})["memory"].(string)
			memoryCapacity = strings.Split(memoryCapacity, "Ki")[0]
			memoryCapInt, _ := strconv.Atoi(memoryCapacity)

			clusterCPUCapSum += cpuCapInt
			clusterMemoryCapSum += memoryCapInt

			clMetricURL := "https://" + openmcpURL + "/metrics/nodes/" + nodeName + "?clustername=" + clusterNm
			go CallAPI(token, clMetricURL, ch)
			clMetricResult := <-ch
			clMetricData := clMetricResult.data

			cpuUse := "0n"
			memoryUse := "0Ki"
			fsUse := "0Ki"
			fsCapaUse := "0Ki"

			// fmt.Println("clusterCPUCapSum", clusterCPUCapSum)
			//  cluster CPU Usage, Memroy Usage 확인
			if clMetricData["nodemetrics"] != nil {
				for _, element := range clMetricData["nodemetrics"].([]interface{}) {
					cpuUseCheck := element.(map[string]interface{})["cpu"].(map[string]interface{})["CPUUsageNanoCores"]
					if cpuUseCheck == nil {
						cpuUse = "0n"
					} else {
						cpuUse = cpuUseCheck.(string)
					}
					cpuUse = strings.Split(cpuUse, "n")[0]
					cpuUseInt, _ := strconv.Atoi(cpuUse)

					memoryUseCheck := element.(map[string]interface{})["memory"].(map[string]interface{})["MemoryUsageBytes"]
					if memoryUseCheck == nil {
						memoryUse = "0Ki"
					} else {
						memoryUse = memoryUseCheck.(string)
					}
					memoryUse = strings.Split(memoryUse, "Ki")[0]
					memoryUseInt, _ := strconv.Atoi(memoryUse)

					fsUseCheck := element.(map[string]interface{})["fs"].(map[string]interface{})["FsUsedBytes"]
					if fsUseCheck == nil {
						fsUse = "0Ki"
					} else {
						fsUse = fsUseCheck.(string)
					}
					fsUse = strings.Split(fsUse, "Ki")[0]
					fsUseInt, _ := strconv.Atoi(fsUse)

					fsCapaCheck := element.(map[string]interface{})["fs"].(map[string]interface{})["FsCapacityBytes"]
					if fsCapaCheck == nil {
						fsCapaUse = "0Ki"
					} else {
						fsCapaUse = fsCapaCheck.(string)
					}
					fsCapaUse = strings.Split(fsCapaUse, "Ki")[0]
					fsCapaUseInt, _ := strconv.Atoi(fsCapaUse)

					nodeResCPU[nodeName] = math.Ceil((float64(cpuUseInt)/1000/1000/1000)*10000) / 10000
					nodeResMem[nodeName] = math.Ceil((float64(memoryUseInt)/1000/1000)*10000) / 10000

					nodeResCPUSum += nodeResCPU[nodeName]
					nodeResMemSum += nodeResMem[nodeName]
					nodeResFSSum += fsUseInt
					nodeResFSCapaSum += fsCapaUseInt
				}
			}

		}

		clusterCPURes := make(map[string]float64)
		clusterMemoryRes := make(map[string]float64)
		clusterCPUCapSum = clusterCPUCapSum * 1000 * 1000 * 1000
		for _, res := range namspaceUsage {
			cpuval := PercentChange(float64(res.cpuUsage), float64(clusterCPUCapSum))
			clusterCPURes[res.namespace] = math.Ceil(cpuval*100) / 100
			clusterMemoryRes[res.namespace] = float64(res.memUsage / 1000)
		}

		clusterCPURank := reverseRank(clusterCPURes, 5)
		clusterMemRank := reverseRank(clusterMemoryRes, 5)

		// fmt.Println(clusterCPURank, clusterMemRank)

		// for _, r := range nodeNameList {
		// 	nodeCPUPecent := PercentChange(float64(nodeResCPU[r]), float64(clusterCPUCapSum))
		// 	nodeResCPU[r] = math.Ceil(nodeCPUPecent*100) / 100
		// }

		nodeCPURank := reverseRank(nodeResCPU, 5)
		nodeMemRank := reverseRank(nodeResMem, 5)

		// nodeResCPUSumStr := fmt.Sprintf("%.1f", nodeResCPUSum/1000/1000/1000)
		nodeResCPUSumStr := nodeResCPUSum / 1000 / 1000 / 1000
		// nodeResMemSumStr := fmt.Sprintf("%.1f", nodeResMemSum/1000)
		nodeResMemSumStr := nodeResMemSum / 1000 / 1000 //Gi
		// nodeResFSSumStr := fmt.Sprintf("%.1f", float64(nodeResFSSum)/1000/1000)
		nodeResFSSumStr := float64(nodeResFSSum) / 1000 / 1000
		// nodeResFSCapaSumStr := fmt.Sprintf("%.1f", float64(nodeResFSCapaSum)/1000/1000)
		nodeResFSCapaSumStr := float64(nodeResFSCapaSum) / 1000 / 1000
		// fmt.Println(nodeResCPUSumStr, nodeResMemSumStr, nodeResFSSum, nodeResFSCapaSum)

		//cat ~/.kube/config > config (copy and paste)
		config, _ := buildConfigFromFlags(clusterNm, kubeConfigFile)
		clientset, _ := kubernetes.NewForConfig(config)
		compStatus, _ := clientset.CoreV1().ComponentStatuses().List(v1.ListOptions{})

		nodeStatus := ""
		if unknownNodeCnt > 0 || unhealthyNodeCnt > 0 {
			nodeStatus = "Unhealthy"
		} else {
			nodeStatus = "Healthy"
		}

		var kubeStatus []NameStatus

		for _, r := range compStatus.Items {
			name := r.Name
			conType := r.Conditions
			// status := r.Conditions[1]
			state := "Healthy"
			for _, c := range conType {
				if c.Type == "Healthy" && c.Status == "True" {
					state = "Healthy"
				} else {
					state = "Unhealthy"
					// clusterStatus = "Unhealthy"
				}
			}

			if strings.Contains(name, "etcd") {
				name = "etcd"
			} else if strings.Contains(name, "scheduler") {
				name = "Scheduler"
			} else if strings.Contains(name, "controller-manager") {
				name = "Controller Manager"
			}
			kubeStatus = append(kubeStatus, NameStatus{name, state})
		}
		// fmt.Println(kubeStatus)
		kubeStatus = append(kubeStatus, NameStatus{"Nodes", nodeStatus})

		eventsURL := "https://" + openmcpURL + "/api/v1/events?clustername=" + clusterNm
		go CallAPI(token, eventsURL, ch)
		eventsResult := <-ch
		eventsData := eventsResult.data
		eventsItems := eventsData["items"].([]interface{})
		var events []Event
		if eventsItems != nil {
			for _, element := range eventsItems {
				project := element.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
				typeNm := element.(map[string]interface{})["type"].(string)
				reason := element.(map[string]interface{})["reason"].(string)
				object := element.(map[string]interface{})["involvedObject"].(map[string]interface{})["kind"].(string)
				message := element.(map[string]interface{})["message"].(string)
				time := "-"
				if element.(map[string]interface{})["lastTimestamp"] != nil {
					time = element.(map[string]interface{})["lastTimestamp"].(string)
				}
				events = append(events, Event{project, typeNm, reason, object, message, time})
			}
		}

		var cpuStaus []NameVal
		var memStaus []NameVal
		var fsStaus []NameVal
		cpuStaus = append(cpuStaus, NameVal{"Used", math.Ceil(nodeResCPUSumStr*100) / 100})
		// cpuStaus = append(cpuStaus, NameVal{"Total", fmt.Sprintf("%.1f", float64(clusterCPUCapSum)/1000/1000/1000)})
		cpuStaus = append(cpuStaus, NameVal{"Total", math.Ceil(float64(clusterCPUCapSum)/1000/1000/1000*100) / 100})
		memStaus = append(memStaus, NameVal{"Used", math.Ceil(nodeResMemSumStr*100) / 100})
		// memStaus = append(memStaus, NameVal{"Total", fmt.Sprintf("%.1f", float64(clusterMemoryCapSum)/1000/1000)})
		memStaus = append(memStaus, NameVal{"Total", math.Ceil(float64(clusterMemoryCapSum)/1000/1000*100) / 100})
		fsStaus = append(fsStaus, NameVal{"Used", math.Ceil(nodeResFSSumStr*100) / 100})
		fsStaus = append(fsStaus, NameVal{"Total", math.Ceil(nodeResFSCapaSumStr*100) / 100})
		cpuUnit := Unit{"core", cpuStaus}
		memUnit := Unit{"Gi", memStaus}
		fsUnit := Unit{"Gi", fsStaus}
		cUsage := ClusterResourceUsage{cpuUnit, memUnit, fsUnit}
		info := BasicInfo{clusterNm, provider, kubeVersion, nodeStatus, region, zone}

		pUsageTop5 := ProjectUsageTop5{clusterCPURank, clusterMemRank}
		nUsageTop5 := NodeUsageTop5{nodeCPURank, nodeMemRank}
		responseJSON := ClusterOverView{info, pUsageTop5, nUsageTop5, cUsage, kubeStatus, events}

		json.NewEncoder(w).Encode(responseJSON)

	}
}

func OpenMCPJoin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("OpenMCPJoin")
	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	clusterName := data["clusterName"].(string)
	clusterAddress := data["clusterAddress"].(string)

	if clusterAddress == "" {
		ch := make(chan Resultmap)
		token := GetOpenMCPToken()
		clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
		go CallAPI(token, clusterurl, ch)
		clusters := <-ch
		clusterData := clusters.data

		for _, element := range clusterData["items"].([]interface{}) {
			targetClusterName := GetStringElement(element, []string{"metadata", "name"})

			if clusterName == targetClusterName {
				endpoint := GetStringElement(element, []string{"spec", "apiEndpoint"})
				endpoint = strings.Replace(endpoint, "https://", "", -1)
				endpoint = strings.Replace(endpoint, "http://", "", -1)
				address := strings.Split(endpoint, ":")
				clusterAddress = address[0]
			}
		}
	}

	type JoinInfo struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}

	// [
	// 		{"op": "replace", "path": "/spec/joinStatus","value": "JOINING"},
	// 		{"op": "replace", "path": "/spec/metalLBRange/addressFrom", "value": "$ADDRESSFROM"},
	// 		{"op": "replace", "path": "/spec/metalLBRange/addressTo", "value": "$ADDRESSTO"}
	// ]

	var body []interface{}
	body = append(body, JoinInfo{"replace", "/spec/joinStatus", "JOINING"})
	body = append(body, JoinInfo{"replace", "/spec/metalLBRange/addressFrom", clusterAddress})
	body = append(body, JoinInfo{"replace", "/spec/metalLBRange/addressTo", openmcpAddress})

	fmt.Println(body)

	var jsonErrs []jsonErr

	projectURL := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/" + clusterName + "?clustername=" + openmcpClusterName

	resp, err := CallPatchAPI(projectURL, "application/json-patch+json", body, true)
	var msg jsonErr

	if err != nil {
		msg = jsonErr{503, "failed", "request fail"}
	}

	var dataRes map[string]interface{}
	json.Unmarshal([]byte(resp), &dataRes)
	if dataRes != nil {
		if dataRes["kind"].(string) == "Status" {
			msg = jsonErr{501, "failed", dataRes["message"].(string)}
		} else {
			msg = jsonErr{200, "success", "Cluster Join Completed"}
		}
	}

	jsonErrs = append(jsonErrs, msg)
	json.NewEncoder(w).Encode(jsonErrs)
}

func OpenMCPUnjoin(w http.ResponseWriter, r *http.Request) {
	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	clusterName := data["clusterName"].(string)
	fmt.Println(clusterName)

	type JoinInfo struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}

	// [
	// 		{"op": "replace", "path": "/spec/joinStatus","value": "JOINING"},
	// 		{"op": "replace", "path": "/spec/metalLBRange/addressFrom", "value": "$ADDRESSFROM"},
	// 		{"op": "replace", "path": "/spec/metalLBRange/addressTo", "value": "$ADDRESSTO"}
	// ]

	var body []interface{}
	body = append(body, JoinInfo{"replace", "/spec/metalLBRange/addressFrom", ""})
	body = append(body, JoinInfo{"replace", "/spec/metalLBRange/addressTo", ""})
	body = append(body, JoinInfo{"replace", "/spec/joinStatus", "UNJOIN"})

	var jsonErrs []jsonErr

	//https://192.168.0.152:30000/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/cluster2?clustername=openmcp
	projectURL := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/" + clusterName + "?clustername=" + openmcpClusterName

	resp, err := CallPatchAPI(projectURL, "application/json-patch+json", body, true)
	var msg jsonErr

	if err != nil {
		msg = jsonErr{503, "failed", "request fail"}
	}

	var dataRes map[string]interface{}
	json.Unmarshal([]byte(resp), &dataRes)
	if dataRes != nil {
		if dataRes["kind"].(string) == "Status" {
			msg = jsonErr{501, "failed", dataRes["message"].(string)}
		} else {
			msg = jsonErr{200, "success", "Cluster UnJoin Completed"}
		}
	}

	jsonErrs = append(jsonErrs, msg)

	json.NewEncoder(w).Encode(jsonErrs)
}
