package dist

import (
	"context"
	"fmt"
	"time"

	"openmcp/openmcp/openmcp-globalcache/pkg/utils"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//CreateJob 원하는 노드에 특정 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreateJob(nodeName string, imageName string, tag string, cmdType string) error {
	appName := r.getAppName(nodeName)
	labelName := r.getLabelName(nodeName)

	jobClient := r.clientset.BatchV1().Jobs(utils.ProjectNamespace)

	//이미 존재할 때의 처리 방법.
	job, _ := jobClient.Get(context.TODO(), appName, metav1.GetOptions{})
	if job.ObjectMeta.Name != "" {
		fmt.Printf("job exist : " + job.ObjectMeta.Name + "\n")
		return nil
	}

	imageFullName := utils.GlobalRepo.URI + "/" + imageName + ":" + tag
	imageOriName := imageName + ":" + tag
	cmd := ""
	switch cmdType {
	case "push":
		cmd = "docker image tag " + imageOriName + " " + imageFullName + ";"
		cmd += "docker push " + imageFullName
		cmd = utils.SetGlobalRegistryCommand(cmd)
	case "pull":
		cmd = "docker pull " + imageFullName
		cmd = utils.SetGlobalRegistryCommand(cmd)
	default:
	}

	//1. 해당 push 할 컨테이너 명을 찾는다.
	//2. commit 명령을 내린다.
	//commitCommand := "docker commit -a openmcp -m 'make " + imageName + "' 598bc3e5efe7 " + imageFullName

	job = r.getJobAPI(appName, labelName, cmd)

	// Create Deployment
	fmt.Println("Creating " + cmdType + " Job...")
	result, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created "+cmdType+" Job %q.\n", result.GetObjectMeta().GetName())

	afterFunc := func(old interface{}, new interface{}) {

		//이후 작업 나열
		newJob := new.(*batchv1.Job)
		//oldJob := old.(*batchv1.Job)
		deletePolicy := metav1.DeletePropagationForeground
		fmt.Printf("JobRunCheck : delete job \n")
		err := jobClient.Delete(context.TODO(), newJob.Name, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		if err != nil {
			fmt.Printf("JobRunCheck Error : %v\n", err)
		}

		fmt.Printf("JobRunCheck end \n")
		//return
		//fmt.Print(r.stopper)
		defer close(r.stopper)
		//<-r.stopper
	}

	r.JobRunCheck(batchv1.JobComplete, afterFunc)

	time.Sleep(time.Second * 5)
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
