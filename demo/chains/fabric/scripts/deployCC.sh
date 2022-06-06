#!/usr/bin/env bash
set -eux pipefail

VERSION=${VERSION:="1"}

# import utils
. ./scripts/setEnv.sh

OREDERER_ENDPOINT=localhost:7050
CC_NAME=${CC_NAME:="fabibc"}

packageChaincode() {
    ORG_NAME=$1

    echo "### Package Chaincode ${CC_NAME}"
    pushd ./external-builders/config/${ORG_NAME}/${CC_NAME}
    tar cfz code.tar.gz connection.json
    tar cfz ${ORG_NAME}-${CC_NAME}.tar.gz code.tar.gz metadata.json
    mv ${ORG_NAME}-${CC_NAME}.tar.gz ../../../../build
    popd
}

installChaincode() {
    ORG_NAME=$1

    echo "### Install Chaincode ./build/${ORG_NAME}-${CC_NAME}.tar.gz"
    setGlobals ${ORG_NAME}
    set -x
    peer lifecycle chaincode install ./build/${ORG_NAME}-${CC_NAME}.tar.gz
    set +x
}

queryInstalled() {
    ORG_NAME=$1

    setGlobals ${ORG_NAME}
    set -x
    peer lifecycle chaincode queryinstalled >&log.txt
    set +x
    cat log.txt
    PACKAGE_ID=$(sed -n "/${CC_NAME}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
    echo ${PACKAGE_ID} &> ./build/${ORG_NAME}-${CC_NAME}-ccid.txt
}

approveForMyOrg() {
    ORG_NAME=$1
    CHANNEL_NAME=$2
    SIGNATURE_POLICY=$3

    setGlobals ${ORG_NAME}
    PACKAGE_ID=$(cat ./build/${ORG_NAME}-${CC_NAME}-ccid.txt)
    set -x
    peer lifecycle chaincode approveformyorg \
    -o ${OREDERER_ENDPOINT} \
    --channelID ${CHANNEL_NAME} \
    --name ${CC_NAME} \
    --version ${VERSION} \
    --sequence ${VERSION} \
    --package-id ${PACKAGE_ID} \
    --signature-policy ${SIGNATURE_POLICY}

    peer lifecycle chaincode checkcommitreadiness \
    --channelID ${CHANNEL_NAME} \
    --name ${CC_NAME} \
    --version ${VERSION} \
    --sequence ${VERSION} \
    --signature-policy ${SIGNATURE_POLICY} \
    --output json
    set +x
}

commitChaincode() {
    ORG_NAME=$1
    CHANNEL_NAME=$2
    SIGNATURE_POLICY=$3

    PEER_CONN_PARAMS=""
    for org in ${@:4}; do
      setGlobals ${org}
      PEER_CONN_PARAMS="$PEER_CONN_PARAMS --peerAddresses $CORE_PEER_ADDRESS"
    done

    setGlobals ${ORG_NAME}
    set -x
    peer lifecycle chaincode commit \
    -o ${OREDERER_ENDPOINT} \
    --channelID ${CHANNEL_NAME} \
    --name ${CC_NAME} \
    --version ${VERSION} \
    --sequence ${VERSION} ${PEER_CONN_PARAMS} \
    --signature-policy ${SIGNATURE_POLICY}

    peer lifecycle chaincode querycommitted --channelID ${CHANNEL_NAME} --name ${CC_NAME}
    set +x
}

mkdir -p ./build

CHANNEL_NAME=channel1
SIGNATURE_POLICY="AND('Org1MSP.peer')"

ORG_NAME="Org1"
packageChaincode ${ORG_NAME}
installChaincode ${ORG_NAME}
queryInstalled ${ORG_NAME}
approveForMyOrg ${ORG_NAME} ${CHANNEL_NAME} ${SIGNATURE_POLICY}

ORG_NAME="Org1"
queryInstalled ${ORG_NAME}
commitChaincode ${ORG_NAME} ${CHANNEL_NAME} ${SIGNATURE_POLICY}
