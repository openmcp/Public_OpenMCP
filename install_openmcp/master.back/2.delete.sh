echo "--- openmcp-cluster-manager"
kubectl delete -f openmcp-cluster-manager/.
echo "--- openmcp-analytic-engine"
kubectl delete -f openmcp-analytic-engine/.
echo "--- openmcp-apiserver"
kubectl delete -f openmcp-apiserver/.
echo "--- openmcp-configmap-controller"
kubectl delete -f openmcp-configmap-controller/.
echo "--- openmcp-secret-controller"
kubectl delete -f openmcp-secret-controller/.
echo "--- openmcp-metric-collector"
kubectl delete -f openmcp-metric-collector/.
echo "--- influxdb"
kubectl delete -f influxdb/.
cd influxdb/secret_info
sh secret_info_delete.sh
cd ../..
echo "--- openmcp-deployment-controller"
kubectl delete -f openmcp-deployment-controller/.
echo "--- openmcp-has-controller"
kubectl delete -f openmcp-has-controller/.
echo "--- openmcp-scheduler"
kubectl delete -f openmcp-scheduler/.
echo "--- openmcp-ingress-controller"
kubectl delete -f openmcp-ingress-controller/.
echo "--- openmcp-service-controller"
kubectl delete -f openmcp-service-controller/.
echo "--- openmcp-policy-engine"
kubectl delete -f openmcp-policy-engine/.
echo "   ==> Delete Policy"
echo "--- delete policy"
#kubectl delete -f openmcp-policy-engine/policy/.
echo "--- openmcp-dns-controller"
kubectl delete -f openmcp-dns-controller/.
echo "--- loadbalancing-controller"
kubectl delete -f openmcp-loadbalancing-controller/.
echo "--- sync-controller"
kubectl delete -f openmcp-sync-controller/.
echo "--- metallb"
kubectl delete -f metallb/.
echo "--- delete crds"
kubectl delete -f crds/.

kubectl delete ns openmcp
