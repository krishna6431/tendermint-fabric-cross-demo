GOLANGCI_VERSION=v1.45.2

# version should be updated periodically
.PHONY: install-tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)

.PHONY: lint-all
lint-all:
	make -C cmds/alpha lint
	make -C cmds/beta lint
	make -C contracts/erc20 lint
	make -C demo/chains/fabric/chaincode/fabibc lint
	make -C demo/chains/tendermint lint

.PHONY: lint-fix-all
lint-fix-all:
	make -C cmds/alpha lint-fix
	make -C cmds/beta lint-fix
	make -C contracts/erc20 lint-fix
	make -C demo/chains/fabric/chaincode/fabibc lint-fix
	make -C demo/chains/tendermint lint-fix

.PHONY: prepare-demo-env
prepare-demo-env:
	make -C demo build -j5
	make -C demo network-down
	make -C demo clean
	make -C demo network
	make -C demo run-init
