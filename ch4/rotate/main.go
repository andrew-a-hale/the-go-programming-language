package main

import "fmt"

func main() {
	x := []int{1, 2, 3}
	rotate(x, 2)
	fmt.Printf("%v\n", x)
	rotate_2(x, 2)	
	fmt.Printf("%v\n", x)
}

func rotate(x []int, mv int) []int {
	temp := make([]int, len(x))
	copy(temp, x)
	for i := 0; i < len(x); i++ {
		new_pos := (i + mv) % len(x)
		x[new_pos] = temp[i]
	}
	return x
}

func rotate_2(x []int, mv int) []int {
	mv = mv % len(x)
	temp := make([]int, len(x))
	copy(temp, x[:mv+1])
	copy(x, x[len(x)-mv:])
	copy(x[mv:], temp)
	return x
}
