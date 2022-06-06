package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/cross"
	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/cross/atomic"
)

func crossCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "cross",
		Short:                      "Cross subcommands",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cross.GetCreateContractTransaction(ctx),
		cross.QueryTxAuthStateCmd(ctx),
		cross.NewIBCSignTxCmd(ctx),
		atomic.GetCoordinatorState(ctx),
	)

	return cmd
}
