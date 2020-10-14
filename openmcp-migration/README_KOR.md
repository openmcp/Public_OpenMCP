# migration(kor)

## 마이그레이션

쿠버네티스 기반의 리소스 (Deployment, Service, PVC, PV) 를 소스 클러스터 에서 타겟 클러스터로 배포하는 OpenMCPMigration CRD.

## 시스템 요구사항
Install [OpenMCP] (https://github.com/openmcp/Public_OpenMCP)


## 설치 방법 

1. Build after setting environment variables at pkg/util/config.go
```
$ ./1.build.sh
```

2. Create migration CRD. 
```
$ ./2.create.sh
```

## Sample code

1. Setting migration information at 4.example.yaml
```
apiVersion: nanum.example.com/v1alpha1
kind: Migration
metadata:
  name: migrations
  namespace: openmcp
spec:
  MigrationServiceSource:
  - SourceCluster: cluster1
    TargetCluster: cluster2
    VolumePath: /nfsdir3/
    NameSpace: openmcp
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

2.  Deploy Migration resource 
```
kubectl create -f 4.example.yaml
```

