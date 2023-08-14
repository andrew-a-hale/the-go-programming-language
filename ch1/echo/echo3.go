// Echo2 prints its command-line arugments.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println(strings.Join(os.Args[0:], "\n")) // exercise 1
	fmt.Println(time.Since(start).Nanoseconds())
}
