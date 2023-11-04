package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"unicode"

	"golang.org/x/net/html"
)

const INDENT_SIZE int = 2

var singleLineNode = []string{"title", "meta", "h1", "tt", "a", "code", "br", "p", "i"}
var depth int

func main() {
	if len(os.Args) < 1 {
		log.Fatalf("missing url")
	}

	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to get url: %s, %s\n", url, err)
	}

	doc, err := html.Parse(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("failed to parse html: %s\n", err)
	}

	forEachNode(doc, startElement, endElement)
}

func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) *html.Node {
	if pre != nil {
		ok := pre(n)
		if !ok {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data != "script" && c.Data != "style" {
			forEachNode(c, pre, post)
		}
	}

	if post != nil {
		ok := post(n)
		if !ok {
			return n
		}
	}

	return n
}

func startElement(n *html.Node) bool {
	switch {
	case n.Type == html.CommentNode:
		fmt.Printf("%*s<!-- %s -->\n", depth*INDENT_SIZE, "", n.Data)
		depth++
	case n.Type == html.ElementNode:
		switch {
		case n.FirstChild == nil:
			oneTagNode(n, depth*INDENT_SIZE)
			depth++
		case oneOf(n.Data, singleLineNode):
			oneLineNode(n, depth*INDENT_SIZE)
			depth++
		default:
			handleNodeAttrs(n, depth*INDENT_SIZE)
			depth++
		}
	case n.Type == html.TextNode && hasAlpha(n.Data):
		switch {
		case oneOf(n.Parent.Data, singleLineNode):
			fmt.Printf("%s", strings.TrimSpace(n.Data))
		default:
			fmt.Printf("%*s%s\n", INDENT_SIZE*depth, "", strings.TrimSpace(n.Data))
		}
	}

	return true
}

func endElement(n *html.Node) bool {
	switch {
	case n.Type == html.CommentNode:
		depth--
	case n.Type == html.ElementNode:
		switch {
		case n.FirstChild == nil:
			depth--
		case oneOf(n.Data, singleLineNode):
			depth--
			if oneOf(n.Parent.Data, singleLineNode) {
				fmt.Printf("</%s>", n.Data)
			} else {
				fmt.Printf("</%s>\n", n.Data)
			}
		default:
			depth--
			fmt.Printf("%*s</%s>\n", depth*INDENT_SIZE, "", n.Data)
		}
	}

	return true
}

func handleNodeAttrs(n *html.Node, indent int) {
	if oneOf(n.Parent.Data, singleLineNode) {
		fmt.Printf("<%s", n.Data)
	} else {
		fmt.Printf("%*s<%s", indent, "", n.Data)
	}

	for _, attr := range n.Attr {
		fmt.Printf(" %s='%s'", attr.Key, attr.Val)
	}
	fmt.Println(">")
}

func oneLineNode(n *html.Node, indent int) {
	if oneOf(n.Parent.Data, singleLineNode) {
		fmt.Printf("<%s", n.Data)
	} else {
		fmt.Printf("%*s<%s", indent, "", n.Data)
	}

	for _, attr := range n.Attr {
		switch attr.Key {
		case "src":
			fmt.Printf(" %s=\"%s\"", attr.Key, attr.Val)
		default:
			fmt.Printf(" %s='%s'", attr.Key, attr.Val)
		}
	}

	fmt.Print(">")
}

func oneTagNode(n *html.Node, indent int) {
	if oneOf(n.Parent.Data, singleLineNode) {
		fmt.Printf("<%s", n.Data)
	} else {
		fmt.Printf("%*s<%s", indent, "", n.Data)
	}

	for _, attr := range n.Attr {
		switch attr.Key {
		case "src":
			fmt.Printf(" %s=\"%s\"", attr.Key, attr.Val)
		default:
			fmt.Printf(" %s='%s'", attr.Key, attr.Val)
		}
	}

	if oneOf(n.Parent.Data, singleLineNode) || (n.PrevSibling != nil && oneOf(n.PrevSibling.Data, singleLineNode)) {
		fmt.Print("/>")
	} else {
		fmt.Println("/>")
	}
}

func hasAlpha(s string) bool {
	for _, char := range s {
		if !unicode.IsSpace(char) {
			return true
		}
	}

	return false
}

func oneOf(needle string, haystack []string) bool {
	for _, el := range haystack {
		if el == needle {
			return true
		}
	}

	return false
}

func TestMain(t *testing.T) {
	url := "http://gopl.io"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to get url: %s, %s\n", url, err)
	}

	doc, err := html.Parse(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("failed to parse html: %s\n", err)
	}

	r, w, err := os.Pipe()
	os.Stdout = w
	if err != nil {
		log.Fatal(err)
	}

	forEachNode(doc, startElement, endElement)
	w.Close()
	_, err = html.Parse(r)
	if err != nil {
		t.Error(`forEachNode failed to generate valid html`)
	}
}
