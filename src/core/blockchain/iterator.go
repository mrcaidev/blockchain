package blockchain

import (
	"blockchain/core/block"

	"github.com/boltdb/bolt"
)

// 区块链迭代器结构。
type chainIterator struct {
	curHash []byte   // 当前指向区块的哈希值。
	db      *bolt.DB // 数据库连接。
}

// 创建迭代器。
func (chain *Chain) Iterator() *chainIterator {
	return &chainIterator{chain.rear, chain.db}
}

// 从尾部开始遍历区块链。
func (iter *chainIterator) Next() *block.Block {
	// 获取迭代器当前指向的区块。
	var curBlock *block.Block
	err := iter.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		seq := bucket.Get(iter.curHash)
		curBlock = block.DeserializeBlock(seq)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// 迭代器移向前一个区块。
	iter.curHash = curBlock.PrevBlockHash
	return curBlock
}
