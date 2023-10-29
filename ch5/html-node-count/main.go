package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

var nodeMap = make(map[string]int)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
	}

	for k, v := range recVisit(nodeMap, doc) {
		fmt.Println(k, v)
	}
}

func recVisit(nodeMap map[string]int, n *html.Node) map[string]int {
	if n.Type == html.ElementNode {
		nodeMap[n.Data]++
	}

	if n.FirstChild != nil {
		recVisit(nodeMap, n.FirstChild)
	}

	if n.NextSibling != nil {
		recVisit(nodeMap, n.NextSibling)
	}
	
	return nodeMap
}
