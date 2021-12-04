package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// 区块结构。
type Block struct {
	Timestamp     int64          // 区块时间戳。
	Transactions  []*Transaction // 交易列表。
	PrevBlockHash []byte         // 前一区块摘要值。
	Hash          []byte         // 区块标识。
	Nonce         int            // 区块随机数。
}

// 创建区块。
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	// 录入时间戳、数据与前一区块摘要值。
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}
	// 证明工作量。
	pow := CreatePow(block)
	pow.Run()

	return block
}

// 创建创世块。
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// 序列化区块。
func (block *Block) Serialize() []byte {
	var seq bytes.Buffer
	encoder := gob.NewEncoder(&seq)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return seq.Bytes()
}

// 反序列化区块。
func DeserializeBlock(seq []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(seq))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// 将区块追加进数据库。
func (block *Block) AddToBucket(bucket *bolt.Bucket) {
	// 添加区块。
	err := bucket.Put(block.Hash, block.Serialize())
	if err != nil {
		log.Panic(err)
	}
	// 单独记录其哈希值，作为最后一个区块的哈希值。
	err = bucket.Put([]byte("l"), block.Hash)
	if err != nil {
		log.Panic(err)
	}
}

// 计算区块内交易的哈希值。
func (block *Block) TransactionHash() []byte {
	var hashes [][]byte

	for _, tx := range block.Transactions {
		hashes = append(hashes, tx.ID)
	}
	generalHash := sha256.Sum256(bytes.Join(hashes, []byte{}))

	return generalHash[:]
}

// 打印区块信息。
func (block *Block) Print() {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Hash:      %x\n", block.Hash)
	fmt.Printf("Nonce:     %d\n", block.Nonce)
	fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
	for index, tx := range block.Transactions {
		fmt.Printf("Transaction %d:\n", index)
		tx.Print()
	}
	fmt.Println("--------------------------------------------------------------------------------")
}
