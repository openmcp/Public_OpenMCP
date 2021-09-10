package v1alpha1

import (
	policyv1alpha1 "openmcp/openmcp/apis/policy/v1alpha1"

	"istio.io/api/meta/v1alpha1"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	hpav2beta1 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenMCPDeploymentSpec defines the desired state of OpenMCPDeployment
// +k8s:openapi-gen=true
type OpenMCPDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// changes
	Template OpenMCPDeploymentTemplate `json:"template" protobuf:"bytes,3,opt,name=template"`

	// Added
	Replicas int32               `json:"replicas" protobuf:"varint,1,opt,name=replicas"`
	Clusters []string            `json:"clusters,omitempty" protobuf:"bytes,11,opt,name=clusters"`
	Labels   map[string]string   `json:"labels,omitempty" protobuf:"bytes,11,opt,name=labels"`
	Affinity map[string][]string `json:"affinity,omitempty" protobuf:"bytes,3,opt,name=affinity"`
	Policy   map[string]string   `json:"policy,omitempty" protobuf:"bytes,3,opt,name=policy"`

	//Placement
}

type OpenMCPDeploymentTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// changed
	Spec   OpenMCPDeploymentTemplateSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status appsv1.DeploymentStatus       `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type OpenMCPDeploymentTemplateSpec struct {
	Replicas *int32                `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
	Selector *metav1.LabelSelector `json:"selector" protobuf:"bytes,2,opt,name=selector"`
	// changed
	Template                OpenMCPPodTemplateSpec    `json:"template" protobuf:"bytes,3,opt,name=template"`
	Strategy                appsv1.DeploymentStrategy `json:"strategy,omitempty" patchStrategy:"retainKeys" protobuf:"bytes,4,opt,name=strategy"`
	MinReadySeconds         int32                     `json:"minReadySeconds,omitempty" protobuf:"varint,5,opt,name=minReadySeconds"`
	RevisionHistoryLimit    *int32                    `json:"revisionHistoryLimit,omitempty" protobuf:"varint,6,opt,name=revisionHistoryLimit"`
	Paused                  bool                      `json:"paused,omitempty" protobuf:"varint,7,opt,name=paused"`
	ProgressDeadlineSeconds *int32                    `json:"progressDeadlineSeconds,omitempty" protobuf:"varint,9,opt,name=progressDeadlineSeconds"`
}

type OpenMCPPodTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// changed
	Spec OpenMCPPodSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type OpenMCPPodSpec struct {
	Volumes []corev1.Volume `json:"volumes,omitempty" patchStrategy:"merge,retainKeys" patchMergeKey:"name" protobuf:"bytes,1,rep,name=volumes"`
	// changes
	InitContainers []OpenMCPContainer `json:"initContainers,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,20,rep,name=initContainers"`
	// changes
	Containers                    []OpenMCPContainer            `json:"containers" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=containers"`
	RestartPolicy                 corev1.RestartPolicy          `json:"restartPolicy,omitempty" protobuf:"bytes,3,opt,name=restartPolicy,casttype=RestartPolicy"`
	TerminationGracePeriodSeconds *int64                        `json:"terminationGracePeriodSeconds,omitempty" protobuf:"varint,4,opt,name=terminationGracePeriodSeconds"`
	ActiveDeadlineSeconds         *int64                        `json:"activeDeadlineSeconds,omitempty" protobuf:"varint,5,opt,name=activeDeadlineSeconds"`
	DNSPolicy                     corev1.DNSPolicy              `json:"dnsPolicy,omitempty" protobuf:"bytes,6,opt,name=dnsPolicy,casttype=DNSPolicy"`
	NodeSelector                  map[string]string             `json:"nodeSelector,omitempty" protobuf:"bytes,7,rep,name=nodeSelector"`
	ServiceAccountName            string                        `json:"serviceAccountName,omitempty" protobuf:"bytes,8,opt,name=serviceAccountName"`
	DeprecatedServiceAccount      string                        `json:"serviceAccount,omitempty" protobuf:"bytes,9,opt,name=serviceAccount"`
	AutomountServiceAccountToken  *bool                         `json:"automountServiceAccountToken,omitempty" protobuf:"varint,21,opt,name=automountServiceAccountToken"`
	NodeName                      string                        `json:"nodeName,omitempty" protobuf:"bytes,10,opt,name=nodeName"`
	HostNetwork                   bool                          `json:"hostNetwork,omitempty" protobuf:"varint,11,opt,name=hostNetwork"`
	HostPID                       bool                          `json:"hostPID,omitempty" protobuf:"varint,12,opt,name=hostPID"`
	HostIPC                       bool                          `json:"hostIPC,omitempty" protobuf:"varint,13,opt,name=hostIPC"`
	ShareProcessNamespace         *bool                         `json:"shareProcessNamespace,omitempty" protobuf:"varint,27,opt,name=shareProcessNamespace"`
	SecurityContext               *corev1.PodSecurityContext    `json:"securityContext,omitempty" protobuf:"bytes,14,opt,name=securityContext"`
	ImagePullSecrets              []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,15,rep,name=imagePullSecrets"`
	Hostname                      string                        `json:"hostname,omitempty" protobuf:"bytes,16,opt,name=hostname"`
	Subdomain                     string                        `json:"subdomain,omitempty" protobuf:"bytes,17,opt,name=subdomain"`
	Affinity                      *corev1.Affinity              `json:"affinity,omitempty" protobuf:"bytes,18,opt,name=affinity"`
	SchedulerName                 string                        `json:"schedulerName,omitempty" protobuf:"bytes,19,opt,name=schedulerName"`
	Tolerations                   []corev1.Toleration           `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	HostAliases                   []corev1.HostAlias            `json:"hostAliases,omitempty" patchStrategy:"merge" patchMergeKey:"ip" protobuf:"bytes,23,rep,name=hostAliases"`
	PriorityClassName             string                        `json:"priorityClassName,omitempty" protobuf:"bytes,24,opt,name=priorityClassName"`
	Priority                      *int32                        `json:"priority,omitempty" protobuf:"bytes,25,opt,name=priority"`
	DNSConfig                     *corev1.PodDNSConfig          `json:"dnsConfig,omitempty" protobuf:"bytes,26,opt,name=dnsConfig"`
	ReadinessGates                []corev1.PodReadinessGate     `json:"readinessGates,omitempty" protobuf:"bytes,28,opt,name=readinessGates"`
	RuntimeClassName              *string                       `json:"runtimeClassName,omitempty" protobuf:"bytes,29,opt,name=runtimeClassName"`
	EnableServiceLinks            *bool                         `json:"enableServiceLinks,omitempty" protobuf:"varint,30,opt,name=enableServiceLinks"`
}

type OpenMCPContainer struct {
	Name       string                 `json:"name" protobuf:"bytes,1,opt,name=name"`
	Image      string                 `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	Command    []string               `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	Args       []string               `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
	WorkingDir string                 `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
	Ports      []corev1.ContainerPort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,6,rep,name=ports"`
	EnvFrom    []corev1.EnvFromSource `json:"envFrom,omitempty" protobuf:"bytes,19,rep,name=envFrom"`
	Env        []corev1.EnvVar        `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,7,rep,name=env"`
	// changes
	Resources                OpenMCPResourceRequirements     `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
	VolumeMounts             []corev1.VolumeMount            `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,9,rep,name=volumeMounts"`
	VolumeDevices            []corev1.VolumeDevice           `json:"volumeDevices,omitempty" patchStrategy:"merge" patchMergeKey:"devicePath" protobuf:"bytes,21,rep,name=volumeDevices"`
	LivenessProbe            *corev1.Probe                   `json:"livenessProbe,omitempty" protobuf:"bytes,10,opt,name=livenessProbe"`
	ReadinessProbe           *corev1.Probe                   `json:"readinessProbe,omitempty" protobuf:"bytes,11,opt,name=readinessProbe"`
	Lifecycle                *corev1.Lifecycle               `json:"lifecycle,omitempty" protobuf:"bytes,12,opt,name=lifecycle"`
	TerminationMessagePath   string                          `json:"terminationMessagePath,omitempty" protobuf:"bytes,13,opt,name=terminationMessagePath"`
	TerminationMessagePolicy corev1.TerminationMessagePolicy `json:"terminationMessagePolicy,omitempty" protobuf:"bytes,20,opt,name=terminationMessagePolicy,casttype=TerminationMessagePolicy"`
	ImagePullPolicy          corev1.PullPolicy               `json:"imagePullPolicy,omitempty" protobuf:"bytes,14,opt,name=imagePullPolicy,casttype=PullPolicy"`
	SecurityContext          *corev1.SecurityContext         `json:"securityContext,omitempty" protobuf:"bytes,15,opt,name=securityContext"`
	Stdin                    bool                            `json:"stdin,omitempty" protobuf:"varint,16,opt,name=stdin"`
	StdinOnce                bool                            `json:"stdinOnce,omitempty" protobuf:"varint,17,opt,name=stdinOnce"`
	TTY                      bool                            `json:"tty,omitempty" protobuf:"varint,18,opt,name=tty"`
}

