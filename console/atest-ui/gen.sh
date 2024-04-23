#!/bin/bash

./node_modules/.bin/grpc_tools_node_protoc --plugin=protoc-gen-ts=. --ts_out=../../pkg/server -I ../../pkg/server ../../pkg/server/server.proto