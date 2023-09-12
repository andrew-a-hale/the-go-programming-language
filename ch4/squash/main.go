package main

import (
	"fmt"
	"unicode"
)

func main() {
	s := []byte("a b  c     d e")

	fmt.Printf("%s\n", s)
	fmt.Printf("%s\n", squash(s))
}

func squash(s []byte) []byte {
	length := len(s)

	for i := range s {
		wsCount := 0
		
		if unicode.IsSpace(rune(s[i])) {
			// lookahead
			for _, ws := range s[i:] {
				if !unicode.IsSpace(rune(ws)) {
					break
				}
				wsCount++
			}
			
			copy(s[i+1:], s[i+wsCount:])
			// adjust returned length
			length -= wsCount - 1
		}
	}

	return s[:length]
}
