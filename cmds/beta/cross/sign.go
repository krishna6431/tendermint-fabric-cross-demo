package cross

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	tenderminttypes "github.com/cosmos/ibc-go/modules/light-clients/07-tendermint/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	xcctypes "github.com/datachainlab/cross/x/core/xcc/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
)

type IBCSignTxer interface {
	resolveXCC(initiatorChannel string) (*codectypes.Any, error)
	getSignerID() (authtypes.AccountID, error)
	sendTx(msg *authtypes.MsgIBCSignTx) error
}

// NewIBCSignTxCmd sends NewMsgIBCSignTx transaction for tendermint/fabric
// TODO: switch tendermint/fabric by what?
func NewIBCSignTxCmd(ctx *config.Context) *cobra.Command {
	const (
		flagTxID                  = "tx-id"
		flagInitiatorChainChannel = "initiator-chain-channel"
	)

	cmd := &cobra.Command{
		Use:   "ibc-signtx",
		Short: "Sign the cross-chain transaction on other chain via the chain",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// fabric
			txer, err := NewFabricIBCSignTx(ctx)
			if err != nil {
				return err
			}

			anyXCC, err := txer.resolveXCC(viper.GetString(flagInitiatorChainChannel))
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}

			signer, err := txer.getSignerID()
			if err != nil {
				return err
			}

			txID, err := hex.DecodeString(viper.GetString(flagTxID))
			if err != nil {
				return err
			}

			res, err := ctx.Chain.QueryIBCClientStates()
			if err != nil {
				return err
			}
			if res.GetClientStates().Len() != 1 {
				return errors.New("ClientStates length is invalid")
			}
			var counterparty tenderminttypes.ClientState
			err = ctx.Codec.Unmarshal(res.ClientStates[0].ClientState.Value, &counterparty)
			if err != nil {
				return err
			}
			height := counterparty.GetLatestHeight()
			msg := authtypes.NewMsgIBCSignTx(
				anyXCC,
				txID,
				[]authtypes.AccountID{signer},
				clienttypes.Height{
					RevisionNumber: height.GetRevisionNumber(),
					RevisionHeight: height.GetRevisionHeight() + 100000,
				},
				0,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return txer.sendTx(msg)
		},
	}
	cmd.Flags().String(flagTxID, "", "hex encoding of the TxID")
	cmd.Flags().String(flagInitiatorChainChannel, "", "channel info: '<channelID>:<portID>'")
	cmd.MarkFlagRequired(flagTxID)
	cmd.MarkFlagRequired(flagInitiatorChainChannel)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseChannelInfoFromString(s string) (*xcctypes.ChannelInfo, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil, errors.New("channel format must be follow a format: '<channelID>:<portID>'")
	}
	return &xcctypes.ChannelInfo{Channel: parts[0], Port: parts[1]}, nil
}
