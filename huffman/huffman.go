package huffman

import (
	"bytes"
	"errors"
	"io"
	"os"
)

var FILE_MAGIC = [4]byte{0x34, 0x89, 0x99, 0xff}

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
		Assert(len(bitsData) <= 8, "Must be not greater than 8 bits")
	}

	decodingTable := BuildDecodingTable(codingTable)
	decodedBytes := make([]byte, 0)
	sliceWritter := SliceWritter{sliceWriteTo: &decodedBytes}
	decoder := CreateDecoder(&sliceWritter, decodingTable)

	// Read data by chunks of 1024 (it would be nice to know length of data from header)
	dataProcessed := false
	for !dataProcessed {
		// Read data
		dataChunk := make([]byte, 1024)
		if bytesRead, err := io.ReadFull(inputFile, dataChunk); err != nil {
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

func frequenciesTableToMap(table [256]int) map[byte]int {
	frequenciesAsMap := make(map[byte]int)
	for i := 0; i < 256; i++ {
		if table[i] != 0 {
			frequenciesAsMap[byte(i)] = table[i]
		}
	}
	return frequenciesAsMap
}

func Encode(data []byte, fileName string) {
	// calculate frequencies for each of bytes
	// Since we are interested each byte of entire input dataset, map operations
	// are quite slow at this amount of data, so we use just table for this purposes.
	var frequenciesTable [256]int
	for _, x := range data {
		frequenciesTable[x]++
	}

	writeBuffer := new(bytes.Buffer)
	bufferWriter := BufferWriter{buffer: writeBuffer}

	coding := BuildCodingFromTree(BuildHuffmanTree(frequenciesTableToMap(frequenciesTable)), nil)

	encoder := CreateEncoder(bufferWriter, coding)
	for _, dataByte := range data {
		encoder.EncodeByte(dataByte)
	}
	encoder.Finalize()

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
		bits := coding[byte(i)]
		headerBuffer.WriteByte(byte(len(bits)))
		for _, bit := range bits {
			headerBuffer.WriteByte(bit)
		}
	}

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
}
