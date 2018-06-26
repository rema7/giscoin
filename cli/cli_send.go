package cli

import (
	"fmt"
	"giscoin/blockchain"
	"giscoin/helpers"
	"giscoin/transaction"
	"giscoin/utxo"
	"giscoin/wallets"
	"log"
)

func (cli *CLI) send(from, to string, amount int) {
	if !wallets.ValidateAddress(from) {
		log.Panic("Invalid from address")
	}
	if !wallets.ValidateAddress(to) {
		log.Panic("Invalid to address")
	}
	bc := blockchain.NewBlockchain(from)
	defer bc.DB.Close()
	UTXOSet := utxo.UTXOSet{bc}

	tx := helpers.NewUTXOTransaction(from, to, amount, bc, UTXOSet)
	newBlock := bc.MineBlock([]*transaction.Transaction{tx})

	UTXOSet.Update(newBlock)
	fmt.Println("Success!")
}
