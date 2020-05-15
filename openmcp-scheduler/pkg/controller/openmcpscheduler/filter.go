package openmcpscheduler

import (
	"k8s.io/klog"
)

func NewFiltering(Request_cpu float64, Request_mem float64, c_List []string, c_Resource map[string]node_Info) []string {
	klog.Infof("[SUJUNE] Start Filtering!!!")
	//TODO
	//여러 Filtering Policy 설계 필요

	//1. Pod 요구 자원을 할당할 수 있는 노드만 필터링
	filterList := resourceFiltering(Request_cpu, Request_mem, c_List, c_Resource)

	//TODO
	//추가 알고리즘 넣자...
	//2. ??
	//첫번째 필터링의 결과물인 filterList에 대해서 다음 알고리즘을 적용
	//필요한 파라미터가 있다면 알고리즘 특성에 맞춰서 설계할 것
	//filterList = NewFiltering(filterList)

	return filterList
}

//현재 OpenMCP의 스케줄러는 클러스터 단위로 데이터를 기반으로 계산하므로,
//Node 단위로 Pod를 더 이상 배포할 수 없어도배포가 되는 문제가 있음
//이러한 문제를 해결하기 위해서, Node 단위로 Pod를 배포할 수 있는지 클러스터를 검사하는 함수
func resourceFiltering(Request_cpu float64, Request_mem float64, c_List []string, c_Resource map[string]node_Info) []string {
	klog.Infof("[SUJUNE] Start resourceFiltering!!!")
	posClusterList := []string{}
	for _, c_name := range c_List {
		cluster_Resource, _ := c_Resource[c_name]
		for key1, value1 := range cluster_Resource.resource {
			klog.Infof("[SUJUNE] %s's %s Node Resource Check", c_name, key1)
			klog.Infof("[SUJUNE] Request cpu = %f", Request_cpu)
			klog.Infof("[SUJUNE] Request mem = %f", Request_mem)
			klog.Infof("[SUJUNE] cpu = %f", value1.cpu_idle)
			klog.Infof("[SUJUNE] mem = %f", value1.mem_idle)
			//Pod 1개 단위로 계산하므로
			//요구하는 Pod의 cpu, mem 값을 충족하는 Node가 1개라도 있으면 해당 cluster는 배포 가능
			if value1.cpu_idle-Request_cpu >= 0 && value1.mem_idle-Request_mem >= 0 {
				//배포 가능한 cluster를 리스트에 추가
				posClusterList = append(posClusterList, c_name)
				//배포 가능함을 확인 했으므로, 해당 클러스터는 더 이상 검사 필요 없음
				break
			}
		}
	}
	return posClusterList
}
