package dist

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/oklog/ulid"

	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-cache/pkg/utils"
	"openmcp/openmcp/util/clusterManager"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

//RegistryManager 레지스트리 배포 매니저
type RegistryManager struct {
	clientset   *kubernetes.Clientset
	clusterName string

	watchType batchv1.JobConditionType
	afterFunc func(newJob interface{}, oldJob interface{})
	informer  cache.SharedIndexInformer
	stopper   chan struct{}
}

func (r *RegistryManager) Init(clusterName string, myCluster *clusterManager.ClusterManager) error {
	var clientset *kubernetes.Clientset

	clientset = myCluster.Cluster_kubeClients[clusterName]
	r.clientset = clientset
	r.clusterName = clusterName

	return nil
}

func (r *RegistryManager) getLabelName(nodeName string) string {
	//labelName := utils.ProjectDomain + "_" + r.clusterName + "_" + nodeName
	labelName := r.clusterName + "_" + nodeName
	return labelName
}

func (r *RegistryManager) getAppName(nodeName string) string {
	uid := r.genUlid()
	//appName := "repository-agent-" + r.clusterName + "-" + nodeName
	appName := "ic-" + r.clusterName + "-" + nodeName
	appName = appName + "-" + uid

	if len(appName) >= 64 {
		appName = appName[0:63]
		omcplog.V(4).Info(" appName size error! ")
	}
	return appName
}

func (r *RegistryManager) genUlid() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	omcplog.V(4).Info("ulid: " + id.String())
	return strings.ToLower(id.String())
}

//canCreatedPod pod 생성이 허용된 노드인지 확인한다.  (false 시 이미지 관리 대상 제외)
func (r *RegistryManager) canCreatedPod(nodeInfo *apiv1.Node) bool {
	//Node 정보에서 Taints 내에 node-role.kubernetes.io/master:NoSchedule 로 지정된 노드들은 pod 생성시 오류가 발생한다.
	//이 친구들은 애초에 pod 생성이 허락된 친구들이 아니기 때문에 image 관리 대상에서 제외한다.
	// 마스터도 pod 생성이 가능하려면 kubectl taint nodes --all node-role.kubernetes.io/master- 명령어 기입.
	for _, Taint := range nodeInfo.Spec.Taints {
		if Taint.Key == "node-role.kubernetes.io/master" && Taint.Value == "NoSchedule" {
			return false
		}
	}
	return true
}

//checkNamespace namespace check 하여 없으면 만들어준다.
func (r *RegistryManager) checkNamespace() (bool, error) {

	if r.clientset == nil {
		return false, errors.New("no cluster is specified")
	}

	//*v1.namespaceList
	//namespace := &v1.Namespace{}
	//getErr := kubeerrors.StatusError{}
	namespace, _ := r.clientset.CoreV1().Namespaces().Get(context.TODO(), utils.ProjectNamespace, metav1.GetOptions{})
	//if getErr != nil {
	//	return false, getErr
	//}

	// 없다. 오류인경우도 여기로 들어감.
	if namespace.ObjectMeta.Name == "" {
		//name space 생성

		nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: utils.ProjectNamespace}}
		_, err := r.clientset.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
		if err != nil {
			return false, err
		}
	}
	omcplog.V(4).Info("Create namespace: ", utils.ProjectNamespace)
	return true, nil
}

//checkNamespace namespace check 하여 없으면 만들어준다.
func (r *RegistryManager) getJobAPI(appName string, labelName string, cmd string) *batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: "openmcp",
		},
		Spec: batchv1.JobSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					//"app": "openmcp",
				},
			},
			//ManualSelector: false,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "openmcp-cache",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "docker-cmds",
							Image: "docker:19.03.8",
							Command: []string{
								"/bin/sh",
							},
							Args: []string{
								"-c", cmd,
								//"-c", "echo 'test'; docker commit -a gen -m commitTest 598bc3e5efe7 10.0.0.224:4999/test:2.0; do echo running; sleep 1000;done"
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/var/run",
									Name:      "docker-sock",
								},
								{
									MountPath: "/etc/docker/",
									Name:      "docker-cert",
								},
								{
									MountPath: "/nfs",
									Name:      "nfs-volume",
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
					NodeSelector: map[string]string{
						labelName: "true",
					},
					Volumes: []apiv1.Volume{
						{
							Name: "docker-sock",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/var/run",
								},
							},
						},
						{
							Name: "docker-cert",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/etc/docker/",
								},
							},
						},
						{
							Name: "nfs-volume",
							VolumeSource: apiv1.VolumeSource{
								NFS: &apiv1.NFSVolumeSource{
									Server: utils.ImageCacheNfs,
									Path:   "/home/nfs/images",
								},
							},
						},
					},
				},
			},
		},
	}
	return job
}
