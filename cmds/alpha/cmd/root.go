package cmd

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	sdkconfig "github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	cross "github.com/datachainlab/cross/x/core"
	"github.com/hyperledger-labs/yui-relayer/chains/tendermint"
	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/config"
	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20contract"
	erc20contracttypes "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20contract/types"
	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr"
	erc20types "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
)

var defaultHome = os.ExpandEnv("$HOME/.alpha")

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// context must be configured by svrcmd.Execute()
	if err := svrcmd.Execute(NewRootCmd(), defaultHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}

// NewRootCmd creates a new root command
func NewRootCmd() *cobra.Command {
	// refer to demo/chains/tendermint/simapp/app.go L143
	moduleBasics := module.NewBasicManager(
		auth.AppModuleBasic{},
		staking.AppModuleBasic{},
		bank.AppModuleBasic{},
		cross.AppModuleBasic{},
		erc20mgr.AppModuleBasic{},
		erc20contract.AppModuleBasic{},
	)
	encodingConfig := makeEncodingConfig(moduleBasics)

	initClientCtx := client.Context{}.
		WithJSONCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithSkipConfirmation(true).
		WithHomeDir(defaultHome).
		WithViper("") // In simapp, we don't use any prefix for env variables.

	rootCmd := &cobra.Command{
		Use:   "alphacli",
		Short: "This provides commands for cross on a Tendermint chain",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = client.ReadHomeFlag(initClientCtx, cmd)

			initClientCtx, err := sdkconfig.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd)
		},
		SilenceUsage: true,
	}

	initRootCmd(rootCmd, moduleBasics, encodingConfig)

	return rootCmd
}

func makeEncodingConfig(moduleBasics module.BasicManager) simappparams.EncodingConfig {
	encodingConfig := simappparams.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	moduleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	moduleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

func initRootCmd(rootCmd *cobra.Command, moduleBasics module.BasicManager, encodingConfig simappparams.EncodingConfig) {
	cobra.EnableCommandSorting = false

	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	// Register interfaces
	tendermint.RegisterInterfaces(interfaceRegistry)
	erc20types.RegisterInterfaces(interfaceRegistry)
	erc20contracttypes.RegisterInterfaces(interfaceRegistry)
	cryptotypes.RegisterInterfaces(interfaceRegistry)
	cross.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)

	ctx := &config.Context{Codec: cdc}

	// Register subcommands
	// refer to `demo/chains/tendermint/simapp/simd/cmd/root.go`
	rootCmd.AddCommand(
		genutilcli.InitCmd(moduleBasics, defaultHome),
		keys.Commands(defaultHome),
		addGenesisAccountCmd(defaultHome),
		genutilcli.GenTxCmd(moduleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, defaultHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, defaultHome),
		bankcli.NewTxCmd(),
		queryCommand(moduleBasics),
		tendermintCmd(),
		crossCmd(ctx),
		erc20Cmd(),
	)
}
