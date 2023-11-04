package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

func main() {
	var url string
	if len(os.Args) > 1 {
		url = os.Args[1]
	}
	words, images, err := CountWordsAndImages(url)
	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	fmt.Printf("words: %d\nimages: %d\n", words, images)
}

func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed get %s: %s", url, err)
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("failed parse %s", err)
		return
	}

	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(node *html.Node) (words, images int) {
	if node.Type == html.ElementNode && node.Data == "img" {
		images++
	}

	if node.Type == html.TextNode {
		for _, val := range strings.Fields(node.Data) {
			if hasLetters(val) {
				words++
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		isFilteredNode := c.Type == html.ElementNode && (c.Data == "script" || c.Data == "style")
		if c.Type == html.TextNode || !isFilteredNode {
			ws, is := countWordsAndImages(c)
			words, images = words+ws, images+is
		}
	}
	return
}

func hasLetters(s string) bool {
	for _, char := range s {
		if unicode.IsLetter(char) {
			return true
		}
	}

	return false
}
