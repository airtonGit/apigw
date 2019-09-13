#!/bin/bash
docker run --rm -it -v "$PWD":/go/src/api-gateway -w /go/src/api-gateway -e GOOS=linux -e GOARCH=386 golang:1.12.4 go build -v