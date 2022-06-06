module github.com/datachainlab/fabric-tendermint-cross-demo/contracts/erc20

go 1.16

require (
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/cosmos/cosmos-sdk v0.43.0-beta1
	github.com/datachainlab/cross v0.2.2
	github.com/datachainlab/cross-cdt v0.0.0-20211216051311-b41689652356
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/tendermint/tendermint v0.34.10
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/go-kit/kit => github.com/go-kit/kit v0.8.0
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
)
