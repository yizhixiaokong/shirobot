.PHONY: image lint help run

all: image

lint:
	golangci-lint run --timeout 10m ./...
	@echo "lint down"

image:
	docker build --pull --rm -f "Dockerfile" -t shirobot:latest "." 

run:
	docker run --rm -v /etc/timezone:/etc/timezone:ro -v /etc/localtime:/etc/localtime:ro -v $(PWD)/config:/config shirobot:latest

help:
	@echo "make: build docker image"
	@echo "make lint: golangci-lint run --timeout 10m ./..."
