#/bin/bash
NS=openmcp
controller_name="openmcp-has-controller"

NAME=$(kubectl get pod -n $NS | grep -E $controller_name | awk '{print $1}')

echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
kubectl logs -n $NS $NAME --follow --v=1 --tail=10
