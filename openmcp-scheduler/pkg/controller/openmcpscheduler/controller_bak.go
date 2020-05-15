/*
Copyright 2018 The Multicluster-Controller Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openmcpscheduler // import "admiralty.io/multicluster-controller/examples/openmcpscheduler/pkg/controller/openmcpscheduler"

import (
	"admiralty.io/multicluster-controller/pkg/reference"
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"math/rand"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	"sort"
	"time"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"openmcp-scheduler/pkg/apis"
	ketiv1alpha1 "openmcp-scheduler/pkg/apis/keti/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"
	kubesource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"

	"math"
	"strconv"
	"strings"
)

type ClusterManager struct {
	Fed_namespace   string
	Host_config     *rest.Config
	Host_client     genericclient.Client
	Cluster_list    *fedv1b1.KubeFedClusterList
	Cluster_configs map[string]*rest.Config
	Cluster_clients map[string]genericclient.Client
}

type resource_Info struct {
	cpu_idle  float64
	cpu_total float64
	mem_idle  float64
	mem_total float64
}

type node_Info struct {
	name     string
	resource map[string]node_Resource
}

type node_Resource struct {
	cpu_idle float64
	mem_idle float64
}

type score_Info map[string]float64

var (
	nodeLevelResource map[string]node_Info
	clusterResource   map[string]resource_Info
	posClusterList    []string
	cluster_clientset map[string]*kubernetes.Clientset
	allocatedResource map[string]Resource
	availableNodes    []Node
)

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string) (*controller.Controller, error) {
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}

	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients = append(ghostclients, ghostclient)
	}

	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPscheduler{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	for _, ghost := range ghosts {
		fmt.Printf("%T, %s\n", ghost, ghost.GetClusterName())
		if err := co.WatchResourceReconcileController(ghost, &appsv1.Deployment{}, controller.WatchOptions{}); err != nil {
			return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
		}
	}
	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

var i int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	i += 1
	fmt.Println("********* [", i, "] *********")
	fmt.Println(req.Context, " / ", req.Namespace, " / ", req.Name)
	cm := NewClusterManager()

	// Fetch the OpenMCPDeployment instance
	instance := &ketiv1alpha1.OpenMCPscheduler{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			// ...TODO: multicluster garbage collector
			// Until then...
			err := cm.DeleteDeployments(req.NamespacedName)
			return reconcile.Result{}, err
		}
		fmt.Println("Error1")
		return reconcile.Result{}, err
	}
	if instance.Status.ClusterMaps == nil {
		fmt.Println("Scheduling Start & Create Deployments")

		//Get replicas
		replicas := instance.Spec.Replicas
		//Get policy
		policy := instance.Spec.Policy
		//Get Request_Resource -> CPU, Memory
		cpu_unit := instance.Spec.Template.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()
		quantity_change_cpu := kubesource.MustParse(instance.Spec.Template.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
		quantity_change_mem := kubesource.MustParse(instance.Spec.Template.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
		temp_cpu, _ := quantity_change_cpu.AsInt64()
		temp_mem, _ := quantity_change_mem.AsInt64()

		Request_cpu := float64(temp_cpu)
		Request_mem := float64(temp_mem) / 1024 / 1024 / 1024

		//CPU의 Core 단위가 m 일 경우 m단위로 들어갈 경우, float64 자료형으로 변경해주는 함수가 없는 걸로 판단
		if strings.Contains(cpu_unit, "m") {
			Request_cpu = milliTocoreCPU(cpu_unit)
		}

		//failSelected = 더 이상 배포할 수 있는 클러스터가 없는 경우, 몇개의 Pod가 배포 실패 했는지에 대한 Pod 수
		cluster_replicas_map, failSelected := cm.Scheduling(replicas, Request_cpu, Request_mem, policy)
		klog.Infof("[CHECK] check!!!! %v", cluster_replicas_map)
		//TODO
		//2020-01-13
		// failSelected가 1개 이상일 때, 예외처리 동작 설계 해야함
		klog.Infof("[Result] Can't deploy pod : failSelected = %d", failSelected)
		klog.Infof("[Result] Not have cluster that remaining resources")
		//Deployment를 생성하는 부분..
		for _, cluster := range cm.Cluster_list.Items {

			if cluster_replicas_map[cluster.Name] == 0 {
				continue
			}
			found := &appsv1.Deployment{}
			cluster_client := cm.Cluster_clients[cluster.Name]

			err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-deploy")
			if err != nil && errors.IsNotFound(err) {
				replica := cluster_replicas_map[cluster.Name]
				fmt.Println("Cluster '"+cluster.Name+"' Deployed (", replica, " / ", replicas, ")")
				dep := r.deploymentForOpenMCPscheduler(req, instance, replica)
				err = cluster_client.Create(context.Background(), dep)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
		instance.Status.ClusterMaps = cluster_replicas_map
		instance.Status.Replicas = replicas

		err := r.live.Status().Update(context.TODO(), instance)
		if err != nil {
			fmt.Println("Failed to update instance status", err)
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil

	} else if instance.Status.Replicas != instance.Spec.Replicas {
		fmt.Println("Change Spec Replicas ! ReScheduling Start & Update Deployment")
		cluster_replicas_map := cm.ReScheduling(instance.Spec.Replicas, instance.Status.Replicas, instance.Status.ClusterMaps)

		for _, cluster := range cm.Cluster_list.Items {
			update_replica := cluster_replicas_map[cluster.Name]
			cluster_client := cm.Cluster_clients[cluster.Name]

			dep := r.deploymentForOpenMCPscheduler(req, instance, update_replica)

			found := &appsv1.Deployment{}
			err := cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-deploy")
			if err != nil && errors.IsNotFound(err) {
				// Not Exist Deployment.
				if update_replica != 0 {
					// Create !
					err = cluster_client.Create(context.Background(), dep)
					if err != nil {
						return reconcile.Result{}, err
					}
				}

			} else if err != nil {
				return reconcile.Result{}, err
			} else {
				// Already Exist Deployment.
				if update_replica == 0 {
					// Delete !
					dep := &appsv1.Deployment{}
					err = cluster_client.Delete(context.Background(), dep, req.Namespace, req.Name+"-deploy")

					if err != nil {
						return reconcile.Result{}, err
					}
				} else {
					// Update !
					err = cluster_client.Update(context.TODO(), dep)
					if err != nil {
						return reconcile.Result{}, err
					}

				}

			}

		}
		instance.Status.ClusterMaps = cluster_replicas_map
		instance.Status.Replicas = instance.Spec.Replicas
		err := r.live.Status().Update(context.TODO(), instance)
		if err != nil {
			fmt.Println("Failed to update instance status", err)
			return reconcile.Result{}, err
		}

	} else {
		// Check Deployment in cluster
		for k, v := range instance.Status.ClusterMaps {
			cluster_name := k
			replica := v

			if replica == 0 {
				continue
			}
			found := &appsv1.Deployment{}
			cluster_client := cm.Cluster_clients[cluster_name]
			err = cluster_client.Get(context.TODO(), found, instance.Namespace, instance.Name+"-deploy")
			//		fmt.Println("err = ", err, " found = ", found, " Namespace = ", instance.Namespace)
			if err != nil && errors.IsNotFound(err) {
				// Delete Deployment Detected
				fmt.Println("Cluster '"+cluster_name+"' ReDeployed => ", replica)
				dep := r.deploymentForOpenMCPscheduler(req, instance, replica)
				err = cluster_client.Create(context.Background(), dep)
				if err != nil {
					return reconcile.Result{}, err
				}

			}

		}

	}

	return reconcile.Result{}, nil // err
}

func (r *reconciler) deploymentForOpenMCPscheduler(req reconcile.Request, m *ketiv1alpha1.OpenMCPscheduler, replica int32) *appsv1.Deployment {

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-deploy",
			Namespace: m.Namespace,
		},
		Spec: m.Spec.Template.Spec,
	}
	dep.Spec.Replicas = &replica

	reference.SetMulticlusterControllerReference(dep, reference.NewMulticlusterOwnerReference(m, m.GroupVersionKind(), req.Context))

	return dep
}

func isInObject(child *appsv1.Deployment, parent string) bool {
	refKind_str := child.ObjectMeta.Annotations["multicluster.admiralty.io/controller-reference"]
	refKind_map := make(map[string]interface{})
	err := json.Unmarshal([]byte(refKind_str), &refKind_map)
	if err != nil {
		panic(err)
	}
	if refKind_map["kind"] == parent {
		return true
	}
	return false
}

func (cm *ClusterManager) DeleteDeployments(nsn types.NamespacedName) error {
	dep := &appsv1.Deployment{}

	for _, cluster := range cm.Cluster_list.Items {
		cluster_client := cm.Cluster_clients[cluster.Name]
		err := cluster_client.Get(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")

		if err != nil && errors.IsNotFound(err) {
			fmt.Println("Not Found")
			continue
		}

		if !isInObject(dep, "OpenMCPDeployment") {
			continue
		}

		err = cluster_client.Delete(context.Background(), dep, nsn.Namespace, nsn.Name+"-deploy")
		if err != nil {
			return err
		}
	}
	return nil
}

func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		fmt.Println("Error retrieving list of federated clusters: %+v", err)
	}
	if len(clusterList.Items) == 0 {
		fmt.Println("No federated clusters found")
	}
	return clusterList
}

func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, client genericclient.Client, fedNamespace string) map[string]*rest.Config {
	clusterConfigs := make(map[string]*rest.Config)
	for _, cluster := range clusterList.Items {
		config, _ := util.BuildClusterConfig(&cluster, client, fedNamespace)
		clusterConfigs[cluster.Name] = config
	}
	return clusterConfigs
}
func KubeFedClusterClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {

	cluster_clients := make(map[string]genericclient.Client)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		cluster_config := cluster_configs[clusterName]
		cluster_client := genericclient.NewForConfigOrDie(cluster_config)
		cluster_clients[clusterName] = cluster_client
	}
	return cluster_clients
}

func NewClusterManager() *ClusterManager {
	fed_namespace := "kube-federation-system"
	host_config, _ := rest.InClusterConfig()
	host_client := genericclient.NewForConfigOrDie(host_config)
	cluster_list := ListKubeFedClusters(host_client, fed_namespace)
	cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)
	cluster_clients := KubeFedClusterClients(cluster_list, cluster_configs)

	cm := &ClusterManager{
		Fed_namespace:   fed_namespace,
		Host_config:     host_config,
		Host_client:     host_client,
		Cluster_list:    cluster_list,
		Cluster_configs: cluster_configs,
		Cluster_clients: cluster_clients,
	}
	return cm
}

func (cm *ClusterManager) Scheduling(replicas int32, request_cpu float64, request_mem float64, policy []ketiv1alpha1.PolicyData) (map[string]int32, int64) {
	rand.Seed(time.Now().UTC().UnixNano())
	klog.Infof("[SCHED] : Start Scheduling... ")

	var failSelected int64 = 0
	clusterResource = make(map[string]resource_Info)
	nodeLevelResource = make(map[string]node_Info)

	err := cm.Get_data()
	if err != nil {
		klog.Infof("[ERROR] %v", err)
	}
	replicas_cluster := map[string]int32{}

	for {
		klog.Infof("#####################Start Schduling algorithm######################")
		//clusterList에 존재하는 cluster들만 스케줄링 알고리즘 적용하도록 필터링 실행
		klog.Infof("[SUJUNE] First Phase, Filtering!!")
		//Resource의 값이 갱신되지 않기 때문에, Filtering의 정확도가 떨어짐
		//node단위 정보는 직접 연산하여 리소스 값이 감소됨
		//Issue에 존재하는 내용이므로 확인할 것
		filterClusterList := NewFiltering(request_cpu, request_mem, posClusterList, nodeLevelResource)
		klog.Infof("[TEST] filterClusterList = %s", filterClusterList)
		klog.Infof("[TEST] posClusterList = %s", posClusterList)

		klog.Infof("[SCHED] Second Phase, Scoring!!")

		scoreInfo := []score_Info{}
		resScore := getScoreResource(request_cpu, request_mem, filterClusterList)
		//resScore를 score list에 저장
		scoreInfo = append(scoreInfo, resScore) // score list
		klog.Infof("[SCHED] ResScore = %s", resScore)

		//mostResource 알고리즘
		mostScore := getMostScore(request_cpu, request_mem, filterClusterList)
		//mostScore를 score list에 저장
		scoreInfo = append(scoreInfo, mostScore)
		klog.Infof("[SCHED] MostScore = %f", mostScore)

		//leastResoruce 알고리즘
		leastScore := getLeastScore(request_cpu, request_mem, filterClusterList)
		//leastScore를 score list에 저장
		scoreInfo = append(scoreInfo, leastScore)
		klog.Infof("[SUJUNE] LeastScore = %f", leastScore)
		klog.Infof("#######################################################")
		klog.Infof("[TEST] ScoreSUM = %f", scoreInfo)

		//policy가 존재하는 경우, Calculator 컨트롤러로부터 데이터 가져오기
		//Network만 존재함, 향후 추가할 필요가 있다면 확장할 것
		for _, policy_data := range policy {
			klog.Infof("[Policy_CHECK] policy_data = %s", policy_data.Rule)
			klog.Infof("[Policy_CHECK] policy_data's type = %T", policy_data.Rule)
			if policy_data.Rule == "network" {
				//Calculator 호출
				//우선순위 후보군 가져오기
				policyScore := Clients(policy_data.Rule)
				//policyScore를 score list에 저장
				scoreInfo = append(scoreInfo, policyScore)
				klog.Infof("[SUJUNE] policyScore = %s", policyScore)
				//Policy 점수와 기본 스케줄러의 점수 ResScore를 합산
			}
		}
		//우선순위화 함수 실행
		//data_Info = map[float64]string 타입으로, 내림차순을 위해 key를 float64으로 선언
		//data_Rank = []float64 타입으로, 내림차순으로 정령됨 index = 0이 가장 큰 값
		//data_Rank값을 data_Info의 key값으로 입력하여 클러스터의 이름을 가져온다.
		//scoreInfo에 저장된 score 값들을 합쳐서 계산
		data_Info, data_Rank := getPriority(scoreInfo)
		klog.Infof("[SUJUNE] data_Info = %s", data_Info)
		klog.Infof("[SUJUNE] data_Rank = %s", data_Rank)

		//배포할 클러스터 이름을 반환함
		//exist 값이 1인 경우, 배포할 클러스터가 선택되는 경우를 뜻하며,
		//exist 값이 0인 경우, 배포할 클러스터가 없는 경우를 뜻함
		resultCluster, exist := selectedCluster(request_cpu, request_mem, data_Info, data_Rank, replicas)
		if exist == 1 {
			replicas_cluster[resultCluster] += 1
			replicas--
			klog.Infof("[Check] replicas_cluster = %s", replicas_cluster)
		}
		// exit = 0 일 경우에 대해서
		// 더이상 배포할 곳이 없는 경우, 무한 루프에 빠질 가능성 있음
		// exit = 0 인 경우에 대해서, failSelected 값 증가하고 반환
		if exist == 0 {
			failSelected++
			replicas--
		}

		if replicas == 0 {
			break
		}
	}

	return replicas_cluster, failSelected

}

//input
// spec_replicas -> 요구하는 값, status_replicas -> 현재 값, status_cluster_replicas_map-> 각 클러스터에 들어가는 값
func (cm *ClusterManager) ReScheduling(spec_replicas int32, status_replicas int32, status_cluster_replicas_map map[string]int32) map[string]int32 {
	rand.Seed(time.Now().UTC().UnixNano())

	result_cluster_replicas_map := make(map[string]int32)
	for k, v := range status_cluster_replicas_map {
		result_cluster_replicas_map[k] = v
	}

	action := "dec"
	replica_rate := spec_replicas - status_replicas
	if replica_rate > 0 {
		action = "inc"
	}

	remain_replica := replica_rate

	for remain_replica != 0 {
		cluster_len := len(result_cluster_replicas_map)
		selected_cluster_target_index := rand.Intn(int(cluster_len))

		target_key := keyOf(result_cluster_replicas_map, selected_cluster_target_index)
		if action == "inc" {
			result_cluster_replicas_map[target_key] += 1
			remain_replica -= 1
		} else {
			if result_cluster_replicas_map[target_key] >= 1 {
				result_cluster_replicas_map[target_key] -= 1
				remain_replica += 1
			}
		}
	}
	keys := make([]string, 0)
	for k, _ := range result_cluster_replicas_map {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("ReScheduling Result: ")
	for _, k := range keys {
		v := result_cluster_replicas_map[k]
		prev_v := status_cluster_replicas_map[k]
		fmt.Println("  ", k, ": ", prev_v, " -> ", v)
	}

	return result_cluster_replicas_map

}
func keyOf(my_map map[string]int32, target_index int) string {
	index := 0
	for k, _ := range my_map {
		if index == target_index {
			return k
		}
		index += 1
	}
	return ""

}

//Pod를 배포할 Clutser를 선택하는 함수
//Input : Request Resource(cpu, mem), 우선순위화 후보군(data_Info, data_Rank), Replicas 수
//Output : 배포할 클러스터의 이름, 배포 가능 여부 값(배포할 클러스터가 존재 : 1, 존재하지 않음 : 0)
func selectedCluster(request_cpu float64, request_mem float64, data_Info map[float64]string, data_Rank []float64, replicas int32) (string, int32) {
	klog.Infof("[SUJUNE] : Start selected_cluster...")

	klog.Infof("[SUJUNE] Final data_Info = %s ", data_Info)
	klog.Infof("[SUJUNE] Final data_Rank = %s ", data_Rank)

	for _, cluster_score := range data_Rank {
		//data_Info[cluster_score] -> cluster 이름으로 나옴
		res_info, _ := clusterResource[data_Info[cluster_score]]
		if res_info.cpu_idle-request_cpu >= 0 && res_info.mem_idle-request_mem >= 0 {
			//최소 요구사항을 만족할 경우,
			//자원 최신화 작업 실행
			res_info.cpu_idle -= request_cpu
			res_info.mem_idle -= request_mem
			clusterResource[data_Info[cluster_score]] = res_info

			//TODO
			//Node에 자원이 없는 경우, 알려주는 매커니즘은 없음...
			node_Resource, _ := nodeLevelResource[data_Info[cluster_score]]
			for key1, value1 := range node_Resource.resource {
				klog.Infof("[Check] Check node's resource key1 = %s", key1)
				if value1.cpu_idle-request_cpu >= 0 && value1.mem_idle-request_mem >= 0 {
					value1.cpu_idle -= request_cpu
					value1.mem_idle -= request_mem
					node_Resource.resource[key1] = value1
					nodeLevelResource[data_Info[cluster_score]] = node_Resource
					break
				}
			}

			//배치할 cluster 이름 반환, 배치 완료 ( 1 : completed, 0 : Not completed)
			res_infos, _ := clusterResource[data_Info[cluster_score]]
			klog.Infof("[Selected]-----------------------------------------------------")
			klog.Infof(fmt.Sprintf("[SUJUNE] cluster = %s\n", data_Info[cluster_score]))
			klog.Infof(fmt.Sprintf("[SUJUNE] cpu_idle = %f\n", res_infos.cpu_idle))
			klog.Infof(fmt.Sprintf("[SUJUNE] total_cpu = %f\n", res_infos.cpu_total))
			klog.Infof(fmt.Sprintf("[SUJUNE] mem_idle = %f\n", res_infos.mem_idle))
			klog.Infof(fmt.Sprintf("[SUJUNE] total_mem = %f\n", res_infos.mem_total))
			return data_Info[cluster_score], 1
		}
	}

	return "No", 0
}

//DB로부터 데이터 가져오기
//Input, Output 없음
func (cm *ClusterManager) Get_data() error {

	clusterResource = make(map[string]resource_Info)
	posClusterList = make([]string, 0)

	cluster_clientset = make(map[string]*kubernetes.Clientset)

	//cluster list 생성
	for _, cluster := range cm.Cluster_list.Items {
		posClusterList = append(posClusterList, cluster.Name)
	}

	for _, cluster := range cm.Cluster_list.Items {
		config, _ := util.BuildClusterConfig(&cluster, cm.Host_client, cm.Fed_namespace)
		cluster_clientset[cluster.Name], _ = kubernetes.NewForConfig(config)

		initAllocatedResource(cluster.Name)
		initNodes(cluster.Name)
	}
	return nil

}

/*
func getScoreResource() map[string]float64 {
        var ResourceScore map[string]float64 = make(map[string]float64)
        for _, c_name := range clusterList  {
                temp := clusterResource[c_name]
                //남아 있는 자원이, 많을수록 높은 점수 부여
                cpu_score := temp.cpu_idle / temp.cpu_total
                mem_score := temp.mem_idle / temp.mem_total
                resScore := (0.5) * cpu_score +  (0.5) * mem_score

                ResourceScore[c_name] = resScore * 10
        }
        return ResourceScore
}
*/

