package cmd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
)

func fabricCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "fabric",
		Short:                      "Fabric subcommands",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		idFabricCmd(ctx),
	)

	return cmd
}

func idFabricCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "id",
		Short: "Get id in contract module",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(ctx.Address.String())
			return nil
		},
	}

	return cmd
}
