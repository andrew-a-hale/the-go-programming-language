package main

import "fmt"

func main() {
	x := [5]int{1, 2, 3, 4, 5}
	reverse(&x)
	fmt.Printf("%v\n", x)

	y := []byte("a b c")
	fmt.Printf("%s\n", reverseBytes(y))
}

func reverse(x *[5]int) {
	for i := 0; i < len(x)/2; i++ {
		x[i], x[len(x)-1-i] = x[len(x)-1-i], x[i]
	}
}

func reverseBytes(x []byte) []byte {
	for i := 0; i < len(x)/2; i++ {
		x[i], x[len(x)-1-i] = x[len(x)-1-i], x[i]
	}
	return x
}
