package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibc "github.com/cosmos/ibc-go/modules/core"
	ibctypes "github.com/cosmos/ibc-go/modules/core/types"
	"github.com/datachainlab/cross-cdt/modules/erc20"
	erc20types "github.com/datachainlab/cross-cdt/modules/erc20/types"
	cross "github.com/datachainlab/cross/x/core"
	crosstypes "github.com/datachainlab/cross/x/core/types"
	"github.com/hyperledger-labs/yui-fabric-ibc/app"
	"github.com/hyperledger-labs/yui-fabric-ibc/chaincode"
	"github.com/hyperledger-labs/yui-fabric-ibc/commitment"
	fabrictypes "github.com/hyperledger-labs/yui-fabric-ibc/x/ibc/light-clients/xx-fabric/types"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"

	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20contract"
	erc20contracttypes "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20contract/types"
	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr"
	erc20mgrtypes "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
	fabapp "github.com/datachainlab/fabric-tendermint-cross-demo/demo/chains/chaincode/fabibc/app"
)

// TODO: Genesis state should get from a gensis json.
//  Wait for the `yui-relayer fabric chaincode init` command to be updated.
// value could be found by command `${BETACLI_TOKEN_OWNER} fabric id`
var ERC20Admin = "" // embedded when building

func main() {
	ibcCC := chaincode.NewIBCChaincode(
		"ibc1",
		tmlog.NewTMLogger(os.Stdout),
		commitment.NewDefaultSequenceManager(),
		newApp,
		fabapp.DefaultAnteHandler,
		chaincode.DefaultDBProvider,
		chaincode.DefaultMultiEventHandler(),
	)
	cc, err := contractapi.NewChaincode(ibcCC)
	if err != nil {
		panic(err)
	}

	server := &shim.ChaincodeServer{
		CCID:    os.Getenv("CHAINCODE_CCID"),
		Address: os.Getenv("CHAINCODE_ADDRESS"),
		CC:      cc,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}
	if err = server.Start(); err != nil {
		fmt.Printf("Error starting IBC chaincode: %s", err)
	}
}

func newApp(appName string, logger tmlog.Logger, db tmdb.DB, traceStore io.Writer, seqMgr commitment.SequenceManager, blockProvider app.BlockProvider, anteHandlerProvider app.AnteHandlerProvider) (app.Application, error) {
	ibcApp, err := fabapp.NewIBCApp(
		appName,
		logger,
		db,
		traceStore,
		fabapp.MakeEncodingConfig(),
		seqMgr,
		blockProvider,
		anteHandlerProvider,
	)
	if err != nil {
		return nil, err
	}
	wrapped := &IBCApp{IBCApp: ibcApp}
	ibcApp.SetInitChainer(wrapped.InitChainer)
	return wrapped, nil
}

type IBCApp struct {
	*fabapp.IBCApp
}

func (app *IBCApp) InitChainer(ctx sdk.Context, appStateBytes []byte) error {
	var genesisState simapp.GenesisState
	if err := tmjson.Unmarshal(appStateBytes, &genesisState); err != nil {
		return err
	}
	ibcGenesisState := ibctypes.DefaultGenesisState()
	ibcGenesisState.ClientGenesis.Params.AllowedClients = append(
		ibcGenesisState.ClientGenesis.Params.AllowedClients,
		fabrictypes.Fabric,
	)
	genesisState[ibc.AppModule{}.Name()] = app.AppCodec().MustMarshalJSON(ibcGenesisState)
	genesisState[cross.AppModuleBasic{}.Name()] = app.AppCodec().MustMarshalJSON(crosstypes.DefaultGenesis())
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

	bz, err := tmjson.Marshal(genesisState)
	if err != nil {
		return err
	}
	return app.IBCApp.InitChainer(ctx, bz)
}
