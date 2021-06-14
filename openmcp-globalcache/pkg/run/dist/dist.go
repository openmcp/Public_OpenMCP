package dist

import (
	"context"
	"errors"
	"fmt"
	"time"

	"openmcp/openmcp/openmcp-globalcache/pkg/utils"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

//GetNodeInfoList 노드리스트 가져오는 함수.
func (r *RegistryManager) GetNodeInfoList() ([]string, error) {

	if r.clientset == nil {
		return nil, errors.New("no cluster is specified")
	}

	//*v1.NodeList
	nodeList, err := r.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	//nodeName := "my-node"
	//pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
	//	FieldSelector: "spec.nodeName=" + nodeName,
	//})

	result := make([]string, len(nodeList.Items))
	for i, item := range nodeList.Items {
		result[i] = item.Name
	}

	return result, nil
}

//SetNodeLabelSync 선택된 클러스터의 전체 Sync 를 맞추는 함수.
func (r *RegistryManager) SetNodeLabelSync() error {

	if r.clientset == nil {
		return errors.New("no cluster is specified")
	}

	_, checkErr := r.checkNamespace()
	if checkErr != nil {
		return checkErr
	}

	var listErr error
	nodes, listErr := r.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if listErr != nil {
		return listErr
	}

	for _, node := range nodes.Items {
		//node name으로 node 정보 가져와서 내용 바꿔서 update 하는 기능.
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			//1. Node 정보를 가져온다.
			result, getErr := r.clientset.CoreV1().Nodes().Get(context.TODO(), node.Name, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}
			//2. 라벨을 추가해준다.
			modLabels := result.GetLabels()
			addLabelName := r.getLabelName(node.Name)
			modLabels[addLabelName] = "true"
			result.SetLabels(modLabels)
			//3. 업데이트.
			_, updateErr := r.clientset.CoreV1().Nodes().Update(context.TODO(), result, metav1.UpdateOptions{})
			return updateErr
		})
		if retryErr != nil {
			return retryErr
		}
		fmt.Println("Updated ..." + node.Name)
	}

	return nil
}

//DistributeRegistryAgent 선택된 클러스터의 모든 노드에 중계 서버를 배포하는 기능
func (r *RegistryManager) DistributeRegistryAgent() error {

	if r.clientset == nil {
		return errors.New("no cluster is specified")
	}

	_, checkErr := r.checkNamespace()
	if checkErr != nil {
		return checkErr
	}

	//get Node List
	var listErr error
	nodes, listErr := r.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if listErr != nil {
		return listErr
	}

	for _, node := range nodes.Items {
		//node name으로 node 정보 가져와서 내용 바꿔서 update 하는 기능.
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			//1. Node 정보를 가져온다.
			result, getErr := r.clientset.CoreV1().Nodes().Get(context.TODO(), node.Name, metav1.GetOptions{})
			if getErr != nil {
				panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
			}

			// pod 생성 가능한 노드만 관리.
			if r.canCreatedPod(result) {
				//2. 라벨을 추가해준다.
				modLabels := result.GetLabels()
				addLabelName := r.getLabelName(node.Name)
				modLabels[addLabelName] = "true"
				result.SetLabels(modLabels)
				//3. 업데이트.
				_, updateErr := r.clientset.CoreV1().Nodes().Update(context.TODO(), result, metav1.UpdateOptions{})
				if updateErr != nil {
					return updateErr
				}

				//3. 업데이트 된 노드를 대상으로 하는 Deployment 생성
				updateErr = r.CreateRepositoryAgent(node.Name)
				if updateErr != nil {
					return updateErr
				}
			}

			return nil
		})
		if retryErr != nil {
			return retryErr
		}
		fmt.Println("creating..." + node.Name)
	}

	return nil
}

