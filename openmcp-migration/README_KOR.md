# migration(kor)

## 마이그레이션

쿠버네티스 기반의 리소스 (Deployment, Service, PVC, PV) 를 소스 클러스터 에서 타겟 클러스터로 배포하는 OpenMCPMigration CRD.

## 시스템 요구사항
Install [OpenMCP] (https://github.com/openmcp/Public_OpenMCP)


## 설치 방법 

1. pkg/util/config.go 파일 User Config 부분을 작성 후 빌드. 
```
$ ./1.build.sh
```

2. migration CRD 배포 
```
$ ./2.create.sh
```

## 사용 방법

1. migration cr spec 작성  (예시 4.example.yaml)
```
apiVersion: openmcp.k8s.io/v1alpha1
kind: Migration
metadata:
  name: migrations
  namespace: openmcp
spec:
  MigrationServiceSource:
  - SourceCluster: cluster1
    TargetCluster: cluster2
    NameSpace: testmig
    ServiceName: testim
    MigrationSource:
    - ResourceName: testim-dp
      ResourceType: Deployment
    - ResourceName: testim-sv
      ResourceType: Service
    - ResourceName: testim-pv
      ResourceType: PersistentVolume
    - ResourceName: testim-pvc
      ResourceType: PersistentVolumeClaim
```

2.  작성한 Migration yaml 파일 배포
```
kubectl create -f 4.example.yaml
```

