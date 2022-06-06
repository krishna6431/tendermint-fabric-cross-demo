package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/ibc"
)

func ibcCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "ibc",
		Short:                      "IBC subcommands",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		ibc.QueryIBCChannelCmd(ctx),
	)

	return cmd
}
