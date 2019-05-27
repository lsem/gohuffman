package huffman

import "os"

func DecodeFile(fileName string) (data []byte, err error) {
	inputFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return DecodeStream(inputFile)
}

func Encode(data []byte, fileName string) {
	// Create file
	outFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = EncodeStream(data, outFile)
	if err != nil {
		panic(err)
	}
}
