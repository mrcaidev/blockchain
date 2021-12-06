package transaction

import (
	"blockchain/utils"
	"bytes"
)

// 交易输入结构。
type TXInput struct {
	RefID     []byte // 引用输出所属交易的 ID。
	RefIndex  int    // 引用输出在上一笔交易所有输出中的索引。
	Signature []byte // 发起者的数字签名。
	Pubkey    []byte // 用于锁定的公钥。
}

// 创建交易输入。
func NewTXI(refID []byte, refIndex int, signature []byte, pubkey []byte) *TXInput {
	return &TXInput{
		RefID:     refID,
		RefIndex:  refIndex,
		Signature: signature,
		Pubkey:    pubkey,
	}
}

// 检验交易输入是否被指定公钥锁定。
// 即：判断交易输入内的公钥，与传入的指定公钥，是不是同一把。
func (txi *TXInput) IsLockedWith(pubkeyHash []byte) bool {
	lockingHash := utils.GetPubkeyHash(txi.Pubkey)
	return bytes.Equal(lockingHash, pubkeyHash)
}
