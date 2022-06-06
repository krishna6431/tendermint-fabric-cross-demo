package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
)

var _ types.QueryServer = (*Keeper)(nil)

func (k Keeper) BalanceOf(ctx context.Context, req *types.QueryBalanceOfRequest) (*types.QueryBalanceOfResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	v, err := k.erc20Keeper.BalanceOf(sdkCtx, req.Id.String())
	if err != nil {
		return nil, err
	}
	return &types.QueryBalanceOfResponse{
		Balance: &types.Balance{
			Id:     req.Id,
			Amount: v,
		},
	}, nil
}

func (k Keeper) TotalSupply(ctx context.Context, req *emptypb.Empty) (*types.QueryTotalSupplyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	v, err := k.erc20Keeper.TotalSupply(sdkCtx)
	if err != nil {
		return nil, err
	}
	return &types.QueryTotalSupplyResponse{TotalSupply: v}, nil
}

func (k Keeper) Allowance(ctx context.Context, req *types.QueryAllowanceRequest) (*types.QueryAllowanceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	v, err := k.erc20Keeper.Allowance(sdkCtx, req.Owner.String(), req.Spender.String())
	if err != nil {
		return nil, err
	}
	return &types.QueryAllowanceResponse{Amount: v}, nil
}
