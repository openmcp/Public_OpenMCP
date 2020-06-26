#!/bin/sh


echo "**************************************"
echo "   Start Cluster Ingress Controller   "
echo "**************************************"
./cluster-nginx-start.sh

echo "            Start Test Pod            "
echo "**************************************"

kubectl create -f ./example/IoT-Data-Collector/yaml/IoT-Data-Collector-dep.yaml -n openmcp --context=cluster1
kubectl create -f ./example/IoT-Data-Collector/yaml/IoT-Data-Collector-svc.yaml -n openmcp --context=cluster1
kubectl create -f ./example/IoT-Data-Collector/yaml/IoT-Data-Collector-dep.yaml -n openmcp --context=cluster2
kubectl create -f ./example/IoT-Data-Collector/yaml/IoT-Data-Collector-svc.yaml -n openmcp --context=cluster2
kubectl create -f ./example/IoT-Data-Collector/yaml/IoT-Data-Collector-dep.yaml -n openmcp --context=cluster3
kubectl create -f ./example/IoT-Data-Collector/yaml/IoT-Data-Collector-svc.yaml -n openmcp --context=cluster3

echo "             Start Ingress            "
echo "**************************************"

kubectl create -f ./example/example_ingress.yaml -n openmcp --context=cluster1
kubectl create -f ./example/example_ingress.yaml -n openmcp --context=cluster2
kubectl create -f ./example/example_ingress.yaml -n openmcp --context=cluster3



