package main

import (
	"giscoin/cli"
	"giscoin/network"
)

func main() {

	network.Start()
	cli2 := cli.CLI{}
	cli2.Run()
}
