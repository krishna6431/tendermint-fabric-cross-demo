package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cdttypes "github.com/datachainlab/cross-cdt/x/cdt/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
)

// RegisterInterfaces register the ibc transfer module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgContractCallTx{},
	)
	registry.RegisterImplementations(
		(*authtypes.AuthExtensionVerifier)(nil),
	)
	cdttypes.RegisterInterfaces(registry)
}

// ModuleCdc references the global erc20 module codec. Note, the codec
// should ONLY be used in certain instances of tests and for JSON encoding.
//
// The actual codec used for serialization should be provided to x/ibc-transfer and
// defined at the application level.
var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
