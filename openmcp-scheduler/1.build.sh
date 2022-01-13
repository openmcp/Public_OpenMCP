#!/bin/bash
docker_id="ketidevit2"
image_name="openmcp-scheduler"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$image_name -mod=vendor openmcp/openmcp/$image_name/src/main && \
docker build -t $docker_id/$image_name:v0.0.1 build && \
docker push $docker_id/$image_name:v0.0.1

