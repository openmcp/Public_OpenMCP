package resources

import (
	"encoding/json"
	"fmt"

	"openmcp/openmcp/openmcp-snapshot/pkg/run/etcd/util"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// Service resource
type Service struct {
	apiCaller typedv1.ServiceInterface
}

// convertResourceObj : json String 을 obj 로 변환
func (svc Service) convertResourceObj(resourceInfoJSON string) (*apiv1.Service, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var service *apiv1.Service
	jsonEerr := json.Unmarshal(jsonBytes, &service)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return service, nil
}

// CreateResourceForJSON : deploy 생성
func (svc Service) CreateResourceForJSON(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := svc.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	nameSpace := apiv1.NamespaceDefault
	if resourceInfo.GetObjectMeta().GetNamespace() != "" && resourceInfo.GetObjectMeta().GetNamespace() != apiv1.NamespaceDefault {
		nameSpace = resourceInfo.GetObjectMeta().GetNamespace()
	}
	svc.apiCaller = clientset.CoreV1().Services(nameSpace)

	for key, val := range resourceInfo.GetObjectMeta().GetLabels() {
		fmt.Println("===Label===")
		fmt.Println(key, val)
		resourceInfo.ObjectMeta.Labels[key] = val + SnapshotTailName
	}

	resourceInfo.ObjectMeta.ResourceVersion = ""
	resourceInfo.ObjectMeta.Name = resourceInfo.ObjectMeta.Name + SnapshotTailName

	/*
		service := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceInfo.ObjectMeta.Name + SnapshotTailName,
				Namespace: resourceInfo.ObjectMeta.Namespace,
				Labels:    resourceInfo.ObjectMeta.Labels,
			},

			Spec: apiv1.ServiceSpec{
				Type: resourceInfo.Spec.Type,
				Ports: resourceInfo.Spec.Ports,
				Selector: resourceInfo.Spec.Selector,
				ClusterIP: resourceInfo.Spec.ClusterIP,
			},
		}
	*/

	// Create Service
	fmt.Println("Creating service...")
	result, apiCallErr := svc.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())

	return true, nil
}

// GetJSON : return json string
func (svc Service) GetJSON(clientset *kubernetes.Clientset, resourceName string, resourceNamespace string) (string, error) {
	svc.apiCaller = clientset.CoreV1().Services(resourceNamespace)
	fmt.Printf("Listing Resource in namespace %q:\n", resourceNamespace)

	result, apiCallErr := svc.apiCaller.Get(resourceName, metav1.GetOptions{})
	if apiCallErr != nil {
		return "", apiCallErr
	}

	return util.Obj2JsonString(result)
}
