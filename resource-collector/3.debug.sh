#/bin/bash
NS=rescollect

NAME=$(kubectl get pod -n $NS | grep -E 'openmcp-resource-collector' | awk '{print $1}')

#echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
kubectl logs -f -n $NS $NAME

