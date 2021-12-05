package core

import (
	"blockchain/utils"

	"github.com/boltdb/bolt"
)

// 区块链迭代器结构。
type BlockChainIterator struct {
	curHash []byte   // 当前指向区块的哈希值。
	db      *bolt.DB // 数据库连接。
}

// 创建迭代器。
func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.rear, chain.db}
}

// 从尾部开始遍历区块链。
func (iter *BlockChainIterator) Next() *Block {
	// 获取迭代器当前指向的区块。
	var block *Block
	err := iter.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BlocksBucket))
		blockSeq := bucket.Get(iter.curHash)
		block = DeserializeBlock(blockSeq)
		return nil
	})
	if err != nil {
		panic(err)
	}
	// 迭代器移向前一个区块。
	iter.curHash = block.PrevBlockHash
	return block
}
