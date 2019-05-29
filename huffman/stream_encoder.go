package huffman

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func frequenciesTableToMap(table [256]int) map[byte]int {
	frequenciesAsMap := make(map[byte]int)
	for i := 0; i < 256; i++ {
		if table[i] != 0 {
			frequenciesAsMap[byte(i)] = table[i]
		}
	}
	return frequenciesAsMap
}

func EncodeStream(data []byte, writer io.Writer) (err error) {
	// calculate frequencies for each of bytes
	// Since we are interested each byte of entire input dataset, map operations
	// are quite slow at this amount of data, so we use just table for this purposes.
	var frequenciesTable [256]int
	for _, x := range data {
		frequenciesTable[x]++
	}

	writeBuffer := new(bytes.Buffer)
	bufferWriter := BufferWriter{buffer: writeBuffer}

	huffmanTree := BuildHuffmanTree(frequenciesTableToMap(frequenciesTable))
	coding := BuildCodingFromTree(huffmanTree, nil)

	encoder := CreateEncoder(bufferWriter, coding)
	for _, dataByte := range data {
		encoder.EncodeByte(dataByte)
	}
	encoder.Finalize()

	var headerBuffer bytes.Buffer

	// Write magic
	bytesWritten, err := headerBuffer.Write(FILE_MAGIC[:])
	if err != nil {
		return err
	}
	if bytesWritten != len(FILE_MAGIC[:]) {
		return errors.New("failed writing magic")
	}

	err = binary.Write(&headerBuffer, binary.BigEndian, uint64(encoder.totalBitsCount))
	if err != nil {
		return errors.New("failed writing total bits count")
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

	// TODO: move out header into separate encode method so we can test it independently

	// write header
	_, err = headerBuffer.WriteTo(writer)
	if err != nil {
		panic(err)
	}

	// write data
	_, err = writeBuffer.WriteTo(writer)
	if err != nil {
		panic(err)
	}

	return nil
}
