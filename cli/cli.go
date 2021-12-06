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
		panic("[Error] Command not given! Use command `help` to check out usage.")
	}

	// 区块链创建。
	chainCmd := flag.NewFlagSet("chain", flag.ExitOnError)
	chainAddr := chainCmd.String("address", "", "The address who mined out genesis block.")
	// 钱包创建。
	walletCmd := flag.NewFlagSet("wallet", flag.ExitOnError)
	// 地址展示。
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	// 余额查询。
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	balanceAddr := balanceCmd.String("address", "", "The address being queried.")
	// 交易。
	tradeCmd := flag.NewFlagSet("trade", flag.ExitOnError)
	tradeFrom := tradeCmd.String("from", "", "Source wallet address.")
	tradeTo := tradeCmd.String("to", "", "Destination wallet address.")
	tradeAmount := tradeCmd.String("amount", "0", "Amount of coins to trade.")
	// 区块链打印。
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	// 帮助信息。
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
		err = errors.New("command not supported! Use command `help` to check out usage")
	}
	if err != nil {
		panic(err)
	}

	if chainCmd.Parsed() {
		if *chainAddr == "" {
			chainCmd.Usage()
			os.Exit(1)
		}
		newChain(*chainAddr)

	} else if walletCmd.Parsed() {
		newWallet()

	} else if listCmd.Parsed() {
		list()

	} else if balanceCmd.Parsed() {
		if *balanceAddr == "" {
			balanceCmd.Usage()
			os.Exit(1)
		}
		balance(*balanceAddr)

	} else if tradeCmd.Parsed() {
		amount, err := strconv.Atoi(*tradeAmount)
		if err != nil {
			tradeCmd.Usage()
			os.Exit(1)
		}
		if *tradeFrom == "" || *tradeTo == "" || amount <= 0 {
			tradeCmd.Usage()
			os.Exit(1)
		}
		trade(*tradeFrom, *tradeTo, amount)

	} else if printCmd.Parsed() {
		print()

	} else if helpCmd.Parsed() {
		help()
	}
}
