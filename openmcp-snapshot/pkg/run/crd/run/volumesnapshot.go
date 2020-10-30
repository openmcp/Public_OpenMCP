package run

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	snapshotv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"
)

// DynamicVolumeSnapshot resource
type DynamicVolumeSnapshot struct {
	apiCaller dynamic.ResourceInterface
}

// CreateResource : VolumeSnapshot CRD 를 생성하는 함수.
func (dynamicResource DynamicVolumeSnapshot) CreateResource(clientset dynamic.Interface, namespace string, volumeDataSource snapshotv1alpha1.VolumeDataSource) (bool, error) {
	//namespace := apiv1.NamespaceDefault
	gvk := schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1beta1", Resource: "VolumeSnapshot"}
	dynamicResource.apiCaller = clientset.Resource(gvk).Namespace(namespace)
	name := util.GetVolumeSnapshotName(volumeDataSource)

	//dynamic.apiCaller = clientset.AppsV1().Deployments(nameSpace)

	resourceInfo := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "snapshot.storage.k8s.io/v1beta1",
			"kind":       "VolumeSnapshot",
			"metadata": map[string]interface{}{
				//"name": "snapshot-source-pvc",
				"name": name,
			},
			"spec": map[string]interface{}{
				"volumeSnapshotClassName": volumeDataSource.VolumeSnapshotClassName,
				"source": map[string]interface{}{
					"persistentVolumeClaimName": volumeDataSource.VolumeSnapshotSourceName,
				},
			},
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

// IsRunningVolumeSnapshot : VolumeSnapshot 의 구동 완료 준비가 되었는지 체크하는 함수.
func (dynamicResource DynamicVolumeSnapshot) IsRunningVolumeSnapshot(clientset dynamic.Interface, namespace string, volumeDataSource snapshotv1alpha1.VolumeDataSource) (bool, error) {
	gvk := schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1beta1", Resource: "VolumeSnapshot"}

	dynamicResource.apiCaller = clientset.Resource(gvk).Namespace(namespace)
	name := util.GetVolumeSnapshotName(volumeDataSource)

	result, apiCallErr := dynamicResource.apiCaller.Get(name, metav1.GetOptions{})
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created CRD 2 %q.\n", result)

	/* result 내의 내용.
	apiVersion: snapshot.storage.k8s.io/v1beta1
	kind: VolumeSnapshot
	metadata:
		creationTimestamp: "2020-04-14T02:13:13Z"
		finalizers:
		- snapshot.storage.kubernetes.io/volumesnapshot-as-source-protection
		- snapshot.storage.kubernetes.io/volumesnapshot-bound-protection
		generation: 1
		name: snapshot-source-pvc
		namespace: default
		resourceVersion: "1293962"
		selfLink: /apis/snapshot.storage.k8s.io/v1beta1/namespaces/default/volumesnapshots/snapshot-source-pvc
		uid: 7eced146-7df5-11ea-ab69-42010a80002a
	spec:
		source:
			persistentVolumeClaimName: source-pvc
		volumeSnapshotClassName: default-snapshot-class
	status:
		boundVolumeSnapshotContentName: snapcontent-7eced146-7df5-11ea-ab69-42010a80002a
		readyToUse: true
	*/

	// openmcp snapshot 이 아닌 생성된 VolumeSnapshot 에 대해서 Snapshot ID 를 추출해야 한다.
	//volumeSnapshotContentName := unstructured.NestedString(result.Object, "status", "boundVolumeSnapshotContentName")
	isReady, found, err := unstructured.NestedBool(result.Object, "status", "readyToUse")
	if err != nil || !found {
		fmt.Printf("VolumeSnapshot not found %s: error=%s", result.GetName(), err)
		return false, err
	}
	isReady = true
	return isReady, nil
}

/*
// GetResourceSnapshotID : VolumeSnapshot CRD 에서 데이터를 가져와서 Snapshot ID 를 추출하는 함수.
func (dynamicResource DynamicVolumeSnapshot) GetResourceSnapshotID(clientset dynamic.Interface, namespace string, volumeDataSource nanumv1alpha1.VolumeDataSource) (string, error) {
	gvk := schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1beta1", Resource: "VolumeSnapshot"}

	dynamicResource.apiCaller = clientset.Resource(gvk).Namespace(namespace)
	name := utils.GetVolumeSnapshotName(volumeDataSource)

	result, apiCallErr := dynamicResource.apiCaller.Get(name, metav1.GetOptions{})
	if apiCallErr != nil {
		return "", apiCallErr
	}
	fmt.Printf("Created CRD 2 %q.\n", result)

	// openmcp snapshot 이 아닌 생성된 VolumeSnapshot 에 대해서 Snapshot ID 를 추출해야 한다.
	//volumeSnapshotContentName := unstructured.NestedString(result.Object, "status", "boundVolumeSnapshotContentName")
	unstructured.NestedString(result.Object, "status", "boundVolumeSnapshotContentName")
	//snapshot-source-pvc

	return "", nil
}
*/
