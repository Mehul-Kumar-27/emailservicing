#!/bin/bash

echo "Genetating proto files..."

protoc --go_out=./../logs --go_opt=paths=source_relative \
    --go-grpc_out=./../logs --go-grpc_opt=paths=source_relative \
    logs.proto
