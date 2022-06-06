DOCKER         ?= docker
DOCKER_COMPOSE ?= docker-compose
DOCKER_REPO    ?= ""
DOCKER_BUILD   ?= $(DOCKER) build --rm --no-cache --pull

MAKEFILE_DIR:=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))

FABRIC_VERSION    ?=2.2.0
FABRIC_CA_VERSION ?=1.4.7

TENDERMINT_TAG ?= latest

PROXY_TAG ?= latest

.PHONY: wait-for-launch
wait-for-launch:
	$(MAKEFILE_DIR)scripts/wait-for-launch $(ATTEMPT) $(CONTAINER)
