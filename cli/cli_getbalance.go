package cli

import (
	"fmt"
	"giscoin/blockchain"
	"giscoin/utils"
	"giscoin/utxo"
	"giscoin/wallets"
	"log"
)

func (cli *CLI) getBalance(address string) {
	if !wallets.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := blockchain.NewBlockchain(address)
	defer bc.DB.Close()
	utxoSet := utxo.UTXOSet{bc}

	balance := 0

	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := utxoSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
