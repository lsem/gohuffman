package huffman

import "container/heap"

type TreeNodePriorityQueue []*TreeNode

func (self TreeNodePriorityQueue) Len() int { return len(self) }

func (self TreeNodePriorityQueue) Less(i, j int) bool {
	return self[i].weight < self[j].weight
}

func (self *TreeNodePriorityQueue) Pop() interface{} {
	old := *self
	n := len(old)
	item := old[n-1]
	*self = old[0 : n-1]
	return item
}

func (self *TreeNodePriorityQueue) Push(x interface{}) {
	item := x.(*TreeNode)
	*self = append(*self, item)
}

func (self TreeNodePriorityQueue) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type PriorityQueue struct {
	container TreeNodePriorityQueue
}

func (self *PriorityQueue) Init(nodes []*TreeNode) {
	self.container = make([]*TreeNode, len(nodes))
	copy(self.container, nodes)
	heap.Init(&self.container)
}

func (self *PriorityQueue) Add(node TreeNode) {
	heap.Push(&self.container, &node)
}

func (self *PriorityQueue) Get() *TreeNode {
	return heap.Pop(&self.container).(*TreeNode)
}

func (self PriorityQueue) Len() int {
	return len(self.container)
}
