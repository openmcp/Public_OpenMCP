#!/bin/bash
docker_id="openmcp"
image_name="openmcp-snapshot"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$image_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcp/openmcp/openmcp-snapshot/cmd/manager && \
docker build -t $docker_id/$image_name:v1.1.0.dev build && \
docker push $docker_id/$image_name:v1.1.0.dev\
