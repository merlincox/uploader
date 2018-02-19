#!/bin/bash

source_dir=$(dirname $0)

export GOARCH=amd64
export GOOS=linux

go build -o $source_dir/s3uploader.ubuntu $source_dir/main.go
