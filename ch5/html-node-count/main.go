package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
	}

	var total int
	for k, v := range visit(doc) {
		fmt.Println(k, v)
		total = total + v
	}
	fmt.Printf("total: %d\n", total)
}

func recVisit(n *html.Node) map[string]int {
	var nodeMap = make(map[string]int)
	if n.Type == html.ElementNode {
		nodeMap[n.Data]++
	}

	if n.FirstChild != nil {
		nodes := recVisit(n.FirstChild)
		for k, v := range nodes {
			nodeMap[k] = nodeMap[k] + v
		}
	}

	if n.NextSibling != nil {
		nodes := recVisit(n.NextSibling)
		for k, v := range nodes {
			nodeMap[k] = nodeMap[k] + v
		}
	}
	
	return nodeMap
}

func visit(n *html.Node) map[string]int {
	var nodeMap = make(map[string]int)
	if n.Type == html.ElementNode {
		nodeMap[n.Data]++
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes := visit(c)
		for k, v := range nodes {
			nodeMap[k] = nodeMap[k] + v
		}
	}
	
	return nodeMap
}