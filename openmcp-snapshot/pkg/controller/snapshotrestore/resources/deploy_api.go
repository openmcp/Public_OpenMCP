package resources

import (
	"encoding/json"
	"openmcp/openmcp/omcplog"

	appsv1 "k8s.io/api/apps/v1"
)

// json2DeployInit : json String 을 obj 로 변환
func json2DeployInit(resourceInfoJSON string) (*appsv1.Deployment, error) {

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

// JSON2Deploy : deploy 조작
func JSON2Deploy(resourceInfoJSON string) (*appsv1.Deployment, error) {
	resourceInfo, convertErr := json2DeployInit(resourceInfoJSON)
	if convertErr != nil {
		return nil, convertErr
	}

	// tail 은 붙이지 않는다.
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
	//	for key, val := range resourceInfo.Spec.Template.Labels {
	//		fmt.Println("===spec_template_Label===")
	//		fmt.Println(key, val)
	//		resourceInfo.Spec.Template.Labels[key] = val + SnapshotTailName
	//	}
	//	for i, val := range resourceInfo.Spec.Template.Spec.Volumes {
	//		fmt.Println("===spec_template_volumes===")
	//		fmt.Println(i, val)
	//		// TODO 여기는 pvc랑 매핑되는 부분인데, 외부 모듈 연결시 추가 작업이 필요함.
	//		resourceInfo.Spec.Template.Spec.Volumes[i].VolumeSource.PersistentVolumeClaim.ClaimName = val.VolumeSource.PersistentVolumeClaim.ClaimName + SnapshotTailName
	//	}
	//resourceInfo.ObjectMeta.Name = resourceInfo.ObjectMeta.Name + SnapshotTailName

	resourceInfo.ObjectMeta.ResourceVersion = ""
	resourceInfo.ObjectMeta.UID = ""
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
	omcplog.V(2).Info("Creating deployment Obj...")
	return resourceInfo, nil
}
