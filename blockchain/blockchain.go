package blockchain

// 区块链结构。
type BlockChain struct {
	Blocks []*Block
}

// 创建区块链。
func CreateBlockChain() *BlockChain {
	return &BlockChain{[]*Block{CreateGenesisBlock()}}
}

// 添加区块。
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.ID)
	chain.Blocks = append(chain.Blocks, newBlock)
}
