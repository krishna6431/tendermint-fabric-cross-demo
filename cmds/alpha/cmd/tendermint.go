package cmd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	alphaflags "github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/flags"
)

func tendermintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tendermint",
		Short:                      "Tendermint subcommands",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		accountIDTendermintCmd(),
	)

	return cmd
}

func accountIDTendermintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-id",
		Short: "Get account id",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fmt.Println(clientCtx.GetFromAddress().String())
			return nil
		},
	}
	alphaflags.AddDefaultFlagsToCmd(cmd)

	return cmd
}
