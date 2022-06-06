package atomic

import (
	"context"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	atomictypes "github.com/datachainlab/cross/x/core/atomic/types"
	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/config"
)

func GetCoordinatorState(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coordinator-state [tx-id]",
		Short: "Query the state of a coordinator in a given path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			q := atomictypes.NewQueryClient(clientCtx)

			if txId, err := hex.DecodeString(args[0]); err != nil {
				return err
			} else if res, err := q.CoordinatorState(
				context.Background(),
				&atomictypes.QueryCoordinatorStateRequest{
					TxId: txId,
				}); err != nil {
				return err
			} else if bz, err := ctx.Codec.MarshalJSON(&res.CoodinatorState); err != nil {
				return err
			} else if err := clientCtx.PrintString(string(bz)); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
