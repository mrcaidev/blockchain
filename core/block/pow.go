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
func newPow(block *Block) *pow {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	return &pow{block, target}
}

// 计算当前区块内交易的哈希值。
func (block *Block) hashTX() []byte {
	var IDs [][]byte

	for _, TX := range block.Transactions {
		IDs = append(IDs, TX.ID)
	}
	generalHash := sha256.Sum256(bytes.Join(IDs, []byte{}))

	return generalHash[:]
}

// 计算当前区块的哈希值。
func (proof *pow) hashBlock(nonce int) [32]byte {
	blockBytes := bytes.Join(
		[][]byte{
			utils.IntToBytes(proof.block.Timestamp),
			proof.block.hashTX(),
			proof.block.PrevBlockHash,
			utils.IntToBytes(int64(nonce)),
			utils.IntToBytes(int64(difficulty)),
		},
		[]byte{},
	)

	return sha256.Sum256(blockBytes)
}

// 供自身判断工作量是否被证明。
func (proof *pow) isProved(hash [32]byte) bool {
	hashInt := utils.BytesToBigInt(hash[:])
	return hashInt.Cmp(proof.target) == -1
}

// 开始证明工作量。
func (proof *pow) Run() {
	for nonce := 0; nonce < maxNonce; {
		hash := proof.hashBlock(nonce)

		if proof.isProved(hash) {
			proof.block.Hash = hash[:]
			proof.block.Nonce = nonce
			break
		}
		nonce++
	}
}
