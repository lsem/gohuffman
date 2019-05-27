package huffman

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRandomLengths(t *testing.T) {
	rand.Seed(1)

	lengths := []byte{75, 255, 127, 128, 129 }

	for _, length := range lengths {

		t.Run(fmt.Sprintf("%d", length), func(t*testing.T) {
			randomData := make([]byte, length)
			for i := 0; i < len(randomData); i++ {
				//randomData[i] = byte((i + 1) * 31 % 255)
				randomData[i] = byte(rand.Intn(255))
			}

			Encode(randomData, "/tmp/TestRandomLengths.dat")
			decodedData, err := DecodeFile("/tmp/TestRandomLengths.dat")

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

