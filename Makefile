.PHONY: build

VERSION=latest

build:
	DOCKER_BUILDKIT_NAME=1 docker build -f build/Dockerfile -t zura:$(VERSION) .