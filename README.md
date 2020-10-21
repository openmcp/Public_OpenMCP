**Table of Contents**
- [Introduction of OpenMCP&reg;](#introduction-of-openmcp)
- [How To Install](#how-to-install)
  - [1. Rename OpenMCP cluster name](#1-Rename-openmcp-cluster-name)
  - [2. Deploy init modules for OpenMCP](#2-Deploy-init-modules-for-OpenMCP)
  - [3. Register OpenMCP server on external storage](#3-Register-OpenMCP-server-on-external-storage)
- [How To Join Cluster](#how-to-join-cluster)
  - [1. How to join Kubernetes Cluster to OpenMCP](#1-How-to-join-Kubernetes-Cluster-to-OpenMCP)
    - [(1) (option) Rename cluster name [In sub-cluster]](#1-option-Rename-cluster-name-In-sub-cluster)
    - [(2) Register sub-cluster to external storage [In sub-cluster]](#2-Register-sub-cluster-to-external-storage-In-sub-cluster)
    - [(3) Join sub-cluster to OpenMCP [In OpenMCP]](#3-Join-sub-cluster-to-OpenMCP-In-OpenMCP)
    - [(4) Register Region and Zone to Master Node of sub-cluster [In OpenMCP]](#4-Register-Region-and-Zone-to-Master-Node-of-sub-cluster-In-OpenMCP)
    - [(5) Create MetalLB Config in sub-cluster [In OpenMCP]](#5-Create-MetalLB-Config-in-sub-cluster-In-OpenMCP)
  - [2. How to join GKE Cluster to OpenMCP](#2-How-to-join-GKE-Cluster-to-OpenMCP)
    - [(1) Install Cloud SDK](#1-Install-Cloud-SDK)
    - [(2) gcloud init](#2-gcloud-init)
    - [(3) gcloud container clusters list](#3-gcloud-container-clusters-list)
    - [(4) Join GKE cluster to OpenMCP](#4-Join-GKE-cluster-to-OpenMCP)
  - [3. How to join EKS Cluster to OpenMCP](#3-How-to-join-EKS-Cluster-to-OpenMCP)
    - [(1) Install AWS CLI](#1-Install-AWS-CLI)
    - [(2) aws configure](#2-aws-configure)
    - [(3) aws eks list-clusters](#3-aws-eks-list-clusters)
    - [(4) Join EKS cluster to OpenMCP](#4-Join-EKS-cluster-to-OpenMCP)
- [OpenMCP EXAMPLE](#openmcp-example)
  - [Deploy OpenMCPDeployment resource](#Deploy-OpenMCPDeployment-resource)
  - [Deploy OpenMCPService resource](#Deploy-OpenMCPService-resource)
  - [Deploy OpenMCPIngress resource](#Deplo-OpenMCPIngress-resource)
  - [Deploy OpenMCPDomain,OpenMCPServiceDNSRecord,OpenMCPIngressDNSRecord resource](#Deploy-OpenMCPDomainOpenMCPServiceDNSRecordOpenMCPIngressDNSRecord-resource)
  - [Deploy OpenMCPHybridAutoScaler resource](#Deploy-OpenMCPHybridAutoScaler-resource)

# Introduction of OpenMCP&reg;

> Container control and management platform that provides flexible service movement and seamless automatic scaling of computing resources through mutual collaboration between locally isolated containers

![Architecture of the openmcp](/images/openmcp_architecture_eng.png)

# How To Install

## Requirement

Before you install OpenMCP, 'federation' and 'external server' as NFS server are required.
1. GO v1.14
1. Install [federation](https://github.com/kubernetes-sigs/kubefed/blob/master/docs/userguide.md) (Version: 0.1.0-rc6)
1. Install [nfs server](https://github.com/openmcp/external)

-----------------------------------------------------------------------------------------------
[ Test Environment ]
```
OpenMCP   Master IP : 10.0.3.30  
Cluster1  Master IP : 10.0.3.40  
Cluster2  Master IP : 10.0.3.50  
NFS       Server IP : 10.0.3.12  
```

## 1. Rename OpenMCP cluster name

Rename cluster name to 'openmcp' in kubeconfig file.

> Default directory of kubeconfig file : $HOME/.kube/config

```bash
$ vi $HOME/.kube/config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ...
    server: https://10.0.3.30:6443
  name: openmcp
contexts:
- context:
    cluster: openmcp
    user: openmcp-admin
  name: openmcp
current-context: openmcp
kind: Config
preferences: {}
users:
- name: openmcp-admin
  user:
    client-certificate-data: ...
    client-key-data: ...
```

## 2. Deploy init modules for OpenMCP

Before you deploy modules, set environment variables first.

If you are using a private network, enter the IP and Port through port forwarding for the Public input box.

If you are not using a private network, enter the IP of the OpenMCP Cluster Node in the Public input box.

Enter the IP of the OpenMCP Master Cluster Node in the non-public IP text box except NFS. (It doesn't matter if it's not Public)

PowerDNS Server Port means the PowerDNS query/response port, and enter the value that port port number 53 is forwarded to.

Enter the previously set value for PowerDNS Server API Key. (https://github.com/openmcp/external#how-to-install)
```bash
$ cd ./install_openmcp
$ ./SETTING.sh
OpenMCP Analytic Engine GRPC Server IP -> 10.0.3.20
OpenMCP Analytic Engine GRPC Server Port -> 32050
OpenMCP Metric Collector GRPC Server IP(Public) -> 119.65.195.180
OpenMCP Metric Collector GRPC Server Port(Public) -> 32051
InfluxDB Server IP -> 10.0.3.20
InfluxDB Server Port -> 31051
InfluxDB User Name -> root
InfluxDB User Password -> root
NFS & PowerDNS Server IP -> 10.0.3.12
PowerDNS Server IP(Public) -> 119.65.195.180
PowerDNS Server Port(Public) -> 5353
PowerDNS Server API Key -> 1234
OpenMCP MetalLB Address IP Range (FROM) -> 10.0.3.241
OpenMCP MetalLB Address IP Range (TO) -> 10.0.3.250
```

Deploy all init modules for OpenMCP.

```bash
$ cd master/
$ ./1.create.sh
```
> list of installed Deployment

> - Sync Controller
> - Resource Controller (Deployment, HybridAutoScaler, Ingress, Service, Configmap, Secret)
> - LoadBalancing Controller
> - Scheduler
> - Resource Manager (Analytic Engine, Metric Collector)
> - Policy Engine
> - DNS Controller
> - API Server
> - InfluxDB

Check state of pods.
```bash
$ kubectl get pods -n openmcp
NAME                                                READY   STATUS    RESTARTS   AGE
influxdb-68bff77cbd-f9p6k                           1/1     Running   0          3m11s
openmcp-analytic-engine-cbbcfd7f4-fjrvb             1/1     Running   0          3m12s
openmcp-apiserver-56bf4c7bd-7q4vw                   1/1     Running   0          3m12s
openmcp-configmap-controller-6c97c8cd57-ww7x8       1/1     Running   0          3m11s
openmcp-deployment-controller-747cf6d76-2xr52       1/1     Running   0          3m9s
openmcp-dns-controller-78ff9bcdd5-5fkpq             1/1     Running   0          3m2s
openmcp-has-controller-ccbcd86c4-bxtn8              1/1     Running   0          3m9s
openmcp-ingress-controller-7fc4489594-zmcsl         1/1     Running   0          3m8s
openmcp-loadbalancing-controller-867b79b8d6-bhrkk   1/1     Running   0          3m1s
openmcp-metric-collector-77c5f94759-79tjq           1/1     Running   0          3m10s
openmcp-policy-engine-7c7b5fb7d5-st5qs              1/1     Running   0          3m7s
openmcp-scheduler-75f4bc655-4mzdl                   1/1     Running   0          3m8s
openmcp-secret-controller-6d7c5bf4fc-5mlm2          1/1     Running   0          3m11s
openmcp-service-controller-776cc6574-b8wsm          1/1     Running   0          3m8s
openmcp-sync-controller-769b85d4b4-crnxc            1/1     Running   0          3m1s

$ kubectl get openmcppolicy -n openmcp
NAME                           AGE
analytic-metrics-weight        3m16s
has-target-cluster             3m15s
hpa-minmax-distribution-mode   3m15s
log-level                      3m16s
metric-collector-period        3m16s
```

### OpenMCP Architecture
![Architecture of the openmcp](/images/openmcp_architecture_2.png)


## 3. Register OpenMCP server on external storage

Before you register OpenMCP, install [omcpctl](https://github.com/openmcp/openmcp/tree/master/omcpctl) and add a external server with /etc/resolv.conf.
```bash
$ vi /etc/resolv.conf
nameserver 10.0.3.12
```

And then, register OpenMCP server on NFS server.
```bash
$ omcpctl register openmcp
Success OpenMCP Master Register '10.0.3.30'
```

---

# How To Join Cluster

1. How to join Kubernetes Cluster to OpenMCP
2. How to join GKE Cluster to OpenMCP
3. How to join EKS Cluster to OpenMCP

--------------------------------------------------------------------------------------------
## 1. How to join Kubernetes Cluster to OpenMCP

### (1) (option) Rename cluster name [In sub-cluster]

Before you join the sub-cluster to OpenMCP, check for duplicate cluster names so it does not overlap.

> Default directory of kubeconfig file : $HOME/.kube/config

```bash
$ vi $HOME/.kube/config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ...
    server: https://10.0.3.40:6443
  name: cluster1
contexts:
- context:
    cluster: cluster1
    user: cluster1-admin
  name: cluster1
current-context: cluster1
kind: Config
preferences: {}
users:
- name: cluster1-admin
  user:
    client-certificate-data: ...
    client-key-data: ...
```

### (2) Register sub-cluster to external storage [In sub-cluster]

After installing [omctl](https://github.com/openmcp/openmcp-cli) on each cluster, use omctl to register sub-cluster to nfs server.

```bash
$ OPENMCP_IP="10.0.3.30"
$ omctl register member ${OPENMCP_IP}
Success Register '10.0.3.40' in OpenMCP Master: 10.0.3.30
```

### (3) Join sub-cluster to OpenMCP [In OpenMCP]

From OpenMPC server, use omcptl to join sub-cluster you want to use.
```bash
$ CLUSTER_IP="10.0.3.40"
$ omcpctl join cluster ${CLUSTER_IP}
```

### (4) Register Region and Zone to Master Node of sub-cluster [In OpenMCP]

Tag labels(Region, Zone) on sub-cluster.
```bash
$ kubectl label nodes <node-name> failure-domain.beta.kubernetes.io/region=<region> --context=<cluster-name>
$ kubectl label nodes <node-name> failure-domain.beta.kubernetes.io/zone=<zone> --context=<cluster-name>
```
> Region    
> - AS (Asia)  
> - AF (Africa)  
> - AN (Antarctica)    
> - EU (Europe)    
> - NA (North America)    
> - SA (South America)    

> Zone (ISO 3166-1 alpha-2)  
> - https://ko.wikipedia.org/wiki/ISO_3166-1
---

### (5) Create MetalLB Config in sub-cluster [In OpenMCP]

Create a new file 'metallb_config.yaml' to set the assigned IP range for the LoadBalancer of sub-cluster, and deploy it to the sub-cluster.

```
$ vim metallb_config.yaml 
apiVersion: v1
kind: ConfigMap
metadata:
 namespace: metallb-system
 name: config
data:
 config: |
   address-pools:
   - name: default
     protocol: layer2
     addresses:
     - <<REPLACE_ADDRESS_FROM>>-<<REPLACE_ADDRESS_TO>>


$ kubectl create -f metallb_config.yaml --context=<cluster-name>
```

## 2. How to join GKE Cluster to OpenMCP

### (1) Install Cloud SDK
[https://cloud.google.com/sdk/docs/downloads-apt-get?hl=ko](https://cloud.google.com/sdk/docs/downloads-apt-get?hl=ko)

### (2) gcloud init
```
$ gcloud init
```

### (3) gcloud container clusters list
```
$ gcloud container clusters list
NAME         LOCATION       MASTER_VERSION  MASTER_IP       MACHINE_TYPE  NODE_VERSION   NUM_NODES  STATUS
cluster3     asia-east1-a   1.16.13-gke.1   35.201.135.105  e2-medium     1.16.13-gke.1  2          RUNNING
```


### (4) Join GKE cluster to OpenMCP
```
$ omcpctl join gke-cluster cluster3
...
***** [End] Cluster Join Completed - cluster3 *****

$ kubectl get kubefedclusters -n kube-federation-system
NAME       READY   AGE
cluster1   True    15m
cluster2   True    15m
cluster3   True    10s
```


## 3. How to join EKS Cluster to OpenMCP

### (1) Install AWS CLI
[https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/getting-started-console.html](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/getting-started-console.html)

### (2) aws configure
```
$ aws configure
AWS Access Key ID [****************RCGA]: AKIAJAAK64B5XVB2RCGA
AWS Secret Access Key [****************PivT]: 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
Default region name [us-east-2]: us-east-2
Default output format [json]: json
```

### (3) aws eks list-clusters
```
{
    "clusters": [
        "cluster4"
    ]
}
```

### (4) Join EKS cluster to OpenMCP
```
$ omcpctl join eks-cluster cluster4
...
***** [End] Cluster Join Completed - cluster4 *****

$ kubectl get kubefedclusters -n kube-federation-system
NAME       READY   AGE
cluster1   True    15m
cluster2   True    15m
cluster3   True    80s
cluster4   True    11s
```


# OpenMCP EXAMPLE

Execute 'EXAMPLE TEST' with cluster1, cluster2 joined to OpenMCP.

```bash
$ kubectl get kubefedcluster -n kube-federation-system
NAME       READY   AGE
cluster1   True    23h
cluster2   True    23h
```

Create 'openmcp' namespaces to each clusters.
```bash
$ kubectl create ns openmcp --context=<cluster-name>
```

## Deploy OpenMCPDeployment resource

As you deploy OpenMCPDeployment resource, pods will be scheduled to cluster1, cluster2 as Deployment resource.

```bash
$ kubectl create -f sample/openmcpdeployment/.
```
```bash
$ kubectl get openmcpdeployment -n openmcp
NAME                 AGE
openmcp-deployment   72s

$ kubectl get deploy -n openmcp --context cluster1
NAME                  READY   UP-TO-DATE   AVAILABLE   AGE
openmcp-deployment    2/2     2            2           79s

$ kubectl get deploy -n openmcp --context cluster2
NAME                  READY   UP-TO-DATE   AVAILABLE   AGE
openmcp-deployment    2/2     2            2           80s
```

## Deploy OpenMCPService resource

As you deploy OpenMCPService, service resource will be deployed to cluster1, cluster2.
```bash
$ kubectl create -f sample/openmcpservice/.
```
```bash
$ kubectl get openmcpservice -n openmcp
NAME              AGE
openmcp-service   18s

$ kubectl get service -n openmcp --context cluster1
NAME                    TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
openmcp-service         LoadBalancer   10.110.248.190   10.0.3.200    80:31890/TCP     36s

$ kubectl get service -n openmcp --context cluster2
NAME                    TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
openmcp-service         LoadBalancer   10.103.151.112   10.0.3.180    80:30569/TCP     34s
```

## Deploy OpenMCPIngress resource

When you deploy OpenMCPIngress, it will search for cluster with Target Service and deploy Ingress resource to the cluster.
```bash
$ kubectl create -f sample/openmcpingress/.
```
```bash
$ kubectl get openmcpingress -n openmcp
NAME              AGE
openmcp-ingress   4s

$ kubectl get ingress -n openmcp --context cluster1
NAME              HOSTS                                ADDRESS   PORTS   AGE
openmcp-ingress   cluster1.loadbalancing.openmcp.org             80      18s

$ kubectl get ingress -n openmcp --context cluster2
NAME              HOSTS                                ADDRESS   PORTS   AGE
openmcp-ingress   cluster2.loadbalancing.openmcp.org             80      18s
```

## Deploy OpenMCPDomain,OpenMCPServiceDNSRecord,OpenMCPIngressDNSRecord resource
```bash
$ kubectl create -f sample/openmcpdns/.
```
```bash
$ kubectl get openmcpdomain -n kube-federation-system
NAME                     AGE
openmcp-service-domain   16h
```
```bash
$ kubectl get openmcpservicednsrecord -n openmcp
NAME              AGE
openmcp-service   16h
```
```bash
$ kubectl get openmcpingressdnsrecord -n openmcp
NAME              AGE
openmcp-ingress   16h
```
```bash
$ kubectl get openmcpdnsendpoint -n openmcp
NAME                      AGE
ingress-openmcp-ingress   16h
service-openmcp-service   16h
```
```bash
$ curl -L http://openmcp.service.org
```
## Deploy OpenMCPHybridAutoScaler resource

When you deploy OpenMCPHybridAutoScaler, it will search for cluster with Target Deployment and 
deploy HorizontalPodAutoscaler, VerticalPodAutoscaler resource to the cluster.

```bash
$ kubectl create -f sample/openmcphybridautoscaler/.
```
```bash
$ kubectl get openmcphybridautoscaler -n openmcp
NAME          AGE
openmcp-has   6m51s

$ kubectl get hpa,vpa -n openmcp --context cluster1
NAME                                              REFERENCE                       TARGETS          MINPODS   MAXPODS   REPLICAS   AGE
horizontalpodautoscaler.autoscaling/openmcp-has   Deployment/openmcp-deployment   56/100, 0%/50%   2         4         2          12m

NAME                                                   AGE
verticalpodautoscaler.autoscaling.k8s.io/openmcp-has   11m

$ kubectl get hpa,vpa -n openmcp --context cluster2
NAME                                              REFERENCE                       TARGETS          MINPODS   MAXPODS   REPLICAS   AGE
horizontalpodautoscaler.autoscaling/openmcp-has   Deployment/openmcp-deployment   42/100, 0%/50%   2         4         2          11m

NAME                                                   AGE
verticalpodautoscaler.autoscaling.k8s.io/openmcp-has   11m
```

# Governance

This project was supported by Institute of Information & communications Technology Planning & evaluation (IITP) grant funded by the Korea government (MSIT)
(No.2019-0-00052, Development of Distributed and Collaborative Container Platform enabling Auto Scaling and Service Mobility)
