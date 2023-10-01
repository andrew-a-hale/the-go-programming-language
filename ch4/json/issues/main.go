package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch4/github"
)

var (
	monthAgo = time.Now().AddDate(0, -1, 0).Unix()
	yearAgo  = time.Now().AddDate(-1, 0, 0).Unix()
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	printByAgeLessThanAMonthAgo(result)
	printByAgeLessThanAYearAgo(result)
	printByAgeGreaterThanAYearAgo(result)
}

func printByAgeLessThanAMonthAgo(result *github.IssuesSearchResult) {
	fmt.Printf("issues created in the last month:\n")
	issues := 0
	for _, item := range result.Items {
		if item.CreatedAt.Unix() > monthAgo {
			fmt.Printf("%v -- #%-5d %9.9s %.55s\n", item.CreatedAt, item.Number, item.User.Login, item.Title)
			issues++
		}
	}

	fmt.Printf("issues: %d\n", issues)
}

func printByAgeLessThanAYearAgo(result *github.IssuesSearchResult) {
	fmt.Printf("issues created in the last year:\n")
	issues := 0
	for _, item := range result.Items {
		if item.CreatedAt.Unix() > yearAgo {
			fmt.Printf("%v -- #%-5d %9.9s %.55s\n", item.CreatedAt, item.Number, item.User.Login, item.Title)
			issues++
		}
	}

	fmt.Printf("issues: %d\n", issues)
}

func printByAgeGreaterThanAYearAgo(result *github.IssuesSearchResult) {
	fmt.Printf("issues at least a year old:\n")
	issues := 0
	for _, item := range result.Items {
		if item.CreatedAt.Unix() < yearAgo {
			fmt.Printf("%v -- #%-5d %9.9s %.55s\n", item.CreatedAt, item.Number, item.User.Login, item.Title)
			issues++
		}
	}

	fmt.Printf("issues: %d\n", issues)
}
