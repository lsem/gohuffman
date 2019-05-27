package huffman

import "sort"

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

type DecodingTable struct {
	tableRecords PlainCodingDataRecordsCollection
}

func BuildDecodingTable(codingMap CodingMap) DecodingTable {
	newTable := DecodingTable{}
	newTable.tableRecords = make(PlainCodingDataRecordsCollection, 0)
	for k, v := range codingMap {
		newSequence := make([]byte, len(v))
		copy(newSequence, v)
		newRecord := PlainCodingDataRecord{symbol: k, sequence: newSequence}
		newTable.tableRecords = append(newTable.tableRecords, newRecord)
	}

	sort.Slice(newTable.tableRecords, func(i, j int) bool {
		return PlainCodingDataRecordLess(newTable.tableRecords[i], newTable.tableRecords[j])
	})

	return newTable
}

func (self DecodingTable) IndexOf(sequence []byte) int {
	return self.tableRecords.IndexOfSequence(sequence)
}

func (self DecodingTable) At(index int) PlainCodingDataRecord {
	return self.tableRecords[index]
}
