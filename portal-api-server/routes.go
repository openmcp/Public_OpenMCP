package main

import (
	"net/http"
	"portal-api-server/handler"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"changeekstype",
		"POST",
		"/apis/changeekstype",
		handler.ChangeEKSInstanceType,
	},
	Route{
		"starteksnode",
		"POST",
		"/apis/starteksnode",
		handler.StartEKSNode,
	},
	Route{
		"stopeksnode",
		"POST",
		"/apis/stopeksnode",
		handler.StopEKSNode,
	},
	Route{
		"geteksclusterinfo",
		"POST",
		"/apis/geteksclusterinfo",
		GetEKSClusterInfo,
	},
	Route{
		"deletekvmnode",
		"POST",
		"/apis/deletekvmnode",
		handler.DeleteKVMNode,
	},
	Route{
		"createkvmnode",
		"POST",
		"/apis/createkvmnode",
		handler.CreateKVMNode,
	},
	Route{
		"changekvmnode",
		"POST",
		"/apis/changekvmnode",
		handler.ChangeKVMNode,
	},
	Route{
		"stopkvmnode",
		"POST",
		"/apis/stopkvmnode",
		handler.StopKVMNode,
	},
	Route{
		"startkvmnode",
		"POST",
		"/apis/startkvmnode",
		handler.StartKVMNode,
	},
	Route{
		"getkvmnodes",
		"GET",
		"/apis/getkvmnodes",
		handler.GetKVMNodes,
	},
	Route{
		"getgkeclusters",
		"POST",
		"/apis/getgkeclusters",
		handler.GetGKEClusters,
	},
	Route{
		"gkechangenodecount",
		"POST",
		"/apis/gkechangenodecount",
		handler.GKEChangeNodeCount,
	},
	Route{
		"akschangevmss",
		"POST",
		"/apis/akschangevmss",
		handler.AKSChangeVMSS,
	},
	Route{
		"aksgetallres",
		"POST",
		"/apis/aksgetallres",
		handler.AKSGetAllResources,
	},
	Route{
		"stopaksnode",
		"POST",
		"/apis/stopaksnode",
		handler.StopAKSNode,
	},
	Route{
		"startaksnode",
		"POST",
		"/apis/startaksnode",
		handler.StartAKSNode,
	},
	Route{
		"addaksnode",
		"POST",
		"/apis/addaksnode",
		handler.AddAKSnode,
	},
	Route{
		"yamlapply",
		"POST",
		"/apis/yamlapply",
		YamlApply,
	},

	Route{
		"changeeksnode",
		"POST",
		"/apis/changeeksnode",
		ChangeEKSnode,
	},

	Route{
		"migration",
		"POST",
		"/apis/migration",
		Migration,
	},

	Route{
		"addec2node",
		"POST",
		"/apis/addec2node",
		Addec2node,
	},

	Route{
		"dashboard",
		"GET",
		"/apis/dashboard",
		Dashboard,
	},

	Route{
		"clusters",
		"GET",
		"/apis/clusters",
		handler.GetJoinedClusters,
	},
	Route{
		"joinableclusters",
		"GET",
		"/apis/joinableclusters",
		handler.GetJoinableClusters,
	},
	Route{
		"cluster-overview",
		"GET",
		"/apis/clusters/overview",
		handler.ClusterOverview,
	},

	Route{
		"clusterJoin",
		"POST",
		"/apis/clusters/join",
		handler.OpenMCPJoin,
	},

	Route{
		"clusterUnjoin",
		"POST",
		"/apis/clusters/unjoin",
		handler.OpenMCPUnjoin,
	},

	Route{
		"nodes",
		"GET",
		"/apis/clusters/{clusterName}/nodes",
		handler.NodesInCluster,
	},

	Route{
		"nodes",
		"GET",
		"/apis/nodes",
		handler.Nodes,
	},

	Route{
		"node-overview",
		"GET",
		"/apis/nodes/{nodeName}",
		handler.NodeOverview,
	},

	Route{
		"node-metric",
		"GET",
		"/apis/nodes_metric",
		handler.NodesMetric,
	},

	Route{
		"projects",
		"GET",
		"/apis/projects",
		handler.Projects,
	},

	Route{
		"projectOverview",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}",
		handler.GetProjectOverview,
	},

	Route{
		"AddProject",
		"POST",
		"/apis/clusters/projects/create",
		handler.AddProject,
	},

	Route{
		"deployments",
		"GET",
		"/apis/deployments",
		handler.GetDeployments,
	},
	Route{
		"deploymentsInProject",
		"GET",
		"/apis/clsuters/{clusterName}/projects/{projectName}/deployments",
		handler.GetDeploymentsInProject,
	},
	Route{
		"deploymentOverview",
		"GET",
		"/apis/clsuters/{clusterName}/projects/{projectName}/deployments/{deploymentName}",
		handler.GetDeploymentOverview,
	},
	Route{
		"replicaStatus",
		"GET",
		"/apis/clsuters/{clusterName}/projects/{projectName}/deployments/{deploymentName}/replica_status",
		handler.GetDeploymentReplicaStatus,
	},

	Route{
		"statefulsets",
		"GET",
		"/apis/statefulsets",
		handler.GetStatefulsets,
	},
	Route{
		"statefulsetsInProject",
		"GET",
		"/apis/clsuters/{clusterName}/projects/{projectName}/statefulsets",
		handler.GetStatefulsetsInProject,
	},
	Route{
		"statefulsetOverview",
		"GET",
		"/apis/clsuters/{clusterName}/projects/{projectName}/statefulsets/{statefulsetName}",
		handler.GetStatefulsetOverview,
	},

	Route{
		"dns",
		"GET",
		"/apis/dns",
		handler.Dns,
	},

	Route{
		"services",
		"GET",
		"/apis/services",
		handler.Services,
	},

	Route{
		"servicesInProject",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/services",
		handler.GetServicesInProject,
	},

	Route{
		"serviceOverview",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/services/{serviceName}",
		handler.GetServiceOverview,
	},

	Route{
		"ingress",
		"GET",
		"/apis/ingress",
		handler.Ingress,
	},

	Route{
		"ingressInProject",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/ingress",
		handler.GetIngressInProject,
	},

	Route{
		"ingressOverview",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/ingress/{ingressName}",
		handler.GetIngressOverview,
	},

	Route{
		"pods",
		"GET",
		"/apis/pods",
		handler.GetPods,
	},
	Route{
		"podOverview",
		"GET",
		"/apis/pods/{podName}",
		handler.GetPodOverview,
	},

	Route{
		"podPhysicalRes",
		"GET",
		"/apis/pods/{podName}/physicalResPerMin",
		handler.GetPodPhysicalRes,
	},

	Route{
		"vpa",
		"GET",
		"/apis/vpa",
		handler.GetVPAs,
	},

	Route{
		"hpa",
		"GET",
		"/apis/hpa",
		handler.GetHPAs,
	},

	Route{
		"podsInCluster",
		"GET",
		"/apis/clusters/{clusterName}/pods",
		handler.GetPodsInCluster,
	},

	Route{
		"podsInProject",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/pods",
		handler.GetPodsInProject,
	},

	Route{
		"pvcInProject",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/volumes",
		handler.GetVolumes,
	},

	Route{
		"pvcOverview",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/volumes/{volumeName}",
		handler.GetVolumeOverview,
	},

	Route{
		"secretInProject",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/secrets",
		handler.GetSecrets,
	},

	Route{
		"secretOverview",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/secrets/{secretName}",
		handler.GetSecretOverView,
	},

	Route{
		"configmapInProject",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/configmaps",
		handler.GetConfigmaps,
	},

	Route{
		"configmapOverview",
		"GET",
		"/apis/clusters/{clusterName}/projects/{projectName}/configmaps/{configmapName}",
		handler.GetConfigmapOverView,
	},

	Route{
		"settings",
		"GET",
		"/apis/policy/openmcp",
		handler.GetOpenmcpPolicy,
	},

	Route{
		"settings",
		"POST",
		"/apis/policy/openmcp/edit",
		handler.UpdateOpenmcpPolicy,
	},
}