//CPU,Memory 리소스 점수 계산 알고리즘
//DRF 방식을 이용하여, 노드에 남아 있는 자원의 비율과 가장 유사한 노드에게 높은 점수를 부여
//유사한 노드를 계산하기 위해 CosSimility 계산식 사용
//Input : Request Resource(cpu, mem), Join된 클러스터의 리스트
//Output : map[string]float64 타입으로, 클러스터 이름을 key값으로, 클러스터의 점수를 저장
func getScoreResource(request_cpu float64, request_mem float64, c_List []string) map[string]float64 {
	var ResourceScore map[string]float64 = make(map[string]float64)
	//clusterList에만 있는 Cluster의 데이터만 가져오고, 저장함
	for _, c_name := range c_List {
		request_data := []float64{}
		resource_data := []float64{}

		//요청 받은 Pod의 request(cpu, mem) 데이터 수집
		request_data = append(request_data, request_cpu)
		request_data = append(request_data, request_mem)
		temp := clusterResource[c_name]

		//현재 Cluster에서 가용 자원 데이터 수집
		resource_data = append(resource_data, temp.cpu_idle)
		resource_data = append(resource_data, temp.mem_idle)

		//CosSimility 함수 호출
		Cos_result := cosSimility(request_data, resource_data)
		ResourceScore[c_name] = Cos_result * 10
	}
	return ResourceScore
}

