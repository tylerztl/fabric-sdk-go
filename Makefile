# Copyright zhigui Corp All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# -------------------------------------------------------------
# This makefile defines the following targets
#
#   - protos - generate all protobuf artifacts based on .proto files
#   - networkUp - start the fabric network
#   - networkDown - teardown the fabric network and clean the containers and intermediate images
#   - satrt - start the fabric-sdk-go server

.PHONY: protos
protos :
	./scripts/compile_protos.sh

.PHONY: networkUp
networkUp :
	./scripts/start_network.sh -m up

.PHONY: networkDown
networkDown :
	./scripts/start_network.sh -m down

.PHONY: start
start :
	go run main.go start
