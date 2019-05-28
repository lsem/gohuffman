package tests

import (
	"fmt"
	"gohuffman/huffman"
	"math/rand"
	"os"
	"testing"
)

func TestEncodingDecodingRandomLengths(t *testing.T) {
	const fileName = "TestEncodingDecodingRandomLengths.dat"

	rand.Seed(1)

	lengths := []byte{ /*1, 2, 3*/ 75, 127, 128, 129, 254, 255}

	for _, length := range lengths {

		t.Run(fmt.Sprintf("%d", length), func(t *testing.T) {
			randomData := make([]byte, length)
			for i := 0; i < len(randomData); i++ {
				randomData[i] = byte(rand.Intn(255))
			}

			huffman.Encode(randomData, fileName)
			decodedData, err := huffman.DecodeFile(fileName)
			_ = os.Remove(fileName)

			if err != nil {
				t.Error("Failed decoding data")
			}

			if len(decodedData) != len(randomData) {
				t.Errorf("Input and output data amount differ: %d vs %d\n", len(decodedData), len(randomData))
			}

			for bi := 0; bi < len(randomData); bi++ {
				if randomData[bi] != decodedData[bi] {
					t.Errorf("Diff at: %v\b", bi)
				}
			}
		})
	}

}
