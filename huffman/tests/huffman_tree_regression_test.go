package tests

import (
	. "gohuffman/huffman"
	"io/ioutil"
	"math/rand"
	"testing"
)

const fileName = "regression_tests.huffmanTree.previous"

func SaveHuffmanTreeToFile(tree string) error {
	if err := ioutil.WriteFile(fileName, []byte(tree), 0644); err != nil {
		return err
	}
	return nil
}
func LoadHuffmanTreeFromFile() (string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func TestHuffmanTreeRegressions(t *testing.T) {
	// This test intended to catch errors because of non-deterministic
	// tree building. It uses current working directory for storing results
	// from previous run.

	frequencies := make(map[byte]int)

	rand.Seed(1)
	for i := byte(0); i < 255; i++ {
		frequencies[i] = rand.Intn(299) + 1
	}

	var previous, current string

	current = BuildHuffmanTree(frequencies).String()

	previous, err := LoadHuffmanTreeFromFile()
	if err != nil {
		t.Errorf("Failed loading previous file: %v. If it is first run, than just rerun.\n", err)
	}

	if current != previous {
		if testing.Verbose() {
			t.Errorf("Current and previous are different.\nCurrent:\n%v\n\nPrevious:\n %v\n", current, previous)
		} else {
			t.Errorf("Current and previous are different")
		}
	}

	if err := SaveHuffmanTreeToFile(current); err != nil {
		t.Errorf("Failed writing current file: %v", err)
	}

}
