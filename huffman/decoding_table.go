package huffman

type PlainCodingDataRecord struct {
	symbol   byte
	sequence []byte
}

type PlainCodingDataRecordsCollection []PlainCodingDataRecord

func PlainCodingDataRecordLess(left, right PlainCodingDataRecord) bool {
	commonLen := Min(len(left.sequence), len(right.sequence))
	for idx := 0; idx < commonLen; idx++ {
		if left.sequence[idx] != right.sequence[idx] {
			return left.sequence[idx] < right.sequence[idx]
		}
	}
	return len(left.sequence) < len(right.sequence)
}

func PlainCodingDataRecordEqual(left, right PlainCodingDataRecord) bool {
	return !PlainCodingDataRecordLess(left, right) && !PlainCodingDataRecordLess(right, left)
}

func (self PlainCodingDataRecordsCollection) IndexOfSequence(sequence []byte) int {
	fixtureRecord := PlainCodingDataRecord{symbol: 0, sequence: sequence}
	lb := self.LowerBound(fixtureRecord)
	if lb < len(self) && PlainCodingDataRecordEqual(self[lb], fixtureRecord) {
		return lb
	}
	return -1
}

func (self PlainCodingDataRecordsCollection) LowerBound(x PlainCodingDataRecord) int {
	lo, hi := 0, len(self)
	for lo < hi {
		mid := (lo + hi) / 2
		xLessOrEqualThenMid := !PlainCodingDataRecordLess(self[mid], x)
		if xLessOrEqualThenMid {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}
