package app

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	clientkeeper "github.com/cosmos/ibc-go/modules/core/02-client/keeper"
	connectionkeeper "github.com/cosmos/ibc-go/modules/core/03-connection/keeper"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/modules/core/keeper"
	tenderminttypes "github.com/cosmos/ibc-go/modules/light-clients/07-tendermint/types"
	"github.com/hyperledger-labs/yui-fabric-ibc/commitment"
	authtypes "github.com/hyperledger-labs/yui-fabric-ibc/x/auth/types"
	fabrictypes "github.com/hyperledger-labs/yui-fabric-ibc/x/ibc/light-clients/xx-fabric/types"
)

func overrideIBCClientKeeper(k ibckeeper.Keeper, cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace, seqMgr commitment.SequenceManager) *ibckeeper.Keeper {
	clientKeeper := NewClientKeeper(k.ClientKeeper, seqMgr)
	k.ConnectionKeeper = connectionkeeper.NewKeeper(cdc, key, paramSpace, clientKeeper)
	return &k
}

var (
	_ connectiontypes.ClientKeeper = (*ClientKeeper)(nil)
	_ channeltypes.ClientKeeper    = (*ClientKeeper)(nil)
)

// ClientKeeper override `GetSelfConsensusState` and `ValidateSelfClient` in the keeper of ibc-client
// original method doesn't yet support a consensus state for general client
type ClientKeeper struct {
	clientkeeper.Keeper

	seqMgr commitment.SequenceManager
}

func NewClientKeeper(k clientkeeper.Keeper, seqMgr commitment.SequenceManager) ClientKeeper {
	return ClientKeeper{Keeper: k, seqMgr: seqMgr}
}

// GetSelfConsensusState introspects the (self) past historical info at a given height
// and returns the expected consensus state at that height.
// For now, can only retrieve self consensus states for the current version
func (k ClientKeeper) GetSelfConsensusState(ctx sdk.Context, height exported.Height) (exported.ConsensusState, bool) {
	seq, err := k.seqMgr.GetSequence(authtypes.StubFromContext(ctx), height.GetRevisionHeight())
	if err != nil {
		return nil, false
	}
	return &fabrictypes.ConsensusState{
		Timestamp: seq.Timestamp,
	}, true
}

func (k ClientKeeper) ValidateSelfClient(ctx sdk.Context, clientState exported.ClientState) error {
	switch cs := clientState.(type) {
	case *tenderminttypes.ClientState:
		return k.Keeper.ValidateSelfClient(ctx, cs)
	case *fabrictypes.ClientState:
		// just skip validation
		return nil
	default:
		return fmt.Errorf("unexpected client state type: %T", cs)
	}
}
