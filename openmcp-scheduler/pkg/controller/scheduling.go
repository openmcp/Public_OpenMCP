package openmcpscheduler

import (
	"time"
	"strings"
	"k8s.io/klog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ketiv1alpha1 "openmcp/openmcp/openmcp-scheduler/pkg/apis/keti/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/controller/resourceinfo"
	ketiframework "openmcp/openmcp/openmcp-scheduler/pkg/controller/framework/v1alpha1"
)

type OpenMCPScheduler struct {
	ClusterConfigs 	map[string]*rest.Config
	ClusterClients	map[string]*kubernetes.Clientset
	ClusterInfos	map[string]*ketiresource.Cluster
	// Framework runs scheduler plugins (Filtering & Scoring)
	Framework		ketiframework.OpenmcpFramework
}

// Returns ketiresource.Resource if specified
func newPodFromOpenMCPDeployment(dep *ketiv1alpha1.OpenMCPDeployment) *ketiresource.Pod {
	res := ketiresource.NewResource()
	needRes := make(map[string]bool)
	affinities := make(map[string][]string)

	for _, container := range dep.Spec.Template.Spec.Template.Spec.Containers {
		for rName, rQuant := range container.Resources.Requests {
			switch rName {
			case corev1.ResourceCPU:
				res.MilliCPU = rQuant.MilliValue()		
			case corev1.ResourceMemory:
				res.Memory = rQuant.Value()
			case corev1.ResourceEphemeralStorage:
				res.EphemeralStorage = rQuant.Value() 
			default:
				klog.Infof("cannot use resource : %s", rName.String())
			}
		}

		for rName, rNeeded := range container.Resources.Needs {
			val, exist := needRes[rName.String()]

			if !exist {
				needRes[rName.String()] = rNeeded
			} else {
				if val == false && rNeeded == true {
					needRes[rName.String()] = true
				}
			}
		}

		for key, values := range dep.Spec.Affinity {
			for _, value :=  range values {
				affinities[key] = append(affinities[key], value)
			}
		}
	}

	return &ketiresource.Pod {
		Pod:				&corev1.Pod{
			Spec:		openmcpPodSpecToPodSpec(dep.Spec.Template.Spec.Template.Spec),
		},
		RequestedResource: 	res,
		Affinity: 			affinities,
		IsNeedResourceMap:	needRes,
	}
}

func(cm *ClusterManager) Scheduling (dep *ketiv1alpha1.OpenMCPDeployment) map[string]int32 {
	klog.Infof("*********** Scheduling ***********")
	depReplicas := dep.Spec.Replicas

	// Return scheduling result (ex. cluster1:2, cluster2:1)
	totalSchedulingResult := map[string]int32{}

	var sched *OpenMCPScheduler
	sched = &OpenMCPScheduler {
		ClusterConfigs:		cm.Cluster_configs,
		ClusterClients:		make(map[string]*kubernetes.Clientset),
		ClusterInfos:		make(map[string]*ketiresource.Cluster),
		Framework:			ketiframework.NewFramework(),
	}

	// Get Data from Node&Pod Spec
	sched.SetupResources()

	// Make resource to schedule pod into cluster
	startTime := time.Now()
	newPod := newPodFromOpenMCPDeployment(dep)
	elapsedTime := time.Since(startTime)
	klog.Infof("*********** [TIME] %s ***********", elapsedTime)

	// Scheduling one pod
	for i := int32(0); i < depReplicas; i++ {
		klog.Info("*********** [START SCHEDULING POD]  ***********")
		schedulingResult := sched.ScheduleOne(newPod)
		klog.Infof("*********** [END SCHEDULING POD] %v ***********", schedulingResult)

		_, exists := totalSchedulingResult[schedulingResult]
		if !exists {
			totalSchedulingResult[schedulingResult] = 1
		} else{
			totalSchedulingResult[schedulingResult] += 1
		}
	}
	
	return totalSchedulingResult
}

func (sched *OpenMCPScheduler) ScheduleOne (newPod *ketiresource.Pod) string {
	startTime := time.Now()
	klog.Infof("*********** [START FILTERING] ***********")
	filterdResult := sched.Framework.RunFilterPluginsOnClusters(newPod, sched.ClusterInfos)
	elapsedTime := time.Since(startTime)
	klog.Infof("*********** [END FILTERING] %v ***********", filterdResult)
	klog.Infof("*********** [TOTAL FILTERING TIME] %s ***********", elapsedTime)


	filteredCluster := make(map[string]*ketiresource.Cluster)
	for clusterName, isfiltered := range filterdResult {
		if isfiltered {
			filteredCluster[clusterName] = sched.ClusterInfos[clusterName]
		}
	}

	startTime = time.Now()
	klog.Infof("*********** [START SCORING] ***********")
	scoreResult := sched.Framework.RunScorePluginsOnClusters(newPod, filteredCluster)
	elapsedTime = time.Since(startTime)
	klog.Infof("*********** [END SCORING] %v ***********", scoreResult)
	klog.Infof("*********** [TOTAL SCORING TIME] %s ***********", elapsedTime)

	selectedCluster := selectCluster(scoreResult)

	return selectedCluster
}

