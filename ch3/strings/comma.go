// Comma delimit large numbers
package main

import (
	"bytes"
	"strings"
)

func main() {
	println(buffercomma("123"))
	println(buffercomma("-123"))
	println(buffercomma("1234.12"))
	println(buffercomma("1.12"))
	println(buffercomma("12.12"))
	println(buffercomma("123.12"))
	println(buffercomma("+123422.12"))
	println(buffercomma("1234.133452"))
	println(buffercomma("-123422.12"))
	println(buffercomma("-12.12"))
}

func recursivecomma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return recursivecomma(s[:n-3]) + "," + s[n-3:]
}

func buffercomma(s string) string {
	var buf bytes.Buffer
	var n int

	if strings.Contains(s, ".") {
		n = len(s[:strings.Index(s, ".")])
	} else {
		n = len(s)
	}

	if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
		n -= 1
		buf.WriteByte(s[0])
		s = s[1:]
	}

	for i, j := range s {
		if j == '.' {
			buf.WriteString(s[i:])
			return buf.String()
		}
		if (n-i)%3 == 0 && i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune(j)
	}
	return buf.String()
}
