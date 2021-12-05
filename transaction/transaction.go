package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/big"

	"blockchain/utils"
)

// 交易结构。
type Transaction struct {
	ID      []byte // 该笔交易的ID。
	Inputs  []TXI  // 该笔交易的输入。
	Outputs []TXO  // 该笔交易的输出。
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
	txCopy := *tx
	txCopy.ID = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// 深拷贝交易。
func (tx *Transaction) Deepcopy() *Transaction {
	var (
		inputsCopy  []TXI
		outputsCopy []TXO
	)
	// 拷贝输入。
	for _, input := range tx.Inputs {
		inputsCopy = append(inputsCopy, TXI{input.RefID, input.RefIndex, nil, nil})
	}
	// 拷贝输出。
	for _, output := range tx.Outputs {
		outputsCopy = append(outputsCopy, TXO{output.Value, output.PubkeyHash})
	}
	return &Transaction{tx.ID, inputsCopy, outputsCopy}
}

// 对每个输入签名。
func (tx *Transaction) Sign(privkey ecdsa.PrivateKey, prevTxs map[string]*Transaction) {
	// 如果是coinbase交易，就不用签名。
	if tx.IsCoinbase() {
		return
	}

	// 检查各输入的ID是否正确。
	for _, input := range tx.Inputs {
		if prevTxs[hex.EncodeToString(input.RefID)].ID == nil {
			panic("[Error] Previous transaction incorrect.")
		}
	}

	txCopy := tx.Deepcopy()
	for index, input := range txCopy.Inputs {
		prevTx := prevTxs[hex.EncodeToString(input.RefID)]
		txCopy.Inputs[index].Signature = nil
		txCopy.Inputs[index].Pubkey = prevTx.Outputs[input.RefIndex].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[index].Pubkey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privkey, txCopy.ID)
		if err != nil {
			panic(err)
		}
		tx.Inputs[index].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

// 检验交易输入的数字签名。
func (tx *Transaction) Verify(prevTxs map[string]*Transaction) bool {
	// 如果是coinbase交易，就不用验证。
	if tx.IsCoinbase() {
		return true
	}

	// 检查各输入的ID是否正确。
	for _, input := range tx.Inputs {
		if prevTxs[hex.EncodeToString(input.RefID)].ID == nil {
			panic("[Error] Previous transaction incorrect.")
		}
	}

	txCopy := tx.Deepcopy()
	curve := elliptic.P256()

	for index, input := range txCopy.Inputs {
		prevTx := prevTxs[hex.EncodeToString(input.RefID)]
		txCopy.Inputs[index].Signature = nil
		txCopy.Inputs[index].Pubkey = prevTx.Outputs[input.RefIndex].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[index].Pubkey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(input.Signature)
		r.SetBytes(input.Signature[:(sigLen / 2)])
		s.SetBytes(input.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(input.Pubkey)
		x.SetBytes(input.Pubkey[:(keyLen / 2)])
		y.SetBytes(input.Pubkey[(keyLen / 2):])

		rawPubkey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubkey, txCopy.ID, &r, &s) {
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

// 创建一笔coinbase交易。
func NewCoinbaseTX(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	// 没有输入，只有一个输出。
	txi := TXI{
		RefID:     []byte{},
		RefIndex:  -1,
		Signature: nil,
		Pubkey:    []byte(data),
	}
	txo := NewTXO(utils.Subsidy, to)
	tx := Transaction{
		ID:      nil,
		Inputs:  []TXI{txi},
		Outputs: []TXO{txo},
	}
	tx.ID = tx.Hash()
	return &tx
}
