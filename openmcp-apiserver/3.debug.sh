#!/bin/bash
NS=openmcp
controller_name="openmcp-apiserver"

NAME=$(kubectl get pod -n $NS | grep -E $controller_name | awk '{print $1}')

echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
for ((;;))
do
kubectl logs -n $NS $NAME --follow --tail=10
done

