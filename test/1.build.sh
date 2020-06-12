#!/bin/bash
docker_id="atyx300"
image_name="kchtest"

export GO111MODULE=on
go mod vendor

go build -o `pwd`/build/_output/bin/$image_name -mod=vendor `pwd`/cmd/main/main.go && \
docker build -t $docker_id/$image_name:v0.0.5 build && \
docker push $docker_id/$image_name:v0.0.5
