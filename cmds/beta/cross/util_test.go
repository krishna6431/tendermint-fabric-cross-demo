package cross

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// "mnemonic": "math razor capable expose worth grape metal sunset metal sudden usage scheme"
// "path": "m/44'/60'/0'/0/0"
// "ethereum-address" : "0xa89F47C6b463f74d87572b058427dA0A13ec5425"
var hexString = "e517af47112e4f501afb26e4f34eadc8b0ad8eadaf4962169fc04bc8ddbfe091"

func TestHexToSecp256k1PrivKey(t *testing.T) {
	// geth *ecdsa.PrivateKey
	privKey, err := crypto.HexToECDSA(hexString)
	require.NoError(t, err)

	address := crypto.PubkeyToAddress(privKey.PublicKey).Bytes()

	// cosmos *secp256k1.PrivKey
	privKey2, err := hexToSecp256k1PrivKey(hexString)
	require.NoError(t, err)

	// secp256k1.PubKey -> ecdsa.PublicKey
	pubKey2, err := crypto.DecompressPubkey(privKey2.PubKey().Bytes())
	require.NoError(t, err)

	address2 := crypto.PubkeyToAddress(*pubKey2).Bytes()

	require.Equal(t, address, address2)
}
