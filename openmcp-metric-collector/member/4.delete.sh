#!/bin/bash

kubectl delete -f deploy/ --context cluster1
kubectl delete -f deploy/ --context cluster2
kubectl delete -f deploy/ --context cluster3
kubectl delete -f deploy/ --context cluster4
kubectl delete -f deploy/ --context cluster5
kubectl delete -f deploy/ --context cluster6
