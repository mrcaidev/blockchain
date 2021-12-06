package cli

import (
	"blockchain/core/chain"
	"blockchain/transaction"
	"blockchain/utils"
	"blockchain/wallet"
	"fmt"
)

// 创建钱包。
func newWallet() {
	wallets := wallet.LoadWallets()
	address := wallets.AddWallet()
	wallets.Persist()

	fmt.Printf("New wallet address: %s\n", address)
}

// 列出地址。
func listAddresses() {
	wallets := wallet.LoadWallets()

	for index, address := range wallets.Addresses() {
		fmt.Printf("Address %d: %s\n", index, address)
	}
}

// 创建区块链。
func newChain(address string) {
	if !utils.IsValidAddress(address) {
		panic("invalid address")
	}

	chain := chain.NewChain(address)
	chain.Close()

	fmt.Println("New chain created.")
}

// 查询余额。
func queryBalance(address string) {
	if !utils.IsValidAddress(address) {
		panic("invalid address")
	}

	chain := chain.LoadChain()
	defer chain.Close()

	balance := chain.GetBalance(address)

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// 发起交易。
func startTrade(from string, to string, amount int) {
	if !utils.IsValidAddress(from) {
		panic("invalid address <from>")
	}
	if !utils.IsValidAddress(to) {
		panic("invalid address <to>")
	}

	chain := chain.LoadChain()
	defer chain.Close()

	tx := chain.NewUTXOTX(from, to, amount)
	chain.AddBlock([]*transaction.Transaction{tx})

	fmt.Println("Trade completed.")
}

// 打印区块链。
func printChain() {
	chain := chain.LoadChain()
	defer chain.Close()

	chain.Print()
}

// 显示帮助。
func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  wallet                                               Create a new wallet.")
	fmt.Println("  list                                                 List the addresses of all wallets.")
	fmt.Println("  chain      -address <address>                        Create a new blockchain mined out by <address>.")
	fmt.Println("  balance    -address <address>                        Query balance of <address>.")
	fmt.Println("  trade      -from <from> -to <to> -amount <amount>    Trade <amount> of coins from <from> to <to>.")
	fmt.Println("  print                                                Print blockchain information.")
	fmt.Println("  help                                                 Show help of commands.")
}
