# OpenMCP 설치

## 사전준비  

OpenMCP 설치를 위해서는 `federation`, `ketikubecli` 그리고 nfs를 위한 `외부 서버`가 구축되어 있어야 합니다.

1. [federation](https://github.com/kubernetes-sigs/kubefed/blob/master/docs/userguide.md) 설치
1. ketikubecli 설치
1. nfs 서버

## 1. 기본 모듈 배포  

먼저, OpenMCP 동작에 필요한 기본 모듈을 배포합니다.

```bash
./install_openmcp/master/1.create.sh
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

## 2. 클러스터 이름 수정

kubeconfig 파일에서 클러스터 이름을 `opemncp`로 수정합니다.
> 기본 경로 : $HOME/.kube/config

```bash
vi $HOME/.kube/config
```

```bash
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ...
    server: https://10.0.3.40:6443
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
![Architecture of the openmcp](/openmcp_architecture_2.png)
