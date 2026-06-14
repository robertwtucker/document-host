SHELL := /bin/bash
VERSION=$(shell git describe --tags --always | sed 's/v//;s/-.*//')
IMAGE="registry.sptcloud.com/spt/docuhost"

COSIGN_DIR ?= $(HOME)/secrets/cosign/document-host
COSIGN_KEY ?= $(COSIGN_DIR)/cosign.key
COSIGN_PASSWD_ENC ?= $(COSIGN_DIR)/passwd.enc

build:
	pnpm run build

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

cosign:
	@test -f "$(COSIGN_KEY)" || { echo "missing key: $(COSIGN_KEY)"; exit 1; }
	@test -f "$(COSIGN_PASSWD_ENC)" || { echo "missing passwd: $(COSIGN_PASSWD_ENC)"; exit 1; }
	@digest=$$(docker buildx imagetools inspect ${IMAGE}:${VERSION} --format '{{json .Manifest}}' 2>/dev/null | jq -r '.digest'); \
	 if [ -z "$$digest" ] || [ "$$digest" = "null" ]; then echo "could not resolve digest for ${IMAGE}:${VERSION} -- run 'make docker' first"; exit 1; fi; \
	 echo "Signing ${IMAGE}@$$digest"; \
	 COSIGN_PASSWORD="$$(keybase decrypt -i $(COSIGN_PASSWD_ENC))" \
	   cosign sign --yes --key $(COSIGN_KEY) ${IMAGE}@$$digest

verify:
	@digest=$$(docker buildx imagetools inspect ${IMAGE}:${VERSION} --format '{{json .Manifest}}' 2>/dev/null | jq -r '.digest'); \
	 if [ -z "$$digest" ] || [ "$$digest" = "null" ]; then echo "could not resolve digest for ${IMAGE}:${VERSION}"; exit 1; fi; \
	 cosign verify --key cosign.pub ${IMAGE}@$$digest

release: docker cosign
	@echo "Released and signed ${IMAGE}:${VERSION}"
