-include .make.env

REGISTRY ?=
APP_NAME ?=
TAG ?=

IMAGE_NAME := $(REGISTRY)/$(APP_NAME)

PROTO_SRC := proto
PROTO_VENDOR := third_party
DOCKER_CONTEXT := .

.PHONY: all build-app push-app generate-proto clean-proto check-tag

all: build-app

check-tag:
ifndef TAG
	$(error TAG is required. Usage: make TAG=v1.2.3)
endif

build-app: check-tag
	docker build -t $(IMAGE_NAME):$(TAG) $(DOCKER_CONTEXT)

push-app: build-app
	docker push $(IMAGE_NAME):$(TAG)

generate-proto:
	protoc \
		-I $(PROTO_SRC) \
		-I $(PROTO_VENDOR) \
		--go_out=paths=source_relative:$(PROTO_SRC) \
		--go-grpc_out=paths=source_relative:$(PROTO_SRC) \
		--validate_out=lang=go,paths=source_relative:$(PROTO_SRC) \
		$$(find $(PROTO_SRC) -name "*.proto")

clean-proto:
	find $(PROTO_SRC) -name "*.pb.go" -delete
	find $(PROTO_SRC) -name "*.pb.validate.go" -delete

