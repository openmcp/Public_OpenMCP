# ketikubecli

## Introduction of KetiKubeCli

> KETI에서 개발한 OpenMCP 플랫폼의 Cluster 자동 Join/Unjoin 기능을 효과적으로 제어할 수 있는 명령어 인터페이스

## Requirement
OpenMCP Master Node에서 ketikubecli를 사용한다면, openmcp 설치가 완료되어야 함.
OpenMCP 하위 Cluster Node들에서 ketikubecli를 사용한다면, 별도의 openmcp 설치가 필요하지 않음. 


## How to Install
```
# 실행 프로그램 만들기
1.build.sh

# 실행 프로그램 경로 지정 및 Config 파일 경로 지정
2.install.sh
```

## Config 파일 설정

> KetiKubeCli는 다음과 같은 설정값(/var/lib/ketikubecli/config.yaml)이 필요합니다.
```
# OpenMCP 설치 경로 지정(OpenMCP Master인 경우만)
openmcpDir: "/root/workspace/openmcp"

# External(nfs) 서버 지정
nfsServer: "10.0.3.12"
```

## Governance

본 프로젝트는 정보통신기술진흥센터(IITP)에서 지원하는 '19년 정보통신방송연구개발사업으로, "컴퓨팅 자원의 유연한 확장 및 서비스 이동을 제공하는 분산·협업형 컨테이너 플랫폼 기술 개발 과제" 임.
