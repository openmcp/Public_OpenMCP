package init

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// DynamicVolumeSnapshotClass resource
type DynamicVolumeSnapshotClass struct {
	apiCaller dynamic.ResourceInterface
}

// CreateResource : Resource 생성
func (dynamicResource DynamicVolumeSnapshotClass) CreateResource(clientset dynamic.Interface, name string) (bool, error) {
	namespace := apiv1.NamespaceDefault
	gvk := schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1beta1", Resource: "VolumeSnapshotClass"}

	dynamicResource.apiCaller = clientset.Resource(gvk).Namespace(namespace)

	//dynamic.apiCaller = clientset.AppsV1().Deployments(nameSpace)

	resourceInfo := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "snapshot.storage.k8s.io/v1beta1",
			"kind":       "VolumeSnapshotClass",
			"metadata": map[string]interface{}{
				//"name": "default-snapshot-class",
				"name": name,
			},
			"driver":       "pd.csi.storage.gke.io",
			"deletePolicy": "Delete",
		},
	}

	// Create Deployment
	fmt.Println("Creating CRD 2...")

	//resourceInfo.SetResourceVersion("")
	//oldName := resourceInfo.GetName()
	//newName := oldName + "-snapshot"
	//resourceInfo.SetName(newName)
	result, apiCallErr := dynamicResource.apiCaller.Create(resourceInfo, metav1.CreateOptions{})
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created CRD 2 %q.\n", result)

	return true, nil

}
