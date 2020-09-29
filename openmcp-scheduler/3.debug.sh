#/bin/bash
NS=openmcp

NAME=$(kubectl get pod -n $NS | grep -E 'openmcp-scheduler' | awk '{print $1}')

echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
for ((;;))
do
	kubectl logs -n $NS $NAME --follow
done
