package handler

type ClustersRes struct {
	Clusters []ClusterInfo `json:"clusters"`
}

type ClusterInfo struct {
	Name          string               `json:"name"`
	Region        string               `json:"region"`
	Zones         string               `json:"zone"`
	Status        string               `json:"status"`
	Provider      string               `json:"provider"`
	Nodes         int                  `json:"nodes"`
	Cpu           string               `json:"cpu"`
	Ram           string               `json:"ram"`
	Disk          string               `json:"disk"`
	Network       string               `json:"network"`
	ResourceUsage ClusterResourceUsage `json:"resourceUsage"`
}

type NodeRes struct {
	Nodes []NodeInfo `json:"nodes"`
}

type NodeInfo struct {
	Name          string `json:"name"`
	Cluster       string `json:"cluster"`
	Status        string `json:"status"`
	Role          string `json:"role"`
	SystemVersion string `json:"system_version"`
	Cpu           string `json:"cpu"`
	Ram           string `json:"memory"`
	Pods          string `json:"pods"`
	Provider      string `json:"provider"`
	Region        string `json:"region"`
	Zone          string `json:"zone"`
}

type NodeOverView struct {
	Info              NodeBasicInfo     `json:"basic_info"`
	KbNodeStatus      []NameStatus      `json:"kubernetes_node_status"`
	NodeResourceUsage NodeResourceUsage `json:"node_resource_usage"`
}

type NodeBasicInfo struct {
	Name            string `json:"name"`
	Status          string `json:"status"`
	Role            string `json:"role"`
	Kubernetes      string `json:"kubernetes"`
	KubernetesProxy string `json:"kubernetes_proxy"`
	IP              string `json:"ip"`
	OS              string `json:"os"`
	Docker          string `json:"docker"`
	CreatedTime     string `json:"created_time"`
	Taint           Taint  `json:"taint"`
	Provider        string `json:"provider"`
	Cluster         string `json:"cluster"`
}

type NodeResourceUsage2 struct {
	Cluster string `json:"cluster"`
	Node    string `json:"node"`
	Cpu     Unit   `json:"cpu"`
	Memory  Unit   `json:"memory"`
	Storage Unit   `json:"storage"`
}

type NodeResourceUsage struct {
	Cpu     Unit `json:"cpu"`
	Memory  Unit `json:"memory"`
	Storage Unit `json:"storage"`
	Pods    Unit `json:"pods"`
}

type Taint struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Taint string `json:"taint"`
}

type ProjectRes struct {
	Projects []ProjectInfo `json:"projects"`
}

type ProjectInfo struct {
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Cluster     string                 `json:"cluster"`
	CreatedTime string                 `json:"created_time"`
	Labels      map[string]interface{} `json:"labels"`
	UID         string                 `json:"uid"`
}

type ProjectOverview struct {
	Info              ProjectInfo           `json:"basic_info"`
	ProjectResource   []ProjectResourceType `json:"project_resource"`
	UsageTop5         UsageTop5             `json:"usage_top5"`
	PhysicalResources PhysicalResources     `json:"physical_resources"`
}

type ProjectResourceType struct {
	Name     string `json:"resource"`
	Total    int    `json:"total"`
	Abnormal int    `json:"abnormal"`
}

type UsageTop5 struct {
	CPU    []UsageType `json:"cpu"`
	Memory []UsageType `json:"memory"`
}

type UsageType struct {
	Name string `json:"name"`
	// Type  string `json:"type"`
	Usage string `json:"usage"`
}

type PodRes struct {
	Pods []PodInfo `json:"pods"`
}

type PodInfo struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Cluster     string `json:"cluster,omitempty"`
	Project     string `json:"project,omitempty"`
	PodIP       string `json:"pod_ip,omitempty"`
	Node        string `json:"node,omitempty"`
	NodeIP      string `json:"node_ip,omitempty"`
	Cpu         string `json:"cpu,omitempty"`
	Ram         string `json:"memory,omitempty"`
	CreatedTime string `json:"created_time,omitempty"`
}

type HPARes struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Cluster     string `json:"cluster"`
	Reference   string `json:"reference"`
	MinRepl     string `json:"min_repl"`
	MaxRepl     string `json:"max_repl"`
	CurrentRepl string `json:"current_repl"`
}

type VPARes struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Cluster    string `json:"cluster"`
	Reference  string `json:"reference"`
	UpdateMode string `json:"update_mode"`
}

type DeploymentRes struct {
	Deployments []DeploymentInfo `json:"deployments"`
}

type DeploymentInfo struct {
	Name        string                 `json:"name"`
	Status      string                 `json:"status,omitempty"`
	Cluster     string                 `json:"cluster,omitempty"`
	Project     string                 `json:"project"`
	Image       string                 `json:"image,omitempty"`
	CreatedTime string                 `json:"created_time"`
	Uid         string                 `json:"uid,omitempty"`
	Labels      map[string]interface{} `json:"labels"`
}

