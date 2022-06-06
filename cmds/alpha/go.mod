module github.com/datachainlab/fabric-tendermint-cross-demo/cmds/alpha

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.43.0-beta1
	github.com/cosmos/ibc-go v1.0.0-beta1
	github.com/datachainlab/cross v0.2.2
	github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20 v0.0.0
	github.com/hyperledger-labs/yui-relayer v0.1.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/tendermint/tendermint v0.34.10
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20 => ../../contracts/erc20
	github.com/datachainlab/fabric-tendermint-cross-demo/demo/chains/tendermint => ../../demo/chains/tendermint
	github.com/go-kit/kit => github.com/go-kit/kit v0.8.0
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
)
