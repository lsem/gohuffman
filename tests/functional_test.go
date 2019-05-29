package tests

import (
	"fmt"
	"gohuffman/huffman"
	"math/rand"
	"os"
	"testing"
)

func TestEncodingDecodingRandomLengths(t *testing.T) {
	rand.Seed(1)

	lengths := []int{2, 3, 75, 127, 128, 129, 254, 255}

	if !testing.Short() {
		for i := 76; i < 127; i++ {
			lengths = append(lengths, i)
		}
		for i := 256; i < 2048; i++ {
			lengths = append(lengths, i)
		}
	}

	for _, length := range lengths {
		rand.Seed(1)
		fileName := fmt.Sprintf("TestEncodingDecodingRandomLengths_%d.dat", length);
		t.Run(fmt.Sprintf("%d", length), func(t *testing.T) {
			inData := make([]byte, length)
			for i := 0; i < len(inData); i++ {
				inData[i] = byte(rand.Intn(255))
			}

			huffman.Encode(inData, fileName)
			outData, err := huffman.DecodeFile(fileName)
			_ = os.Remove(fileName)

			if err != nil {
				t.Error("Failed decoding data")
			}

			if len(outData) != len(inData) {
				t.Errorf("Input and output data amount differ: %d vs %d\n", len(outData), len(inData))
				t.Fatal()
			}

			for bi := 0; bi < len(inData); bi++ {
				if inData[bi] != outData[bi] {
					t.Errorf("Diff at: %v\b", bi)
				}
			}
		})
	}

}
