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
func (cli *CLI) validate() {
	if len(os.Args) < 2 {
		fmt.Println("[Error] No command given!")
		cli.showHelp()
		os.Exit(1)
	}
}

// 运行命令行实例。
func (cli *CLI) Run() {
	cli.validate()
	// 自定义命令行指令。
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	newAddress := newCmd.String("address", "", "The address to which rewards should be sent")
	balanceAddress := balanceCmd.String("address", "", "The address whose balance is being queried.")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.String("amount", "0", "Amount of coins")

	// 解析命令行参数。
	var err error
	switch os.Args[1] {
	case "new":
		err = newCmd.Parse(os.Args[2:])
	case "balance":
		err = balanceCmd.Parse(os.Args[2:])
	case "send":
		err = sendCmd.Parse(os.Args[2:])
	case "print":
		err = printCmd.Parse(os.Args[2:])
	case "help":
		err = helpCmd.Parse(os.Args[2:])
	default:
		cli.showHelp()
		os.Exit(1)
	}
	if err != nil {
		log.Panic(err)
	}

	// 如果是添加区块命令。
	if newCmd.Parsed() {
		if *newAddress == "" {
			newCmd.Usage()
			os.Exit(1)
		}
		cli.newChain(*newAddress)

	} else if balanceCmd.Parsed() {
		if *balanceAddress == "" {
			balanceCmd.Usage()
			os.Exit(1)
		}
		cli.queryBalance(*balanceAddress)

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
		cli.send(*sendFrom, *sendTo, amount)

	} else if printCmd.Parsed() {
		cli.print()
	} else if helpCmd.Parsed() {
		cli.showHelp()
	}
}
