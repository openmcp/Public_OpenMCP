kubectl create -f ns.yaml
kubectl apply -f policy/scheduling-policy.yaml
kubectl create -f bookinfo
kubectl create -f gw.yaml
