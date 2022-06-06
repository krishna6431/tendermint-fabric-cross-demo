# Alpha CLI for Tendermint

This CLI is used in [sample-scenario](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/demo/scripts/scenario/sample-scenario).

## Command List

### InitCommand

It refers to [cosmos-sdk: genutil cli](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/x/genutil/client/cli/init.go)

- `init` ... Initialize private validator, p2p, genesis, and application configuration files.
  Generated`.alpha` directory includes keys in the `demo` directory.

### KeysCommand

It refers to [cosmos-sdk: client keys](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/client/keys/root.go)

- `keys` ... Manage application keys

### AddGenesisAccountCommand

- `add-genesis-account` ... Add a genesis account to genesis.json

### GenTxCommand

It refers to [cosmos-sdk: genutil cli](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/x/genutil/client/cli/gentx.go)

- `gentx` ... Generate a genesis tx carrying a self delegation

### CollectGenTxsCommand

It refers to [cosmos-sdk: genutil cli](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/x/genutil/client/cli/collect.go)

- `collect-gentxs` ... Collect genesis txs and output a genesis.json file

### NewTxCommand

It refers to [cosmos-sdk: bank cli](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/x/bank/client/cli/tx.go)

- `bank` ... Bank transaction subcommands

### QueryCommand

- `query` ... It contains various query sub commands
  - `account` ... Query account
  - `tendermint-validator-set` ... Get the full tendermint validator set at given height
  - `block` ... Get verified data for the block at given height
  - `txs` ... Query for paginated transactions that match a set of events
  - `tx` ... Query for a transaction by hash in a committed block

### CrossCommand

- `cross` ... Cross Framework related commands
  - `create-contract-tx` ... Create a new contract transaction
  - `tx-auth-state` ... Query the state of a client in a given path
  - `create-initiate-tx` ... Create and submit a NewMsgInitiateTx transaction for a simple commit
  - `coordinator-state` ... Query the state of a coordinator in a given path

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

### TendermintCommand

- `tendermint`
  - `account-id` ... Get account id

## Development

Not all of them need to be developed from scratch. Project-specific commands are `TendermintCommand` and `ERC20Command`.
But `CrossCommand` needs to be tweaked as a specification.
