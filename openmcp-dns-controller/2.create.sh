#!/bin/bash
kubectl create ns openmcp
kubectl create -f deploy/crds
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml

