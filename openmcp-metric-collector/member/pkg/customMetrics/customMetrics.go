package customMetrics

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/storage"
	"os"
	"strconv"
	"strings"
)

var prev_networkrxusage = make(map[string]int)
var prev_networktxusage = make(map[string]int)
var diff_networkrxusage = make(map[string]int)
var diff_networktxusage = make(map[string]int)

func AddToDeployCustomMetricServer(data *storage.Collection, token string, host string, cluster_client *kubernetes.Clientset) {
	fmt.Println("AddToDeployCustomMetricServer Called")
	podList := make([]storage.PodMetricsPoint, 0)
	for i := 0; i < len(data.Metricsbatchs); i++ {
		podList = append(podList, data.Metricsbatchs[i].Pods...)
	}

	rs, err := cluster_client.AppsV1().ReplicaSets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	fmt.Println("[List] ReplicaSets")
	if err == nil {
		for _, replicaset := range rs.Items {
			check_exist := 0
			sum_cpuusage := 0
			sum_memoryusage := 0
			sum_networkrxusage := 0
			sum_networktxusage := 0
			sum_fsusage := 0

			if podList != nil {
				for _, value := range podList {
					if value.Name != "" {
						if strings.HasPrefix(value.Name, replicaset.Name) {
							check_exist += 1

							if len(value.CPUUsageNanoCores.String()) > 1 {
								tmp_cpu, _ := strconv.Atoi(value.CPUUsageNanoCores.String()[:len(value.CPUUsageNanoCores.String())-1])
								sum_cpuusage += tmp_cpu
							}

							if len(value.MemoryUsageBytes.String()) > 1 {
								tmp_mem, _ := strconv.Atoi(value.MemoryUsageBytes.String()[:len(value.MemoryUsageBytes.String())-2])
								sum_memoryusage += tmp_mem
							}

							_, rxexists := prev_networkrxusage[value.Name]
							if !rxexists {
								prev_networkrxusage[value.Name] = 0
							}
							_, txexists := prev_networktxusage[value.Name]
							if !txexists {
								prev_networkrxusage[value.Name] = 0
							}
							_, rxexists2 := diff_networkrxusage[value.Name]
							if !rxexists2 {
								prev_networkrxusage[value.Name] = 0
							}
							_, txexists2 := diff_networktxusage[value.Name]
							if !txexists2 {
								prev_networkrxusage[value.Name] = 0
							}

							tmp_rx, _ := strconv.Atoi(value.NetworkRxBytes.String())
							diff_networkrx := tmp_rx - prev_networkrxusage[value.Name]
							//fmt.Println("[rx] ",value.Name, " - ", diff_networkrx)
							if prev_networkrxusage[value.Name] == 0 {
								//fmt.Println(".. 1 init")
								diff_networkrx = 0
							}else if tmp_rx == 0 || diff_networkrx == 0 {
								//fmt.Println(".. 2 not change")
								diff_networkrx = diff_networkrxusage[value.Name]
							}
							sum_networkrxusage += diff_networkrx
							prev_networkrxusage[value.Name] = tmp_rx
							diff_networkrxusage[value.Name] = diff_networkrx

							tmp_tx, _ := strconv.Atoi(value.NetworkTxBytes.String())
							diff_networktx := tmp_tx - prev_networktxusage[value.Name]
							//fmt.Println("[tx] ",value.Name, " - ", diff_networktx)
							if prev_networktxusage[value.Name] == 0 {
								//fmt.Println(".. 1 init")
								diff_networktx = 0
							}else if tmp_tx == 0 || diff_networktx == 0 {
								//fmt.Println(".. 2 not change")
								diff_networktx = diff_networktxusage[value.Name]
							}
							sum_networktxusage += diff_networktx
							prev_networktxusage[value.Name] = tmp_tx
							diff_networktxusage[value.Name] = diff_networktx

							if len(value.FsUsedBytes.String()) > 1 {
								tmp_fs, _ := strconv.Atoi(value.FsUsedBytes.String()[:len(value.FsUsedBytes.String())-2])
								sum_fsusage += tmp_fs
							}
						}

					} else {
						fmt.Println("err : value.Name nil")
					}
				}
			} else {
				fmt.Println("Fail : Cannot load podList")
			}

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}

			if check_exist > 0 {
				namespace := replicaset.Namespace
				name := replicaset.Name[:strings.LastIndexAny(replicaset.Name, "-")]
				fmt.Println("[",name,"/",namespace,"]")
				fmt.Println("--------------------------")
				fmt.Println("Post CpuUsage :", strconv.Itoa(sum_cpuusage/check_exist)+"n")
				PostData(host, token, client, namespace, name, "CpuUsage", strconv.Itoa(sum_cpuusage/check_exist)+"n")
				fmt.Println("--------------------------")
				fmt.Println("Post MemoryUsage :", strconv.Itoa(sum_memoryusage/check_exist)+"Ki")
				PostData(host, token, client, namespace, name, "MemoryUsage", strconv.Itoa(sum_memoryusage/check_exist)+"Ki")
				fmt.Println("--------------------------")
				fmt.Println("Post NetworkRxUsage :", strconv.Itoa(sum_networkrxusage/check_exist))
				PostData(host, token, client, namespace, name, "NetworkRxUsage", strconv.Itoa(sum_networkrxusage/check_exist))
				fmt.Println("--------------------------")
				fmt.Println("Post NetworkTxUsage :", strconv.Itoa(sum_networktxusage/check_exist))
				PostData(host, token, client, namespace, name, "NetworkTxUsage", strconv.Itoa(sum_networktxusage/check_exist))
				fmt.Println("--------------------------")
				fmt.Println("Post FsUsage :", strconv.Itoa(sum_fsusage/check_exist)+"Ki")
				PostData(host, token, client, namespace, name, "FsUsage", strconv.Itoa(sum_fsusage/check_exist)+"Ki")
				fmt.Println("--------------------------")
			}

		}
	} else {
		fmt.Println("Fail : Cannot load RS ", err)
	}
}

