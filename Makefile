# const
PROJECT_NAME := gospy

# paths
ENV_PATH := .env.dev
DOCKER_COMPOSE_PATH = ./build/compose.yaml

# exec
GO := go
DOCKER_COMPOSE := docker compose -f $(DOCKER_COMPOSE_PATH) --env-file $(ENV_PATH) -p $(PROJECT_NAME)

.PHONY: build
build: clean build/bin/gospy

.PHONY: build/bin/gospy
build/bin/gospy:
	$(GO) build -o $(@) ./cmd/gospy/main.go

.PHONY: clean
clean:
	rm -rf build/bin/*

.PHONY: tidyvendor
tidyvendor:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: docker-build
docker-build:
	$(DOCKER_COMPOSE) build

.PHONY: docker-up
docker-up:
	$(DOCKER_COMPOSE) up -d

.PHONY: docker-stop
docker-stop:
	$(DOCKER_COMPOSE) stop

.PHONY: docker-down
docker-down:
	$(DOCKER_COMPOSE) down -v
