package transaction

import (
	"blockchain/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

// 交易结构。
type Transaction struct {
	ID      []byte      // 该笔交易的ID。
	Inputs  []*TXInput  // 该笔交易的输入。
	Outputs []*TXOutput // 该笔交易的输出。
}

// 判断交易是否为 coinbase 交易。
// 同时满足以下三个条件，说明是 coinbase 交易：
// 1. 只有一个输入；
// 2. 这个输入没有引用之前的交易；
// 3. 这个输入在之前的输出里索引为-1。
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].RefID) == 0 && tx.Inputs[0].RefIndex == -1
}

// 哈希化交易。
func (tx *Transaction) Hash() []byte {
	txCopy := *tx
	txCopy.ID = []byte{}
	hash := sha256.Sum256(txCopy.serialize())
	return hash[:]
}

// 创建交易的无签名副本。
func (tx *Transaction) noSigCopy() *Transaction {
	var (
		txiCopy []*TXInput
		txoCopy []*TXOutput
	)
	// 拷贝输入。
	for _, txi := range tx.Inputs {
		txiCopy = append(txiCopy, &TXInput{txi.RefID, txi.RefIndex, nil, nil})
	}
	// 拷贝输出。
	for _, txo := range tx.Outputs {
		txoCopy = append(txoCopy, &TXOutput{txo.Value, txo.PubkeyHash})
	}
	return &Transaction{tx.ID, txiCopy, txoCopy}
}

// 对每笔交易输入签名。
func (tx *Transaction) Sign(privkey ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
	// 如果当前交易是 coinbase 交易，就不用签名。
	if tx.IsCoinbase() {
		return
	}

	// 检查交易输入所属的交易的 ID 是否正确。
	for _, txi := range tx.Inputs {
		if prevTXs[hex.EncodeToString(txi.RefID)].ID == nil {
			panic("previous transaction ID incorrect")
		}
	}

	// 对交易的每一笔输入签名。
	txCopy := tx.noSigCopy()
	for txiIndex, txi := range txCopy.Inputs {
		// 获取签名的对象。
		prevTx := prevTXs[hex.EncodeToString(txi.RefID)]
		txCopy.Inputs[txiIndex].Signature = nil
		txCopy.Inputs[txiIndex].Pubkey = prevTx.Outputs[txi.RefIndex].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[txiIndex].Pubkey = nil

		// 存储 ECDSA 数字签名。
		r, s, err := ecdsa.Sign(rand.Reader, &privkey, txCopy.ID)
		if err != nil {
			panic(err)
		}
		tx.Inputs[txiIndex].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

// 检验交易输入的签名。
func (tx *Transaction) Verify(prevTXs map[string]*Transaction) bool {
	// 如果当前交易是 coinbase 交易，就不用验证。
	if tx.IsCoinbase() {
		return true
	}

	// 检查交易输入所属的交易的 ID 是否正确。
	for _, txi := range tx.Inputs {
		if prevTXs[hex.EncodeToString(txi.RefID)].ID == nil {
			panic("previous transaction ID incorrect")
		}
	}

	// 验证交易的每一笔输入的签名。
	txCopy := tx.noSigCopy()
	curve := elliptic.P256()
	for txiIndex, txi := range tx.Inputs {
		// 获取与签名时相同的对象。
		prevTX := prevTXs[hex.EncodeToString(txi.RefID)]
		txCopy.Inputs[txiIndex].Signature = nil
		txCopy.Inputs[txiIndex].Pubkey = prevTX.Outputs[txi.RefIndex].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[txiIndex].Pubkey = nil

		// 解析签名数据。
		sigLen := len(txi.Signature)
		r := utils.BytesToBigInt(txi.Signature[:(sigLen / 2)])
		s := utils.BytesToBigInt(txi.Signature[(sigLen / 2):])

		// 解析公钥数据。
		keyLen := len(txi.Pubkey)
		x := utils.BytesToBigInt(txi.Pubkey[:(keyLen / 2)])
		y := utils.BytesToBigInt(txi.Pubkey[(keyLen / 2):])

		supposedPubkey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
		if !ecdsa.Verify(&supposedPubkey, txCopy.ID, r, s) {
			return false
		}
	}
	return true
}

// 打印交易信息。
func (tx *Transaction) Print() {
	fmt.Printf("\n  ID: %x\n", tx.ID)
	for txiIndex, txi := range tx.Inputs {
		fmt.Printf("  Input %d:\n", txiIndex)
		fmt.Printf("    RefID:        %x\n", txi.RefID)
		fmt.Printf("    RefIndex:     %d\n", txi.RefIndex)
	}
	for txoIndex, txo := range tx.Outputs {
		fmt.Printf("  Output %d:\n", txoIndex)
		fmt.Printf("    Value:        %d\n", txo.Value)
		fmt.Printf("    PubkeyHash:   %x\n", txo.PubkeyHash)
	}
}

// 序列化交易。
func (tx *Transaction) serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
