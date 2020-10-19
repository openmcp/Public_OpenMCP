# migration

## Introduction of migration

This is an example of migrating resources (Deployment, Service, PVC, PV) from Source Cluster to Target Cluster by deploying OpenMCPMigration CRD.

## Requirement
Install [OpenMCP](https://github.com/openmcp/Public_OpenMCP)


## How to Install
1. Build after setting environment variables at pkg/util/config.go
```
$ ./1.build.sh
```

2. Create migration CRD. 
```
$ ./2.create.sh
```

## Sample code

1. Setting migration spec information at 4.example.yaml
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

2.  Deploy Migration resource 
```
kubectl create -f 4.example.yaml
```

