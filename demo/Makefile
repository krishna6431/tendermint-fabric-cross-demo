
# run in parallel by `make build -j5`
# `-j` option (short for –jobs) specifies the number of jobs that are run at the same time
.PHONY: build
build:
	make -C ../relayer build && mkdir -p ./bin && cp ../relayer/build/uly ./bin/ & \
	make -C chains/fabric bin & \
	cd chains/tendermint && go mod vendor && make docker-images & \
	make -C ../cmds/alpha build & \
	make -C ../cmds/beta build

.PHONY: network
network:
	make -C chains/fabric network
	make -C chains/tendermint network

.PHONY: network-down
network-down:
	make -C chains/fabric network-down
	make -C chains/tendermint network-down

.PHONY: clean
clean:
	rm -rf .urelayer
	rm -rf .alphacli
	rm -rf .betacli
	rm -rf ./data
	rm -rf ./fixtures
	make -C chains/fabric clean
	make -C chains/tendermint clean

.PHONY: restart-demo
restart-demo: network-down clean network run-init

.PHONY: run-init
run-init:
	@./scripts/fixture
	@./scripts/init-rly
	@./scripts/handshake
	@./scripts/init-alphacli
	@./scripts/init-betacli

.PHONY: run
run:
	@./scripts/scenario/sample-scenario

.PHONY: test-alphacli
test-alphacli:
	@./scripts/scenario/sample-alphacli

.PHONY: test-betacli
test-betacli:
	@./scripts/scenario/sample-betacli
