package openmcpscheduler

import (
	// "fmt"
	kubesource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog"
	ketiv1alpha1 "resource-controller/apis/keti/v1alpha1"
	"time"
	// corev1 "k8s.io/api/core/v1"

	// _ "github.com/influxdata/influxdb1-client"  // this is important because of the buf in go mod
	// client "github.com/influxdata/influxdb1-client/v2"
	ketiresource "resource-controller/controllers/openmcpscheduler/pkg/controller/resourceInfo"
	// collector "resource-controller/controllers/openmcpscheduler/pkg/controller/resourceCollector"
)

func (cm *ClusterManager) Scheduling(pod *ketiv1alpha1.OpenMCPDeployment) map[string]int32 {
	klog.Infof("*********** Scheduling ***********")

	// return value (ex. cluster1:2, cluster2:1)
	replicas_cluster := map[string]int32{}

	// get data from influxdb
	// getResources()

	// make resource to schedule pod into cluster
	newResource := newResourceFromPod(pod)
	klog.Infof("*********** %s ***********", newResource)

	startTime := time.Now()
	elapsedTime := time.Since(startTime)
	klog.Infof("*********** %s ***********", elapsedTime)

	return replicas_cluster
}

func newResourceFromPod(dep *ketiv1alpha1.OpenMCPDeployment) *ketiresource.Resource {
	res := &ketiresource.Resource{}

	//_, container := range dep.Spec.Template.Spec.Template.Spec.Containers
	cpu := kubesource.MustParse(dep.Spec.Template.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	memory := kubesource.MustParse(dep.Spec.Template.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())

	res.MilliCPU = cpu.MilliValue()
	res.Memory = memory.Value()

	return res
}
