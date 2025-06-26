-include .make.env

# ──────────────────────────────
# 📁 Directory Paths
# ──────────────────────────────
PROTO_SRC := proto

# ──────────────────────────────
# 📜 Generate & Clean Protobuf
# ──────────────────────────────
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

# ──────────────────────────────
# 🧪 Unit Tests
# ──────────────────────────────
test-unit:
	go test .../test/unit/...

