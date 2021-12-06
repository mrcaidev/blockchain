package transaction

import (
	"blockchain/utils"
	"bytes"
)

// 交易输出结构。
type TXOutput struct {
	Value      int    // 交易输出的币值。
	PubkeyHash []byte // 公钥哈希值。
}

// 创建输出。
func NewTXO(value int, address string) *TXOutput {
	hash := utils.Base58Decode([]byte(address))
	hash = hash[1 : len(hash)-4]
	txo := TXOutput{value, hash}
	return &txo
}

// 检验输出是否能被公钥解锁。
func (txo *TXOutput) IsUnlockableWith(pubkeyHash []byte) bool {
	return bytes.Equal(txo.PubkeyHash, pubkeyHash)
}
