package cli

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

// 运行命令行实例。
func Run() {
	// 验证是否给出命令。
	if len(os.Args) < 2 {
		panic("use command `help` to check out usage")
	}

	// 钱包创建。
	walletCmd := flag.NewFlagSet("wallet", flag.ExitOnError)
	// 列出地址。
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	// 创建区块链。
	chainCmd := flag.NewFlagSet("chain", flag.ExitOnError)
	chainAddr := chainCmd.String("address", "", "The address who mined out genesis block.")
	// 查询余额。
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	balanceAddr := balanceCmd.String("address", "", "The address being queried.")
	// 发起交易。
	tradeCmd := flag.NewFlagSet("trade", flag.ExitOnError)
	tradeFrom := tradeCmd.String("from", "", "Source wallet address.")
	tradeTo := tradeCmd.String("to", "", "Destination wallet address.")
	tradeAmount := tradeCmd.String("amount", "0", "Amount of coins to trade.")
	// 打印区块链。
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	// 显示帮助。
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

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
	case "trade":
		err = tradeCmd.Parse(os.Args[2:])
	case "print":
		err = printCmd.Parse(os.Args[2:])
	case "help":
		err = helpCmd.Parse(os.Args[2:])
	default:
		err = errors.New("command not supported")
	}
	if err != nil {
		panic(err)
	}

	if chainCmd.Parsed() {
		if *chainAddr == "" {
			chainCmd.Usage()
		} else {
			newChain(*chainAddr)
		}

	} else if walletCmd.Parsed() {
		newWallet()

	} else if listCmd.Parsed() {
		listAddresses()

	} else if balanceCmd.Parsed() {
		if *balanceAddr == "" {
			balanceCmd.Usage()
			os.Exit(1)
		} else {
			queryBalance(*balanceAddr)
		}

	} else if tradeCmd.Parsed() {
		amount, err := strconv.Atoi(*tradeAmount)
		if err != nil {
			tradeCmd.Usage()
		} else if *tradeFrom == "" || *tradeTo == "" || amount <= 0 {
			tradeCmd.Usage()
		} else {
			startTrade(*tradeFrom, *tradeTo, amount)
		}

	} else if printCmd.Parsed() {
		printChain()

	} else if helpCmd.Parsed() {
		showHelp()
	}
}
