package huffman

// SliceWritter implements IDecoderClient interface
type SliceWritter struct {
	sliceWriteTo *[]byte
}

func (self *SliceWritter) ByteDecoded(b byte) {
	*self.sliceWriteTo = append(*self.sliceWriteTo, b)
}
