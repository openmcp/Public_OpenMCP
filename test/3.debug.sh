#/bin/bash
NS=openmcp

NAME=$(kubectl get pod -n $NS | grep -E 'kchtest' | awk '{print $1}')

#echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
for((;;))
do
kubectl logs -f -n $NS $NAME
done

