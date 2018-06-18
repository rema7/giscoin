package block

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		[]byte(data),
		prevBlockHash,
		[]byte{},
		0,
	}

	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	encoder.Encode(b)

	return result.Bytes()
}
