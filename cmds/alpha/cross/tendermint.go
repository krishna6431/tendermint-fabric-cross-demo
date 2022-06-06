package cross

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	ibctmtypes "github.com/cosmos/ibc-go/modules/light-clients/07-tendermint/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	"github.com/datachainlab/cross/x/core/initiator/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	tmtypes "github.com/tendermint/tendermint/types"
)

type tendermintInitiateTx struct {
	clientCtx client.Context
	flags     *flag.FlagSet
}

func NewTendermintInitiateTx(clientCtx client.Context, cmd *cobra.Command) (InitiateTxer, error) {
	return &tendermintInitiateTx{
		clientCtx: clientCtx,
		flags:     cmd.Flags(),
	}, nil
}

// getSigner returns authtypes.Account
func (t *tendermintInitiateTx) getSigner() (authtypes.Account, error) {
	sender := authtypes.AccountIDFromAccAddress(t.clientCtx.GetFromAddress())
	return authtypes.NewAccount(sender, authtypes.NewAuthTypeLocal()), nil
}

func (t *tendermintInitiateTx) getHeight() (*clienttypes.Height, error) {
	return getTendermintHeight(t.clientCtx)
}

func (t *tendermintInitiateTx) sendTxWithEvent(msg *types.MsgInitiateTx) error {
	return tx.GenerateOrBroadcastTxCLI(t.clientCtx, t.flags, msg)
}

func getTendermintHeight(clientCtx client.Context) (*clienttypes.Height, error) {
	h, height, err := queryTendermintHeader(clientCtx)
	if err != nil {
		return nil, err
	}
	version := clienttypes.ParseChainID(h.Header.ChainID)
	newHeight := clienttypes.NewHeight(version, uint64(height)+100)
	return &newHeight, nil
}

// QueryTendermintHeader takes a client context and returns the appropriate
// tendermint header
// Original implementation(but has a little) is here: https://github.com/cosmos/cosmos-sdk/blob/300b7393addba8c162cae929db90b083dcf93bd0/x/ibc/core/02-client/client/utils/utils.go#L123
func queryTendermintHeader(clientCtx client.Context) (ibctmtypes.Header, int64, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return ibctmtypes.Header{}, 0, err
	}

	info, err := node.ABCIInfo(context.Background())
	if err != nil {
		return ibctmtypes.Header{}, 0, err
	}

	height := info.Response.LastBlockHeight

	commit, err := node.Commit(context.Background(), &height)
	if err != nil {
		return ibctmtypes.Header{}, 0, err
	}

	page := 1
	count := 10_000

	validators, err := node.Validators(context.Background(), &height, &page, &count)
	if err != nil {
		return ibctmtypes.Header{}, 0, err
	}

	protoCommit := commit.SignedHeader.ToProto()
	protoValset, err := tmtypes.NewValidatorSet(validators.Validators).ToProto()
	if err != nil {
		return ibctmtypes.Header{}, 0, err
	}

	header := ibctmtypes.Header{
		SignedHeader: protoCommit,
		ValidatorSet: protoValset,
	}

	return header, height, nil
}
