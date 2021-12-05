package cli

import (
	"blockchain/core"
	tx "blockchain/transaction"
	"blockchain/utils"
	"blockchain/wallet"
	"fmt"
	"log"
)

// 创建新区块链。
func newChain(address string) {
	if !utils.ValidateAddress(address) {
		log.Panic("[Error] Invalid address!")
	}
	chain := core.NewBlockChain(address)
	chain.CloseDatabase()
	fmt.Println("Created!")
}

// 创建新钱包。
func newWallet() {
	wallets := wallet.NewWallets()
	address := wallets.AddWallet()
	wallets.Store()
	fmt.Printf("New address: %s\n", address)
}

// 查询余额。
func queryBalance(address string) {
	if !utils.ValidateAddress(address) {
		log.Panic("[Error] Invalid address!")
	}
	chain := core.LoadBlockChain()
	defer chain.CloseDatabase()

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
func listAddresses() {
	wallets := wallet.NewWallets()
	for index, addr := range wallets.Addresses() {
		fmt.Printf("Address %d: %s\n", index, addr)
	}
}

// 发送币。
func send(from string, to string, amount int) {
	if !utils.ValidateAddress(from) {
		log.Panic("[Error] Invalid <from>!")
	}
	if !utils.ValidateAddress(to) {
		log.Panic("[Error] Invalid <to>!")
	}
	chain := core.LoadBlockChain()
	defer chain.CloseDatabase()

	TX := chain.NewUTXOTX(from, to, amount)
	chain.AddBlock([]*tx.Transaction{TX})
	fmt.Println("Success!")
}

// 打印区块链。
func print() {
	chain := core.LoadBlockChain()
	defer chain.CloseDatabase()
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
func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  new        -address <address>                        Create a new blockchain and reward subsidy to <address>.")
	fmt.Println("  balance    -address <address>                        Query balance of <address> stored in the blockchain.")
	fmt.Println("  send       -from <from> -to <to> -amount <amount>    Send <amount> of coins from <from> to <to>.")
	fmt.Println("  print                                                Print blocks info of the blockchain.")
	fmt.Println("  help                                                 Show help for commands.")
}
