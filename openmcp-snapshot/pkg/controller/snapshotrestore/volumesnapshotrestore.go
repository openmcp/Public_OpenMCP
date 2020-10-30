package snapshotrestore

import (
	"context"

	// "openmcp/openmcp/migration/pkg/apis"

	nanumv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"

	corev1 "k8s.io/api/core/v1"
	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/kubefed/pkg/controller/util"
	// "openmcp/openmcp/migration/pkg/controller"
)

//volumeSnapshotRun 내에는 PV 만 들어온다고 가정한다.
func volumeSnapshotRestoreRun(r *reconciler, snapshotRestoreSource *nanumv1alpha1.SnapshotRestoreSource, startTime string) error {
	omcplog.V(4).Info(snapshotRestoreSource)
	client := cm.Cluster_genClients[snapshotRestoreSource.ResourceCluster]
	omcplog.V(3).Info("volumeSnapshot Restore Start")

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
	pvResourceOri := &corev1.PersistentVolume{}

	// Key 추출 (crd)
	snapshotKey := snapshotRestoreSource.SnapshotKey
	// Key로 resource Name 추출
	resourceName := util.GetResourceNameBySnapshotKey(snapshotKey)

	pvGetErr := client.Get(context.TODO(), pvResourceOri, "default", resourceName)
	if pvGetErr != nil {
		omcplog.V(3).Info("get pv_info error")
	}
	pvResource := pvResourceOri.DeepCopy()
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
	expvResource := util.GetExternalNfsPVAPI(snapshotKey, *pvResource)
	expvcResource := util.GetExternalNfsPVCAPI(snapshotKey)
	pvcResource := util.GetPVCAPI(snapshotKey, *pvResource)

	targetErr := client.Create(context.TODO(), expvResource)
	if targetErr != nil {
		omcplog.V(3).Info("expvResource create : " + expvResource.Name)
	}
	targetErr = client.Create(context.TODO(), expvcResource)
	if targetErr != nil {
		omcplog.V(3).Info("expvcResource create : " + expvcResource.Name)
	}
	targetErr = client.Create(context.TODO(), pvcResource)
	if targetErr != nil {
		omcplog.V(3).Info("pvcResource create : " + pvcResource.Name)
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
	snapshotCmd := util.GetSnapshotRestoreTemplate(startTime)

	// 잡생성
	jobResource := util.GetJobAPI(snapshotKey, snapshotCmd)
	targetErr = client.Create(context.TODO(), jobResource)
	if targetErr != nil {
		omcplog.V(3).Info("job create : " + jobResource.Name)
	}

	return nil
}
