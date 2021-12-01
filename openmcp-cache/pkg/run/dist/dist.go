package dist

import (
	"context"
	"errors"
	"strings"
	"time"

	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-cache/pkg/utils"

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

	result := make([]string, len(nodeList.Items))
	for i, item := range nodeList.Items {
		result[i] = item.Name
	}

	return result, nil
}

//SetNodeLabelSync 선택된 클러스터의 전체 Sync 를 맞추는 함수.
func (r *RegistryManager) SetNodeLabelSync() error {
	omcplog.V(3).Info("Node SetNodeLabelSync")

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
		omcplog.V(3).Info(" error ", listErr)
		return listErr
	}

	for _, node := range nodes.Items {
		//node name으로 node 정보 가져와서 내용 바꿔서 update 하는 기능.
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

			omcplog.V(4).Info("   get Node : " + node.Name)
			//1. Node 정보를 가져온다.
			result, getErr := r.clientset.CoreV1().Nodes().Get(context.TODO(), node.Name, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}

			omcplog.V(4).Info("   get Node mod...")
			omcplog.V(4).Info("   addLabelName : " + r.getLabelName(node.Name))

			//2. 라벨을 추가해준다.
			modLabels := result.GetLabels()
			addLabelName := r.getLabelName(node.Name)
			modLabels[addLabelName] = "true"
			result.SetLabels(modLabels)

			omcplog.V(4).Info("   get Node update")
			//3. 업데이트.
			_, updateErr := r.clientset.CoreV1().Nodes().Update(context.TODO(), result, metav1.UpdateOptions{})
			return updateErr
		})
		if retryErr != nil {
			return retryErr
		}
		omcplog.V(3).Info(" Node Sync Updated ..." + node.Name)
	}
	omcplog.V(3).Info("Node SetNodeLabelSync end")

	return nil
}

//DeleteRegistryJob ㅇ
func (r *RegistryManager) DeleteJob() error {
	jobClient := r.clientset.BatchV1().Jobs(utils.ProjectNamespace)

	// Delete Deployment
	omcplog.V(3).Info("Delete job...")
	deletePolicy := metav1.DeletePropagationForeground
	err := jobClient.DeleteCollection(context.TODO(), metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy}, metav1.ListOptions{LabelSelector: "app=openmcp-cache"})
	if err != nil {
		return err
	}
	omcplog.V(3).Info("Deleted job %q.\n")

	return nil
}

//CreatePushJobForCluster 원하는 클러스터에 push 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreatePushJobForCluster(imageName string) error {
	return r.CreateJobForCluster(imageName, "push")
}

//CreatePullJobForCluster 원하는 클러스터에 Pull 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreatePullJobForCluster(imageName string) error {
	return r.CreateJobForCluster(imageName, "pull")
}

//CreateJobForCluster 원하는 노드에 특정 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreateJobForCluster(imageName string, cmdType string) error {
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
		createErr := r.CreateJob(result.Name, imageName, cmdType)
		if createErr != nil {
			return createErr
		}
		//return createErr
		//})
		//if retryErr != nil {
		//	return retryErr
		//}
		//omcplog.V(3).Info("Updated ..." + node.Name)
	}

	return nil
}

//CreatePushJob 원하는 노드에 Push 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreatePushJob(nodeName string, imageName string) error {
	return r.CreateJob(nodeName, imageName, "push")
}

//CreatePullJob 원하는 노드에 Push 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreatePullJob(nodeName string, imageName string) error {
	return r.CreateJob(nodeName, imageName, "pull")
}

//CreatePullJob 원하는 노드에 Push 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreateDeleteJob(imageName string) error {
	return r.CreateDelete(imageName)
}
func (r *RegistryManager) CreateDelete(imageName string) error {
	omcplog.V(4).Info("Create Delete Command")
	//nfs 해당 이미지 삭제
	return nil
}

//CreateJob 원하는 노드에 특정 명령을 내리는 job을 생성하는 기능.
func (r *RegistryManager) CreateJob(nodeName string, imageName string, cmdType string) error {
	appName := r.getAppName(nodeName)
	labelName := r.getLabelName(nodeName)

	jobClient := r.clientset.BatchV1().Jobs(utils.ProjectNamespace)
	//이미 존재할 때의 처리 방법.
	job, _ := jobClient.Get(context.TODO(), appName, metav1.GetOptions{})
	if job.ObjectMeta.Name != "" {
		omcplog.V(3).Info("job exist : " + job.ObjectMeta.Name + "\n")
		err := jobClient.Delete(context.TODO(), appName, metav1.DeleteOptions{})
		if err != nil {
			omcplog.V(3).Info("exist job delete error!")
			return nil
		} else {
			omcplog.V(3).Info("job : " + job.ObjectMeta.Name + "delete complete!")
		}
	}
	cmd := ""

	imagedir := "/nfs"
	switch cmdType {
	case "delete":
		// docker rm images
		cmd += "rm -r " + imagedir + "/" + imageName + ".tar" + ";"
		cmd += "docker rmi -f " + imageName + ";"

		cmd = utils.SetImageCommand(cmd)
	case "push":
		//docker save images
		if strings.Contains(imageName, "/") {
			dirlist := strings.SplitAfter(imageName, "/")
			dirName := strings.Replace(imageName, dirlist[len(dirlist)-1], "", -1)
			cmd += "mkdir -p " + dirName + ";"
			omcplog.V(4).Info("mkdir -p " + dirName)
		}
		cmd += "docker save -o " + imageName + ".tar" + " " + imageName + ";"
		omcplog.V(4).Info("docker save -o " + imageName + ".tar" + " " + imageName + ";")
		cmd = utils.SetImageCommand(cmd)
	case "pull":
		cmd = "docker load -i " + imageName + ".tar" + ";"
		cmd = utils.SetImageCommand(cmd)
	default:
	}

	//1. 해당 push 할 컨테이너 명을 찾는다.
	//2. commit 명령을 내린다.
	//commitCommand := "docker commit -a openmcp -m 'make " + imageName + "' 598bc3e5efe7 " + imageFullName

	job = r.getJobAPI(appName, labelName, cmd)

	// Create Deployment
	omcplog.V(3).Info("Creating " + cmdType + " Job...")
	_, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	omcplog.V(3).Info("Created "+cmdType+" Job %q.\n", job.GetObjectMeta().GetName())

	//r.JobRunCheck(batchv1.JobComplete, afterFunc)

	time.Sleep(time.Second * 5)
	return nil
}
func int32Ptr(i int32) *int32 { return &i }
