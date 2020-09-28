# OpenMCP API Server

## Introduction of OpenMCP API Server

> KETI에서 개발한 OpenMCP 플랫폼의 연합클러스터 제어를 위한 API 서버
>
> 리소스 GET / PUT / DELETE / POST 요청을 받아드리고 멤버클러스터의 데이터까지 조회 가능
>
> Metric 콜렉터에서 수집한 데이터 조회 API 제공
>
> API-Server는 인증을 위한 JWT 발급 가능
> 
> API-Server는 발급된 JWT를 통해 Token 인증이 가능 (Default: 1시간 후 만료)

## Requirement

1. [OpenMCP 플랫폼](https://github.com/openmcp/openmcp) 설치


## How to Use
OpenMCP 플랫폼 설치 후 가동중인 openmcp-apiserver Pod, Service 리소스 확인
```
kubectl get pod -n openmcp
NAME                                                READY   STATUS    RESTARTS   AGE
...
openmcp-apiserver-bcd8bc7fd-v44pl                   1/1     Running   0          5h27m

NAME                                       TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
...
service/openmcp-apiserver                  LoadBalancer   10.97.128.224   10.0.3.242    8080:31635/TCP   5h27m
...

```


> 리소스 API 참조: (https://v1-17.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/)
>
> Curl 명령어를 통한 Pod 리소스 조회 API Example (NodePort 이용)


```
# Input:
USERNAME="openmcp"
PASSWORD="keti"
IP="10.0.3.20"
PORT="31635"
URL="api/v1/namespaces/default/pods"
CLUSTER="cluster1"

TOKEN_JSON=`curl -XGET -H "Content-type: application/json" "http://$IP:$PORT/token?username=$USERNAME&password=$PASSWORD"`
TOKEN=`echo $TOKEN_JSON | jq .token`
TOKEN=`echo "$TOKEN" | tr -d '"'`

curl -XGET -H "Authorization: Bearer $TOKEN" $IP:$PORT/$URL?clustername=$CLUSTER | jq .

# Reuslt:
{
  "kind": "PodList",
  "apiVersion": "v1",
  "metadata": {
    "selfLink": "/api/v1/namespaces/default/pods",
    "resourceVersion": "12777072"
  },
  "items": [
    {
      "metadata": {
        "name": "nginx-deployment-55fbd9fd6d-h7d8t",
        "generateName": "nginx-deployment-55fbd9fd6d-",
        "namespace": "default",
        "selfLink": "/api/v1/namespaces/default/pods/nginx-deployment-55fbd9fd6d-h7d8t",
        "uid": "3b89fa8b-06b8-4a28-967e-0b22e29e7325",
        "resourceVersion": "9754819",
        "creationTimestamp": "2020-08-18T04:09:41Z",
        "labels": {
          "app": "nginx",
          "pod-template-hash": "55fbd9fd6d"
        },
        "annotations": {
          "cni.projectcalico.org/podIP": "10.244.41.133/32",
          "cni.projectcalico.org/podIPs": "10.244.41.133/32"
        },
        "ownerReferences": [
          {
            "apiVersion": "apps/v1",
            "kind": "ReplicaSet",
            "name": "nginx-deployment-55fbd9fd6d",
            "uid": "d04104a8-277f-464d-8481-e23270820ab9",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
....................
....................
....................
....................
```
> Cluster1의 Node 및 Pod 메트릭 데이터 조회 API Example
```
# Node Metric 조회 Example
> http://10.0.3.20:31635/metrics/nodes/kube1-worker1?clustername=cluster1
> http://10.0.3.20:31635/metrics/nodes/kube1-worker1?clustername=cluster1&timeStart=2020-09-03_09:00:00
> http://10.0.3.20:31635/metrics/nodes/kube1-worker1?clustername=cluster1&timeEnd=2020-09-03_09:00:15
> http://10.0.3.20:31635/metrics/nodes/kube1-worker1?clustername=cluster1&timeStart=2020-09-03_09:00:00&timeEnd=2020-09-03_09:00:15

# Pod Metric 조회 Example
> http://10.0.3.20:31635/metrics/namespaces/default/pods/nginx-deployment-55fbd9fd6d-h7d8t?clustername=cluster1
> http://10.0.3.20:31635/metrics/namespaces/default/pods/nginx-deployment-55fbd9fd6d-h7d8t?clustername=cluster1&timeStart=2020-09-03_09:00:00
> http://10.0.3.20:31635/metrics/namespaces/default/pods/nginx-deployment-55fbd9fd6d-h7d8t?clustername=cluster1&timeEnd=2020-09-03_09:00:15
> http://10.0.3.20:31635/metrics/namespaces/default/pods/nginx-deployment-55fbd9fd6d-h7d8t?clustername=cluster1&timeStart=2020-09-03_09:00:00&timeEnd=2020-09-03_09:00:15
```
> Curl 명령어를 통한 Cluster1의 Node(kube1-worker1) 메트릭 데이터 조회 API Example (NodePort 이용)
```
# Input:
USERNAME="openmcp"
PASSWORD="keti"
IP="10.0.3.20"
PORT="31635"
URL="metrics/nodes/kube1-worker1"
TIME_START="2020-09-03_09:00:00"
TIME_END="2020-09-03_09:00:15"
CLUSTER="cluster1"

TOKEN_JSON=`curl -XGET -H "Content-type: application/json" "http://$IP:$PORT/token?username=$USERNAME&password=$PASSWORD&timeStart=$TIMNE_START&timeEnd=$TIME_END"`
TOKEN=`echo $TOKEN_JSON | jq .token`
TOKEN=`echo "$TOKEN" | tr -d '"'

curl -XGET -H "Authorization: Bearer $TOKEN" $IP:$PORT/$URL?clustername=$CLUSTER  | jq .


#Reuslt:
{
  "nodemetrics": [
    {
      "time": "2020-07-14T07:23:08.917378367Z",
      "cluster": "cluster1",
      "node": "kube1-worker1",
      "cpu": {
        "CPUUsageNanoCores": "218250189n"
      },
      "memory": {
        "MemoryAvailableBytes": "14692232Ki",
        "MemoryUsageBytes": "3595328Ki",
        "MemoryWorkingSetBytes": "1730704Ki"
      },
      "fs": {
        "FsAvailableBytes": "51281400Ki",
        "FsCapacityBytes": "102094168Ki",
        "FsUsedBytes": "45603580Ki"
      },
      "network": {
        "NetworkRxBytes": "32699045919",
        "NetworkTxBytes": "28999023941"
      }
    }
  ]
}

```
 

## Governance

본 프로젝트는 정보통신기술진흥센터(IITP)에서 지원하는 '19년 정보통신방송연구개발사업으로, "컴퓨팅 자원의 유연한 확장 및 서비스 이동을 제공하는 분산·협업형 컨테이너 플랫폼 기술 개발 과제" 임.
