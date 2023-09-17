package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	const (
		space   = "space"
		digit   = "digit"
		letter  = "letter"
		punct   = "punct"
		symbol  = "symbol"
		invalid = "invalid"
	)
	typeCounts := make(map[string]int)
	counts := make(map[rune]int)
	var utflen [utf8.UTFMax + 1]int

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune()

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stdout, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			typeCounts[invalid]++
			continue
		}

		switch {
		case unicode.IsSpace(r):
			typeCounts[space]++
		case unicode.IsLetter(r):
			typeCounts[letter]++
		case unicode.IsDigit(r):
			typeCounts[digit]++
		case unicode.IsPunct(r):
			typeCounts[punct]++
		case unicode.IsSymbol(r):
			typeCounts[symbol]++
		}

		counts[r]++
		utflen[n]++
	}

	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}

	fmt.Printf("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}

	fmt.Printf("\nsummary\n")
	for k, v := range typeCounts {
		if v > 0 {
			fmt.Printf("%ss %d\n", k, v)
		}
	}
}
