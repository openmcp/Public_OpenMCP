#!/bin/bash

kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f dedploy/operator.yaml
kubectl delete -f deploy/crds/cr.yaml

kubectl delete deploy example-openmcpdeployment --context cluster4 -n openmcp
kubectl delete deploy example-openmcpdeployment --context cluster5 -n openmcp
kubectl delete deploy example-openmcpdeployment --context cluster6 -n openmcp
kubectl delete openmcpdeployments example-openmcpdeployment -n openmcp
#kubectl delete ns openmcp --context cluster1 &
#kubectl delete ns openmcp --context cluster2 &
#kubectl delete ns openmcp --context cluster3 &
#kubectl delete ns openmcp &

#kubectl delete -f deploy/crds/crd.yaml
