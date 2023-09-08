package main

import "fmt"

func main() {
	strings := []string{"a", "a", "b", "c"}
	fmt.Printf("%v\n", adj_dedup(strings))

	strings = []string{"a", "b", "c"}
	fmt.Printf("%v\n", adj_dedup(strings))

	strings = []string{"a", "b", "b", "c"}
	fmt.Printf("%v\n", adj_dedup(strings))

	strings = []string{"a", "b", "c", "c"}
	fmt.Printf("%v\n", adj_dedup(strings))

	strings = []string{"a"}
	fmt.Printf("%v\n", adj_dedup(strings))
}

func adj_dedup(strings []string) []string {
	out := []string{strings[0]}
	for i := 1; i < len(strings); i++ {
		if strings[i] != strings[i-1] {
			out = append(out, strings[i])
		}
	}
	return out
}
