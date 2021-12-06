package chain

import (
	"blockchain/core/block"
	"blockchain/utils"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"

	"blockchain/transaction"

	"github.com/boltdb/bolt"
)

// 数据库。
const dbPath = "blockchain.db"
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
	if utils.HasFile(dbPath) {
		panic("blockchain already exists")
	}

	// 打开数据库。
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		panic(err)
	}

	// 创建 coinbase 交易和相应的创世块。
	coinbaseTX := NewCoinbaseTX(address, genesisCoinbase)
	genesisBlock := block.NewGenesisBlock(coinbaseTX)
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
	if !utils.HasFile(dbPath) {
		panic("blockchain not found")
	}

	// 打开数据库。
	db, err := bolt.Open(dbPath, 0600, nil)
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
		if !c.VerifyTX(tx) {
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

// 找到未消费的交易。
func (c *Chain) FindUTX(pubkeyHash []byte) []*transaction.Transaction {
	var utxs []*transaction.Transaction
	stxoIndexes := make(map[string][]int)

	// 遍历区块链中的每一个区块。
	iter := c.Iterator()
	for {
		// 遍历区块中的每一笔交易。
		curBlock := iter.Next()
		for _, tx := range curBlock.Transactions {
			txID := hex.EncodeToString(tx.ID)
			// 遍历交易的输出。
			for txoIndex, txo := range tx.Outputs {
				// 如果这笔输出已经被消费了，就检查下一笔。
				if stxoIndexes[txID] != nil {
					for _, stxoIndex := range stxoIndexes[txID] {
						if stxoIndex == txoIndex {
							break
						}
					}
				} else if txo.IsUnlockableWith(pubkeyHash) {
					utxs = append(utxs, tx)
				}
			}
			// 如果这笔交易不是 coinbase 交易。
			if !tx.IsCoinbase() {
				// 遍历这笔交易的输入。
				for _, txi := range tx.Inputs {
					// 如果能解锁，说明被花费过。
					if txi.IsLockedWith(pubkeyHash) {
						refID := hex.EncodeToString(txi.RefID)
						stxoIndexes[refID] = append(stxoIndexes[refID], txi.RefIndex)
					}
				}
			}
		}
		if len(curBlock.PrevBlockHash) == 0 {
			break
		}
	}
	return utxs
}

// 找到未消费的交易输出。
func (c *Chain) FindUTXO(pubkeyHash []byte) []*transaction.TXOutput {
	var utxos []*transaction.TXOutput
	utxs := c.FindUTX(pubkeyHash)
	// 遍历每一笔未消费交易。
	for _, tx := range utxs {
		// 遍历交易的每一笔输出。
		for _, txo := range tx.Outputs {
			// 如果这笔输出可以被解锁，就追加进列表。
			if txo.IsUnlockableWith(pubkeyHash) {
				utxos = append(utxos, txo)
			}
		}
	}
	return utxos
}

// 找到用于当次支付的未消费交易输出。
func (c *Chain) FindUTXOToPay(pubkeyHash []byte, amount int) (int, map[string][]int) {
	utxoToPay := make(map[string][]int)
	utxs := c.FindUTX(pubkeyHash)
	atHand := 0
	// 遍历未花费的交易。
	for _, tx := range utxs {
		txID := hex.EncodeToString(tx.ID)
		// 遍历交易的输出。
		for txoIndex, txo := range tx.Outputs {
			// 如果输出能被解锁，就追加进字典。
			if txo.IsUnlockableWith(pubkeyHash) {
				atHand += txo.Value
				utxoToPay[txID] = append(utxoToPay[txID], txoIndex)
				// 如果钱够了，就不再检查之后的输出。
				if atHand >= amount {
					return atHand, utxoToPay
				}
			}
		}
	}
	// 如果遍历完所有的输出，钱还是不够，就返回空值。
	return 0, map[string][]int{}
}

// 凭 ID 查找交易。
func (c *Chain) FindTX(ID []byte) *transaction.Transaction {
	iter := c.Iterator()
	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return tx
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	panic("transaction not found")
}

// 对交易进行数字签名。
func (c *Chain) SignTX(tx *transaction.Transaction, privkey ecdsa.PrivateKey) {
	refTxs := make(map[string]*transaction.Transaction)

	// 找到交易每一笔输入引用的交易。
	for _, txi := range tx.Inputs {
		refTx := c.FindTX(txi.RefID)
		refTxs[hex.EncodeToString(refTx.ID)] = refTx
	}
	tx.Sign(privkey, refTxs)
}

// 验证交易的数字签名。
func (c *Chain) VerifyTX(tx *transaction.Transaction) bool {
	refTxs := make(map[string]*transaction.Transaction)
	for _, txi := range tx.Inputs {
		refTx := c.FindTX(txi.RefID)
		refTxs[hex.EncodeToString(refTx.ID)] = refTx
	}
	return tx.Verify(refTxs)
}

// 获得区块链属于某地址的余额。
func (c *Chain) GetBalance(address string) int {
	// 获得地址内蕴含的公钥哈希。
	pubkeyHash := utils.Base58Decode([]byte(address))
	pubkeyHash = pubkeyHash[1 : len(pubkeyHash)-4]

	// 使用该公钥哈希，寻找每一笔未消费的输出。
	UTXOs := c.FindUTXO(pubkeyHash)

	// 累加这些输出内的余额。
	balance := 0
	for _, utxo := range UTXOs {
		balance += utxo.Value
	}

	return balance
}
