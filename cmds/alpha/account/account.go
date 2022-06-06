package account

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
)

func GetAccountIDFromAccAddress(addr sdk.AccAddress) authtypes.AccountID {
	return authtypes.AccountIDFromAccAddress(addr)
}

func GetAccountIDFromCtx(clientCtx client.Context) authtypes.AccountID {
	return authtypes.AccountIDFromAccAddress(clientCtx.GetFromAddress())
}

func GetAccountFromCtx(clientCtx client.Context) authtypes.Account {
	account := authtypes.AccountIDFromAccAddress(clientCtx.GetFromAddress())
	return authtypes.NewAccount(account, authtypes.NewAuthTypeLocal())
}
