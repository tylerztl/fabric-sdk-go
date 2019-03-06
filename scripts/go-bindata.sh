#!/bin/bash
#
# Copyright Zhigui Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -eux

# convert the Swagger UI to Go source code
go-bindata --nocompress -pkg swagger -o third_party/swagger-ui/datafile.go third_party/swagger-ui/...
