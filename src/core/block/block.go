package block

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"blockchain/core/transaction"
)

// 区块结构。
type Block struct {
	Timestamp     int64                      // 时间戳。
	Transactions  []*transaction.Transaction // 交易列表。
	PrevBlockHash []byte                     // 前一区块哈希值。
	Hash          []byte                     // 本区块哈希值。
	Nonce         int                        // 随机数。
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

	fmt.Println("Mining new block...")
	block.proofWork()

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

	for index, tx := range b.Transactions {
		fmt.Printf("\nTransaction %d:\n", index)
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
