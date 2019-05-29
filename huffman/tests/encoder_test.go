package tests

import (
	. "gohuffman/huffman"
	"testing"
)

type EncodedData []byte

func (encodedData *EncodedData) ByteReady(b byte) {
	*encodedData = append(*encodedData, b)
}

func TestEncodedDataIsShorterThenOriginal(t *testing.T) {
	encodedData := EncodedData{}

	data := []byte{1, 2, 3, 4, 5}

	frequencies := map[byte]int{}
	for _, b := range data {
		frequencies[byte(b)]++
	}

	huffmanTree := BuildHuffmanTree(frequencies)
	coding := BuildCodingFromTree(huffmanTree, nil)
	encoder := CreateEncoder(&encodedData, coding)
	for _, dataByte := range data {
		encoder.EncodeByte(dataByte)
	}
	encoder.Finalize()

	if len(encodedData) >= len(data) {
		t.Errorf("Encoded data is not shorter. Original Len: %d, Encoded Len: %d", len(data), len(encodedData))
	}
}

func TestDataConsistingOfSingleValueShouldBeEncodedToOneBit(t *testing.T) {
	encodedData := EncodedData{}

	var data []byte
	for i := 0; i < 500; i++ {
		data = append(data, 7, 200)
	}

	frequencies := map[byte]int{}
	for _, b := range data {
		frequencies[byte(b)]++
	}

	huffmanTree := BuildHuffmanTree(frequencies)
	coding := BuildCodingFromTree(huffmanTree, nil)
	encoder := CreateEncoder(&encodedData, coding)
	for _, dataByte := range data {
		encoder.EncodeByte(dataByte)
	}
	encoder.Finalize()

	// Since we have only two numbers, we expect it one will be encoded with '0' and one with '1'
	// 1 bit for each byte. Since we have 1000 numbers and each number expected to be encoded to 1/8 of byte (one bit)
	// total encoded array length should be 1000 / 8 = 125.

	if len(encodedData) != 125 {
		t.Errorf("Expected 125 but, got: %d", len(encodedData))
	}
}

func TestCaseWhenAllBytesHasSameProbabilityExceptOne(t *testing.T) {
	encodedData := EncodedData{}

	var data []byte
	for i := 0; i < 255; i++ {
		data = append(data, byte(i))
	}
	// Add 4 zeroes more (total is 5).
	data = append(data, 0, 0, 0, 0)

	frequencies := map[byte]int{}
	for _, b := range data {
		frequencies[byte(b)]++
	}

	huffmanTree := BuildHuffmanTree(frequencies)
	coding := BuildCodingFromTree(huffmanTree, nil)
	encoder := CreateEncoder(&encodedData, coding)
	for _, dataByte := range data {
		encoder.EncodeByte(dataByte)
	}
	encoder.Finalize()

	if len(encodedData) != len(data) {
		t.Errorf("Expected %d got %d", len(data), len(encodedData))
	}
}