//CosSimilty를 계산하는 함수
//Input = []float64 타입으로, a -> Request Resource(cpu, mem순으로 저장)
//                            b -> Idle  Resource(cpu, mem순으로 저장)
//Output = float64 자료형으로 CosSimilty 계산 결과를 반환
func cosSimility(a []float64, b []float64) float64 {
	count := 0
	length_a := len(a)
	length_b := len(b)
	if length_a > length_b {
		count = length_a
	} else {
		count = length_b
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= length_a {
			s2 += math.Pow(b[k], 2)
			continue
		}
		if k >= length_b {
			s1 += math.Pow(a[k], 2)
			continue
		}
		sumA += a[k] * b[k]
		s1 += math.Pow(a[k], 2)
		s2 += math.Pow(b[k], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0
	}
	return sumA / (math.Sqrt(s1) * math.Sqrt(s2))
}

//map 자료형을 변환해주는 함수
// map[string]float64 -> map[float64]string
// key값을 float64으로 변경하여 내림차순으로 map의 내용을 정렬하기 위함
func mapStringToFloat(maps map[string]float64) map[float64]string {
	temp := make(map[float64]string)

	for keys, value := range maps {
		temp[value] = keys
	}
	return temp
}

//우선순위화 함수
//map[float64]string 자료형의 key 값을 기반으로 내림차순으로 []float64 자료형을 생성
//Input : map[string]float64으로 Score map
//Output : map[float64]string과 []float64 자료형을 반환하며, []float64 자료형 데이터가 내림차순 정렬된 배열
func getPriority(scoreList []score_Info) (map[float64]string, []float64) {
	final_score := make(map[string]float64)
	for _, score_data := range scoreList {
		for c_name, score := range score_data {
			klog.Infof("[Priority] c_name = %s, score = %s", c_name, score)
			final_score[c_name] += score
		}
	}

	for key1, value1 := range final_score {
		klog.Infof("[SUM] key1 = %s", key1)
		klog.Infof("[SUM] value1 = %f", value1)
		final_score[key1] = value1 / float64(len(scoreList))
	}
	klog.Infof("[SUM] final_score = %f and lens = %d", final_score, len(scoreList))
	priority_list := mapStringToFloat(final_score)
	//내림차순 함수 호출
	rank_list := getRank_priority(priority_list)

	//cluster 및 점수를 갖는 map 변수와 내림차순으로 정의된 배열 변수를 반환
	return priority_list, rank_list
}

//내림차순으로 정렬된 배열을 생성하는 함수
//Input : map[float64]string으로 Score 점수를 갖는 map 자료형
//Output : []float64 자료형으로 내림차순 정렬된 배열을 반환
func getRank_priority(scores map[float64]string) []float64 {
	keys := make([]float64, 0, len(scores))
	for i := range scores {
		keys = append(keys, i)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))
	// keys는 Score 점수 이며, 이를 내림 차순으로 keys 배열을 생성해서 반환
	return keys
}

