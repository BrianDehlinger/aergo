package p2putil

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func ConvertPKToBTCEC(pk crypto.PrivKey) *btcec.PrivateKey {
	raw, err := pk.Raw()
	if err != nil {
		return nil
	}
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), raw)
	return priv
}

func ConvertPubKeyToBTCEC(pk crypto.PubKey) *btcec.PublicKey {
	raw, err := pk.Raw()
	if err != nil {
		return nil
	}
	pub, _ := btcec.ParsePubKey(raw, btcec.S256())
	return pub
}
func ConvertPKToLibP2P(pk  *btcec.PrivateKey) crypto.PrivKey {
	libp2pKey, err := crypto.UnmarshalSecp256k1PrivateKey(pk.Serialize())
	if err != nil {
		return nil
	}
	return libp2pKey
}


func ConvertPubToLibP2P(pk  *btcec.PublicKey) crypto.PubKey {
	libp2pKey, err := crypto.UnmarshalSecp256k1PublicKey(pk.SerializeCompressed())
	if err != nil {
		return nil
	}
	return libp2pKey
}