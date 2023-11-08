package main

import (
	"fmt"
	"log"
)

var prereqs = map[string][]string{
	"algorithms":            {"data structures"},
	"calculus":              {"linear algebra"},
	"linear algebra":        {"calculus 2"},
	"calculus 2":            {"calculus"},
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
		val, err := getByValue(ordering, i)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d:\t%s\n", i, val)
	}
}

func topoSort(m map[string][]string) map[string]int {
	var visitAll func(item string)
	order := make(map[string]int)
	var pos int
	var head string

	visitAll = func(item string) {
		for _, dep := range m[item] {
			if head == dep {
				fmt.Printf("cycle detected for key: %s\n", head)
				continue
			}
			if !contains(order, dep) {
				visitAll(dep)
			}
		}
		if !contains(order, item) {
			pos++
			order[item] = pos
		}
	}

	for key := range m {
		head = key
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

func contains(haystack map[string]int, needle string) bool {
	for cand := range haystack {
		if cand == needle {
			return true
		}
	}

	return false
}
