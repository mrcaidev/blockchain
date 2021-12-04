package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// 挖矿的奖励。
const subsidy = 10

// 交易输入结构。
type TXI struct {
	RefID     []byte // 引用交易的ID。
	OutIndex  int    // 交易输出的索引。
	ScriptSig string // 密钥。
}

// 交易输出结构。
type TXO struct {
	Value        int    // 交易输出的币值。
	ScriptPubKey string // 公钥。
}

// 交易结构。
type Transaction struct {
	ID      []byte // 该笔交易的ID。
	Inputs  []TXI  // 该笔交易的输入。
	Outputs []TXO  // 该笔交易的输出。
}

// 创建一笔coinbase交易。
func NewCoinbaseTransaction(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	// 没有输入，只有一个输出。
	txi := TXI{[]byte{}, -1, data}
	txo := TXO{subsidy, to}
	tx := Transaction{nil, []TXI{txi}, []TXO{txo}}
	tx.SetID()
	return &tx
}

// 创建一笔UTXO交易。
func NewUTXOTransaction(from string, to string, cost int, chain *BlockChain) *Transaction {
	var (
		inputs  []TXI
		outputs []TXO
	)
	deposit, validOutputs := chain.FindUTXOToPay(from, cost)

	// 如果发起方的钱不够了，就报错退出。
	if deposit < cost {
		log.Panic("[Error] Not enough money.")
	}

	// 创建输入。
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			inputs = append(inputs, TXI{txID, out, from})
		}
	}

	// 创建输出。
	outputs = append(outputs, TXO{cost, to})

	// 如果需要找零，就多加一笔记录。
	if deposit > cost {
		outputs = append(outputs, TXO{deposit - cost, from})
	}

	// 将输入、输出存储进该次交易内。
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}

// 判断该输入能否在指定地址解锁输出。
func (txi *TXI) CanUnlockOutputWith(unlockingData string) bool {
	return txi.ScriptSig == unlockingData
}

// 判断该输出能否被指定地址解锁。
func (txo *TXO) CanBeUnlockedWith(unlockingData string) bool {
	return txo.ScriptPubKey == unlockingData
}

// 判断该笔交易是否为coinbase交易。
// 同时满足一下三个条件，说明是coinbase交易：
// 1. 只有一个输入；
// 2. 输入没有引用之前的交易；
// 3. 输入在之前的输出里索引为-1。
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].RefID) == 0 && tx.Inputs[0].OutIndex == -1
}

// 计算交易的ID。
func (tx *Transaction) SetID() {
	var id bytes.Buffer

	encoder := gob.NewEncoder(&id)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	idHash := sha256.Sum256(id.Bytes())
	tx.ID = idHash[:]
}

func (tx *Transaction) Print() {
	fmt.Printf("  ID: %x\n", tx.ID)
	for index, txi := range tx.Inputs {
		fmt.Printf("  Input %d:\n", index)
		fmt.Printf("    RefID:     %x\n", txi.RefID)
		fmt.Printf("    OutIndex:  %d\n", txi.OutIndex)
		fmt.Printf("    ScriptSig: %s\n", txi.ScriptSig)
	}
	for index, txo := range tx.Outputs {
		fmt.Printf("  Output %d:\n", index)
		fmt.Printf("    Value:        %d\n", txo.Value)
		fmt.Printf("    ScriptPubKey: %s\n", txo.ScriptPubKey)
	}
}
