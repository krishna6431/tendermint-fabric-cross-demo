package config

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hyperledger-labs/yui-relayer/chains/fabric"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

type Context struct {
	Chain    *Chain
	Codec    codec.ProtoCodecMarshaler
	Config   *fabric.ChainConfig
	ClientID *cid.ClientID
	Address  sdk.AccAddress
}
