// Comma delimit large numbers
package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
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
	println(buffercomma(""))
	println(buffercomma("abcd"))
	println(buffercomma("abc1"))
	println(buffercomma("123,123.2"))
}

func recursivecomma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return recursivecomma(s[:n-3]) + "," + s[n-3:]
}

func buffercomma(s string) (string, error) {
	var buf bytes.Buffer
	var size int // number of digits before a decimal point, otherwise len(s)

	// validate
	regex := `^[\d-+\.]*$`
	matched, err := regexp.Match(regex, []byte(s))
	switch {
	case len(s) == 0:
		return "", errors.New(fmt.Sprint("received empty string\n"))
	case err != nil:
		return "", errors.New(fmt.Sprintf("invalid regex used, regex was: %v", regex))
	case !matched:
		return "", errors.New(fmt.Sprintf("string contained invalid characters expected decimal-like strings, got: %v\n", s))
	}

	size = strings.LastIndex(s, ".")
	if size == -1 {
		size = len(s)
	}

	// peel off the +/- and write to buf
	if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
		size -= 1
		buf.WriteByte(s[0])
		s = s[1:]
	}

	for i, j := range s {
		if j == '.' {
			buf.WriteString(s[i:])
			return buf.String(), nil
		}
		if (size-i)%3 == 0 && i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune(j)
	}
	return buf.String(), nil
}
