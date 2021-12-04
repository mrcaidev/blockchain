package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// 区块结构。
type Block struct {
	Timestamp     int64  // 区块时间戳。
	Data          []byte // 数据。
	PrevBlockHash []byte // 前一区块摘要值。
	Hash          []byte // 区块标识。
	Nonce         int    // 区块随机数。
}

// 序列化区块。
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
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

// 创建区块。
func CreateBlock(data string, prevBlockHash []byte) *Block {
	// 录入时间戳、数据与前一区块摘要值。
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
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

// 将区块追加进数据库。
func (block *Block) AddToBucket(bucket *bolt.Bucket) error {
	// 添加区块。
	err := bucket.Put(block.Hash, block.Serialize())
	if err != nil {
		return err
	}
	// 单独记录其哈希值，作为最后一个区块的哈希值。
	err = bucket.Put([]byte("l"), block.Hash)
	if err != nil {
		return err
	}

	return nil
}

// 打印区块信息。
func (block *Block) Print() {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("ID: %x\n", block.Hash)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("Data: %s\n", block.Data)
	fmt.Println("--------------------------------------------------------------------------------")
}