type DeploymentOverview struct {
	Info   DeploymentInfo `json:"basic_info"`
	Pods   []PodInfo      `json:"pods"`
	Ports  []PortInfo     `json:"ports"`
	Events []Event        `json:"events"`
	// PhysicalResources PhysicalResources `json:"physical_resources"`
	// ReplicaStatus     []ReplicaInfo     `json:"replica_status"`
}

type ReplicaStatus struct {
	Cluster       string `json:"cluster"`
	Project       string `json:"project"`
	Deployment    string `json:"deployment"`
	Replicas      int    `json:"replicas"`
	ReadyReplicas int    `json:"ready_replicas"`
	// UnavailableReplicas int    `unavailable_replicas`
}

type StatefulsetRes struct {
	Statefulsets []StatefulsetInfo `json:"statefulsets"`
}

type StatefulsetInfo struct {
	Name        string                 `json:"name"`
	Status      string                 `json:"status,omitempty"`
	Cluster     string                 `json:"cluster,omitempty"`
	Project     string                 `json:"project"`
	Image       string                 `json:"image,omitempty"`
	CreatedTime string                 `json:"created_time"`
	Uid         string                 `json:"uid,omitempty"`
	Labels      map[string]interface{} `json:"labels"`
}

type StatefulsetOverview struct {
	Info   StatefulsetInfo `json:"basic_info"`
	Pods   []PodInfo       `json:"pods"`
	Ports  []PortInfo      `json:"ports"`
	Events []Event         `json:"events"`
	// PhysicalResources PhysicalResources `json:"physical_resources"`
	// ReplicaStatus     []ReplicaInfo     `json:"replica_status"`
}

type DNSRes struct {
	DNS []DNSInfo `json:"dns"`
}

type DNSInfo struct {
	Name    string `json:"name"`
	Project string `json:"project"`
	DnsName string `json:"dns_name"`
	IP      string `json:"ip"`
}

type ReplicaInfo struct {
	Cluster string    `json:"cluster"`
	Pods    []PodInfo `json:"pods"`
}

type PortInfo struct {
	Name     string `json:"port_name"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

type ServicesRes struct {
	Services []ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Name        string `json:"name"`
	Cluster     string `json:"cluster"`
	Project     string `json:"project"`
	Type        string `json:"type"`
	Selector    string `json:"selector"`
	Port        string `json:"port"`
	CreatedTime string `json:"created_time"`
	ClusterIP   string `json:"cluster_ip"`
	ExternalIP  string `json:"external_ip"`
}

type ServiceOverview struct {
	Info   ServiceBasicInfo `json:"basic_info"`
	Pods   []PodInfo        `json:"pods"`
	Events []Event          `json:"events"`
}

type ServiceBasicInfo struct {
	Name            string `json:"name"`
	Project         string `json:"project"`
	Type            string `json:"type"`
	Cluster         string `json:"cluster"`
	ClusterIP       string `json:"cluster_ip"`
	ExternalIP      string `json:"external_ip"`
	SessionAffinity string `json:"session_affinity"`
	Selector        string `json:"selector"`
	Endpoints       string `json:"endpoints"`
	CreatedTime     string `json:"created_time"`
}

//Ingress List
type IngerssRes struct {
	Ingress []IngerssInfo `json:"ingress"`
}

type IngerssInfo struct {
	Name        string `json:"name"`
	Cluster     string `json:"cluster"`
	Project     string `json:"project"`
	Address     string `json:"address"`
	CreatedTime string `json:"created_time"`
}

//Ingress Overview
type IngressOverView struct {
	Info   IngerssInfo `json:"basic_info"`
	Rules  []Rules     `json:"rules"`
	Events []Event     `json:"events"`
}

type Rules struct {
	Domain   string `json:"domain"`
	Protocol string `json:"protocol"`
	Path     string `json:"path"`
	Services string `json:"services"`
	Port     string `json:"port"`
	Secret   string `json:"secret"`
}

//Cluster Overview
type ClusterOverView struct {
	Info             BasicInfo            `json:"basic_info"`
	PusageTop5       ProjectUsageTop5     `json:"project_usage_top5"`
	NusageTop5       NodeUsageTop5        `json:"node_usage_top5"`
	CUsage           ClusterResourceUsage `json:"cluster_resource_usage"`
	KubernetesStatus []NameStatus         `json:"kubernetes_status"`
	Events           []Event              `json:"events"`
}

type ClusterResourceUsage struct {
	Cpu     Unit `json:"cpu"`
	Memory  Unit `json:"memory"`
	Storage Unit `json:"storage"`
}

type Unit struct {
	Unit    string    `json:"unit"`
	NameVal []NameVal `json:"status"`
}

type NameVal struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type ProjectUsageTop5 struct {
	CPU    PairList `json:"cpu"`
	Memory PairList `json:"memory"`
}

type NodeUsageTop5 struct {
	CPU    PairList `json:"cpu"`
	Memory PairList `json:"memory"`
}

type Pair struct {
	Name  string  `json:"name"`
	Usage float64 `json:"usage"`
}

type PairList []Pair

type NameStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type BasicInfo struct {
	Name              string `json:"name"`
	Provider          string `json:"provider"`
	KubernetesVersion string `json:"kubernetes_version"`
	Status            string `json:"status"`
	Region            string `json:"region"`
	Zone              string `json:"zone"`
}

type Event struct {
	Project string `json:"project,omitempty"`
	Typenm  string `json:"type"`
	Reason  string `json:"reason"`
	Object  string `json:"object,omitempty"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

type PodOverviewInfo struct {
	Name              string `json:"name"`
	Status            string `json:"status"`
	Cluster           string `json:"cluster"`
	Project           string `json:"project"`
	PodIP             string `json:"pod_ip"`
	Node              string `json:"node"`
	NodeIP            string `json:"node_ip"`
	Namespace         string `json:"namespace"`
	TotalRestartCount string `json:"total_restart_count"`
	CreatedTime       string `json:"created_time"`
}

type PodOverviewContainer struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	RestartCount int    `json:"restart_count"`
	Port         string `json:"port"`
	Image        string `json:"image"`
}

