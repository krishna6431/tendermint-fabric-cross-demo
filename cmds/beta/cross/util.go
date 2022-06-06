package cross

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

func hexToSecp256k1PrivKey(hexString string) (*secp256k1.PrivKey, error) {
	bz, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	return hd.Secp256k1.Generate()(bz).(*secp256k1.PrivKey), nil
}
