#!/bin/bash
project_root_dir=/root/workspace/openmcp/openmcp-apiserver
docker_id="atyx300"
controller_name="openmcp-apiserver"

org_dir=`pwd`
cd $project_root_dir

export GO111MODULE=on
go mod vendor

go build -o `pwd`/build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcp-apiserver/pkg/main && \
docker build -t $docker_id/$controller_name:v0.0.1 `pwd`/build && \
docker push $docker_id/$controller_name:v0.0.1

cd $org_dir
