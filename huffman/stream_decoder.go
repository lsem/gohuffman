package huffman

import (
	"encoding/binary"
	"errors"
	"io"
)

//var FILE_MAGIC = [4]byte{0x34, 0x89, 0x99, 0xff}

func DecodeStream(reader io.Reader) (data []byte, err error) {
	// Read Magic
	var buff [4]byte
	if _, err := io.ReadFull(reader, buff[:]); err != nil {
		return nil, errors.New("failed reading Magic")
	}
	if buff != FILE_MAGIC {
		return nil, errors.New("invalid Magic")
	}

	var totalBitsCount uint64
	if err := binary.Read(reader, binary.BigEndian, &totalBitsCount); err != nil {
		return nil, errors.New("failed reading total bits count")
	}

	var codingTable = make(CodingMap)
	// Read coding table
	for i := 0; i < 256; i++ {
		var bitsNumberBuff [1]byte
		if _, err := io.ReadFull(reader, bitsNumberBuff[:]); err != nil {
			return nil, errors.New("Failed reading coding table: bits number for record " + string(i))
		}

		bitsNumber := bitsNumberBuff[0]
		if bitsNumber == 0 {
			// We no need empty sequences which means there is no going to be such byte in encoded sequence
			continue
		}

		var bitsData = make([]byte, bitsNumber)
		if _, err := io.ReadFull(reader, bitsData); err != nil {
			return nil, errors.New("Failed reading coding table: bits data for record " + string(i))
		}
		codingTable[byte(i)] = bitsData
		Assert(len(bitsData) <= 8, "Must be not greater than 8 bits")
	}

	decodingTable := BuildDecodingTable(codingTable)
	decodedBytes := make([]byte, 0)
	sliceWritter := SliceWritter{sliceWriteTo: &decodedBytes}
	decoder := CreateDecoder(&sliceWritter, decodingTable, totalBitsCount)

	// Read data by chunks of 1024 (it would be nice to know length of data from header)
	dataProcessed := false
	for !dataProcessed {
		// Read data
		dataChunk := make([]byte, 1024)
		if bytesRead, err := io.ReadFull(reader, dataChunk); err != nil {
			if err == io.ErrUnexpectedEOF {
				dataChunk = dataChunk[:bytesRead]
				// todo: here last byte might require special handling taking into account its size
				//dataChunk = append(dataChunk, lastByteSizeAndByte[1])
				dataProcessed = true
			} else if err == io.EOF {
				dataChunk = dataChunk[:bytesRead]
				dataProcessed = true
			} else {
				panic(err)
			}
		}
		// Process data
		for _, b := range dataChunk {
			decoder.DecodeByte(b)
		}
	}

	return decodedBytes, nil
}
