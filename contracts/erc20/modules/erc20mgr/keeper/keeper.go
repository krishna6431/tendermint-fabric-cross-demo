package keeper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	erc20keeper "github.com/datachainlab/cross-cdt/modules/erc20/keeper"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	txtypes "github.com/datachainlab/cross/x/core/tx/types"
	"github.com/pkg/errors"

	"github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
)

type Keeper struct {
	m           codec.Codec
	erc20Keeper erc20keeper.Keeper
	paramSpace  paramtypes.Subspace
}

// NewKeeper creates a new keeper instance
func NewKeeper(m codec.Codec, erc20Keeper erc20keeper.Keeper, paramSpace paramtypes.Subspace) Keeper {
	paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	return Keeper{m: m, erc20Keeper: erc20Keeper, paramSpace: paramSpace}
}

// HandleContractCall is called by ContractModule
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
	case "mint":
		return k.handleMint(ctx, id, req)
	case "burn":
		return k.handleBurn(ctx, id, req)
	case "transfer":
		return k.handleTransfer(ctx, id, req)
	case "approve":
		return k.handleApprove(ctx, id, req)
	case "transferFrom":
		return k.HandleTransferFrom(ctx, id, req)
	default:
		panic(fmt.Sprintf("unknown method '%v'", req.Method))
	}
}

func (k Keeper) handleMint(ctx sdk.Context, id sdk.AccAddress, req types.ContractCallRequest) (*txtypes.ContractCallResult, error) {
	admin, err := k.getAdmin(ctx)
	if err != nil {
		return nil, err
	}
	if !id.Equals(admin) {
		return nil, errors.Errorf("failed to authorize. id: %s, admin id: %s", id.String(), admin.String())
	}

	if len(req.Args) != 2 {
		return nil, errors.New("invalid argument")
	}

	account, err := sdk.AccAddressFromBech32(req.Args[0])
	if err != nil {
		return nil, err
	}
	amount, err := getAmount(req.Args[1])
	if err != nil {
		return nil, err
	}
	err = k.erc20Keeper.Mint(ctx, account.String(), amount)
	if err != nil {
		return nil, err
	}

	// TODO: set return values
	return &txtypes.ContractCallResult{Data: []byte("")}, nil
}

func (k Keeper) handleBurn(ctx sdk.Context, id sdk.AccAddress, req types.ContractCallRequest) (*txtypes.ContractCallResult, error) {
	admin, err := k.getAdmin(ctx)
	if err != nil {
		return nil, err
	}
	if !id.Equals(admin) {
		return nil, errors.Errorf("failed to authorize. id: %s", id)
	}

	if len(req.Args) != 2 {
		return nil, errors.New("invalid argument")
	}
	account, err := sdk.AccAddressFromBech32(req.Args[0])
	if err != nil {
		return nil, err
	}
	amount, err := getAmount(req.Args[1])
	if err != nil {
		return nil, err
	}
	err = k.erc20Keeper.Burn(ctx, account.String(), amount)
	if err != nil {
		return nil, err
	}

	// TODO: set return values
	return &txtypes.ContractCallResult{Data: []byte("")}, nil
}

func (k Keeper) handleTransfer(ctx sdk.Context, sender sdk.AccAddress, req types.ContractCallRequest) (*txtypes.ContractCallResult, error) {
	if len(req.Args) != 2 {
		return nil, errors.New("invalid argument")
	}
	recipient, err := sdk.AccAddressFromBech32(req.Args[0])
	if err != nil {
		return nil, err
	}
	amount, err := getAmount(req.Args[1])
	if err != nil {
		return nil, err
	}
	return k.Transfer(ctx, sender, recipient, amount)
}

func (k Keeper) Transfer(ctx sdk.Context, sender, recipient sdk.AccAddress, amount int64) (*txtypes.ContractCallResult, error) {
	err := k.erc20Keeper.Transfer(ctx, sender.String(), recipient.String(), amount)
	if err != nil {
		return nil, err
	}

	// TODO: set return values
	return &txtypes.ContractCallResult{Data: []byte("")}, nil
}

func (k Keeper) handleApprove(ctx sdk.Context, owner sdk.AccAddress, req types.ContractCallRequest) (*txtypes.ContractCallResult, error) {
	if len(req.Args) != 2 {
		return nil, errors.New("invalid argument")
	}
	spender, err := sdk.AccAddressFromBech32(req.Args[0])
	if err != nil {
		return nil, err
	}
	amount, err := getAmount(req.Args[1])
	if err != nil {
		return nil, err
	}
	err = k.erc20Keeper.Approve(ctx, owner.String(), spender.String(), amount)
	if err != nil {
		return nil, err
	}

	// TODO: set return values
	return &txtypes.ContractCallResult{Data: []byte("")}, nil
}

func (k Keeper) HandleTransferFrom(ctx sdk.Context, spender sdk.AccAddress, req types.ContractCallRequest) (*txtypes.ContractCallResult, error) {
	if len(req.Args) != 3 {
		return nil, errors.New("invalid argument")
	}
	owner, err := sdk.AccAddressFromBech32(req.Args[0])
	if err != nil {
		return nil, err
	}
	recipient, err := sdk.AccAddressFromBech32(req.Args[1])
	if err != nil {
		return nil, err
	}
	amount, err := getAmount(req.Args[2])
	if err != nil {
		return nil, err
	}
	err = k.erc20Keeper.TransferFrom(ctx, owner.String(), spender.String(), recipient.String(), amount)
	if err != nil {
		return nil, err
	}

	// TODO: set return values
	return &txtypes.ContractCallResult{Data: []byte("")}, nil
}

func (k Keeper) Codec() codec.Codec {
	return k.m
}

func (k Keeper) SetAdmin(ctx sdk.Context, admin sdk.AccAddress) error {
	var isSetAdmin bool
	k.paramSpace.GetIfExists(ctx, types.KeyIsSetAdmin, &isSetAdmin)
	if isSetAdmin {
		return errors.New("failed to set admin")
	}
	k.paramSpace.Set(ctx, types.KeyAdmin, admin.String())
	k.paramSpace.Set(ctx, types.KeyIsSetAdmin, true)

	return nil
}

func (k Keeper) getAdmin(ctx sdk.Context) (sdk.AccAddress, error) {
	var admin string
	k.paramSpace.GetIfExists(ctx, types.KeyAdmin, &admin)
	if admin == "" {
		return sdk.AccAddress{}, errors.New("failed to get admin")
	}

	acc, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	return acc, nil
}

func getAmount(target string) (int64, error) {
	return strconv.ParseInt(target, 10, 64)
}
