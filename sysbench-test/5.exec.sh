#!/bin/bash
NS=openmcp

NAME=$(kubectl get pod -n $NS --context cluster1 | grep -E 'sysbench-test' | awk '{print $1}')

echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
kubectl exec -it $NAME -n $NS --context cluster1 bash

