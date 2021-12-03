package main

import (
	"blockchain/blockchain"
)

func main() {
	chain := blockchain.CreateBlockChain()
	chain.AddBlock("Send 1 BTC to Ivan")
	chain.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range chain.Blocks {
		block.Print()
	}
}
