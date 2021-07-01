#!/bin/bash
docker_registry_ip="10.0.3.40:5005"
docker_id="openmcp"
controller_name="openmcp-migration"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcp/openmcp/openmcp-migration/cmd/manager && \

docker build -t $docker_registry_ip/$docker_id/$controller_name:v0.0.6 build && \
docker push $docker_registry_ip/$docker_id/$controller_name:v0.0.6\
