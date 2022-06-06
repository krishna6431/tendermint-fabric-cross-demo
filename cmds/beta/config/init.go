package config

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/hyperledger-labs/yui-relayer/chains/fabric"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ConfigPath   = "config"
	ConfigFile   = "config.json"
	WalletPath   = "wallet"
	WalletSuffix = ".id"
)

func InitConfigCmd(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			home, err := cmd.Flags().GetString(flags.FlagHome)
			if err != nil {
				return err
			}

			cfgDir := path.Join(home, ConfigPath)
			cfgPath := path.Join(cfgDir, ConfigFile)

			// Check if config path is not exist to initialize
			if _, err := os.Stat(cfgPath); err == nil {
				return errors.New("file already exists")
			} else if !errors.Is(err, fs.ErrNotExist) {
				return err
			} else {
				if err = os.MkdirAll(cfgDir, os.ModePerm); err != nil {
					return err
				}

				cf := viper.GetString(flagConfigFile)
				certPath := viper.GetString(flagFabClientCertPath)
				privKeyPath := viper.GetString(flagFabClientPrivateKeyPath)
				cConfig := fabric.ChainConfig{}

				if bz, err := ioutil.ReadFile(cf); err != nil {
					return err
				} else if err := ctx.Codec.UnmarshalJSON(bz, &cConfig); err != nil {
					return err
				} else if iChain, err := cConfig.Build(); err != nil {
					return err
				} else if chain, ok := iChain.(*fabric.Chain); !ok {
					return fmt.Errorf("unrecognized chain type: %T", iChain)
				} else if err := chain.Init(home, 0, ctx.Codec, false); err != nil {
					return err
				} else if err := chain.PopulateWallet(certPath, privKeyPath); err != nil {
					return err
				}

				f, err := os.Create(cfgPath)
				if err != nil {
					return err
				}
				defer f.Close()

				if bz, err := ctx.Codec.MarshalJSON(&cConfig); err != nil {
					return err
				} else if _, err = f.Write(bz); err != nil {
					return err
				}

				return nil
			}
		},
	}

	configFileFlag(cmd)
	populateWalletFlag(cmd)

	return cmd
}
