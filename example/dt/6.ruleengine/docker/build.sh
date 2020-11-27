#!/bin/bash
docker_id="openmcp"
image_name="keti-ruleengine"

#kubectl delete -f ../yaml

docker build -t $docker_id/$image_name:v0.0.1 . && \
docker push $docker_id/$image_name:v0.0.1

#kubectl create -f ../yaml
