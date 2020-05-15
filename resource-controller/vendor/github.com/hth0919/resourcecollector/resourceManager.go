package CRM

import (
	"crypto/tls"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"strconv"
	"strings"
)

func initmeticmap() map[string]float64 {
	mm := map[string]float64{
		"cpu_cfs_periods_total":                  0,
		"cpu_cfs_throttled_periods_total":        0,
		"cpu_cfs_throttled_seconds_total":        0,
		"cpu_load_average_10s":                   0,
		"cpu_system_seconds_total":               0,
		"cpu_usage_seconds_total":                0,
		"cpu_user_seconds_total":                 0,
		"fs_inodes_free":                         0,
		"fs_inodes_total":                        0,
		"fs_io_current":                          0,
		"fs_io_time_seconds_total":               0,
		"fs_io_time_weighted_seconds_total":      0,
		"fs_limit_bytes":                         0,
		"fs_read_seconds_total":                  0,
		"fs_reads_bytes_total":                   0,
		"fs_reads_merged_total":                  0,
		"fs_reads_total":                         0,
		"fs_sector_reads_total":                  0,
		"fs_sector_writes_total":                 0,
		"fs_usage_bytes":                         0,
		"fs_write_seconds_total":                 0,
		"fs_writes_bytes_total":                  0,
		"fs_writes_merged_total":                 0,
		"fs_writes_total":                        0,
		"last_seen":                              0,
		"memory_cache":                           0,
		"memory_failcnt":                         0,
		"memory_failures_total":                  0,
		"memory_mapped_file":                     0,
		"memory_max_usage_bytes":                 0,
		"memory_rss":                             0,
		"memory_swap":                            0,
		"memory_usage_bytes":                     0,
		"memory_working_set_bytes":               0,
		"network_receive_bytes_total":            0,
		"network_receive_errors_total":           0,
		"network_receive_packets_dropped_total":  0,
		"network_receive_packets_total":          0,
		"sockets":0,
		"network_transmit_bytes_total":           0,
		"network_transmit_errors_total":          0,
		"network_transmit_packets_dropped_total": 0,
		"network_transmit_packets_total":         0,
		"scrape_error 0":                         0,
		"spec_cpu_period":                        0,
		"spec_cpu_quota":                         0,
		"spec_cpu_shares":                        0,
		"spec_memory_limit_bytes":                0,
		"spec_memory_reservation_limit_bytes":    0,
		"spec_memory_swap_limit_bytes":           0,
		"start_time_seconds":                     0,
		"tasks_state":                            0,
		"machine_cpu_cores":                      0,
		"machine_memory_bytes":                   0,
	}
	return mm
}
var metricValue []string
var cs *ClientSet = &ClientSet{clientSet:nil}
var strbuf []rune
var ch rune = -1
var tempmericvalue string
var temppodname string
var tempvalue float64
var podflag bool = false
var quatesflag bool = false
var valueflag bool= false
var sharpflag bool= false

func (ci *ClusterInfo)NewClusterClient(masterUri string) {
	ci.ClusterMetricSum =initmeticmap()

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	ci.Host = config.Host
	cs.clientSet, err = kubernetes.NewForConfig(config)
	ci.Pods = make([]string,0,1)
	if err != nil {
		panic(err.Error())
	}
	secrets, _ := cs.clientSet.CoreV1().Secrets(metav1.NamespaceAll).List(metav1.ListOptions{})
	P,_:=cs.clientSet.CoreV1().Pods("rescollect").List(metav1.ListOptions{})
	ci.Pods = append(ci.Pods,P.Items[0].Name)
	if err != nil {
		panic(err.Error())
	}
	ci.AdminToken = string(secrets.Items[0].Data["token"])
	ci.AdminToken = strings.TrimSpace(ci.AdminToken)
	for i:=0;i<len(secrets.Items);i++ {
		if strings.Contains(secrets.Items[i].Name, "monitoring") {
			ci.AdminToken = string(secrets.Items[i].Data["token"])
			ci.AdminToken = strings.TrimSpace(ci.AdminToken)
		}
	}

}

