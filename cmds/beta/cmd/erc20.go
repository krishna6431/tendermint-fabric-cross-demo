package cmd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/datachainlab/cross/x/core/auth/types"
	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/flags"
	erc20mgrtypes "github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20/modules/erc20mgr/types"
)

func erc20Cmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "erc20",
		Short:                      "ERC20 subcommands",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		mintERC20Cmd(ctx),
		approveERC20Cmd(ctx),
		allowanceERC20Cmd(ctx),
		balanceOfERC20Cmd(ctx),
		totalSupplyERC20Cmd(ctx),
		transferERC20Cmd(ctx),
	)

	return cmd
}

func mintERC20Cmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint",
		Short: "Mint token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.Chain.Connect(); err != nil {
				return err
			}

			receiver, err := cmd.Flags().GetString(flags.FlagReceiverID)
			if err != nil {
				return err
			}

			amount, err := cmd.Flags().GetString(flags.FlagAmount)
			if err != nil {
				return err
			}

			msg := erc20mgrtypes.NewMsgContractCallTx(
				&erc20mgrtypes.ContractCallRequest{
					Method: "mint",
					Args:   []string{receiver, amount},
				},
				[]authtypes.AccountID{
					authtypes.AccountIDFromAccAddress(ctx.Address),
				},
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			_, err = ctx.Chain.SendMsgs([]sdk.Msg{msg})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flags.FlagReceiverID, "", "id")
	_ = cmd.MarkFlagRequired(flags.FlagReceiverID)
	cmd.Flags().String(flags.FlagAmount, "", "amount")
	_ = cmd.MarkFlagRequired(flags.FlagAmount)

	return cmd
}

func approveERC20Cmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.Chain.Connect(); err != nil {
				return err
			}

			spender, err := cmd.Flags().GetString(flags.FlagSpenderID)
			if err != nil {
				return err
			}
			amount, err := cmd.Flags().GetString(flags.FlagAmount)
			if err != nil {
				return err
			}

			msg := erc20mgrtypes.NewMsgContractCallTx(
				&erc20mgrtypes.ContractCallRequest{
					Method: "approve",
					Args:   []string{spender, amount},
				},
				[]authtypes.AccountID{authtypes.AccountIDFromAccAddress(ctx.Address)},
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			_, err = ctx.Chain.SendMsgs([]sdk.Msg{msg})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flags.FlagSpenderID, "", "id")
	_ = cmd.MarkFlagRequired(flags.FlagSpenderID)
	cmd.Flags().String(flags.FlagAmount, "", "amount")
	_ = cmd.MarkFlagRequired(flags.FlagAmount)

	return cmd
}

func allowanceERC20Cmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allowance",
		Short: "Get allowance",
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, err := cmd.Flags().GetString(flags.FlagOwnerID)
			if err != nil {
				return err
			}
			ownerAddr, err := sdk.AccAddressFromBech32(owner)
			if err != nil {
				return err
			}

			spender, err := cmd.Flags().GetString(flags.FlagSpenderID)
			if err != nil {
				return err
			}
			spenderAddr, err := sdk.AccAddressFromBech32(spender)
			if err != nil {
				return err
			}

			res, err := ctx.Chain.QueryAllowance(ownerAddr, spenderAddr)
			if err != nil {
				return err
			}

			j, err := ctx.Chain.Codec().MarshalJSON(res)
			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		},
	}

	cmd.Flags().String(flags.FlagOwnerID, "", "id")
	_ = cmd.MarkFlagRequired(flags.FlagOwnerID)
	cmd.Flags().String(flags.FlagSpenderID, "", "id")
	_ = cmd.MarkFlagRequired(flags.FlagSpenderID)

	return cmd
}

func balanceOfERC20Cmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance-of",
		Short: "Get balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, err := cmd.Flags().GetString(flags.FlagOwnerID)
			if err != nil {
				return err
			}
			ownerAddr, err := sdk.AccAddressFromBech32(owner)
			if err != nil {
				return err
			}

			res, err := ctx.Chain.QueryBalanceOf(ownerAddr)
			if err != nil {
				return err
			}

			fmt.Println(res.Balance.Amount)
			return nil
		},
	}

	cmd.Flags().String(flags.FlagOwnerID, "", "id")
	_ = cmd.MarkFlagRequired(flags.FlagOwnerID)

	return cmd
}

func totalSupplyERC20Cmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-supply",
		Short: "Get totalSupply",
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := ctx.Chain.QueryTotalSupply()
			if err != nil {
				return err
			}

			fmt.Println(res.TotalSupply)
			return nil
		},
	}

	return cmd
}

func transferERC20Cmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Transfer token from owner account to recipient",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ctx.Chain.Connect(); err != nil {
				return err
			}

			receiver, err := cmd.Flags().GetString(flags.FlagReceiverID)
			if err != nil {
				return err
			}
			amount, err := cmd.Flags().GetString(flags.FlagAmount)
			if err != nil {
				return err
			}

			msg := erc20mgrtypes.NewMsgContractCallTx(
				&erc20mgrtypes.ContractCallRequest{
					Method: "transfer",
					Args:   []string{receiver, amount},
				},
				[]authtypes.AccountID{
					authtypes.AccountIDFromAccAddress(ctx.Address),
				},
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			_, err = ctx.Chain.SendMsgs([]sdk.Msg{msg})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flags.FlagReceiverID, "", "id")
	_ = cmd.MarkFlagRequired(flags.FlagReceiverID)
	cmd.Flags().String(flags.FlagAmount, "", "amount")
	_ = cmd.MarkFlagRequired(flags.FlagAmount)

	return cmd
}
