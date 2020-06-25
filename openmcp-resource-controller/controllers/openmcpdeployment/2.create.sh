#!/bin/bash


kubectl create ns openmcp
kubectl create ns openmcp
kubectl create ns openmcp --context cluster1
kubectl create ns openmcp --context cluster2
kubectl create ns openmcp --context cluster3
kubectl create -f deploy/crds/crd.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml
#kubectl create -f deploy/crds/cr.yaml

