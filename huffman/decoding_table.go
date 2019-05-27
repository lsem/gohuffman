package huffman

import "sort"

type DecodingTableRecord struct {
	symbol   byte
	sequence []byte
}

func Less(left, right DecodingTableRecord) bool {
	commonLen := Min(len(left.sequence), len(right.sequence))
	for idx := 0; idx < commonLen; idx++ {
		if left.sequence[idx] != right.sequence[idx] {
			return left.sequence[idx] < right.sequence[idx]
		}
	}
	return len(left.sequence) < len(right.sequence)
}

func Equal(left, right DecodingTableRecord) bool {
	return !Less(left, right) && !Less(right, left)
}

func LowerBound(records []DecodingTableRecord, x DecodingTableRecord) int {
	lo, hi := 0, len(records)
	for lo < hi {
		mid := (lo + hi) / 2
		xLessOrEqualThenMid := !Less(records[mid], x)
		if xLessOrEqualThenMid {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

type DecodingTable struct {
	tableRecords []DecodingTableRecord
}

func BuildDecodingTable(codingMap CodingMap) DecodingTable {
	newTable := DecodingTable{}
	newTable.tableRecords = make([]DecodingTableRecord, 0)
	for k, v := range codingMap {
		newSequence := make([]byte, len(v))
		copy(newSequence, v)
		newRecord := DecodingTableRecord{symbol: k, sequence: newSequence}
		newTable.tableRecords = append(newTable.tableRecords, newRecord)
	}

	sort.Slice(newTable.tableRecords, func(i, j int) bool {
		return Less(newTable.tableRecords[i], newTable.tableRecords[j])
	})

	return newTable
}

func (self DecodingTable) IndexOf(sequence []byte) int {
	record := DecodingTableRecord{symbol: 0, sequence: sequence}
	lb := LowerBound(self.tableRecords, record)
	if lb < len(self.tableRecords) && Equal(self.tableRecords[lb], record) {
		return lb
	}
	return -1
}

func (self DecodingTable) At(index int) DecodingTableRecord {
	return self.tableRecords[index]
}
