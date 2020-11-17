package resources

import (
	"encoding/json"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
)

// json2ServiceInit : json String 을 obj 로 변환
func json2ServiceInit(resourceInfoJSON string) (*apiv1.Service, error) {

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

// JSON2Service : service 조작
func JSON2Service(resourceInfoJSON string) (*apiv1.Service, error) {
	resourceInfo, convertErr := json2ServiceInit(resourceInfoJSON)
	if convertErr != nil {
		return nil, convertErr
	}

	//	for key, val := range resourceInfo.GetObjectMeta().GetLabels() {
	//		fmt.Println("===Label===")
	//		fmt.Println(key, val)
	//		resourceInfo.ObjectMeta.Labels[key] = val + SnapshotTailName
	//	}
	//	resourceInfo.ObjectMeta.Name = resourceInfo.ObjectMeta.Name + SnapshotTailName
	resourceInfo.ObjectMeta.ResourceVersion = ""

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
	return resourceInfo, nil
}
