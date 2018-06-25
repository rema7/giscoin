package transaction

import (
	"bytes"
	"giscoin/wallets"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallets.HashPubKey(pubKeyHash)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
