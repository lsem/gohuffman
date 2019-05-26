package huffman

func Assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
