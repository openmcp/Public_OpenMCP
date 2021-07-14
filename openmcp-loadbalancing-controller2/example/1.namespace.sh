kubectl apply -f - <<EOF
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPNamespace
metadata:
  name: bookinfo
EOF
