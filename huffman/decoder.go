package huffman

type Decoder struct {
	decodingTable   DecodingTable
	client          IDecoderClient
	currentSequence []byte
}

type IDecoderClient interface {
	ByteDecoded(b byte)
}

func CreateDecoder(client IDecoderClient, decodingTable DecodingTable) Decoder {
	newDecoder := Decoder{}
	newDecoder.currentSequence = make([]byte, 0, 8)
	newDecoder.decodingTable = decodingTable
	newDecoder.client = client
	return newDecoder
}

func (self *Decoder) DecodeByte(b byte) {
	for bitNum := 0; bitNum < 8; bitNum++ {
		// Take MSB due to encoding order.
		msb := (b & 0x80) >> 7
		self.currentSequence = append(self.currentSequence, msb)
		b <<= 1
		index := self.decodingTable.IndexOf(self.currentSequence)
		if index != -1 {
			self.currentSequence = self.currentSequence[:0]
			self.client.ByteDecoded(self.decodingTable.At(index).symbol)
		}
	}
}
