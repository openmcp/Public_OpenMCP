#!/bin/bash
docker_id="atyx300"
controller_name="openmcp-loadbalancing-controller"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcp/openmcp/openmcp-loadbalancing-controller/cmd/manager && \
docker build -t $docker_id/$controller_name:v0.0.1 build && \
docker push $docker_id/$controller_name:v0.0.1

