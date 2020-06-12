#!/bin/bash
NS=openmcp

NAME=$(kubectl get pod -n $NS | grep -E 'sync-controller' | awk '{print $1}')

echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
for ((;;))
do
kubectl logs --follow -n $NS $NAME
done
