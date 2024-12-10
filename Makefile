SHELL := /bin/bash
VERSION=$(shell git describe --tags --always | sed 's/v//;s/-.*//')
IMAGE="registry.sptcloud.com/spt/docuhost"

build:
	npm run build

docker:
	docker buildx build \
		--tag ${IMAGE}:latest \
		--tag ${IMAGE}:${VERSION} \
		--platform linux/arm64,linux/amd64 \
		--push .

docker-amd:
	docker buildx build \
		--tag ${IMAGE}:latest \
		--tag ${IMAGE}:${VERSION} \
		--platform linux/amd64 \
		--load .

docker-arm:
	docker buildx build \
		--tag ${IMAGE}:latest \
		--tag ${IMAGE}:${VERSION} \
		--platform linux/arm64 \
		--load .
