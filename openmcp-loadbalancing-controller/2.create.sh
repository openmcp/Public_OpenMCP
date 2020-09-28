#!/bin/bash
cd deploy

kubectl create -f service_account.yaml
kubectl create -f role_binding.yaml
kubectl create -f operator.yaml


#kubectl expose deployment/openmcp-loadbalancing-controller -n openmcp --port 80 --type=LoadBalancer
#kubectl edit svc loadbalancing-controller -n openmcp

cd ..
