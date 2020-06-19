#!/bin/bash
project_root_dir=/root/workspace/usr/lyn/openmcpscheduler
docker_id="atyx300"
resource_name="openmcpscheduler"
controller_name="openmcp-scheduler"

org_dir=`pwd`
cd $project_root_dir

export GO111MODULE=on
go mod vendor

go build -o `pwd`/build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcpscheduler/cmd && \
docker build -t $docker_id/$controller_name:v0.0.5 build && \
docker push $docker_id/$controller_name:v0.0.5

cd $org_dir
