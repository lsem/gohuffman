package huffman

type CodingMap map[byte][]byte

func cloneAndExtend(slice []byte, nums ...byte) []byte {
	sliceCopy := make([]byte, len(slice))
	copy(sliceCopy, slice)
	sliceCopy = append(sliceCopy, nums...)
	return sliceCopy
}

func BuildCodingFromTree(node TreeNode, coding []byte) CodingMap {
	if node.IsLeave() {
		return CodingMap{*node.symbol: coding}
	}

	var leftCoding, rightCoding CodingMap
	if node.left != nil {
		leftCoding = BuildCodingFromTree(*node.left,
			cloneAndExtend(coding, 0))
	}
	if node.right != nil {
		rightCoding = BuildCodingFromTree(*node.right,
			cloneAndExtend(coding, 1))
	}

	// merge left to right so right is union of both subtrees
	for symbol, coding := range leftCoding {
		rightCoding[symbol] = coding
	}
	return rightCoding
}
