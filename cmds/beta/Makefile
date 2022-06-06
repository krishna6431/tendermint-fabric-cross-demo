
.PHONY: lint
lint:
	golangci-lint run --config ../../.golangci.yml

.PHONY: lint-fix
lint-fix:
	golangci-lint run --config ../../.golangci.yml --fix

.PHONY: build
build:
	go build -o ../../demo/bin/betacli .
