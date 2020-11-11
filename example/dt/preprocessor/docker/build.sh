#!/bin/bash
docker_id="openmcp"
image_name="keti-preprocessor"

docker build -t $docker_id/$image_name:v0.0.1 .
docker push $docker_id/$image_name:v0.0.1
