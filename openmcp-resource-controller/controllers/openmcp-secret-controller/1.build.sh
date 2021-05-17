#!/bin/bash
docker_id="ketidevit"
controller_name="openmcp-secret-controller"

export GO111MODULE=on
go mod vendor

go build -o `pwd`/build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcp/openmcp/openmcp-resource-controller/controllers/openmcp-secret-controller/pkg/main && \

docker build -t $docker_id/$controller_name:v0.0.1 build && \
docker push $docker_id/$controller_name:v0.0.1

