
#kubectl label --context=cluster1 namespace default istio-injection=enabled --overwrite
#kubectl label --context=cluster2 namespace default istio-injection=enabled --overwrite

NS=bookinfo
CTX=cluster1
kubectl create secret generic regcred -n $NS --context $CTX  \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson

NS=bookinfo
CTX=cluster2
kubectl create secret generic regcred -n $NS --context $CTX  \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson


kubectl apply -f - <<EOF
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPService
metadata:
  name: productpage
  namespace: bookinfo
  labels:
    app: productpage
spec:
  labelselector:
    app: productpage
  template:
    spec:
      ports:
      - port: 9080
        name: http
      selector:
        app: productpage
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: productpage-v1
  namespace: bookinfo
spec:
  replicas: 2
  clusters:
  - cluster1
  - cluster2
  labels:
    app: productpage
    affinity: "yes"
  template:
    spec:
      template:
        spec:
          imagePullSecrets:
          - name: regcred
          containers:
          - name: productpage
            image: istio/examples-bookinfo-productpage-v1:1.10.0
            imagePullPolicy: IfNotPresent
            ports:
            - containerPort: 9080
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPService
metadata:
  name: details
  namespace: bookinfo
  labels:
    app: details
spec:
  labelselector:
    app: details
  template:
    spec:
      ports:
      - port: 9080
        name: http
      selector:
        app: details
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: details-v1
  namespace: bookinfo
spec:
  replicas: 1
  clusters:
  - cluster1
  labels:
    app: details
    affinity: "yes"
  template:
    spec:
      template:
        spec:
          imagePullSecrets:
          - name: regcred
          containers:
          - name: details
            image: istio/examples-bookinfo-details-v1:1.10.0
            imagePullPolicy: IfNotPresent
            ports:
            - containerPort: 9080
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPService
metadata:
  name: reviews
  namespace: bookinfo
  labels:
    app: reviews
spec:
  labelselector:
    app: reviews
  template:
    spec:
      ports:
      - port: 9080
        name: http
      selector:
        app: reviews
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: reviews-v1
  namespace: bookinfo
spec:
  replicas: 1
  clusters:
  - cluster1
  labels:
    app: reviews
    affinity: "yes"
  template:
    spec:
      template:
        spec:
          imagePullSecrets:
          - name: regcred
          containers:
          - name: reviews
            image: istio/examples-bookinfo-reviews-v1:1.10.0
            imagePullPolicy: IfNotPresent
            ports:
            - containerPort: 9080
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPService
metadata:
  name: ratings
  namespace: bookinfo
  labels:
    app: ratings
spec:
  labelselector:
    app: ratings
  template:
    spec:
      ports:
      - port: 9080
        name: http
      selector:
        app: ratings
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: ratings-v1
  namespace: bookinfo
spec:
  replicas: 1
  clusters:
  - cluster2
  labels:
    app: ratings
    affinity: "yes"
  template:
    spec:
      template:
        spec:
          imagePullSecrets:
          - name: regcred
          containers:
          - name: ratings
            image: istio/examples-bookinfo-ratings-v1:1.10.0
            imagePullPolicy: IfNotPresent
            ports:
            - containerPort: 9080
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: reviews-v2
  namespace: bookinfo
spec:
  replicas: 1
  clusters:
  - cluster2
  labels:
    app: reviews
    affinity: "yes"
  template:
    spec:
      template:
        spec:
          imagePullSecrets:
          - name: regcred
          containers:
          - name: reviews
            image: istio/examples-bookinfo-reviews-v2:1.10.0
            imagePullPolicy: IfNotPresent
            ports:
            - containerPort: 9080
---
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: reviews-v3
  namespace: bookinfo
spec:
  replicas: 1
  clusters:
  - cluster2
  labels:
    app: reviews
    affinity: "yes"
  template:
    spec:
      template:
        spec:
          imagePullSecrets:
          - name: regcred
          containers:
          - name: reviews
            image: istio/examples-bookinfo-reviews-v3:1.10.0
            imagePullPolicy: IfNotPresent
            ports:
            - containerPort: 9080
EOF
