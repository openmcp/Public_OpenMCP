package resources

import (
	"encoding/json"
	"openmcp/openmcp/omcplog"

	apiv1 "k8s.io/api/core/v1"
)

// json2PvcInit : json String 을 obj 로 변환
func json2PvcInit(resourceInfoJSON string) (*apiv1.PersistentVolumeClaim, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var persistentVolumeClaim *apiv1.PersistentVolumeClaim
	jsonEerr := json.Unmarshal(jsonBytes, &persistentVolumeClaim)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return persistentVolumeClaim, nil
}

// JSON2Pvc : pvc 조작
func JSON2Pvc(resourceInfoJSON string) (*apiv1.PersistentVolumeClaim, error) {
	resourceInfo, convertErr := json2PvcInit(resourceInfoJSON)
	if convertErr != nil {
		return nil, convertErr
	}

	//	for key, val := range resourceInfo.GetObjectMeta().GetLabels() {
	//		fmt.Println("===Label===")
	//		fmt.Println(key, val)
	//		resourceInfo.ObjectMeta.Labels[key] = val + SnapshotTailName
	//	}
	//	for key, val := range resourceInfo.Spec.Selector.MatchLabels {
	//		fmt.Println("===spec_MatchLabel===")
	//		fmt.Println(key, val)
	//		resourceInfo.Spec.Selector.MatchLabels[key] = val + SnapshotTailName
	//	}
	//	resourceInfo.ObjectMeta.Name = resourceInfo.ObjectMeta.Name + SnapshotTailName
	//	resourceInfo.Spec.VolumeName = resourceInfo.Spec.VolumeName + SnapshotTailName

	resourceInfo.ObjectMeta.ResourceVersion = ""

	/*
		pvc := &apiv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				//Name: "demo",
				Name:      resourceInfo.ObjectMeta.Name + SnapshotTailName,
				Namespace: resourceInfo.ObjectMeta.Namespace,
				Labels:    resourceInfo.ObjectMeta.Labels,
			},

			Spec: resourceInfo.Spec,
			//Spec: appsv1.PersistentVolumeClaimSpec{
			//	AccessModes: resourceInfo.Spec.AccessModes,
			//	Resources: resourceInfo.Spec.Resources,
			//	Selector: resourceInfo.Spec.Selector,
			//},
		}
	*/
	// Create PersistentVolumeClaim
	omcplog.V(2).Info("Creating pvc Obj...")
	return resourceInfo, nil
}
