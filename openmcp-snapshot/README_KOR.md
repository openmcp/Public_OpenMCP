# snapshot(kor)

## 스냅샷

쿠버네티스 기반의 리소스 (Deployment, Service, PVC, PV) 들을 스냅샷 또는 스냅샷 복구 할 수 있는 OpenMCPSnapshot, OpenMCPSnapshotRestore CRD.

## 시스템 요구사항
Install [OpenMCP] (https://github.com/openmcp/Public_OpenMCP)


## 설치 방법 

1. pkg/util/config.go 파일 User Config 부분을 작성 후 빌드. 
```
$ ./1.build.sh
```

2. snapshot CRD 배포 
```
$ ./2.create.sh
```

## 사용 방법

##### 스냅샷

1. snapshot cr spec 작성  (예시 4.example_snapshot.yaml)
```
apiVersion: openmcp.k8s.io/v1alpha1
kind: Snapshot
metadata:
  name: example-snapshot
spec:
  snapshotSources:
  - resourceCluster: cluster1
    resourceNamespace: default
    resourceType: pv
    resourceName: test-pv
  - resourceCluster: cluster1
     resourceNamespace: default
    resourceType: pvc
    resourceName: test-pvc
  - resourceCluster: cluster1
    resourceNamespace: default
    resourceType: deployment
    resourceName: test-dp
```

2.  작성한 snapshot yaml 파일 배포
```
kubectl create -f 4.example_snapshot.yaml
```


##### 스냅샷 복구

1. snapshot restore cr spec 작성  (예시 4.example_snapshotrestore.yaml)
```
apiVersion: openmcp.k8s.io/v1alpha1
kind: SnapshotRestore
metadata:
  name: example-snapshotrestore
spec:
  snapshotSources:
  - resourceCluster: cluster1
    resourceNamespace: default
    resourceType: pv
    snapshotKey: pv-test-pv-snapshot (Key)
  - resourceCluster: cluster1
    resourceNamespace: default
    resourceType: pvc
    resourceName: pvc-test-pvc-snapshot
    snapshotKey: pvc-test-pvc-snapshot (Key)
  - resourceCluster: cluster1
    resourceNamespace: default
    resourceType: deployment
    resourceName: dp-test-dp-snapshot
    snapshotKey: dp-test-dp-snapshot (Key)
```

2.  작성한 snapshot yaml 파일 배포
```
kubectl create -f 4.example_snapshot.yaml
```

