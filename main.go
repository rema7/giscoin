package main

import (
	"fmt"
	"giscoin/block"
	"giscoin/blockchain"
	"giscoin/concensus"
	"strconv"
)

func GenerateBlock(bc *blockchain.Blockchain, data string) *block.Block {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := block.NewBlock(data, prevBlock.Hash)
	pow := concensus.NewProofOfWork(newBlock)
	nonce, hash := pow.Run()
	newBlock.Hash = hash
	newBlock.Nonce = nonce

	return newBlock
}

func main() {
	bc := blockchain.NewBlockchain()

	newBlock := GenerateBlock(bc, "First block")
	bc.AddBlock(newBlock)

	newBlock = GenerateBlock(bc, "Second block")
	bc.AddBlock(newBlock)

	for _, b := range bc.Blocks {
		fmt.Printf("Previous hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)
		pow := concensus.NewProofOfWork(b)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
