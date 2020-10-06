package openmcpscheduler

import (
	"context"
	"fmt"
	"openmcp/openmcp/omcplog"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	ketiframework "openmcp/openmcp/openmcp-scheduler/pkg/framework/v1alpha1"
	"openmcp/openmcp/openmcp-scheduler/pkg/protobuf"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
	"openmcp/openmcp/util/clusterManager"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type OpenMCPScheduler struct {
	ClusterClients map[string]*kubernetes.Clientset
	ClusterInfos   map[string]*ketiresource.Cluster
	Framework      ketiframework.OpenmcpFramework
	ClusterManager *clusterManager.ClusterManager
	GRPC_Client    protobuf.RequestAnalysisClient
	IsNetwork      bool
	PostPods       []*ketiresource.Pod
}

func NewScheduler(cm *clusterManager.ClusterManager, grpcClient protobuf.RequestAnalysisClient) *OpenMCPScheduler {
	sched := &OpenMCPScheduler{}
	sched.ClusterClients = make(map[string]*kubernetes.Clientset)
	sched.ClusterInfos = make(map[string]*ketiresource.Cluster)
	sched.Framework = ketiframework.NewFramework(grpcClient)
	sched.ClusterManager = cm
	sched.GRPC_Client = grpcClient
	sched.PostPods = make([]*ketiresource.Pod, 0)
	if !sched.IsNetwork {
		go sched.LocalNetworkAnalysis()
		omcplog.V(0).Infof("LocalNetworkAnalysis Start")

	}
	return sched
}

func (sched *OpenMCPScheduler) Scheduling(dep *resourcev1alpha1.OpenMCPDeployment) (map[string]int32, error) {
	startTime := time.Now()
	cm := sched.ClusterManager

	// Get CLusterClients from clusterManager
	sched.ClusterClients = cm.Cluster_kubeClients

	// Return scheduling result (ex. cluster1:2, cluster2:1)
	totalSchedulingResult := map[string]int32{}

	if sched.IsNetwork == false {
		sched.SetupResources()
		omcplog.V(0).Infof("SetupResources")
		sched.IsNetwork = true
	}
	// Get Data from Node&Pod Spec
	//sched.SetupResources()

	depReplicas := dep.Spec.Replicas

	// Make resource to schedule pod into cluster
	newPod := newPodFromOpenMCPDeployment(dep)


	// Scheduling one pod
	for i := int32(0); i < depReplicas; i++ {

		// If there is no proper cluster to deploy Pod,
		// stop scheduling and return scheduling result
		schedulingResult, err := sched.ScheduleOne(newPod, depReplicas)
		if err != nil {
			return totalSchedulingResult, fmt.Errorf("There is no proper cluster to deploy Pod(%d)~Pod(%d)", i, depReplicas)
		}
		if schedulingResult == "" {
			continue
		}

		_, exists := totalSchedulingResult[schedulingResult]
		if !exists {
			totalSchedulingResult[schedulingResult] = 1
		} else {
			totalSchedulingResult[schedulingResult] += 1
		}

		sched.UpdateResources(newPod, schedulingResult)
	}
	sched.Framework.EndPod()
	elapsedTime := time.Since(startTime)
	omcplog.V(0).Infof("    => Scheduling Time [%v]", elapsedTime)


	return totalSchedulingResult, nil
}
func (sched *OpenMCPScheduler) ScheduleOne(newPod *ketiresource.Pod, replicas int32) (string, error) {
	filterdResult := sched.Framework.RunFilterPluginsOnClusters(newPod, sched.ClusterInfos)

	filteredCluster := make(map[string]*ketiresource.Cluster)

	for clusterName, isfiltered := range filterdResult {
		if isfiltered {
			filteredCluster[clusterName] = sched.ClusterInfos[clusterName]
		}
	}

	if len(filteredCluster) == 0 {
		postresult := sched.Framework.RunPostFilterPluginsOnClusters(newPod, sched.ClusterInfos, sched.PostPods)

		if postresult["unscheduable"] {
			return "", fmt.Errorf("There is postpods error")
		}
		if postresult["error"] {
			return "", fmt.Errorf("There is PostFilter Error")
		}
		if postresult["success"] {
			sched.PostPods = append(sched.PostPods, newPod)
			return "", nil
		}
		return "", fmt.Errorf("There is no cluster")
	}

	selectedCluster := sched.Framework.RunScorePluginsOnClusters(newPod, filteredCluster, sched.ClusterInfos, replicas)

	return selectedCluster, nil
}

func (sched *OpenMCPScheduler) LocalNetworkAnalysis() {
	clusters := sched.ClusterInfos
	for {

		for _, cluster := range clusters {
			for _, node := range cluster.Nodes {
				var nodeScore int64
				node_info := &protobuf.NodeInfo{ClusterName: cluster.ClusterName, NodeName: node.NodeName}
				client := sched.GRPC_Client
				result, err := client.SendNetworkAnalysis(context.TODO(), node_info)
				if result.RX == -1 || result.TX == -1 {
					continue
				}
				if err != nil || result == nil {
					omcplog.V(0).Infof("cannot get %v's data from openmcp-analytic-engine", node.NodeName)
					continue
				}
				node.UpdateRX = result.RX
				node.UpdateTX = result.TX
				if node.UpdateRX == 0 && node.UpdateTX == 0 {
					nodeScore = 100
				} else {
					nodeScore = int64((1 / float64(node.UpdateRX+node.UpdateTX)) * float64(100))
				}
				node.NodeScore = nodeScore
			}
		}
		time.Sleep(1 * time.Second)
	}

}

func (sched *OpenMCPScheduler) UpdateResources(newPod *ketiresource.Pod, schedulingResult string) {

	var maxScoreNode *ketiresource.NodeInfo
	maxScore := int64(0)

	for _, node := range sched.ClusterInfos[schedulingResult].Nodes {
		if maxScore < node.NodeScore {
			maxScoreNode = node
			maxScore = node.NodeScore
		}
	}

	maxScoreNode.RequestedResource = ketiresource.AddResources(maxScoreNode.RequestedResource, newPod.RequestedResource)
	maxScoreNode.AllocatableResource = ketiresource.GetAllocatable(maxScoreNode.CapacityResource, maxScoreNode.RequestedResource)
}

// Returns ketiresource.Resource if specified
func newPodFromOpenMCPDeployment(dep *resourcev1alpha1.OpenMCPDeployment) *ketiresource.Pod {
	res := ketiresource.NewResource()
	additionalResource := make([]string, 0)
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
				// Casting from ResourceName to stirng because rName is ResourceName type
				resourceName := fmt.Sprintf("%s", rName)
				additionalResource = append(additionalResource, resourceName)
			}
		}
		for key, values := range dep.Spec.Affinity {
			for _, value := range values {
				affinities[key] = append(affinities[key], value)
			}
		}
	}

	return &ketiresource.Pod{
		Pod: &corev1.Pod{
			Spec: openmcpPodSpecToPodSpec(dep.Spec.Template.Spec.Template.Spec),
		},
		RequestedResource:  res,
		AdditionalResource: additionalResource,
		Affinity:           affinities,
	}
}
func (sched *OpenMCPScheduler) SetupResources() error {
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
			pod_request := &ketiresource.Resource{0, 0, 0}
			pod_additionalResource := make([]string, 0)

			for _, container := range pod.Spec.Containers {
				for rName, rQuant := range container.Resources.Requests {
					switch rName {
					case corev1.ResourceCPU:
						pod_request.MilliCPU = rQuant.MilliValue()
					case corev1.ResourceMemory:
						pod_request.Memory = rQuant.Value()
					case corev1.ResourceEphemeralStorage:
						pod_request.EphemeralStorage = rQuant.Value()
					default:
						// Casting from ResourceName to stirng because rName is ResourceName type
						resourceName := fmt.Sprintf("%s", rName)
						pod_additionalResource = append(pod_additionalResource, resourceName)
					}
				}
			}

			newPod := &ketiresource.Pod{
				Pod:                &pod,
				ClusterName:        clusterName,
				NodeName:           pod.Spec.NodeName,
				PodName:            pod.Name,
				RequestedResource:  pod_request,
				AdditionalResource: pod_additionalResource,
			}
			allPods = append(allPods, newPod)
		}

		// Setup Nodes
		nodes, _ := sched.ClusterClients[clusterName].CoreV1().Nodes().List(metav1.ListOptions{})
		for _, node := range nodes.Items {

			// Get v1.Pod, corev1.ContainerPort and RequestResource
			podsInNode := make([]*ketiresource.Pod, 0)
			node_request := ketiresource.NewResource()

			for _, pod := range allPods {
				if strings.Compare(pod.NodeName, node.Name) == 0 {
					podsInNode = append(podsInNode, pod)
					node_request = ketiresource.AddResources(node_request, pod.RequestedResource)
				}
			}

			// Get capacity, Additional Resource from node Spec
			node_additionalResource := make([]string, 0)
			node_capacity := &ketiresource.Resource{}

			for rName, rQuant := range node.Status.Capacity {
				switch rName {
				case corev1.ResourceCPU:
					node_capacity.MilliCPU = rQuant.MilliValue()
				case corev1.ResourceMemory:
					node_capacity.Memory = rQuant.Value()
				case corev1.ResourceEphemeralStorage:
					node_capacity.EphemeralStorage = rQuant.Value()
				default:
					// Casting from ResourceName to stirng because rName is ResourceName type
					resourceName := fmt.Sprintf("%s", rName)
					node_additionalResource = append(node_additionalResource, resourceName)
				}
			}

			// Get allocatable Resource based on capacity and request
			node_allocatable := ketiresource.GetAllocatable(node_capacity, node_request)

			// Get Affinity
			node_affinity := make(map[string]string)

			for key, value := range node.Labels {
				switch key {
				case "failure-domain.beta.kubernetes.io/region":
					if _, ok := node_affinity["region"]; !ok {
						node_affinity["region"] = value
					}

				case "failure-domain.beta.kubernetes.io/zone":
					if _, ok := node_affinity["zone"]; !ok {
						node_affinity["zone"] = value
					}
				}
			}

			// make new Node
			newNode := &ketiresource.NodeInfo{
				ClusterName: clusterName,
				NodeName:    node.Name,
				Node:        &node,
				Pods:        podsInNode,
				// UsedPorts:				node_usedPorts,
				CapacityResource:    node_capacity,
				RequestedResource:   node_request,
				AllocatableResource: node_allocatable,
				AdditionalResource:  node_additionalResource,
				Affinity:            node_affinity,
				NodeScore:           0,
			}
			allNodes = append(allNodes, newNode)
			cluster_request = ketiresource.AddResources(cluster_request, node_request)
			cluster_allocatable = ketiresource.AddResources(cluster_allocatable, node_allocatable)
		}

		// Setup Cluster
		sched.ClusterInfos[clusterName] = &ketiresource.Cluster{
			ClusterName:         clusterName,
			Nodes:               allNodes,
			RequestedResource:   cluster_request,
			AllocatableResource: cluster_allocatable,
		}
	}

	return nil
}

func openmcpContainersToContainers(containers []resourcev1alpha1.OpenMCPContainer) []corev1.Container {
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

func openmcpPodSpecToPodSpec(spec resourcev1alpha1.OpenMCPPodSpec) corev1.PodSpec {
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

func openmcpPodTemplateSpecToPodTemplateSpec(template resourcev1alpha1.OpenMCPPodTemplateSpec) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: template.ObjectMeta,
		Spec:       openmcpPodSpecToPodSpec(template.Spec),
	}
}

func openmcpDeploymentTemplateSpecToDeploymentSpec(templateSpec resourcev1alpha1.OpenMCPDeploymentTemplateSpec) appsv1.DeploymentSpec {
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