//클러스터 단위의 할당된 Pod의 자원 정보를 수집
//Input : 클러스터 이름
//Output : 없음
func initAllocatedResource(c_name string) {
	klog.Infof("[SUJUNE] Start KETI_initAllocatedResource...")
	pods := getRunningPods(c_name)
	allocatedResource = make(map[string]Resource)
	for _, pod := range pods {
		nodeName := pod.NodeName
		var res Resource
		res, ok := allocatedResource[nodeName]
		if ok {
			res.MilliCpu += float64(pod.RequestMilliCpu)
			res.Memory += float64(pod.RequestMemory)
			// 총 사용량 계산

			allocatedResource[nodeName] = res
		} else {
			res.MilliCpu = float64(pod.RequestMilliCpu)
			res.Memory = float64(pod.RequestMemory)
			// 총 사용량 계산
			allocatedResource[nodeName] = res
		}
	}
	klog.Infof("[SUJUNE] AllocatedResource initialization is completed.")
}

//노드에 실행중인 Pod의 정보를 수집
//Pod의 자원 크기 등을 수집하여 저장
//Input : 클러스터 이름
func getRunningPods(c_name string) []Pod {
	klog.Infof("[SUJUNE] Start KETI_getRunningPods...")
	pods, err := cluster_clientset[c_name].CoreV1().Pods("").List(metav1.ListOptions{})

	if err != nil {
		klog.Infof("[SUJUNE] %v ", err)
	}

	runningPods := make([]Pod, 0)
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodRunning {
			var requestsMilliCpu, requestsMemory int64
			for _, ctn := range pod.Spec.Containers {
				requestsMilliCpu += ctn.Resources.Requests.Cpu().MilliValue()
				//requestsMilliCpu += ctn.Resources.Limits.Cpu().MilliValue()
				requestsMemory += ctn.Resources.Requests.Memory().Value() / 1024 / 1024
				//requestsMemory += ctn.Resources.Limits.Memory().Value() / 1024 / 1024
			}
			newPod := Pod{
				Name:            pod.Name,
				Uid:             pod.Namespace,
				NodeName:        pod.Spec.NodeName,
				RequestMilliCpu: requestsMilliCpu,
				RequestMemory:   requestsMemory,
			}
			runningPods = append(runningPods, newPod)
		}
	}
	// runningPods는 실행중인 Pod의 정보들을 담아서 저장한 곳
	return runningPods

}

