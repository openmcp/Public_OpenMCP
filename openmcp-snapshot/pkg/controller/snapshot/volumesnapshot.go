package snapshot

import (
	"context"
	"encoding/json"

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
func volumeSnapshotRun(r *reconciler, snapshotSource *nanumv1alpha1.SnapshotSource, startTime string) error {
	omcplog.V(4).Info(json.Marshal(snapshotSource))
	client := cm.Cluster_genClients[snapshotSource.ResourceCluster]
	omcplog.V(3).Info("volumeSnapshot Start")

	/*

		1. 정보추출
		입력값
			PVNAME (예 : demo-app-v0) : pv 의 name
			DATE : 리눅스시간으로 표기 (배쉬는 date '+%s', golang은 time.Now().Unix() -리턴값 int)  - 전체 스냅샷 시작 전에 실행
			PV 전체 yaml 정보 (그대로 다른곳에 띄워야함, job, pvc 연결필요.)

		얻어내는 값
			KEY : DATE-CLUSTER_NAME-pv-PVNAME 로 Key로 사용
			이 Key 를 crd의 status 에 업데이트

		상수
			NFS 용 더미 PVC 정보
			NFS 용 더미 Job 정보  (/data 에 고정 매핑) + externalNFS 도 삽입할 수 있어야함.

	*/
	pvResourceOri := &corev1.PersistentVolume{}

	omcplog.V(4).Info("get PV { Namespace : " + snapshotSource.ResourceNamespace + ", ResourceName : " + snapshotSource.ResourceName)
	pvGetErr := client.Get(context.TODO(), pvResourceOri, snapshotSource.ResourceNamespace, snapshotSource.ResourceName)
	if pvGetErr != nil {
		omcplog.V(3).Info("get pv_info error")
	}
	pvResource := pvResourceOri.DeepCopy()
	//get Date : startTime
	//get PVNAME : snapshotSource.ResourceName
	//get PV yaml Info (mountPath) : pvResource

	// Key 생성 후 snapshotSource.volumeDataSource.VolumeSnapshotID 에 넣기. - 로직 끝난 뒤 reconcile 에서 업데이트.
	omcplog.V(4).Info("makeSnapshotKey { Cluster : " + snapshotSource.ResourceCluster + ", ResourceName : " + snapshotSource.ResourceName)
	snapshotKey := util.MakeVolumeSnapshotKey(startTime, snapshotSource.ResourceCluster, snapshotSource.ResourceName)
	omcplog.V(4).Info("snapshotKey : " + snapshotKey)
	snapshotSource.VolumeDataSource.VolumeSnapshotKey = snapshotKey

	/*
		2. dummy job 생성 및 PV, external NFS 연결

		CLUSTER_NAME, NAMESPACE, PVNAME 를 이용하여
		/home/nfs/storage/CLUSTERNAME/volume/PVNAME/ -> job의 /storage 에 마운트
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
		omcplog.V(3).Info("expvResource create error : " + expvResource.Name)
	}
	targetErr = client.Create(context.TODO(), expvcResource)
	if targetErr != nil {
		omcplog.V(3).Info("expvcResource create error : " + expvcResource.Name)
	}
	targetErr = client.Create(context.TODO(), pvcResource)
	if targetErr != nil {
		omcplog.V(3).Info("pvcResource create error : " + pvcResource.Name)
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
	snapshotCmd := util.GetSnapshotTemplate(string(startTime))

	// 잡생성
	jobResource := util.GetJobAPI(snapshotKey, snapshotCmd)
	targetErr = client.Create(context.TODO(), jobResource)
	if targetErr != nil {
		omcplog.V(3).Info("job create : " + jobResource.Name)
	}

	return nil
}
