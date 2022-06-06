#!/usr/bin/env bash
set -eu pipefail

SCRIPT_DIR=$(cd $(dirname $0);pwd)
PROJECT_DIR=$SCRIPT_DIR/..

CHAINCODE_CCID_ORG1=$(cat ${PROJECT_DIR}/build/Org1-fabibc-ccid.txt)

set -x
CHAINCODE_CCID_ORG1=${CHAINCODE_CCID_ORG1} \
docker-compose -f docker-compose-chaincode.yaml up -d \
