package main

import (
	"blockchain/blockchain"
	"blockchain/cli"
)

func main() {
	chain := blockchain.CreateBlockChain()
	defer chain.DB.Close()

	cli := cli.CreateCLI(chain)
	cli.Run()
}
