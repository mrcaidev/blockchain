package transaction

import (
	"blockchain/utils"
	"bytes"
)

// 交易输入结构。
type TXI struct {
	RefID     []byte // 引用交易的ID。
	RefIndex  int    // 在引用交易输出中的索引。
	Signature []byte // 发起者的数字签名。
	Pubkey    []byte // 用于解锁的公钥。
}

func (txi *TXI) UsesKey(pubkeyHash []byte) bool {
	lockingHash := utils.HashPubKey(txi.Pubkey)
	return bytes.Equal(lockingHash, pubkeyHash)
}
