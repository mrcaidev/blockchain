package block

import (
	"blockchain/core/merkle"
	"blockchain/utils"
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

// 难度系数。
const difficulty = 8
const maxNonce = math.MaxInt64

// 工作量证明的目标。
var target = big.NewInt(0).Lsh(big.NewInt(1), uint(256-difficulty))

// 获取区块内交易的 Merkle 树根结点值。
func (b *Block) hashTx() []byte {
	var txs [][]byte
	for _, tx := range b.Transactions {
		txs = append(txs, tx.Serialize())
	}
	tree := merkle.NewMerkleTree(txs)
	return tree.Root.Data
}

// 获取区块的哈希值。
func (b *Block) hash(nonce int) [32]byte {
	blockBytes := bytes.Join(
		[][]byte{
			utils.Int64ToBytes(b.Timestamp),
			b.hashTx(),
			b.PrevBlockHash,
			utils.Int64ToBytes(int64(nonce)),
			utils.Int64ToBytes(int64(difficulty)),
		},
		[]byte{},
	)

	return sha256.Sum256(blockBytes)
}

// 判断工作量是否被证明。
func isProved(hash [32]byte) bool {
	return utils.BytesToBigInt(hash[:]).Cmp(target) == -1
}

// 开始证明工作量。
func (b *Block) proofWork() {
	for nonce := 0; nonce < maxNonce; {
		hash := b.hash(nonce)
		if isProved(hash) {
			b.Hash = hash[:]
			b.Nonce = nonce
			break
		}
		nonce++
	}
}
