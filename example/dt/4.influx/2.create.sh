#!/bin/bash
sh secret_info.sh

kubectl create -f deploy --context cluster3
kubectl create -f sdr.yaml
