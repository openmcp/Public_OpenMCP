# omcpctl

## Introduction of omcpctl

> Command Interface for Master Cluster for Combined Cluster Control of OpenMCP Platform developed by KETI
>
> Process commands such as Join / Create / Delete / Update by requesting data from OpenMCP API Server

## Requirement
1. [External Server](https://github.com/openmcp/external) to Store Cluster Information

2. Install [OpenMCP](https://github.com/openmcp/openmcp)

3. Install nfs-common (apt-get install nfs-common)

4. Install golang (Version: 1.14.2)

## How to Install
1.Build after setting environment variables at build.sh
```
$ vim 1.build.sh
...
OPENMCP_APISERVER="10.0.3.20:31635"                # Specify OpenMCP API Server (ref. kubectl get svc -n openmcp)
OPENMCP_DIR="\/root\/workspace\/openmcp\/openmcp"  # Specifying the OpenMCP installation directory
EXTERNAL_IP="10.0.3.12"                            # Specifying an External(nfs) server
...

$ ./1.build.sh
```

## How to Use
Join process using [omcpctl](https://github.com/openmcp/openmcp/tree/master/omcpctl) for OpenMCP master cluster and [omcpctl](https://github.com/openmcp/openmcp-cli) for OpneMCP member cluster after [installation of OpenMCP master](https://github.com/openmcp/openmcp)
```
1. Registering Openmcp in an OpenMCP Master
  Master) omcpctl register openmcp

2. Register an OpenMCP Member Cluster
  Member) omcpctl register member <OpenMCP_Master_IP>

3. Query currently OpenMCP Joined Clusters
  Master) omcpctl get cluster -n kube-federation-system
            NS           | CLUSTERNAME  | STATUS |   REGION    |     ZONES     |      APIENDPOINT         | PLATFORM |  AGE
+------------------------+--------------+--------+-------------+---------------+--------------------------+----------+-------+
  kube-federation-system | cluster1     | True   | AS          | CN,KR         | https://CLUSTER1_IP:6443 |          |
  kube-federation-system | cluster2     | True   | AS          | CN,KR         | https://CLUSTER2_IP:6443 |          |

4. Current OpenMCP Unjoin Cluster Lookup
  Master) omcpctl joinable list

  CLUSTERNAME  |                               APIENDPOINT                                | PLATFORM
+--------------+--------------------------------------------------------------------------+----------+
  cluster3     | https://CLUSTERIP3_IP:6443                                               |
  eks-cluster1 | https://EKS_CLUSTERIP_IP                                                 | eks

5. Cluster Join and Deploy Base Modules
  Master) omcpctl join cluster <OpenMCP_Member_IP>

6. Query currently OpenMCP Joined Clusters
  Master) omcpctl get cluster -n kube-federation-system

            NS           | CLUSTERNAME  | STATUS |   REGION    |     ZONES     |      APIENDPOINT         | PLATFORM |  AGE
+------------------------+--------------+--------+-------------+---------------+--------------------------+----------+-------+
  kube-federation-system | cluster1     | True   | AS          | CN,KR         | https://CLUSTER1_IP:6443 |          |
  kube-federation-system | cluster2     | True   | AS          | CN,KR         | https://CLUSTER2_IP:6443 |          |
  kube-federation-system | cluster3     | True   | AS          | KR            | https://CLUSTER3_IP:6443 |          |

```
> Get
```
omcpctl get pod
omcpctl get pod -A
omcpctl get pod,svc
omcpctl get pod --context cluster1
omcpctl get pod -o yaml
omcpctl get node
omcpctl get cluster -n kube-federation-system

# OpenMCP Custom Resource
omcpctl get odeploy -A # OpenMCP Deployment
omcpctl get osvc -A # OpenMCP Service
omcpctl get oing -A # OpenMCP Ingress
omcpctl get ohas -A # OpenMCP Hybrid Auto Scaler
omcpctl get opol -A # OpenMCP Policy
omcpctl get ocm -A # OpenMCP ConfigMap
omcpctl get osec -A # OpenMCP Secret
```

> Create
```
omcpctl create -f pod.yaml
omcpctl create -f pod.yaml --context cluster1
```

> Apply
```
omcpctl apply -f pod.yaml
omcpctl apply -f pod.yaml --context cluster1
```

> Delete
```
omcpctl delete pod PODNAME
omcpctl delete pod PODNAME -n NAMESPACE
omcpctl delete pod PODNAME -n NAMESPACE --context cluster1
```



## Governance

This project was supported by Institute of Information & communications Technology Planning & evaluation (IITP) grant funded by the Korea government (MSIT)
(No.2019-0-00052, Development of Distributed and Collaborative Container Platform enabling Auto Scaling and Service Mobility)
