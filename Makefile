# Variables
PROJECT_NAME=order-service
GRPC_PORT=50051
PROTO_PATH=api
GEN_PATH=pkg/api/test
GOOGLEAPIS_PATH=third_party/googleapis

.PHONY: all build run clean gen test

# Default target
all: build

# Generate gRPC code
gen:
	protoc -I. -I$(GOOGLEAPIS_PATH) \
    		--go_out=. \
    		--go-grpc_out=. \
    		--grpc-gateway_out=. \
    		$(PROTO_PATH)/order.proto

# Build binary
build:
	go build -o bin/server cmd/server/main.go

# Run server
run:
	go run cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin/ $(GEN_PATH)/*.pb.go $(GEN_PATH)/*.gw.go

# Test with grpcurl
test:
	grpcurl -plaintext -d '{"item": "Test", "quantity": 1}' localhost:${GRPC_PORT} api.OrderService/CreateOrder

# Install dependencies
deps:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Development setup
dev: deps gen build

help:
	@echo "Available targets:"
	@echo "  gen    - Generate gRPC code"
	@echo "  build  - Build binary"
	@echo "  run    - Run server"
	@echo "  test   - Test with grpcurl"
	@echo "  clean  - Clean build artifacts"
	@echo "  dev    - Full development setup"