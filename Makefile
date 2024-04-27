SHELL := /bin/bash
VERSION=$(shell git describe --tags --always | sed 's/v//;s/-.*//')
IMAGE="registry.sptcloud.com/spt/docuhost"

build:
	npm run build

docker:
	docker build -t ${IMAGE}:latest -t ${IMAGE}:${VERSION} .
