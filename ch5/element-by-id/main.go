package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing url")
	}

	url := os.Args[1]
	id := os.Args[2]
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to get url: %s, %s\n", url, err)
	}

	doc, err := html.Parse(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("failed to parse html: %s\n", err)
	}

	forEachNode(doc, id, ElementById)
}

func forEachNode(n *html.Node, id string, pre func(n *html.Node, id string) *html.Node) (found *html.Node) {
	if pre != nil {
		if found := pre(n, id); found != nil {
			return found
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := forEachNode(c, id, pre); found != nil {
			return found
		}
	}

	return found
}

func ElementById(doc *html.Node, id string) *html.Node {
	for _, attrs := range doc.Attr {
		if attrs.Key == "id" && attrs.Val == id {
			return doc
		}
	}

	return nil
}
