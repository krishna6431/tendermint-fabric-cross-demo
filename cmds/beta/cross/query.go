package cross

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	"github.com/datachainlab/cross/x/core/initiator/types"
	xcctypes "github.com/datachainlab/cross/x/core/xcc/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
)

func GetCreateContractTransaction(ctx *config.Context) *cobra.Command {
	const (
		flagSignerID              = "signer-id"
		flagCallInfo              = "call-info"
		flagInitiatorChainChannel = "initiator-chain-channel"
		flagOutput                = "output"
	)

	cmd := &cobra.Command{
		Use:   "create-contract-tx",
		Short: "Create a new contract transaction",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var anyXCC *codectypes.Any

			initiatorChannel := viper.GetString(flagInitiatorChainChannel)
			ci, err := parseChannelInfoFromString(initiatorChannel)
			if err != nil {
				return err
			}
			anyXCC, err = xcctypes.PackCrossChainChannel(ci)
			if err != nil {
				return err
			}

			id := viper.GetString(flagSignerID)
			idAccAddress, err := sdk.AccAddressFromBech32(id)
			if err != nil {
				return err
			}
			signer := authtypes.Account{
				Id:       authtypes.AccountIDFromAccAddress(idAccAddress),
				AuthType: authtypes.NewAuthTypeChannelWithAny(anyXCC),
			}

			callInfo := []byte(viper.GetString(flagCallInfo))
			cTx := types.ContractTransaction{
				CrossChainChannel: anyXCC,
				Signers:           []authtypes.Account{signer},
				CallInfo:          callInfo,
			}
			// prepare output document
			closeFunc, err := setOutputFile(cmd)
			if err != nil {
				return err
			}
			defer closeFunc()

			bz, err := ctx.Codec.MarshalJSON(&cTx)
			if err != nil {
				return err
			}

			if _, err := cmd.OutOrStdout().Write(bz); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flagSignerID, "", "ID to sign to initiateTx")
	_ = cmd.MarkFlagRequired(flagSignerID)

	cmd.Flags().String(flagCallInfo, "", "A contract call info")
	_ = cmd.MarkFlagRequired(flagCallInfo)

	cmd.Flags().String(flagInitiatorChainChannel, "", "The channel info of the counterparty: '<channelID>:<portID>'")
	_ = cmd.MarkFlagRequired(flagInitiatorChainChannel)

	cmd.Flags().String(flags.FlagOutputDocument, "", "The document will be written to the given file instead of STDOUT")
	cmd.Flags().StringP(flagOutput, "o", "text", "Output format (text|json)")

	ethSignKeyFlag(cmd)

	return cmd
}

func setOutputFile(cmd *cobra.Command) (func(), error) {
	outputDoc, err := cmd.Flags().GetString(flags.FlagOutputDocument)
	if err != nil {
		return nil, err
	}
	if outputDoc == "" {
		cmd.SetOut(cmd.OutOrStdout())
		return nil, nil
	}

	dir := filepath.Dir(outputDoc)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
			return nil, err
		}
	}

	fp, err := os.OpenFile(outputDoc, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}

	cmd.SetOut(fp)

	return func() { fp.Close() }, nil
}

func QueryTxAuthStateCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx-auth-state [tx-id]",
		Short: "Query the state of a client in a given path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if txId, err := hex.DecodeString(args[0]); err != nil {
				return err
			} else if res, err := ctx.Chain.QueryTxAuthState(txId); err != nil {
				return err
			} else if bz, err := ctx.Chain.Codec().MarshalJSON(res); err != nil {
				return err
			} else {
				fmt.Println(string(bz))
			}

			return nil
		},
	}

	return cmd
}
