package main

import (
	"flag"
	"fmt"
	"giscoin/blockchain"
	"giscoin/concensus"
	"giscoin/transaction"
	"giscoin/utils"
	"giscoin/wallets"
	"log"
	"os"
	"strconv"
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

func (cli *CLI) initBlockchain(address string) {
	if !wallets.ValidateAddress(address) {
		log.Panic("Invalid address")
	}
	bc := blockchain.InitBlockchain(address)
	bc.Db.Close()
	fmt.Println("Done!")
}

func (cli *CLI) createWallet() {
	_wallets, _ := wallets.NewWallets()
	address := _wallets.CreateWallet()
	_wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}

func (cli *CLI) getBalance(address string) {
	bc := blockchain.NewBlockchain(address)
	defer bc.Db.Close()

	balance := 0

	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	bc := blockchain.NewBlockchain(from)
	defer bc.Db.Close()

	tx := bc.NewUTXOTransaction(from, to, amount)
	bc.MineBlock([]*transaction.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bc := blockchain.NewBlockchain("")
	defer bc.Db.Close()
	bci := bc.Iterator()

	for {
		block := bci.Next()

		block.ToSting()
		pow := concensus.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	initChainCmd := flag.NewFlagSet("init", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

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
	default:
		cli.printUsage()
		os.Exit(1)
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

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
