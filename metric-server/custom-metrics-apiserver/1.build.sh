export REGISTRY="openmcp"

make test-adapter-container



docker push $REGISTRY/k8s-test-metrics-adapter-amd64:latest