//DeleteRegistryAgentAll : Deprecated / 선택된 클러스터의 모든 노드의 중계 서버를 삭제하는 기능
func (r *RegistryManager) DeleteRegistryAgentAll() error {

	if r.clientset == nil {
		return errors.New("no cluster is specified")
	}

	//get Node List
	var listErr error
	nodes, listErr := r.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if listErr != nil {
		return listErr
	}

	for _, node := range nodes.Items {
		//node name으로 node 정보 가져와서 내용 바꿔서 update 하는 기능.
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			//1. Node 정보를 가져온다.
			result, getErr := r.clientset.CoreV1().Nodes().Get(context.TODO(), node.Name, metav1.GetOptions{})
			if getErr != nil {
				panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
			}

			// pod 생성 가능한 노드만 관리.
			if r.canCreatedPod(result) {
				//2. 라벨을 변경해준다.
				modLabels := result.GetLabels()
				addLabelName := r.getLabelName(node.Name)
				modLabels[addLabelName] = "false"
				result.SetLabels(modLabels)
				//3. 업데이트.
				_, updateErr := r.clientset.CoreV1().Nodes().Update(context.TODO(), result, metav1.UpdateOptions{})
				if updateErr != nil {
					return updateErr
				}

				//3. 업데이트 된 노드를 대상으로 하는 Deployment 제거
				updateErr = r.DeleteRegistryAgent(node.Name)
				if updateErr != nil {
					return updateErr
				}
			}

			return nil
		})
		if retryErr != nil {
			return retryErr
		}
		fmt.Println("deleted ..." + node.Name)
	}

	return nil
}

/*
apiVersion: v1
kind: Pod
metadata:
	name: dood
spec:
	containers:
	- name: docker-cmds
		image: docker:19.03.8
		command: ["/bin/sh"]
		args: ["-c", "while true; do echo running; sleep 1000;done"]
		resources:
			requests:
				cpu: 10m
				memory: 256Mi
		volumeMounts:
		- mountPath: /var/run
			name: docker-sock
	volumes:
	- name: docker-sock
		hostPath:
			path: /var/run
*/
//CreateRepositoryAgent : Deprecated / 중계모듈을 배포하는 기능. 이미지는 openmcp/repository-agent 프로젝트에 있다. -> 이미지 위치 조정 중..
func (r *RegistryManager) CreateRepositoryAgent(nodeName string) error {
	appName := r.getAppName(nodeName)
	labelName := r.getLabelName(nodeName)

	deploymentsClient := r.clientset.AppsV1().Deployments(utils.ProjectNamespace)
	deployment, _ := deploymentsClient.Get(context.TODO(), appName, metav1.GetOptions{})
	if deployment.ObjectMeta.Name != "" {
		fmt.Printf("deployment exist : " + deployment.ObjectMeta.Name + "\n")
		return nil
	}

	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": appName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
					NodeSelector: map[string]string{
						labelName: "true",
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return nil
}

//DeleteRegistryAgent d
func (r *RegistryManager) DeleteRegistryAgent(nodeName string) error {
	appName := r.getAppName(nodeName)

	deploymentsClient := r.clientset.AppsV1().Deployments(utils.ProjectNamespace)
	deployment, _ := deploymentsClient.Get(context.TODO(), appName, metav1.GetOptions{})
	if deployment.ObjectMeta.Name == "" {
		fmt.Printf("deployment not exist : " + deployment.ObjectMeta.Name + "\n")
		return nil
	}

	// Delete Deployment
	fmt.Println("Delete deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentsClient.Delete(context.TODO(), appName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted deployment %q.\n", appName)

	return nil
}

//DeleteRegistryJob ㅇ
func (r *RegistryManager) DeleteRegistryJob(nodeName string) error {
	appName := r.getAppName(nodeName)

	jobClient := r.clientset.BatchV1().Jobs(utils.ProjectNamespace)
	job, _ := jobClient.Get(context.TODO(), appName, metav1.GetOptions{})
	if job.ObjectMeta.Name == "" {
		fmt.Printf("job not exist : " + job.ObjectMeta.Name + "\n")
		return nil
	}

	// Delete Deployment
	fmt.Println("Delete job...")
	deletePolicy := metav1.DeletePropagationForeground
	err := jobClient.Delete(context.TODO(), appName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted job %q.\n", appName)

	return nil
}

//CreatePushJobForCluster 원하는 클러스터에 push 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreatePushJobForCluster(imageName string, tag string) error {
	return r.CreateJobForCluster(imageName, tag, "push")
}

//CreatePullJobForCluster 원하는 클러스터에 Pull 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreatePullJobForCluster(imageName string, tag string) error {
	return r.CreateJobForCluster(imageName, tag, "pull")
}

