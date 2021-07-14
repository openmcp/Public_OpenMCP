kubectl create ns openmcp
kubectl create ns metallb-system
kubectl create ns istio-system

kubectl create secret generic REPLACE_DOCKERSECRETNAME \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson \
    --namespace=openmcp


echo "--- deploy crds"
kubectl create -f crds/.
echo "--- openmcp-cluster-manager"
kubectl create -f openmcp-cluster-manager/.
echo "--- openmcp-analytic-engine"
kubectl create -f openmcp-analytic-engine/.
echo "--- openmcp-apiserver"
kubectl create -f openmcp-apiserver/.
echo "--- openmcp-configmap-controller"
kubectl create -f openmcp-configmap-controller/.
echo "--- openmcp-secret-controller"
kubectl create -f openmcp-secret-controller/.
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
echo "--- configmap"
kubectl apply -f configmap/coredns/.
echo "--- ingress gateway"
kubectl create ns istio-system --context openmcp
# istio 클러스터간 접근을 위한 인증서 만들기
cd istio
export PATH=$PWD/bin:$PATH

mkdir -p certs
pushd certs

make -f ../tools/certs/Makefile.selfsigned.mk root-ca
make -f ../tools/certs/Makefile.selfsigned.mk openmcp-cacerts

kubectl create secret generic cacerts -n istio-system \
      --from-file=openmcp/ca-cert.pem \
      --from-file=openmcp/ca-key.pem \
      --from-file=openmcp/root-cert.pem \
      --from-file=openmcp/cert-chain.pem
popd



# istio-system 네임 스페이스가 이미 생성 된 경우 여기에 클러스터의 네트워크를 설정해야합니다
kubectl --context=openmcp get namespace istio-system && \
kubectl --context=openmcp label namespace istio-system topology.istio.io/network=network-openmcp

# openmcp에 대한 Istio configuration 을 만듭니다.
cat <<EOF > openmcp.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  meshConfig:
   defaultConfig:
     proxyMetadata:
       ISTIO_META_DNS_CAPTURE: "true"
  values:
    global:
      meshID: mesh-openmcp
      multiCluster:
        clusterName: openmcp
      network: network-openmcp
EOF

# openmcp에 configuration 적용
istioctl install --context=openmcp -f openmcp.yaml -y

# openmcp에 east-west traffic 전용 게이트웨이를 설치합니다.
samples/multicluster/gen-eastwest-gateway.sh \
    --mesh mesh-openmcp --cluster openmcp --network network-openmcp | \
    istioctl --context=openmcp install -y -f -

# East-west 게이트웨이에 외부 IP 주소가 할당 될 때까지 기다립니다.
for ((;;))
do
        status=`kubectl --context=openmcp get svc istio-eastwestgateway -n istio-system | grep istio-eastwestgateway | awk '{print $4}'`
        if [ "$status" != "<none>" ]; then
                break
        fi
        echo "Wait LB IP Allocate"
        sleep 1
done

# Expose the control plane in openmcp
kubectl apply --context=openmcp -f \
    samples/multicluster/expose-istiod.yaml

# Expose services in openmcp
kubectl --context=openmcp apply -n istio-system -f \
    samples/multicluster/expose-services.yaml

#istio 인증서 복사
cp -r certs ../member/istio/
