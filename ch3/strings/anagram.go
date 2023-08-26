// Anagrams
package main

import (
	"sort"
	"strings"
)

func main() {
	// println(isAnagram("poo", "poo"))
	// println(isAnagram("poo", "po"))
	// println(isAnagram("poo", "opo"))
	// println(isAnagram("poo", ""))
	// println(isAnagram("", ""))

	println(recursiveIsAnagram("poo", "poo"))
	println(recursiveIsAnagram("poo", "po"))
	println(recursiveIsAnagram("poo", "opo"))
	println(recursiveIsAnagram("poo", ""))
	println(recursiveIsAnagram("", ""))
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
	for i, r := range asplit {
		if r != bsplit[i] {
			return false
		}
	}

	return true
}

func recursiveIsAnagram(a, b string) bool {
	if a == b {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	return recursiveAnagramHelper(a, b)
}

func recursiveAnagramHelper(a, b string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}

	if !strings.ContainsRune(b, rune(a[0])) {
		return false
	}

	return recursiveAnagramHelper(a[1:], removeRune(b, rune(a[0])))
}

func removeRune(str string, remove rune) string {
	for i, c := range str {
		if c == remove {
			return str[0:i] + str[i+1:]
		}
	}

	return str
}