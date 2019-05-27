package huffman

import "bytes"

// BufferWriter implements IEncoderClient
type BufferWriter struct {
	buffer *bytes.Buffer
}

func (self BufferWriter) ByteReady(b byte) {
	self.buffer.WriteByte(b)
}
