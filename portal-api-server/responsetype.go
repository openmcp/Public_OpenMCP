package main

type NameVal struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Attributes struct {
	Status string `json:"status"`
	Region string `json:"region"`
	Zone   string `json:"zone"`
	// Attributes struct {
	// 	Status string `json:"status"`
	// } `json:"attributes"`
}
type ChildNode struct {
	Name       string     `json:"name"`
	Attributes Attributes `json:"attributes"`
}

type Region struct {
	Name       string      `json:"name"`
	Attributes Attributes  `json:"attributes"`
	Children   []ChildNode `json:"children"`
}

type JoinedClusters struct {
	Name       string      `json:"name"`
	Attributes Attributes  `json:"attributes"`
	Children   []ChildNode `json:"children"`
}

type DashboardRes struct {
	Clusters struct {
		ClustersCnt    int       `json:"counts"`
		ClustersStatus []NameVal `json:"status"`
	} `json:"clusters"`
	Nodes struct {
		NodesCnt    int       `json:"counts"`
		NodesStatus []NameVal `json:"status"`
	} `json:"nodes"`
	Pods struct {
		PodsCnt    int       `json:"counts"`
		PodsStatus []NameVal `json:"status"`
	} `json:"pods"`
	Projects struct {
		ProjectsCnt    int       `json:"counts"`
		ProjectsStatus []NameVal `json:"status"`
	} `json:"projects"`
	Regions        []Region       `json:"regions"`
	JoinedClusters JoinedClusters `json:"joined_clusters"`
}

type ManagedCluster struct {
	Name               string      `json:"name"`
	ResourceGroup      string      `json:"resourcegroup"`
	NodeResourceGrouop string      `json:"noderesourcegroup"`
	AgentPool          []AgentPool `json:"agentpools"`
	Location           string      `json:"location"`
	// VmssNames          []string `json:"vmssnames"`
}

type AgentPool struct {
	Name     string `json:"name"`
	VmssName string `json:"vmssname"`
}

type EKSInstance struct {
	InstanceId string `json:"instance_id"`
}

type EKSNodegroup struct {
	NGName           string        `json:"name"`
	InstanceType     string        `json:"instance_type"`
	DesiredSize      int64         `json:"desired_size"`
	MaxSize          int64         `json:"max_size"`
	MinSize          int64         `json:"min_size"`
	AutoscalingGroup string        `json:"auto_scaling_group"`
	Instances        []EKSInstance `json:"intances"`
}

type EKSCluster struct {
	ClusterName string         `json:"name"`
	Nodegroups  []EKSNodegroup `json:"nodegroups"`
}

// type ClusterData struct {
// 	Name       string            `json:"name"`
// 	attributes ClusterAttributes `json:"attributes`
// }

// type ClusterAttributes struct {
// 	status string `json:"status"`
// 	region string `json:"region"`
// 	zone   string `json:"zone"`
// }
