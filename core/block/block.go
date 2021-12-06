package block

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"blockchain/transaction"
)

// 区块结构。
type Block struct {
	Timestamp     int64                      // 区块时间戳。
	Transactions  []*transaction.Transaction // 交易列表。
	PrevBlockHash []byte                     // 前一区块摘要值。
	Hash          []byte                     // 区块标识。
	Nonce         int                        // 区块随机数。
}

// 创建区块。
func NewBlock(txs []*transaction.Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  txs,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	// 证明工作量。
	fmt.Println("Mining new block...")
	pow := newPow(block)
	pow.Run()

	return block
}

// 创建创世块。
func NewGenesisBlock(coinbaseTx *transaction.Transaction) *Block {
	return NewBlock([]*transaction.Transaction{coinbaseTx}, []byte{})
}

// 打印区块信息。
func (b *Block) Print() {
	fmt.Println("--------------------------------------------------------------------------------")

	fmt.Printf("Hash:      %x\n", b.Hash)
	fmt.Printf("Nonce:     %d\n", b.Nonce)
	fmt.Printf("Prev hash: %x\n", b.PrevBlockHash)
	fmt.Println()

	for index, tx := range b.Transactions {
		fmt.Printf("Transaction %d:\n", index)
		tx.Print()
	}

	fmt.Println("--------------------------------------------------------------------------------")
}

// 序列化区块。
func (b *Block) Serialize() []byte {
	var seq bytes.Buffer

	encoder := gob.NewEncoder(&seq)
	err := encoder.Encode(b)
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
