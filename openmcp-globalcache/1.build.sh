#!/bin/bash
docker_id="openmcp"
image_name="openmcp-globalcache"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$image_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcp/openmcp/openmcp-globalcache/cmd/manager && \
docker build -t $docker_id/$image_name:v0.0.2 build && \
docker push $docker_id/$image_name:v0.0.2
