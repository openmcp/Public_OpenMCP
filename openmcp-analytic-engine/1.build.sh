#!/bin/bash
docker_id="atyx300"
image_name="openmcp-analytic-engine"

export GO111MODULE=on
go mod vendor
kubectl config view >> `pwd`/build/bin/config

go build -o `pwd`/build/_output/bin/$image_name -mod=vendor `pwd`/cmd/main/main.go && \
docker build -t $docker_id/$image_name:v0.0.5 build && \
docker push $docker_id/$image_name:v0.0.5
