include ./linting.mk

IMAGE_TAG := $(shell git rev-parse HEAD)
IMAGE_NAME="ff-bot"

.PHONY: deps
deps:
	go mod tidy
	go mod download
	go mod vendor
	go mod verify

.PHONY: dockerize
dockerize:
	docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" -f Dockerfile .

.PHONY: build
build:
	go build -a -o ./bin/svc
