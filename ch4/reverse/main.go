package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	x := [5]int{1, 2, 3, 4, 5}
	reverse(&x)
	fmt.Printf("%v\n", x)

	y := []byte("a 界 b 世 c")
	fmt.Printf("%s\n", y)
	fmt.Printf("%s\n", reverseByteSlice(y))
}

func reverse(x *[5]int) {
	for i := 0; i < len(x)/2; i++ {
		x[i], x[len(x)-1-i] = x[len(x)-1-i], x[i]
	}
}

func reverseNaive(x []byte) []byte {
	for i := 0; i < len(x)/2; i++ {
		x[i], x[len(x)-1-i] = x[len(x)-1-i], x[i]
	}

	return x
}

// this works by reversing the runes of multiple bytes twice to keep them in the same order
func reverseByteSlice(x []byte) []byte {
	// reverse runes that are multiple bytes
	for i := 0; i < len(x); i++ {
		_, s := utf8.DecodeRune(x[i:])
		if s > 1 {
			copy(x[i:i+s], reverseNaive(x[i:i+s]))
			i += s
		}
	}

	// naive reverse
	x = reverseNaive(x)

	return x
}
