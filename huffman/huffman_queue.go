package huffman

type HuffmanTreeNodePriorityQueue []*HuffmanTreeNode

func (self HuffmanTreeNodePriorityQueue) Len() int { return len(self) }

func (self HuffmanTreeNodePriorityQueue) Less(i, j int) bool {
	return self[i].weight < self[j].weight
}

func (self *HuffmanTreeNodePriorityQueue) Pop() interface{} {
	old := *self
	n := len(old)
	item := old[n-1]
	*self = old[0 : n-1]
	return item
}

func (self *HuffmanTreeNodePriorityQueue) Push(x interface{}) {
	item := x.(*HuffmanTreeNode)
	*self = append(*self, item)
}

func (self HuffmanTreeNodePriorityQueue) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
