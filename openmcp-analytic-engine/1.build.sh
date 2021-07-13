#!/bin/bash
docker_id="ketidevit"
controller_name="openmcp-analytic-engine"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$controller_name -mod=vendor openmcp/openmcp/openmcp-analytic-engine/src/main && \

docker build -t $docker_id/$controller_name:v0.0.1 build && \
docker push $docker_id/$controller_name:v0.0.1
