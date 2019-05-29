package huffman

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrorReadMagic     = errors.New("failed reading Magic")
	ErrorInvalidMagic  = errors.New("invalid Magic")
	ErrorReadTotalBits = errors.New("failed reading total bits count")
	ErrorReadBitsNum   = func(i int) error { return errors.New("Failed reading coding table, bits num for record:  " + string(i)) }
	ErrorReadBitsData  = func(i int) error { return errors.New("Failed reading coding table: bits data for record: " + string(i)) }
)

func DecodeStream(reader io.Reader) (data []byte, err error) {
	// Read Magic
	var buff [4]byte
	if _, err := io.ReadFull(reader, buff[:]); err != nil {
		return nil, ErrorReadMagic
	}
	if buff != FILE_MAGIC {
		return nil, ErrorInvalidMagic
	}

	var totalBitsCount uint64
	if err := binary.Read(reader, binary.BigEndian, &totalBitsCount); err != nil {
		return nil, ErrorReadTotalBits
	}

	var codingTable = make(CodingMap)
	// Read coding table
	for i := 0; i < 256; i++ {
		var bitsNumberBuff [1]byte
		if _, err := io.ReadFull(reader, bitsNumberBuff[:]); err != nil {
			return nil, ErrorReadBitsNum(i)
		}

		bitsNumber := bitsNumberBuff[0]
		if bitsNumber == 0 {
			// We no need empty sequences which means there is no going to be such byte in encoded sequence
			continue
		}

		var bitsData = make([]byte, bitsNumber)
		if _, err := io.ReadFull(reader, bitsData); err != nil {
			return nil, ErrorReadBitsData(i)
		}
		codingTable[byte(i)] = bitsData
	}

	decodingTable := BuildDecodingTable(codingTable)
	decodedBytes := make([]byte, 0)
	sliceWriter := SliceWritter{sliceWriteTo: &decodedBytes}
	decoder := CreateDecoder(&sliceWriter, decodingTable, totalBitsCount)

	done := false
	for !done {
		dataChunk := make([]byte, 1024)
		if bytesRead, err := io.ReadFull(reader, dataChunk); err != nil {
			if err == io.ErrUnexpectedEOF || err == io.EOF {
				dataChunk = dataChunk[:bytesRead]
				done = true
			} else {
				return nil, err
			}
		}
		for _, b := range dataChunk {
			decoder.DecodeByte(b)
		}
	}

	return decodedBytes, nil
}
