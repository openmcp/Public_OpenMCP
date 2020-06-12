export REGISTRY="atyx300"
make test-adapter-container

docker push $REGISTRY/k8s-test-metrics-adapter-amd64:latest

#kubectl apply -f test-adapter-deploy/testing-adapter.yaml --context cluster1
#kubectl apply -f test-adapter-deploy/testing-adapter.yaml --context cluster2
kubectl apply -f test-adapter-deploy/testing-adapter.yaml --context cluster3
#kubectl apply -f test-adapter-deploy/testing-adapter.yaml --context cluster4
#kubectl apply -f test-adapter-deploy/testing-adapter.yaml --context cluster5
#kubectl apply -f test-adapter-deploy/testing-adapter.yaml --context cluster6
