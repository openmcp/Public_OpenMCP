#!/bin/bash

kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/operator.yaml
#kubectl delete -f deploy/crds/cr.yaml

#kubectl delete service example-openmcpservice --context cluster1 -n openmcp
#kubectl delete service example-openmcpservice --context cluster2 -n openmcp
#kubectl delete service example-openmcpservice --context cluster3 -n openmcp
#kubectl delete openmcpservice example-openmcpservice -n openmcp
#kubectl delete ns openmcp --context cluster1 &
#kubectl delete ns openmcp --context cluster2 &
#kubectl delete ns openmcp --context cluster3 &
#kubectl delete ns openmcp &

#kubectl delete -f deploy/crds/crd.yaml
