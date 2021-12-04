package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

// 数据库路径。
const dbFile = "blockchain.db"
const blocksBucket = "blocks"

const genesisCoinbase = "Lorem ipsum"

// 区块链结构。
type BlockChain struct {
	rear []byte   // 最后一个记录的哈希值。
	db   *bolt.DB // 数据库连接。
}

// 区块链迭代器结构。
type BlockChainIterator struct {
	curHash []byte   // 当前指向区块的哈希值。
	db      *bolt.DB // 数据库连接。
}

// 创建区块链。
func NewBlockChain(address string) *BlockChain {
	// 如果数据库已经存在，就报错退出。
	if hasFile(dbFile) {
		fmt.Println("There is already a blockchain.")
		os.Exit(1)
	}

	// 打开数据库。
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// 初始化数据库。
	var rear []byte
	err = db.Update(func(tx *bolt.Tx) error {
		// 新建bucket。
		bucket, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}
		// 创建创世块并录入数据库。
		cbtx := NewCoinbaseTransaction(address, genesisCoinbase)
		genesis := NewGenesisBlock(cbtx)
		genesis.AddToBucket(bucket)
		rear = genesis.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{rear, db}
}

// 读取区块链。
func LoadBlockChain() *BlockChain {
	// 如果找不到数据库，就报错退出。
	if !hasFile(dbFile) {
		fmt.Println("No chain found. Create one first.")
		os.Exit(1)
	}

	// 打开数据库。
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// 从数据库读取目前的区块链信息。
	var rear []byte
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		rear = bucket.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{rear, db}
}

// 添加区块。
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	// 从数据库获取最后一个区块的哈希值。
	var lastHash []byte
	err := chain.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	// 创建新区块。
	newBlock := NewBlock(transactions, lastHash)
	err = chain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		newBlock.AddToBucket(bucket)
		chain.rear = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 关闭数据库连接。
func (chain *BlockChain) CloseDatabase() {
	chain.db.Close()
}

// 找到UTX。
func (chain *BlockChain) FindUTX(address string) []*Transaction {
	var UTXs []*Transaction
	STXs := make(map[string][]int)

	// 遍历区块链。
	iter := chain.Iterator()
	for {
		// 遍历区块中的交易。
		block := iter.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			// 遍历交易中的输出。
			for index, output := range tx.Outputs {
				// 如果被花费了。
				if STXs[txID] != nil {
					for _, spentOutput := range STXs[txID] {
						if spentOutput == index {
							continue Outputs
						}
					}
				}
				// 如果可以被解锁，说明可以被花费，追加进结果。
				if output.CanBeUnlockedWith(address) {
					UTXs = append(UTXs, tx)
				}
			}

			// 如果这笔交易不是coinbase交易。
			if !tx.IsCoinbase() {
				// 遍历这笔交易的输入。
				for _, input := range tx.Inputs {
					if input.CanUnlockOutputWith(address) {
						refID := hex.EncodeToString(input.RefID)
						STXs[refID] = append(STXs[refID], input.OutIndex)
					}
				}
			}
		}
		// 直到遍历结束。
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return UTXs
}

// 找到所有的UTXO。
func (chain *BlockChain) FindUTXO(address string) []TXO {
	var UTXOs []TXO
	UTXs := chain.FindUTX(address)

	for _, tx := range UTXs {
		for _, txo := range tx.Outputs {
			if txo.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, txo)
			}
		}
	}
	return UTXOs
}

// 找到用于当次支付的UTXO。
func (chain *BlockChain) FindUTXOToPay(address string, amount int) (int, map[string][]int) {
	outputsToPay := make(map[string][]int)
	unspentTxs := chain.FindUTX(address)
	atHand := 0
TxTraverse:
	// 遍历未花费的交易。
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		// 遍历交易的输出。
		for index, output := range tx.Outputs {
			if output.CanBeUnlockedWith(address) {
				atHand += output.Value
				outputsToPay[txID] = append(outputsToPay[txID], index)
				if atHand >= amount {
					break TxTraverse
				}
			}
		}
	}
	return atHand, outputsToPay
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
		bucket := tx.Bucket([]byte(blocksBucket))
		blockSeq := bucket.Get(iter.curHash)
		block = DeserializeBlock(blockSeq)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	// 迭代器移向前一个区块。
	iter.curHash = block.PrevBlockHash
	return block
}
