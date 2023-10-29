package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
	}

	for _, link := range recVisit(nil, doc) {
		fmt.Println(link)
	}
}

func recVisit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode {
		var tmp string
		tmp = extractHref(n)
		if tmp != "" {
			links = append(links, tmp)
		}

		tmp = extractSrc(n)
		if tmp != "" {
			links = append(links, tmp)
		}
	}

	if n.FirstChild != nil {
		links = recVisit(links, n.FirstChild)
	}

	if n.NextSibling != nil {
		links = recVisit(links, n.NextSibling)
	}

	return links
}

func extractSrc(n *html.Node) string {
	var ret string
	for _, node := range n.Attr {
		if node.Key == "src" && hasLetters(node.Val) {
			ret = fmt.Sprintf("%s: %s", n.Data, strings.TrimSpace(node.Val))
		}
	}

	return ret
}

func extractHref(n *html.Node) string {
	var ret string
	for _, node := range n.Attr {
		if node.Key == "href" && hasLetters(node.Val) {
			ret = fmt.Sprintf("%s: %s", n.Data, strings.TrimSpace(node.Val))
		}
	}

	return ret
}

func hasLetters(s string) bool {
	for _, char := range s {
		if unicode.IsLetter(char) {
			return true
		}
	}

	return false
}
