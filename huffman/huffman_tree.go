package huffman

// TreeNode represents data structure for storing one item
// of huffman tree which is just ordinary binary tree. If node's
// symbol is not nil that means it is leave node.
type TreeNode struct {
	left, right *TreeNode
	weight      uint32
	symbol      *byte // nil means node is not leave node.
}

func (self *TreeNode) IsLeave() bool {
	return self.symbol != nil
}

func (self *TreeNode) Serialize() {
	// TODO: Implement
}

func (self *TreeNode) Height() int {
	if self == nil {
		return 0
	} else {
		return 1 + Max(self.left.Height(), self.right.Height())
	}
}

func BuildHuffmanTree(frequencies map[byte]int) TreeNode {
	var nodes []*TreeNode
	for k, v := range frequencies {
		Assert(v != 0, "There should be no zeroes in map")
		symbol := byte(k)
		nodes = append(nodes,
			&TreeNode{left: nil, right: nil, weight: uint32(v), symbol: &symbol})
	}

	var queue PriorityQueue
	queue.Init(nodes)

	for queue.Len() > 1 {
		p, q := queue.Get(), queue.Get()
		queue.Add(TreeNode{q, p, uint32(p.weight + q.weight), nil})
	}

	root := queue.Get()
	return *root
}
