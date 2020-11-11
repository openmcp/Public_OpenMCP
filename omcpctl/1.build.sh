#!/bin/bash

export GO111MODULE=on
go mod vendor

OPENMCP_APISERVER="10.0.0.226:31635"
OPENMCP_DIR="\/root\/workspace\/public-openmcp\/openmcp"
EXTERNAL_IP="211.45.109.210"

go build -o omcpctl && \
cp omcpctl /usr/local/bin/omcpctl && \
mkdir -p /var/lib/omcpctl && \
cp config.yaml /var/lib/omcpctl/config.yaml &&
sed -i 's/<YOUR_OPENMCP_APISERVER>/'${OPENMCP_APISERVER}'/g' /var/lib/omcpctl/config.yaml &&
sed -i 's/<YOUR_OPENMCP_DIR>/'${OPENMCP_DIR}'/g' /var/lib/omcpctl/config.yaml &&
sed -i 's/<YOUR_EXTERNAL_IP>/'${EXTERNAL_IP}'/g' /var/lib/omcpctl/config.yaml
