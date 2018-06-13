package blockchain

import (
	"giscoin/block"
)

type Blockchain struct {
	Blocks []*block.Block
}

func (bc *Blockchain) AddBlock(newBlock *block.Block) {
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock() *block.Block {
	return block.NewBlock("Genesis block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*block.Block{NewGenesisBlock()}}
}
