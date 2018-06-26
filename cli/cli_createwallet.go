package cli

import (
	"giscoin/wallets"
	"fmt"
)

func (cli *CLI) createWallet() {
	_wallets, _ := wallets.NewWallets()
	address := _wallets.CreateWallet()
	_wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}


