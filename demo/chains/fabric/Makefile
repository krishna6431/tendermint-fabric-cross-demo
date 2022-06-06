include ../../var.mk

FABRIC_BIN_DIR ?= ${CURDIR}/.bin
export PATH := ${FABRIC_BIN_DIR}:${PATH}
export FABRIC_CFG_PATH=${CURDIR}/configtx

SCRIPT_DIR ?= ${CURDIR}/scripts

.PHONY: network
network:
	make cryptogen
	make configtx
	make network-common
	make channel
	make chaincode
	make docker-images
	make network-chaincode

.PHONY: network-common
network-common:
	IMAGE_TAG=${FABRIC_VERSION} ${DOCKER_COMPOSE} up -d \
	fabric-orderer.fabric-tendermint-cross-demo.com \
	fabric-peer0.org1.fabric-tendermint-cross-demo.com

.PHONY: network-chaincode
network-chaincode:
	IMAGE_TAG=${FABRIC_VERSION} ${DOCKER_COMPOSE} up -d \
	fabric-chaincode.org1.fabric-tendermint-cross-demo.com

.PHONY: network-down
network-down:
	IMAGE_TAG=${FABRIC_VERSION} ${DOCKER_COMPOSE} down -v --remove-orphans

.PHONY: clean
clean:
	rm -rf artifacts
	rm -rf backup
	rm -rf build
	rm -rf ${CURDIR}/external-builders/config/**/code.tar.gz
	rm -rf organizations
	make -C chaincode/fabibc clean

.PHONY: bin
bin:
	rm -rf ${CURDIR}/.bin && \
	curl -sSL https://bit.ly/2ysbOFE | bash -s -- ${FABRIC_VERSION} ${FABRIC_CA_VERSION} -d -s && \
	rm -rf ${CURDIR}/config && \
	mv ${CURDIR}/bin ${CURDIR}/.bin

.PHONY: cryptogen
cryptogen:
	${FABRIC_BIN_DIR}/cryptogen generate --config ${CURDIR}/cryptogen/crypto-config.yaml --output ./organizations

.PHONY: configtx
configtx:
	${SCRIPT_DIR}/configtx.sh

.PHONY: channel
channel:
	${SCRIPT_DIR}/createChannel.sh

.PHONY: chaincode
chaincode:
	${SCRIPT_DIR}/deployCC.sh

.PHONY: docker-images
docker-images:
	make prepare-beta-address
	make -C chaincode/fabibc docker-chaincode

# TODO: temporarily generated address is passed through docker image build for chaincode
.PHONY: prepare-beta-address
prepare-beta-address:
	@../../scripts/init-betacli
	@../../scripts/util/beta-admin-address ./chaincode/fabibc/address
