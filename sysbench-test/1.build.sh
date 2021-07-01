#!/bin/bash
docker_registry_ip="10.0.3.20:5005"
docker_id="openmcp"
controller_name="sysbench-test"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcp/openmcp/sysbench-test/cmd/manager && \

docker build -t $docker_registry_ip/$docker_id/$controller_name:v0.0.1 build && \
docker push $docker_registry_ip/$docker_id/$controller_name:v0.0.1


#docker build -t $docker_id/$controller_name:v0.0.1 build && \
#docker push $docker_id/$controller_name:v0.0.1
