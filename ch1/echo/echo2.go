// Echo2 prints its command-line arugments.
package main

import (
	"fmt"
	"os"
	"strconv" // exercise 1.2
	"time"
)

func main() {
	start := time.Now()

	s, sep := "", "\n"

	s += os.Args[0] + sep
	for i, arg := range os.Args[1:] {
		s += strconv.Itoa(i) + " " + arg
		s += sep
	}
	fmt.Println(s)
	fmt.Println(time.Since(start).Nanoseconds())
}
