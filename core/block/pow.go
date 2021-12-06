package block

import (
	"blockchain/utils"
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

// 难度系数。
const difficulty = 24
const maxNonce = math.MaxInt64

// 工作量证明结构。
type pow struct {
	block  *Block   // 要证明工作量的区块。
	target *big.Int // 证明目标。
}

// 创建工作量证明事件。
func newPow(b *Block) *pow {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	return &pow{b, target}
}

// 计算当前区块内交易 ID 的哈希值。
func (b *Block) hashTxID() []byte {
	var IDs [][]byte

	for _, tx := range b.Transactions {
		IDs = append(IDs, tx.ID)
	}
	generalHash := sha256.Sum256(bytes.Join(IDs, []byte{}))

	return generalHash[:]
}

// 计算当前区块的哈希值。
func (p *pow) hashBlock(nonce int) [32]byte {
	blockBytes := bytes.Join(
		[][]byte{
			utils.Int64ToBytes(p.block.Timestamp),
			p.block.hashTxID(),
			p.block.PrevBlockHash,
			utils.Int64ToBytes(int64(nonce)),
			utils.Int64ToBytes(int64(difficulty)),
		},
		[]byte{},
	)

	return sha256.Sum256(blockBytes)
}

// 供自身判断工作量是否被证明。
func (p *pow) isProved(hash [32]byte) bool {
	hashInt := utils.BytesToBigInt(hash[:])
	return hashInt.Cmp(p.target) == -1
}

// 开始证明工作量。
func (p *pow) Run() {
	for nonce := 0; nonce < maxNonce; {
		hash := p.hashBlock(nonce)
		if p.isProved(hash) {
			p.block.Hash = hash[:]
			p.block.Nonce = nonce
			break
		}
		nonce++
	}
}