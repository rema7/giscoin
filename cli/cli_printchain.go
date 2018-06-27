package cli

import (
	"fmt"
	"giscoin/blockchain"
	"giscoin/concensus"
	"strconv"
)

func (cli *CLI) printChain() {
	bc := blockchain.NewBlockchain("")
	defer bc.DB.Close()
	bci := bc.Iterator()

	for {
		block := bci.Next()
		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := concensus.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
