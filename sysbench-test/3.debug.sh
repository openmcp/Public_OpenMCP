#!/bin/bash
NS=openmcp

NAME=$(kubectl get pod -n $NS --context cluster1 | grep -E 'sysbench-test' | awk '{print $1}')

echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
for ((;;))
do
kubectl logs --follow -n $NS $NAME --context cluster1
done
