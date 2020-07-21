#!/bin/bash

kubectl create -f deploy/service_account.yaml --context cluster1
kubectl create -f deploy/role_binding.yaml --context cluster1
kubectl create -f deploy/operator.yaml --context cluster1

kubectl create -f deploy/service_account.yaml --context cluster2
kubectl create -f deploy/role_binding.yaml --context cluster2
kubectl create -f deploy/operator.yaml --context cluster2

kubectl create -f deploy/service_account.yaml --context cluster3
kubectl create -f deploy/role_binding.yaml --context cluster3
kubectl create -f deploy/operator.yaml --context cluster3
