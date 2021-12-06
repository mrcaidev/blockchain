package block

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	tx "blockchain/transaction"
)

// 区块结构。
type Block struct {
	Timestamp     int64             // 区块时间戳。
	Transactions  []*tx.Transaction // 交易列表。
	PrevBlockHash []byte            // 前一区块摘要值。
	Hash          []byte            // 区块标识。
	Nonce         int               // 区块随机数。
}

// 创建区块。
func NewBlock(TXs []*tx.Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  TXs,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	fmt.Println("Mining block of this transaction...")
	proof := newPow(block)
	proof.Run()

	return block
}

// 创建创世块。
func NewGenesisBlock(coinbaseTX *tx.Transaction) *Block {
	return NewBlock([]*tx.Transaction{coinbaseTX}, []byte{})
}

// 打印区块信息。
func (block *Block) Print() {
	fmt.Println("--------------------------------------------------------------------------------")

	fmt.Printf("Hash:      %x\n", block.Hash)
	fmt.Printf("Nonce:     %d\n", block.Nonce)
	fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
	fmt.Println()

	for index, tx := range block.Transactions {
		fmt.Printf("Transaction %d:\n", index)
		tx.Print()
	}

	fmt.Println("--------------------------------------------------------------------------------")
}

// 序列化区块。
func (block *Block) Serialize() []byte {
	var seq bytes.Buffer

	encoder := gob.NewEncoder(&seq)
	err := encoder.Encode(block)
	if err != nil {
		panic(err)
	}

	return seq.Bytes()
}

// 反序列化区块。
func DeserializeBlock(seq []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(seq))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}

	return &block
}
