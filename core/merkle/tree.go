package merkle

type merkleTree struct {
	root *merkleNode
}

func NewMerkleTree(dataList [][]byte) *merkleTree {
	var baseNodes []merkleNode

	// 如果有奇数份数据，就复制最后一份数据。
	if len(dataList)%2 != 0 {
		dataList = append(dataList, dataList[len(dataList)-1])
	}

	// 为每一份数据创建其 Merkle 树底层结点。
	for _, data := range dataList {
		node := newMerkleNode(nil, nil, data)
		baseNodes = append(baseNodes, *node)
	}

	// 逐层创建上层结点。
	for level := 0; level < len(dataList)/2; level++ {
		var newLevel []merkleNode

		// 将当前底层结点两两合并。
		for index := 0; index < len(baseNodes); index += 2 {
			node := newMerkleNode(&baseNodes[index], &baseNodes[index+1], nil)
			newLevel = append(newLevel, *node)
		}
		baseNodes = newLevel
	}

	return &merkleTree{&baseNodes[0]}
}
