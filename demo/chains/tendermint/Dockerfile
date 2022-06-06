ARG GOLAGN_VER=1.15.8
ARG ALPINE_VER=3.13

FROM golang:${GOLAGN_VER}-alpine${ALPINE_VER} AS scaffold-initializer
RUN apk add --no-cache jq=1.6-r1 bash=5.1.16-r0

WORKDIR /root

COPY ./ ./

# First, build code to generate keys
RUN go build -v -mod=vendor -o /usr/bin/simd ./simapp/simd

# Generate seeds
ARG CHAINID
ARG CHAINDIR=./data
ARG RPCPORT=26657
ARG P2PPORT=26656
ARG PROFPORT=6060
ARG GRPCPORT=9090

# Second, build with address embedded to run tendermint chain
RUN ./scripts/tm-chain simd $CHAINID $CHAINDIR $RPCPORT $P2PPORT $PROFPORT $GRPCPORT && \
 ./scripts/build-go-embedded simd $CHAINID $CHAINDIR

FROM alpine:${ALPINE_VER}

WORKDIR /root

ARG CHAINID
ARG CHAINDIR=./data
ARG RPCPORT=26657
ARG P2PPORT=26656
ARG PROFPORT=6060
ARG GRPCPORT=9090

ENV CHAINID=$CHAINID
ENV CHAINDIR=$CHAINDIR
ENV GRPCPORT=$GRPCPORT

COPY --from=scaffold-initializer /usr/bin/simd /usr/bin/simd
COPY --from=scaffold-initializer /root/$CHAINDIR /root/$CHAINDIR
COPY ./scripts/entrypoint.sh /root/entrypoint.sh

EXPOSE $RPCPORT $P2PPORT $PROFPORT $GRPCPORT

ENTRYPOINT ["./entrypoint.sh"]
