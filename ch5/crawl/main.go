package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	var url string
	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	urls := crawl(url)
	for _, url := range urls {
		fmt.Printf("%s\n", url)
	}
}

func extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as html: %v", url, err)
	}

	domain := strings.Split(url, "/")[2]
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}

				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}

				linkString, _ := strings.CutSuffix(link.String(), "/")
				segments := strings.Split(linkString, "/")[2:]
				if segments[0] == domain && len(segments) > 1 {
					// make dir for domain
					dir := strings.Join(segments[:len(segments)-1], "/")
					err := os.MkdirAll(dir, 0755)
					if err != nil {
						fmt.Printf("failed to create file for %s: %s\n", dir, err)
						continue
					}

					// get html for url with name as last segment of url
					resp, err := http.Get(link.String())
					if err != nil || resp.StatusCode != http.StatusOK {
						continue
					}

					content, _ := io.ReadAll(resp.Body)
					fileExt, _ := mime.ExtensionsByType(resp.Header.Get("Content-Type"))
					err = os.WriteFile(strings.Join(segments, "/")+fileExt[0], content, 0755)
					if err != nil {
						fmt.Printf("failed to write file %s: %s\n", segments, err)
					}

					resp.Body.Close()
				}

				links = append(links, link.String())
			}
		}
	}

	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(doc *html.Node, pre func(*html.Node), post func(*html.Node)) {
	if pre != nil {
		pre(doc)
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(doc)
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}
