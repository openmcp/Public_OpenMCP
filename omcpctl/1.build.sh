#!/bin/bash

export GO111MODULE=on
go mod vendor

OPENMCP_APISERVER="10.0.3.20:31635"
OPENMCP_DIR="\/root\/workspace\/openmcp\/openmcp"
EXTERNAL_IP="10.0.3.12"

go build -o omcpctl && \
cp omcpctl /usr/local/bin && \
mkdir -p /var/lib/omcpctl && \
cp config.yaml /var/lib/omcpctl/config.yaml &&
sed -i 's/<YOUR_OPENMCP_APISERVER>/'${OPENMCP_APISERVER}'/g' /var/lib/omcpctl/config.yaml &&
sed -i 's/<YOUR_OPENMCP_DIR>/'${OPENMCP_DIR}'/g' /var/lib/omcpctl/config.yaml &&
sed -i 's/<YOUR_EXTERNAL_IP>/'${EXTERNAL_IP}'/g' /var/lib/omcpctl/config.yaml
