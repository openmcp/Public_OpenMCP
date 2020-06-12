# OpenMCP 설치

## 사전준비  

OpenMCP 설치를 위해서는 `federation`, `ketikubecli` 그리고 nfs를 위한 `외부 서버`가 구축되어 있어야 합니다.

1. [federation](https://github.com/kubernetes-sigs/kubefed/blob/master/docs/userguide.md) 설치
1. [ketikubecli](https://github.com/openmcp/openmcp/tree/master/ketikubecli) 설치
1. nfs 서버

## 1. ketikubecli 사용을 위한 환경 설정 

### (1) kubeconfig 파일 수정

kubeconfig 파일에서 클러스터 이름을 `opemncp`로 수정합니다.
> kubeconfig 기본 경로 : $HOME/.kube/config

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

### (2) 외부 스토리지에 OpenMCP 서버 등록
ketikubecli를 사용하여 nfs 서버에 OpenMCP를 등록합니다.
```bash
ketikubecli regist openmcp
```


## 2. 기본 모듈 배포  

OpenMCP 동작에 필요한 기본 모듈을 배포합니다.

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


## 3. 클러스터 Join
### (1) (선택) kubeconfig 파일 수정 - 하위 클러스터에서 수행
OpenMCP에 하위 클러스터를 join하기 전에 클러스터의 이름을 사용자가 원하는 이름으로 변경합니다.
> kubeconfig 기본 경로 : $HOME/.kube/config

```bash
vi $HOME/.kube/config
```

```bash
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

### (2) 외부 스토리지에 Join 클러스터 서버 등록
ketikubecli를 사용하여 nfs 서버에 하위 클러스터를 등록합니다.
```bash
ketikubecli regist cluster1
```

### (3) OpenMCP에 하위 클러스터 Join
ketikubecli를 사용하여 OpenMCP에 하위 클러스터를 join합니다.
```bash
ketikubecli join cluster1
```

## OpenMCP TEST


![Architecture of the openmcp](/openmcp_architecture_2.png)
