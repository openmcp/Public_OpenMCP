package resources

import (
	"context"
	"encoding/json"
	"openmcp/openmcp/omcplog"

	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"

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
func JSON2Service(resourceInfoJSON string, client genericclient.Client) (*apiv1.Service, error) {
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

	//update 전에 알아낼 revision 번호를 가져오기 위함으로 결국 resource Version 알아내기
	currentObj := &apiv1.Service{}
	err := client.Get(context.TODO(), currentObj, resourceInfo.Namespace, resourceInfo.Name)
	if err != nil {
		omcplog.V(2).Info("Pre Get Resource for JSON error", err)
		return nil, err
	}
	//resourceInfo.ObjectMeta.ResourceVersion = currentObj.ObjectMeta.ResourceVersion
	resourceInfo.ObjectMeta.UID = currentObj.ObjectMeta.UID
	resourceInfo.ObjectMeta.ResourceVersion = currentObj.ObjectMeta.ResourceVersion
	resourceInfo.Spec.ClusterIP = currentObj.Spec.ClusterIP
	//resourceInfo.Spec = currentObj.DeepCopy().Spec
	//resourceInfo.ObjectMeta.ResourceVersion = ""

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
	omcplog.V(2).Info("Creating service Obj...")
	return resourceInfo, nil
}
