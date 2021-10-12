package handler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type jsonErr struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
	Text   string `json:"text"`
}

type Resultmap struct {
	secs float64
	url  string
	data map[string]interface{}
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetOpenMCPToken() string {
	// var client http.Client

	client := &http.Client{
		Transport: &http.Transport{
			// TLSClientConfig: &tls.Config{
			// 	RootCAs: caCertPool,
			// },
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	account := Account{"openmcp", "keti"}

	pbytes, _ := json.Marshal(account)
	buff := bytes.NewBuffer(pbytes)

	// resp, err := client.Get("https://" + openmcpURL + "/token?username=openmcp&password=keti")
	resp, err := client.Post("https://"+openmcpURL+"/token", "application/json", buff)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	var data map[string]interface{}
	token := ""

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal([]byte(bodyBytes), &data)
		token = data["token"].(string)

	} else {
		fmt.Println("failed")
	}

	return token
}

func CallAPI(token string, url string, ch chan<- Resultmap) {
	start := time.Now()
	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	// var client http.Client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)

	if err != nil {
		// log.Fatal(err)
		fmt.Print(err)
	}
	var data map[string]interface{}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close() // 리소스 누출 방지
	if err != nil {
		// ch <- fmt.Sprintf("while reading %s: %v", url, err)
		// return
		fmt.Print(err)
	}
	json.Unmarshal([]byte(bodyBytes), &data)

	secs := time.Since(start).Seconds()

	// ch <- fmt.Sprintf("%.2fs %s %v", secs, url, data)

	ch <- Resultmap{secs, url, data}
}

func PostYaml(url string, yaml io.Reader) ([]byte, error) {
	token := GetOpenMCPToken()
	// fmt.Println("yaml   :", yaml)
	var bearer = "Bearer " + token
	req, err := http.NewRequest("POST", url, yaml)

	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	// var client http.Client

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	str := string(respBody)
	fmt.Println(str)
	return respBody, nil

}

func CallPostAPI(url string, headtype string, body interface{}) ([]byte, error) {
	token := GetOpenMCPToken()
	// fmt.Println("yaml   :", yaml)
	var bearer = "Bearer " + token

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	req, err := http.NewRequest("POST", url, payloadBuf)

	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	// var client http.Client

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	str := string(respBody)
	fmt.Println(str)
	return respBody, nil
}

func CallPatchAPI(url string, headtype string, body []interface{}, bodyIsArray bool) ([]byte, error) {
	token := GetOpenMCPToken()
	// fmt.Println("yaml   :", yaml)
	var bearer = "Bearer " + token

	payloadBuf := new(bytes.Buffer)
	if bodyIsArray {
		json.NewEncoder(payloadBuf).Encode(body)
	} else {
		json.NewEncoder(payloadBuf).Encode(body[0])
	}

	req, err := http.NewRequest("PATCH", url, payloadBuf)

	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", headtype)
	// Send req using http Client
	// var client http.Client

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	str := string(respBody)
	fmt.Println(str)
	return respBody, nil
}

func NodeHealthCheck(condType string) string {
	result := ""

	return result
}

func ClusterHealthCheck(condType string) string {
	result := ""

	return result
}

func GetInfluxPodsMetric(clusterName string, in *Influx) []client.Result {
	q := client.Query{}
	// q = client.NewQuery("SELECT last(*) FROM Pods WHERE cluster = '"+clusterName+"' ORDER BY DESC LIMIT 1", "Metrics", "")
	// select last(*) from Pods where time > now() - 5m and cluster='cluster1' group by namespace,pod order by desc limit 1
	q = client.NewQuery("select last(*) from Pods where time > now() - 5m and cluster='"+clusterName+"' group by namespace,pod order by desc limit 1", "Metrics", "")

	//select last(*) from Pods where time > now() -1m and cluster='cluster1' group by namespace,pod,time(1m)
	response, err := in.inClient.Query(q)
	// fmt.Println(response)

	if err == nil && response.Error() == nil {

		return response.Results
	}

	return nil
}

// func GetInfluxPod10mMetric(clusterName string, namespace string, pod string, in *Influx) []client.Result {
func GetInfluxPod10mMetric(clusterName string, namespace string, pod string) PhysicalResources {
	nowTime := time.Now().UTC() //.Add(time.Duration(offset) * time.Second)
	// startTime := nowTime.Add(time.Duration(-10) * time.Minute)
	endTime := nowTime
	startTime := nowTime.Add(time.Duration(-11) * time.Minute)
	// endTime := nowTime.Add(time.Duration(-1) * time.Minute)
	_, offset := time.Now().Zone()
	start := startTime.Format("2006-01-02_15:04") + ":00"
	end := endTime.Format("2006-01-02_15:04") + ":00"
	// fmt.Println(start, "///", end)

	ch := make(chan Resultmap)
	token := GetOpenMCPToken()
	// http://192.168.0.152:31635/metrics/namespaces/kube-system/pods/kube-flannel-ds-nn5p5?clustername=cluster1&timeStart=2020-09-03_09:00:00&timeEnd=2020-09-03_09:00:15

	podMetricURL := "https://" + openmcpURL + "/metrics/namespaces/" + namespace + "/pods/" + pod + "?clustername=" + clusterName + "&timeStart=" + start + "&timeEnd=" + end

	go CallAPI(token, podMetricURL, ch)

	podMetricResult := <-ch
	podMetricData := podMetricResult.data["podmetrics"]
	var podCPUUsageMins []PodCPUUsageMin
	var podMemoryUsageMins []PodMemoryUsageMin
	var podNetworkUsageMins []PodNetworkUsageMin
	if podMetricData != nil {

		metricsPerMin := make(map[string][]interface{})
		for _, m := range podMetricData.([]interface{}) {
			times := m.(map[string]interface{})["time"].(string)
			ind := strings.Index(times, ":")
			timeHM := times[ind-2 : ind+3]
			timeHM = timeHM + ":00"
			t1, _ := time.Parse("15:04:05", timeHM)
			t1 = t1.Add(time.Duration(offset) * time.Second)
			timeHM = t1.Format("15:04:05")

			// fmt.Println(timeHM)

			metricsPerMin[timeHM] = append(metricsPerMin[timeHM], m)
		}

		for k, m := range metricsPerMin {
			// fmt.Println(k)
			// fmt.Println(m)
			// fmt.Println()
			cpuSum := 0
			memorySum := 0
			oldNtTxUseInt := 0
			oldNtRxUseInt := 0
			maxTxUseInt := 0
			minTxUseInt := 0
			maxRxUseInt := 0
			minRxUseInt := 0

			for index, v := range m {
				if v.(map[string]interface{})["cpu"].(map[string]interface{})["CPUUsageNanoCores"] != nil {
					cpuUse := v.(map[string]interface{})["cpu"].(map[string]interface{})["CPUUsageNanoCores"].(string)
					cpuUse = strings.Split(cpuUse, "n")[0]
					cpuUseInt, _ := strconv.Atoi(cpuUse)
					cpuSum += cpuUseInt
				}

				if v.(map[string]interface{})["memory"].(map[string]interface{})["MemoryUsageBytes"] != nil {
					memoryUse := v.(map[string]interface{})["memory"].(map[string]interface{})["MemoryUsageBytes"].(string)
					memoryUse = strings.Split(memoryUse, "Ki")[0]
					memoryUseInt, _ := strconv.Atoi(memoryUse)
					memorySum += memoryUseInt
				}
				ntTxUseInt := 0
				ntRxUseInt := 0
				if v.(map[string]interface{})["network"].(map[string]interface{})["NetworkTxBytes"] != nil {
					ntTxUse := v.(map[string]interface{})["network"].(map[string]interface{})["NetworkTxBytes"].(string)
					ntTxUseInt, _ = strconv.Atoi(ntTxUse)
				}

				if v.(map[string]interface{})["network"].(map[string]interface{})["NetworkRxBytes"] != nil {
					ntRxUse := v.(map[string]interface{})["network"].(map[string]interface{})["NetworkRxBytes"].(string)
					ntRxUseInt, _ = strconv.Atoi(ntRxUse)
				}
				// fmt.Println(v.(map[string]interface{})["time"], ntTxUseInt, ntRxUseInt)

				if index == 0 {
					// fmt.Println()
					// fmt.Println("index0:", index, ntTxUseInt, oldNtTxUseInt)
					oldNtTxUseInt = ntTxUseInt
					oldNtRxUseInt = ntRxUseInt
					minTxUseInt = ntTxUseInt
					minRxUseInt = ntRxUseInt
					maxTxUseInt = ntTxUseInt
					maxRxUseInt = ntRxUseInt
				} else {
					// fmt.Println("indexN:", ntTxUseInt, oldNtTxUseInt)
					if oldNtTxUseInt <= ntTxUseInt {
						maxTxUseInt = ntTxUseInt
					}
					if oldNtRxUseInt <= ntRxUseInt {
						maxRxUseInt = ntRxUseInt
					}

					oldNtTxUseInt = ntTxUseInt
					oldNtRxUseInt = ntRxUseInt
				}
			}

			cpuAvg := float64(cpuSum) / float64(len(m)) / 1000 / 1000 / 1000
			memoryAvg := float64(memorySum) / float64(len(m)) / 1000
			// fmt.Println(maxTxUseInt, minTxUseInt)
			// fmt.Println(maxRxUseInt, minRxUseInt)
			inBps := (maxTxUseInt - minTxUseInt) / 60
			outBps := (maxRxUseInt - minRxUseInt) / 60
			//fmt.Println(k, "cpu: ", cpuAvg)
			podCPUUsageMins = append(podCPUUsageMins, PodCPUUsageMin{math.Ceil(cpuAvg*1000) / 1000, k})
			podMemoryUsageMins = append(podMemoryUsageMins, PodMemoryUsageMin{math.Ceil(memoryAvg*10) / 10, k})
			podNetworkUsageMins = append(podNetworkUsageMins, PodNetworkUsageMin{"Bps", inBps, outBps, k})

		}
		sort.Slice(podCPUUsageMins, func(i, j int) bool {
			return podCPUUsageMins[i].Time < podCPUUsageMins[j].Time
		})
		sort.Slice(podMemoryUsageMins, func(i, j int) bool {
			return podMemoryUsageMins[i].Time < podMemoryUsageMins[j].Time
		})
		sort.Slice(podNetworkUsageMins, func(i, j int) bool {
			return podNetworkUsageMins[i].Time < podNetworkUsageMins[j].Time
		})

		if len(podCPUUsageMins) > 10 {
			podCPUUsageMins = podCPUUsageMins[1:]
			podMemoryUsageMins = podMemoryUsageMins[1:]
			podNetworkUsageMins = podNetworkUsageMins[1:]
		}
		result := PhysicalResources{podCPUUsageMins, podMemoryUsageMins, podNetworkUsageMins}
		return result

	} else {

		podCPUUsageMins = append(podCPUUsageMins, PodCPUUsageMin{float64(0), ""})
		podMemoryUsageMins = append(podMemoryUsageMins, PodMemoryUsageMin{float64(0), ""})
		podNetworkUsageMins = append(podNetworkUsageMins, PodNetworkUsageMin{"Bps", 0, 0, ""})
		return PhysicalResources{podCPUUsageMins, podMemoryUsageMins, podNetworkUsageMins}
	}
}

func GetInfluxDBPod10mMetric(clusterName string, projectName string) PhysicalResources {

	var podCPUUsageMins []PodCPUUsageMin
	var podMemoryUsageMins []PodMemoryUsageMin
	var podNetworkUsageMins []PodNetworkUsageMin

	InitInfluxConfig()
	inf := NewInflux(InfluxConfig.Influx.Ip, InfluxConfig.Influx.Port, InfluxConfig.Influx.Username, InfluxConfig.Influx.Username)

	nowTime := time.Now().UTC() //.Add(time.Duration(offset) * time.Second)
	endTime := nowTime
	startTime := nowTime.Add(time.Duration(-12) * time.Minute)
	start := startTime.Format("2006-01-02T15:04") + ":00.0Z"
	end := endTime.Format("2006-01-02T15:04") + ":00.0Z"
	// start = "2020-12-16T10:21:00.0Z"
	// end = "2020-12-16T10:32:00.0Z"

	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// 	},
	// }

	q := client.Query{}
	// t := time.Now().Format(time.RFC3339)
	// select time, CPUUsageNanoCores as cpuUsage, MemoryUsageBytes as memoryUsage, NetworkRxBytes as Rx, NetworkTxBytes as Tx, pod  from Pods where time < '2020-12-16T10:31:00.0Z' and time > '2020-12-16T10:21:00.0Z' and cluster='cluster2' and namespace='ingress-nginx'
	query := "select time, CPUUsageNanoCores as cpuUsage, MemoryUsageBytes as memoryUsage, NetworkRxBytes as Rx, NetworkTxBytes as Tx, pod from Pods where time < '" + end + "' and time > '" + start + "' and cluster='" + clusterName + "' and namespace='" + projectName + "' order by time asc"

	q = client.NewQuery(query, "Metrics", "")
	response, err := inf.inClient.Query(q)

	if err != nil {
		fmt.Println("ERR : ", err)
	}

	// var queryResult []client.Result
	queryResult := response.Results[0]

	if len(queryResult.Series) == 0 {
		podCPUUsageMins = append(podCPUUsageMins, PodCPUUsageMin{math.Ceil(0), ""})
		podMemoryUsageMins = append(podMemoryUsageMins, PodMemoryUsageMin{math.Ceil(0), ""})
		podNetworkUsageMins = append(podNetworkUsageMins, PodNetworkUsageMin{"Bps", 0, 0, ""})
		result := PhysicalResources{podCPUUsageMins, podMemoryUsageMins, podNetworkUsageMins}
		return result
	}

	ser := queryResult.Series[0]
	_, offset := time.Now().Zone()

	type NtUsage struct {
		RxAvg float64 `json:"rx_avg"`
		TxAvg float64 `json:"tx_avg"`
		RxMin float64 `json:"rx_min"`
		TxMin float64 `json:"tx_min"`
	}

	type PodNt map[string]*NtUsage

	type ResourceInfo struct {
		Count  float64 `json:"count"`
		CPU    float64 `json:"cpu"`
		Memory float64 `json:"memory"`
		NtRx   float64 `json:"ntRx"`
		NtTx   float64 `json:"ntTx"`
		Pod    PodNt   `json:"pod_nt"`
	}

	type MetricByMin map[string]*ResourceInfo

	metrics := make(MetricByMin)
	preTimeHM := ""
	podlist := []string{}
	for i, value := range ser.Values {
		// 	[
		// 		"2020-12-16T04:22:02.112243605Z", [0][0]
		// 		"1177304n",  [0][1] cpuUsage
		// 		"170028Ki",  [0][2] memoryUsage
		// 		"926886456", [0][3] Rx
		// 		"279574980"  [0][4] Tx
		//		"pods-1"		 [0][5] pod
		//  ],
		// value := ser.Values[r][c].(string)
		timeValue := value[0].(string)
		ind := strings.Index(timeValue, ":")
		timeHM := timeValue[ind-2 : ind+3]
		timeHM = timeHM + ":00"
		t1, _ := time.Parse("15:04:05", timeHM)
		t1 = t1.Add(time.Duration(offset) * time.Second)

		timeHM = t1.Format("15:04:05")
		cpuValue, _ := strconv.ParseFloat(strings.Split(value[1].(string), "n")[0], 64)
		memValue, _ := strconv.ParseFloat(strings.Split(value[2].(string), "Ki")[0], 64)
		rxValue, _ := strconv.ParseFloat(value[3].(string), 64)
		txValue, _ := strconv.ParseFloat(value[4].(string), 64)
		podName, _ := value[5].(string)

		if i == 0 {
			preTimeHM = timeHM
			metrics[timeHM] = &ResourceInfo{}
			metrics[timeHM].Pod = make(PodNt)
		} else if timeHM != preTimeHM || len(ser.Values)-1 == i {
			metrics[preTimeHM].CPU = metrics[preTimeHM].CPU / metrics[preTimeHM].Count / 1000 / 1000 / 1000
			metrics[preTimeHM].Memory = metrics[preTimeHM].Memory / metrics[preTimeHM].Count / 1000

			for _, pod := range podlist {
				metrics[preTimeHM].NtRx += metrics[preTimeHM].Pod[pod].RxAvg / float64(len(podlist))
				metrics[preTimeHM].NtTx += metrics[preTimeHM].Pod[pod].TxAvg / float64(len(podlist))
			}

			metrics[timeHM] = &ResourceInfo{}
			metrics[timeHM].Pod = make(PodNt)
			preTimeHM = timeHM
			podlist = []string{}
		}

		metrics[timeHM].CPU = metrics[timeHM].CPU + cpuValue
		metrics[timeHM].Memory = metrics[timeHM].Memory + memValue

		if FindInStrArr(podlist, podName) {
			metrics[timeHM].Pod[podName].RxAvg = (rxValue - metrics[timeHM].Pod[podName].RxMin) / 60
			metrics[timeHM].Pod[podName].TxAvg = (txValue - metrics[timeHM].Pod[podName].TxMin) / 60
		} else {
			podlist = append(podlist, podName)
			metrics[timeHM].Pod[podName] = &NtUsage{}
			metrics[timeHM].Pod[podName].RxMin = rxValue
			metrics[timeHM].Pod[podName].TxMin = txValue
		}

		metrics[timeHM].Count++
	}

	for key, element := range metrics {
		podCPUUsageMins = append(podCPUUsageMins, PodCPUUsageMin{math.Ceil(element.CPU*1000) / 1000, key})
		podMemoryUsageMins = append(podMemoryUsageMins, PodMemoryUsageMin{math.Ceil(element.Memory*1000) / 1000, key})
		podNetworkUsageMins = append(podNetworkUsageMins, PodNetworkUsageMin{"Bps", int(element.NtRx), int(element.NtTx), key})
	}

	sort.Slice(podCPUUsageMins, func(i, j int) bool {
		return podCPUUsageMins[i].Time < podCPUUsageMins[j].Time
	})
	sort.Slice(podMemoryUsageMins, func(i, j int) bool {
		return podMemoryUsageMins[i].Time < podMemoryUsageMins[j].Time
	})
	sort.Slice(podNetworkUsageMins, func(i, j int) bool {
		return podNetworkUsageMins[i].Time < podNetworkUsageMins[j].Time
	})

	if len(podCPUUsageMins) > 10 {
		podCPUUsageMins = podCPUUsageMins[1 : len(podCPUUsageMins)-1]
		podMemoryUsageMins = podMemoryUsageMins[1 : len(podCPUUsageMins)-1]
		podNetworkUsageMins = podNetworkUsageMins[1 : len(podCPUUsageMins)-1]
	}

	result := PhysicalResources{podCPUUsageMins, podMemoryUsageMins, podNetworkUsageMins}
	return result
}

func GetInfluxPodTop5(clusterName string, projectName string) UsageTop5 {

	nowTime := time.Now().UTC() //.Add(time.Duration(offset) * time.Second)
	// startTime := nowTime.Add(time.Duration(-10) * time.Minute)
	startTime := nowTime.Add(time.Duration(-5) * time.Minute)
	start := startTime.Format("2006-01-02T15:04") + ":00.0Z"

	var usageTop5 UsageTop5

	InitInfluxConfig()
	inf := NewInflux(InfluxConfig.Influx.Ip, InfluxConfig.Influx.Port, InfluxConfig.Influx.Username, InfluxConfig.Influx.Password)

	q := client.Query{}

	query := "select time, last(CPUUsageNanoCores) as cpuUsage, MemoryUsageBytes as memoryUsage, namespace, cluster, pod  from Pods where cluster='" + clusterName + "' and namespace='" + projectName + "' and time > '" + start + "' group by pod"

	fmt.Println(query)
	q = client.NewQuery(query, "Metrics", "")
	response, _ := inf.inClient.Query(q)

	// var queryResult []client.Result
	// fmt.Println(response)
	if response == nil {
		usageTop5.CPU = []UsageType{}
		usageTop5.Memory = []UsageType{}

		result := usageTop5
		return result
	}

	queryResult := response.Results

	if len(queryResult[0].Series) == 0 {
		usageTop5.CPU = []UsageType{}
		usageTop5.Memory = []UsageType{}

		result := usageTop5
		return result
	}

	for _, qRes := range queryResult {
		for _, ser := range qRes.Series {
			for _, value := range ser.Values {
				cpuUsage := UsageType{}
				memUsage := UsageType{}
				podName := ser.Tags["pod"]
				cpuUsage.Name = podName
				// cpuUsage.Type = podName

				intCpu, _ := strconv.Atoi(strings.Split(value[1].(string), "n")[0])
				floatCpu := float64(intCpu) / 1000 / 1000 / 1000
				strCpu := fmt.Sprintf("%.5g", floatCpu) + " core"
				cpuUsage.Usage = strCpu
				// cpuUsage.Usage = strings.Split(value[1].(string), "n")[0]

				memUsage.Name = podName
				// memUsage.Type = podName
				intMem, _ := strconv.Atoi(strings.Split(value[2].(string), "Ki")[0])
				floatMem := float64(intMem) / 1000 / 1000 //Gi
				strMem := fmt.Sprintf("%.5g", floatMem) + " Gi"
				memUsage.Usage = strMem
				// memUsage.Usage = strings.Split(value[2].(string), "Ki")[0]

				usageTop5.CPU = append(usageTop5.CPU, cpuUsage)
				usageTop5.Memory = append(usageTop5.Memory, memUsage)
			}
		}
	}

	sort.Slice(usageTop5.CPU, func(i, j int) bool {
		a, _ := strconv.ParseFloat(strings.Split(usageTop5.CPU[i].Usage, " core")[0], 64)
		b, _ := strconv.ParseFloat(strings.Split(usageTop5.CPU[j].Usage, " core")[0], 64)
		return a > b
	})

	sort.Slice(usageTop5.Memory, func(i, j int) bool {
		a, _ := strconv.ParseFloat(strings.Split(usageTop5.Memory[i].Usage, " Gi")[0], 64)
		b, _ := strconv.ParseFloat(strings.Split(usageTop5.Memory[j].Usage, " Gi")[0], 64)
		return a > b
	})

	if len(usageTop5.CPU) > 5 {
		usageTop5.CPU = usageTop5.CPU[0:5]
	}

	if len(usageTop5.Memory) > 5 {
		usageTop5.Memory = usageTop5.Memory[0:5]
	}

	result := usageTop5
	return result
}

func reverseRank(data map[string]float64, top int) PairList {
	pl := make(PairList, len(data))

	if top > len(data) {
		top = len(data)
	}
	i := 0
	for k, v := range data {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl[:top]
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Usage < p[j].Usage }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func GetStringElement(nMap interface{}, keys []string) string {
	result := ""

	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})[0].(string)
				} else {
					result = childMap.(string)
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = "-"
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = "-"
					break
				}
			}
		}
	} else {
		result = "-"
	}
	return result
}

func GetIntElement(nMap interface{}, keys []string) int {
	result := 0
	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})[0].(int)
				} else {
					result = childMap.(int)
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = 0
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = 0
					break
				}
			}
		}
	} else {
		result = 0
	}
	return result
}

func GetFloat64Element(nMap interface{}, keys []string) float64 {
	var result float64 = 0.0
	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})[0].(float64)
				} else {
					result = childMap.(float64)
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = 0.0
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = 0.0
					break
				}
			}
		}
	} else {
		result = 0.0
	}
	return result
}

func GetInterfaceElement(nMap interface{}, keys []string) interface{} {
	var result interface{}
	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})[0]
				} else {
					result = childMap
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = nil
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = nil
					break
				}
			}
		}
	} else {
		result = nil
	}
	return result
}

func GetArrayElement(nMap interface{}, keys []string) []interface{} {
	var result []interface{}
	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})
				} else {
					result = childMap.([]interface{})
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = nil
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = nil
					break
				}
			}
		}
	} else {
		result = nil
	}
	return result
}

func GetJsonBody(rbody io.Reader) map[string]interface{} {
	bodyBytes, err := ioutil.ReadAll(rbody)

	var data map[string]interface{}

	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(bodyBytes), &data)
	return data
}
