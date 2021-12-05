package core

import (
	"blockchain/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

// 工作量证明结构。
type Pow struct {
	block  *Block
	target *big.Int
}

// 创建工作量证明事件。
func CreatePow(b *Block) *Pow {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-utils.Difficulty))
	return &Pow{b, target}
}

// 将区块数据与随机数拼接。
func (pow *Pow) joinBlockNonce(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			utils.IntToBytes(pow.block.Timestamp),
			pow.block.TransactionHash(),
			pow.block.PrevBlockHash,
			utils.IntToBytes(int64(nonce)),
			utils.IntToBytes(int64(utils.Difficulty)),
		},
		[]byte{},
	)
}

// 计算当前区块的哈希值。
func (pow *Pow) calcHash(nonce int) [32]byte {
	blockData := pow.joinBlockNonce(nonce)
	return sha256.Sum256(blockData)
}

// 判断工作量是否被证明，即当前哈希值是否小于目标值。
func (pow *Pow) isProved(hash [32]byte) bool {
	hashInt := utils.BytesToBigInt(hash[:])
	return hashInt.Cmp(pow.target) == -1
}

// 开始证明工作量。
func (pow *Pow) Run() {
	// 遍历寻找能让证明成功的随机数。
	for nonce := 0; nonce < utils.MaxNonce; {
		// 计算当前随机数对应的哈希值。
		result := pow.calcHash(nonce)
		// 如果结果比目标小，则工作量证明成功。
		if pow.isProved(result) {
			pow.block.Hash = result[:]
			pow.block.Nonce = nonce
			break
		}
		// 否则，使用下一个随机数重试。
		nonce++
	}
}

// 检验区块是否通过工作量证明。
func (pow *Pow) Validate() bool {
	hashInt := utils.BytesToBigInt(pow.block.Hash)
	return hashInt.Cmp(pow.target) == -1
}
