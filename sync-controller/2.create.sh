#!/bin/bash
cd deploy

#kubectl create ns openmcp
#kubectl create ns openmcp --context cluster1
#kubectl create ns openmcp --context cluster2
#kubectl create ns openmcp --context cluster3
kubectl create -f crds/crd.yaml
kubectl create -f service_account.yaml
kubectl create -f role_binding.yaml
kubectl create -f operator.yaml
#kubectl create -f crds/cr_deploy_create.yaml
#kubectl create -f crds/cr_svc_create.yaml

cd ..
