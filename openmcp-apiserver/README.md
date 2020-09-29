# OpenMCP API Server

## Introduction of OpenMCP API Server

> API server for combined cluster control of OpenMCP platform developed by KETI
>
> Resource GET / PUT / DELETE / POST requests can be accepted and data from the member cluster can be viewed.
>
> Provides an API for data retrieval collected from Metric Collector
> 
> API-Server can issue JWT for authentication
> 
> API-Server can authenticate to Tokyo via JWT issued (Default: Expires in 1 hour)

## Requirement

1. Install [OpenMCP Platform](https://github.com/openmcp/openmcp)


## How to Use
Check openmcp-appiserver Pod, Service resources in operation after installing OpenMCP platform
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


> Resource API Reference: (https://v1-17.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/)
>
> Pod Resource Lookup API Example with Curl Command (using NodePort)


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
> Cluster1 Node and Pod Metric Data Lookup API Example
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
> Cluster1 Node (cube1-worker1) Metric Data Lookup API Example (using NodePort)
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

This project was supported by Institute of Information & communications Technology Planning & evaluation (IITP) grant funded by the Korea government (MSIT) (No.2019-0-00052, Development of Distributed and Collaborative Container Platform enabling Auto Scaling and Service Mobility)
