package cmd

import (
	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
)

func configCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "manage configuration file",
	}

	cmd.AddCommand(
		config.InitConfigCmd(ctx),
	)

	return cmd
}
