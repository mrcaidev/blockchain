package blockchain

type BlockChain struct {
	Blocks []*Block
}

func CreateBlockChain() *BlockChain {
	return &BlockChain{[]*Block{CreateGenesisBlock()}}
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}
