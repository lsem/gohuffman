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




