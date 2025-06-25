-include .make.env

# ----------------------------
# Application image settings
# ----------------------------
REGISTRY ?=
APP_NAME ?=
TAG ?=

APP_IMAGE := $(REGISTRY)/$(APP_NAME)

# ----------------------------
# Migration image settings
# ----------------------------
MIGRATE_REGISTRY ?= $(REGISTRY)
MIGRATE_NAME ?= $(APP_NAME)-migrate
MIGRATE_TAG ?= $(TAG)

MIGRATE_IMAGE := $(MIGRATE_REGISTRY)/$(MIGRATE_NAME)

# ----------------------------
# Directory paths
# ----------------------------
PROTO_SRC := proto
DOCKER_CONTEXT := ./docker

# ----------------------------
# Default target
# ----------------------------
.PHONY: all
all: build-app

# ----------------------------
#  Run the application
#  ----------------------------
.PHONY: run
run:
	go run ./cmd/server

# ----------------------------
# Build & push app image
# ----------------------------
.PHONY: build-app push-app
build-app: check-tag
	docker build -t $(APP_IMAGE):$(TAG) -f $(DOCKER_CONTEXT)/app.Dockerfile .

push-app: build-app
	docker push $(APP_IMAGE):$(TAG)

# ----------------------------
# Build & push migration image
# ----------------------------
.PHONY: build-migrate push-migrate
build-migrate: check-tag
	docker build -t $(MIGRATE_IMAGE):$(MIGRATE_TAG) -f $(DOCKER_CONTEXT)/migrate.Dockerfile .

push-migrate: build-migrate
	docker push $(MIGRATE_IMAGE):$(MIGRATE_TAG)

# ----------------------------
# Generate and clean proto files
# ----------------------------
.PHONY: generate-proto clean-proto
generate-proto:
	protoc \
		-I $(PROTO_SRC) \
		--go_out=paths=source_relative:$(PROTO_SRC) \
		--go-grpc_out=paths=source_relative:$(PROTO_SRC) \
		--validate_out=lang=go,paths=source_relative:$(PROTO_SRC) \
		$$(find $(PROTO_SRC) -name "*.proto")

clean-proto:
	find $(PROTO_SRC) -name "*.pb.go" -delete
	find $(PROTO_SRC) -name "*.pb.validate.go" -delete

# ----------------------------
# Ent schema generation
# ----------------------------
.PHONY: generate-ent
generate-ent:
	@echo "Generating Ent schema..."
	go generate ./ent

# ----------------------------
# Run unit tests
# ----------------------------
.PHONY: test
test:
	go test ./test/...

# ----------------------------
# Check that the image tag is set
# ----------------------------
.PHONY: check-tag
check-tag:
ifndef TAG
	$(error TAG is required. Usage: make TAG=v1.2.3)
endif
