package transaction

import (
	"blockchain/utils"
	"bytes"
)

// 交易输出结构。
type TXO struct {
	Value      int    // 交易输出的币值。
	PubkeyHash []byte // 公钥哈希值。
}

// 锁定输出。
func (txo *TXO) Lock(address []byte) {
	pubkeyHash := utils.Base58Decode(address)
	txo.PubkeyHash = pubkeyHash[1 : len(pubkeyHash)-4]
}

// 检验输出是否能被公钥解锁。
func (txo *TXO) IsUnlockableWith(pubkeyHash []byte) bool {
	return bytes.Equal(txo.PubkeyHash, pubkeyHash)
}

// 创建输出。
func NewTXO(value int, address string) TXO {
	txo := TXO{value, nil}
	txo.Lock([]byte(address))
	return txo
}
