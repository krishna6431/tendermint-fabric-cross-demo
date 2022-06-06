package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
)

// msg types
const (
	TypeContractCallTx = "ContractCallTx"
)

var _ sdk.Msg = (*MsgContractCallTx)(nil)

func NewMsgContractCallTx(
	request *ContractCallRequest,
	signers []authtypes.AccountID,
) *MsgContractCallTx {
	return &MsgContractCallTx{
		Request: request,
		Signers: signers,
	}
}

// Route implements sdk.Msg
func (MsgContractCallTx) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgContractCallTx) Type() string {
	return TypeContractCallTx
}

// ValidateBasic performs a basic check of the MsgContractCallTx fields.
func (msg MsgContractCallTx) ValidateBasic() error {
	if msg.Request == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing request")
	}
	if len(msg.Signers) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "missing signers")
	}
	return nil
}

// GetSignBytes implements sdk.Msg.
func (msg MsgContractCallTx) GetSignBytes() []byte {
	panic("not support amino")
}

// GetSigners implements sdk.Msg
// GetSigners returns the addresses that must sign the transaction.
// Addresses are returned in a deterministic order.
// Duplicate addresses will be omitted.
func (msg MsgContractCallTx) GetSigners() []sdk.AccAddress {
	seen := map[string]bool{}
	signers := []sdk.AccAddress{}

	for _, s := range msg.Signers {
		addr := s.AccAddress().String()
		if !seen[addr] {
			signers = append(signers, s.AccAddress())
			seen[addr] = true
		}
	}

	return signers
}

func (msg MsgContractCallTx) GetAccountIDs() []authtypes.AccountID {
	accountIDs := make([]authtypes.AccountID, 0, len(msg.Signers))
	return append(accountIDs, msg.Signers...)
}