//CreateJobForCluster 원하는 노드에 특정 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreateJobForCluster(imageName string, tag string, cmdType string) error {
	if r.clientset == nil {
		return errors.New("no cluster is specified")
	}

	_, checkErr := r.checkNamespace()
	if checkErr != nil {
		return checkErr
	}

	var listErr error
	nodes, listErr := r.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if listErr != nil {
		return listErr
	}

	for _, node := range nodes.Items {
		///node name으로 node 정보 가져와서 데이터를 추가해준다.
		//retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		//1. Node 정보를 가져온다.
		result, getErr := r.clientset.CoreV1().Nodes().Get(context.TODO(), node.Name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}

		//2. Job 생성한다. (삭제 메소드도 있음)
		createErr := r.CreateJob(result.Name, imageName, tag, cmdType)
		if createErr != nil {
			return createErr
		}
		//return createErr
		//})
		//if retryErr != nil {
		//	return retryErr
		//}
		//fmt.Println("Updated ..." + node.Name)
	}

	return nil
}

//CreatePushJob 원하는 노드에 Push 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreatePushJob(nodeName string, imageName string, tag string) error {
	return r.CreateJob(nodeName, imageName, tag, "push")
}

//CreatePullJob 원하는 노드에 Push 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreatePullJob(nodeName string, imageName string, tag string) error {
	return r.CreateJob(nodeName, imageName, tag, "pull")
}

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

//PullTest 원하는 노드에 특정 명령을 내리는 job을 생성하는 기능.
// func (r *RegistryManager) PullTest(nodeName string, imageName string, tag string) error {
// 	cmdType := "pull"
// 	appName := r.getAppName(nodeName)
// 	labelName := r.getLabelName(nodeName)

// 	jobClient := r.clientset.BatchV1().Jobs(utils.ProjectNamespace)

// 	//이미 존재할 때의 처리 방법.
// 	job, _ := jobClient.Get(appName, metav1.GetOptions{})
// 	if job.ObjectMeta.Name != "" {
// 		fmt.Printf("job exist : " + job.ObjectMeta.Name + "\n")
// 		return nil
// 	}

// 	imageFullName := utils.GlobalRepo.URI + "/" + imageName + ":" + tag
// 	cmd := ""
// 	switch cmdType {
// 	case "push":
// 		cmd = "docker push " + imageFullName
// 		cmd = utils.SetGlobalRegistryCommand(cmd)
// 	case "pull":
// 		cmd = "docker pull " + imageFullName
// 		cmd = utils.SetGlobalRegistryCommand(cmd)
// 	default:
// 	}

// 	//1. 해당 push 할 컨테이너 명을 찾는다.
// 	//2. commit 명령을 내린다.
// 	//commitCommand := "docker commit -a openmcp -m 'make " + imageName + "' 598bc3e5efe7 " + imageFullName

// 	job = r.getJobAPI(appName, labelName, cmd)

// 	// Create Deployment
// 	fmt.Println("Creating " + cmdType + " Job...")
// 	result, err := jobClient.Create(job)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("Created "+cmdType+" Job %q.\n", result.GetObjectMeta().GetName())

// 	afterFunc := func(old interface{}, new interface{}) {

// 		//이후 작업 나열
// 		newJob := new.(*batchv1.Job)
// 		//oldJob := old.(*batchv1.Job)
// 		deletePolicy := metav1.DeletePropagationForeground
// 		fmt.Printf("JobRunCheck : delete job \n")
// 		err := jobClient.Delete(newJob.Name, &metav1.DeleteOptions{
// 			PropagationPolicy: &deletePolicy,
// 		})
// 		if err != nil {
// 			fmt.Printf("JobRunCheck Error : %v\n", err)
// 		}

// 		fmt.Printf("JobRunCheck end \n")
// 		//return
// 		close(r.stopper)
// 	}
// 	r.JobRunCheck(batchv1.JobComplete, afterFunc)

// 	time.Sleep(time.Second * 5)
// 	return nil
// }

func int32Ptr(i int32) *int32 { return &i }
