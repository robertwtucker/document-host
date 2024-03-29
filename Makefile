MODULE = $(shell go list -m)
SHELL := /bin/bash
LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
GO_PKGS=$(foreach pkg, $(shell go list ./...), $(if $(findstring /vendor/, $(pkg)), , $(pkg)))
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')

VERSION=$(shell git describe --tags --always | sed 's/v//;s/-.*//')
REVISION=$(shell git rev-parse --short=7 HEAD)
PACKAGE="github.com/robertwtucker/document-host/internal/config"
IMAGE="registry.sptcloud.com/spt/docuhost"

OUTPUT_DIR=bin/server

ENV := local
ifdef $$APP_ENV
ENV := $$APP_ENV
endif

export PROJECT = github.com/robertwtucker/document-host

build:
	env GOOS=linux GOARCH=amd64 go build -o bin/server $(PROJECT)/cmd/server \
		-ldflags "-X ${PACKAGE}.appVersion=${VERSION} -X ${PACKAGE}.revision=${REVISION}" \
		-o ${OUTPUT_DIR} $(PROJECT)/cmd/server
	chmod +x bin/server

build-mac:
	env GOOS=darwin GOARCH=amd64 go build \
		-ldflags "-X ${PACKAGE}.appVersion=${VERSION} -X ${PACKAGE}.revision=${REVISION}" \
		-o ${OUTPUT_DIR} $(PROJECT)/cmd/server
	chmod +x bin/server

docker:
	docker build -t ${IMAGE}:latest -t ${IMAGE}:${VERSION} \
		--build-arg BUILD_VERSION=${VERSION} --build-arg BUILD_REVISION=${REVISION} .
run:
	go run ./cmd/server/main.go

start:
	./bin/server

test:
	go test ./... -count=1

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy -compat=1.17
	go mod vendor

deps-upgrade:
	go get -u -t -d -v ./...

deps-cleancache:
	go clean -modcache

fmt:
	@go fmt $(GO_PKGS)
	@goimports -w -l $(GO_FILES)

$(LINT_TOOL):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.43.0

qc: $(LINT_TOOL)
	$(LINT_TOOL) run --config=.golangci.yaml ./...
	staticcheck ./...
