kubectl create ns openmcp
kubectl create ns metallb-system

echo "--- openmcp-analytic-engine"
kubectl create -f openmcp-analytic-engine/.
echo "--- openmcp-metric-collector"
kubectl create -f openmcp-metric-collector/.
echo "--- influxdb"
kubectl create -f influxdb/.
cd influxdb/secret_info
sh secret_info.sh
cd ../..
echo "--- openmcp-deployment-controller"
kubectl create -f openmcp-deployment-controller/.
echo "--- openmcp-has-controller"
kubectl create -f openmcp-has-controller/.
echo "--- openmcp-scheduler"
kubectl create -f openmcp-scheduler/.
echo "--- openmcp-ingress-controller"
kubectl create -f openmcp-ingress-controller/.
echo "--- openmcp-service-controller"
kubectl create -f openmcp-service-controller/.
echo "--- openmcp-policy-engine"
kubectl create -f openmcp-policy-engine/.
echo "   ==> CREATE Policy"
echo "--- create policy"
kubectl create -f openmcp-policy-engine/policy/.
echo "--- openmcp-dns-controller"
kubectl create -f openmcp-dns-controller/.
echo "--- openmcp-loadbalancing-controller"
kubectl create -f openmcp-loadbalancing-controller/.
kubectl expose deployment/openmcp-loadbalancing-controller -n openmcp --port 80 --type=LoadBalancer
echo "--- openmcp-sync-controller"
kubectl create -f openmcp-sync-controller/.
echo "--- metallb"
kubectl create -f metallb/.
