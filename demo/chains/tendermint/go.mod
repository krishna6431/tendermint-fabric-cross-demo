module github.com/datachainlab/fabric-tendermint-cross-demo/demo/chains/tendermint

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.43.0-beta1
	github.com/cosmos/ibc-go v1.0.0-beta1
	github.com/datachainlab/cross v0.2.2
	github.com/datachainlab/cross-cdt v0.0.0-20211216051311-b41689652356
	github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20 v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/hyperledger-labs/yui-fabric-ibc v0.2.1-0.20220124085331-d9981e90b43b
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.10
	github.com/tendermint/tm-db v0.6.4
)

replace (
	github.com/cosmos/ibc-go => github.com/datachainlab/ibc-go v0.0.0-20210706141244-07dc9d32d9e7
	github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20 => ./../../../contracts/erc20
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)
