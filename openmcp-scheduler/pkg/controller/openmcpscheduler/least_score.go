package openmcpscheduler

import (
	"k8s.io/klog"
)

//포드의 요청된 자원이 적은 노드를 선호함
//즉, 노드에 사용된 자원이 적을 수록 높은 점수를 부여함
//여기서, 포드의 요청된 자원은 노드에 설치된 포드들의 할당 자원 크기를 의미함
func getLeastScore(cpu float64, mem float64, c_List []string) map[string]float64 {
	var leastScore map[string]float64 = make(map[string]float64)
	klog.Infof("[SUJUNE] Start getLeastScore")
	for _, c_name := range c_List {
		temp := clusterResource[c_name]

		if temp.cpu_idle < 0 && temp.mem_idle < 0 {
			continue
		}

		if cpu > temp.cpu_idle && mem > temp.mem_idle {
			continue
		}

		cluster_Score := temp.cpu_idle / temp.cpu_total
		cluster_Score += temp.mem_idle / temp.mem_total
		// Resource 종류가 CPU, Memory 이므로 Weight = 2
		cluster_Score = cluster_Score / 2 //(Weight = 2)

		leastScore[c_name] = cluster_Score * 10
	}
	return leastScore
}
