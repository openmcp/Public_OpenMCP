#!/bin/bash


kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/operator.yaml 

#kubectl delete -f deploy/service_account.yaml --context cluster2
#kubectl delete -f deploy/role_binding.yaml --context cluster2
#kubectl delete -f deploy/operator.yaml --context cluster2

#kubectl delete -f deploy/service_account.yaml --context cluster3
#kubectl delete -f deploy/role_binding.yaml --context cluster3
#kubectl delete -f deploy/operator.yaml --context cluster3

