#!/bin/bash

echo "Genetating proto files..."

protoc --go_out=./../api/logs --go_opt=paths=source_relative \
    --go-grpc_out=./../api/logs --go-grpc_opt=paths=source_relative \
    logs.proto
