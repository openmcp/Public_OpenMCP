**Table of Contents**
- [Introduction of OpenMCP&reg;](#introduction-of-openmcp)
- [How To Install](#how-to-install)
  - [1. Rename OpenMCP cluster name](#1-Rename-openmcp-cluster-name)
  - [2. Deploy init modules for OpenMCP](#2-Deploy-init-modules-for-OpenMCP)
- [How To Join Cluster](#how-to-join-cluster)
  - [1. How to join Kubernetes Cluster to OpenMCP](#1-How-to-join-Kubernetes-Cluster-to-OpenMCP)
    - [(1) (option) Rename cluster name [In sub-cluster]](#1-option-Rename-cluster-name-In-sub-cluster)
    - [(2) Check Status OpenMCP API Server in openmcp-cluster [In openmcp-cluster]](#2-Check-Status-OpenMCP-API-Server-In-openmcp-cluster)
    - [(3) Register DNS Server at "/etc/resolv.conf" [In sub-cluster]](#3-Register-DNS-Server-at-etc/resolv.conf-In-sub-cluster)
    - [(4) Register Region and Zone to All Nodes of sub-cluster [In sub-cluster]](#4-Register-Region-and-Zone-to-All-Nodes-of-sub-cluster-In-sub-cluster)
    - [(5) Register sub-cluster to OpenMCP [In sub-cluster]](#5-Register-sub-cluster-to-OpenMCP-In-sub-cluster)
    - [(6) Check Registered OpenMCPCluster [In openmcp-cluster]](#6-Check-Registered-OpenMCPCluster-In-openmcp-cluster)
    - [(7) Join sub-cluster to OpenMCP [In openmcp-cluster]](#7-Join-sub-cluster-to-OpenMCP-In-openmcp-cluster)
    
  - [2. How to join GKE Cluster to OpenMCP](#2-How-to-join-GKE-Cluster-to-OpenMCP)
    - [(1) Install Cloud SDK](#1-Install-Cloud-SDK)
    - [(2) gcloud init](#2-gcloud-init)
    - [(3) gcloud container clusters list](#3-gcloud-container-clusters-list)
    - [(4) Register GKE cluster to OpenMCP](#4-Register-GKE-cluster-to-OpenMCP)
    - [(5) Join GKE cluster to OpenMCP](#5-Join-GKE-cluster-to-OpenMCP)
  - [3. How to join EKS Cluster to OpenMCP](#3-How-to-join-EKS-Cluster-to-OpenMCP)
    - [(1) Install AWS CLI](#1-Install-AWS-CLI)
    - [(2) aws configure](#2-aws-configure)
    - [(3) aws eks list-clusters](#3-aws-eks-list-clusters)
    - [(4) Register EKS cluster to OpenMCP](#4-Register-EKS-cluster-to-OpenMCP)
    - [(5) Join EKS cluster to OpenMCP](#5-Join-EKS-cluster-to-OpenMCP)
  - [4. How to join AKS Cluster to OpenMCP](#4-How-to-join-AKS-Cluster-to-OpenMCP)
    - [(1) Install Azure CLI](#1-Install-AWS-CLI)
    - [(2) az login](#2-az-login)
    - [(3) az aks list](#3-az-aks-list)
    - [(4) Register AKS cluster to OpenMCP](#4-Register-AKS-cluster-to-OpenMCP)
    - [(5) Join AKS cluster to OpenMCP](#5-Join-AKS-cluster-to-OpenMCP)
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

Before you install OpenMCP, Federation and OpenMCP ExternalServer are required.
1. GO v1.14
1. [Install federation](https://github.com/kubernetes-sigs/kubefed/blob/master/docs/userguide.md) (Version: 0.8.1)
1. [ExternalServer](https://github.com/openmcp/external)
1. [flannel CNI](https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml)

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
$ vim settings.yaml
default:
  # learning or debug
  installType: debug
  docker:
    imagePullSecretName: regcred
    imagePullPolicy: Always
master:
  internal:
    ip: 10.0.3.30
    ports:
      apiServerPort: 30000
      clusterManagerPort: 30001
      analyticEnginePort: 30002
      metricCollectorPort: 30003
      influxDBPort: 30004
  public:
    ip: 119.65.195.180
    ports:
      metricCollectorPort: 3212
  APIServer:
      AppKey: openmcp-apiserver
      UserName: openmcp
      UserPW: keti
  metalLB:
    rangeStartIP: 10.0.3.191
    rangeEndIP: 10.0.3.200

powerDNS:
  apiKey: 1234
  internal:
    ip: 10.0.3.12
    ports:
      pdnsPort: 53
  public:
    ip: 119.65.195.180
    ports:
      pdnsPort: 5353


$ ./create.sh
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

### (2) Check Status OpenMCP API Server [In openmcp-cluster]
```bash
$ kubectl get pod,svc -n openmcp
ME                                                    READY   STATUS    RESTARTS   AGE
...
pod/openmcp-apiserver-ddf89465f-b77t9                   1/1     Running   0          54s
...

NAME                                       TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
...
service/openmcp-apiserver                  LoadBalancer   10.96.65.179     XX.XX.XX.XX   8080:30000/TCP   54s
```

### (3) Register DNS Server at "/etc/resolv.conf" [In sub-cluster]

Write the [<EXTERNAL_SERVER_IP>](https://github.com/openmcp/external) of the external server at the top of the '/etc/resolv.conf' file.
This is to join the OpenMCP cluster through DNS Domain on the API server.
  
```bash
$ vim /etc/resolv.conf
nameserver <EXTERNAL_SERVER_IP>
```

### (4) Register Region and Zone to All Nodes of sub-cluster [In sub-cluster]

Tag labels(Region, Zone) on sub-cluster.
```bash
$ kubectl label nodes <node-name> topology.kubernetes.io/region=<region> 
$ kubectl label nodes <node-name> topology.kubernetes.io/zone=<zone>
$ kubectl label nodes <node-name> topology.istio.io/subzone=<cluster-name>
```
> Region (ISO 3166-1 alpha-2) https://ko.wikipedia.org/wiki/ISO_3166-1
> - KR (Korea)  
> - US (USA)  
> - CH (China)    
> - JP (Japan)    
> - IN (India)    

> Zone (locationInfo.csv)  
> - https://github.com/openmcp/Public_OpenMCP/blob/master/locationInfo.csv

### (5) Register sub-cluster to OpenMCP [In sub-cluster]

Install 'kubectl request-join' plugin on sub-cluster.
Before execute join command, you must set KUBECONFIG file and ~/.hosts file. 

```bash
$ cd kubectl_plugin
$ chmod +x kubectl-request_join
$ cp kubectl-request_join /usr/local/bin
$ kubectl request-join
```
  
### (6) Check Registered OpenMCPCluster [In openmcp-cluster]
  
```bash
$ kubectl get openmcpcluster -n openmcp
  
NAME       STATUS
cluster1   UNJOIN
cluster2   UNJOIN

```
### (7) Join sub-cluster to OpenMCP [In openmcp-cluster]

Install 'kubectl join' plugin on sub-cluster.
Before execute join command, you must set KUBECONFIG file and ~/.hosts file.

```bash
$ cd kubectl_plugin
$ chmod +x kubectl-join
$ cp kubectl-join /usr/local/bin
$ kubectl join <CLUSTERNAME> <PLATFORMTYPE>
```

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

Install 'kubectl regist' plugin on OpenMCP.
```
$ cd kubectl_plugin
$ chmod +x kubectl-regist_join
$ cp kubectl-regist_join /usr/local/bin
$ kubectl regist-join GKE ${GKE_CLUSTER_NAME} ${OPENMCP_IP_PORT}
```

### (5) Join GKE cluster to OpenMCP

Install 'kubectl join' plugin on OpenMCP.
```
$ cd kubectl_plugin
$ chmod +x kubectl-join
$ cp kubectl-join /usr/local/bin
$ kubectl join <CLUSTERNAME> <PLATFORMTYPE>
```

## 3. How to join EKS Cluster to OpenMCP

### (1) Install AWS CLI
[https://docs.aws.amazon.com/cli/latest/userguide/install-linux.html](https://docs.aws.amazon.com/cli/latest/userguide/install-linux.html)

### (2) aws configure

```
$ aws configure
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

Install 'kubectl regist' plugin on OpenMCP.
```
$ cd kubectl_plugin
$ chmod +x kubectl-regist_join
$ cp kubectl-regist_join /usr/local/bin
$ kubectl regist-join EKS ${EKS_CLUSTER_NAME} ${OPENMCP_IP_PORT}
```

### (5) Join EKS cluster to OpenMCP

Before trying to join EKS cluster, you should execute 'openmcp-cluster-manager' container, and set aws configure.
```
$ kubectl exec -it openmcp-cluster-manager-69c9ccc499-wjcqt -n openmcp bash
bash-4.2# aws configure
AWS Access Key ID [None]: AKIAVJTB7UPJPEMHUAJR
AWS Secret Access Key [None]: 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
Default region name [None]: eu-west-2
Default output format [None]: json
```

Install 'kubectl join' plugin on OpenMCP.
```
$ cd kubectl_plugin
$ chmod +x kubectl-join
$ cp kubectl-join /usr/local/bin
$ kubectl join <CLUSTERNAME> <PLATFORMTYPE>
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

Install 'kubectl regist' plugin on OpenMCP.
```
$ cd kubectl_plugin
$ chmod +x kubectl-regist_join
$ cp kubectl-regist_join /usr/local/bin
$ kubectl regist-join AKS ${AKS_CLUSTER_NAME} ${OPENMCP_IP_PORT}
```

### (5) Join AKS cluster to OpenMCP

Install 'kubectl join' plugin on OpenMCP.
```
$ cd kubectl_plugin
$ chmod +x kubectl-join
$ cp kubectl-join /usr/local/bin
$ kubectl join <CLUSTERNAME>
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
