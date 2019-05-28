package huffman

type Decoder struct {
	decodingTable     DecodingTable
	client            IDecoderClient
	currentSequence   []byte
	totalBitsNumber   uint64
	totalBitsDecoded  uint64
	totalBytesDecoded int
}

type IDecoderClient interface {
	ByteDecoded(b byte)
}

func CreateDecoder(client IDecoderClient, decodingTable DecodingTable, totalBitsNumber uint64) Decoder {
	newDecoder := Decoder{}
	newDecoder.currentSequence = make([]byte, 0, 8)
	newDecoder.decodingTable = decodingTable
	newDecoder.client = client
	newDecoder.totalBitsNumber = totalBitsNumber
	newDecoder.totalBitsDecoded = 0
	return newDecoder
}

func (self *Decoder) DecodeByte(b byte) {
	// Last byte requires special handling so we should calculate how many
	// bits to decode in last byte.
	bitsToProcess := 8
	if (self.totalBitsNumber - self.totalBitsDecoded) < 8 {
		bitsToProcess = int(self.totalBitsNumber % 8)
	}
	for bitNum := 0; bitNum < bitsToProcess; bitNum++ {
		self.totalBitsDecoded++
		// Take MSB due to encoding order.
		msb := (b & 0x80) >> 7
		self.currentSequence = append(self.currentSequence, msb)
		b <<= 1
		index := self.decodingTable.IndexOf(self.currentSequence)
		if index != -1 {
			self.currentSequence = self.currentSequence[:0]
			self.client.ByteDecoded(self.decodingTable.At(index).symbol)
			self.totalBytesDecoded++
		}
	}
}