func selectCluster (scoreResult ketiframework.OpenmcpPluginToClusterScores) string{
	var selectedCluster string
	var maxScore int64

	for clusterName, scoreList := range scoreResult {
		var clusterScore int64
		for _, score := range scoreList {
			clusterScore += score.Score
		}

		if clusterScore > maxScore {
			selectedCluster = clusterName
			maxScore = clusterScore
		}
	}
	return selectedCluster
}

func (sched *OpenMCPScheduler)SetupResources() error {
	// get ClusterClients
	for clusterName, config := range sched.ClusterConfigs {
		sched.ClusterClients[clusterName], _ = kubernetes.NewForConfig(config)
	}

	// Setup Clusters
	for clusterName, _ := range sched.ClusterClients {
		pods, _ := sched.ClusterClients[clusterName].CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{})

		// informations on cluster level
		allPods := make([]*ketiresource.Pod, 0)
		allNodes := make([]*ketiresource.NodeInfo, 0)
		cluster_request := ketiresource.NewResource()
		cluster_allocatable := ketiresource.NewResource()

		// Setup Pods
		for _, pod := range pods.Items {
			// add Stroage
			var pod_requestMilliCpu, pod_requestMemory, pod_requestEStorage int64

			for _, ctn := range pod.Spec.Containers {
				pod_requestMilliCpu += ctn.Resources.Requests.Cpu().MilliValue()
				pod_requestMemory += ctn.Resources.Requests.Memory().Value()
				pod_requestEStorage += ctn.Resources.Requests.StorageEphemeral().Value()
			}

			newPod := &ketiresource.Pod {
				Pod:					&pod,
				PodName:           		pod.Name,
				NodeName:       		pod.Spec.NodeName,
				RequestedResource: 		&ketiresource.Resource {
					MilliCPU:	pod_requestMilliCpu,
					Memory:		pod_requestMemory,
					EphemeralStorage:	pod_requestEStorage,
				},
			}
			allPods = append(allPods, newPod)
		}

		// Setup Nodes
		nodes, _ := sched.ClusterClients[clusterName].CoreV1().Nodes().List(metav1.ListOptions{})
		for _, node := range nodes.Items {
			node_request := ketiresource.NewResource()
			var node_usedPorts []*corev1.ContainerPort

			// Get Pods belonging to this node
			podsInNode := make([]*ketiresource.Pod, 0)
			for _, pod := range allPods {
				if strings.Compare(pod.NodeName, node.Name) == 0 {
					podsInNode = append(podsInNode, pod)
					node_request = ketiresource.AddResources(node_request, pod.RequestedResource)
					
					for j := range pod.Pod.Spec.Containers {
						container := &pod.Pod.Spec.Containers[j]
						for k := range container.Ports {
							node_usedPorts = append(node_usedPorts, &container.Ports[k])
						}
					}
				}
			}

			// Get capacity From node Spec
			node_capacity := &ketiresource.Resource {
				MilliCPU:			node.Status.Capacity.Cpu().MilliValue(),
				Memory:				node.Status.Capacity.Memory().Value(),
				EphemeralStorage:	node.Status.Capacity.StorageEphemeral().Value(),
			}

			// Get allocatable Resource based on capacity and request
			node_allocatable := ketiresource.GetAllocatable(node_capacity, node_request)

			// Get Hardware spec
			node_needRes := make(map[string]bool)
			// for test.. temporary
			node_needRes["gpu"] = true
			node_needRes["hba"] = false

			// Get Affinity
			node_affinity := make(map[string]string)

			for key, value := range node.Labels {
				switch key {
				case "failure-domain.beta.kubernetes.io/region":
					if _, ok := node_affinity["region"]; !ok{
						node_affinity["region"] = value
					}

				case "failure-domain.beta.kubernetes.io/zone":
					if _, ok := node_affinity["zone"]; !ok{
						node_affinity["zone"] = value
					}
				}
			}

			// make new Node 
			newNode := &ketiresource.NodeInfo {
				Node:					&node,
				NodeName:				node.Name,
				Pods:					podsInNode,
				UsedPorts:				node_usedPorts,
				CapacityResource:		node_capacity,
				RequestedResource:		node_request,
				AllocatableResource:	node_allocatable,
				IsNeedResourceMap:		node_needRes,
				Affinity:				node_affinity,
			}
			allNodes = append(allNodes, newNode)
			cluster_request = ketiresource.AddResources(cluster_request, node_request)
			cluster_allocatable = ketiresource.AddResources(cluster_allocatable, node_allocatable)
		}

		// Setup Cluster
		sched.ClusterInfos[clusterName] = &ketiresource.Cluster {
			ClusterName:			clusterName,
			Nodes:					allNodes,
			RequestedResource: 		cluster_request,
			AllocatableResource:	cluster_allocatable,
		}
	}

	klog.Infof("[CHECK INFORMATIONS] %v", sched.ClusterInfos)

	return nil
}