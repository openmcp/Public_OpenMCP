# omcpctl

## Introduction of omcpctl

> KETI에서 개발한 OpenMCP 플랫폼의 연합클러스터 제어를 위한 Master Cluster용 명령어 인터페이스
>
>OpenMCP API Server에 데이터를 요청하여 

## Requirement
[1. 클러스터 정보를 저장할 External 서버](https://github.com/openmcp/external)

[2. OpenMCP 플랫폼 설치](https://github.com/openmcp/openmcp)

## How to Install
```
# 실행 프로그램 빌드 & 경로 지정 및 Config 파일 경로 지정
./1.build.sh
```

## Config 파일 설정

> omcpctl는 다음과 같은 설정값(/var/lib/omcpctl/config.yaml)이 필요합니다.
```
# OpenMCP API Server 지정 (kubectl get svc -n openmcp 참고)
openmcpAPIServer: "10.0.3.20:31635"

# OpenMCP 설치 디렉토리 지정
openmcpDir: "/root/workspace/openmcp/openmcp"

# External(nfs) 서버 지정
nfsServer: "10.0.3.12"
```

## How to Use
[OpenMCP Master 설치](https://github.com/openmcp/openmcp) 후 OpenMCP Master Cluster용 [omcpctl](https://github.com/openmcp/openmcp/tree/master/omcpctl)과 OpneMCP Member Cluster용 [omcpctl](https://github.com/openmcp/openmcp-cli)을 이용한 Join 과정
```
1. OpenMCP Master Cluster에서 Openmcp 등록 
  Master) omcpctl register master

2. OpenMCP Member Cluster 등록
  Member) omcpctl register member <OpenMCP_Master_IP>

3. 현재 OpenMCP Join된 클러스터 조회
  Master) omcpctl join list

   CLUSTERNAME | STATUS | REGION | ZONES |       APIENDPOINT        | AGE  
 +-------------+--------+--------+-------+--------------------------+-----+
   cluster1    | True   | AS     | CN,KR | https://CLUSTER1_IP:6443 |      
   cluster2    | True   | AS     | CN,KR | https://CLUSTER2_IP:6443 |      


4. 현재 OpenMCP Unjoin(조인가능한) 클러스터 조회 
  Master) omcpctl unjoin list

  CLUSTERNAME |        APIENDPOINT        
+-------------+--------------------------+
  cluster3    | https://CLUSTERIP3_IP:6443  

5. Cluster Join 및 기본 모듈 배포
  Master) omcpctl join cluster <OpenMCP_Member_IP>

6. 현재 OpenMCP Join된 클러스터 조회
  Master) omcpctl join list

   CLUSTERNAME | STATUS | REGION | ZONES |       APIENDPOINT        | AGE  
 +-------------+--------+--------+-------+--------------------------+-----+
   cluster1    | True   | AS     | CN,KR | https://CLUSTER1_IP:6443 |      
   cluster2    | True   | AS     | CN,KR | https://CLUSTER2_IP:6443 |      
   cluster3    | True   | AS     | KR    | https://CLUSTER3_IP:6443 |      
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

본 프로젝트는 정보통신기술진흥센터(IITP)에서 지원하는 '19년 정보통신방송연구개발사업으로, "컴퓨팅 자원의 유연한 확장 및 서비스 이동을 제공하는 분산·협업형 컨테이너 플랫폼 기술 개발 과제" 임.