func (ci *ClusterInfo)NodeListInit() {
	var ni *NodeInfo

	metricValue = []string{"cpu_cfs_periods_total", "cpu_cfs_throttled_periods_total", "cpu_cfs_throttled_seconds_total", "cpu_load_average_10s", "cpu_system_seconds_total", "cpu_usage_seconds_total", "cpu_user_seconds_total", "fs_inodes_free", "fs_inodes_total", "fs_io_current", "fs_io_time_seconds_total", "fs_io_time_weighted_seconds_total", "fs_limit_bytes", "fs_read_seconds_total", "fs_reads_bytes_total", "fs_reads_merged_total", "fs_reads_total", "fs_sector_reads_total", "fs_sector_writes_total", "fs_usage_bytes", "fs_write_seconds_total", "fs_writes_bytes_total", "fs_writes_merged_total", "fs_writes_total", "last_seen", "memory_cache", "memory_failcnt", "memory_failures_total", "memory_mapped_file", "memory_max_usage_bytes", "memory_rss", "memory_swap", "memory_usage_bytes", "memory_working_set_bytes", "network_receive_bytes_total", "network_receive_errors_total", "network_receive_packets_dropped_total", "network_receive_packets_total", "sockets", "network_transmit_bytes_total", "network_transmit_errors_total", "network_transmit_packets_dropped_total", "network_transmit_packets_total", "scrape_error", "spec_cpu_period", "spec_cpu_quota", "spec_cpu_shares", "spec_memory_limit_bytes", "spec_memory_reservation_limit_bytes", "spec_memory_swap_limit_bytes", "start_time_seconds", "tasks_state", "machine_cpu_cores", "machine_memory_bytes"}
	nodes, err :=cs.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	ci.NodeList = make([]*NodeInfo,0,len(nodes.Items))
	for i:=0;i<len(nodes.Items);i++ {
		ni = &NodeInfo{
			NodeName:        "",
			PodList:         []*PodInfo{},
			NodeMetricSum:   map[string]float64{},
			NodeAllocatable: map[string]int64{},
			NodeCapacity:    map[string]int64{},
			GeoInfo:		 map[string]string{},
			CpuCores:        0,
			MemoryTotal:     0,
			ScrapeError:     0,
		}
		ni.NodeMetricSum = initmeticmap()
		ci.NodeList = append(ci.NodeList, ni)
	}
	for i:=0;i<len(nodes.Items);i++ {
		ci.NodeList[i].NodeName = nodes.Items[i].Name
		ci.NodeList[i].NodeCapacity["CPU"] = nodes.Items[i].Status.Capacity.Cpu().Value()
		ci.NodeList[i].NodeAllocatable["CPU"] = nodes.Items[i].Status.Allocatable.Cpu().Value()
		ci.NodeList[i].NodeCapacity["Memory"] = nodes.Items[i].Status.Capacity.Memory().Value()
		ci.NodeList[i].NodeAllocatable["Memory"] = nodes.Items[i].Status.Allocatable.Memory().Value()
		ci.NodeList[i].NodeCapacity["EphemeralStorage"] = nodes.Items[i].Status.Capacity.StorageEphemeral().Value()
		ci.NodeList[i].NodeAllocatable["EphemeralStorage"] = nodes.Items[i].Status.Allocatable.StorageEphemeral().Value()
		ci.NodeList[i].GeoInfo["Region"] = nodes.Items[i].Labels["failure-domain.beta.kubernetes.io/region"]
		ci.NodeList[i].GeoInfo["Zone"] = nodes.Items[i].Labels["failure-domain.beta.kubernetes.io/zone"]
		responseTokenizer(ci,ci.NodeList[i].NodeName,i)
	}

}

func (ci *ClusterInfo)CalculateClusterMetricSum() map[string]float64{
	var sum float64
	sum = 0

	for i:=0;i<len(metricValue);i++ {
		sum = 0
		for j := 0;j<len(ci.NodeList);j++ {
			if metricValue[i] == "scrape_error" {
				continue
			} else if metricValue[i] == "machine_cpu_cores"{
				continue
			}else if metricValue[i] == "machine_memory_bytes"{
				continue
			}
			sum += ci.NodeList[j].NodeMetricSum[metricValue[i]]
			ci.ClusterMetricSum[metricValue[i]] = sum
		}
	}
	return ci.ClusterMetricSum
}

