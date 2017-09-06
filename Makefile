SHELL := /bin/bash

DOCKERFILE := Dockerfile
DOCKER_IMAGE := nawa/backfriend
DOCKERFILE_DIR := docker
DOCKER_CONTAINER := backfriend
PORTS := -p 8888:8080
ENV ?= dev
VOLUMES := -v $(PWD)/config/env/$(ENV).yml:/app/config/config.yml

.PHONY: docker-build docker-run docker-stop docker-start docker-rmf docker-rmi functional-tests

docker-build:
	docker rmi -f $(DOCKER_IMAGE):bak || true
	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE):bak || true
	docker rmi -f $(DOCKER_IMAGE) || true
	docker build -f $(DOCKERFILE_DIR)/$(DOCKERFILE) -t $(DOCKER_IMAGE) .

docker-run:
	docker rm $(DOCKER_CONTAINER) || true
	docker run -d --name $(DOCKER_CONTAINER) $(PORTS) $(VOLUMES) $(DOCKER_IMAGE)

docker-stop:
	docker stop $(DOCKER_CONTAINER)

docker-start:
	docker start $(DOCKER_CONTAINER)

docker-rmf:
	docker rm -f $(DOCKER_CONTAINER)

docker-rmi:
	docker rmi $(DOCKER_IMAGE)

functional-tests:
	@/bin/bash functional-tests/run_tests.sh
