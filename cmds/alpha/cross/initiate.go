package cross

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	"github.com/datachainlab/cross/x/core/initiator/types"
	txtypes "github.com/datachainlab/cross/x/core/tx/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/config"
	alphaflags "github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/flags"
)

type InitiateTxer interface {
	getSigner() (authtypes.Account, error)
	getHeight() (*clienttypes.Height, error)
	sendTxWithEvent(msg *types.MsgInitiateTx) error
}

// NewInitiateTxCmd sends a NewMsgInitiateTx transaction
func NewInitiateTxCmd(ctx *config.Context) *cobra.Command {
	const (
		flagContractTransactions = "contract-txs"
	)

	cmd := &cobra.Command{
		Use:   "create-initiate-tx",
		Short: "Create and submit a MsgInitiateTx transaction for a simple commit",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			clientCtx = clientCtx.WithOutputFormat("json")

			// tendermint
			txer, err := NewTendermintInitiateTx(clientCtx, cmd)
			if err != nil {
				return err
			}

			ctxs, err := readContractTransactions(clientCtx.JSONCodec, viper.GetStringSlice(flagContractTransactions))
			if err != nil {
				return err
			}

			signer, err := txer.getSigner()
			if err != nil {
				return err
			}

			msg := types.NewMsgInitiateTx(
				[]authtypes.Account{signer},
				clientCtx.ChainID,
				uint64(time.Now().Unix()),
				txtypes.COMMIT_PROTOCOL_SIMPLE,
				ctxs,
				clienttypes.Height{},
				0,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return txer.sendTxWithEvent(msg)
		},
	}

	cmd.Flags().StringSlice(flagContractTransactions, nil, "A file path to includes a contract transaction")
	_ = cmd.MarkFlagRequired(flagContractTransactions)

	alphaflags.AddDefaultFlagsToCmd(cmd)

	return cmd
}

func readContractTransactions(m codec.JSONCodec, pathList []string) ([]types.ContractTransaction, error) {
	var cTxs []types.ContractTransaction
	for _, path := range pathList {
		bz, err := ioutil.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
		var cTx types.ContractTransaction
		if err := m.UnmarshalJSON(bz, &cTx); err != nil {
			return nil, err
		}
		cTxs = append(cTxs, cTx)
	}
	return cTxs, nil
}