func (ci *ClusterInfo)CalculateNodeMetricSum(index int) map[string]float64 {
	var sum float64
	sum = 0
	for i:=0;i<len(metricValue);i++ {
		sum = 0
		for j := 0;j<len(ci.NodeList[index].PodList);j++ {
			if metricValue[i] == "scrape_error" {
				continue
			} else if metricValue[i] == "machine_cpu_cores"{
				continue
			}else if metricValue[i] == "machine_memory_bytes"{
				continue
			}
			sum += ci.NodeList[index].PodList[j].PodMetrics[metricValue[i]]
			ci.NodeList[index].NodeMetricSum[metricValue[i]] = sum
		}
	}
	return ci.NodeList[index].NodeMetricSum
}

func FindOrMakePodInfo(name string, pil []*PodInfo) (*PodInfo,int, ) {
	var result PodInfo
	result = PodInfo{
		PodName: name,
		PodNamespace: "",
		PodMetrics: initmeticmap(),
	}

	for i:=0;i<len(pil);i++ {
		if pil[i].PodName == name {
			return pil[i], i
		}
	}
	po, err := cs.clientSet.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for i := 0;i<len(po.Items);i++ {
		if name == po.Items[i].Name {
			ns := po.Items[i].Namespace
			result.PodNamespace = ns
		}
	}
	pil = append(pil, &result)
	return &result, -1
}
func parser(ci *ClusterInfo, str string, indexnum int) []*PodInfo{
	var pil []*PodInfo
	var ch rune
	var zeroslice []rune
	zeroslice = []rune{}
	pil = make([]*PodInfo,0,0)
	for i := 0; i<len(str);i++ {
		ch = rune(str[i])
		if ch == '{'{
			if sharpflag {
				strbuf = zeroslice
				continue
			}
			yylex(string(strbuf))
			strbuf = zeroslice
		} else if ch == ',' {
			if sharpflag {
				strbuf = zeroslice
				continue
			}
			strbuf = zeroslice
		} else if ch == '=' {
			if sharpflag {
				strbuf = zeroslice
				continue
			}
			yylex(string(strbuf))
			strbuf = zeroslice
		} else if ch == '"' {
			if sharpflag {
				strbuf = zeroslice
				continue
			}
			if quatesflag {
				quatesflag = false
				if podflag {
					temppodname = string(strbuf)

					podflag = false
					strbuf = zeroslice
				}
			} else {
				quatesflag = true
			}
		} else if ch == ' ' {
			if sharpflag {
				strbuf = zeroslice
				continue
			}
			yylex(string(strbuf))

			if valueflag {
				tempvalue, _ = strconv.ParseFloat(string(strbuf),64)
				strbuf = zeroslice
				valueflag = false
			} else {
				strbuf = zeroslice
				valueflag = true
			}
			strbuf = zeroslice

		} else if ch == '\n' {

			if sharpflag {
				strbuf = zeroslice
				sharpflag = false
				continue
			}
			if tempmericvalue == "machine_cpu_cores"{
				tempvalue, _ = strconv.ParseFloat(string(strbuf),64)
				strbuf = zeroslice
				ci.NodeList[indexnum].CpuCores = tempvalue
				valueflag = false
				continue
			} else if tempmericvalue == "machine_memory_bytes"{
				tempvalue, _ = strconv.ParseFloat(string(strbuf),64)
				strbuf = zeroslice
				ci.NodeList[indexnum].MemoryTotal = tempvalue
				valueflag = false
				continue
			}else if tempmericvalue == "scrape_error" {
				tempvalue, _ = strconv.ParseFloat(string(strbuf),64)
				strbuf = zeroslice
				ci.NodeList[indexnum].ScrapeError = tempvalue
				valueflag = false
				continue
			}
			strbuf = zeroslice
			pi, index := FindOrMakePodInfo(temppodname, pil)

			temp := 0.0
			temp = pi.PodMetrics[tempmericvalue]
			temp += tempvalue
			pi.PodMetrics[tempmericvalue] = temp

			if index == -1 {
				pil = append(pil, pi)
			} else {
				pil[index] = pi
			}
			valueflag = false
		} else if ch == '#' {
			sharpflag = true
		} else if ch == '_' {
			if sharpflag {
				strbuf = zeroslice
				continue
			}
			if string(strbuf) == "container" {
				strbuf = zeroslice
			} else {strbuf = append(strbuf, ch)}
		} else{
			if sharpflag {
				strbuf = zeroslice
				continue
			}
			strbuf = append(strbuf, ch)
		}

	}
	return pil
}
func yylex(buff string) {
	switch buff {
	case "cadvisor_version_info":
		sharpflag = true
	case "cpu_cfs_periods_total":
		tempmericvalue = buff
	case "cpu_cfs_throttled_periods_total":
		tempmericvalue = buff
	case "cpu_cfs_throttled_seconds_total":
		tempmericvalue = buff
	case "cpu_load_average_10s":
		tempmericvalue = buff
	case "cpu_system_seconds_total":
		tempmericvalue = buff
	case "cpu_usage_seconds_total":
		tempmericvalue = buff
	case "cpu_user_seconds_total":
		tempmericvalue = buff
	case "fs_inodes_free":
		tempmericvalue = buff
	case "fs_inodes_total":tempmericvalue = buff
	case "fs_io_current":tempmericvalue = buff
	case "fs_io_time_seconds_total":tempmericvalue = buff
	case "fs_io_time_weighted_seconds_total":tempmericvalue = buff
	case "fs_limit_bytes":tempmericvalue = buff
	case "fs_read_seconds_total":tempmericvalue = buff
	case "fs_reads_bytes_total":tempmericvalue = buff
	case "fs_reads_merged_total":tempmericvalue = buff
	case "fs_reads_total":tempmericvalue = buff
	case "fs_sector_reads_total":tempmericvalue = buff
	case "fs_sector_writes_total":tempmericvalue = buff
	case "fs_usage_bytes":tempmericvalue = buff
	case "fs_write_seconds_total":tempmericvalue = buff
	case "fs_writes_bytes_total":tempmericvalue = buff
	case "fs_writes_merged_total":tempmericvalue = buff
	case "fs_writes_total":tempmericvalue = buff
	case "last_seen":tempmericvalue = buff
	case "memory_cache":tempmericvalue = buff
	case "memory_failcnt":tempmericvalue = buff
	case "memory_failures_total":tempmericvalue = buff
	case "memory_mapped_file":tempmericvalue = buff
	case "memory_max_usage_bytes":tempmericvalue = buff
	case "memory_rss":tempmericvalue = buff
	case "memory_swap":tempmericvalue = buff
	case "memory_usage_bytes":tempmericvalue = buff
	case "memory_working_set_bytes":tempmericvalue = buff
	case "network_receive_bytes_total":tempmericvalue = buff
	case "network_receive_errors_total":tempmericvalue = buff
	case "network_receive_packets_dropped_total":tempmericvalue = buff
	case "network_receive_packets_total":tempmericvalue = buff
	case "network_transmit_bytes_total":tempmericvalue = buff
	case "network_transmit_errors_total":tempmericvalue = buff
	case "network_transmit_packets_dropped_total":tempmericvalue = buff
	case "network_transmit_packets_total":tempmericvalue = buff
	case "scrape_error 0":tempmericvalue = buff
	case "spec_cpu_period":tempmericvalue = buff
	case "spec_cpu_quota":tempmericvalue = buff
	case "spec_cpu_shares":tempmericvalue = buff
	case "spec_memory_limit_bytes":tempmericvalue = buff
	case "spec_memory_reservation_limit_bytes":tempmericvalue = buff
	case "spec_memory_swap_limit_bytes":tempmericvalue = buff
	case "start_time_seconds":tempmericvalue = buff
	case "tasks_state":tempmericvalue = buff
	case "sockets":tempmericvalue = buff
	case "machine_cpu_cores":tempmericvalue = buff
	case "machine_memory_bytes":tempmericvalue = buff
	case "scrape_error":tempmericvalue = buff
	case "pod": podflag = true
	}
}



func responseTokenizer(ci *ClusterInfo, nodename string,indexnum int) {

	url := "" + ci.Host + "/api/v1/nodes/" + nodename + "/proxy/metrics/cadvisor"
	url = strings.TrimSpace(url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	//필요시 헤더 추가 가능
	req.Header.Add("Authorization", "Bearer " + ci.AdminToken)

	// Client객체에서 Request 실행
	client := &http.Client{Transport:tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()


	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	filename := nodename + ".txt"
	err = ioutil.WriteFile(filename, bytes, 0)
	if err != nil {
		panic(err)
	}
	ci.NodeList[indexnum].PodList = parser(ci, str, indexnum)
	strbuf = []rune{}
	ch = -1
	tempmericvalue =""
	temppodname =""
	tempvalue =0
	podflag  = false
	quatesflag = false
	valueflag = false
	sharpflag = false
}