#!/bin/bash
#
# Copyright Zhigui Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -eux

# Protoc Plugin
brew install protobuf

go get -u github.com/golang/protobuf/protoc-gen-go
