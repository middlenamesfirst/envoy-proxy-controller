DOCKER_ORG := middlenamesfirst
NAME := docker-envoy-proxy
GIT_COMMIT := $(shell git rev-parse --short=10 HEAD 2>/dev/null)

BASE_IMAGE_URL := $(DOCKER_ORG)/$(NAME)
IMAGE_URL := $(BASE_IMAGE_URL):$(GIT_COMMIT)

.PHONY: docker-build
docker-build:
	docker build --pull -t ${IMAGE_URL} .

.PHONY: docker-run
docker-run:
	docker run -d --name envoy -p 9901:9901 ${IMAGE_URL}
