#!/bin/bash

NS=$1
CLUSTERNAME=$2

kubectl get namespace $NS -o json --context $CLUSTERNAME > temp.json

sed -i -e 's/"kubernetes"//' temp.json

kubectl replace --raw "/api/v1/namespaces/$NS/finalize" -f ./temp.json --context $CLUSTERNAME
