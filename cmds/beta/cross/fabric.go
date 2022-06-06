package cross

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	xcctypes "github.com/datachainlab/cross/x/core/xcc/types"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
)

type fabricIBCSignTx struct {
	ctx *config.Context
}

func NewFabricIBCSignTx(ctx *config.Context) (IBCSignTxer, error) {
	return &fabricIBCSignTx{
		ctx: ctx,
	}, nil
}

func (f *fabricIBCSignTx) resolveXCC(initiatorChannel string) (*codectypes.Any, error) {
	if initiatorChannel == "" {
		// Query self-XCC to chaincode
		res, err := f.ctx.Chain.QuerySelfXCC()
		if err != nil {
			return nil, err
		}
		return res.Xcc, nil
	}
	ci, err := parseChannelInfoFromString(initiatorChannel)
	if err != nil {
		return nil, err
	}
	return xcctypes.PackCrossChainChannel(ci)
}

// getSignerID returns authtypes.AccountID
func (f *fabricIBCSignTx) getSignerID() (authtypes.AccountID, error) {
	return authtypes.AccountIDFromAccAddress(f.ctx.Address), nil
}

func (f *fabricIBCSignTx) sendTx(msg *authtypes.MsgIBCSignTx) error {
	err := f.ctx.Chain.Connect()
	if err != nil {
		return err
	}
	_, err = f.ctx.Chain.SendMsgs([]sdk.Msg{msg})
	return err
}
