package chain

import (
	"blockchain/core/block"
	"blockchain/utils"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"

	tx "blockchain/transaction"

	"github.com/boltdb/bolt"
)

// 数据库路径。
const dbPath = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbase = "Genesis Coinbase"
const lastHashKey = "l"

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

// 添加区块。
func (chain *Chain) AddBlock(TXs []*tx.Transaction) {
	// 验证每笔交易。
	for _, TX := range TXs {
		if !chain.VerifyTX(TX) {
			panic("invalid transaction")
		}
	}

	// 从数据库获取最后一个区块的哈希值。
	var lastHash []byte
	err := chain.db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	if err != nil {
		panic(err)
	}

	// 创建新区块。
	newBlock := block.NewBlock(TXs, lastHash)
	chain.rear = newBlock.Hash

	// 将区块录入 bucket。
	err = chain.db.Update(func(t *bolt.Tx) error {
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

// 关闭数据库连接。
func (chain *Chain) CloseDatabase() {
	chain.db.Close()
}

// 找到未消费的交易。
func (chain *Chain) FindUTX(pubkeyHash []byte) []*tx.Transaction {
	var utxs []*tx.Transaction
	stxoIndexes := make(map[string][]int)

	// 遍历区块链中的每一个区块。
	iter := chain.Iterator()
	for {
		// 遍历区块中的每一笔交易。
		curBlock := iter.Next()
		for _, TX := range curBlock.Transactions {
			TXID := hex.EncodeToString(TX.ID)
			// 遍历交易的输出。
			for txoIndex, txo := range TX.Outputs {
				// 如果这笔输出已经被消费了，就检查下一笔。
				if stxoIndexes[TXID] != nil {
					for _, stxoIndex := range stxoIndexes[TXID] {
						if stxoIndex == txoIndex {
							break
						}
					}
				} else if txo.IsUnlockableWith(pubkeyHash) {
					utxs = append(utxs, TX)
				}
			}
			// 如果这笔交易不是coinbase交易。
			if !TX.IsCoinbase() {
				// 遍历这笔交易的输入。
				for _, txi := range TX.Inputs {
					// 如果能解锁，说明这是被花费过的。
					if txi.UsesKey(pubkeyHash) {
						refID := hex.EncodeToString(txi.RefID)
						stxoIndexes[refID] = append(stxoIndexes[refID], txi.RefIndex)
					}
				}
			}
		}
		// 直到遍历结束。
		if len(curBlock.PrevBlockHash) == 0 {
			break
		}
	}
	return utxs
}

// 找到未消费的交易输出。
func (chain *Chain) FindUTXO(pubkeyHash []byte) []*tx.TXOutput {
	var utxos []*tx.TXOutput
	utxs := chain.FindUTX(pubkeyHash)
	// 遍历每一笔未消费交易。
	for _, TX := range utxs {
		// 遍历交易的每一笔输出。
		for _, txo := range TX.Outputs {
			// 如果这笔输出可以被解锁，就追加进列表。
			if txo.IsUnlockableWith(pubkeyHash) {
				utxos = append(utxos, txo)
			}
		}
	}
	return utxos
}

// 找到用于当次支付的未消费交易输出。
func (chain *Chain) FindUTXOToPay(pubkeyHash []byte, amount int) (int, map[string][]int) {
	utxoToPay := make(map[string][]int)
	utxs := chain.FindUTX(pubkeyHash)
	atHand := 0
	// 遍历未花费的交易。
	for _, TX := range utxs {
		TXID := hex.EncodeToString(TX.ID)
		// 遍历交易的输出。
		for txoIndex, txo := range TX.Outputs {
			// 如果输出能被解锁，就追加进字典。
			if txo.IsUnlockableWith(pubkeyHash) {
				atHand += txo.Value
				utxoToPay[TXID] = append(utxoToPay[TXID], txoIndex)
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
func (chain *Chain) FindTX(ID []byte) *tx.Transaction {
	iter := chain.Iterator()
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
func (chain *Chain) SignTX(TX *tx.Transaction, privkey ecdsa.PrivateKey) {
	prevTXs := make(map[string]*tx.Transaction)

	// 找到交易每一笔输入引用的交易。
	for _, txi := range TX.Inputs {
		TX := chain.FindTX(txi.RefID)
		prevTXs[hex.EncodeToString(TX.ID)] = TX
	}
	TX.Sign(privkey, prevTXs)
}

// 验证交易的数字签名。
func (chain *Chain) VerifyTX(TX *tx.Transaction) bool {
	prevTXs := make(map[string]*tx.Transaction)
	for _, txi := range TX.Inputs {
		prevTx := chain.FindTX(txi.RefID)
		prevTXs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	return TX.Verify(prevTXs)
}
