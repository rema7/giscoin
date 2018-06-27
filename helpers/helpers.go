package helpers

import (
	"encoding/hex"
	"giscoin/blockchain"
	"giscoin/transaction"
	"giscoin/utxo"
	"giscoin/wallets"
	"log"
)

func NewUTXOTransaction(from, to string, amount int, bc *blockchain.Blockchain, utxoSet utxo.UTXOSet) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	_wallets, err := wallets.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := _wallets.GetWallet(from)
	pubKeyHash := wallets.HashPubKey(wallet.PublicKey)
	acc, validOutputs := utxoSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("EROR: Not enough funds")
	}

	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := transaction.TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *transaction.NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *transaction.NewTXOutput(acc-amount, from))
	}

	tx := transaction.Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}
