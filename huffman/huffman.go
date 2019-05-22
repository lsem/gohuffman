package huffman

import (
	"bytes"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

var FILE_MAGIC = [4]byte{0x34, 0x89, 0x99, 0xff}

// Node for HuffmanTree
type HuffmanTreeNode struct {
	left, right *HuffmanTreeNode
	weight      uint32
	symbol      *byte // nil means node is not leave node.
}

type HuffmanTreeNodePriorityQueue []*HuffmanTreeNode

type CodingMap map[byte][]byte

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

func Decode(fileName string) (data []byte, err error) {

	inputFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	// Read Magic
	var buff [4]byte
	if _, err := io.ReadFull(inputFile, buff[:]); err != nil {
		return nil, errors.New("Failed reading Magic")
		panic(err)
	}
	//log.Printf("magic: %v", buff)

	if buff != FILE_MAGIC {
		return nil, errors.New("Invalid Magic")
	}

	type PlainCodingDataRecord struct {
		symbol   byte
		sequence []byte
	}

	var plainCodingData = make([]PlainCodingDataRecord, 0, 256)

	var codingTable = make(CodingMap)
	// Read coding table
	for i := 0; i < 256; i++ {
		var bitsNumberBuff [1]byte
		if _, err := io.ReadFull(inputFile, bitsNumberBuff[:]); err != nil {
			return nil, errors.New("Failed reading coding table: bits number for record " + string(i))
		}

		bitsNumber := bitsNumberBuff[0]

		if bitsNumber == 0 {
			// We no need empty sequences which means there is no going to be such byte in encoded sequence
			continue
		}

		var bitsData = make([]byte, bitsNumber)
		if _, err := io.ReadFull(inputFile, bitsData); err != nil {
			return nil, errors.New("Failed reading coding table: bits data for record " + string(i))
		}
		codingTable[byte(i)] = bitsData

		if len(bitsData) > 8 {
			panic(errors.New("Must be not greater than 8 bits"))
		}

		bitsDataCopy := make([]byte, len(bitsData))
		copy(bitsDataCopy, bitsData)

		var plainRecord PlainCodingDataRecord
		plainRecord.symbol = byte(i)
		plainRecord.sequence = bitsDataCopy
		plainCodingData = append(plainCodingData, plainRecord)
	}

	// Read last byte and its size
	lastByteSizeAndByte := make([]byte, 2)
	if _, err := io.ReadFull(inputFile, lastByteSizeAndByte); err != nil {
		return nil, errors.New("Failed reading last byte information")
	}

	// Sort plain data records lexicographically
	plainCodingData_Less := func(left, right PlainCodingDataRecord) bool {
		commonPrefixLen := min(len(left.sequence), len(right.sequence))
		for idx := 0; idx < commonPrefixLen; idx++ {
			leftSequence, rightSequence := left.sequence, right.sequence
			if leftSequence[idx] != rightSequence[idx] {
				return leftSequence[idx] < rightSequence[idx]
			}
		}
		return len(left.sequence) < len(right.sequence)
	}
	plainCodingData_Equal := func(left, right PlainCodingDataRecord) bool {
		return !plainCodingData_Less(left, right) && !plainCodingData_Less(right, left)
	}
	sort.Slice(plainCodingData, func(i, j int) bool {
		return plainCodingData_Less(plainCodingData[i], plainCodingData[j])
	})
	//fmt.Println(plainCodingData)

	indexOfSequence := func(plainCodingData []PlainCodingDataRecord, sequence []byte) int {
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

	// Build reverse coding table. We can consider building prefix tree for this purposes
	// so we can navigate bit by bit. (todo)

	decodedBytes := make([]byte, 0)

	// Read data by chunks of 1024 (it would be nice to know length of data from header)
	currentSequence := make([]byte, 0, 8)
	dataProceessed := false
	for !dataProceessed {
		// Read data
		dataChunk := make([]byte, 1024)
		if bytesRead, err := io.ReadFull(inputFile, dataChunk); err != nil {
			if err == io.ErrUnexpectedEOF {
				dataChunk = dataChunk[:bytesRead]
				// todo: here last byte might require special handling taking into account its size
				dataChunk = append(dataChunk, lastByteSizeAndByte[1])
				dataProceessed = true
			} else {
				panic(err)
			}
		}
		// Process data
		for _, byte := range dataChunk {
			for bitNum := 0; bitNum < 8; bitNum++ {
				// Take MSB due to encoding order.
				msb := (byte & 0x80) >> 7
				currentSequence = append(currentSequence, msb)
				byte <<= 1
				index := indexOfSequence(plainCodingData, currentSequence)
				if index != -1 {
					// match
					//fmt.Printf("Match(index: %v, byte: %v,  sequence: %v)\n", index, plainCodingData[index].symbol, currentSequence)
					currentSequence = currentSequence[:0]
					decodedBytes = append(decodedBytes, plainCodingData[index].symbol)
				}
			}
		}
	}
	return decodedBytes, nil
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func treeHeight(node *HuffmanTreeNode) int {
	if node == nil {
		return 0
	} else {
		return 1 + max(treeHeight(node.left), treeHeight(node.right))
	}
}

func Encode(data []byte, fileName string) {
	// calculate frequencies for each of bytes
	println("Freq")

	var frequenciesTable [256]int
	for _, x := range data {
		frequenciesTable[x]++
	}
	frequencies := make(map[byte]int)
	for i := 0; i < 256; i++ {
		frequencies[byte(i)] = frequenciesTable[i]
	}

	// create tree nodes and push all them to the priority queue
	println("Queue Prep")
	var queue HuffmanTreeNodePriorityQueue
	for k, v := range frequencies {
		if v == 0 {
			// this would not exist in original frequencies map.
			continue
		}
		symbol := k
		newNode := HuffmanTreeNode{nil, nil, uint32(v), &symbol}
		queue = append(queue, &newNode)
	}
	heap.Init(&queue)

	// build Huffman tree
	println("Tree Prep")
	var lastNode *HuffmanTreeNode = nil
	for len(queue) > 1 {
		p := heap.Pop(&queue).(*HuffmanTreeNode)
		q := heap.Pop(&queue).(*HuffmanTreeNode)
		if p.weight == 0 {
			//panic(errors.New("pizdec"));
		}
		// todo: does it matter left is p or left is q?
		newNode := HuffmanTreeNode{q, p, uint32(p.weight + q.weight), nil}
		if newNode.weight == 0 {
			//panic(errors.New("pizdec"));
		}
		heap.Push(&queue, &newNode)
		lastNode = &newNode
	}
	//fmt.Printf("Tree height: %v\n", treeHeight(lastNode))

	println("Build Coding")
	// having Huffman tree, create coding where keys contain symbols from
	// file and values corresponding sequence of 0 and 1 for given symbol
	codingTable := buildCodingFromTree(*lastNode, make([]byte, 0))

	// Build fast coding table optimized for numerous looks up
	fastCodingTable := make([][]byte, 256)
	for bIndex := 0; bIndex < 256; bIndex++ {
		fastCodingTable[bIndex] = codingTable[byte(bIndex)]
	}

	var writeBuffer bytes.Buffer

	// Encode data
	println("Encode")
	var totalBytesCount uint32
	var currentByte byte
	var bitsSetCount uint32
	for _, dataByte := range data {
		for _, bit := range fastCodingTable[dataByte] {
			currentByte |= bit
			bitsSetCount++
			if bitsSetCount%8 == 0 {
				//fmt.Printf("%3d: %b\n", totalBytesCount, currentByte)
				writeBuffer.WriteByte(currentByte)
				currentByte = 0
				totalBytesCount++
			} else { // make place for next bit
				currentByte <<= 1
			}
		}
	}

	// Process last byte; pad zeroes on right side
	var lastByte byte
	var lastByteBitsSize byte

	lasByteBitsProcessed := bitsSetCount % 8
	lastByteBitsSize = byte(lasByteBitsProcessed)
	if lasByteBitsProcessed > 0 {
		needsPadding := 8 - lasByteBitsProcessed
		for i := 0; i < int(needsPadding-1); i++ {
			currentByte |= 0
			currentByte <<= 1
		}
		lastByte = currentByte
		totalBytesCount++
	}

	var headerBuffer bytes.Buffer

	// Write magic
	_, err := headerBuffer.Write(FILE_MAGIC[:])
	if err != nil {
		panic(err)
	}

	// Write coding table. Table consists of 256 records, each record has next format:
	//	1 byte: length
	//  [length] bytes of data.
	for i := 0; i < 256; i++ {
		bits := codingTable[byte(i)]
		headerBuffer.WriteByte(byte(len(bits)))
		for _, bit := range bits {
			headerBuffer.WriteByte(bit)
		}
	}

	// Then leave space for last byte, 1 byte for bit-length and one for data.
	headerBuffer.WriteByte(lastByteBitsSize)
	headerBuffer.WriteByte(lastByte)

	// Create file
	outFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// write header
	_, err = headerBuffer.WriteTo(outFile)
	if err != nil {
		panic(err)
	}

	// write data
	_, err = writeBuffer.WriteTo(outFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Total bytes count: %d. Compress ratio: %v%%\n", totalBytesCount, int(100*float64(totalBytesCount)/float64(len(data))))
}
