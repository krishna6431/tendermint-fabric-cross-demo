package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	atomictypes "github.com/datachainlab/cross/x/core/atomic/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	initiatortypes "github.com/datachainlab/cross/x/core/initiator/types"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/hyperledger-labs/yui-fabric-ibc/app"
	"github.com/hyperledger-labs/yui-fabric-ibc/chaincode"
	"github.com/hyperledger-labs/yui-relayer/chains/fabric"
	abcitypes "github.com/tendermint/tendermint/abci/types"

	erc20mgrtypes "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
)

type Chain struct {
	*fabric.Chain
}

const (
	queryFunc = "query"
)

func (c *Chain) Query(req app.RequestQuery) (*app.ResponseQuery, error) {
	bz, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r, err := c.Contract().EvaluateTransaction(queryFunc, string(bz))
	if err != nil {
		return nil, err
	}
	var res app.ResponseQuery
	if err := json.Unmarshal(r, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Chain) QueryCoordinatorState(txId []byte) (*atomictypes.QueryCoordinatorStateResponse, error) {
	req := &atomictypes.QueryCoordinatorStateRequest{
		TxId: txId,
	}
	var res atomictypes.QueryCoordinatorStateResponse
	if err := c.query("/cross.core.atomic.Query/CoordinatorState", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Chain) QueryTxAuthState(txId []byte) (*authtypes.QueryTxAuthStateResponse, error) {
	req := &authtypes.QueryTxAuthStateRequest{
		TxID: txId,
	}
	var res authtypes.QueryTxAuthStateResponse
	if err := c.query("/cross.core.auth.Query/TxAuthState", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Chain) QuerySelfXCC() (*initiatortypes.QuerySelfXCCResponse, error) {
	req := &initiatortypes.QuerySelfXCCRequest{}

	var res initiatortypes.QuerySelfXCCResponse
	if err := c.query("/cross.core.initiator.Query/SelfXCC", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Chain) QueryIBCChannels() (*channeltypes.QueryChannelsResponse, error) {
	req := &channeltypes.QueryChannelsRequest{}
	var res channeltypes.QueryChannelsResponse
	if err := c.query("/ibc.core.channel.v1.Query/Channels", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Chain) QueryIBCClientStates() (*clienttypes.QueryClientStatesResponse, error) {
	req := &clienttypes.QueryClientStatesRequest{}
	var res clienttypes.QueryClientStatesResponse
	if err := c.query("/ibc.core.client.v1.Query/ClientStates", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Chain) QueryAllowance(owner sdk.AccAddress, spender sdk.AccAddress) (*erc20mgrtypes.QueryAllowanceResponse, error) {
	req := &erc20mgrtypes.QueryAllowanceRequest{Owner: owner, Spender: spender}
	var res erc20mgrtypes.QueryAllowanceResponse
	if err := c.query("/erc20mgr.Query/Allowance", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Chain) QueryBalanceOf(id sdk.AccAddress) (*erc20mgrtypes.QueryBalanceOfResponse, error) {
	req := &erc20mgrtypes.QueryBalanceOfRequest{Id: id}
	var res erc20mgrtypes.QueryBalanceOfResponse
	if err := c.query("/erc20mgr.Query/BalanceOf", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Chain) QueryTotalSupply() (*erc20mgrtypes.QueryTotalSupplyResponse, error) {
	var res erc20mgrtypes.QueryTotalSupplyResponse
	if err := c.query("/erc20mgr.Query/TotalSupply", &types.Empty{}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Chain) query(path string, req proto.Message, res interface{ Unmarshal(bz []byte) error }) error {
	bz, err := proto.Marshal(req)
	if err != nil {
		return err
	}
	r, err := c.Query(app.RequestQuery{
		Data: chaincode.EncodeToString(bz),
		Path: path,
	})
	if err != nil {
		return err
	}
	bz, err = chaincode.DecodeString(r.Value)
	if err != nil {
		return err
	}
	return res.Unmarshal(bz)
}

func (c *Chain) OutputTxIDFromEvent(res []byte) error {
	var txRes app.ResponseTx
	if err := json.Unmarshal(res, &txRes); err != nil {
		fmt.Println("can't unmarshal to ResponseTx")
		return err
	}
	var events []abcitypes.Event
	if err := json.Unmarshal([]byte(txRes.Events), &events); err != nil {
		fmt.Println("can't unmarshal Events")
		return err
	}

	exists := false
	for _, ev := range events {
		if ev.Type == "tx" {
			for _, attr := range ev.Attributes {
				if bytes.Equal(attr.Key, []byte("id")) {
					fmt.Printf("%s\n", string(attr.Value))
					exists = true
				}
			}
		}
	}

	if !exists {
		return fmt.Errorf("the id property does not exist: %v", events)
	}

	return nil
}
