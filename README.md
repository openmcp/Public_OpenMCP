**Table of Contents**
- [Introduction of OpenMCP&reg;](#introduction-of-openmcp)
- [How To Install](#how-to-install)
  - [1. Rename OpenMCP cluster name](#1-Rename-openmcp-cluster-name)
  - [2. Deploy init modules for OpenMCP](#2-Deploy-init-modules-for-OpenMCP)
- [How To Join Cluster](#how-to-join-cluster)
  - [1. How to join Kubernetes Cluster to OpenMCP](#1-How-to-join-Kubernetes-Cluster-to-OpenMCP)
    - [(1) (option) Rename cluster name [In sub-cluster]](#1-option-Rename-cluster-name-In-sub-cluster)
    - [(2) Register sub-cluster to OpenMCP [In sub-cluster]](#2-Register-sub-cluster-to-OpenMCP-In-sub-cluster)
    - [(3) Join sub-cluster to OpenMCP [In OpenMCP]](#3-Join-sub-cluster-to-OpenMCP-In-OpenMCP)
    - [(4) Register Region and Zone to Master Node of sub-cluster [In OpenMCP]](#4-Register-Region-and-Zone-to-Master-Node-of-sub-cluster-In-OpenMCP)
  - [2. How to join GKE Cluster to OpenMCP](#2-How-to-join-GKE-Cluster-to-OpenMCP)
    - [(1) Install Cloud SDK](#1-Install-Cloud-SDK)
    - [(2) gcloud init](#2-gcloud-init)
    - [(3) gcloud container clusters list](#3-gcloud-container-clusters-list)
    - [(4) Join GKE cluster to OpenMCP](#4-Join-GKE-cluster-to-OpenMCP)
    - [(5) Register GKE cluster to OpenMCP](#5-Register-GKE-cluster-to-OpenMCP)
  - [3. How to join EKS Cluster to OpenMCP](#3-How-to-join-EKS-Cluster-to-OpenMCP)
    - [(1) Install AWS CLI](#1-Install-AWS-CLI)
    - [(2) aws configure](#2-aws-configure)
    - [(3) aws eks list-clusters](#3-aws-eks-list-clusters)
    - [(4) Join EKS cluster to OpenMCP](#4-Join-EKS-cluster-to-OpenMCP)
    - [(5) Register EKS cluster to OpenMCP](#5-Register-EKS-cluster-to-OpenMCP)
   - [4. How to join AKS Cluster to OpenMCP](#4-How-to-join-AKS-Cluster-to-OpenMCP)
    - [(1) Install Azure CLI](#1-Install-AWS-CLI)
    - [(2) az login](#2-az-login)
    - [(3) az aks list](#3-az-aks-list)
    - [(4) Join AKS cluster to OpenMCP](#4-Join-AKS-cluster-to-OpenMCP)
    - [(5) Register AKS cluster to OpenMCP](#5-Register-AKS-cluster-to-OpenMCP)
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

Before you install OpenMCP, Federation is required.
1. GO v1.14
1. Install [federation](https://github.com/kubernetes-sigs/kubefed/blob/master/docs/userguide.md) (Version: 0.1.0-rc6)

-----------------------------------------------------------------------------------------------
[ Test Environment ]
```
OpenMCP   Master IP : 10.0.3.30  
Cluster1  Master IP : 10.0.3.40  
Cluster2  Master IP : 10.0.3.50  
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
```bash
$ cd ./install_openmcp
$ ./SETTING.sh
OpenMCP Install Type [debug/learning] -> debug
OpenMCP Server IP -> 10.0.3.30
Docker Secret Name for Authentication -> openmcp-private-registry
Docker Registry IP:PORT -> 10.0.3.30:5005
Docker imagePullPolicy [Always/IfNotPresent] -> IfNotPresent
OpenMCP API Server Port -> 30000
OpenMCP Cluster Manager Port ->30001
OpenMCP Analytic Engine GRPC Server Port -> 30002
OpenMCP Metric Collector GRPC Server IP (Public) -> 10.0.3.30
OpenMCP Metric Collector GRPC Server Port (Public) -> 30003
InfluxDB Server Port -> 30004
InfluxDB User Name -> root
InfluxDB User Password -> root
PowerDNS Server IP (Public) -> 10.0.3.30
PowerDNS Server Port (Public) -> 53
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

> - Cluster Manager
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
openmcp-cluster-manager-69c9ccc499-wjcqt            1/1     Running   0          3m12s
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
NAME                            AGE
analytic-metrics-weight         3m16s
has-target-cluster              3m15s
hpa-minmax-distribution-mode    3m15s
loadbalancing-controller-policy 3m15s
log-level                       3m16s
metric-collector-period         3m16s
```

### OpenMCP Architecture
![Architecture of the openmcp](/images/openmcp_architecture_2.png)

---

# How To Join Cluster

1. How to join Kubernetes Cluster to OpenMCP
2. How to join GKE Cluster to OpenMCP
3. How to join EKS Cluster to OpenMCP
4. How to join AKS Cluster to OpenMCP

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

### (2) Register sub-cluster to OpenMCP [In sub-cluster]

Install 'kubectl join' plugin on sub-cluster.
Before execute join command, you must set KUBECONFIG file and ~/.hosts file. 

```bash
$ cd kubectl-plugin
$ ./install_kubectl_join
$ kubectl join ${OPENMCP_IP_PORT}
```

### (3) Join sub-cluster to OpenMCP [In OpenMCP]

Install 'kubectl join-status' plugin on OpenMCP.
```bash
$ cd kubectl-plugin
$ ./install_kubectl_joinstatus
$ kubectl join-status ${CLUSTER_NAME} JOIN
Input MetalLB IP Address Range (FROM) : 10.0.3.251
Input MetalLB IP Address Range (TO) : 10.0.3.260
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


### (4) Register GKE cluster to OpenMCP

Install 'kubectl register' plugin on OpenMCP.
```
$ cd kubectl-plugin
$ ./install_kubectl_register
$ kubectl register GKS ${GKE_CLUSTER_NAME} ${OPENMCP_IP_PORT}
```

### (5) Join GKE cluster to OpenMCP

Install 'kubectl join-status' plugin on OpenMCP.
```
$ cd kubectl-plugin
$ ./install_kubectl_joinstatus
$ kubectl join-status ${GKE_CLUSTER_NAME} JOIN
```

## 3. How to join EKS Cluster to OpenMCP

### (1) Install AWS CLI
[https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/getting-started-console.html](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/getting-started-console.html)

### (2) aws configure

You should execute 'openmcp-cluster-manager' container, and set aws configure.
```
$ kubectl exec -it openmcp-cluster-manager-69c9ccc499-wjcqt -n openmcp bash
bash-4.2# aws configure
AWS Access Key ID [None]: AKIAVJTB7UPJPEMHUAJR
AWS Secret Access Key [None]: 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
Default region name [None]: eu-west-2
Default output format [None]: json
```

### (3) aws eks list-clusters
```
$ aws eks list-clusters
{
    "clusters": [
        "cluster4"
    ]
}
```

### (4) Register EKS cluster to OpenMCP

Install 'kubectl register' plugin on OpenMCP.
```
$ cd kubectl-plugin
$ ./install_kubectl_register
$ kubectl register EKS ${EKS_CLUSTER_NAME} ${OPENMCP_IP_PORT}
```

### (5) Join EKS cluster to OpenMCP

Install 'kubectl join-status' plugin on OpenMCP.
```
$ cd kubectl-plugin
$ ./install_kubectl_joinstatus
$ kubectl join-status ${EKS_CLUSTER_NAME} JOIN
```

## 4. How to join AKS Cluster to OpenMCP

### (1) Install Azure CLI
[https://docs.microsoft.com/ko-kr/cli/azure/install-azure-cli-linux?pivots=apt](https://docs.microsoft.com/ko-kr/cli/azure/install-azure-cli-linux?pivots=apt)

### (2) az login
```
$ az login
To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code S36KSFXTS to authenticate.
[
  {
    "cloudName": "AzureCloud",
    "homeTenantId": "e2bc6150-a5fd-481b-9617-3effff354b44",
    "id": "aa2aa30f-ca48-45d2-aed8-11a3850565bd",
    "isDefault": true,
    "managedByTenants": [],
    "name": "Microsoft Azure",
    "state": "Enabled",
    "tenantId": "e2bc6150-a5fd-481b-9617-3effff354b44",
    "user": {
      "name": "admin@openmcp.onmicrosoft.com",
      "type": "user"
    }
  }
]
```

### (3) az aks list
```
$ az aks list
[
  {
    ...
    "kubernetesVersion": "1.17.11",
    "linuxProfile": null,
    "location": "australiaeast",
    "maxAgentPools": 10,
    "name": "aks-cluster1",
    ...
  }
]
```

### (4) Register AKS cluster to OpenMCP

Install 'kubectl register' plugin on OpenMCP.
```
$ cd kubectl-plugin
$ ./install_kubectl_register
$ kubectl register AKS ${AKS_CLUSTER_NAME} ${OPENMCP_IP_PORT}
```

### (5) Join AKS cluster to OpenMCP

Install 'kubectl join-status' plugin on OpenMCP.
```
$ cd kubectl-plugin
$ ./install_kubectl_joinstatus
$ kubectl join-status ${AKS_CLUSTER_NAME} JOIN
```

# OpenMCP EXAMPLE

```bash
$ kubectl get kubefedcluster -n kube-federation-system
NAME       READY   AGE
cluster1   True    23h
cluster2   True    23h
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
