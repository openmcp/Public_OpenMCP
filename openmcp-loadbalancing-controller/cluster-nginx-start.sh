#!/bin/sh

kubectl apply -f ./deploy/cluster/cluster-nginx.yaml --context=cluster1
kubectl apply -f ./deploy/cluster/cluster-nginx.yaml --context=cluster2
kubectl apply -f ./deploy/cluster/cluster-nginx.yaml --context=cluster3
