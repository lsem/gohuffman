package huffman

import "container/heap"

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

type HuffmanPriorityQueue struct {
	container HuffmanTreeNodePriorityQueue
}

func (self *HuffmanPriorityQueue) Init(nodes []HuffmanTreeNode) {
	for _, node := range nodes {
		n := node // we must copy, otherwise pointer to local node will be added
		self.container = append(self.container, &n)
	}
	heap.Init(&self.container)
}

func (self *HuffmanPriorityQueue) Add(node HuffmanTreeNode) {
	heap.Push(&self.container, &node)
}

func (self *HuffmanPriorityQueue) Get() *HuffmanTreeNode {
	return heap.Pop(&self.container).(*HuffmanTreeNode)
}

func (self HuffmanPriorityQueue) Len() int {
	return len(self.container)
}
