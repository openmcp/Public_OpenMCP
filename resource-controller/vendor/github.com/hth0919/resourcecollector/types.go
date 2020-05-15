package CRM




type ClusterInfo struct {
	MetricValue          []string           `protobuf:"bytes,1,rep,name=MetricValue,proto3" json:"MetricValue,omitempty"`
	Clustername          string             `protobuf:"bytes,2,opt,name=Clustername,proto3" json:"Clustername,omitempty"`
	KubeConfig           string             `protobuf:"bytes,3,opt,name=KubeConfig,proto3" json:"KubeConfig,omitempty"`
	AdminToken           string             `protobuf:"bytes,4,opt,name=AdminToken,proto3" json:"AdminToken,omitempty"`
	NodeList             []*NodeInfo        `protobuf:"bytes,5,rep,name=NodeList,proto3" json:"NodeList,omitempty"`
	ClusterMetricSum     map[string]float64 `protobuf:"bytes,6,rep,name=ClusterMetricSum,proto3" json:"ClusterMetricSum,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
	Host                 string             `protobuf:"bytes,9,opt,name=Host,proto3" json:"Host,omitempty"`
	Pods                 []string           `protobuf:"bytes,10,rep,name=Pods,proto3" json:"Pods,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}
type NodeInfo struct {
	NodeName             string             `protobuf:"bytes,1,opt,name=NodeName,proto3" json:"NodeName,omitempty"`
	PodList              []*PodInfo         `protobuf:"bytes,2,rep,name=PodList,proto3" json:"PodList,omitempty"`
	NodeMetricSum        map[string]float64 `protobuf:"bytes,3,rep,name=NodeMetricSum,proto3" json:"NodeMetricSum,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
	NodeCapacity         map[string]int64   `protobuf:"bytes,4,rep,name=NodeCapacity,proto3" json:"NodeCapacity,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	NodeAllocatable      map[string]int64   `protobuf:"bytes,5,rep,name=NodeAllocatable,proto3" json:"NodeAllocatable,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	GeoInfo              map[string]string  `protobuf:"bytes,6,rep,name=GeoInfo,proto3" json:"GeoInfo,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CpuCores             float64            `protobuf:"fixed64,7,opt,name=CpuCores,proto3" json:"CpuCores,omitempty"`
	MemoryTotal          float64            `protobuf:"fixed64,8,opt,name=MemoryTotal,proto3" json:"MemoryTotal,omitempty"`
	ScrapeError          float64            `protobuf:"fixed64,9,opt,name=ScrapeError,proto3" json:"ScrapeError,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

type PodInfo struct {
	PodName              string             `protobuf:"bytes,1,opt,name=PodName,proto3" json:"PodName,omitempty"`
	PodNamespace         string             `protobuf:"bytes,2,opt,name=PodNamespace,proto3" json:"PodNamespace,omitempty"`
	PodMetrics           map[string]float64 `protobuf:"bytes,3,rep,name=PodMetrics,proto3" json:"PodMetrics,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

type ReturnValue struct {
	Tick                 int64    `protobuf:"varint,1,opt,name=Tick,proto3" json:"Tick,omitempty"`
	ClusterName          string   `protobuf:"bytes,2,opt,name=ClusterName,proto3" json:"ClusterName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}