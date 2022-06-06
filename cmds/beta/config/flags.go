package config

import (
	"github.com/spf13/cobra"
)

const (
	flagFabClientCertPath       = "cert"
	flagFabClientPrivateKeyPath = "key"
	flagConfigFile              = "config"
)

func populateWalletFlag(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().String(flagFabClientCertPath, "", "a path of client cert file")
	cmd.Flags().String(flagFabClientPrivateKeyPath, "", "a path of client private key file")
	_ = cmd.MarkFlagRequired(flagFabClientCertPath)
	_ = cmd.MarkFlagRequired(flagFabClientPrivateKeyPath)
	return cmd
}

func configFileFlag(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringP(flagConfigFile, "c", "", "config file")
	_ = cmd.MarkFlagRequired(flagConfigFile)
	return cmd
}
