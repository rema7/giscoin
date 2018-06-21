package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"giscoin/transaction"
	"time"
)

type Block struct {
	Timestamp     int64
	Transactions  []*transaction.Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(transactions []*transaction.Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		transactions,
		prevBlockHash,
		[]byte{},
		0,
	}

	return block
}

func NewGenesisBlock(coinbase *transaction.Transaction) *Block {
	return NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	encoder.Encode(b)

	return result.Bytes()
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (b Block) ToSting() {
	fmt.Printf("Previous hash: %x\n", b.PrevBlockHash)
	fmt.Printf("Data: %v\n", b.Transactions)
	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Println()
}