type PodOverviewStatus struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	LastUpdate string `json:"last_update"`
	Reason     string `json:"reason"`
	Message    string `json:"message"`
}

type PodCPUUsageMin struct {
	CPU  float64 `json:"cpu"`
	Time string  `json:"time"`
}
type PodMemoryUsageMin struct {
	Memory float64 `json:"memory"`
	Time   string  `json:"time"`
}

type PodNetworkUsageMin struct {
	Unit string `json:"unit"`
	In   int    `json:"in"`
	Out  int    `json:"out"`
	Time string `json:"time"`
}

type PhysicalResources struct {
	CPU                []PodCPUUsageMin     `json:"cpu"`
	Memory             []PodMemoryUsageMin  `json:"memory"`
	PodNetworkUsageMin []PodNetworkUsageMin `json:"network"`
}

type PodOverviewRes struct {
	BasicInfo         PodOverviewInfo        `json:"basic_info"`
	Containers        []PodOverviewContainer `json:"containers"`
	Status            []PodOverviewStatus    `json:"pod_status"`
	PhysicalResources PhysicalResources      `json:"physical_resources"`
	Event             []Event                `json:"events"`
}
type VolumeRes struct {
	Volumes []VolumeInfo `json:"volumes"`
}

type VolumeInfo struct {
	Name         string `json:"name"`
	Project      string `json:"project"`
	Status       string `json:"status"`
	Capacity     string `json:"capacity"`
	CreatedTime  string `json:"created_time"`
	StorageClass string `json:"storage_class,omitempty"`
	AccessMode   string `json:"access_mode,omitempty"`
}

// VolumeOverview
type VolumeOverview struct {
	Info      VolumeInfo `json:"basic_info"`
	MountedBy []PodInfo  `json:"mounted_by"`
	Events    []Event    `json:"events"`
}

// GetVolumeOverview
// Secret List
type SecretRes struct {
	Secrets []SecretInfo `json:"secrets"`
}
type SecretInfo struct {
	Name        string `json:"name"`
	Project     string `json:"project"`
	Type        string `json:"type"`
	CreatedTime string `json:"created_time"`
}

//Secret Overview
type SecretOverView struct {
	Info SecretInfo `json:"basic_info"`
	Data []Data     `json:"data"`
}

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//Configmap List
type ConfigmapRes struct {
	Configmaps []ConfigmapInfo `json:"configmaps"`
}

type ConfigmapInfo struct {
	Name        string `json:"name"`
	Project     string `json:"project"`
	Keys        string `json:"keys"`
	CreatedTime string `json:"created_time"`
}

//Configmap Overview
type ConfigmapOverView struct {
	Info ConfigmapInfo `json:"basic_info"`
	Data []Data        `json:"data"`
}

type ManagedCluster struct {
	Name               string      `json:"name"`
	ResourceGroup      string      `json:"resourcegroup"`
	NodeResourceGrouop string      `json:"noderesourcegroup"`
	AgentPool          []AgentPool `json:"agentpools"`
	Location           string      `json:"location"`
	ProvisionState     string      `json:"pvstate"`
	// VmssNames          []string `json:"vmssnames"`
}

type AgentPool struct {
	Name     string `json:"name"`
	VmssName string `json:"vmssname"`
	Count    int32  `json:"nodecount"`
}

type OpenmcpPolicyRes struct {
	OpenmcpPolicy []OpenmcpPolicy `json:"openmcp_policy"`
}

type OpenmcpPolicy struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Value  string `json:"value"`
}

type EditPolicyRes struct {
	EditPolicy []interface{} `json:"edit_policy"`
}

type EditInfo struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}
