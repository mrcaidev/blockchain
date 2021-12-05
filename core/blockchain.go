package core

import (
	"blockchain/utils"
	"blockchain/wallet"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	tx "blockchain/transaction"

	"github.com/boltdb/bolt"
)

// 区块链结构。
type BlockChain struct {
	rear []byte   // 最后一个记录的哈希值。
	db   *bolt.DB // 数据库连接。
}

// 创建区块链。
func NewBlockChain(address string) *BlockChain {
	// 如果数据库已经存在，就报错退出。
	if utils.HasFile(utils.DBFile) {
		fmt.Println("There is already a blockchain.")
		os.Exit(1)
	}

	// 打开数据库。
	db, err := bolt.Open(utils.DBFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	// 创建coinbase交易和相应的创世块。
	cbtx := tx.NewCoinbaseTX(address, utils.GenesisCoinbase)
	genesis := NewGenesisBlock(cbtx)
	rear := genesis.Hash

	// 将创世块录入新bucket。
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(utils.BlocksBucket))
		if err != nil {
			return err
		}
		genesis.StoreInBucket(bucket)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return &BlockChain{rear, db}
}

// 读取区块链。
func LoadBlockChain() *BlockChain {
	// 如果找不到数据库，就报错退出。
	if !utils.HasFile(utils.DBFile) {
		fmt.Println("No chain found. Create one first.")
		os.Exit(1)
	}

	// 打开数据库。
	db, err := bolt.Open(utils.DBFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	// 从数据库读取目前的区块链信息。
	var rear []byte
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BlocksBucket))
		rear = bucket.Get([]byte("l"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	return &BlockChain{rear, db}
}

// 添加区块。
func (chain *BlockChain) AddBlock(TXs []*tx.Transaction) {
	// 从数据库获取最后一个区块的哈希值。
	var lastHash []byte

	for _, TX := range TXs {
		if !chain.VerifyTX(TX) {
			panic("[Error] Invalid transaction.")
		}
	}

	err := chain.db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(utils.BlocksBucket))
		lastHash = bucket.Get([]byte("l"))
		return nil
	})
	if err != nil {
		panic(err)
	}

	// 创建新区块。
	newBlock := NewBlock(TXs, lastHash)
	err = chain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BlocksBucket))
		newBlock.StoreInBucket(bucket)
		chain.rear = newBlock.Hash
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// 关闭数据库连接。
func (chain *BlockChain) CloseDatabase() {
	chain.db.Close()
}

// 找到UTX。
func (chain *BlockChain) FindUTX(pubkeyHash []byte) []*tx.Transaction {
	var UTXs []*tx.Transaction
	STXOs := make(map[string][]int)

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
				if STXOs[txID] != nil {
					for _, spentOutput := range STXOs[txID] {
						if spentOutput == index {
							continue Outputs
						}
					}
				}
				// 如果可以被解锁，说明可以被花费，追加进结果。
				if output.IsUnlockableWith(pubkeyHash) {
					UTXs = append(UTXs, tx)
				}
			}

			// 如果这笔交易不是coinbase交易。
			if !tx.IsCoinbase() {
				// 遍历这笔交易的输入。
				for _, input := range tx.Inputs {
					if input.UsesKey(pubkeyHash) {
						refID := hex.EncodeToString(input.RefID)
						STXOs[refID] = append(STXOs[refID], input.RefIndex)
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
func (chain *BlockChain) FindUTXO(pubkeyHash []byte) []tx.TXO {
	var UTXOs []tx.TXO
	UTXs := chain.FindUTX(pubkeyHash)

	for _, tx := range UTXs {
		for _, output := range tx.Outputs {
			if output.IsUnlockableWith(pubkeyHash) {
				UTXOs = append(UTXOs, output)
			}
		}
	}
	return UTXOs
}

// 找到用于当次支付的UTXO。
func (chain *BlockChain) FindUTXOToPay(pubkeyHash []byte, amount int) (int, map[string][]int) {
	UTXOToPay := make(map[string][]int)
	UTXs := chain.FindUTX(pubkeyHash)
	atHand := 0
TxTraverse:
	// 遍历未花费的交易。
	for _, tx := range UTXs {
		txID := hex.EncodeToString(tx.ID)
		// 遍历交易的输出。
		for index, output := range tx.Outputs {
			if output.IsUnlockableWith(pubkeyHash) {
				atHand += output.Value
				UTXOToPay[txID] = append(UTXOToPay[txID], index)
				if atHand >= amount {
					break TxTraverse
				}
			}
		}
	}
	return atHand, UTXOToPay
}

func (chain *BlockChain) FindTX(ID []byte) (*tx.Transaction, error) {
	iter := chain.Iterator()
	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return tx, nil
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return &tx.Transaction{}, errors.New("transaction not found")
}

func (chain *BlockChain) SignTX(TX *tx.Transaction, privkey ecdsa.PrivateKey) {
	prevTxs := make(map[string]*tx.Transaction)

	for _, input := range TX.Inputs {
		prevTx, err := chain.FindTX(input.RefID)
		if err != nil {
			panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	TX.Sign(privkey, prevTxs)
}

func (chain *BlockChain) VerifyTX(TX *tx.Transaction) bool {
	prevTxs := make(map[string]*tx.Transaction)
	for _, input := range TX.Inputs {
		prevTx, err := chain.FindTX(input.RefID)
		if err != nil {
			panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	return TX.Verify(prevTxs)
}

// 创建一笔UTXO交易。
func (chain *BlockChain) NewUTXOTX(from string, to string, cost int) *tx.Transaction {
	var (
		inputs  []tx.TXI
		outputs []tx.TXO
	)

	wallets := wallet.NewWallets()
	wallet := wallets.GetWallet(from)
	pubkeyHash := utils.HashPubKey(wallet.PublicKey)

	deposit, UTXOToPay := chain.FindUTXOToPay(pubkeyHash, cost)

	// 如果发起方的钱不够了，就报错退出。
	if deposit < cost {
		panic("[Error] Not enough money.")
	}

	// 创建输入。
	for txid, output := range UTXOToPay {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			panic(err)
		}
		for _, out := range output {
			inputs = append(inputs, tx.TXI{
				RefID:     txID,
				RefIndex:  out,
				Signature: nil,
				Pubkey:    wallet.PublicKey,
			})
		}
	}

	// 创建输出。
	outputs = append(outputs, tx.NewTXO(cost, to))

	// 如果需要找零，就多加一笔记录。
	if deposit > cost {
		outputs = append(outputs, tx.NewTXO(deposit-cost, from))
	}

	// 将输入、输出存储进该次交易内。
	newTX := tx.Transaction{
		ID:      nil,
		Inputs:  inputs,
		Outputs: outputs,
	}
	newTX.ID = newTX.Hash()
	chain.SignTX(&newTX, wallet.PrivateKey)
	return &newTX
}
