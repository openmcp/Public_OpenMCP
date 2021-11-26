package snapshotrestore

import (
	"context"
	"fmt"
	"regexp"
	"time"

	// "openmcp/openmcp/migration/pkg/apis"

	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"
	config "openmcp/openmcp/openmcp-snapshot/pkg/util"

	"openmcp/openmcp/openmcp-snapshot/pkg/controller/snapshotrestore/resources"
	"openmcp/openmcp/openmcp-snapshot/pkg/util/etcd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/kubefed/pkg/controller/util"
	// "openmcp/openmcp/migration/pkg/controller"
)

//volumeSnapshotRun 내에는 PV 만 들어온다고 가정한다.
func volumeSnapshotRestoreRun(r *reconciler, resourceCluster string, resourceSnapshotKey string, groupSnapshotKey string, volumeSnapshotKey string, pvIdx int) (error, error) {
	//omcplog.V(4).Info(snapshotRestoreSource)
	client := cm.Cluster_genClients[resourceCluster]
	omcplog.V(3).Info("volumeSnapshot Restore Start")

	runType := util.RunTypeSnapshotRestore
	/*

		1. 정보추출

		```
		입력값
			DATE (리눅스시간) : crd의 snapshot key 를 이용하여 DATE 추출 ( str dateStr = strings.Split(key, "-")[0] )
			FINDPATH (예 : /volumeData) : pv 의 spec.nfs.path
			PVNAME (예 : demo-app-v0) : pv 의 name

		얻어내는 값

		상수
			NFS 용 더미 PVC 정보
			NFS 용 더미 Job 정보  (/data 에 고정 매핑)
		```

	*/

	// Key 추출 (crd)
	//resourceSnapshotKey := snapshotRestoreSource.ResourceSnapshotKey
	//resourceCluster := snapshotRestoreSource.ResourceCluster
	// Key로 resource Name 추출
	resourceName, getResourceNameErr := util.GetResourceNameBySnapshotKey(resourceSnapshotKey)
	if getResourceNameErr != nil {
		omcplog.Error(getResourceNameErr)
		return fmt.Errorf("getResourceName error"), getResourceNameErr
	}
	if len(resourceName) > 63 {
		resourceName = resourceName[0:60]
	}
	//snapshotKey = util.GetStartTimeBySnapshotKey(snapshotKey)
	omcplog.V(3).Info("  * resourceName : " + resourceName)
	omcplog.V(3).Info("  * resourceSnapshotKey : " + resourceSnapshotKey)

	// TODO ETCD Get 으로 변경.
	// pvResourceOri := &apiv1.PersistentVolume{}
	// pvGetErr := util.GetPVAPIOri(snapshotKey, pvResource)
	// if pvGetErr != nil {
	// 	omcplog.V(3).Info("get pv_info error")
	// }

	//pvResource := pvResourceOri.DeepCopy()
	//pvResource.Name = resourceName

	pvResource, getEtcdPVErr := getEtcdSnapshotRestoreForPV(r, resourceSnapshotKey, groupSnapshotKey)
	if getEtcdPVErr != nil {
		omcplog.Error(getEtcdPVErr)
		return fmt.Errorf("getEtcdPVErr error"), getEtcdPVErr
	}
	pvResource.ClusterName = resourceCluster
	//pvcResource := getEtcdSnapshotRestoreForPVC(r, snapshotRestoreSource, snapshotKey)
	//get Date : startTime
	//get PVNAME : snapshotRestoreSource.ResourceName
	//get PV yaml Info (mountPath) : pvResource

	/*
		2. dummy job 생성 및 PV, external NFS 연결

		CLUSTER_NAME, NAMESPACE, PVNAME 를 이용하여
		/home/nfs/storage/CLUSTERNAME/pv/PVNAME/ -> job의 /storage 에 마운트
		# externalNFS 와 연결되는 job, pvc, pv 의 이름은 각각 sns-DATE-CLUSTERNAME-job, sns-DATE-CLUSTERNAME-pv, sns-DATE-CLUSTERNAME-pvc

		PV 와 연결할 job, pvc 는 pv 정보를 토대로 연결  (미리 job, pvc yaml 을 상수화시켜놓는다.)
		# PV 와 연결되는 job, pvc 의 이름은 각각  sns-DATE-CLUSTERNAME-volume-job, sns-DATE-CLUSTERNAME-volume-pvc
		path (바인딩은 /data/)
	*/
	expvResource, mountPath := util.GetExternalNfsPVAPI(volumeSnapshotKey, *pvResource, runType)
	expvcResource := util.GetExternalNfsPVCAPI(volumeSnapshotKey, runType)
	pvcResource := util.GetPVCAPI(volumeSnapshotKey, *pvResource, runType)
	oriPvResource := util.GetPVAPI(volumeSnapshotKey, *pvResource, runType)

	targetErr := client.Create(context.TODO(), expvResource)
	if targetErr != nil {
		omcplog.Error(targetErr)
		return fmt.Errorf("expvResource create error[" + expvResource.Name + "]"), targetErr
	} else {
		omcplog.V(3).Info("expvResource create")
	}

	targetErr = client.Create(context.TODO(), expvcResource)
	if targetErr != nil {
		omcplog.Error(targetErr)
		return fmt.Errorf("expvcResource create error[" + expvcResource.Name + "]"), targetErr
	} else {
		omcplog.V(3).Info("expvcResource create")
	}

	targetErr = client.Create(context.TODO(), pvcResource)
	if targetErr != nil {
		omcplog.Error(targetErr)
		return fmt.Errorf("pvcResource create error[" + pvcResource.Name + "]"), targetErr
	} else {
		omcplog.V(3).Info("pvcResource create")
	}

	targetErr = client.Create(context.TODO(), oriPvResource)
	if targetErr != nil {
		omcplog.Error(targetErr)
		return fmt.Errorf("oriPvResource create error[" + oriPvResource.Name + "]"), targetErr
	} else {
		omcplog.V(3).Info("oriPvResource create")
	}

	/*
	   3. externalNFS 가 붙은 Deploy 에서 스냅샷 명령실행

	   ```
	   # 1. externalNFS 에서 해당 deploy 로 지정된 스냅샷 폴더로 이동한다,
	   cd /storage    # externalNFS 의 /home/nfs/storage/CLUSTERNAME/volume/PVNAME/ 와 마운트됨
	   export lastDir=`ls -tr | tail -1`  #가장 최근 스냅샷 폴더

	   # 2. newerthan, 을 구한다. 폴더가 비어있을 경우 newerthan 는 1970년1월1일이다.
	   export newerthan=`date +"%F %T" --date @0`  #초기화
	   if [ -n "$lastDir" ]; then
	     export newerthan=`date +"%F %T" --date @$lastDir`   #가장 최근에 스냅샷한 시간
	   fi

	   # 3. olderthan 을 구한다. olderthan 은 현재 시간(리눅스시간) 이다. -> 이것은 코드상에서 계산에서 넣도록한다.
	   export olderthan=`date '+%F %T' --date @!DATE`      # 스냅샷 시작 시간

	   # 4. newerthan, olderthan 을 이용하여 파일 찾아서 압축   #/data 인 이유는 PV 에 연결된 /data가 여깃음
	   find /data -type f -newermt "$newerthan" ! -newermt "$olderthan" | xargs tar cvf !DATE
	   ```
	*/

	// startTime 를 이용하여 cmd 내용 작성
	snapshotCmd, getTmpErr := util.GetSnapshotRestoreTemplate(groupSnapshotKey, mountPath)
	if getTmpErr != nil {
		omcplog.Error(getTmpErr)
		return fmt.Errorf("get GetSnapshotRestoreTemplate error"), getTmpErr
	}
	snapshotTmpCmd, getTmpErr := util.GetLoopForSuccessTemplate()
	if getTmpErr != nil {
		omcplog.Error(getTmpErr)
		return fmt.Errorf("get GetLoopForSuccessTemplate error"), getTmpErr
	}

	// 잡생성
	jobResource := util.GetJobAPI(volumeSnapshotKey, snapshotTmpCmd, runType)
	targetErr = client.Create(context.TODO(), jobResource)
	if targetErr != nil {
		omcplog.Error(targetErr)
		return fmt.Errorf("job create error : " + jobResource.Name), targetErr
	} else {
		omcplog.V(3).Info("jobResource create")
	}

	targetListClient := *cm.Cluster_kubeClients[resourceCluster]
	timeoutcheck := 0

	checkResourceName := jobResource.Name
	isSnapshotRestoreCompleted := false
	omcplog.V(4).Info("connecting... : " + checkResourceName)
	podName := ""
	for !isSnapshotRestoreCompleted {
		var regexpErr error
		matchResult := false
		pods, _ := targetListClient.CoreV1().Pods(config.JOB_NAMESPACE).List(context.TODO(), metav1.ListOptions{})
		for _, pod := range pods.Items {
			matchResult, regexpErr = regexp.MatchString(checkResourceName, pod.Name)

			if regexpErr == nil && matchResult == true {
				if pod.Name != "" {
					podName = pod.Name
				} else {
					podName = checkResourceName + "-unknown"
				}
				omcplog.V(3).Info("TargetCluster PodName : " + podName + "/" + string(pod.Status.Phase))
				if pod.Status.Phase == corev1.PodRunning {
					isSnapshotRestoreCompleted = true
					omcplog.V(3).Info(podName + " is Running.")
				}
			}
		}

		if timeoutcheck == 30 {
			//시간초과 - 오류 루틴으로 진입
			omcplog.V(3).Info("long time error...")

			//1. 이벤트 리소스를 통한 오류 검출. 이벤트에서 해당 오류 찾아서 도출. Reason에 SuccessfulCreate가 포함된 경우는 CMD 에서 오류를 찾아야한다.
			errDetail, eventErr := config.FindErrorForEvent(&targetListClient, jobResource.Name)
			if errDetail == "" {
				return fmt.Errorf("VolumeSnapshotRestore Failed. FindErrorForEvent error"), eventErr
			}

			omcplog.Error(errDetail)
			omcplog.Error("VolumeSnapshotRestore Failed")
			return fmt.Errorf("Job runs long. VolumeSnapshotRestore Failed"), fmt.Errorf(errDetail)
		}

		if !isSnapshotRestoreCompleted {
			timeoutcheck = timeoutcheck + 5
			time.Sleep(time.Second * 5)
			omcplog.V(4).Info("connecting...")
		}
	}

	//Job Command 실행
	restconfig := cm.Cluster_configs[resourceCluster]
	cmdResult, commandErr := config.RunCommand(&targetListClient, restconfig, podName, snapshotCmd, config.JOB_NAMESPACE)
	if nil != commandErr {
		return fmt.Errorf("Command run  error"), commandErr
	}

	omcplog.V(4).Info("=====in======")
	omcplog.V(4).Info(fmt.Sprintf("%s", cmdResult.Stdin.Bytes()))
	omcplog.V(4).Info("=====out======")
	omcplog.V(4).Info(fmt.Sprintf("%s", cmdResult.Stdout.Bytes()))
	omcplog.V(4).Info("=====error=======")
	omcplog.V(4).Info(fmt.Sprintf("%s", cmdResult.Stderr.Bytes()))

	return nil, nil
}

