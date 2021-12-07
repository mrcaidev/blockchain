package blockchain

import (
	"blockchain/core/block"
	"blockchain/core/transaction"
	"blockchain/utils"
	"os"

	"github.com/boltdb/bolt"
)

// 数据库。
const chainDbPath = "blockchain.db"
const blocksBucket = "blocks"
const lastHashKey = "l"

// 创世块 coinbase 内含数据。
const genesisCoinbase = "Genesis Coinbase"

// 区块链结构。
type Chain struct {
	rear []byte   // 最后一个记录的哈希值。
	db   *bolt.DB // 数据库连接。
}

// 创建区块链。
func NewChain(address string) *Chain {
	// 如果数据库已经存在，就报错退出。
	if !chainDbNotExists() {
		panic("blockchain already exists")
	}

	// 打开数据库。
	db, err := bolt.Open(chainDbPath, 0600, nil)
	if err != nil {
		panic(err)
	}

	// 创建 coinbase 交易和相应的创世块。
	coinbaseTx := NewCoinbaseTx(address, genesisCoinbase)
	genesisBlock := block.NewGenesisBlock(coinbaseTx)
	rear := genesisBlock.Hash

	// 将创世块录入新 bucket。
	err = db.Update(func(t *bolt.Tx) error {
		bucket, err := t.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}

		err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(lastHashKey), genesisBlock.Hash)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	return &Chain{rear, db}
}

// 读取区块链。
func LoadChain() *Chain {
	// 如果数据库不存在，就报错退出。
	if chainDbNotExists() {
		panic("blockchain not found")
	}

	// 打开数据库。
	db, err := bolt.Open(chainDbPath, 0600, nil)
	if err != nil {
		panic(err)
	}

	// 从数据库读取目前的区块链信息。
	var rear []byte
	err = db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		rear = bucket.Get([]byte(lastHashKey))
		return nil
	})
	if err != nil {
		panic(err)
	}

	return &Chain{rear, db}
}

// 向区块链添加区块。
func (c *Chain) AddBlock(txs []*transaction.Transaction) {
	// 验证每笔交易。
	for _, tx := range txs {
		if !c.VerifyTx(tx) {
			panic("invalid transaction")
		}
	}

	// 从数据库获取最后一个区块的哈希值。
	var lastHash []byte
	err := c.db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	if err != nil {
		panic(err)
	}

	// 创建新区块。
	newBlock := block.NewBlock(txs, lastHash)
	c.rear = newBlock.Hash

	// 将区块录入 bucket。
	err = c.db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))

		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(lastHashKey), newBlock.Hash)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	c.Update(newBlock)
}

// 关闭区块链的数据库连接。
func (c *Chain) Close() {
	c.db.Close()
}

// 打印区块链信息。
func (c *Chain) Print() {
	iter := c.Iterator()
	for {
		block := iter.Next()
		block.Print()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

// 获得区块链属于某地址的余额。
func (c *Chain) GetBalance(address string) int {
	// 获得地址内蕴含的公钥哈希。
	pubkeyHash := utils.Base58Decode([]byte(address))
	pubkeyHash = pubkeyHash[1 : len(pubkeyHash)-4]

	// 使用该公钥哈希，寻找每一笔未消费的交易输出。
	UTXOs := c.FindPayableUtxos(pubkeyHash)

	// 累加这些输出内的余额。
	balance := 0
	for _, utxo := range UTXOs {
		balance += utxo.Value
	}

	return balance
}

// 判断区块链数据库是否存在。
func chainDbNotExists() bool {
	_, err := os.Stat(chainDbPath)
	return os.IsNotExist(err)
}
