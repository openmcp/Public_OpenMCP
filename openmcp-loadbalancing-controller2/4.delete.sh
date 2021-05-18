#!/bin/sh
cd deploy
kubectl delete -f service_account.yaml
kubectl delete -f role_binding.yaml
kubectl delete -f operator.yaml
kubectl delete -f service.yaml

cd ..
