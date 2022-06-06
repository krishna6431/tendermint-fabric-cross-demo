include ../../../../var.mk

.PHONY: lint
lint:
	golangci-lint run --config ../../../../../.golangci.yml

.PHONY: lint-fix
lint-fix:
	golangci-lint run --config ../../../../../.golangci.yml --fix

# just making sure build works
.PHONY: build
build:
	go build -o ./build/fabibc .
	rm -rf ./build

# need to be built with latest address file
# see Makefile `docker-images` in chains/fabric directory
.PHONY: docker-chaincode
docker-chaincode:
	go mod vendor
	${DOCKER_BUILD} \
		--build-arg ADDR=$(shell cat address) \
		--tag fabric-chaincode-fabibc:latest .

.PHONY: clean
clean:
	rm -rf vendor
	rm -rf address