type OpenMCPResourceRequirements struct {
	Limits   corev1.ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList,castkey=ResourceName"`
	Requests corev1.ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList,castkey=ResourceName"`
	// Added
	Needs isNeedResourceList `json:"needs,omitempty" protobuf:"bytes,2,rep,name=needs,casttype=isNeedResourceList,castkey=ResourceName"`
}

type isNeedResourceList map[corev1.ResourceName]bool

// OpenMCPDeploymentStatus defines the observed state of OpenMCPDeployment
// +k8s:openapi-gen=true
type OpenMCPDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas                  int32                 `json:"replicas"`
	ClusterMaps               map[string]int32      `json:"clusterMaps"`
	LastSpec                  OpenMCPDeploymentSpec `json:"lastSpec"`
	SchedulingNeed            bool                  `json:"schedulingNeed"`
	SchedulingComplete        bool                  `json:"schedulingComplete"`
	CreateSyncRequestComplete bool                  `json:"createSyncRequestComplete"`
	SyncRequestName           string                `json:"syncRequestName"`
	CheckSubResource          bool                  `json:"checkSubResource" protobuf:"bytes,3,opt,name=checkSubResource"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPDeployment is the Schema for the openmcpdeployments API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   OpenMCPDeploymentSpec   `json:"spec,omitempty"`
	Status OpenMCPDeploymentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPDeploymentList contains a list of OpenMCPDeployment
type OpenMCPDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPDeployment `json:"items"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenMCPIngressSpec defines the desired state of OpenMCPIngress
// +k8s:openapi-gen=true
type OpenMCPIngressSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	IngressForClientFrom string          `json:"ingressForClientFrom" protobuf:"bytes,1,opt,name=ingressForClientFrom"`
	Template             extv1b1.Ingress `json:"template" protobuf:"bytes,3,opt,name=template"`
}

// OpenMCPIngressStatus defines the observed state of OpenMCPIngress
// +k8s:openapi-gen=true
type OpenMCPIngressStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps      map[string]int32   `json:"clusterMaps"`
	ChangeNeed       bool               `json:"changeNeed"`
	LastSpec         OpenMCPIngressSpec `json:"lastSpec"`
	CheckSubResource bool               `json:"checkSubResource" protobuf:"bytes,3,opt,name=checkSubResource"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPIngress is the Schema for the openmcpingresss API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPIngressSpec   `json:"spec,omitempty"`
	Status OpenMCPIngressStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPIngressList contains a list of OpenMCPIngress
type OpenMCPIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPIngress `json:"items"`
}

// OpenMCPServiceSpec defines the desired state of OpenMCPService
// +k8s:openapi-gen=true
type OpenMCPServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	LabelSelector map[string]string `json:"labelselector" protobuf:"bytes,1,opt,name=labelselector"`
	Template      corev1.Service    `json:"template" protobuf:"bytes,2,opt,name=template"`
}

// OpenMCPServiceStatus defines the observed state of OpenMCPService
// +k8s:openapi-gen=true
type OpenMCPServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps      map[string]int32   `json:"clusterMaps"`
	LastSpec         OpenMCPServiceSpec `json:"lastSpec"`
	ChangeNeed       bool               `json:"changeNeed"`
	CheckSubResource bool               `json:"checkSubResource" protobuf:"bytes,3,opt,name=checkSubResource"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPService is the Schema for the openmcpservices API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPServiceSpec   `json:"spec,omitempty"`
	Status OpenMCPServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPServiceList contains a list of OpenMCPService
type OpenMCPServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPService `json:"items"`
}

type OpenMCPHybridAutoScalerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	MainController string         `json:"mainController" protobuf:"bytes,1,opt,name=maincontroller"`
	ScalingOptions ScalingOptions `json:"scalingOptions,omitempty" protobuf:"bytes,2,opt,name=scalingoptions"`
}

type ScalingOptions struct {
	CpaTemplate CpaTemplate                        `json:"cpaTemplate,omitempty" protobuf:"bytes,1,opt,name=cpatemplate"`
	HpaTemplate hpav2beta1.HorizontalPodAutoscaler `json:"hpaTemplate,omitempty" protobuf:"bytes,2,opt,name=hpatemplate"`
	VpaTemplate string                             `json:"vpaTemplate,omitempty" protobuf:"bytes,3,opt,name=vpatemplate"`
	//VpaTemplate vpav1beta2.VerticalPodAutoscaler `json:"vpaTemplate" protobuf:"bytes,3,opt,name=vpaTemplate"`
}

type CpaTemplate struct {
	ScaleTargetRef ScaleTargetRef `json:"scaleTargetRef,omitempty" protobuf:"bytes,1,opt,name=scaletargetref"`
	MinReplicas    int32          `json:"minReplicas,omitempty" protobuf:"varint,2,opt,name=minreplicas"`
	MaxReplicas    int32          `json:"maxReplicas,omitempty" protobuf:"varint,3,opt,name=maxreplicas"`
}

type ScaleTargetRef struct {
	// Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds"
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	// Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
	// API version of the referent
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,3,opt,name=apiVersion"`
}

// OpenMCPHybridAutoScalerStatus defines the observed state of OpenMCPHybridAutoScaler
// +k8s:openapi-gen=true
type OpenMCPHybridAutoScalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	//Nodes []string `json:"nodes"`
	LastSpec         OpenMCPHybridAutoScalerSpec      `json:"lastSpec"`
	Policies         []policyv1alpha1.OpenMCPPolicies `json:"policies"`
	RebalancingCount map[string]int32                 `json:"rebalancingCount"`
	SyncRequestName  string                           `json:"syncRequestName"`
	ChangeNeed       bool                             `json:"changeNeed"`
	CheckSubResource bool                             `json:"checkSubResource"`
	ClusterMinMaps   map[string]int32                 `json:"clusterMinMaps"`
	ClusterMaxMaps   map[string]int32                 `json:"clusterMaxMaps"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPHybridAutoScaler is the Schema for the openmcphybridautoscalers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPHybridAutoScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPHybridAutoScalerSpec   `json:"spec,omitempty"`
	Status OpenMCPHybridAutoScalerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPHybridAutoScalerList contains a list of OpenMCPHybridAutoScaler
type OpenMCPHybridAutoScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPHybridAutoScaler `json:"items"`
}

// OpenMCPConfigMapSpec defines the desired state of OpenMCPConfigMap
// +k8s:openapi-gen=true
type OpenMCPConfigMapSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template corev1.ConfigMap `json:"template" protobuf:"bytes,3,opt,name=template"`
}

// OpenMCPConfigMapStatus defines the observed state of OpenMCPConfigMap
// +k8s:openapi-gen=true
type OpenMCPConfigMapStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps      map[string]int32     `json:"clusterMaps"`
	SyncRequestName  string               `json:"syncRequestName"`
	LastSpec         OpenMCPConfigMapSpec `json:"lastSpec"`
	CheckSubResource bool                 `json:"checkSubResource" protobuf:"bytes,3,opt,name=checkSubResource"`
}

// OpenMCPConfigMap is the Schema for the openmcpconfigmaps API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPConfigMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPConfigMapSpec   `json:"spec,omitempty"`
	Status OpenMCPConfigMapStatus `json:"status,omitempty"`
}

// OpenMCPConfigMapList contains a list of OpenMCPConfigMap
type OpenMCPConfigMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPConfigMap `json:"items"`
}

type OpenMCPSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template corev1.Secret `json:"template" protobuf:"bytes,3,opt,name=template"`
}

// OpenMCPSecretStatus defines the observed state of OpenMCPSecret
// +k8s:openapi-gen=true
type OpenMCPSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps      map[string]int32  `json:"clusterMaps"`
	LastSpec         OpenMCPSecretSpec `json:"lastSpec"`
	CheckSubResource bool              `json:"checkSubResource" protobuf:"bytes,3,opt,name=checkSubResource"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPSecret is the Schema for the openmcpsecrets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPSecretSpec   `json:"spec,omitempty"`
	Status OpenMCPSecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPSecretList contains a list of OpenMCPSecret
type OpenMCPSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPSecret `json:"items"`
}

type OpenMCPJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template batchv1.Job `json:"template" protobuf:"bytes,3,opt,name=template"`
}

// OpenMCPJobStatus defines the observed state of OpenMCPJob
// +k8s:openapi-gen=true
type OpenMCPJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps      map[string]int32 `json:"clusterMaps"`
	LastSpec         OpenMCPJobSpec   `json:"lastSpec"`
	CheckSubResource bool             `json:"checkSubResource" protobuf:"bytes,3,opt,name=checkSubResource"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPJob is the Schema for the openmcpjobs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPJobSpec   `json:"spec,omitempty"`
	Status OpenMCPJobStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPJobList contains a list of OpenMCPJob
type OpenMCPJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPJob `json:"items"`
}

type OpenMCPNamespaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template appsv1.Deployment `json:"template" protobuf:"bytes,3,opt,name=template"`
	//Placement
}

// OpenMCPNamespaceStatus defines the observed state of OpenMCPNamespace
// +k8s:openapi-gen=true
type OpenMCPNamespaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps      map[string]int32     `json:"clusterMaps"`
	LastSpec         OpenMCPNamespaceSpec `json:"lastSpec"`
	CheckSubResource bool                 `json:"checkSubResource" protobuf:"bytes,3,opt,name=checkSubResource"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPNamespace is the Schema for the openmcpnamespaces API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPNamespaceSpec   `json:"spec,omitempty"`
	Status OpenMCPNamespaceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPNamespaceList contains a list of OpenMCPNamespace
type OpenMCPNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPNamespace `json:"items"`
}

//func init() {
//	SchemeBuilder.Register(&OpenMCPDeployment{}, &OpenMCPDeploymentList{})
//	SchemeBuilder.Register(&OpenMCPIngress{}, &OpenMCPIngressList{})
//	SchemeBuilder.Register(&OpenMCPService{}, &OpenMCPServiceList{})
//	SchemeBuilder.Register(&OpenMCPHybridAutoScaler{}, &OpenMCPHybridAutoScalerList{})
//	SchemeBuilder.Register(&OpenMCPPolicyEngine{}, &OpenMCPPolicyEngineList{})
//
//}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenMCPVirtualServiceSpec defines the desired state of OpenMCPVirtualService
// +k8s:openapi-gen=true
type OpenMCPVirtualService struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec networkingv1alpha3.VirtualService `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	Status v1alpha1.IstioStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualServiceList is a collection of VirtualServices.
type OpenMCPVirtualServiceList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []OpenMCPVirtualService `json:"items" protobuf:"bytes,2,rep,name=items"`
}
