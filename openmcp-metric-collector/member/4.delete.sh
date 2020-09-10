#!/bin/bash

kubectl delete -f deploy/ --context cluster1
kubectl delete -f deploy/ --context cluster2
kubectl delete -f deploy/ --context cluster3
kubectl delete -f deploy/ --context cluster4
kubectl delete -f deploy/ --context cluster5
kubectl delete -f deploy/ --context cluster6

kubectl delete -f deploy/operator/operator-cluster1.yaml --context cluster1
kubectl delete -f deploy/operator/operator-cluster2.yaml --context cluster2
kubectl delete -f deploy/operator/operator-cluster3.yaml --context cluster3
kubectl delete -f deploy/operator/operator-cluster4.yaml --context cluster4
kubectl delete -f deploy/operator/operator-cluster5.yaml --context cluster5
kubectl delete -f deploy/operator/operator-cluster6.yaml --context cluster6
