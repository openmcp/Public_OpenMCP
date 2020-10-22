# 고가용성 클러스터 구축

## 구성 환경

Master1 IP : 10.0.3.100 
Master2 IP : 10.0.3.40  
Master3 IP : 10.0.3.30  
LoadBalancer IP : 10.0.3.99

## 구축 순서
1. hostname 설정
1. docker 설치
1. kubeadm/kubectl/kubelet 설치
1. HAProxy 설치
1. Keepalived 설치
1. 클러스터 생성 (kubeadm init)
1. 마스터 노드 추가 (kubeadm join)
1. 마스터 노드 복제 확인

## 1. hostname 설정
### Master1
```
$ cd create_cluster
$ ./1.set_host.sh 
Your Hostname? master1
Please Reconnect shell to Change Hostname!
```
### Master2
```
$ cd create_cluster
$ ./1.set_host.sh 
Your Hostname? master2
Please Reconnect shell to Change Hostname!
```
### Master3
```
$ cd create_cluster
$ ./1.set_host.sh 
Your Hostname? master3
Please Reconnect shell to Change Hostname!
```
## 2. docker 설치
### Master1 / Master2 / Master3
```
$ apt-get update
$ apt-get install apt-transport-https ca-certificates curl gnupg-agent software-properties-common
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
$ add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
$ apt-cache madison docker-ce
$ apt-get install docker-ce=17.03.2~ce-0~ubuntu-xenial

$ docker version
Client:
 Version:      17.03.2-ce
 API version:  1.27
 Go version:   go1.7.5
 Git commit:   f5ec1e2
 Built:        Tue Jun 27 03:35:14 2017
 OS/Arch:      linux/amd64

Server:
 Version:      17.03.2-ce
 API version:  1.27 (minimum version 1.12)
 Go version:   go1.7.5
 Git commit:   f5ec1e2
 Built:        Tue Jun 27 03:35:14 2017
 OS/Arch:      linux/amd64
 Experimental: true
```
## 3. kubeadm/kubectl/kubelet 설치
### Master1 / Master2 / Master3
```
$ apt-get install kubeadm=1.17.3-00
$ apt-get install kubectl=1.17.3-00
$ apt-get install kubelet=1.17.3-00
```
## 4. HAProxy 설치
### Master1 / Master2 / Master3
```
```
## 5. Keepalived 설치

## 6. 클러스터 생성 (kubeadm init)

## 7. 마스터 노드 추가 (kubeadm join)

## 8. 마스터 노드 복제 확인
