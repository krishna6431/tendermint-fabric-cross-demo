# Fabric

This directory is for Hyperledger Fabric related including application code and configuration files for the network.

## Directory structure

This directory consists of the following files/directories

- chaincode
  - fabibc ... fabric chaincode application code.
- configtx ... these files are used for peer/orderer and channel configurations. See ["The Operations Service"](https://hyperledger-fabric.readthedocs.io/en/release-2.2/operations_service.html) and ["Using configtx.yaml to build a channel configuration"](https://hyperledger-fabric.readthedocs.io/en/release-2.2/create_channel/create_channel_config.html).
- cryptogen ... `cryptogen generate` commands generates Hyperledger Fabric key using this configuration file. See [docs](https://hyperledger-fabric.readthedocs.io/en/release-2.2/commands/cryptogen.html).
- external-builders ... configurations for [External Builders and Launchers](https://hyperledger-fabric.readthedocs.io/en/release-2.2/cc_launcher.html).
- scripts ... fabric network configurations similar to [fabric test-network](https://github.com/hyperledger/fabric-samples/tree/main/test-network).
- docker-compose.yaml ... fabric network consists of 1 orderer, 1 peer, and 1 chaincode server.

## Fabibc Application

See [README.md](https://github.com/datachainlab/fabric-tendermint-cross-demo/tree/main/demo/chains/fabric/chaincode/fabibc).
