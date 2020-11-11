package dist

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

//JobRunCheck run 체크... sync를 확인한다.
func (r *RegistryManager) JobRunCheck(watchType batchv1.JobConditionType, afterFunc func(old interface{}, new interface{})) {
	fmt.Printf("start JobRunCheck\n")
	r.watchType = watchType
	r.afterFunc = afterFunc
	factory := informers.NewSharedInformerFactory(r.clientset, 0)
	r.informer = factory.Batch().V1().Jobs().Informer()
	r.stopper = make(chan struct{})
	//defer close(r.stopper)
	defer func() {
		recover() // recover 함수로 패닉 복구
		//fmt.Println(msg)
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
	//	fmt.Printf("It has the label!")
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
	//fmt.Println("-- old Job --")
	//fmt.Println(oldJob.Status.Conditions)
	//fmt.Println("-- new Job --")
	//fmt.Println(newJob.Status.Conditions)
	fmt.Print("work...")

	isComplete := false
	for _, condition := range newJob.Status.Conditions {
		if r.stopper == nil {
			fmt.Println("\n-- r.stopper is nil --")
			return
		}
		fmt.Println(r.stopper)
		if condition.Type == r.watchType {
			//if condition.Type == batchv1.JobComplete {
			// 성공하면
			isComplete = true
			fmt.Println("\n-- " + r.watchType + " --")
			fmt.Println(condition)
		} else {
			fmt.Print(".")
		}
	}

	if isComplete {
		fmt.Printf("%s job... job name : %s\n", r.watchType, newJob.Name)
		r.afterFunc(old, new)
	}
	return
}
