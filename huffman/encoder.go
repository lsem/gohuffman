package huffman

import (
	"bytes"
)

type Encoder struct {
	currentByte     byte
	totalBitsCount  int // Total bits set count
	totalBytesCount int
	fastCodingTable [][]byte
	client          IEncoderClient
	finalized		bool
}

type IEncoderClient interface {
	ByteReady(b byte)
}

func CreateFastCodingTable(codingTable CodingMap) [][]byte {
	fastCodingTable := make([][]byte, 256)
	for bIndex := 0; bIndex < 256; bIndex++ {
		fastCodingTable[bIndex] = codingTable[byte(bIndex)]
	}
	return fastCodingTable
}

func CreateEncoder(client IEncoderClient, codingTable CodingMap) Encoder {
	return Encoder{fastCodingTable: CreateFastCodingTable(codingTable), client: client}
}

func (self* Encoder) EncodeByte(b byte) {
	Assert(!self.finalized, "Encoder already finalized and cannot be reused")
	for _, bit := range self.fastCodingTable[b] {
		self.currentByte |= bit
		self.totalBitsCount++
		if self.totalBitsCount%8 == 0 {
			self.client.ByteReady(self.currentByte)
			self.currentByte = 0
			self.totalBytesCount++
		} else { // make place for next bit
			self.currentByte <<= 1
		}
	}
}

func (self* Encoder) Finalize() {
	Assert(!self.finalized, "Encoder already finalized and cannot be reused")
	lasByteBitsProcessed := byte(self.totalBitsCount%8)
	if lasByteBitsProcessed > 0 {
		needsPadding := 8 - lasByteBitsProcessed
		for i := 0; i < int(needsPadding-1); i++ {
			self.currentByte |= 0
			self.currentByte <<= 1
		}
		self.client.ByteReady(self.currentByte)
		self.totalBytesCount++
	}
	self.finalized = true
}

// TODO: Move to separate file
type BufferWritingClient struct {
	buffer* bytes.Buffer
}

func (self BufferWritingClient) ByteReady(b byte) {
	self.buffer.WriteByte(b)
}
