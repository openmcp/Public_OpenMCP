#!/bin/bash
controller_name="etcd-http-server"

export GO111MODULE=on
go mod vendor

go build -o `pwd`/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor $controller_name/pkg/main
