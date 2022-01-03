#/bin/bash
NS=openmcp


#kubectl exec -it $NAME -n $NS /bin/sh
for ((;;))
do
	NAME=$(kubectl get pod -n $NS | grep -E 'openmcp-scheduler' | awk '{print $1}')

	echo "Exec Into '"$NAME"'"

	kubectl logs -n $NS $NAME --follow 
done
