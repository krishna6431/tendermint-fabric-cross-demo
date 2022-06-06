package cross

import (
	"context"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	initiatortypes "github.com/datachainlab/cross/x/core/initiator/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/config"
)

func GetCreateContractTransaction(ctx *config.Context) *cobra.Command {
	const (
		FlagSignerAddress = "signer-address"
		flagCallInfo      = "call-info"
		flagOutput        = "output"
	)

	cmd := &cobra.Command{
		Use:   "create-contract-tx",
		Short: "Create a new contract transaction",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Query self-XCC to chaincode
			q := initiatortypes.NewQueryClient(clientCtx)
			res, err := q.SelfXCC(
				context.Background(),
				&initiatortypes.QuerySelfXCCRequest{},
			)
			if err != nil {
				return err
			}
			var anyXCC *codectypes.Any = res.Xcc

			// signer
			signerAddress, err := cmd.Flags().GetString(FlagSignerAddress)
			if err != nil {
				return err
			}
			signerAccAddress, err := sdk.AccAddressFromBech32(signerAddress)
			if err != nil {
				return err
			}
			signerAccountID := authtypes.AccountIDFromAccAddress(signerAccAddress)

			signer := authtypes.Account{
				Id:       signerAccountID,
				AuthType: authtypes.NewAuthTypeLocal(),
			}

			callInfo := []byte(viper.GetString(flagCallInfo))
			cTx := initiatortypes.ContractTransaction{
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

	cmd.Flags().String(FlagSignerAddress, "", "signer address to sign to initiateTx")
	_ = cmd.MarkFlagRequired(FlagSignerAddress)

	cmd.Flags().String(flagCallInfo, "", "A contract call info")
	_ = cmd.MarkFlagRequired(flagCallInfo)

	cmd.Flags().String(flags.FlagOutputDocument, "", "The document will be written to the given file instead of STDOUT")
	cmd.Flags().StringP(flagOutput, "o", "text", "Output format (text|json)")

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
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			q := authtypes.NewQueryClient(clientCtx)

			if txId, err := hex.DecodeString(args[0]); err != nil {
				return err
			} else if res, err := q.TxAuthState(
				context.Background(),
				&authtypes.QueryTxAuthStateRequest{
					TxID: txId,
				}); err != nil {
				return err
			} else if bz, err := ctx.Codec.MarshalJSON(res); err != nil {
				return err
			} else {
				return clientCtx.PrintString(string(bz))
			}
		},
	}

	return cmd
}
