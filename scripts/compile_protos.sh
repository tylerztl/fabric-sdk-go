#!/bin/bash
#
# Copyright Zhigui Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -eux

# To set a proto root for a set of protos, create a .protoroot file in one of the parent directories
# which you wish to use as the proto root.  If no .protoroot file exists within ./<your_proto>
# then the proto root for that proto is inferred to be its containing directory.

# Find explicit proto roots
PROTO_ROOT_FILES="$(find . -name ".protoroot")"
PROTO_ROOT_DIRS="$(dirname $PROTO_ROOT_FILES)"

# As this is a proto root, and there may be subdirectories with protos, compile the protos for each sub-directory which contains them
for protos in $(find "$PROTO_ROOT_DIRS" -name '*.proto' -exec dirname {} \; | sort | uniq) ; do
    protoc --proto_path="$PROTO_ROOT_DIRS" \
            --go_out=plugins=grpc:$GOPATH/src \
            "$protos"/*.proto
done
