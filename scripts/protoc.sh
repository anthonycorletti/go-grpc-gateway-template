#!/bin/sh -ex

protoc -I ./proto \
--go_out ./proto \
--go_opt paths=source_relative \
--go-grpc_out ./proto \
--go-grpc_opt paths=source_relative \
--grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative --grpc-gateway_opt logtostderr=true \
--openapiv2_out ./docs --openapiv2_opt logtostderr=true --openapiv2_opt use_go_templates=true \
./proto/api/api.proto
