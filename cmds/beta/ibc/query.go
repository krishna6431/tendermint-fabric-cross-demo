package ibc

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/datachainlab/fabric-tendermint-cross-demo/cmds/beta/config"
)

func QueryIBCChannelCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "Query the ChannelState of IBC",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if res, err := ctx.Chain.QueryIBCChannels(); err != nil {
				return err
			} else if len(res.Channels) == 0 {
				return errors.New("ibc channel does not exist")
			} else {
				channel := res.Channels[0]
				fmt.Printf("%s:%s\n", channel.ChannelId, channel.PortId)
				return nil
			}
		},
	}

	return cmd
}
