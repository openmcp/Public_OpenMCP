#!/bin/bash
docker_id="openmcp"
image_name="openmcp-analytic-engine"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$image_name -mod=vendor openmcp/openmcp/openmcp-analytic-engine/cmd/main && \
docker build -t $docker_id/$image_name:v0.0.1 build && \
docker push $docker_id/$image_name:v0.0.1
