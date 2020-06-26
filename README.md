**Table of Contents**
- [Introduction of OpenMCP&reg;](#introduction-of-openmcp)
- [How To Install](#how-to-install)
  - [1. ketikubecli를 이용한 OpenMCP 서버 등록](#1-ketikubecli를-이용한-openmcp-서버-등록)
    - [(1) `openmcp` namespaces 리소스 생성](#1-openmcp-namespaces-리소스-생성)
    - [(2) cluster 이름 변경](#2-cluster-이름-변경)
    - [(3) 외부 스토리지에 OpenMCP 서버 등록](#3-외부-스토리지에-openmcp-서버-등록)   
  - [2. OpenMCP 기본 모듈 배포](#2-openmcp-기본-모듈-배포)
- [How To Join Cluster](#how-to-join-cluster)
  - [1. (선택) cluster 이름 변경 [하위 클러스터에서 수행]](#1-선택-cluster-이름-변경-하위-클러스터에서-수행)
  - [2. 외부 스토리지에 Join하고자 하는 클러스터 서버 등록 [하위 클러스터에서 수행]](#2-외부-스토리지에-join하고자-하는-클러스터-서버-등록-하위-클러스터에서-수행)
  - [3. 외부 스토리지에 등록된 하위 클러스터를 OpenMCP에 Join [OpenMCP에서 수행]](#3-외부-스토리지에-등록된-하위-클러스터를-openmcp에-join-openmcp에서-수행)
  - [4. 하위 클러스터 Master Node에 Region, Zone 등록](#4-하위-클러스터-master-node에-region-zone-등록)
  - [5. 하위 클러스터 MetalLB Config 생성 (LoadBalancer IP 할당)](#5-하위-클러스터-metallb-config-생성-loadbalancer-ip-할당)
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

OpenMCP 설치를 위해서는 먼저 `federation`, `ketikubecli` 그리고 nfs를 위한 `외부 서버`가 구축되어 있어야 합니다.

1. [federation](https://github.com/kubernetes-sigs/kubefed/blob/master/docs/userguide.md) 설치
1. [ketikubecli](https://github.com/openmcp/openmcp-cli) 설치
1. [nfs 서버](https://github.com/openmcp/external) 설치

## 1. ketikubecli를 이용한 OpenMCP 서버 등록

### (1) `openmcp` namespaces 리소스 생성

```bash
$ kubectl create ns openmcp
```

### (2) cluster 이름 변경

kubeconfig 파일에서 클러스터 이름을 `opemncp`로 수정합니다.
> kubeconfig 기본 경로 : $HOME/.kube/config

```bash
$ vi $HOME/.kube/config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ...
    server: https://10.0.3.20:6443
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

### (3) 외부 스토리지에 OpenMCP 서버 등록
ketikubecli를 사용하여 nfs 서버에 OpenMCP 서버를 등록합니다.
```bash
$ ketikubecli regist openmcp
Success OpenMCP Master Regist '10.0.3.30'
```

## 2. OpenMCP 기본 모듈 배포  

모듈을 배포하기 전 환경변수 설정을 해줍니다.
```bash
$ cd ./install_openmcp
$ ./SETTING.sh
GRPC Server IP -> 10.0.3.30
GRPC Server Port -> 32050
InfluxDB Server IP -> 10.0.3.30
InfluxDB Server Port -> 31051
InfluxDB User Name -> root
InfluxDB User Password -> root
NFS Server IP -> 10.0.3.12
PowerDNS Server IP -> 10.0.3.12
PowerDNS Server Port -> 8081
PowerDNS Server API Key -> 1234
OpenMCP MetalLB Address IP Range (FROM) -> 10.0.3.241
OpenMCP MetalLB Address IP Range (TO) -> 10.0.3.250
```

OpenMCP 동작에 필요한 기본 모듈을 배포합니다.

```bash
$ cd ./install_openmcp/master
$ ./1.create.sh
```
> 설치 항목
> - Sync Controller
> - Resource Controller (Deployment, HybridAutoScaler, Ingress, Service)
> - LoadBalancing Controller
> - Scheduler
> - Resource Manager (Analytic Engine, Metric Collector)
> - Policy Engine
> - DNS Controller
> - InfluxDB

설치 확인
```bash
$ kubectl get pods -n openmcp
NAME                                            READY   STATUS    RESTARTS   AGE
influxdb-68bff77cbd-kdcs4                       1/1     Running   0          21h
loadbalancing-controller-bb7547df8-fpbbj        1/1     Running   0          21h
openmcp-analytic-engine-67dc4b7d9d-kxpb8        1/1     Running   0          21h
openmcp-deployment-controller-747cf6d76-tvm64   1/1     Running   0          21h
openmcp-dns-controller-78ff9bcdd5-lkcx8         1/1     Running   0          21h
openmcp-has-controller-8688867566-bklhw         1/1     Running   0          21h
openmcp-ingress-controller-7fc4489594-jmccz     1/1     Running   0          21h
openmcp-metric-collector-79dc4b466b-5h9wp       1/1     Running   0          21h
openmcp-policy-engine-7c7b5fb7d5-4m4tl          1/1     Running   0          21h
openmcp-scheduler-65794548ff-92fql              1/1     Running   0          21h
openmcp-service-controller-776cc6574-xfd8c      1/1     Running   0          21h
sync-controller-67b4d858d9-4zwnk                1/1     Running   0          21h

$ kubectl get openmcppolicyengine -n openmcp
NAME                           AGE
analytic-metrics-weight        2m1s
hpa-minmax-distribution-mode   2m10s
has-target-cluster             2m6s
```

### OpenMCP Architecture
![Architecture of the openmcp](/images/openmcp_architecture_2.png)

---

# How To Join Cluster
## 1. (선택) cluster 이름 변경 [하위 클러스터에서 수행]
OpenMCP에 하위 클러스터를 join하기 전에 다른 클러스터와 이름이 겹치지 않도록 하위 클러스터의 이름을 변경합니다.
> kubeconfig 기본 경로 : $HOME/.kube/config

```bash
$ vi $HOME/.kube/config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ...
    server: https://10.0.3.30:6443
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

## 2. 외부 스토리지에 Join하고자 하는 클러스터 서버 등록 [하위 클러스터에서 수행]
ketikubecli를 사용하여 nfs 서버에 join 하고자 하는 클러스터를 등록합니다.
```bash
$ OPENMCP_IP="10.0.3.30"
$ ketikubecli regist member --ip ${OPENMCP_IP}
Success Regist '10.0.3.40' in OpenMCP Master: 10.0.3.30
```

## 3. 외부 스토리지에 등록된 하위 클러스터를 OpenMCP에 Join [OpenMCP에서 수행]
OpenMCP 서버에서 ketikubecli를 사용하여 특정 클러스터를 join합니다.
```bash
$ CLUSTER_IP="10.0.3.40"
$ ketikubecli join cluster --ip ${CLUSTER_IP}
```

## 4. 하위 클러스터 Master Node에 Region, Zone 등록
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

## 5. 하위 클러스터 MetalLB Config 생성 (LoadBalancer IP 할당)
하위 클러스터 LoadBalancer의 할당 IP 범위를 설정합니다.

> vim metallb_config.yaml 
```
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
```
> 위에서 만든 Metallb의 Loadbalancer IP 할당범위를 등록합니다.
```  
kubectl create -f metallb_config.yaml
```

# OpenMCP EXAMPLE
OpenMCP에 cluster1, cluster2가 조인된 상태에서 EXAMPLE TEST를 진행합니다.
```bash
$ kubectl get kubefedcluster -n kube-federation-system
NAME       READY   AGE
cluster1   True    23h
cluster2   True    23h
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
