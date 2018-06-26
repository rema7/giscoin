package cli

import (
	"fmt"
	"giscoin/blockchain"
	"giscoin/utxo"
	"giscoin/wallets"
	"log"
)

func (cli *CLI) initBlockchain(address string) {
	if !wallets.ValidateAddress(address) {
		log.Panic("Invalid address")
	}
	bc := blockchain.InitBlockchain(address)
	defer bc.DB.Close()

	UTXOSet := utxo.UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}
