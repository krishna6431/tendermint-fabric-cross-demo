package atomic

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
)

func GetCoordinatorState(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coordinator-state [tx-id]",
		Short: "Query the state of a coordinator in a given path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if txId, err := hex.DecodeString(args[0]); err != nil {
				return err
			} else if res, err := ctx.Chain.QueryCoordinatorState(txId); err != nil {
				return err
			} else if bz, err := ctx.Chain.Codec().MarshalJSON(&res.CoodinatorState); err != nil {
				return err
			} else {
				fmt.Println(string(bz))
			}

			return nil
		},
	}
	return cmd
}
