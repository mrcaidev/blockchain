package transaction

import (
	"blockchain/utils"
	"bytes"
)

// 交易输出结构。
type TxOutput struct {
	Value      int    // 交易输出存储的价值。
	PubkeyHash []byte // 公钥哈希值。
}

// 创建交易输出。
func NewTxo(value int, address string) *TxOutput {
	payload := utils.Base58Decode([]byte(address))
	pubkeyHash := payload[1 : len(payload)-utils.ChecksumLen]

	txo := TxOutput{value, pubkeyHash}
	return &txo
}

// 检验交易输出是否能被指定公钥解锁。
// 即：判断交易输出内的公钥，与传入的指定公钥，是不是同一把。
func (txo *TxOutput) IsUnlockableWith(pubkeyHash []byte) bool {
	return bytes.Equal(txo.PubkeyHash, pubkeyHash)
}
