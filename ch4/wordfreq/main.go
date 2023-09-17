package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	wordFreqs := make(map[string]int)

	var in *os.File
	var err error
	if len(os.Args) > 1 {
		file := os.Args[1]
		in, err = os.Open(file)
	} else if in == nil || err != nil {
		in = os.Stdin
	}

	input := bufio.NewScanner(in)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		word := input.Text()
		word = strings.ToLower(word)

		// remove punct
		r, _ := utf8.DecodeLastRuneInString(word)
		if unicode.IsPunct(r) {
			word = word[:len(word)-1]
		}

		if len(word) > 0 {
			wordFreqs[word]++
		}
	}

	for k, v := range wordFreqs {
		fmt.Printf("%s: %d\n", k, v)
	}
}
