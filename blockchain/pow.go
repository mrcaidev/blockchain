package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

// 难度系数。
const difficulty = 24
const maxNonce = math.MaxInt64

// 工作量证明结构。
type Pow struct {
	block  *Block
	target *big.Int
}

// 创建工作量证明事件。
func CreatePow(b *Block) *Pow {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	return &Pow{b, target}
}

// 将区块数据与随机数拼接。
func (pow *Pow) joinBlockNonce(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			IntToBytes(pow.block.Timestamp),
			IntToBytes(int64(nonce)),
			IntToBytes(int64(difficulty)),
			pow.block.Data,
		},
		[]byte{},
	)
}

// 计算当前区块的SHA值。
func (pow *Pow) calcSHA(nonce int) [32]byte {
	blockData := pow.joinBlockNonce(nonce)
	return sha256.Sum256(blockData)
}

// 判断工作量是否被证明，即当前SHA值是否小于目标值。
func (pow *Pow) isProved(sha [32]byte) bool {
	shaInt := BytesToBigInt(sha[:])
	return shaInt.Cmp(pow.target) == -1
}

// 开始证明工作量。
func (pow *Pow) Run() {
	// 遍历寻找能让证明成功的随机数。
	for nonce := 0; nonce < maxNonce; {
		// 计算当前随机数对应的SHA值。
		result := pow.calcSHA(nonce)
		// 如果结果比目标小，则工作量证明成功。
		if pow.isProved(result) {
			pow.block.ID = result[:]
			pow.block.Nonce = nonce
			break
		}
		// 否则，使用下一个随机数重试。
		nonce++
	}
}

// 检验区块是否通过工作量证明。
func (pow *Pow) Validate() bool {
	shaInt := BytesToBigInt(pow.block.ID)
	return shaInt.Cmp(pow.target) == -1
}
