#!/bin/bash

kubectl delete -f deploy/ --context cluster1
kubectl delete -f deploy/ --context cluster2
kubectl delete -f deploy/ --context cluster3
