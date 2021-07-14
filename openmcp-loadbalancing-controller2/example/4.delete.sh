kubectl delete openmcpdeployment -n bookinfo --all
kubectl delete openmcpservice -n bookinfo --all
kubectl delete gw,dr,se,ovs,vs -n bookinfo --all
kubectl delete ons bookinfo
