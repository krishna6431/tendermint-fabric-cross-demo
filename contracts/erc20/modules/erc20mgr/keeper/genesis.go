package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
)

// InitGenesis initializes Params of genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	admin, err := sdk.AccAddressFromBech32(state.Params.Admin)
	if err != nil {
		panic(err)
	}
	err = k.SetAdmin(ctx, admin)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis exports empty genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{}
}
