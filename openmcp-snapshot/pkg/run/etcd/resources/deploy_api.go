package resources

import (
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"openmcp/openmcp/openmcp-snapshot/pkg/run/etcd/util"
)

// Deployment resource
type Deployment struct {
	apiCaller typedv1.DeploymentInterface
}

// convertResourceObj : json String 을 obj 로 변환
func (deploy Deployment) convertResourceObj(resourceInfoJSON string) (*appsv1.Deployment, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var deployment *appsv1.Deployment
	jsonEerr := json.Unmarshal(jsonBytes, &deployment)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return deployment, nil
}

// CreateResourceForJSON : deploy 생성
func (deploy Deployment) CreateResourceForJSON(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := deploy.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	nameSpace := apiv1.NamespaceDefault
	if resourceInfo.GetObjectMeta().GetNamespace() != "" && resourceInfo.GetObjectMeta().GetNamespace() != apiv1.NamespaceDefault {
		nameSpace = resourceInfo.GetObjectMeta().GetNamespace()
	}

	deploy.apiCaller = clientset.AppsV1().Deployments(nameSpace)

	for key, val := range resourceInfo.GetObjectMeta().GetLabels() {
		fmt.Println("===Label===")
		fmt.Println(key, val)
		resourceInfo.ObjectMeta.Labels[key] = val + SnapshotTailName
	}
	for key, val := range resourceInfo.Spec.Selector.MatchLabels {
		fmt.Println("===spec_MatchLabel===")
		fmt.Println(key, val)
		resourceInfo.Spec.Selector.MatchLabels[key] = val + SnapshotTailName
	}
	for key, val := range resourceInfo.Spec.Template.Labels {
		fmt.Println("===spec_template_Label===")
		fmt.Println(key, val)
		resourceInfo.Spec.Template.Labels[key] = val + SnapshotTailName
	}
	for i, val := range resourceInfo.Spec.Template.Spec.Volumes {
		fmt.Println("===spec_template_volumes===")
		fmt.Println(i, val)
		// TODO 여기는 pvc랑 매핑되는 부분인데, 외부 모듈 연결시 추가 작업이 필요함.
		resourceInfo.Spec.Template.Spec.Volumes[i].VolumeSource.PersistentVolumeClaim.ClaimName = val.VolumeSource.PersistentVolumeClaim.ClaimName + SnapshotTailName
	}

	resourceInfo.ObjectMeta.ResourceVersion = ""
	resourceInfo.ObjectMeta.Name = resourceInfo.ObjectMeta.Name + SnapshotTailName
	/*
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				//Name: "demo",
				Name:      resourceInfo.GetObjectMeta().GetName() + SnapshotTailName,
				Namespace: resourceInfo.ObjectMeta.Namespace,
				Labels:    resourceInfo.ObjectMeta.Labels,
			},

			Spec: resourceInfo.Spec,

			Spec: appsv1.DeploymentSpec{
				Replicas: resourceInfo.Spec.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: resourceInfo.Spec.Selector.MatchLabels,
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: resourceInfo.Spec.Template.ObjectMeta.Labels,
					},
					Spec: apiv1.PodSpec{
						Containers: resourceInfo.Spec.Template.Spec.Containers,
					},
				},
			},
		}
	*/

	// Create Deployment
	fmt.Println("Creating deployment...")

	result, apiCallErr := deploy.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return true, nil
}

// GetJSON : return json string
func (deploy Deployment) GetJSON(clientset *kubernetes.Clientset, resourceName string, resourceNamespace string) (string, error) {
	deploy.apiCaller = clientset.AppsV1().Deployments(resourceNamespace)

	fmt.Printf("Listing Resource in namespace %q:\n", resourceNamespace)

	result, apiCallErr := deploy.apiCaller.Get(resourceName, metav1.GetOptions{})
	if apiCallErr != nil {
		return "", apiCallErr
	}

	return util.Obj2JsonString(result)
}
