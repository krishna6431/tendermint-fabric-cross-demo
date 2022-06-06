package keeper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	txtypes "github.com/datachainlab/cross/x/core/tx/types"

	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20contract/types"
	erc20mgrkeeper "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/keeper"
)

type Keeper struct {
	m              codec.Codec
	storeKey       sdk.StoreKey
	erc20mgrKeeper erc20mgrkeeper.Keeper
}

// NewKeeper creates a new keeper instance
func NewKeeper(m codec.Codec, key sdk.StoreKey, erc20mgrKeeper erc20mgrkeeper.Keeper) Keeper {
	return Keeper{m: m, storeKey: key, erc20mgrKeeper: erc20mgrKeeper}
}

// FIXME: HandleContractCall is called by ContractModule
func (k Keeper) HandleContractCall(goCtx context.Context, signers []authtypes.Account, callInfo txtypes.ContractCallInfo) (*txtypes.ContractCallResult, error) {
	var req types.ContractCallRequest
	if err := k.m.UnmarshalJSON(callInfo, &req); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if len(signers) == 0 {
		return nil, fmt.Errorf("can't get signer")
	}
	id := signers[0].Id.AccAddress()

	switch req.Method {
	case "transfer":
		amount, err := getAmount(req.Args[1])
		if err != nil {
			return nil, err
		}
		recipient, err := sdk.AccAddressFromBech32(req.Args[0])
		if err != nil {
			return nil, err
		}
		return k.erc20mgrKeeper.Transfer(ctx, id, recipient, amount)
	default:
		panic(fmt.Sprintf("unknown method '%v'", req.Method))
	}
}

func (k Keeper) Codec() codec.Codec {
	return k.m
}

func getAmount(target string) (int64, error) {
	return strconv.ParseInt(target, 10, 64)
}
