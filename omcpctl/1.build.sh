#!/bin/bash

export GO111MODULE=on
go mod vendor

OPENMCP_DIR="\/root\/workspace\/usr\/kch\/openmcp\/openmcp"
EXTERNAL_IP="10.0.3.12"

go build -o omcpctl && \
cp omcpctl /usr/local/bin/omcpctl && \
mkdir -p /var/lib/omcpctl && \
cp config.yaml /var/lib/omcpctl/config.yaml &&
sed -i 's/<YOUR_OPENMCP_DIR>/'${OPENMCP_DIR}'/g' /var/lib/omcpctl/config.yaml &&
sed -i 's/<YOUR_EXTERNAL_IP>/'${EXTERNAL_IP}'/g' /var/lib/omcpctl/config.yaml
