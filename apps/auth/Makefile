-include .make.env

# ----------------------------
# Application image settings
# ----------------------------
REGISTRY ?=
APP_NAME ?=
TAG ?=

MIGRATE_NAME ?=
MIGRATE_DEV_URL ?=

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

MOCK_DST_DIR := test/mock
MOCK_TARGET_DIR ?=

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
.PHONY: generate-migrate apply-migrate build-migrate push-migrate
generate-migrate: check-migrate-name
	atlas migrate diff $(MIGRATE_NAME) \
		--dev-url $(MIGRATE_DEV_URL) \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema"

apply-migrate:
	atlas migrate apply \
		--url $(MIGRATE_DEV_URL) \
		--dir "file://ent/migrate/migrations"

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

.PHONY: check-migrate-name
check-migrate-name:
ifndef MIGRATE_NAME
	$(error MIGRATE_NAME is required. Usage: make MIGRATE_NAME=migrate)
endif

# ----------------------------
#  Generate mocks for interfaces
#  ----------------------------
.PHONY: generate-mocks check-mock-target-dir
generate-mocks: check-mock-target-dir
	@echo "🔍 Generating mocks for interfaces under: '$(MOCK_TARGET_DIR)'"
	@MODULE_PATH=$$(go list -m); \
	find $(MOCK_TARGET_DIR) -type f -name '*.go' -not -path '*/mocks/*' | while read src; do \
		src_rel=$$(realpath --relative-to=. $$src); \
		pkg_dir=$$(dirname $$src_rel); \
		ifaces=$$(grep -E 'type [A-Z][a-zA-Z0-9_]+ interface' $$src | awk '{print $$2}'); \
		import_path=$${MODULE_PATH}/$$pkg_dir; \
		for iface in $$ifaces; do \
			out_dir=$(MOCK_DST_ROOT)/$$pkg_dir; \
			mkdir -p $$out_dir; \
			out_file=$$out_dir/mock_$$iface.go; \
			echo "⚙️  Generating $$iface → $$out_file"; \
			mockgen $$import_path $$iface > $$out_file || echo "❌ Failed to generate mock for $$iface"; \
		done; \
	done

check-mock-target-dir:
ifndef MOCK_TARGET_DIR
	$(error MOCK_TARGET_DIR is required. Usage: make MOCK_TARGET_DIR=path/to/target)
endif
