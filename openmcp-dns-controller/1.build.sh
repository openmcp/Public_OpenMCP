#!/bin/bash
docker_id="atyx300"
controller_name="openmcp-dns-controller"

export GO111MODULE=on
go mod vendor

go build -o `pwd`/build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcp-dns-controller/cmd/manager && \
docker build -t $docker_id/$controller_name:v0.0.1 `pwd`/build && \
docker push $docker_id/$controller_name:v0.0.1

