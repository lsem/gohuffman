package huffman

import (
	"fmt"
	"sort"
)

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

type VisitorFunc func(TreeNode)

func preorderTree(node TreeNode, visitor VisitorFunc) {
	visitor(node)
	if node.left != nil {
		preorderTree(*node.left, visitor)
	}
	if node.right != nil {
		preorderTree(*node.right, visitor)
	}
}

func (self TreeNode) String() string {
	// TODO: Implement
	result := ""
	sep := ""
	visitor := func(node TreeNode) {
		symbolStr := ""
		if node.symbol != nil {
			symbolStr = fmt.Sprint(*node.symbol)
		}
		result += fmt.Sprintf("%v %v:%v", sep, node.weight, symbolStr)
		sep = ", "
	}
	preorderTree(self, visitor)
	return result
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

	// This sort is needed to have deterministic order. Since queue.Init() will establish
	// heap invariant by reordering which is some sort of stable.
	sort.Slice(nodes, func(i, j int) bool {
		return *nodes[i].symbol < *nodes[j].symbol
	})

	var queue PriorityQueue
	queue.Init(nodes)

	for queue.Len() > 1 {
		p, q := queue.Get(), queue.Get()
		queue.Add(TreeNode{q, p, uint32(p.weight + q.weight), nil})
	}

	root := queue.Get()
	return *root
}
