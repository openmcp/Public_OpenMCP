#!/bin/bash
project_root_dir=/root/workspace/usr/lyn
resource_name="openmcpscheduler"

org_dir=`pwd`
cd $project_root_dir/$resource_name/deploy

kubectl create ns openmcp
kubectl create ns openmcp
kubectl create ns openmcp --context cluster1
kubectl create ns openmcp --context cluster2
kubectl create -f crds/crd.yaml
kubectl create -f service_account.yaml
kubectl create -f role_binding.yaml
kubectl create -f operator.yaml

cd $org_dir