//Node 단위, 클러스터 단위의 정보를 수집하는 함수
//OpenMCP는 클러스터 단위로 자원 정보를 수집하도록 설계되어 있음
//스케줄링의 필터링에서 사용될 Node 단위의 요구자원 확인을 위해 Node 단위의 정보를 수집하도록 설계되어 있음
//Input : 클러스터 이름
func initNodes(c_name string) {
	var cpu_sum, mem_sum, total_csum, total_msum float64
	cpu_sum = 0
	mem_sum = 0
	total_csum = 0
	total_msum = 0
	klog.Infof("[SUJUNE] Start initNodes...")
	availableNodes = make([]Node, 0)
	nodes, err := cluster_clientset[c_name].CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		klog.Infof("[ERROR] %v", err)
	}
	for _, node := range nodes.Items {
		//              KETIfile.WriteString(fmt.Sprintf("node Info = %s\n", node.Status))
		Ip_info := node.Status.Addresses[0]
		Host_info := node.Status.Addresses[1]
		klog.Infof("[KETI] node address = %s", Ip_info.Address)
		klog.Infof("[KETI] host info = %s", Host_info)
		newNode_cpu_mem := Node{
			Name: node.Name,
			Resource: Resource{
				MilliCpu: float64(node.Status.Capacity.Cpu().MilliValue()),
				Memory:   float64(node.Status.Capacity.Memory().Value() / 1024 / 1024),
			},
		}
		//Node 들의 전체 자원 크기 수집
		availableNodes = append(availableNodes, newNode_cpu_mem)
	}
	var Node_Info node_Info
	Node_Info.resource = make(map[string]node_Resource)
	for Node_name, Node_res := range allocatedResource {
		if !strings.Contains(Node_name, "master") {
			// Node의 전체 자원 크기 가져오기
			cpu_total, mem_total := getTotalResource(Node_name)
			// 데이터 여부 체크
			cpu_sum += cpu_total - (Node_res.MilliCpu / 1000)
			mem_sum += mem_total - (Node_res.Memory / 1024)
			total_csum += cpu_total
			total_msum += mem_total

			//NodeLevel data 저장
			//Golang의 고질적인 문제점
			//멀티 Struct 자료형의 경우, 하위 Struct를 직접 정의하고 데이터 선언해야함
			var Node_Resource node_Resource
			Node_Resource.cpu_idle = cpu_total - (Node_res.MilliCpu / 1000)
			Node_Resource.mem_idle = mem_total - (Node_res.Memory / 1024)

			Node_Info.name = Node_name
			Node_Info.resource[Node_name] = Node_Resource
		}
	}
	nodeLevelResource[c_name] = Node_Info
	//klog.Infof("#######################################################")
	//klog.Infof(fmt.Sprintf("[TEST] nodeLevelResource = %s \n", nodeLevelResource))
	//klog.Infof("#######################################################")
	var temp resource_Info
	temp.cpu_idle = cpu_sum
	temp.cpu_total = total_csum
	temp.mem_idle = mem_sum
	temp.mem_total = total_msum
	clusterResource[c_name] = temp
	klog.Infof("#######################################################")
	klog.Infof(fmt.Sprintf("[SUJUNE] %s's cluster\n", c_name))
	klog.Infof(fmt.Sprintf("[SUJUNE] cpu_idle = %f\n", temp.cpu_idle))
	klog.Infof(fmt.Sprintf("[SUJUNE] total_cpu = %f\n", temp.cpu_total))
	klog.Infof(fmt.Sprintf("[SUJUNE] mem_idle = %f\n", temp.mem_idle))
	klog.Infof(fmt.Sprintf("[SUJUNE] total_mem = %f\n", temp.mem_total))

	fmt.Println("AvailableNodes initialization is completed..")
}

//Node의 전체 자원 크기를 반환하는 함수
//Input : Node 이름
//Output : Total CPU와 Memory 크기(Core, GiB 단위)
func getTotalResource(N_name string) (float64, float64) {
	for _, node := range availableNodes {
		if node.Name == N_name {
			return node.Resource.MilliCpu / 1000, node.Resource.Memory / 1024
		}
	}
	return 0, 0
}

//CPU Core가 m(milli)단위로 입력될 경우, Core단위로 변환해주는 함수
//1000m -> 1 core, 500m -> 0.5 core
//Input : CPU의 자원량(Milli 단위)
//Ouput : CPU의 자원량(Core 단위)
func milliTocoreCPU(value string) float64 {
	temp := strings.Split(value, "m")
	milli_value, _ := strconv.Atoi(temp[0])
	core_value := float64(milli_value)
	Request_cpu := core_value / 1000

	return Request_cpu

}
