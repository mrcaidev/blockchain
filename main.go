package main

import (
	"blockchain/blockchain"
)

func main() {
	chain := blockchain.CreateBlockChain()
	chain.AddBlock("Send 1 BTC to Ivan")
	chain.AddBlock("Send 2 more BTC to Ivan")

	iter := chain.Iterator()
	for {
		if len(iter.CurHash) == 0 {
			break
		}
		curBlock := iter.Next()
		curBlock.Print()
	}
}
