# ERC20 Contract

Each of the modules is in the `modules` directory.  
These contract modules are developed based on [Cosmos SDK](https://docs.cosmos.network/v0.44/).
Note that these modules depend on version `v0.43.0-beta1` for now. See [go.mod](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/contracts/erc20/go.mod).

Each of modules needs to meet specific [Interfaces](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/contracts/erc20/modules/erc20contract/module.go).

```go
_ module.AppModule             = AppModule{}
_ module.AppModuleBasic        = AppModuleBasic{}
_ contracttypes.ContractModule = AppModule{}
```

Refer to source code

- [module.AppModule](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/types/module/module.go#L156-L183)
- [module.AppModuleBasic](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/types/module/module.go#L47-L60)
- [contracttypes.ContractModule](https://github.com/datachainlab/cross/blob/v0.2.2/x/core/contract/types/types.go#L13-L15)

See section [Introduction to SDK Modules](https://docs.cosmos.network/v0.44/building-modules/intro.html) to know Cosmos SDK Modules in detail.

## Erc20mgr module

This module has ERC20 functionalities based on [CDT](https://github.com/datachainlab/cross-cdt) according to [EIP-20: Token Standard](https://eips.ethereum.org/EIPS/eip-20).

### Functionalities based on [CDT](https://github.com/datachainlab/cross-cdt)

1. Mint(account string, amount int64)

- create `amount` tokens and assigns them to `account`, increasing the total supply

2. Burn(account string, amount int64)

- destroy `amount` tokens from `account`, reducing the total supply

3. Transfer(spender, recipient string, amount int64)

- move `amount` of tokens from `sender` to `recipient`

4. Approve(owner string, spender string, amount int64)

- approve sets `amount` as the allowance of `spender` over the caller's tokens

5. TransferFrom(owner string, spender string, recipient string, amount int64)

- move `amount` tokens from `sender` to `recipient` using the allowance mechanism. `amount` is then deducted from the caller's allowance

6. Allowance(owner string, spender string)

- return the remaining number of tokens that `spender` will be allowed to spend on behalf of `owner` through `transferFrom`. This is zero by default.

7. BalanceOf(account string)

- return the amount of tokens owned by `account`

8. TotalSupply()

- return the amount of tokens in existence

### Handler and process flow

This module includes `handler.go`

1. Message [MsgContractCallTx](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/contracts/erc20/modules/erc20mgr/types/msgs.pb.go) is submitted from CLI.
2. `NewHandler()` is called from `Route()` in `module.go`
3. `handleContractCallTx()` in `NewHandler()` is called
4. `HandleContractCall()` in `keepr.go` is called
5. Then it calls function according to request method

- See [HandleContractCall function](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/contracts/erc20/modules/erc20mgr/keeper/keeper.go).

## Erc20contract module

This module is a [Contract Module](https://datachainlab.github.io/cross-docs/architecture/overview/#contract-module) of [Cross Framework](https://github.com/datachainlab/cross) called from [Contract Manager](https://datachainlab.github.io/cross-docs/architecture/overview/#contract-manager).

### Functionalities

1. Transfer

- call `Transfer` of erc20mgr contract

### Handler and process flow

This module is called from only `Cross framework`.

1. [OnContractCall()](https://datachainlab.github.io/cross-docs/architecture/overview#contract-module) in `module.go` is called from `Cross framework`
2. `contractHandler()` set by `keeper.HandleContractCall` in called
3. Then it calls function according to request method

- See [HandleContractCall function](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/contracts/erc20/modules/erc20contract/keeper/keeper.go).
