package cli

import (
	"blockchain/core/blockchain"
	"blockchain/core/transaction"
	"blockchain/core/wallet"
	"blockchain/utils"
	"bytes"
	"fmt"
)

// 判断是否是有效地址。
func isValidAddress(address string) bool {
	pubkeyHash := utils.Base58Decode([]byte(address))
	actualChecksum := pubkeyHash[len(pubkeyHash)-utils.ChecksumLen:]
	supposedChecksum := utils.GetChecksum(pubkeyHash[:len(pubkeyHash)-utils.ChecksumLen])
	return bytes.Equal(actualChecksum, supposedChecksum)
}

// 创建钱包。
func newWallet() {
	wallets := wallet.LoadWallets()

	address := wallets.AddWallet()
	wallets.Persist()

	fmt.Printf("New wallet created: %s\n", address)
}

// 列出地址。
func listAddresses() {
	wallets := wallet.LoadWallets()

	for index, address := range wallets.Addresses() {
		fmt.Printf("Wallet %d address: %s\n", index, address)
	}
}

// 创建区块链。
func newChain(address string) {
	if !isValidAddress(address) {
		panic("invalid address")
	}

	chain := blockchain.NewChain(address)
	defer chain.Close()

	chain.Reindex()

	fmt.Println("New chain created.")
}

// 查询余额。
func queryBalance(address string) {
	if !isValidAddress(address) {
		panic("invalid address")
	}

	chain := blockchain.LoadChain()
	defer chain.Close()

	balance := chain.GetBalance(address)

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// 发起交易。
func startTrade(from string, to string, amount int) {
	if !isValidAddress(from) {
		panic("invalid address <from>")
	}
	if !isValidAddress(to) {
		panic("invalid address <to>")
	}
	if from == to {
		panic("invalid trade cycle")
	}

	chain := blockchain.LoadChain()
	defer chain.Close()

	tx := chain.NewUtxoTx(from, to, amount)
	coinbaseTx := blockchain.NewCoinbaseTx(from, "")
	chain.AddBlock([]*transaction.Transaction{coinbaseTx, tx})

	fmt.Println("Trade completed.")
}

// 重新索引区块链。
func reindexChain() {
	chain := blockchain.LoadChain()
	defer chain.Close()

	chain.Reindex()
	cnt := chain.CountTx()

	fmt.Printf("Reindex completed: %d transactions in chain.", cnt)
}

// 打印区块链。
func printChain() {
	chain := blockchain.LoadChain()
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
	fmt.Println("  reindex                                              Reindex the transactions in chain.")
	fmt.Println("  print                                                Print blockchain information.")
	fmt.Println("  help                                                 Show help of commands.")
}
