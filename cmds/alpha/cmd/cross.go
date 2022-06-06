package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/config"
	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/cross"
	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/cross/atomic"
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
		cross.NewInitiateTxCmd(ctx),
		atomic.GetCoordinatorState(ctx),
	)

	return cmd
}
