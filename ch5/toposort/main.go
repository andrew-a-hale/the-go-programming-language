package main

import (
	"fmt"
	"log"
)

var prereqs = map[string][]string{
	"algorithms":            {"data structures"},
	"calculus":              {"linear algebra"},
	"compilers":             {"data structures", "formal languages", "computer organization"},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operatings systems":    {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	ordering := topoSort(prereqs)
	for i := 1; i < len(ordering); i++ {
		val, err  := getByValue(ordering, i)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d:\t%s\n", i, val)
	}
}

func topoSort(m map[string][]string) map[string]int {
	var visitAll func(item string)
	order := make(map[string]int)
	seen := make(map[string]bool)
	var size int

	visitAll = func(item string) {
		if !seen[item] {
			seen[item] = true
			for _, dep := range m[item] {
				visitAll(dep)
			}
			size++
			order[item] = size
		}
	}

	for key := range m {
		visitAll(key)
	}

	return order
}

func getByValue(m map[string]int, i int) (string, error) {
	for k, v := range m {
		if v == i {
			return k, nil
		}
	}

	return "", fmt.Errorf("value %d not found in map", i)
}