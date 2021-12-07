package blockchain

import (
	"blockchain/core/block"
	"blockchain/transaction"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/boltdb/bolt"
)

// 挖出新块的奖励。
const subsidy = 10

// 数据库。
const utxoBucket = "utxo"

// 凭 ID 查找交易。
func (c *Chain) FindTx(ID []byte) *transaction.Transaction {
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

// 找到所有未消费的交易输出。
func (c *Chain) FindUtxos() map[string]transaction.TXOutputs {
	utxos := make(map[string]transaction.TXOutputs)
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
				} else {
					utxo := utxos[txID]
					utxo.Outputs = append(utxo.Outputs, txo)
					utxos[txID] = utxo
				}
			}
			// 如果这笔交易不是 coinbase 交易。
			if !tx.IsCoinbase() {
				// 遍历这笔交易的输入。
				for _, txi := range tx.Inputs {
					refID := hex.EncodeToString(txi.RefID)
					stxoIndexes[refID] = append(stxoIndexes[refID], txi.RefIndex)
				}
			}
		}
		if len(curBlock.PrevBlockHash) == 0 {
			break
		}
	}
	return utxos
}

// 找到指定公钥可解锁的未消费交易输出。
func (c *Chain) FindPayableUtxos(pubkeyHash []byte) []*transaction.TXOutput {
	var utxos []*transaction.TXOutput

	err := c.db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(utxoBucket))
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			txos := transaction.DeserializeTXOutputs(value)

			for _, txo := range txos.Outputs {
				if txo.IsUnlockableWith(pubkeyHash) {
					utxos = append(utxos, txo)
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return utxos
}

// 找到指定公钥可解锁的、用于当次支付的未消费交易输出。
func (c *Chain) FindUtxosToPay(pubkeyHash []byte, amount int) (int, map[string][]int) {
	utxoToPay := make(map[string][]int)
	atHand := 0

	err := c.db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(utxoBucket))
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			txID := hex.EncodeToString(key)
			txos := transaction.DeserializeTXOutputs(value)

			for txoIndex, txo := range txos.Outputs {
				if txo.IsUnlockableWith(pubkeyHash) {
					atHand += txo.Value
					utxoToPay[txID] = append(utxoToPay[txID], txoIndex)
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return atHand, utxoToPay
}

// 对交易进行数字签名。
func (c *Chain) SignTx(tx *transaction.Transaction, privkey ecdsa.PrivateKey) {
	refTxs := make(map[string]*transaction.Transaction)

	// 找到交易每一笔输入引用的交易。
	for _, txi := range tx.Inputs {
		refTx := c.FindTx(txi.RefID)
		refTxs[hex.EncodeToString(refTx.ID)] = refTx
	}
	tx.Sign(privkey, refTxs)
}

// 验证交易的数字签名。
func (c *Chain) VerifyTx(tx *transaction.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	refTxs := make(map[string]*transaction.Transaction)
	for _, txi := range tx.Inputs {
		refTx := c.FindTx(txi.RefID)
		refTxs[hex.EncodeToString(refTx.ID)] = refTx
	}
	return tx.Verify(refTxs)
}

// 获取区块链内的交易数量。
func (c *Chain) CountTx() int {
	cnt := 0

	err := c.db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(utxoBucket))
		cursor := bucket.Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			cnt++
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return cnt
}

func (c *Chain) Reindex() {
	bucketName := []byte(utxoBucket)

	// 删除之前的 bucket。
	err := c.db.Update(func(t *bolt.Tx) error {
		err := t.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			return err
		}

		_, err = t.CreateBucket(bucketName)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	// 找到所有未花费的交易输出。
	utxos := c.FindUtxos()

	// 重新构建 bucket。
	err = c.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)

		for txIDString, utxo := range utxos {
			txID, err := hex.DecodeString(txIDString)
			if err != nil {
				panic(err)
			}

			err = b.Put(txID, utxo.Serialize())
			if err != nil {
				panic(err)
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (c *Chain) Update(block *block.Block) {
	err := c.db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, txi := range tx.Inputs {
					updatedTxos := transaction.TXOutputs{}
					txosBytes := bucket.Get(txi.RefID)
					txos := transaction.DeserializeTXOutputs(txosBytes)

					for txoIndex, txo := range txos.Outputs {
						if txoIndex != txi.RefIndex {
							updatedTxos.Outputs = append(updatedTxos.Outputs, txo)
						}
					}

					if len(updatedTxos.Outputs) == 0 {
						err := bucket.Delete(txi.RefID)
						if err != nil {
							return err
						}
					} else {
						err := bucket.Put(txi.RefID, updatedTxos.Serialize())
						if err != nil {
							return err
						}
					}
				}
			}

			newTxos := transaction.TXOutputs{}
			newTxos.Outputs = append(newTxos.Outputs, tx.Outputs...)

			err := bucket.Put(tx.ID, newTxos.Serialize())
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
