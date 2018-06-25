package transaction

import (
	"bytes"
	"giscoin/utils"
)

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, publicKeyHash) == 0
}
