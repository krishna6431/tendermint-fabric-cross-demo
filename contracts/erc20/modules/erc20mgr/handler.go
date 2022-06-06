package erc20mgr

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	contracttypes "github.com/datachainlab/cross/x/core/contract/types"

	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/keeper"
	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
)

// NewHandler returns sdk.Handler for IBC token transfer module messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		ctx = setupContractContext(ctx, msg)

		switch msg := msg.(type) {
		case *types.MsgContractCallTx:
			return handleContractCallTx(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized erc20mgr message type: %T", msg)
		}
	}
}

func handleContractCallTx(ctx sdk.Context, k keeper.Keeper, msg sdk.Msg) (*sdk.Result, error) {
	callInfo := msg.(*types.MsgContractCallTx).Request.ContractCallInfo(k.Codec())
	// this should be the same decorators to App
	chd := types.CDTContractHandleDecorators()
	goCtx, err := chd.Handle(ctx.Context(), callInfo)
	if err != nil {
		return nil, err
	}
	ctx = ctx.WithContext(goCtx)
	signers := []authtypes.Account{}
	for _, s := range msg.GetSigners() {
		signers = append(signers, authtypes.Account{
			Id: authtypes.AccountIDFromAccAddress(s),
		})
	}

	res, err := k.HandleContractCall(sdk.WrapSDKContext(ctx), signers, callInfo)
	return sdk.WrapServiceResult(ctx, res, err)
}

func setupContractContext(ctx sdk.Context, _ sdk.Msg) sdk.Context {
	return contracttypes.SetupContractContext(
		ctx,
		contracttypes.ContractRuntimeInfo{
			CommitMode:           contracttypes.BasicMode,
			ExternalCallResolver: nil,
		},
	)
}
