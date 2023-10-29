package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

var nodeMap = make(map[string]int)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
	}

	recVisit(doc)
}

func recVisit(n *html.Node) {
	if n.Data == "script" || n.Data == "style" {
		return
	}

	if n.Type == html.TextNode && hasLetters(n.Data) {
		fmt.Printf("%s\n", strings.TrimSpace(n.Data))
	}

	if n.FirstChild != nil {
		recVisit(n.FirstChild)
	}

	if n.NextSibling != nil {
		recVisit(n.NextSibling)
	}
}

func hasLetters(s string) bool {
	for _, char := range s {
		if unicode.IsLetter(char) {
			return true
		}
	}

	return false
}
