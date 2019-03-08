# Copyright zhigui Corp All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# -------------------------------------------------------------
# This makefile defines the following targets
#
#   - protos - generate all protobuf artifacts based on .proto files

.PHONY: protos
protos :
	./scripts/compile_protos.sh
