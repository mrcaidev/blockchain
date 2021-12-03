package blockchain

import (
	"fmt"
	"time"
)

// 区块结构。
type Block struct {
	ID            []byte // 区块标识。
	PrevBlockHash []byte // 前一区块摘要值。
	Timestamp     int64  // 区块时间戳。
	Nonce         int    // 区块随机数。
	Data          []byte // 数据。
}

// 创建区块。
func CreateBlock(data string, prevBlockHash []byte) *Block {
	// 录入时间戳、数据与前一区块摘要值。
	block := &Block{
		ID:            []byte{},
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().Unix(),
		Nonce:         0,
		Data:          []byte(data),
	}
	// 证明工作量。
	fmt.Printf("Running PoW of block '%s'\n", block.Data)
	pow := CreatePow(block)
	pow.Run()

	return block
}

// 创建创世块。
func CreateGenesisBlock() *Block {
	return CreateBlock("Genesis Block", []byte{})
}

// 打印区块信息。
func (block *Block) Print() {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("ID: %x\n", block.ID)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("Data: %s\n", block.Data)
	fmt.Println("--------------------------------------------------------------------------------")
}
