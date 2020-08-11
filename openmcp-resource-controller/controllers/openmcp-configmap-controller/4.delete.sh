#!/bin/bash
cd deploy


kubectl delete -f service_account.yaml
kubectl delete -f role_binding.yaml
kubectl delete -f operator.yaml
#kubectl delete -f crds/cr.yaml

#kubectl delete deploy example-openmcpconfigmap-deploy --context cluster1 -n openmcp
#kubectl delete deploy example-openmcpconfigmap-deploy --context cluster2 -n openmcp
#kubectl delete deploy example-openmcpconfigmap-deploy --context cluster3 -n openmcp
#kubectl delete openmcpconfigmaps example-openmcpconfigmap -n openmcp
#kubectl delete ns openmcp --context cluster1 &
#kubectl delete ns openmcp --context cluster2 &
#kubectl delete ns openmcp --context cluster3 &
#kubectl delete ns openmcp &

#kubectl delete -f crds/crd.yaml
cd ..
