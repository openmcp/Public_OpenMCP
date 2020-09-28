**Table of Contents**
- [Introduction of OpenMCP&reg;](#introduction-of-openmcp)
- [How To Install](#how-to-install)
  - [1. OpenMCP cluster 이름 변경](#1-openmcp-cluster-이름-변경)
  - [2. OpenMCP 기본 모듈 배포](#2-openmcp-기본-모듈-배포)
  - [3. 외부 스토리지에 OpenMCP 서버 등록](#3-외부-스토리지에-openmcp-서버-등록)
- [How To Join Cluster](#how-to-join-cluster)
  - [1. OpenMCP에 On-premise cluster 조인](#1-OpenMCP에-On-premise-cluster-조인)
    - [(1) (선택) cluster 이름 변경 [하위 클러스터에서 수행]](#1-선택-cluster-이름-변경-하위-클러스터에서-수행)
    - [(2) 외부 스토리지에 Join하고자 하는 클러스터 서버 등록 [하위 클러스터에서 수행]](#2-외부-스토리지에-join하고자-하는-클러스터-서버-등록-하위-클러스터에서-수행)
    - [(3) 외부 스토리지에 등록된 하위 클러스터를 OpenMCP에 Join [OpenMCP에서 수행]](#3-외부-스토리지에-등록된-하위-클러스터를-openmcp에-join-openmcp에서-수행)
    - [(4) 하위 클러스터 Master Node에 Region, Zone 등록[OpenMCP에서 수행]](#4-하위-클러스터-master-node에-region-zone-등록-openmcp에서-수행)
    - [(5) 하위 클러스터 MetalLB Config 생성 [OpenMCP에서 수행]](#5-하위-클러스터-metallb-config-생성-openmcp에서-수행)
  - [2. OpenMCP에 GKE cluster 조인](#2-OpenMCP에-GKE-cluster-조인)
    - [(1) Cloud SDK 설치](#1-Cloud-SDK-설치)
    - [(2) gcloud init](#2-gcloud-init)
    - [(3) gcloud container clusters list](#3-gcloud-container-clusters-list)
    - [(4) OpenMCP에 GKE cluster 조인](#4-OpenMCP에-GKE-cluster-조인)
  - [3. OpenMCP에 EKS cluster 조인](#3-OpenMCP에-EKS-cluster-조인)
    - [(1) AWS CLI 설치](#1-AWS-CLI-설치)
    - [(2) aws configure](#2-aws-configure)
    - [(3) aws eks list-clusters](#3-aws-eks-list-clusters)
    - [(4) OpenMCP에 EKS cluster 조인](#4-OpenMCP에-EKS-cluster-조인)
- [OpenMCP EXAMPLE](#openmcp-example)
  - [OpenMCPDeployment 배포](#openmcpdeployment-배포)
  - [OpenMCPService 배포](#openmcpservice-배포)
  - [OpenMCPIngress 배포](#openmcpingress-배포)
  - [OpenMCPDomain,OpenMCPServiceDNSRecord,OpenMCPIngressDNSRecord 배포](#openmcpdomainopenmcpservicednsrecordopenmcpingressdnsrecord-배포)
  - [OpenMCPHybridAutoScaler 배포](#openmcphybridautoscaler-배포)

# Introduction of OpenMCP&reg;

> 지역적(Region/Zone)으로 격리된 컨테이너 간 상호협업하여 유연한 서비스 이동(Migration) 및 컴퓨팅 자원의 끊김없이 자동 확장(Seamless Auto Scaling) 기능을 제공하는 컨테이너 제어·관리 플랫폼

![Architecture of the openmcp](/images/openmcp_architecture.png)

# How To Install

## Requirement

OpenMCP 설치를 위해서는 먼저 `federation`과 nfs를 위한 `외부 서버`가 구축되어 있어야 합니다.

1. [federation](https://github.com/kubernetes-sigs/kubefed/blob/master/docs/userguide.md) 설치
1. [nfs 서버](https://github.com/openmcp/external) 설치

-----------------------------------------------------------------------------------------------
실행 환경
```
OpenMCP   Master IP : 10.0.3.30  
Cluster1  Master IP : 10.0.3.40  
Cluster2  Master IP : 10.0.3.50  
NFS       Server IP : 10.0.3.12  
```

## 1. OpenMCP cluster 이름 변경

kubeconfig 파일에서 클러스터 이름을 `opemncp`로 수정합니다.
> kubeconfig 기본 경로 : $HOME/.kube/config

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

## 2. OpenMCP 기본 모듈 배포  

모듈을 배포하기 전 환경변수 설정을 해줍니다.
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

OpenMCP 동작에 필요한 기본 모듈을 배포합니다.

```bash
$ cd master/
$ ./1.create.sh
```
> 설치 항목
> - Sync Controller
> - Resource Controller (Deployment, HybridAutoScaler, Ingress, Service, Configmap, Secret)
> - LoadBalancing Controller
> - Scheduler
> - Resource Manager (Analytic Engine, Metric Collector)
> - Policy Engine
> - DNS Controller
> - API Server
> - InfluxDB

설치 확인
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
log-version                    3m16s
metric-collector-period        3m16s
```

### OpenMCP Architecture
![Architecture of the openmcp](/images/openmcp_architecture_2.png)


## 3. 외부 스토리지에 OpenMCP 서버 등록

등록하기 전, [omcpctl](https://github.com/openmcp/openmcp/tree/master/omcpctl)를 설치하고 /etc/resolv.conf에 외부 서버를 등록합니다.
```bash
$ vi /etc/resolv.conf
nameserver 10.0.3.12
```

omcpctl 사용하여 nfs 서버에 OpenMCP 서버를 등록합니다.
```bash
$ omcpctl register openmcp
Success OpenMCP Master Register '10.0.3.30'
```

---

# How To Join Cluster

1. OpenMCP에 On-premise cluster 조인
2. OpenMCP에 GKE cluster 조인
3. OpenMCP에 EKS cluster 조인

--------------------------------------------------------------------------------------------
## 1. OpenMCP에 On-premise cluster 조인

### (1) (선택) cluster 이름 변경 [하위 클러스터에서 수행]
OpenMCP에 하위 클러스터를 join하기 전에 다른 클러스터와 이름이 겹치지 않도록 하위 클러스터의 이름을 변경합니다.
> kubeconfig 기본 경로 : $HOME/.kube/config

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

### (2) 외부 스토리지에 Join하고자 하는 클러스터 서버 등록 [하위 클러스터에서 수행]
각 클러스터에 [omctl](https://github.com/openmcp/openmcp-cli) 설치 후, omctl 사용하여 nfs 서버에 join 하고자 하는 클러스터를 등록합니다.
```bash
$ OPENMCP_IP="10.0.3.30"
$ omctl register member ${OPENMCP_IP}
Success Register '10.0.3.40' in OpenMCP Master: 10.0.3.30
```

### (3) 외부 스토리지에 등록된 하위 클러스터를 OpenMCP에 Join [OpenMCP에서 수행]
OpenMCP 서버에서 omcpctl 사용하여 특정 클러스터를 join합니다.
```bash
$ CLUSTER_IP="10.0.3.40"
$ omcpctl join cluster ${CLUSTER_IP}
```

### (4) 하위 클러스터 Master Node에 Region, Zone 등록 [OpenMCP에서 수행]
하위 클러스터의 Label에 Region과 Zone을 등록합니다.
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

### (5) 하위 클러스터 MetalLB Config 생성 [OpenMCP에서 수행]
'metallb_config.yaml' 파일을 새로 생성하여 하위 클러스터 LoadBalancer의 할당 IP 범위를 설정 후 각 클러스터에 배포합니다.
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

## 2. OpenMCP에 GKE cluster 조인

### (1) Cloud SDK 설치
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

### (4) OpenMCP에 GKE cluster 조인
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

## 3. OpenMCP에 EKS cluster 조인

### (1) AWS CLI 설치
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

### (4) OpenMCP에 EKS cluster 조인

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
OpenMCP에 cluster1, cluster2가 조인된 상태에서 EXAMPLE TEST를 진행합니다.
```bash
$ kubectl get kubefedcluster -n kube-federation-system
NAME       READY   AGE
cluster1   True    23h
cluster2   True    23h
```
각 클러스터에 'openmcp' namespaces를 생성합니다.
```bash
$ kubectl create ns openmcp --context=<cluster-name>
```

## OpenMCPDeployment 배포
OpenMCPDeployment를 배포하면 Pod는 스케줄링 되어 Deployment 리소스로 cluster1, cluster2에 배포됩니다.
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

## OpenMCPService 배포
OpenMCPService를 배포하면 Service 리소스가 cluster1, cluster2에 배포됩니다.
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

## OpenMCPIngress 배포
OpenMCPIngress를 배포하면 Target Service가 있는 클러스터를 탐색하여 해당 클러스터에 Ingress 리소스를 배포합니다.
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

## OpenMCPDomain,OpenMCPServiceDNSRecord,OpenMCPIngressDNSRecord 배포
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
## OpenMCPHybridAutoScaler 배포
OpenMCPHybridAutoScaler를 배포하면 Target Deployment가 있는 클러스터를 탐색하여 해당 클러스터에 HorizontalPodAutoscaler, VerticalPodAutoscaler 리소스를 배포합니다.
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

본 프로젝트는 정보통신기술진흥센터(IITP)에서 지원하는 '19년 정보통신방송연구개발사업으로, "컴퓨팅 자원의 유연한 확장 및 서비스 이동을 제공하는 분산·협업형 컨테이너 플랫폼 기술 개발 과제" 임.
