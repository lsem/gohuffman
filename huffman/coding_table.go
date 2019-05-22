package huffman

///////////////////////////////////////////////////////////////////////////////
// Types
///////////////////////////////////////////////////////////////////////////////

// Node for HuffmanTree
type HuffmanTreeNode struct {
	left, right *HuffmanTreeNode
	weight      uint32
	symbol      *byte // nil means node is not leave node.
}

type HuffmanTreeNodePriorityQueue []*HuffmanTreeNode

type CodingMap map[byte][]byte

///////////////////////////////////////////////////////////////////////////////

type PlainCodingDataRecord struct {
	symbol   byte
	sequence []byte
}

func plainCodingData_Less(left, right PlainCodingDataRecord) bool {
	commonPrefixLen := min(len(left.sequence), len(right.sequence))
	for idx := 0; idx < commonPrefixLen; idx++ {
		leftSequence, rightSequence := left.sequence, right.sequence
		if leftSequence[idx] != rightSequence[idx] {
			return leftSequence[idx] < rightSequence[idx]
		}
	}
	return len(left.sequence) < len(right.sequence)
}

func plainCodingData_Equal(left, right PlainCodingDataRecord) bool {
	return !plainCodingData_Less(left, right) && !plainCodingData_Less(right, left)
}

type PlainCodingDataRecordsCollection []PlainCodingDataRecord

func (self PlainCodingDataRecordsCollection) IndexOfSequence(sequence []byte) int {
	return indexOfSequenceImpl(self, sequence)
}

func indexOfSequenceImpl(plainCodingData PlainCodingDataRecordsCollection, sequence []byte) int {
	lowerBound := func(x PlainCodingDataRecord) int {
		lo, hi := 0, len(plainCodingData)
		for lo < hi {
			mid := (lo + hi) / 2
			xLessOrEqualThenMid := !plainCodingData_Less(plainCodingData[mid], x)
			if xLessOrEqualThenMid {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		return lo
	}

	var fixtureRecord PlainCodingDataRecord
	fixtureRecord.sequence = sequence
	fixtureRecord.symbol = 0
	lb := lowerBound(fixtureRecord)
	if lb < len(plainCodingData) && plainCodingData_Equal(plainCodingData[lb], fixtureRecord) {
		// found sequence
		return lb
	} else {
		return -1
	}
}

///////////////////////////////////////////////////////////////////////////////
// Heap interface
///////////////////////////////////////////////////////////////////////////////

func (pq HuffmanTreeNodePriorityQueue) Len() int { return len(pq) }
func (pq HuffmanTreeNodePriorityQueue) Less(i, j int) bool {
	return pq[i].weight < pq[j].weight
}
func (pq *HuffmanTreeNodePriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (pq *HuffmanTreeNodePriorityQueue) Push(x interface{}) {
	item := x.(*HuffmanTreeNode)
	*pq = append(*pq, item)
}

func (pq HuffmanTreeNodePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func buildCodingFromTree(node HuffmanTreeNode, coding []byte) CodingMap {
	// check if this node is terminal or not
	if node.symbol != nil {
		// reached terminal node, coding slice contains result.
		return CodingMap{*node.symbol: coding}
	} else {
		var leftCoding, rightCoding CodingMap
		if node.left != nil {
			tmpCoding := make([]byte, len(coding))
			copy(tmpCoding, coding)
			tmpCoding = append(tmpCoding, 0)

			leftCoding = buildCodingFromTree(*node.left, tmpCoding)
		}
		if node.right != nil {
			tmpCoding := make([]byte, len(coding))
			copy(tmpCoding, coding)
			tmpCoding = append(tmpCoding, 1)

			rightCoding = buildCodingFromTree(*node.right, tmpCoding)
		}
		// merge left to right so right is union of both subtrees
		for symbol, coding := range leftCoding {
			rightCoding[symbol] = coding
		}
		return rightCoding
	}
}
