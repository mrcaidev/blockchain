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

// 判断该笔交易是否为coinbase交易。
// 同时满足一下三个条件，说明是coinbase交易：
// 1. 只有一个输入；
// 2. 输入没有引用之前的交易；
// 3. 输入在之前的输出里索引为-1。
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].RefID) == 0 && tx.Inputs[0].RefIndex == -1
}

// 序列化交易。
func (tx *Transaction) Serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

// 哈希化交易。
func (tx *Transaction) Hash() []byte {
	TXCopy := *tx
	TXCopy.ID = []byte{}
	hash := sha256.Sum256(TXCopy.Serialize())
	return hash[:]
}

// 深拷贝交易。
func (tx *Transaction) Deepcopy() *Transaction {
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

// 对每个输入签名。
func (tx *Transaction) Sign(privkey ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
	// 如果是coinbase交易，就不用签名。
	if tx.IsCoinbase() {
		return
	}

	// 检查各输入的ID是否正确。
	for _, txi := range tx.Inputs {
		if prevTXs[hex.EncodeToString(txi.RefID)].ID == nil {
			panic("previous transaction incorrect")
		}
	}

	txCopy := tx.Deepcopy()
	for txiIndex, txi := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(txi.RefID)]
		txCopy.Inputs[txiIndex].Signature = nil
		txCopy.Inputs[txiIndex].Pubkey = prevTX.Outputs[txi.RefIndex].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[txiIndex].Pubkey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privkey, txCopy.ID)
		if err != nil {
			panic(err)
		}
		tx.Inputs[txiIndex].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

// 检验交易输入的数字签名。
func (tx *Transaction) Verify(prevTXs map[string]*Transaction) bool {
	// 如果是coinbase交易，就不用验证。
	if tx.IsCoinbase() {
		return true
	}

	// 检查各输入的ID是否正确。
	for _, txi := range tx.Inputs {
		if prevTXs[hex.EncodeToString(txi.RefID)].ID == nil {
			panic("previous transaction id incorrect")
		}
	}

	txCopy := tx.Deepcopy()
	txCopy.Print()
	curve := elliptic.P256()

	for txiIndex, txi := range tx.Inputs {
		prevTX := prevTXs[hex.EncodeToString(txi.RefID)]
		txCopy.Inputs[txiIndex].Signature = nil
		txCopy.Inputs[txiIndex].Pubkey = prevTX.Outputs[txi.RefIndex].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[txiIndex].Pubkey = nil

		sigLen := len(txi.Signature)
		r := utils.BytesToBigInt(txi.Signature[:(sigLen / 2)])
		s := utils.BytesToBigInt(txi.Signature[(sigLen / 2):])

		keyLen := len(txi.Pubkey)
		x := utils.BytesToBigInt(txi.Pubkey[:(keyLen / 2)])
		y := utils.BytesToBigInt(txi.Pubkey[(keyLen / 2):])

		rawPubkey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
		if !ecdsa.Verify(&rawPubkey, txCopy.ID, r, s) {
			return false
		}
	}
	return true
}

// 打印交易信息。
func (tx *Transaction) Print() {
	fmt.Printf("  ID: %x\n", tx.ID)
	for index, txi := range tx.Inputs {
		fmt.Printf("  Input %d:\n", index)
		fmt.Printf("    RefID:        %x\n", txi.RefID)
		fmt.Printf("    RefIndex:     %d\n", txi.RefIndex)
		fmt.Printf("    Signature:    %x\n", txi.Signature)
		fmt.Printf("    Pubkey:       %x\n", txi.Pubkey)
	}
	for index, txo := range tx.Outputs {
		fmt.Printf("  Output %d:\n", index)
		fmt.Printf("    Value:        %d\n", txo.Value)
		fmt.Printf("    PubkeyHash:   %x\n", txo.PubkeyHash)
	}
}
