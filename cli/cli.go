package cli

import (
	"blockchain/blockchain"
	"flag"
	"log"
	"os"
)

// 命令行结构。
type CLI struct {
	chain *blockchain.BlockChain
}

// 创建命令行实例。
func CreateCLI(chain *blockchain.BlockChain) *CLI {
	return &CLI{chain}
}

// 运行命令行实例。
func (cli *CLI) Run() {
	// 定义命令行格式。
	addDataCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addData := addDataCmd.String("data", "", "Block data")

	// 解析命令行参数。
	var err error
	switch os.Args[1] {
	case "add":
		err = addDataCmd.Parse(os.Args[2:])
	case "print":
		err = printChainCmd.Parse(os.Args[2:])
	default:
		os.Exit(1)
	}
	if err != nil {
		log.Panic(err)
	}

	// 如果是添加区块命令。
	if addDataCmd.Parsed() {
		if *addData == "" {
			addDataCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addData)
	}

	// 如果是打印区块链命令。
	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
