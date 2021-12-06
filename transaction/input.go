package transaction

import (
	"blockchain/utils"
	"bytes"
)

// 交易输入结构。
type TXInput struct {
	RefID     []byte // 引用交易的ID。
	RefIndex  int    // 在引用交易输出中的索引。
	Signature []byte // 发起者的数字签名。
	Pubkey    []byte // 用于解锁的公钥。
}

func NewTXI(refID []byte, refIndex int, signature []byte, pubkey []byte) *TXInput {
	return &TXInput{
		RefID:     refID,
		RefIndex:  refIndex,
		Signature: signature,
		Pubkey:    pubkey,
	}
}

// 检验输入是否被指定公钥锁定。
func (txi *TXInput) UsesKey(pubkeyHash []byte) bool {
	lockingHash := utils.GetPubkeyHash(txi.Pubkey)
	return bytes.Equal(lockingHash, pubkeyHash)
}
