#!/bin/bash
helm install . --name mu-rabbit --kube-context=cluster3
kubectl create -f sdr.yaml
