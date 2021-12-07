package merkle

import "crypto/sha256"

// Merkle 树结点结构。
type merkleNode struct {
	Left  *merkleNode // 左结点。
	Right *merkleNode // 右结点。
	Data  []byte      // 数据。
}

// 创建 Merkle 树结点。
func newMerkleNode(left, right *merkleNode, data []byte) *merkleNode {
	node := merkleNode{}

	// 生成该结点哈希值。
	var hash [32]byte
	if left == nil && right == nil {
		hash = sha256.Sum256(data)
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash = sha256.Sum256(prevHashes)
	}
	node.Data = hash[:]

	// 绑定左右结点。
	node.Left = left
	node.Right = right

	return &node
}
