package cmd

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/account"
	alphaflags "github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha/flags"
	erc20mgrtypes "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
)

func erc20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "erc20",
		Short:                      "ERC20 subcommands",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		mintERC20Cmd(),
		approveERC20Cmd(),
		allowanceERC20Cmd(),
		balanceOfERC20Cmd(),
		totalSupplyERC20Cmd(),
		transferERC20Cmd(),
	)

	return cmd
}

// bin mint --receiver-address --amount
func mintERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint",
		Short: "Mint token",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			clientCtx = clientCtx.WithOutputFormat("json")

			// self account
			selfAccountID := account.GetAccountIDFromCtx(clientCtx)

			// receiver
			receiverAddress, err := cmd.Flags().GetString(alphaflags.FlagReceiverAddress)
			if err != nil {
				return err
			}
			receiver, err := sdk.AccAddressFromBech32(receiverAddress)
			if err != nil {
				return err
			}

			// amount
			amount, err := cmd.Flags().GetString(alphaflags.FlagAmount)
			if err != nil {
				return err
			}

			msg := erc20mgrtypes.NewMsgContractCallTx(
				&erc20mgrtypes.ContractCallRequest{
					Method: "mint",
					Args:   []string{receiver.String(), amount},
				},
				[]authtypes.AccountID{
					selfAccountID,
				},
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(alphaflags.FlagReceiverAddress, "", "receiver address")
	_ = cmd.MarkFlagRequired(alphaflags.FlagReceiverAddress)
	cmd.Flags().String(alphaflags.FlagAmount, "", "amount")
	_ = cmd.MarkFlagRequired(alphaflags.FlagAmount)
	alphaflags.AddDefaultFlagsToCmd(cmd)

	return cmd
}

// bin approve --spender-address --amount
func approveERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve token",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			clientCtx = clientCtx.WithOutputFormat("json")

			// self account
			selfAccountID := account.GetAccountIDFromCtx(clientCtx)

			// spender
			spenderAddress, err := cmd.Flags().GetString(alphaflags.FlagSpenderAddress)
			if err != nil {
				return err
			}
			spender, err := sdk.AccAddressFromBech32(spenderAddress)
			if err != nil {
				return err
			}

			// amount
			amount, err := cmd.Flags().GetString(alphaflags.FlagAmount)
			if err != nil {
				return err
			}

			msg := erc20mgrtypes.NewMsgContractCallTx(
				&erc20mgrtypes.ContractCallRequest{
					Method: "approve",
					Args:   []string{spender.String(), amount},
				},
				[]authtypes.AccountID{selfAccountID},
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(alphaflags.FlagSpenderAddress, "", "spender address")
	_ = cmd.MarkFlagRequired(alphaflags.FlagSpenderAddress)
	cmd.Flags().String(alphaflags.FlagAmount, "", "amount")
	_ = cmd.MarkFlagRequired(alphaflags.FlagAmount)
	alphaflags.AddDefaultFlagsToCmd(cmd)

	return cmd
}

// bin allowance --owner-address --spender-address
func allowanceERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allowance",
		Short: "Get allowance",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// owner
			ownerAddress, err := cmd.Flags().GetString(alphaflags.FlagOwnerAddress)
			if err != nil {
				return err
			}
			owner, err := sdk.AccAddressFromBech32(ownerAddress)
			if err != nil {
				return err
			}

			// spender
			spenderAddress, err := cmd.Flags().GetString(alphaflags.FlagSpenderAddress)
			if err != nil {
				return err
			}
			spender, err := sdk.AccAddressFromBech32(spenderAddress)
			if err != nil {
				return err
			}

			// query
			queryClient := erc20mgrtypes.NewQueryClient(clientCtx)
			res, err := queryClient.Allowance(
				context.Background(),
				&erc20mgrtypes.QueryAllowanceRequest{
					Owner:   owner,
					Spender: spender,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintString(strconv.FormatInt(res.Amount, 10))
		},
	}

	cmd.Flags().String(alphaflags.FlagOwnerAddress, "", "owner address")
	_ = cmd.MarkFlagRequired(alphaflags.FlagOwnerAddress)
	cmd.Flags().String(alphaflags.FlagSpenderAddress, "", "spender address")
	_ = cmd.MarkFlagRequired(alphaflags.FlagSpenderAddress)
	alphaflags.AddDefaultFlagsToCmd(cmd)

	return cmd
}

// bin balanceOf --owner-address
func balanceOfERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance-of",
		Short: "Get balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// owner
			ownerAddress, err := cmd.Flags().GetString(alphaflags.FlagOwnerAddress)
			if err != nil {
				return err
			}
			owner, err := sdk.AccAddressFromBech32(ownerAddress)
			if err != nil {
				return err
			}

			// query
			queryClient := erc20mgrtypes.NewQueryClient(clientCtx)
			res, err := queryClient.BalanceOf(
				context.Background(),
				&erc20mgrtypes.QueryBalanceOfRequest{
					Id: owner,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(strconv.FormatInt(res.Balance.Amount, 10))
		},
	}

	cmd.Flags().String(alphaflags.FlagOwnerAddress, "", "owner address")
	_ = cmd.MarkFlagRequired(alphaflags.FlagOwnerAddress)
	alphaflags.AddDefaultFlagsToCmd(cmd)

	return cmd
}

// bin totalSupply
func totalSupplyERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-supply",
		Short: "Get totalSupply",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// query
			queryClient := erc20mgrtypes.NewQueryClient(clientCtx)
			res, err := queryClient.TotalSupply(
				context.Background(),
				&emptypb.Empty{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(strconv.FormatInt(res.TotalSupply, 10))
		},
	}
	alphaflags.AddDefaultFlagsToCmd(cmd)

	return cmd
}

func transferERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Transfer token from owner account to recipient",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			clientCtx = clientCtx.WithOutputFormat("json")

			// self account
			selfAccountID := account.GetAccountIDFromCtx(clientCtx)

			// receiver
			receiverAddress, err := cmd.Flags().GetString(alphaflags.FlagReceiverAddress)
			if err != nil {
				return err
			}
			receiver, err := sdk.AccAddressFromBech32(receiverAddress)
			if err != nil {
				return err
			}

			// amount
			amount, err := cmd.Flags().GetString(alphaflags.FlagAmount)
			if err != nil {
				return err
			}

			msg := erc20mgrtypes.NewMsgContractCallTx(
				&erc20mgrtypes.ContractCallRequest{
					Method: "transfer",
					Args:   []string{receiver.String(), amount},
				},
				[]authtypes.AccountID{
					selfAccountID,
				},
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(alphaflags.FlagReceiverAddress, "", "receiver address")
	_ = cmd.MarkFlagRequired(alphaflags.FlagReceiverAddress)
	cmd.Flags().String(alphaflags.FlagAmount, "", "amount")
	_ = cmd.MarkFlagRequired(alphaflags.FlagAmount)
	alphaflags.AddDefaultFlagsToCmd(cmd)

	return cmd
}
