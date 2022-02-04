mkdir -p /home/nfs/pv/influxdb
mkdir -p /home/nfs/pv/api-server/cert

echo "--- API Server Generation Key File"
#openssl genrsa -out server.key 2048
#echo \r\n ; echo \r\n; echo \r\n; echo \r\n; echo \r\n; echo openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org; echo \r\n) | openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
openssl req \
    -x509 \
    -nodes \
    -newkey rsa:2048 \
    -keyout server.key \
    -out server.crt \
    -days 3650 \
    -subj "/C=KR/ST=Seoul/L=Seoul/O=Global Company/OU=IT Department/CN=openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"

mv server.key /home/nfs/pv/api-server/cert
mv server.crt /home/nfs/pv/api-server/cert

kubectl create ns openmcp
kubectl create ns metallb-system
kubectl create ns istio-system
# kubectl create ns nginx-ingress

echo "Input Your Docker ID(No Pull Limit Plan)"
docker login

kubectl create secret generic "regcred" \
    --from-file=.dockerconfigjson=$HOME/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson \
    --namespace=openmcp

echo "--- deploy crds"
kubectl create -f ../../crds/.
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
echo "--- openmcp-sync-controller"
kubectl create -f openmcp-sync-controller/.
echo "--- openmcp-job-controller"
kubectl apply -f openmcp-job-controller/.
echo "--- openmcp-namespace-controller"
kubectl apply -f openmcp-namespace-controller/.
echo "--- openmcp-pv-controller"
kubectl apply -f openmcp-pv-controller/.
echo "--- openmcp-pvc-controller"
kubectl apply -f openmcp-pvc-controller/.
echo "--- openmcp-daemonset-controller"
kubectl apply -f openmcp-daemonset-controller/.
echo "--- openmcp-statefulset-controller"
kubectl apply -f openmcp-statefulset-controller/.
echo "--- metallb"
kubectl create -f metallb/.
echo "--- configmap"
kubectl apply -f configmap/coredns/.
# echo "--- ingress gateway"
# kubectl create -f nginx-ingress-controller

kubectl create ns istio-system --context openmcp
# istio 클러스터간 접근을 위한 인증서 만들기
cd istio
export PATH=$PWD/bin:$PATH

mkdir -p certs
pushd certs

CTX=openmcp

make -f ../tools/certs/Makefile.selfsigned.mk root-ca
make -f ../tools/certs/Makefile.selfsigned.mk openmcp-cacerts

kubectl create secret generic cacerts -n istio-system \
      --from-file=openmcp/ca-cert.pem \
      --from-file=openmcp/ca-key.pem \
      --from-file=openmcp/root-cert.pem \
      --from-file=openmcp/cert-chain.pem
popd

chmod 755 bin/istioctl
sudo cp bin/istioctl /usr/local/bin
chmod 755 samples/multicluster/gen-eastwest-gateway.sh

# istio-system 네임 스페이스가 이미 생성 된 경우 여기에 클러스터의 네트워크를 설정해야합니다
kubectl --context=$CTX get namespace istio-system && \
kubectl --context=$CTX label namespace istio-system topology.istio.io/network=network-$CTX

# openmcp에 대한 Istio configuration 을 만듭니다.
cat <<EOF > $CTX.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  meshConfig:
   defaultConfig:
     proxyMetadata:
       ISTIO_META_DNS_CAPTURE: "true"
  values:
    global:
      meshID: mesh-$CTX
      multiCluster:
        clusterName: $CTX
      network: network-$CTX
  components:
    ingressGateways:
      - name: istio-ingressgateway
        enabled: true
        k8s:
          service:
            ports:
              # We have to specify original ports otherwise it will be erased
              - name: status-port
                nodePort: 31022
                port: 15022
                protocol: TCP
                targetPort: 15021
              - name: http2
                nodePort: 31080
                port: 80
                protocol: TCP
                targetPort: 8080
              - name: https
                nodePort: 31443
                port: 443
                protocol: TCP
                targetPort: 8443
              - name: tcp-istiod
                nodePort: 31013
                port: 15013
                protocol: TCP
                targetPort: 15012
              - name: tls
                nodePort: 31444
                port: 15444
                protocol: TCP
                targetPort: 15443
EOF

# openmcp에 configuration 적용
istioctl install --context=$CTX -f $CTX.yaml -y

# openmcp에 east-west traffic 전용 게이트웨이를 설치합니다.
samples/multicluster/gen-eastwest-gateway.sh \
    --mesh mesh-$CTX --cluster $CTX --network network-$CTX | \
    istioctl --context=$CTX install -y -f -

# East-west 게이트웨이에 외부 IP 주소가 할당 될 때까지 기다립니다.
for ((;;))
do
        status=`kubectl --context=$CTX get svc istio-eastwestgateway -n istio-system | grep istio-eastwestgateway | awk '{print $4}'`
        if [ "$status" != "<none>" ]; then
                break
        fi
        echo "Wait LB IP Allocate"
        sleep 1
done

# Expose the control plane in openmcp
kubectl apply --context=$CTX -f \
    samples/multicluster/expose-istiod.yaml

# Expose services in openmcp
kubectl --context=$CTX apply -n istio-system -f \
    samples/multicluster/expose-services.yaml

kubectl apply -f patch_istio_configmap.yaml

#istio 인증서 복사
rm -r ../../member/istio/certs
cp -r certs ../../member/istio/

# Core DNS 리스타트
kubectl delete pod --namespace kube-system --selector k8s-app=kube-dns
