#!/bin/bash


kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/operator.yaml

