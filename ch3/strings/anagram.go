// Anagrams
package main

import (
	"sort"
	"strings"
)

func main() {
	a := "poo"
	b := "poo"
	println(isAnagram(a, b))
}

func isAnagram(a string, b string) bool {
	if a == b {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	asplit, bsplit := strings.Split(a, ""), strings.Split(b, "")
	sort.Strings(asplit)
	sort.Strings(bsplit)
	if strings.Join(asplit, "") == strings.Join(bsplit, "") {
		return true
	}
	return false
}