func AddToPodCustomMetricServer(data *storage.Collection, token string, host string) {
	fmt.Println("AddToPodCustomMetricServer Called")
	for i := 0; i < len(data.Metricsbatchs); i++ {
		podList := data.Metricsbatchs[i].Pods
		if podList != nil {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}
			for _, value := range podList {
				if value.Name != "" {
					namespace := value.Namespace
					name := value.Name

					fmt.Println("Post CpuUsage :", value.CPUUsageNanoCores.String())
					PostData(host, token, client, namespace, name, "CpuUsage", value.CPUUsageNanoCores.String())

					fmt.Println("Post MemoryUsage :", value.MemoryUsageBytes.String())
					PostData(host, token, client, namespace, name, "MemoryUsage", value.MemoryUsageBytes.String())

					fmt.Println("Post NetworkRxUsage :", value.NetworkRxBytes.String())
					PostData(host, token, client, namespace, name, "NetworkRxUsage", value.NetworkRxBytes.String())

					fmt.Println("Post NetworkTxUsage :", value.NetworkTxBytes.String())
					PostData(host, token, client, namespace, name, "NetworkTxUsage", value.NetworkTxBytes.String())

					fmt.Println("Post FsUsage :", value.FsUsedBytes.String())
					PostData(host, token, client, namespace, name, "FsUsage", value.FsUsedBytes.String())

				} else {
					fmt.Println("Fail : Cannot load resources")
				}
			}
		} else {
			fmt.Println("Fail : Cannot load Pod list")
		}
	}
}

func PostData(host string, token string, client *http.Client, resourceNamespace string, resourceName string, resourceMetricName string, resourceMetricValue string) {
	fmt.Println("PostData Called")
	apiserver := host
	baselink := "/api/v1/namespaces/custom-metrics/services/custom-metrics-apiserver:http/proxy/"
	basepath := "write-metrics"
	resourceKind := "pods"

	url := "" + apiserver + baselink + basepath + "/namespaces/" + resourceNamespace + "/" + resourceKind + "/" + resourceName + "/" + resourceMetricName
	buff := bytes.NewBufferString(resourceMetricValue)

	req, err := http.NewRequest("POST", os.ExpandEnv(url), buff)

	if err != nil {
		// handle err
		fmt.Println("Fail NewRequest")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", os.ExpandEnv("Bearer "+token))


	resp, err := client.Do(req)
	if err != nil {
		// handle err
		fmt.Println("Fail POST")
	} else {
	}
	defer resp.Body.Close()
}
