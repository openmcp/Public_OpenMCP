#!/bin/bash
docker_id="atyx300"
image_name="openmcp-scheduler"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$image_name -mod=vendor openmcp/openmcp/openmcp-scheduler/cmd/main && \
docker build -t $docker_id/$image_name:v0.0.5 build && \
docker push $docker_id/$image_name:v0.0.5

