package snapshotrestore

// // Run : 실제 로직단
// func (r *ReconcileSnapshotRestore) Run(instance *snapshotv1alpha1.SnapshotRestore) (bool, error) {

// 	// Resource 별로 for 문 진행
// 	for idx, resource := range instance.Spec.SnapshotRestoreSource {
// 		fmt.Println("\n[" + strconv.Itoa(idx) + "] : Resource")
// 		switch resource.ResourceType {
// 		//case "PersistentVolume", "persistentvolume", "pv":
// 		// 아무것도 하지 않는다.
// 		//	volumeSnapshot()
// 		case "PersistentVolumeClaim", "persistentvolumeclaim", "pvc":
// 			isMakeVSnapshot, vSnapshotErr := volumeSnapshotRestore(resource)
// 			if !isMakeVSnapshot || vSnapshotErr != nil {
// 				return false, vSnapshotErr
// 			}
// 			fallthrough // 이어서 default 실행
// 		default:
// 			isMakeESnapshot, eSnapshotErr := etcdSnapshotRestore(resource)
// 			if !isMakeESnapshot || eSnapshotErr != nil {
// 				return false, eSnapshotErr
// 			}
// 			//setHistoryEtcdSnapshot(resource)
// 		}
// 	}
// 	//updateInfo(instance)

// 	return true, nil
// }

// // updateInfo 는 성공한 결과를 status 에 반영하는 함수입니다.
// func updateInfo(instance *snapshotv1alpha1.Snapshot) (bool, error) {
// 	ctrl := crdResource.DynamicVolumeSnapshot{}
// 	clientset := crd.DynamicInitKube()
// 	isSuccessSnapshot, snapshotErr := ctrl.UpdateResource(clientset, resource.ResourceNamespace, resource.VolumeDataSource)
// 	if snapshotErr != nil {
// 		return false, snapshotErr
// 	}
// 	if !isSuccessSnapshot {
// 		return false, fmt.Errorf(resource.VolumeDataSource.VolumeSnapshotSourceName + "에 대한 VolumeSnapshot 생성에 실패하였습니다.")
// 	}

// 	return isSuccessSnapshot, snapshotErr
// }

// // volumeSnapshot :
// func volumeSnapshot(resource nanumv1alpha1.SnapshotSource) (bool, error) {

// 	// ------------------------------
// 	ctrl := crdResource.DynamicVolumeSnapshot{}

// 	// Volume Snapshot
// 	clientset := crd.DynamicInitKube()
// 	// volumeSnapshotSourceKind 는 PersistentVolumeClaim 고정. volumeSnapshotSourceName 는 pvc 이름
// 	isSuccessSnapshot, snapshotErr := ctrl.CreateResource(clientset, resource.ResourceNamespace, resource.VolumeDataSource)
// 	if snapshotErr != nil {
// 		return false, snapshotErr
// 	}
// 	if !isSuccessSnapshot {
// 		return false, fmt.Errorf(resource.VolumeDataSource.VolumeSnapshotSourceName + "에 대한 VolumeSnapshot 생성에 실패하였습니다.")
// 	}

// 	// Snapshot 상태 계속 기다리면서 가져오기. (go routin)
// 	var getError error
// 	var isRunning bool
// 	maxCount := 30
// 	for i := 0; i < maxCount; i++ {
// 		isRunning, getError = ctrl.IsRunningVolumeSnapshot(clientset, resource.ResourceNamespace, resource.VolumeDataSource)
// 		//if getError != nil {
// 		//	return false, getError
// 		//}
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

// // setHistoryVolumeSnapshot 는 Volume Snapshot 한 결과를 etcd 데이터에 넣는 함수입니다.
// func setHistoryVolumeSnapshot(resource nanumv1alpha1.SnapshotSource) (bool, error) {
// 	return true, nil
// }

// func etcdSnapshot(resource nanumv1alpha1.SnapshotSource) (bool, error) {
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

// func setHistoryEtcdSnapshot(resource nanumv1alpha1.SnapshotSource) (bool, error) {
// 	return true, nil
// }
