package app

import (
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"github.com/cosmos/ibc-go/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/modules/core"
	porttypes "github.com/cosmos/ibc-go/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/modules/core/keeper"
	ibcmock "github.com/cosmos/ibc-go/testing/mock"
	"github.com/datachainlab/cross-cdt/modules/erc20"
	erc20keeper "github.com/datachainlab/cross-cdt/modules/erc20/keeper"
	erc20types "github.com/datachainlab/cross-cdt/modules/erc20/types"
	cdtkeeper "github.com/datachainlab/cross-cdt/x/cdt/keeper"
	cdttypes "github.com/datachainlab/cross-cdt/x/cdt/types"
	cross "github.com/datachainlab/cross/x/core"
	crossatomic "github.com/datachainlab/cross/x/core/atomic"
	atomickeeper "github.com/datachainlab/cross/x/core/atomic/keeper"
	atomictypes "github.com/datachainlab/cross/x/core/atomic/types"
	contractkeeper "github.com/datachainlab/cross/x/core/contract/keeper"
	crosskeeper "github.com/datachainlab/cross/x/core/keeper"
	"github.com/datachainlab/cross/x/core/router"
	crosstypes "github.com/datachainlab/cross/x/core/types"
	xcctypes "github.com/datachainlab/cross/x/core/xcc/types"
	"github.com/datachainlab/cross/x/packets"
	mockclient "github.com/datachainlab/ibc-mock-client/modules/light-clients/xx-mock"
	ethmultisig "github.com/datachainlab/ibc-proxy-solidity/modules/light-clients/xx-ethmultisig"
	"github.com/hyperledger-labs/yui-fabric-ibc/app"
	"github.com/hyperledger-labs/yui-fabric-ibc/commitment"
	yuifabricauthante "github.com/hyperledger-labs/yui-fabric-ibc/x/auth/ante"
	fabric "github.com/hyperledger-labs/yui-fabric-ibc/x/ibc/light-clients/xx-fabric"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20contract"
	erc20contractkeeper "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20contract/keeper"
	erc20contracttypes "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20contract/types"
	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr"
	erc20mgrkeeper "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/keeper"
	erc20mgrtypes "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
)

var (
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		ibc.AppModuleBasic{},
		mockclient.AppModuleBasic{},
		fabric.AppModuleBasic{},
		transfer.AppModuleBasic{},
		ibcmock.AppModuleBasic{},
		cross.AppModuleBasic{},
		crossatomic.AppModuleBasic{},
		erc20.AppModuleBasic{},
		erc20mgr.AppModuleBasic{},
		erc20contract.AppModuleBasic{},
		ethmultisig.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:  nil,
		distrtypes.ModuleName:       nil,
		ibctransfertypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}
)

var _ app.Application = (*IBCApp)(nil)

type IBCApp struct {
	*app.BaseApp
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	// invCheckPeriod uint // not used anywhere

	// keys to access the substores
	keys map[string]*sdk.KVStoreKey
	// tkeys   map[string]*sdk.TransientStoreKey // not used anywhere
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper       authkeeper.AccountKeeper
	BankKeeper          bankkeeper.Keeper
	CapabilityKeeper    *capabilitykeeper.Keeper
	StakingKeeper       stakingkeeper.Keeper
	SlashingKeeper      slashingkeeper.Keeper
	MintKeeper          mintkeeper.Keeper
	DistrKeeper         distrkeeper.Keeper
	GovKeeper           govkeeper.Keeper
	CrisisKeeper        crisiskeeper.Keeper
	UpgradeKeeper       upgradekeeper.Keeper
	ParamsKeeper        paramskeeper.Keeper
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper      evidencekeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper
	CrossKeeper         crosskeeper.Keeper
	AtomicKeeper        atomickeeper.Keeper
	ERC20Keeper         erc20keeper.Keeper
	ERC20mgrKeeper      erc20mgrkeeper.Keeper
	ERC20contractKeeper erc20contractkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedCrossKeeper    capabilitykeeper.ScopedKeeper
	ScopedIBCMockKeeper  capabilitykeeper.ScopedKeeper

	// other modules
	XCCResolver xcctypes.XCCResolver

	// the module manager
	mm *module.Manager

	// the configurator
	configurator module.Configurator
}

