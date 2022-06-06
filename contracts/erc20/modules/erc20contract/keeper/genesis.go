package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20contract/types"
)

// InitGenesis does nothing.
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
}

// ExportGenesis exports empty genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{}
}
