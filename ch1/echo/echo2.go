// Echo2 prints its command-line arugments.
package main

import (
	"fmt"
	"os"
	"strconv" // exercise 1.2
)

func main() {
	s, sep := "", ""
	for i, arg := range os.Args[1:] {
		s += sep + arg 
		s += sep + strconv.Itoa(i) + arg // exercise 1.2
		sep = "\n" // exercise 1.2
	}
	fmt.Println(s)
}
