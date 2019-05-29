package main

import (
	"fmt"
	"gohuffman/huffman"
	"io/ioutil"
	"log"
	"math/rand"
)

func compareData(arr1, arr2 []byte, n int) bool {
	if n == 0 {
		if len(arr1) != len(arr2) {
			fmt.Printf("Arrays sizes differ. Left is %v while right is %v\n", len(arr1), len(arr2))
			return false
		}
		n = len(arr1)
	}

	for i := 0; i < n; i++ {
		if arr1[i] != arr2[i] {
			fmt.Printf("Arrays differ at %v position\n", i)
			return false
		}
	}

	return true
}

func main() {

	data, err := ioutil.ReadFile("book.txt")
	if err != nil {
		panic(err)
	}

	println("Prep")
	//var data = make([]byte, 10000000)

	for i := 0; i < len(data); i++ {
		data[i] = byte(rand.Intn(76))
		if i < 30 {
			//fmt.Printf("%v: %v\n", i, data[i])
		}
	}

	println("Start")
	huffman.Encode(data, "out.huffman")

	decodedData, err := huffman.DecodeFile("out.huffman")
	if err != nil {
		log.Fatal(err)
		return
	}

	if !compareData(data, decodedData, 0) {
		//if !compareData(data, decodedData, 607391) {
		fmt.Println("Input and Output Arrays differ")
	} else {
		fmt.Println("Elements are the same!")
	}
}
