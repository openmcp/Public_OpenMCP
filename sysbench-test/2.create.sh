#!/bin/bash

kubectl create -f deploy/service_account.yaml --context openmcp
kubectl create -f deploy/role_binding.yaml --context openmcp
kubectl create -f deploy/operator.yaml --context openmcp

#kubectl create -f deploy/service_account.yaml --context cluster2
#kubectl create -f deploy/role_binding.yaml --context cluster2
#kubectl create -f deploy/operator.yaml --context cluster2

#kubectl create -f deploy/service_account.yaml --context cluster3
#kubectl create -f deploy/role_binding.yaml --context cluster3
#kubectl create -f deploy/operator.yaml --context cluster3
