package openmcpscheduler

import (
	"time"
	"strings"
	"k8s.io/klog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
	ketiframework "openmcp/openmcp/openmcp-scheduler/pkg/framework/v1alpha1"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/openmcp-scheduler/pkg/protobuf"
)

type OpenMCPScheduler struct {
	ClusterConfigs 	map[string]*rest.Config
	ClusterClients	map[string]*kubernetes.Clientset
	ClusterInfos	map[string]*ketiresource.Cluster
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
				klog.V(0).Info("cannot use resource : ", rName.String())
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

func Scheduling (cm *clusterManager.ClusterManager, dep *ketiv1alpha1.OpenMCPDeployment, grpcServer protobuf.RequestAnalysisClient) map[string]int32 {
	depReplicas := dep.Spec.Replicas

	// Return scheduling result (ex. cluster1:2, cluster2:1)
	totalSchedulingResult := map[string]int32{}

	var sched *OpenMCPScheduler
	sched = &OpenMCPScheduler {
		ClusterConfigs:		cm.Cluster_configs,
		ClusterClients:		make(map[string]*kubernetes.Clientset),
		ClusterInfos:		make(map[string]*ketiresource.Cluster),
		Framework:			ketiframework.NewFramework(grpcServer),
	}

	// Get Data from Node&Pod Spec
	sched.SetupResources()

	// Make resource to schedule pod into cluster
	// startTime := time.Now()
	newPod := newPodFromOpenMCPDeployment(dep)
	// elapsedTime := time.Since(startTime)
	// klog.Infof("*********** [TIME] %s ***********", elapsedTime)

	klog.Infof("***** [Start] Scheduling *****")
	startTime := time.Now()
	// Scheduling one pod
	for i := int32(0); i < depReplicas; i++ {
		startTime2 := time.Now()
		// klog.Info("*********** [START SCHEDULING POD]  ***********")
		schedulingResult := sched.ScheduleOne(newPod)
		// klog.Infof("*********** [END SCHEDULING POD] ***********")

		_, exists := totalSchedulingResult[schedulingResult]
		if !exists {
			totalSchedulingResult[schedulingResult] = 1
		} else{
			totalSchedulingResult[schedulingResult] += 1
		}

		// sched.UpdateResources(newPod, schedulingResult)
		elapsedTime2 := time.Since(startTime2)
		klog.Infof("=> %d. Filtering & Scoring Time [%v] ", i, elapsedTime2)
	}

	elapsedTime := time.Since(startTime)
	klog.Infof("=> Scheduling Result : ", totalSchedulingResult)
	klog.Infof("=> Scheduling Time [%v]", elapsedTime)
	klog.Infof("***** [End] Scheduling *****")
	
	return totalSchedulingResult
}

func (sched *OpenMCPScheduler) UpdateResources (newPod *ketiresource.Pod, schedulingResult string) {
	// startTime := time.Now()
	startTime := time.Now()

	var maxScoreNode *ketiresource.NodeInfo
	maxScore := int64(0)

	for _, node := range sched.ClusterInfos[schedulingResult].Nodes {
		if maxScore < node.NodeScore {
			maxScoreNode = node
			maxScore = node.NodeScore
		}
	}

	maxScoreNode.RequestedResource = ketiresource.AddResources(maxScoreNode.RequestedResource, newPod.RequestedResource)

	elapsedTime := time.Since(startTime)
	klog.Infof("=> Total updateResource Time : %v", elapsedTime)

	// elapsedTime := time.Since(startTime)
	// klog.Infof("=> Update Resource time [%v]", elapsedTime)
}

func (sched *OpenMCPScheduler) ScheduleOne (newPod *ketiresource.Pod) string {
	// startTime := time.Now()
	// klog.Infof("*********** [START FILTERING] ***********")
	filterdResult := sched.Framework.RunFilterPluginsOnClusters(newPod, sched.ClusterInfos)
	// elapsedTime := time.Since(startTime)
	// klog.Infof("*********** [END FILTERING] %v ***********", elapsedTime)


	filteredCluster := make(map[string]*ketiresource.Cluster)
	for clusterName, isfiltered := range filterdResult {
		if isfiltered {
			filteredCluster[clusterName] = sched.ClusterInfos[clusterName]
		}
	}

	// startTime = time.Now()
	// klog.Infof("*********** [START SCORING] ***********")
	scoreResult := sched.Framework.RunScorePluginsOnClusters(newPod, filteredCluster)
	// elapsedTime = time.Since(startTime)
	// klog.Infof("*********** [END SCORING] %v ***********", elapsedTime)

	selectedCluster := selectCluster(scoreResult)

	return selectedCluster
}

func selectCluster (scoreResult ketiframework.OpenmcpPluginToClusterScores) string{
	var selectedCluster string
	var maxScore int64

	startTime := time.Now()

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

	elapsedTime := time.Since(startTime)
	klog.Infof("=> Total selectCluster Time : %v", elapsedTime)

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
				NodeScore:				0,
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

	// klog.Infof("[CHECK INFORMATIONS] %v", sched.ClusterInfos)

	return nil
}

func openmcpContainersToContainers(containers []ketiv1alpha1.OpenMCPContainer) []corev1.Container {
	var newContainers []corev1.Container

	for _, container := range containers {
		newContainer := corev1.Container{
			Name:       container.Name,
			Image:      container.Image,
			Command:    container.Command,
			Args:       container.Args,
			WorkingDir: container.WorkingDir,
			Ports:      container.Ports,
			EnvFrom:    container.EnvFrom,
			Env:        container.Env,
			Resources: corev1.ResourceRequirements{
				Limits:   container.Resources.Limits,
				Requests: container.Resources.Requests,
			},
			VolumeMounts:             container.VolumeMounts,
			VolumeDevices:            container.VolumeDevices,
			LivenessProbe:            container.LivenessProbe,
			ReadinessProbe:           container.ReadinessProbe,
			Lifecycle:                container.Lifecycle,
			TerminationMessagePath:   container.TerminationMessagePath,
			TerminationMessagePolicy: container.TerminationMessagePolicy,
			ImagePullPolicy:          container.ImagePullPolicy,
			SecurityContext:          container.SecurityContext,
			Stdin:                    container.Stdin,
			StdinOnce:                container.StdinOnce,
			TTY:                      container.TTY,
		}
		newContainers = append(newContainers, newContainer)
	}

	return newContainers
}


func openmcpPodSpecToPodSpec(spec ketiv1alpha1.OpenMCPPodSpec) corev1.PodSpec {
	return corev1.PodSpec{
		Volumes:                       spec.Volumes,
		InitContainers:                openmcpContainersToContainers(spec.InitContainers),
		Containers:                    openmcpContainersToContainers(spec.Containers),
		RestartPolicy:                 spec.RestartPolicy,
		TerminationGracePeriodSeconds: spec.TerminationGracePeriodSeconds,
		ActiveDeadlineSeconds:         spec.ActiveDeadlineSeconds,
		DNSPolicy:                     spec.DNSPolicy,
		NodeSelector:                  spec.NodeSelector,
		ServiceAccountName:            spec.ServiceAccountName,
		DeprecatedServiceAccount:      spec.DeprecatedServiceAccount,
		AutomountServiceAccountToken:  spec.AutomountServiceAccountToken,
		NodeName:                      spec.NodeName,
		HostNetwork:                   spec.HostNetwork,
		HostPID:                       spec.HostPID,
		HostIPC:                       spec.HostIPC,
		ShareProcessNamespace:         spec.ShareProcessNamespace,
		SecurityContext:               spec.SecurityContext,
		ImagePullSecrets:              spec.ImagePullSecrets,
		Hostname:                      spec.Hostname,
		Subdomain:                     spec.Subdomain,
		Affinity:                      spec.Affinity,
		SchedulerName:                 spec.SchedulerName,
		Tolerations:                   spec.Tolerations,
		HostAliases:                   spec.HostAliases,
		PriorityClassName:             spec.PriorityClassName,
		Priority:                      spec.Priority,
		DNSConfig:                     spec.DNSConfig,
		ReadinessGates:                spec.ReadinessGates,
		RuntimeClassName:              spec.RuntimeClassName,
		EnableServiceLinks:            spec.EnableServiceLinks,
	}
}

func openmcpPodTemplateSpecToPodTemplateSpec(template ketiv1alpha1.OpenMCPPodTemplateSpec) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: template.ObjectMeta,
		Spec:       openmcpPodSpecToPodSpec(template.Spec),
	}
}

func openmcpDeploymentTemplateSpecToDeploymentSpec(templateSpec ketiv1alpha1.OpenMCPDeploymentTemplateSpec) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas:                templateSpec.Replicas,
		Selector:                templateSpec.Selector,
		Template:                openmcpPodTemplateSpecToPodTemplateSpec(templateSpec.Template),
		Strategy:                templateSpec.Strategy,
		MinReadySeconds:         templateSpec.MinReadySeconds,
		RevisionHistoryLimit:    templateSpec.RevisionHistoryLimit,
		Paused:                  templateSpec.Paused,
		ProgressDeadlineSeconds: templateSpec.ProgressDeadlineSeconds,
	}
}
