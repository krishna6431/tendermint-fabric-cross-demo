ARG GO_VER=1.16
ARG ALPINE_VER=3.13

FROM golang:${GO_VER}-alpine${ALPINE_VER} as builder
RUN apk add --no-cache \
    gcc=10.2.1_pre1-r3 \
    libc-dev=0.7.2-r3

WORKDIR /root
COPY ./ ./

ARG ADDR
RUN go build -v \
-ldflags "-X main.ERC20Admin=${ADDR}" \
-mod=vendor -o /root/build/chaincode .

FROM alpine:${ALPINE_VER}

COPY ./entrypoint.sh /entrypoint.sh
COPY --from=builder /root/build/chaincode /usr/bin/chaincode

WORKDIR /root

CMD ["/entrypoint.sh"]
