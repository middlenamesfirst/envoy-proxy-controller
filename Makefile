DOCKER_ORG := middlenamesfirst
NAME := envoy-proxy-controller
GIT_COMMIT := $(shell git rev-parse --short=10 HEAD 2>/dev/null)

BASE_IMAGE_URL := $(DOCKER_ORG)/$(NAME)
IMAGE_URL := $(BASE_IMAGE_URL):$(GIT_COMMIT)

.PHONY: docker-build
docker-build:
	docker build --pull -t ${IMAGE_URL} .

.PHONY: vendor
version:
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod vendor

start-local:
	docker-compose up --build

stop-local:
	docker-compose down
