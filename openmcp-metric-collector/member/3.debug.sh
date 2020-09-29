#/bin/bash
NS=openmcp
CLUSTER=eks-cluster1
NAME=$(kubectl get pod -n $NS --context $CLUSTER | grep -E 'cluster-metric-collector' | awk '{print $1}')

#echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
for ((;;))
do
kubectl logs -f -n $NS $NAME --context $CLUSTER
done

