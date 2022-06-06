package types

import (
	"context"

	cdttypes "github.com/datachainlab/cross-cdt/x/cdt/types"
	contracttypes "github.com/datachainlab/cross/x/core/contract/types"
	txtypes "github.com/datachainlab/cross/x/core/tx/types"
)

func CDTContractHandleDecorators() contracttypes.ContractHandleDecorator {
	return contracttypes.ContractHandleDecorators{
		CDTContractHandleDecorator{},
	}
}

var _ contracttypes.ContractHandleDecorator = (*CDTContractHandleDecorator)(nil)

type CDTContractHandleDecorator struct{}

func (cd CDTContractHandleDecorator) Handle(ctx context.Context, callInfo txtypes.ContractCallInfo) (newCtx context.Context, err error) {
	opmgr := cdttypes.NewOPManager()
	newCtx = cdttypes.ContextWithOPManager(ctx, opmgr)
	return newCtx, nil
}
