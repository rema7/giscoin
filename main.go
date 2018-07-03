package main

import (
	"fmt"
	"giscoin/node"
	"giscoin/wallets"
	"github.com/urfave/cli"
	"log"
	"os"
)

var port int
var connectToNode string

var (
	portFlag = cli.IntFlag{
		Name:        "port",
		Usage:       "Port number",
		Value:       30300,
		Destination: &port,
	}
	connectTo = cli.StringFlag{
		Name:        "connectTo",
		Usage:       "enode address",
		Value:       "",
		Destination: &connectToNode,
	}
	printWallet = cli.Command{
		Name:  "printwallet",
		Usage: "Port number",
		Action: func(c *cli.Context) {
			_wallets, _ := wallets.NewWallets()
			fmt.Println(_wallets)
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		portFlag,
		connectTo,
	}
	app.Commands = []cli.Command{
		printWallet,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	srv := node.Start(port)
	if connectToNode != "" {
		err := node.ConnectToPeer(srv, connectToNode)
		if err != nil {
			log.Printf("Failed to connect to peer with err: %v", err)
			return
		}

		communicated := make(chan bool)
		go node.SubscribeToEvents(srv, communicated)

		// Sent and received message, stopping the server since work is done
		<-communicated
		srv.Stop()
	} else {
		select {}
	}
	//cli2 := cli.CLI{}
	//cli2.Run()
}
