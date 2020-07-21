#!/bin/bash

export GO111MODULE=on
go mod vendor

go build -o omcpctl && cp omcpctl /usr/local/bin

