# Fabibc Application

These are packages used by go implementations of Fabric chaincode based on [fabric-ibc](https://github.com/datachainlab/public-docs/tree/master/fabric-ibc).  
Refer to [yui-fabric-ibc](https://github.com/hyperledger-labs/yui-fabric-ibc), [Fabric chaincode lifecycle](https://hyperledger-fabric.readthedocs.io/en/release-2.2/chaincode_lifecycle.html).

## Deployment

Built binary is deployed as Docker container using [Dockerfile](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/demo/chains/fabric/chaincode/fabibc/Dockerfile) in this directory refered from [docker-compose.yaml](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/demo/chains/fabric/docker-compose.yaml).

### Command and timing

The actual deployment timing is below. See [Makefile](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/demo/Makefile).

```
- [make: demo network] -> [make: chains/fabric network] -> [make:  docker-images] -> [make: chaincode/fabibc docker-chaincode]
- [make: demo network] -> [make: chains/fabric network] -> [make: network-chaincode]
```

## How to add ERC20 modules into fabibc

### [main.go](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/demo/chains/fabric/chaincode/fabibc/main.go)

- Add `genesisState` in InitChainer(). [code](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/demo/chains/fabric/chaincode/fabibc/main.go)

```
// erc20mgr module
erc20mgrGenesisState := erc20mgrtypes.DefaultGenesis()
erc20mgrGenesisState.Params = erc20mgrtypes.NewParams(
	ERC20Admin,
	false,
)
genesisState[erc20mgr.AppModuleBasic{}.Name()] = app.AppCodec().MustMarshalJSON(erc20mgrGenesisState)
// erc20 module in cross-cdt
genesisState[erc20.AppModuleBasic{}.Name()] = app.AppCodec().MustMarshalJSON(erc20types.DefaultGenesis())
// erc20cotract module
genesisState[erc20contract.AppModuleBasic{}.Name()] = app.AppCodec().MustMarshalJSON(erc20contracttypes.DefaultGenesis())
```

### [app/app.go](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/demo/chains/fabric/chaincode/fabibc/app/app.go)

- Create `BasicManager` with additional `AppModuleBasic{}`.

````go
ModuleBasics = module.NewBasicManager(
	...
	erc20.AppModuleBasic{},
	erc20mgr.AppModuleBasic{},
	erc20contract.AppModuleBasic{},
	```
)
````

- Add `Keeper`s into `IBCApp` struct.

```go
type IBCApp struct {
	...
	ERC20Keeper         erc20keeper.Keeper
	ERC20mgrKeeper      erc20mgrkeeper.Keeper
	ERC20contractKeeper erc20contractkeeper.Keeper
	...
}
```

- Add `StoreKeys` when creating `KVStoreKeys` in `NewIBCApp()`.

```go
keys := sdk.NewKVStoreKeys(
	...
	erc20types.StoreKey, erc20mgrtypes.StoreKey, erc20contracttypes.StoreKey,
)
```

- Create modules in `NewIBCApp()`.
  - Note. [cross-cdt](https://github.com/datachainlab/cross-cdt) data type is used for `Store` in this demo. So some code depend on CDT.

```go
// Create CDT Store
schema := cdttypes.NewSchema()
schemaERC20Prefix := []byte(erc20mgrtypes.ModuleName + "/")
schema.Set(schemaERC20Prefix, cdttypes.CDT_TYPE_INT64)
cdtStore := cdtkeeper.NewStore(appCodec, keys[crosstypes.StoreKey], schema)

// Create ERC20mgr module
erc20Int64Store := cdtStore.GetInt64Store(schemaERC20Prefix)
app.ERC20Keeper = erc20keeper.NewKeeper(erc20Int64Store)
erc20Module := erc20.NewAppModule(app.ERC20Keeper)
app.ERC20mgrKeeper = erc20mgrkeeper.NewKeeper(appCodec, app.ERC20Keeper, app.GetSubspace(erc20mgrtypes.ModuleName))
erc20mgrModule := erc20mgr.NewAppModule(app.ERC20mgrKeeper)

// Create erc20contract module
app.ERC20contractKeeper = erc20contractkeeper.NewKeeper(appCodec, keys[erc20contracttypes.StoreKey], app.ERC20mgrKeeper)
erc20contractModule := erc20contract.NewAppModule(app.ERC20contractKeeper)
```

- Setup a cross module using CDT in `NewIBCApp()`.

```go
// Setup a cross module
app.XCCResolver = xcctypes.NewChannelInfoResolver(app.IBCKeeper.ChannelKeeper)
cmgr := contractkeeper.NewContractManager(
	appCodec,
	crosstypes.NewPrefixStoreKey(keys[crosstypes.StoreKey], crosstypes.ContractManagerPrefix),
	erc20contractModule,
	cdtStore,
	erc20mgrtypes.CDTContractHandleDecorators(),
)
```

- Create `Manager` with additional modules in `NewIBCApp()`.

```go
	app.mm = module.NewManager(
		...
		erc20Module,
		erc20mgrModule,
		erc20contractModule,
	)
```

- Add `ModuleName` when calling `SetOrderInitGenesis` in `NewIBCApp()`.

```go
app.mm.SetOrderInitGenesis(
	...
	erc20types.ModuleName, erc20mgrtypes.ModuleName,
)
```

- Add `Subspace` in `initParamsKeeper()` if used.

```go
paramsKeeper.Subspace(erc20mgrtypes.ModuleName)
```
