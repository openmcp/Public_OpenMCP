#!/bin/bash

#kubectl create ns openmcp --context cluster1
#kubectl create ns openmcp --context cluster2
#kubectl create ns openmcp --context cluster3

kubectl create -f deploy/ --context cluster1
kubectl create -f deploy/ --context cluster2
kubectl create -f deploy/ --context cluster3
kubectl create -f deploy/ --context cluster4
kubectl create -f deploy/ --context cluster5
kubectl create -f deploy/ --context cluster6
