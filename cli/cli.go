package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// 命令行结构。
type CLI struct{}

// 检验参数。
func validateArgs() {
	if len(os.Args) < 2 {
		fmt.Println("[Error] No command given!")
		showHelp()
		os.Exit(1)
	}
}

// 运行命令行实例。
func (cli *CLI) Run() {
	validateArgs()

	// 自定义命令行指令。
	chainCmd := flag.NewFlagSet("chain", flag.ExitOnError)
	walletCmd := flag.NewFlagSet("wallet", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	newChainAddr := chainCmd.String("address", "", "The address to which rewards should be sent")
	balanceAddr := balanceCmd.String("address", "", "The address whose balance is being queried.")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.String("amount", "0", "Amount of coins")

	// 解析命令行参数。
	var err error
	switch os.Args[1] {
	case "chain":
		err = chainCmd.Parse(os.Args[2:])
	case "wallet":
		err = walletCmd.Parse(os.Args[2:])
	case "list":
		err = listCmd.Parse(os.Args[2:])
	case "balance":
		err = balanceCmd.Parse(os.Args[2:])
	case "send":
		err = sendCmd.Parse(os.Args[2:])
	case "print":
		err = printCmd.Parse(os.Args[2:])
	case "help":
		err = helpCmd.Parse(os.Args[2:])
	default:
		showHelp()
		os.Exit(1)
	}
	if err != nil {
		log.Panic(err)
	}

	if chainCmd.Parsed() {
		if *newChainAddr == "" {
			chainCmd.Usage()
			os.Exit(1)
		}
		newChain(*newChainAddr)

	} else if walletCmd.Parsed() {
		newWallet()

	} else if listCmd.Parsed() {
		listAddresses()

	} else if balanceCmd.Parsed() {
		if *balanceAddr == "" {
			balanceCmd.Usage()
			os.Exit(1)
		}
		queryBalance(*balanceAddr)

	} else if sendCmd.Parsed() {
		amount, err := strconv.Atoi(*sendAmount)
		if err != nil {
			sendCmd.Usage()
			os.Exit(1)
		}
		if *sendFrom == "" || *sendTo == "" || amount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		send(*sendFrom, *sendTo, amount)

	} else if printCmd.Parsed() {
		print()

	} else if helpCmd.Parsed() {
		showHelp()
	}
}