//volumeSnapshotRun 내에는 PV 만 들어온다고 가정한다.
func getEtcdSnapshotRestoreForPV(r *reconciler, resourceSnapshotKey string, startTime string) (*apiv1.PersistentVolume, error) {
	omcplog.V(4).Info("# getEtcdSnapshotRestoreForPV")
	snapshotKeyAllPath := resourceSnapshotKey

	//ETCD 에서 데이터 가져오기.
	etcdCtl, etcdInitErr := etcd.InitEtcd()
	if etcdInitErr != nil {
		omcplog.Error("etcdsnapshot.go : Etcd Init Err")
		return &corev1.PersistentVolume{}, etcdInitErr
	}
	omcplog.V(2).Info("snapshotKeyAllPath : " + snapshotKeyAllPath)
	resp, etcdGetErr := etcdCtl.Get(snapshotKeyAllPath)
	if etcdGetErr != nil {
		omcplog.Error("etcdsnapshotresource.go : Etcd Get Err")
		return &corev1.PersistentVolume{}, etcdGetErr
	}
	resourceJSONString := string(resp.Kvs[0].Value)

	resourceObj, err := resources.JSON2Pv(resourceJSONString)
	if err != nil {
		omcplog.Error("CreateResource for JSON error")
		return &corev1.PersistentVolume{}, err
	}
	return resourceObj, nil
}

