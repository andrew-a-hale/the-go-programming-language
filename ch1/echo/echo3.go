// Echo2 prints its command-line arugments.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args[0:], " ")) // exercise 1
}
