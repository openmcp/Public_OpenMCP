kubectl create ns openmcp

echo "--- analytic-engine"
kubectl create -f analytic-engine/.
echo "--- metric-collector"
kubectl create -f metric-collector/.
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
echo "--- loadbalancing-controller"
kubectl create -f loadbalancing-controller/.
kubectl expose deployment/loadbalancing-controller -n openmcp --port 80 --type=LoadBalancer
echo "--- sync-controller"
kubectl create -f sync-controller/.
echo "--- metallb"
kubectl create -f metallb/.