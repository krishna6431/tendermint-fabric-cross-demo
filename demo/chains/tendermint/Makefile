include ../../var.mk

###############################################################################
# Development
###############################################################################
.PHONY: lint
lint:
	golangci-lint run --config ../../../.golangci.yml

.PHONY: lint-fix
lint-fix:
	golangci-lint run --config ../../../.golangci.yml --fix

# just making sure build works
.PHONY: build
build:
	go build -o ./build/simd ./simapp/simd
	rm -rf ./build

###############################################################################
# Network, Container
###############################################################################

.PHONE: network
network:
	TAG=${TENDERMINT_TAG} $(DOCKER_COMPOSE) \
		-f docker-compose.yaml \
		up -d \
		tendermint-chain0

.PHONY: network-down
network-down:
	TAG=${TENDERMINT_TAG} $(DOCKER_COMPOSE) \
		-f ./docker-compose.yaml \
		down --volumes --remove-orphans

vendor:
	go mod vendor

.PHONY: docker-images
docker-images: vendor
	$(DOCKER_BUILD) \
		--build-arg CHAINID=ibc0 \
		--tag $(DOCKER_REPO)tendermint-chain0:$(TENDERMINT_TAG) .

.PHONY: clean
clean:
	rm -rf vendor
