package cli

import (
	"fmt"
	"giscoin/wallets"
)

func (cli *CLI) printWallet() {
	_wallets, _ := wallets.NewWallets()

	fmt.Println(_wallets)
}
