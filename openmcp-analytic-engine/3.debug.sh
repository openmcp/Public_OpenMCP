#/bin/bash
NS=openmcp

NAME=$(kubectl get pod -n $NS | grep -E 'openmcp-analytic-engine' | awk '{print $1}')

#echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh

for ((;;))
do
kubectl logs -f -n $NS $NAME --tail 10
done
