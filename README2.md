# OpenMCP 설치

## 사전준비  

OpenMCP 설치를 위해서는 federation, ketikubecli 그리고 nfs를 위한 외부 서버가 구축되어 있어야 합니다.

1. [federation](https://github.com/kubernetes-sigs/kubefed/blob/master/docs/userguide.md) 설치
1. ketikubecli 설치
1. nfs 서버
  
## 1. 기본 모듈 배포  

먼저, OpenMCP 동작에 필요한 기본 모듈을 배포합니다.

```
kubectl create ./install_openmcp/master/1.create.sh
```
