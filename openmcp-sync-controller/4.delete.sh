#!/bin/bash


kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/operator.yaml
kubectl delete -f deploy/crds/cr_deploy_create.yaml
kubectl delete -f deploy/crds/cr_svc_create.yaml

kubectl delete deploy nginx-deployment --context cluster1 -n openmcp
kubectl delete deploy nginx-deployment --context cluster2 -n openmcp
kubectl delete deploy nginx-deployment --context cluster3 -n openmcp
kubectl delete svc nginx-svc --context cluster1 -n openmcp
kubectl delete svc nginx-svc --context cluster2 -n openmcp
kubectl delete svc nginx-svc --context cluster3 -n openmcp

kubectl delete syncs example-sync -n openmcp

kubectl delete -f deploy/crds/crd.yaml

