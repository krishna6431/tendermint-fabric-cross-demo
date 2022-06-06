package flags

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

const (
	FlagAmount          = "amount"
	FlagOwnerAddress    = "owner-address"
	FlagReceiverAddress = "receiver-address"
	FlagSpenderAddress  = "spender-address"
)

func AddDefaultFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test|memory)")
}
