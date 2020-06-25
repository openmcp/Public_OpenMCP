#!/bin/bash
project_root_dir=/root/workspace/openmcp/resource-controller
resource_name="openmcphpa"

org_dir=`pwd`
cd $project_root_dir/controllers/$resource_name/deploy

kubectl delete -f service_account.yaml
kubectl delete -f role_binding.yaml
kubectl delete -f operator.yaml
#kubectl delete -f crds/cr.yaml
#kubectl delete -f crds/cr_2.yaml

#kubectl delete hpa openmcphpa-request-hpa --context cluster1 -n ria
#kubectl delete hpa openmcphpa-request-hpa --context cluster2 -n ria
#kubectl delete hpa openmcphpa-request-hpa --context cluster3 -n ria

#kubectl delete hpa openmcphpa-norequest-hpa --context cluster1 -n ria
#kubectl delete hpa openmcphpa-norequest-hpa --context cluster2 -n ria
#kubectl delete hpa openmcphpa-norequest-hpa --context cluster3 -n ria

#kubectl delete vpa openmcphpa-request-vpa --context cluster1 -n ria
#kubectl delete vpa openmcphpa-request-vpa --context cluster2 -n ria
#kubectl delete vpa openmcphpa-request-vpa --context cluster3 -n ria

#kubectl delete vpa openmcphpa-norequest-vpa --context cluster1 -n ria
#kubectl delete vpa openmcphpa-norequest-vpa --context cluster2 -n ria
#kubectl delete vpa openmcphpa-norequest-vpa --context cluster3 -n ria

#kubectl delete deploy example-openmcpdeployment-deploy --context cluster1 -n openmcp
#kubectl delete deploy example-openmcpdeployment-deploy --context cluster2 -n openmcp
#kubectl delete deploy example-openmcpdeployment-deploy --context cluster3 -n openmcp
#kubectl delete openmcpdeployments example-openmcpdeployment -n openmcp
#kubectl delete ns openmcp --context cluster1 &
#kubectl delete ns openmcp --context cluster2 &
#kubectl delete ns openmcp --context cluster3 &
#kubectl delete ns openmcp &

#kubectl delete -f crds/crd.yaml
cd $org_dir
