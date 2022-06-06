package cross

import "github.com/spf13/cobra"

const (
	// flagLightHeight = "light-height"  // not used anywhere
	flagEthSignKey = "eth-sign-key"
)

// nolint: unparam
func ethSignKeyFlag(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().String(flagEthSignKey, "", "the Ethereum Chain private key used by the importer for signing")
	return cmd
}
