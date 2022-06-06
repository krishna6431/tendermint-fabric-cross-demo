.PHONY: protoc
protoc:
	@echo "Generating Protobuf files"
	docker run -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh

.PHONY: lint
lint:
	golangci-lint run --config ../../.golangci.yml

.PHONY: lint-fix
lint-fix:
	golangci-lint run --config ../../.golangci.yml --fix

.PHONY: test
test:
	go test -race -v ./...
