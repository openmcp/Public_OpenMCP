# ketikubecli

## Introduction of KetiKubeCli;

> KETI에서 개발한 OpenMCP 플랫폼의 Cluster 자동 Join/Unjoin 기능을 효과적으로 제어할 수 있는 명령어 인터페이스

## Requirement



## How to Install
<pre>
<code>
# 실행 프로그램 만들기
1.build.sh

# 실행 프로그램 경로 지정 및 Config 파일 경로 지정
2.install.sh
</code>
</pre>

## Config 파일 설정

> KetiKubeCli는 다음과 같은 설정값(/var/lib/ketikubecli/config.yaml)이 필요합니다.
```
# OpenMCP 설치 경로 지정
openmcpDir: "/root/workspace/openmcp"

# External(nfs) 서버 지정
nfsServer: "10.0.3.12"
```

## Governance

본 프로젝트는 정보통신기술진흥센터(IITP)에서 지원하는 '19년 정보통신방송연구개발사업으로, "컴퓨팅 자원의 유연한 확장 및 서비스 이동을 제공하는 분산·협업형 컨테이너 플랫폼 기술 개발 과제" 임.
