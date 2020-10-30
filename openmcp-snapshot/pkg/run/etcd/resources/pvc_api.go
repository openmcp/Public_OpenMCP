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

// PersistentVolumeClaim resource
type PersistentVolumeClaim struct {
	apiCaller typedv1.PersistentVolumeClaimInterface
}

// convertResourceObj : json String 을 obj 로 변환
func (pvc PersistentVolumeClaim) convertResourceObj(resourceInfoJSON string) (*apiv1.PersistentVolumeClaim, error) {

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

// CreateResourceForJSON : deploy 생성
func (pvc PersistentVolumeClaim) CreateResourceForJSON(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := pvc.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	nameSpace := apiv1.NamespaceDefault
	if resourceInfo.GetObjectMeta().GetNamespace() != "" && resourceInfo.GetObjectMeta().GetNamespace() != apiv1.NamespaceDefault {
		nameSpace = resourceInfo.GetObjectMeta().GetNamespace()
	}

	pvc.apiCaller = clientset.CoreV1().PersistentVolumeClaims(nameSpace)

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

	resourceInfo.ObjectMeta.ResourceVersion = ""
	resourceInfo.ObjectMeta.Name = resourceInfo.ObjectMeta.Name + SnapshotTailName
	resourceInfo.Spec.VolumeName = resourceInfo.Spec.VolumeName + SnapshotTailName

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
	fmt.Println("Creating pvc...")
	result, apiCallErr := pvc.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created pvc %q.\n", result.GetObjectMeta().GetName())

	return true, nil
}

// GetJSON : return json string
func (pvc PersistentVolumeClaim) GetJSON(clientset *kubernetes.Clientset, resourceName string, resourceNamespace string) (string, error) {
	pvc.apiCaller = clientset.CoreV1().PersistentVolumeClaims(resourceNamespace)
	fmt.Printf("Listing Resource in namespace %q:\n", resourceNamespace)

	result, apiCallErr := pvc.apiCaller.Get(resourceName, metav1.GetOptions{})
	if apiCallErr != nil {
		return "", apiCallErr
	}

	return util.Obj2JsonString(result)
}
