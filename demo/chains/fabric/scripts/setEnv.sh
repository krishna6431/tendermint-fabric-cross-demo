#!/usr/bin/env bash

set -eu

setGlobals() {
  USING_ORG=$1
  echo "Using organization ${USING_ORG}"
  if [[ ${USING_ORG} = "Org1" ]]; then
    export CORE_PEER_ID=peer0.org1.fabric-tendermint-cross-demo.com
    export CORE_PEER_LOCALMSPID=Org1MSP
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.fabric-tendermint-cross-demo.com/users/Admin@org1.fabric-tendermint-cross-demo.com/msp
  else
    echo "================== ERROR !!! ORG Unknown =================="
    exit 1
  fi
}

