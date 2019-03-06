#!/bin/bash
#
# Copyright Zhigui Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -eux

# Protoc Plugin
go get -u github.com/golang/protobuf/protoc-gen-go

# grpc-gateway is a plugin of protoc. It reads gRPC service definition,
# and generates a reverse-proxy server which translates a RESTful JSON API into gRPC.
# This server is generated according to custom options in your gRPC definition.
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

# Swagger Plugins
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/jteeuwen/go-bindata/...
go get -u github.com/elazarl/go-bindata-assetfs/...