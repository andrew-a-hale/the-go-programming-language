package main

import (
	"fmt"
	"log"
)

var prereqs = map[string][]string{
	"algorithms":            {"data structures"},
	"calculus":              {"linear algebra"},
	"calculus 2":            {"calculus"},
	"compilers":             {"data structures", "formal languages", "computer organization"},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

var files = map[string][]string{
	"d1": {"d4", "f2", "f3"},
	"d4": {"f10", "f11"},
	"d5": {"f1", "f5"},
	"d2": {"f4", "d5", "f6"},
	"d3": {"f7", "f8", "f9"},
}

func main() {
	ordering := topoSortDepth(files)
	for i := 1; i < len(ordering); i++ {
		val, err := getByValue(ordering, i)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d:\t%s\n", i, val)
	}

	topoSortBreadth(files)
}

func topoSortBreadth(m map[string][]string) {
	seen := make(map[string]bool)

	// search prereqs
	var f func(item string) []string = func(item string) []string {
		if !seen[item] {
			fmt.Println(item)
		}
		for _, el := range m[item] {
			if !seen[el] {
				seen[el] = true
				fmt.Println(el)
			}
		}
		return m[item]
	}

	// invert map
	inverted := make(map[string][]string)
	var worklist []string
	for k := range m {
		for _, v := range m[k] {
			inverted[v] = append(inverted[v], k)
		}
		worklist = append(worklist, k)
	}

	// breadthfirst
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			worklist = append(worklist, f(item)...)
		}
	}
}

func topoSortDepth(m map[string][]string) map[string]int {
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
