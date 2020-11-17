package snapshot

// import (
// 	"context"
// 	"fmt"
// 	"strconv"
// 	"time"

// 	snapshotv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"

// 	nanumetcd "openmcp/openmcp/openmcp-snapshot/pkg/run/etcd"

// 	v1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// /* OLD 코드 */
// // Run : 실제 로직단
// func (r *Reconcile) Run(instance *snapshotv1alpha1.Snapshot) (bool, error) {

// 	// Resource 별로 for 문 진행
// 	for idx, resource := range instance.Spec.SnapshotSources {
// 		fmt.Println("\n[" + strconv.Itoa(idx) + "] : Resource")
// 		switch resource.ResourceType {
// 		//case "PersistentVolume", "persistentvolume", "pv":
// 		// 아무것도 하지 않는다.
// 		//	volumeSnapshot()
// 		case "PersistentVolumeClaim", "persistentvolumeclaim", "pvc":
// 			isMakeVSnapshot, vSnapshotErr := volumeSnapshot(resource)
// 			if !isMakeVSnapshot || vSnapshotErr != nil {
// 				return false, vSnapshotErr
// 			}
// 			fallthrough // 이어서 default 실행
// 		default:
// 			isMakeESnapshot, eSnapshotErr := r.etcdSnapshot(resource)
// 			if !isMakeESnapshot || eSnapshotErr != nil {
// 				return false, eSnapshotErr
// 			}
// 			setHistoryEtcdSnapshot(resource)
// 		}
// 	}
// 	r.updateInfo(instance)

// 	return true, nil
// }

// // updateInfo 는 성공한 결과를 status 에 반영하는 함수입니다.
// func (r *Reconcile) updateInfo(instance *snapshotv1alpha1.Snapshot) (bool, error) {
// 	instance.Status.Status = true

// 	snapshotErr := r.client.Status().Update(context.TODO(), instance)

// 	if snapshotErr != nil {
// 		return false, snapshotErr
// 	}

// 	return true, nil
// }

// // volumeSnapshot :
// func (r *Reconcile) volumeSnapshot(resource snapshotv1alpha1.SnapshotSource) (bool, error) {

// 	// ------------------------------
// 	//ctrl := crdResource.DynamicVolumeSnapshot{}

// 	// Volume Snapshot
// 	clientset := nanumetcd.InitKube()
// 	// volumeSnapshotSourceKind 는 PersistentVolumeClaim 고정. volumeSnapshotSourceName 는 pvc 이름

// 	isSuccessSnapshot, snapshotErr := r.CreateResource(clientset, resource.ResourceNamespace, resource.VolumeDataSource)

// 	if snapshotErr != nil {
// 		return false, snapshotErr
// 	}
// 	if !isSuccessSnapshot {
// 		return false, fmt.Errorf(resource.VolumeDataSource.VolumeSnapshotSourceName + "에 대한 VolumeSnapshot 생성에 실패하였습니다.")
// 	}

// 	// Snapshot 상태 계속 기다리면서 가져오기.
// 	var getError error
// 	var isRunning bool
// 	maxCount := 30
// 	for i := 0; i < maxCount; i++ {
// 		isRunning, getError = r.IsRunningVolumeSnapshot(resource)
// 		if isRunning {
// 			break
// 		}
// 		time.Sleep(1 * time.Second)
// 	}
// 	if getError != nil {
// 		return false, getError
// 	}
// 	if !isRunning {
// 		return false, fmt.Errorf(resource.VolumeDataSource.VolumeSnapshotSourceName + "가 구동 중이 아닙니다.")
// 	}
// 	return isRunning, nil
// }

// // IsRunningVolumeSnapshot 볼륨 구동중인지 확인 - pvc 의 Status.Phase 값을 비교하여 bound 일때 true 반환
// func (r *Reconcile) IsRunningVolumeSnapshot(resource snapshotv1alpha1.SnapshotSource) (bool, error) {
// 	//clientset := crd.DynamicInitKube()
// 	clientset := nanumetcd.InitKube()
// 	apiCaller := clientset.CoreV1().PersistentVolumeClaims(resource.ResourceNamespace)
// 	fmt.Printf("Listing Resource in namespace %q:\n", resource.ResourceNamespace)

// 	result, apiCallErr := apiCaller.Get(resource.ResourceName, metav1.GetOptions{})
// 	if apiCallErr != nil {
// 		return false, apiCallErr
// 	}
// 	if result.Status.Phase == v1.ClaimBound {
// 		return true, nil
// 	} else if result.Status.Phase == v1.ClaimPending {
// 		return false, fmt.Errorf(resource.VolumeDataSource.VolumeSnapshotSourceName + "가 구동 진행 중입니다.")
// 	} else {
// 		return false, fmt.Errorf(resource.VolumeDataSource.VolumeSnapshotSourceName + "가 구동되지 않았습니다.")
// 	}

// }

// // setHistoryVolumeSnapshot 는 Volume Snapshot 한 결과를 etcd 데이터에 넣는 함수입니다.
// func (r *Reconcile) setHistoryVolumeSnapshot(resource snapshotv1alpha1.SnapshotSource) (bool, error) {
// 	return true, nil
// }

// func (r *Reconcile) etcdSnapshot(resource snapshotv1alpha1.SnapshotSource) (bool, error) {
// 	//TODO 가져온 ID ETCD 에 저장하는 부분.
// 	// Etcd Snapshot
// 	isSuccessSnapshot, snapshotErr := nanumetcd.CreateSnapshot(resource)
// 	if snapshotErr != nil {
// 		return false, snapshotErr
// 	}
// 	if !isSuccessSnapshot {
// 		return false, fmt.Errorf("Etcd snapshot failed")
// 	}

// 	return isSuccessSnapshot, snapshotErr
// }

// func setHistoryEtcdSnapshot(resource snapshotv1alpha1.SnapshotSource) (bool, error) {
// 	return true, nil
// }
