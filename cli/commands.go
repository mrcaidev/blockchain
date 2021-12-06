package cli

import (
	"blockchain/core/chain"
	tx "blockchain/transaction"
	"blockchain/utils"
	"blockchain/wallet"
	"fmt"
)

// 创建新区块链。
func newChain(address string) {
	if !utils.IsValidAddress(address) {
		panic("invalid address")
	}
	chain := chain.NewChain(address)
	chain.Close()
	fmt.Println("Created!")
}

// 创建新钱包。
func newWallet() {
	wallets := wallet.LoadWallets()
	address := wallets.AddWallet()
	wallets.Persist()
	fmt.Printf("New address: %s\n", address)
}

// 查询余额。
func balance(address string) {
	if !utils.IsValidAddress(address) {
		panic("invalid address")
	}
	chain := chain.LoadChain()
	defer chain.Close()

	balance := 0
	pubkeyHash := utils.Base58Decode([]byte(address))
	pubkeyHash = pubkeyHash[1 : len(pubkeyHash)-4]
	UTXOs := chain.FindUTXO(pubkeyHash)

	for _, utxo := range UTXOs {
		balance += utxo.Value
	}
	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// 列出地址列表。
func list() {
	wallets := wallet.LoadWallets()
	for index, addr := range wallets.Addresses() {
		fmt.Printf("Address %d: %s\n", index, addr)
	}
}

// 发送币。
func trade(from string, to string, amount int) {
	if !utils.IsValidAddress(from) {
		panic("invalid address <from>")
	}
	if !utils.IsValidAddress(to) {
		panic("invalid address <to>")
	}
	chain := chain.LoadChain()
	defer chain.Close()

	TX := chain.NewUTXOTX(from, to, amount)
	chain.AddBlock([]*tx.Transaction{TX})
	fmt.Println("Success!")
}

// 打印区块链。
func print() {
	chain := chain.LoadChain()
	defer chain.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()
		block.Print()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

// 帮助。
func help() {
	fmt.Println("Usage:")
	fmt.Println("  chain      -address <address>                        Create a new blockchain and reward subsidy to <address>.")
	fmt.Println("  wallet     -address <address>                        Create a new wallet at <address>.")
	fmt.Println("  list                                                 List the addresses of all wallets.")
	fmt.Println("  balance    -address <address>                        Query balance of <address> stored in the blockchain.")
	fmt.Println("  trade      -from <from> -to <to> -amount <amount>    Send <amount> of coins from <from> to <to>.")
	fmt.Println("  print                                                Print blocks info of the blockchain.")
	fmt.Println("  help                                                 Show help for commands.")
}
