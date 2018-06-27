package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	initChainCmd := flag.NewFlagSet("init", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	printWalletCmd := flag.NewFlagSet("printwallet", flag.ExitOnError)

	initChainData := initChainCmd.String("address", "", "Address")
	getBalanceData := getBalanceCmd.String("address", "", "Address")
	sendFrom := sendCmd.String("from", "", "Source wallets address")
	sendTo := sendCmd.String("to", "", "Destination wallets address")
	sendAmount := sendCmd.Int("amount", 0, "Amount of money")

	switch os.Args[1] {
	case "init":
		err := initChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printwallet":
		err := printWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			os.Exit(1)
		}
		cli.getBalance(*getBalanceData)
	}

	if initChainCmd.Parsed() {
		if *initChainData == "" {
			os.Exit(1)
		}
		cli.initBlockchain(*initChainData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if printWalletCmd.Parsed() {
		cli.printWallet()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

}
