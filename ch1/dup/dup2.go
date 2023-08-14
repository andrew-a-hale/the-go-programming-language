// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 10.
//!+

// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	lineLoc := make(map[string][]string)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts, lineLoc)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, lineLoc)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("loc: %s -- %d\t%s\n", strings.Join(lineLoc[line], ", "), n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int, loc map[string][]string) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		txt := input.Text()
		counts[txt]++
		if !contains(loc[txt], f.Name()) {
			loc[txt] = append(loc[txt], f.Name())
		}
	}
	// NOTE: ignoring potential errors from input.Err()
}

func contains[T comparable](s []T, e T) bool {
	for _, t := range s {
		if t == e {
			return true
		}
	}
	return false
}

//!-
