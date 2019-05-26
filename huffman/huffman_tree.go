package huffman

// HuffmanTreeNode represents data structure for storing one item
// of huffman tree which is just ordinary binary tree. If node's
// symbol is not nil that means it is leave node.
type HuffmanTreeNode struct {
	left, right *HuffmanTreeNode
	weight      uint32
	symbol      *byte // nil means node is not leave node.
}

func (self *HuffmanTreeNode) IsLeave() bool {
	return self.symbol != nil
}

func (self *HuffmanTreeNode) Serialize() {
	// TODO: Implement
}

func (self *HuffmanTreeNode) Height() int {
	if self == nil {
		return 0
	} else {
		return 1 + max(self.left.Height(), self.right.Height())
	}
}

func BuildHuffmanTree(frequencies map[byte]int) HuffmanTreeNode {
	println("Queue Prep")
	var queue HuffmanPriorityQueue
	var nodes []HuffmanTreeNode
	for k, v := range frequencies {
		Assert(v != 0, "There should be no zeroes in map")
		symbol := byte(k)
		nodes = append(nodes, HuffmanTreeNode{left: nil, right: nil,
			weight: uint32(v), symbol: &symbol})
	}
	queue.Init(nodes)

	println("Tree Prep")
	for queue.Len() > 1 {
		p, q := queue.Get(), queue.Get()
		queue.Add(HuffmanTreeNode{q, p, uint32(p.weight + q.weight), nil})
	}

	root := queue.Get()
	return *root
}
