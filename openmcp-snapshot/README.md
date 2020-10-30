# snapshot

## Introduction of snapshot

This is an example of Snapshot resources or or snapshot recovery (Deployment, Service, PVC, PV) by deploying OpenMCPMSnapshot CRD.

## Requirement
Install [OpenMCP](https://github.com/openmcp/Public_OpenMCP)
Install [OpenMCP_external](https://github.com/openmcp/Public_OpenMCP_external) etcd server


## How to Install
1. Build after setting environment variables at pkg/util/config.go
```
$ ./1.build.sh
```

2. Create snapshot CRD. 
```
$ ./2.create.sh
```

## Sample code

##### Snapshot

1. Setting snapshot spec information at 4.example_snapshot.yaml
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

2.  Deploy Snapshot resource 
```
kubectl create -f 4.example_snapshot.yaml
```

##### Snapshot restore

1. Setting snapshotRestore spec information at 4.example_snapshotrestore.yaml
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

2.  Deploy SnapshotRestore resource 
```
kubectl create -f 4.example_snapshotrestore.yaml
```