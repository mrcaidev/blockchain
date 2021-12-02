package main

import (
	"blockchain/blockchain"
	"fmt"
	"strconv"
)

func main() {
	chain := blockchain.CreateBlockChain()

	chain.AddBlock("Send 1 BTC to Ivan")
	chain.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range chain.Blocks {
		fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.CreatePoW(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
