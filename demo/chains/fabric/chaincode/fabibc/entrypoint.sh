#!/bin/sh

set -e

export CHAINCODE_CCID=$(cat /root/ccid.txt)

chaincode
