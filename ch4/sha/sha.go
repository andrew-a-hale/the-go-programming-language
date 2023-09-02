package main

import (
	"crypto/sha256"
)

func main() {
	a, b := []byte("x"), []byte("X")
	x, y := sha256.Sum256(a), sha256.Sum256(b)
	println(bitMatchCount(x, y))
}

func bitMatchCount(x, y [32]byte) (count int) {
	for i := 0; i < len(x); i++ {
		notmatched := x[i] ^ y[i]
		for notmatched != 0 {
			notmatched &= notmatched - 1
			count++
		}
	}

	return count
}
