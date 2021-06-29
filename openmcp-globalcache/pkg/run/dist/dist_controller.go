package dist

import (
	"fmt"
	"openmcp/openmcp/omcplog"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

//JobRunCheck run 체크... sync를 확인한다.
func (r *RegistryManager) JobRunCheck(watchType batchv1.JobConditionType, afterFunc func(old interface{}, new interface{})) {
	omcplog.V(4).Info("start JobRunCheck\n")
	r.watchType = watchType
	r.afterFunc = afterFunc
	factory := informers.NewSharedInformerFactory(r.clientset, 0)
	r.informer = factory.Batch().V1().Jobs().Informer()
	r.stopper = make(chan struct{})
	//defer close(r.stopper)
	defer func() {
		recover() // recover 함수로 패닉 복구
		//omcplog.V(3).Info(msg)
	}()
	defer runtime.HandleCrash()
	r.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		//AddFunc: onAdd,
		UpdateFunc: func(old, new interface{}) {
			r.onUpdate(old, new)
		},
		//DeleteFunc: controller.handleObject,controller.handleObject,
	})
	go r.informer.Run(r.stopper)

	if !cache.WaitForCacheSync(r.stopper, r.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-r.stopper

}

// onAdd is the function executed when the kubernetes informer notified the
// presence of a new kubernetes node in the cluster
func (r *RegistryManager) onAdd(obj interface{}) {
	// Cast the obj as node
	// node := obj.(*corev1.Node)
	//_, ok := node.GetLabels()[K8S_LABEL_AWS_REGION]
	//if ok {
	//	omcplog.V(3).Info("It has the label!")
	//}
}

// onUpdate
func (r *RegistryManager) onUpdate(old interface{}, new interface{}) {
	// Cast the obj as node

	newJob := new.(*batchv1.Job)
	oldJob := old.(*batchv1.Job)
	if newJob.ResourceVersion == oldJob.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return
	}
	//omcplog.V(3).Info("-- old Job --")
	//omcplog.V(3).Info(oldJob.Status.Conditions)
	//omcplog.V(3).Info("-- new Job --")
	//omcplog.V(3).Info(newJob.Status.Conditions)
	omcplog.V(3).Info("work...")

	isComplete := false
	for _, condition := range newJob.Status.Conditions {
		if r.stopper == nil {
			omcplog.V(3).Info("\n-- r.stopper is nil --")
			return
		}
		omcplog.V(3).Info(r.stopper)
		if condition.Type == r.watchType {
			//if condition.Type == batchv1.JobComplete {
			// 성공하면
			isComplete = true
			omcplog.V(3).Info("\n-- " + r.watchType + " --")
			omcplog.V(3).Info(condition)
		} else {
			omcplog.V(3).Info(".")
		}
	}

	if isComplete {
		omcplog.V(3).Info(string(r.watchType) + "job... job name : " + newJob.Name)
		r.afterFunc(old, new)
	}
	return
}
