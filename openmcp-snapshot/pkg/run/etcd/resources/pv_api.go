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

// PersistentVolume resource
type PersistentVolume struct {
	apiCaller typedv1.PersistentVolumeInterface
}

// convertResourceObj : json String 을 obj 로 변환
func (pv PersistentVolume) convertResourceObj(resourceInfoJSON string) (*apiv1.PersistentVolume, error) {

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

// CreateResourceForJSON : deploy 생성
func (pv PersistentVolume) CreateResourceForJSON(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := pv.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	pv.apiCaller = clientset.CoreV1().PersistentVolumes()

	for key, val := range resourceInfo.GetObjectMeta().GetLabels() {
		fmt.Println("===Label===")
		fmt.Println(key, val)
		resourceInfo.ObjectMeta.Labels[key] = val + SnapshotTailName
	}

	resourceInfo.ObjectMeta.ResourceVersion = ""
	resourceInfo.ObjectMeta.Name = resourceInfo.ObjectMeta.Name + SnapshotTailName
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
	fmt.Println("Creating pv...")
	result, apiCallErr := pv.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created pv %q.\n", result.GetObjectMeta().GetName())

	return true, nil
}

// GetJSON : return json string
func (pv PersistentVolume) GetJSON(clientset *kubernetes.Clientset, resourceName string, resourceNamespace string) (string, error) {
	pv.apiCaller = clientset.CoreV1().PersistentVolumes()

	result, getErr := pv.apiCaller.Get(resourceName, metav1.GetOptions{})
	if getErr != nil {
		return "", getErr
	}

	return util.Obj2JsonString(result)
}
