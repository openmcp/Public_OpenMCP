kubectl delete pvc --context cluster2 --all
kubectl delete pvc --context cluster3 --all

kubectl delete -f pv-zookeeper.yaml --context cluster2
kubectl delete -f pv-rabbitmq.yaml --context cluster3

