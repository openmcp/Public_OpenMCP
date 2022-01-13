CTX_OPENMCP=openmcp

kubectl apply --context=$CTX_OPENMCP -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: istio-ingress-gateway
  namespace: default
spec:
  selector:
    istio: ingressgateway # use istio default controller
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPVirtualService
metadata:
  name: web
spec:
  hosts:
  - "keti.web.openmcp.in"
  - "keti.web.openmcp.out"
  gateways:
  - istio-ingress-gateway
  http:
  - match:
    - uri:
        prefix: /
    route:
    - destination:
        #host: web.default.svc.cluster.local
        host: web
        port:
          number: 8080
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPVirtualService
metadata:
  name: user-tool
spec:
  hosts:
  - "keti.user-tool.openmcp"
  gateways:
  - istio-ingress-gateway
  http:
  - match:
    - uri:
        prefix: /
    route:
    - destination:
        #host: user-tool.default.svc.cluster.local
        host: user-tool
        port:
          number: 8083

EOF

