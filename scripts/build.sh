#!/bin/bash

GOOS=darwin GOARCH=amd64 go build -v -o $GOPATH/bin/terraform-provider-kubernetes_darwin-amd64
GOOS=linux GOARCH=amd64 go build -v -o $GOPATH/bin/terraform-provider-kubernetes_linux-amd64
GOOS=windows GOARCH=amd64 go build -v -o $GOPATH/bin/terraform-provider-kubernetes_windows-amd64

