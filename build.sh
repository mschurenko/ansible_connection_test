#!/bin/bash

set -ex

# GOOS=linux GOARCH=arm go build -o library/check_port check_port.go

pushd .
cd check_port
GOOS=darwin GOARCH=amd64 go build -o ../library/check_port_osx -race main.go
GOOS=linux GOARCH=amd64 go build -o ../library/check_port main.go

popd
cd check_http
GOOS=darwin GOARCH=amd64 go build -o ../library/check_http_osx -race main.go
GOOS=linux GOARCH=amd64 go build -o ../library/check_http main.go
