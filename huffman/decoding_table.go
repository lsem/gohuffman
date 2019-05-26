package huffman


type PlainCodingDataRecord struct {
	symbol   byte
	sequence []byte
}

func plainCodingData_Less(left, right PlainCodingDataRecord) bool {
	commonLen := Min(len(left.sequence), len(right.sequence))
	for idx := 0; idx < commonLen; idx++ {
		if left.sequence[idx] != right.sequence[idx] {
			return left.sequence[idx] < right.sequence[idx]
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

