#!/bin/bash
cd deploy


kubectl delete -f service_account.yaml
kubectl delete -f role_binding.yaml
kubectl delete -f operator.yaml
kubectl delete -f crds/cr_deploy_create.yaml
kubectl delete -f crds/cr_svc_create.yaml

kubectl delete deploy nginx-deployment --context cluster1 -n openmcp
kubectl delete deploy nginx-deployment --context cluster2 -n openmcp
kubectl delete deploy nginx-deployment --context cluster3 -n openmcp
kubectl delete svc nginx-svc --context cluster1 -n openmcp
kubectl delete svc nginx-svc --context cluster2 -n openmcp
kubectl delete svc nginx-svc --context cluster3 -n openmcp

kubectl delete syncs example-sync -n openmcp
#kubectl delete ns openmcp --context cluster1 &
#kubectl delete ns openmcp --context cluster2 &
#kubectl delete ns openmcp --context cluster3 &
#kubectl delete ns openmcp &

kubectl delete -f crds/crd.yaml

cd ..