//volumeSnapshotRun 내에는 PV 만 들어온다고 가정한다.
func getEtcdSnapshotRestoreForPVC(r *reconciler, resourceSnapshotKey string, startTime string) (*apiv1.PersistentVolumeClaim, error) {
	omcplog.V(4).Info("# getEtcdSnapshotRestoreForPVC")
	snapshotKeyAllPath := resourceSnapshotKey

	//ETCD 에서 데이터 가져오기.
	etcdCtl, etcdInitErr := etcd.InitEtcd()
	if etcdInitErr != nil {
		omcplog.Error("etcdsnapshot.go : Etcd Init Err")
		return &corev1.PersistentVolumeClaim{}, etcdInitErr
	}
	omcplog.V(2).Info("snapshotKeyAllPath : " + snapshotKeyAllPath)
	resp, etcdGetErr := etcdCtl.Get(snapshotKeyAllPath)
	if etcdGetErr != nil {
		omcplog.Error("etcdsnapshotresource.go : Etcd Get Err")
		return &corev1.PersistentVolumeClaim{}, etcdGetErr
	}
	resourceJSONString := string(resp.Kvs[0].Value)

	resourceObj, err := resources.JSON2Pvc(resourceJSONString)
	if err != nil {
		omcplog.Error("CreateResource for JSON error")
		return &corev1.PersistentVolumeClaim{}, err
	}
	return resourceObj, nil
}
