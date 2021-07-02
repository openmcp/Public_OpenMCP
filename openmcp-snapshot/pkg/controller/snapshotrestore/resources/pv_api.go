package resources

import (
	"encoding/json"
	"openmcp/openmcp/omcplog"

	apiv1 "k8s.io/api/core/v1"
)

// json2PvInit : json String 을 obj 로 변환
func json2PvInit(resourceInfoJSON string) (*apiv1.PersistentVolume, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var persistentVolume *apiv1.PersistentVolume
	jsonEerr := json.Unmarshal(jsonBytes, &persistentVolume)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return persistentVolume, nil
}

// JSON2Pv : pv 생성
func JSON2Pv(resourceInfoJSON string) (*apiv1.PersistentVolume, error) {
	resourceInfo, convertErr := json2PvInit(resourceInfoJSON)
	if convertErr != nil {
		return nil, convertErr
	}

	//	for key, val := range resourceInfo.GetObjectMeta().GetLabels() {
	//		fmt.Println("===Label===")
	//		fmt.Println(key, val)
	//		resourceInfo.ObjectMeta.Labels[key] = val + SnapshotTailName
	//	}
	//resourceInfo.ObjectMeta.Name = resourceInfo.ObjectMeta.Name + SnapshotTailName
	resourceInfo.ObjectMeta.ResourceVersion = ""
	resourceInfo.Spec.ClaimRef = nil
	/*
		pvc := &apiv1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				//Name: "demo",
				Name:   resourceInfo.ObjectMeta.Name + SnapshotTailName,
				Labels: resourceInfo.ObjectMeta.Labels,
			},

			Spec: resourceInfo.Spec,
			//Spec: appsv1.PersistentVolumeSpec{
			//	AccessModes: resourceInfo.Spec.AccessModes,
			//	Resources: resourceInfo.Spec.Resources,
			//	Selector: resourceInfo.Spec.Selector,
			//},
		}
	*/

	// Create PersistentVolume
	omcplog.V(2).Info("Creating pv Obj...")
	return resourceInfo, nil
}
