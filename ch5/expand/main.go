package main

import (
	"fmt"
	"strings"
)

func main() {
	str := " $foo -- $foo "
	fmt.Println(str, Expand(str, strings.ToUpper))
}

func Expand(s string, f func(string) string) string {
	replacement := f("foo")
	return strings.ReplaceAll(s, "$foo", replacement)
}