func NewIBCApp(appName string, logger log.Logger, db dbm.DB, traceStore io.Writer, encodingConfig simappparams.EncodingConfig, seqMgr commitment.SequenceManager, blockProvider app.BlockProvider, anteHandlerProvider app.AnteHandlerProvider) (*IBCApp, error) {
	// TODO: Remove cdc in favor of appCodec once all modules are migrated.
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := app.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxJSONDecoder())
	bApp.SetInterfaceRegistry(interfaceRegistry)
	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey,
		stakingtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		crosstypes.StoreKey, erc20types.StoreKey, erc20mgrtypes.StoreKey, erc20contracttypes.StoreKey,
	)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	app := &IBCApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		memKeys:           memKeys,
	}

	// init params keeper and subspaces
	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedCrossKeeper := app.CapabilityKeeper.ScopeToModule(crosstypes.ModuleName)

	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	scopedIBCMockKeeper := app.CapabilityKeeper.ScopeToModule(ibcmock.ModuleName)

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)
	// Create IBC Keeper
	ibcKeeper := ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)
	app.IBCKeeper = overrideIBCClientKeeper(*ibcKeeper, appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), seqMgr)
	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)

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

	// Setup a cross module
	app.XCCResolver = xcctypes.NewChannelInfoResolver(app.IBCKeeper.ChannelKeeper)
	cmgr := contractkeeper.NewContractManager(
		appCodec,
		crosstypes.NewPrefixStoreKey(keys[crosstypes.StoreKey], crosstypes.ContractManagerPrefix),
		erc20contractModule,
		cdtStore,
		erc20mgrtypes.CDTContractHandleDecorators(),
	)
	app.AtomicKeeper = atomickeeper.NewKeeper(
		appCodec, crosstypes.NewPrefixStoreKey(keys[crosstypes.StoreKey], crosstypes.AtomicKeyPrefix),
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper, scopedCrossKeeper,
		cmgr, app.XCCResolver, packets.NewNOPPacketMiddleware(),
	)
	crossAtomicModule := crossatomic.NewAppModule(appCodec, app.AtomicKeeper)

	router := router.NewRouter()
	crossAtomicModule.RegisterPacketRoutes(router)

	// Create Cross Keepers
	app.CrossKeeper = crosskeeper.NewKeeper(
		appCodec,
		crosstypes.NewPrefixStoreKey(keys[crosstypes.StoreKey], crosstypes.InitiatorKeyPrefix),
		crosstypes.NewPrefixStoreKey(keys[crosstypes.StoreKey], crosstypes.AuthKeyPrefix),
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		scopedCrossKeeper,
		packets.NewNOPPacketMiddleware(),
		app.XCCResolver,
		app.AtomicKeeper,
		router,
	)
	crossModule := cross.NewAppModule(appCodec, app.CrossKeeper)

	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	mockModule := ibcmock.NewAppModule(scopedIBCMockKeeper, &app.IBCKeeper.PortKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	ibcRouter.AddRoute(crosstypes.ModuleName, crossModule)
	ibcRouter.AddRoute(ibcmock.ModuleName, mockModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		crossModule,
		crossAtomicModule,
		erc20Module,
		erc20mgrModule,
		erc20contractModule,
		mockModule,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName, distrtypes.ModuleName, stakingtypes.ModuleName,
		ibchost.ModuleName, ibctransfertypes.ModuleName, ibcmock.ModuleName, erc20types.ModuleName, erc20mgrtypes.ModuleName,
		crosstypes.ModuleName, atomictypes.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)

	app.SetAnteHandler(anteHandlerProvider(ibckeeper.Keeper{}, nil))
	app.SetBlockProvider(blockProvider)

	if err := app.LoadLatestVersion(); err != nil {
		return nil, err
	}

	// Initialize and seal the capability keeper so all persistent capabilities
	// are loaded in-memory and prevent any further modules from creating scoped
	// sub-keepers.
	ctx, writer := app.MakeCacheContext(tmproto.Header{})
	app.CapabilityKeeper.InitializeAndSeal(ctx)
	writer()

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedCrossKeeper = scopedCrossKeeper

	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	app.ScopedIBCMockKeeper = scopedIBCMockKeeper

	return app, nil
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(erc20mgrtypes.ModuleName)

	return paramsKeeper
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *IBCApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *IBCApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns IBCApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *IBCApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns IBCApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *IBCApp) AppCodec() codec.Codec {
	return app.appCodec
}

func (app *IBCApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

func (app *IBCApp) InitChainer(ctx sdk.Context, appStateBytes []byte) error {
	var genesisState simapp.GenesisState
	if err := tmjson.Unmarshal(appStateBytes, &genesisState); err != nil {
		return err
	}
	app.mm.InitGenesis(ctx, app.appCodec, genesisState)
	return nil
}

func (app *IBCApp) GetIBCKeeper() ibckeeper.Keeper {
	return *app.IBCKeeper
}

// DefaultAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func DefaultAnteHandler(
	ibcKeeper ibckeeper.Keeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		yuifabricauthante.NewFabricIDVerificationDecorator(),
	)
}

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func MakeEncodingConfig() simappparams.EncodingConfig {
	encodingConfig := simappparams.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() simapp.GenesisState {
	encCfg := MakeEncodingConfig()
	return ModuleBasics.DefaultGenesis(encCfg.Marshaler)
}
