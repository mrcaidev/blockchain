package cli

import (
	"blockchain/blockchain"
	"fmt"
)

// 创建新区块链。
func (cli *CLI) newChain(address string) {
	chain := blockchain.NewBlockChain(address)
	chain.CloseDatabase()
	fmt.Println("Created!")
}

// 查询余额。
func (cli *CLI) queryBalance(address string) {
	chain := blockchain.LoadBlockChain()
	defer chain.CloseDatabase()

	balance := 0
	UTXOs := chain.FindUTXO(address)
	for _, output := range UTXOs {
		balance += output.Value
	}
	fmt.Printf("Balance of '%s': %d", address, balance)
}

// 发送币。
func (cli *CLI) send(from string, to string, amount int) {
	chain := blockchain.LoadBlockChain()
	defer chain.CloseDatabase()
	tx := blockchain.NewUTXOTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Send success!")
}

// 打印区块链。
func (cli *CLI) print() {
	chain := blockchain.LoadBlockChain()
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
func (cli *CLI) showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  new        -address <address>                        Create a new blockchain and reward subsidy to <address>.")
	fmt.Println("  balance    -address <address>                        Query balance of <address> stored in the blockchain.")
	fmt.Println("  send       -from <from> -to <to> -amount <amount>    Send <amount> of coins from <from> to <to>.")
	fmt.Println("  print                                                Print blocks info of the blockchain.")
	fmt.Println("  help                                                 Show help for commands.")
}
