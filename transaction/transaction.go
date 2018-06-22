package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"giscoin/utils"
	"giscoin/wallet"
	"log"
)

const subsidy = 10

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

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature string
	PubKey    []byte
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(pubKeyHash)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (in *TXInput) CanUnlockOutputWith(unblockingData string) bool {
	return in.Signature == unblockingData
}

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TXInput{[]byte{}, -1, nil, []byte(data)}
	txOut := NewTXOutput(subsidy, to)
	tx := &Transaction{nil, []TXInput{txIn}, []TXOutput{*txOut}}
	tx.SetID()

	return tx
}

func (tx *Transaction) ToString() {
	fmt.Printf("ID %s", hex.EncodeToString(tx.ID))
	fmt.Println()
	fmt.Println("Vin:")
	for _, vin := range tx.Vin {
		fmt.Printf("\tTxId: %s VoutIdx: %d Signature: %s", hex.EncodeToString(vin.Txid), vin.Vout, vin.Signature)
		fmt.Println()
	}
	fmt.Println("Vout:")
	for _, vout := range tx.Vout {
		fmt.Printf("\tValue: %d, Signature: %s", vout.Value, vout.PubKeyHash)
		fmt.Println()
	}
}
