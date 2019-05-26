package huffman

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sort"
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

	var plainCodingData = make(PlainCodingDataRecordsCollection, 0, 256)

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

	//// Read last byte and its size
	//lastByteSizeAndByte := make([]byte, 2)
	//if _, err := io.ReadFull(inputFile, lastByteSizeAndByte); err != nil {
	//	return nil, errors.New("Failed reading last byte information")
	//}

	sort.Slice(plainCodingData, func(i, j int) bool {
		return PlainCodingDataRecordLess(plainCodingData[i], plainCodingData[j])
	})
	//fmt.Println(plainCodingData)

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
				//dataChunk = append(dataChunk, lastByteSizeAndByte[1])
				dataProceessed = true
			} else if err == io.EOF {
				dataChunk = dataChunk[:bytesRead]
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
				index := plainCodingData.IndexOfSequence(currentSequence)
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
	bufferWriter := BufferWritingClient{buffer: writeBuffer}

	coding := BuildCodingFromTree( BuildHuffmanTree(
		frequenciesTableToMap(frequenciesTable)), make([]byte, 0))

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
