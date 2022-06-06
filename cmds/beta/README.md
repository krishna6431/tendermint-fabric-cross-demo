# Beta CLI for Fabric

This CLI is used in [sample-scenario](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/demo/scripts/scenario/sample-scenario).

## Command List

### ConfigCommand

- `config` ... Manage configuration file
  - `init` ... Initialize account config, msp.Id

### CrossCommand

- `cross` ... Cross Framework related commands
  - `create-contract-tx` ... Create a new contract transaction
  - `tx-auth-state` ... Query the state of a client in a given path
  - `ibc-signtx` ... Sign the cross-chain transaction on other chain via the chain
  - `coordinator-state` ... Query the state of a coordinator in a given path

### IBCCommand

- `ibc`
  - `channel` ... Query the ChannelState of IBC

### ERC20Command

This command corresponds to `QueryClient` in
[./contracts/erc20/modules/erc20mgr/types/query.pb.go](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/contracts/erc20/modules/erc20mgr/types/query.pb.go)
and [req.Method in HandleContractCall()](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/contracts/erc20/modules/erc20mgr/keeper/keeper.go)

- `erc20` ... ERC20 commands
  - `mint` ... Mint token
  - `approve` ... Approve token
  - `allowance` ... Get allowance
  - `balance-of` ... Get balance
  - `total-supply` ... Get totalSupply
  - `transfer` ... Transfer token from owner account to recipient

### FabricCommand

- `fabric`
  - `id` ... Get id in contract module

## Development

Not all of them need to be developed from scratch. Project-specific commands are `FabricCommand` and `ERC20Command`.
But `CrossCommand` needs to be tweaked as a specification.
