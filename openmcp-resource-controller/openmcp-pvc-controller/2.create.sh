#!/bin/bash
cd deploy

kubectl create -f service_account.yaml
kubectl create -f role_binding.yaml
kubectl create -f operator.yaml
