OTEL_DOCKER_PROTOBUF ?= otel/build-protobuf:0.4.0
PROTOC := docker run --rm -u ${shell id -u} -v${PWD}:${PWD} -w${PWD} ${OTEL_DOCKER_PROTOBUF} --proto_path=${PWD}

.PHONY: genproto
genproto:
	$(PROTOC) --go_out=plugins=grpc:. internal/testing/testserver.proto
