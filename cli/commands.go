package cli

import (
	"fmt"
)

// 添加区块。
func (cli *CLI) addBlock(data string) {
	cli.chain.AddBlock(data)
	fmt.Println("Success!")
}

// 打印区块链。
func (cli *CLI) printChain() {
	iter := cli.chain.Iterator()
	for {
		block := iter.Next()
		block.Print()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
