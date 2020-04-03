#!/bin/bash

export GOOS=linux
export CGO_ENABLED=0
export GO111MODULE=on
export GOPROXY=https://goproxy.cn

go build -ldflags="-s" -o app

if [ $? -ne 0 ]; then
    echo "build fail."
    exit 1
fi
