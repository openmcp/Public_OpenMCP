package customMetrics

import (
	"bytes"
	"crypto/tls"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"openmcp/openmcp/omcplog"

	//"bytes"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/storage"
	//url1 "net/url"
	"os"
	"strconv"
	"strings"
)

func AddToDeployCustomMetricServer(data *storage.Collection, token string, host string, cluster_client *kubernetes.Clientset) {
	omcplog.V(4).Info("AddToDeployCustomMetricServer Called")
	podList := make([]storage.PodMetricsPoint, 0)
	for i := 0; i < len(data.Matricsbatchs); i++ {
		podList = append(podList, data.Matricsbatchs[i].Pods...)
	}

	rs, err := cluster_client.AppsV1().ReplicaSets(metav1.NamespaceAll).List(metav1.ListOptions{})
	omcplog.V(2).Info("[List] ReplicaSets")
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
							//omcplog.V(0).Info(value.Name, "  ", replicaset.Name)
							check_exist += 1

							tmp_cpu, _ := strconv.Atoi(value.CPUUsageNanoCores.String()[:len(value.CPUUsageNanoCores.String())-1])
							sum_cpuusage += tmp_cpu

							tmp_mem, _ := strconv.Atoi(value.MemoryUsageBytes.String()[:len(value.MemoryUsageBytes.String())-2])
							sum_memoryusage += tmp_mem

							tmp_rx, _ := strconv.Atoi(value.NetworkRxBytes.String())
							sum_networkrxusage += tmp_rx

							tmp_tx, _ := strconv.Atoi(value.NetworkTxBytes.String())
							sum_networktxusage += tmp_tx

							tmp_fs, _ := strconv.Atoi(value.FsUsedBytes.String()[:len(value.FsUsedBytes.String())-2])
							sum_fsusage += tmp_fs
						}

					} else {
						omcplog.V(0).Info("err : value.Name nil")
					}
				}
			} else {
				omcplog.V(0).Info("Fail : Cannot load podList")
			}

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}

			if check_exist > 0 {
				namespace := replicaset.Namespace
				name := replicaset.Name[:strings.LastIndexAny(replicaset.Name, "-")]
				//omcplog.V(0).Info(name, " ",sum_cpuusage," ",sum_cpuusage/check_exist, " ", strconv.Itoa(sum_cpuusage/check_exist))

				omcplog.V(3).Info("Post CpuUsage :", strconv.Itoa(sum_cpuusage/check_exist)+"n")
				PostData(host, token, client, namespace, name, "CpuUsage", strconv.Itoa(sum_cpuusage/check_exist)+"n")

				omcplog.V(3).Info("Post MemoryUsage :", strconv.Itoa(sum_memoryusage/check_exist)+"Ki")
				PostData(host, token, client, namespace, name, "MemoryUsage", strconv.Itoa(sum_memoryusage/check_exist)+"Ki")

				omcplog.V(3).Info("Post NetworkRxUsage :", strconv.Itoa(sum_networkrxusage/check_exist))
				PostData(host, token, client, namespace, name, "NetworkRxUsage", strconv.Itoa(sum_networkrxusage/check_exist))

				omcplog.V(3).Info("Post NetworkTxUsage :", strconv.Itoa(sum_networktxusage/check_exist))
				PostData(host, token, client, namespace, name, "NetworkTxUsage", strconv.Itoa(sum_networktxusage/check_exist))

				omcplog.V(3).Info("Post FsUsage :", strconv.Itoa(sum_fsusage/check_exist)+"Ki")
				PostData(host, token, client, namespace, name, "FsUsage", strconv.Itoa(sum_fsusage/check_exist)+"Ki")
			}

		}
	} else {
		omcplog.V(0).Info("Fail : Cannot load RS ", err)
	}
}

func AddToPodCustomMetricServer(data *storage.Collection, token string, host string) {
	omcplog.V(4).Info("AddToPodCustomMetricServer Called")
	for i := 0; i < len(data.Matricsbatchs); i++ {
		podList := data.Matricsbatchs[i].Pods
		if podList != nil {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}
			for _, value := range podList {
				if value.Name != "" {
					namespace := value.Namespace
					name := value.Name

					omcplog.V(3).Info("Post CpuUsage :", value.CPUUsageNanoCores.String())
					PostData(host, token, client, namespace, name, "CpuUsage", value.CPUUsageNanoCores.String())

					omcplog.V(3).Info("Post MemoryUsage :", value.MemoryUsageBytes.String())
					PostData(host, token, client, namespace, name, "MemoryUsage", value.MemoryUsageBytes.String())

					omcplog.V(3).Info("Post NetworkRxUsage :", value.NetworkRxBytes.String())
					PostData(host, token, client, namespace, name, "NetworkRxUsage", value.NetworkRxBytes.String())

					omcplog.V(3).Info("Post NetworkTxUsage :", value.NetworkTxBytes.String())
					PostData(host, token, client, namespace, name, "NetworkTxUsage", value.NetworkTxBytes.String())

					omcplog.V(3).Info("Post FsUsage :", value.FsUsedBytes.String())
					PostData(host, token, client, namespace, name, "FsUsage", value.FsUsedBytes.String())

				} else {
					omcplog.V(0).Info("Fail : Cannot load resources")
				}
			}
		} else {
			omcplog.V(0).Info("Fail : Cannot load Pod list")
		}
	}
}

func PostData(host string, token string, client *http.Client, resourceNamespace string, resourceName string, resourceMetricName string, resourceMetricValue string) {
	omcplog.V(4).Info("PostData Called")
	apiserver := host
	baselink := "/api/v1/namespaces/custom-metrics/services/custom-metrics-apiserver:http/proxy/"
	basepath := "write-metrics"
	resourceKind := "pods"
	//omcplog.V(0).Info(resourceMetricValue)
	//valueString := strconv.FormatFloat(resourceMetricValue, 'e', 4, 64)

	url := "" + apiserver + baselink + basepath + "/namespaces/" + resourceNamespace + "/" + resourceKind + "/" + resourceName + "/" + resourceMetricName
	buff := bytes.NewBufferString(resourceMetricValue)

	//omcplog.V(0).Info("value : ",buff)

	/*data := url1.Values{}
	data.Set("metrics", "111111")
	omcplog.V(0).Info("value : ",strings.NewReader(data.Encode()))*/

	req, err := http.NewRequest("POST", os.ExpandEnv(url), buff)

	if err != nil {
		// handle err
		omcplog.V(0).Info("Fail NewRequest")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", os.ExpandEnv("Bearer "+token))

	//omcplog.V(0).Info("req", req)

	resp, err := client.Do(req)
	if err != nil {
		// handle err
		omcplog.V(0).Info("Fail POST")
	} else {
		//omcplog.V(0).Info("Success POST")
	}
	defer resp.Body.Close()
}
