#!/bin/bash
cd deploy

#kubectl create ns openmcp
kubectl create ns openmcp
kubectl create ns openmcp --context cluster1
kubectl create ns openmcp --context cluster2
kubectl create ns openmcp --context cluster3
#kubectl create -f crds/crd.yaml
kubectl create -f service_account.yaml
kubectl create -f role_binding.yaml
kubectl create -f operator.yaml


#kubectl create -f crds/cr.yaml
#kubectl expose deployment/loadbalancing-controller -n openmcp --port 80 --type=LoadBalancer
#kubectl edit svc loadbalancing-controller -n openmcp

cd ..
