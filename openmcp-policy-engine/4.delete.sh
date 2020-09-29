#!/bin/bash

kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/operator.yaml

kubectl delete -f deploy/crds/cr.yaml

#kubectl delete hpa openmcphpa-request-hpa --context cluster1 -n ria
#kubectl delete hpa openmcphpa-request-hpa --context cluster2 -n ria
#kubectl delete hpa openmcphpa-request-hpa --context cluster3 -n ria

#kubectl delete ns openmcp --context cluster1 &
#kubectl delete ns openmcp --context cluster2 &
#kubectl delete ns openmcp --context cluster3 &
#kubectl delete ns openmcp &

#kubectl delete -f deploy/crds/crd.yaml
