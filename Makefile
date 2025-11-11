PROJECT_NAME=order-service
GRPC_PORT=50051
PROTO_PATH=api
GEN_PATH=pkg/api/test
GOOGLEAPIS_PATH=third_party/googleapis

DOCKER_COMPOSE = docker-compose
GO = go

.PHONY: all build run clean gen test deps dev help
.PHONY: docker-up docker-down docker-logs docker-build docker-restart
.PHONY: migrate-up migrate-down db-shell redis-cli status nuke


all: build

gen:
	protoc -I. -I$(GOOGLEAPIS_PATH) \
		--go_out=. \
		--go-grpc_out=. \
		--grpc-gateway_out=. \
		$(PROTO_PATH)/order.proto


build:
	$(GO) build -o bin/server cmd/server/main.go
	$(GO) build -o bin/migrate cmd/migrate/main.go


run:
	$(GO) run cmd/server/main.go


migrate-local:
	$(GO) run cmd/migrate/main.go


clean:
	rm -rf bin/ $(GEN_PATH)/*.pb.go $(GEN_PATH)/*.gw.go


test:
	grpcurl -plaintext -d '{"item": "Test", "quantity": 1}' localhost:${GRPC_PORT} api.OrderService/CreateOrder


deps:
	$(GO) mod download
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	$(GO) install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest


dev: deps gen build


docker-build:
	$(DOCKER_COMPOSE) build

docker-up:
	$(DOCKER_COMPOSE) up -d --build

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f app

docker-restart: docker-down docker-up


migrate-up:
	$(DOCKER_COMPOSE) up migrate

migrate-down:
	$(DOCKER_COMPOSE) run --rm migrate ./migrate -path /migrations -database "postgres://postgres:password@postgres:5432/orders?sslmode=disable" down 1

db-shell:
	$(DOCKER_COMPOSE) exec postgres psql -U postgres -d orders

redis-cli:
	$(DOCKER_COMPOSE) exec redis redis-cli


status: ## Show service status
	$(DOCKER_COMPOSE) ps


nuke: docker-down
	docker system prune -f
	docker volume prune -f
	rm -rf bin/


help:
	@echo '$(PROJECT_NAME) Management Commands:'
	@echo ''
	@echo 'Targets:'
	@echo '  build         - Build binaries'
	@echo '  run           - Run server locally'
	@echo '  docker-up     - Start all services with Docker'
	@echo '  docker-logs   - Show application logs'
	@echo '  migrate-local - Run migrations locally'
	@echo '  test          - Test with grpcurl'